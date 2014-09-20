
## flyfox
* autocomplete 'microservice' using Go/Redis
* currently results based on 'score', with further querying flexibility soon

## Get Started
```go build```

## load data
```./flyfox load_data ./data.json```

## run web app
```./flyfox web```

## Planned Improvements
* better error handling (return 500, don't panic)
* add min query length requirement
* groups/types [ability to query on a single type]
* lua scripting to replace sequential redis calls
* sorting on specified fields
* handle no ids [internal id generation]
* stop words
* redis config
