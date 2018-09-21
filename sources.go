package cachier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const TypeJson = "json"

type Source struct {
	Key                 string
	RefreshEverySeconds int
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

var Sources []*Source

func AddSource(source *Source) {
	Sources = append(Sources, source)
}

func RemoveSource(sourceKey string) {
	indexToRemove := -1
	for index, source := range Sources {
		if source.Key == sourceKey {
			indexToRemove = index
		}
	}

	if indexToRemove != -1 {
		Sources = append(Sources[:indexToRemove], Sources[indexToRemove+1])
	}
}

func Purge() {
	Sources = []*Source{}
}

func Fetch(source Source) (SourceResponse, error) {
	switch source.Type {
	case TypeJson:
		return fetchJson(source)
	default:
		return SourceResponse{}, fmt.Errorf("%s type is not supported yet", source.Type)
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
