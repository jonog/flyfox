package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"strings"
)

const (
	ZSET_PREFIX         = "ac:z:"
	STORAGE_HASH_PREFIX = "ac:d:"
	QUERY_SET_PREFIX    = "ac:q:"
	MAX_LIMIT           = 10
)

type Item struct {
	Id     string                 `json:"id"`
	Phrase string                 `json:"phrase"`
	Data   map[string]interface{} `json:"data,omitempty"`
	Score  int                    `json:"score"`
	Type   string                 `json:"type"`
}

func (i *Item) Save() (err error) {

	conn := pool.Get()
	defer conn.Close()

	// TODO - add stopword filtering
	terms := strings.Split(i.Phrase, " ")

	conn.Send("MULTI")

	for _, val := range terms {

		// for type tweet, store the phrase 'Random' in the following (sorted) sets
		// tweet:r
		// tweet:ra
		// tweet:ran
		// tweet:rand
		// tweet:rando
		// tweet:random

		stringForIndexing := strings.ToLower(val)

		for idx, _ := range stringForIndexing {
			zsetKey := ZSET_PREFIX + i.Type + ":" + stringForIndexing[0:idx+1]
			conn.Send("ZADD", zsetKey, i.Score, i.Id)
			conn.Send("ZREMRANGEBYRANK", zsetKey, 0, -MAX_LIMIT-1)
		}

	}

	// store id and score within the json
	// this is to avoid decoding/re-encoding the json received from redis within the Go app
	i.Data["id"] = i.Id
	i.Data["score"] = i.Score
	dataJSON, _ := json.Marshal(i.Data)
	conn.Send("HMSET", STORAGE_HASH_PREFIX+i.Id, "id", i.Id, "json", dataJSON)

	_, err = conn.Do("EXEC")

	return
}

func queryByType(conn redis.Conn, query string, limit int, t string) ([]string, error) {

	sortResult, err := redis.Strings(conn.Do("SORT", ZSET_PREFIX+t+":"+query, "BY", "score", "LIMIT", 0, limit, "GET", STORAGE_HASH_PREFIX+"*->json", "DESC"))
	if err != nil {
		return nil, err
	}
	return sortResult, err
}

func Search(types []string, query string, limit int) ([]byte, error) {

	conn := pool.Get()
	defer conn.Close()
	var err error

	lowerCaseQuery := strings.ToLower(query)

	var allResults []string
	for _, t := range types {
		typeResults, err := queryByType(conn, lowerCaseQuery, limit, t)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, resultsToJSONArray(typeResults))
	}
	jsonResponse := typesAndResultsByTypeToJSON(types, allResults)
	return []byte(jsonResponse), err
}
