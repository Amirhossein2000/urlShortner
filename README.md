## Run

```bash
make run
```

## Endpoints:

| URL               | Method(s) | Description          |
|-------------------|-----------|----------------------|
| /shorten          | POST      | Create a shortURL    |
| /{shortURL}       | Any       | Redirect             |
| /{shortURL}/stats | GET       | Get visit statistics |

## API Usage

For api usage import [this](docs/urlShortner.postman_collection.json) postman_collection. It includes all of the
requests with a
few examples.

## Tests

Unit-Tests exist for all the endpoints you can
check [this](internal/handlers/handlers_test.go)
file to see all.

## URL Generator

Instead of random string we have a Base26 Incremental solution to make the URL as short as possible.