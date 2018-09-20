package cachier

import "testing"

func TestGetActiveSources(t *testing.T) {
	expectedSources := map[string]string{}

	expectedSources["kiwi"] = "kiwi"
	expectedSources["goavio"] = "goavio"

	sources := GetActiveSources()

	for _, source := range sources {
		_, ok := expectedSources[source.Key]

		if !ok {
			t.Errorf("Unexpected source %s", source.Key)
		}
	}
}

func TestKiwiFetch(t *testing.T) {
	expectedDataKeys := map[string]string{}

	expectedDataKeys["awesomekey"] = "awesomekey"
	expectedDataKeys["awesomekey1"] = "awesomekey1"
	expectedDataKeys["awesomekey2"] = "awesomekey2"

	sourceResponse, err := Fetch(Source{Key: SourceKiwi, Type: TypeJson})

	if err != nil {
		t.Errorf("Cannot fetch data from source")
	}

	for _, record := range sourceResponse.Data {
		_, ok := expectedDataKeys[record.Key]

		if !ok {
			t.Errorf("Unexpected data key %s", record.Key)
		}
	}
}
