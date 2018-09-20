package cachier

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type SourceMap struct {
	Data map[string]string
	m    sync.RWMutex
}

var sources []*Source
var cacheMap = map[string]*SourceMap{}
var mtx = sync.RWMutex{}
var ticker = time.Ticker{}

func init() {
	sources = GetActiveSources()

	wg := sync.WaitGroup{}
	for _, source := range sources {
		wg.Add(1)
		go buildForSource(source, &wg)
	}

	wg.Wait()

	StartRefreshingSources(time.Minute * 1)
}

func buildForSource(source *Source, wg *sync.WaitGroup) {
	sourceResponse, err := Fetch(*source)

	if err != nil {
		log.Fatalf("Cannot fetch data from %s", source.Key)
	}

	for _, record := range sourceResponse.Data {
		Set(source.Key, record.Key, record.Value)
	}

	source.RefreshedAt = time.Now()

	wg.Done()
}

func StartRefreshingSources(duration time.Duration) {
	ticker := time.NewTicker(duration)

	go func() {
		for range ticker.C {
			wg := sync.WaitGroup{}
			for _, source := range sources {
				duration := time.Since(source.RefreshedAt)

				if int(duration.Seconds()) >= source.RefreshEverySeconds {
					wg.Add(1)
					go buildForSource(source, &wg)
					log.Printf("Updating source '%s'", source.Key)
				}
			}

			wg.Wait()
		}
	}()
}

func StopRefreshingSources() {
	ticker.Stop()
}

func Set(sourceKey string, dataKey string, dataValue string) {
	mtx.RLock()
	source, ok := cacheMap[sourceKey]
	mtx.RUnlock()

	if !ok {
		mtx.Lock()
		source = &SourceMap{
			map[string]string{},
			sync.RWMutex{},
		}
		cacheMap[sourceKey] = source
		mtx.Unlock()
	}

	source.m.Lock()
	source.Data[dataKey] = dataValue
	source.m.Unlock()
}

func Get(sourceKey string, dataKey string) (string, error) {
	mtx.RLock()
	source, ok := cacheMap[sourceKey]
	mtx.RUnlock()

	if !ok {
		return "", fmt.Errorf("source '%s' does not exist", sourceKey)
	}

	source.m.RLock()
	value, ok := source.Data[dataKey]
	source.m.RUnlock()

	if !ok {
		return "", fmt.Errorf("value with key '%s' does not exist for source '%s'", dataKey, sourceKey)
	}

	return value, nil
}
