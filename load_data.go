package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type StoredItemData struct {
	Collection []Item
}

func LoadData(filename string) {

	RedisInit()

	file, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	data := make([]Item, 0)
	json.Unmarshal(file, &data)

	fmt.Println("\nloading keys -> ")
	for _, item := range data {
		fmt.Println(item)
		item.Save()
	}

}
