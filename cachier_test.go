package cachier

import (
	"sync"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	realSource := Source{Key: "real", RefreshEverySeconds: 60, Type: TypeJson}

	Set(&realSource, "real", "I'M real!")
}

func TestConcurrentSet(t *testing.T) {
	realSource := Source{Key: "real", RefreshEverySeconds: 60, Type: TypeJson}

	wg := sync.WaitGroup{}

	for i := 1; i < 50; i++ {
		wg.Add(1)

		go func() {
			Set(&realSource, "real", "I'm real!")
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestSourceDoesNotExist(t *testing.T) {
	_, err := Get("unreal", "unreal")

	if err == nil {
		t.Error("Source unreal should not exist")
	}
}

func TestDataKeyDoesNotExist(t *testing.T) {
	_, err := Get("real", "unreal")

	if err == nil {
		t.Error("Key unreal should not exist under source real'")
	}
}

func TestGet(t *testing.T) {
	realSource := Source{Key: "real", RefreshEverySeconds: 60, Type: TypeJson}

	Set(&realSource, "real", "I'm real!")

	_, err := Get("real", "real")

	if err != nil {
		t.Error("Key real should exist under the source real")
	}
}

func TestEndToEnd(t *testing.T) {
	Purge()

	goavioRefreshEverySeconds := 3

	kiwiSource := Source{Key: "kiwi", RefreshEverySeconds: 60, Type: TypeJson}
	goavioSource := Source{Key: "goavio", RefreshEverySeconds: goavioRefreshEverySeconds, Type: TypeJson}

	AddSource(&kiwiSource)
	AddSource(&goavioSource)

	Init()

	assertEqual(t, 2, len(Sources), "")

	firstTimeRefreshedAt := goavioSource.RefreshedAt

	StartRefreshingSources(1 * time.Second)

	time.Sleep(5 * time.Second)

	StopRefreshingSources()

	mtx.Lock()
	lastTimeRefreshedAt := goavioSource.RefreshedAt
	mtx.Unlock()

	diff := lastTimeRefreshedAt.Sub(firstTimeRefreshedAt)

	if int(diff.Seconds()) < goavioRefreshEverySeconds {
		t.Errorf("goavio source should be refreshed every %d seconds", goavioRefreshEverySeconds)
	}
}
