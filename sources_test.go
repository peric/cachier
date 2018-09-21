package cachier

import (
	"fmt"
	"testing"
)

func TestPurgeSources(t *testing.T) {
	Purge()

	assertEqual(t, 0, len(Sources), "")
}

func TestEmptySources(t *testing.T) {
	Purge()

	assertEqual(t, 0, len(Sources), "")
}

func TestAddSource(t *testing.T) {
	Purge()

	AddSource(&Source{Key: "kiwi", RefreshEverySeconds: 60, Type: TypeJson})

	assertEqual(t, 1, len(Sources), "")
}

func TestAddSources(t *testing.T) {
	Purge()

	AddSource(&Source{Key: "kiwi", RefreshEverySeconds: 60, Type: TypeJson})
	AddSource(&Source{Key: "goavio", RefreshEverySeconds: 3, Type: TypeJson})

	assertEqual(t, 2, len(Sources), "")
}

func TestAddAndCheckSources(t *testing.T) {
	Purge()

	expectedSources := map[string]string{}

	AddSource(&Source{Key: "kiwi", RefreshEverySeconds: 60, Type: TypeJson})
	AddSource(&Source{Key: "goavio", RefreshEverySeconds: 3, Type: TypeJson})

	expectedSources["kiwi"] = "kiwi"
	expectedSources["goavio"] = "goavio"

	assertEqual(t, 2, len(Sources), "")

	for _, source := range Sources {
		_, ok := expectedSources[source.Key]

		if !ok {
			t.Errorf("Unexpected source %s", source.Key)
		}
	}
}

func TestRemoveSource(t *testing.T) {
	Purge()

	AddSource(&Source{Key: "kiwi", RefreshEverySeconds: 60, Type: TypeJson})
	AddSource(&Source{Key: "goavio", RefreshEverySeconds: 3, Type: TypeJson})

	assertEqual(t, 2, len(Sources), "")

	RemoveSource("kiwi")

	assertEqual(t, 1, len(Sources), "")
}

func TestFetch(t *testing.T) {
	sourceKey := "kiwi"

	expectedDataKeys := map[string]string{}

	expectedDataKeys["awesomekey"] = "awesomekey"
	expectedDataKeys["awesomekey1"] = "awesomekey1"
	expectedDataKeys["awesomekey2"] = "awesomekey2"

	sourceResponse, err := Fetch(Source{Key: sourceKey, Type: TypeJson})

	if err != nil {
		t.Errorf("Cannot fetch data from source %s", sourceKey)
	}

	for _, record := range sourceResponse.Data {
		_, ok := expectedDataKeys[record.Key]

		if !ok {
			t.Errorf("Unexpected data key %s", record.Key)
		}
	}
}

func assertEqual(t *testing.T, expected int, actual int, message string) {
	if expected == actual {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("Expected: %d, Actual: %d", expected, actual)
	}
	t.Error(message)
}
