package main

import (
	"encoding/json"
	"fmt"
	"linn221/Requester/requests"
	"log"
)

func main() {
	results, err := requests.ParseHAR("example.har")
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
