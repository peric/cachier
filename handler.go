package cachier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const SourceKiwi = "kiwi"
const SourceGoAvio = "goavio"

const TypeJson = "json"

type Source struct {
	Key                 string
	RefreshEveryMinutes int
	RefreshedAt         time.Time
	Type                string
	Url                 string
	m                   sync.Mutex
}

type SourceResponse struct {
	Data []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"data"`
}

func GetActiveSources() []*Source {
	var sources []*Source

	sources = append(sources, &Source{Key: SourceKiwi, RefreshEveryMinutes: 60, Type: TypeJson})
	sources = append(sources, &Source{Key: SourceGoAvio, RefreshEveryMinutes: 120, Type: TypeJson})

	return sources
}

func Fetch(source Source) (SourceResponse, error) {
	// for now, we will read only from file. in later implementation, we can fetch directly from an API, db etc
	switch source.Type {
	case TypeJson:
		return fetchJson(source)
	default:
		return SourceResponse{}, fmt.Errorf("%s is not implemented yet", source.Type)
	}
}

func fetchJson(source Source) (SourceResponse, error) {
	fileName := fmt.Sprintf("data/%s.json", source.Key)

	jsonFile, err := os.Open(fileName)

	if err != nil {
		return SourceResponse{}, fmt.Errorf("cannot read file: %s", fileName)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var sourceResponse SourceResponse

	json.Unmarshal(byteValue, &sourceResponse)

	return sourceResponse, nil
}
