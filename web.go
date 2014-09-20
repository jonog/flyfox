package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
)

const (
	DEFAULT_ITEM_LIMIT = 10
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

	items, err := Search(mux.Vars(req)["query"], limit)
	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
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
