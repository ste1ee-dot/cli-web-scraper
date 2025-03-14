// TODO   check GET return code for each link
//		if link is outside domain = skip
//		if link is inside domain = go to it and check for link
//		repeat

package main

import (
	"fmt"
	"net/http"
	"os"
	//	"strings"

	"golang.org/x/net/html"
)

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	return

}

func extractLinks(url string, ch chan string, chFinished chan bool) {

	resp, err := http.Get(url)

	defer func() {
		chFinished <- true

	}()

	if err != nil {
		fmt.Println("ERROR: Failed to get links:", url)
		return
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isA := t.Data == "a"
			if !isA {
				continue
			}

			ok, url := getHref(t)
			if !ok {
				continue
			}

			ch <- url

			//	hasHTTP := strings.Index(url, "http") == 0
			//	if hasHTTP {
			//		ch <- url
			//	}

		}

	}

}

func main() {

	//	url := "https://scrape-me.dreamsofcode.io"

	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]

	chUrls := make(chan string)
	chFinished := make(chan bool)

	for _, url := range seedUrls {
		go extractLinks(url, chUrls, chFinished)
	}

	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
	}

	fmt.Println("\nFound", len(foundUrls), "unique links:\n")

	for url, _ := range foundUrls {
		fmt.Println(" - " + url)
	}

	close(chUrls)
}
