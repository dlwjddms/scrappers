package scrap

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

/*Two Channel
 * main <-> getPage
 * getPage <-> extract Job*/

//Scrape Indeed
func Scrape(term string) {
	fmt.Print(term)
	var baseURL string = "https://cse.cau.ac.kr/sub05/sub0501.php"
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseURL)
	// hit the url
	for i := 0; i < totalPages; i++ {
		//[] + [] + [] = []
		// how can you combine many array? -> use "..."
		go getPage(baseURL, i, c)

	}
	// wait for go routine
	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(baseURL, jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func getPage(baseURL string, page int, mainC chan<- []extractedJob) {
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
		if i > 0 { //if i > 0  그냥 우리학교 홈페이지 떄문에 추가한거지 보통은 아님
			go extractJob(row, c)
		}

	})
	for i := 1; i < searchCards.Length(); i++ { //  i := 1 도 그냥 우리학교 홈페이지 떄문에 추가한거지 보통은 아님
		job := <-c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}

func extractJob(row *goquery.Selection, c chan<- extractedJob) {
	s := row.Find(".aleft  a")
	innerURL, _ := s.Attr("href")
	title := CleanString(s.Text())
	date := CleanString(row.Find(".pc-only").Eq(2).Text())
	views := row.Find(".pc-only").Eq(3).Text()
	c <- extractedJob{
		detailURL: innerURL,
		title:     title,
		date:      date,
		views:     views}
}

// CleanString cleans a String
func CleanString(str string) string {
	//Fields split the string ... make an array only text
	//join take array and put it together with seperator
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

/*fucking error.. I think I just did go mod init github.com/dlwjddms/scrappers .. and do it again ** */
/* return how many page is it*/
func getPages(baseURL string) int {
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

// write it at csv
func writeJobs(baseURL string, jobs []extractedJob) {
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
