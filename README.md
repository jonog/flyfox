
## flyfox
* autocomplete 'microservice' using Go/Redis
* results categorised by 'type' and ordered by 'score'

## Get Started
**install dependencies**
```
go get github.com/gorilla/mux
go get github.com/garyburd/redigo/redis
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
* improve memory efficiency
* sorting on specified fields
* handle no ids [internal id generation]
* stop words
* redis config
