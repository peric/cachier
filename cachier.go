package cachier

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type sourceMap struct {
	Data map[string]string
	M    sync.RWMutex
}

var sources []*Source
var cacheMap = map[string]*sourceMap{}
var mtx = sync.RWMutex{}
var ticker = time.Ticker{}

func Init() {
	sources = Sources

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
		log.Fatalf("Cannot fetch data from source: %s", source.Key)
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
					log.Printf("Refreshing source: %s", source.Key)
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
		source = &sourceMap{
			map[string]string{},
			sync.RWMutex{},
		}
		cacheMap[sourceKey] = source
		mtx.Unlock()
	}

	source.M.Lock()
	source.Data[dataKey] = dataValue
	source.M.Unlock()
}

func Get(sourceKey string, dataKey string) (string, error) {
	mtx.RLock()
	source, ok := cacheMap[sourceKey]
	mtx.RUnlock()

	if !ok {
		return "", fmt.Errorf("source does not exist: %s", sourceKey)
	}

	source.M.RLock()
	value, ok := source.Data[dataKey]
	source.M.RUnlock()

	if !ok {
		return "", fmt.Errorf("key \"%s\" does not exist for source \"%s\"", dataKey, sourceKey)
	}

	return value, nil
}
