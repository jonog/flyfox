package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"strings"
)

const (
	LEX_ZSET_KEY        = "ac"
	STORAGE_HASH_PREFIX = "ac:d:"
	QUERY_SET_PREFIX    = "ac:q:"
)

type Item struct {
	Id      string                 `json:"id"`
	Display string                 `json:"display"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Score   int                    `json:"score"`

	// TODO
	// sortable metrics in data - ability to specify indexing when loading
	// categories - store ids in another set - then union
}

func (i *Item) Save() (err error) {

	conn := pool.Get()
	defer conn.Close()

	// TODO - add stopword filtering
	terms := strings.Split(i.Display, " ")

	conn.Send("MULTI")

	for _, val := range terms {
		conn.Send("ZADD", LEX_ZSET_KEY, 0, strings.ToLower(val)+":"+i.Id)
	}

	// json stored within redis contains id, display, score
	// this is to avoid decoding/re-encoding the json received from redis within the Go app
	i.Data["id"] = i.Id
	i.Data["display"] = i.Display
	i.Data["score"] = i.Score
	dataJSON, _ := json.Marshal(i.Data)
	conn.Send("HMSET", STORAGE_HASH_PREFIX+i.Id, "id", i.Id, "score", i.Score, "json", dataJSON)

	_, err = conn.Do("EXEC")

	return
}

func zsetNameToId(input string) string {
	// TODO: fix - impl is fragile to zset naming prefix
	id := strings.Split(input, ":")[1]
	return id
}

func cacheIdsInSet(conn redis.Conn, query string) error {

	// TODO - skip this match via lex, if the set exists
	matchedIds, err := redis.Strings(conn.Do("ZRANGEBYLEX", LEX_ZSET_KEY, "["+query+"\x00", "["+query+"\xff"))
	if err != nil {
		return err
	}

	// TODO - add a TTL to the set
	idsAsArgs := []string{QUERY_SET_PREFIX + query}
	for _, id := range matchedIds {
		idsAsArgs = append(idsAsArgs, zsetNameToId(id))
	}
	conn.Do("SADD", stringArrayToInterfaceArray(idsAsArgs)...)

	return nil
}

func Search(query string, limit int) ([]byte, error) {

	conn := pool.Get()
	defer conn.Close()

	err := cacheIdsInSet(conn, query)
	if err != nil {
		return nil, err
	}

	sortResult, err := redis.Strings(conn.Do("SORT", QUERY_SET_PREFIX+query, "BY", STORAGE_HASH_PREFIX+"*->score", "LIMIT", 0, limit, "GET", STORAGE_HASH_PREFIX+"*->json", "DESC"))
	if err != nil {
		return nil, err
	}

	jsonReponse := "[" + strings.Join(sortResult, ",") + "]"

	return []byte(jsonReponse), err
}
