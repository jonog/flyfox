package main

import (
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func instructions() {
	fmt.Println("Command: flyfox [load_data <file>|web]")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		instructions()
	}
	switch os.Args[1] {
	case "load_data":
		if len(os.Args) != 3 {
			instructions()
		}
		LoadData(os.Args[2])
	case "web":
		WebInit()
	default:
		instructions()
	}
}
