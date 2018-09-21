package cachier

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type cacheSourceMap struct {
	data   map[string]string
	source *Source
	mtx    sync.RWMutex
}

var cacheMap = map[string]*cacheSourceMap{}
var sources map[string]*Source
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
		Set(source, record.Key, record.Value)
	}

	mtx.Lock()
	source.RefreshedAt = time.Now()
	mtx.Unlock()

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

func Set(source *Source, dataKey string, dataValue string) {
	mtx.RLock()
	cachedSource, ok := cacheMap[source.Key]
	mtx.RUnlock()

	if !ok {
		mtx.Lock()
		cachedSource = &cacheSourceMap{
			map[string]string{},
			source,
			sync.RWMutex{},
		}
		cacheMap[source.Key] = cachedSource
		mtx.Unlock()
	}

	cachedSource.mtx.Lock()
	cachedSource.data[dataKey] = dataValue
	cachedSource.mtx.Unlock()
}

func Get(sourceKey string, dataKey string) (string, error) {
	mtx.RLock()
	source, ok := cacheMap[sourceKey]
	mtx.RUnlock()

	if !ok {
		return "", fmt.Errorf("source does not exist: %s", sourceKey)
	}

	source.mtx.RLock()
	value, ok := source.data[dataKey]
	source.mtx.RUnlock()

	if !ok {
		return "", fmt.Errorf("key \"%s\" does not exist for source \"%s\"", dataKey, sourceKey)
	}

	return value, nil
}
