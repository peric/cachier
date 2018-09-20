package cachier

import (
	"sync"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	Set("real", "real", "I'm real!")
}

func TestConcurrentSet(t *testing.T) {
	wg := sync.WaitGroup{}

	for i := 1; i < 50; i++ {
		wg.Add(1)

		go func() {
			Set("real", "real", "I'm real!")
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestSourceDoesNotExist(t *testing.T) {
	_, err := Get("unreal", "unreal")

	if err == nil {
		t.Errorf("Source 'unreal' should not exist")
	}
}

func TestDataKeyDoesNotExist(t *testing.T) {
	_, err := Get("real", "unreal")

	if err == nil {
		t.Errorf("Key 'unreal' should not exist under source 'real'")
	}
}

func TestGet(t *testing.T) {
	Set("real", "real", "I'm real!")

	_, err := Get("real", "real")

	if err != nil {
		t.Errorf("Key 'real' should exist under the source 'real'")
	}
}

func TestGetAndSleep(t *testing.T) {
	Set("real", "real", "I'm real!")

	_, err := Get("real", "real")

	if err != nil {
		t.Errorf("Key 'real' should exist under the source 'real'")
	}

	StartRefreshingSources(1 * time.Second)

	time.Sleep(5 * time.Second)

	StopRefreshingSources()
}
