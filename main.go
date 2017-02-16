package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("This is a simple scraper based the tutorial: https://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html")

}

// Makes a httpRequest
// accepts a string
func makeHTTPRequest(url string) {
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("HTML:\n\n", string(bytes))

	resp.Body.Close()
}
