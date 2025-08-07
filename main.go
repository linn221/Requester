package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"linn221/Requester/requests"
	"log"
	"strings"
)

func main() {
	// if len(os.Args) <= 0 {
	// 	log.Fatal("supply a har file")
	// }

	data, err := ioutil.ReadFile("a.har")
	if err != nil {
		log.Fatal(err)
	}

	results, err := requests.ParseHAR(data, func(my *requests.MyRequest) string {
		var contentType string
		for _, header := range my.ResHeaders {
			if strings.ToLower(header.Name) == "content-type" {
				contentType = header.Name + ":" + header.Value
				break
			}
		}
		s := fmt.Sprintf("%d %d %s %s",
			my.ResStatus, my.RespSize, my.ResBody, contentType,
		)
		return s
	})

	if err != nil {
		log.Fatal(err)
	}

	s, _ := json.Marshal(&results)

	fmt.Println(s)
}
