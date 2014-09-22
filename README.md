
## flyfox
* autocomplete 'microservice' using Go/Redis
* currently results based on 'score', with further querying flexibility soon

## Get Started
**install dependencies**
```
go get github.com/gorilla/mux
go get github.com/garyburd/redigo/redis
go get github.com/satori/go.uuid
```

**build**
```go build```

**load data**
```./flyfox load_data ./data.json```

**run web server**
```./flyfox web```

**query**
```curl --request GET -g 'http://localhost:3001/search?term=ba&limit=10&types[]=sample&types[]=sample_2'```

## Planned Improvements
* ~~better error handling (return 500, don't panic)~~
* add min query length requirement
* ~~groups/types [ability to query on a single type]~~
* lua scripting to replace sequential redis calls
* sorting on specified fields
* handle no ids [internal id generation]
* stop words
* redis config
