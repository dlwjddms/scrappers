package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	detailURL string
	title     string
	date      string
	views     string
}

var baseURL string = "https://cse.cau.ac.kr/sub05/sub0501.php" // offset =1 means first page

/*Two Channel
 * main <-> getPage
 * getPage <-> extract Job*/

func main() {
	var jobs []extractedJob
	totalPages := getPages()
	// hit the url
	for i := 0; i < totalPages; i++ {
		//[] + [] + [] = []
		// how can you combine many array? -> use "..."
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}
	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

// write it at csv
func writeJobs(jobs []extractedJob) {
	// create file
	file, err := os.Create("cau_notice.csv")
	checkErr(err)
	// 1. create writer 2, give data to writer 3. take data Flush to the file
	w := csv.NewWriter(file)
	defer w.Flush() // write the data to the file

	headers := []string{"Title", "date", "detailURL", "views"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.title, job.date, baseURL + job.detailURL, job.views}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func getPage(page int) []extractedJob {
	//empty slice of jobs
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "?offset=" + strconv.Itoa(page) // int to string
	//fmt.Println("Requesting URL :", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close() // We need to close IO... after function is done prevent memory leak
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body) //Body is basically Byte IO
	checkErr(err)

	//searchCards := doc.Find(".aleft")

	searchCards := doc.Find(".table-basic  tr")

	searchCards.Each(func(i int, row *goquery.Selection) { //s is each card~!
		go extractJob(row, c)

	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	return jobs
}

func extractJob(row *goquery.Selection, c chan<- extractedJob) {
	s := row.Find(".aleft  a")
	innerURL, _ := s.Attr("href")
	title := cleanString(s.Text())
	date := cleanString(row.Find(".pc-only").Eq(2).Text())
	views := row.Find(".pc-only").Eq(3).Text()
	c <- extractedJob{
		detailURL: innerURL,
		title:     title,
		date:      date,
		views:     views}
}

/*make string clear because there would be too many spaces
-> pkg strings -> TrimSpace*/
/*I want to clear the string and make evey word in to seperate things
erase small and big space and only word in array*/
func cleanString(str string) string {
	//Fields split the string ... make an array only text
	//join take array and put it together with seperator
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

/*fucking error.. I think I just did go mod init github.com/dlwjddms/scrappers .. and do it again ** */
/* return how many page is it*/
func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close() // We need to close IO... after function is done prevent memory leak
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body) //Body is basically Byte IO
	checkErr(err)

	//class name of the page : pagination-list
	doc.Find(".paging").Each(func(i int, s *goquery.Selection) {
		// we have next , big next and no 1
		// keep next because we don't count for "1" (first page)
		//fmt.Println(s.Find("a").Length()) // count the link
		pages = s.Find("a").Length()
	})
	//fmt.Println(doc)

	return pages
}
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request Failed sith Status: ", res.StatusCode)
	}
}
