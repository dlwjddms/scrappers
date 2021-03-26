package main

import (
	"errors"
	"fmt"
	"net/http"
)

// we will going to build URL checker
/*
 * we will going to have URL slice
 * that we want to check
 * and check each URL
 * if URL is up we willl going to say OK or FAIL
 */
type rquestResult struct {
	url    string
	status string
}

var errRequestFailed = errors.New("Request Failed")

func main() {
	/*  not initialized.. just empty map
	// -> but you have to do initialize before using
	// var results = map[string]string{}
	// -> or use make (makes a make for me)
	*/
	results := make(map[string]string)
	c := make(chan rquestResult)
	/* we want to do this job all at once */
	urls := []string{
		"https://www.airbnb.com/",
		"https://google.com/",
		"https://reddit.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
	}

	for _, url := range urls {
		go hitURL(url, c)
	}
	for i := 0; i < len(urls); i++ {
		result := <-c
		results[result.url] = result.status
	}

	for url, status := range results {
		fmt.Println(url, status)
	}
}

/*function that will hit the websites*/
/*Receive from the channel,sometime you want to specify
* This func have channel, but this channel you can not send, you can just receive
* "c chan<- result" : send only  so "fmt.Println(<-c)"" won't work
 */
func hitURL(url string, c chan<- rquestResult) {
	/*we have to make a request for hit */
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
	}
	c <- rquestResult{url: url, status: status}
}
