package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	DEFAULT_ITEM_LIMIT = 10
	CLR_0 = "\x1b[30;1m"
	CLR_R = "\x1b[31;1m"
	CLR_G = "\x1b[32;1m"
	CLR_Y = "\x1b[33;1m"
	CLR_B = "\x1b[34;1m"
	CLR_M = "\x1b[35;1m"
	CLR_C = "\x1b[36;1m"
	CLR_W = "\x1b[37;1m"
	CLR_N = "\x1b[0m"
)



func WebInit() {

	RedisInit()

	m := mux.NewRouter()
	m.HandleFunc("/search/{query}", Query).Methods("GET")

	http.Handle("/", m)
	fmt.Println("Listening on Port 3001")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		panic(err)
	}
}

func Query(res http.ResponseWriter, req *http.Request) {

	queryParams := req.URL.Query()
	limit := getLimit(queryParams)

	startTime := time.Now().UTC()
	query := mux.Vars(req)["query"]
	items, err := Search(query, limit)
	if err != nil {
		panic(err)
	}
	fmt.Println(CLR_W, "Search query: ", CLR_Y, query, CLR_R, time.Now().UTC().Sub(startTime), CLR_N)

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Write(items)
}

func getLimit(v url.Values) int {

	if len(v["limit"]) != 1 {
		return DEFAULT_ITEM_LIMIT
	}

	limit, parseErr := strconv.Atoi(v["limit"][0])
	if parseErr != nil {
		return DEFAULT_ITEM_LIMIT
	}
	return limit
}
