package main

func main() {
	// if len(os.Args) <= 0 {
	// 	log.Fatal("supply a har file")
	// }

	// data, err := ioutil.ReadFile("a.har")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// results, err := requests.ParseHAR(data, func(my *requests.MyRequest) (string, string) {
	// 	reqText := my.URL + " " + my.Method + " " + my.ReqBody + " " + my.ReqHeaders.EchoMatcher("cookies")

	// 	respText := fmt.Sprintf("%d %d %s %s",
	// 		my.ResStatus, my.RespSize, my.ResBody, my.ResHeaders.EchoFilter("date"),
	// 	)
	// 	return reqText, respText
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// s, _ := json.Marshal(&results)

	// fmt.Println(s)
}
