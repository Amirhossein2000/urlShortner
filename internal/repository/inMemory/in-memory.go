package inMemory

import (
	"errors"
	"sync"
	"time"
	"urlShortner/internal/consts"
	"urlShortner/internal/help"
)

type InMemory struct {
	shortToLongMap map[string]string
	longToShortMap map[string]string
	stats          map[string]map[string]int
	mutex          sync.Mutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		shortToLongMap: make(map[string]string),
		longToShortMap: make(map[string]string),
		stats:          make(map[string]map[string]int),
		mutex:          sync.Mutex{},
	}
}

func (in *InMemory) SetURLMapping(shortURL string, longURL string) error {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	if _, ok := in.longToShortMap[longURL]; ok {
		return errors.New("Shortened URL already exists")
	}

	in.shortToLongMap[shortURL] = longURL
	in.longToShortMap[longURL] = shortURL

	return nil
}

func (in *InMemory) GetLongURL(shortURL string) (string, error) {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	if longURL, ok := in.shortToLongMap[shortURL]; ok {
		return longURL, nil
	}

	return "", errors.New("Short URL not found")
}

func (in *InMemory) GetShortURL(longURL string) (string, error) {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	if shortURL, ok := in.longToShortMap[longURL]; ok {
		return shortURL, nil
	}

	return "", errors.New("Long URL not found")
}

func (in *InMemory) IncrementStats(url string) {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	if in.stats[url] == nil {
		in.stats[url] = make(map[string]int)
	}
	in.stats[url][time.Now().Format(consts.TimeLayout)]++
}

func (in *InMemory) GetLastThreeDaysStats(url string) map[string]int {
	output := make(map[string]int)
	last3Days := help.LastNDays(3)
	for _, day := range last3Days {
		output[day] = in.stats[url][day]
	}
	return output
}
