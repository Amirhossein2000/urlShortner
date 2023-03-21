package ports

type Storage interface {
	SetURLMapping(shortURL string, longURL string) error
	GetLongURL(shortURL string) (string, error)
	GetShortURL(longURL string) (string, error)
	IncrementStats(url string)
	GetLastThreeDaysStats(url string) map[string]int
}

type Encoder interface {
	GenerateShortURL() string
}
