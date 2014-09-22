package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"strings"
)

const (
	LEX_ZSET_KEY           = "ac"
	TYPE_SCORE_ZSET_PREFIX = "ac:s:"
	STORAGE_HASH_PREFIX    = "ac:d:"
	QUERY_SET_PREFIX       = "ac:q:"
)

type Item struct {
	Id     string                 `json:"id"`
	Phrase string                 `json:"phrase"`
	Data   map[string]interface{} `json:"data,omitempty"`
	Score  int                    `json:"score"`
	Type   string                 `json:"type"`

	// TODO
	// sortable metrics in data - ability to specify indexing when loading
	// categories - store ids in another set - then union
}

func (i *Item) Save() (err error) {

	conn := pool.Get()
	defer conn.Close()

	// TODO - add stopword filtering
	terms := strings.Split(i.Phrase, " ")

	conn.Send("MULTI")

	for _, val := range terms {
		conn.Send("ZADD", LEX_ZSET_KEY, 0, strings.ToLower(val)+":"+i.Id)
	}

	conn.Send("ZADD", TYPE_SCORE_ZSET_PREFIX+i.Type, i.Score, i.Id)

	// store id and score within the json
	// this is to avoid decoding/re-encoding the json received from redis within the Go app
	i.Data["id"] = i.Id
	i.Data["score"] = i.Score
	dataJSON, _ := json.Marshal(i.Data)
	conn.Send("HMSET", STORAGE_HASH_PREFIX+i.Id, "id", i.Id, "json", dataJSON)

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

func intersectAndSort(conn redis.Conn, query string, limit int, t string) ([]string, error) {

	tmpId := uuid.NewV4().String()
	// expensive to zinterstore everything - perhaps could apply limit first to the type score zset
	conn.Do("ZINTERSTORE", tmpId, 2, QUERY_SET_PREFIX+query, TYPE_SCORE_ZSET_PREFIX+t)
	sortResult, err := redis.Strings(conn.Do("SORT", tmpId, "BY", "score", "LIMIT", 0, limit, "GET", STORAGE_HASH_PREFIX+"*->json", "DESC"))
	if err != nil {
		return nil, err
	}
	// this vs del?
	conn.Do("EXPIRE", tmpId, 1)
	return sortResult, err
}

func Search(types []string, query string, limit int) ([]byte, error) {

	conn := pool.Get()
	defer conn.Close()

	err := cacheIdsInSet(conn, query)
	if err != nil {
		return nil, err
	}

	var allResults []string
	for _, t := range types {
		typeResults, err := intersectAndSort(conn, query, limit, t)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, resultsToJSONArray(typeResults))
	}
	jsonResponse := typesAndResultsByTypeToJSON(types, allResults)
	return []byte(jsonResponse), err
}
