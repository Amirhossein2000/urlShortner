package urlGen

type URLGenerator struct {
	counter uint64
}

func NewURLGenerator() *URLGenerator {
	return &URLGenerator{
		counter: 0,
	}
}

func (gen *URLGenerator) GenerateShortURL() string {
	gen.counter++
	shortURL := numberToBase26Incremental(gen.counter)
	if shortURL == "shorten" {
		shortURL = numberToBase26Incremental(gen.counter)
	}
	return shortURL
}

func numberToBase26Incremental(n uint64) string {
	base := uint64(26)
	var result string

	for n > 0 {
		n--
		remainder := n % base
		result = string(remainder+'a') + result
		n = n / base
	}

	return result
}
