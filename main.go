package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// we will going to build URL checker
/*
 * we will going to have URL slice
 * that we want to check
 * and check each URL
 * if URL is up we willl going to say OK or FAIL
 */

var errRequestFailed = errors.New("Request Failed")

func main() {
	/*  not initialized.. just empty map
	// -> but you have to do initialize before using
	// var results = map[string]string{}
	// -> or use make (makes a make for me)
	*/
	var results = make(map[string]string)
	/* we want to do this job all at once */
	urls := []string{
		"https://www.airbnb.com/",
		"https://google.com/",
		"https://reddit.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
	}

	for _, url := range urls {
		result := "OK" //default
		err := hitURL(url)
		if err != nil {
			result = "FAILED"
		}
		results[url] = result
	}

}

/*function that will hit the websites*/
func hitURL(url string) error {
	/*we have to make a request for hit */
	fmt.Println("Checking: ", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		/*status code 100,200,300 redirection 400 things are wrong*/
		fmt.Println(err, resp.StatusCode)
		return errRequestFailed
	}
	return nil
}

/*
 * The way to optimize things in GO is by doing concurrency
 * Goroutines : functions do same time with other functions
 *
 * go do things like make child thread I think..
 * so main func finish go also dies
 * -> how do we communication with main func and Goroutine ?
 * = Channel(how do we communicate with goroutine)
 * main funciton store the result
 *
 * goroutine send message to main fucnction
 */
func sexy() {
	/* what information you will going to send? */
	channel := make(chan string)
	people := [2]string{"nico", "JE"}
	for _, person := range people {
		go isSexy(person, channel)
	}
	/* recieve message from channel, main func wait until he gets one message
	   if he get ine mesg he starts again  */
	//fmt.Println(<-channel)
	/* now we xan receive two msg, go knows how many goroutines are running */
	//fmt.Println(<-channel)
	/*but if people size got bigger? -> use for loops*/
	for i := 0; i < len(people); i++ {
		/*recieving msg is a blocking operation*/
		fmt.Println((<-channel))
	}
	// go sexyCount("nico")
	// go sexyCount("JE")
	//time.Sleep(time.Second * 5)
}

/* channel , the function is having one channel to communicate with main func */
func isSexy(person string, c chan string) {
	time.Sleep(time.Second * 5)
	c <- person + " is sexy"
}

/* goroutine */
func sexyCount(person string) {
	for i := 0; i < 10; i++ {
		fmt.Println(person, "is sexy", i)
		time.Sleep(time.Second)
	}
}
