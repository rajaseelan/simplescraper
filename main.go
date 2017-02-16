package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
	"os"
	"strings"
)

// helper func to pull href attribute from token
func getHref(t html.Token) (ok bool, href string) {

	// Iterate over all the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}

	}

	// bare return func, will return ok & href
	// as defined in func
	// UGLY!!	
	return
}

// extract all http** links from a given webpage
func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	// notify that we're done after this function
	defer func() {
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of doc, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// extract the href value
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// Make sure url begins w/ http**
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				ch <- url
			}
		}
	}
}


func main() {
	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// kick of crawl process concurrently
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
	}


	// Subscribe to both channels
	// what's happening:
	// listen to both channels, block until data received
	// if you get data from chUrls channel - add to the foundUrls map
	// getData from chFinished - increment c and start listening again
	// will only loop for the number of URLs supplied in the cmd line 
	//  - len of seed URLS
	for c := 0; c < len(seedUrls); {
		select {
		case url := <- chUrls:
			foundUrls[url] = true
		case <- chFinished:
			c++
		}
	}

	// parsing done. Print Urls
	

	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

	for url, _ := range foundUrls {
		fmt.Println(" - " + url)
	}

}
