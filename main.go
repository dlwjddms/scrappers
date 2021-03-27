package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

/*To navigate through HTML and to be able to find things in the HTML
we will going to use goquery, Its like jquery for GO
It allows us to navigate through HTML , inside of the HTML document
: install -> go get github.com/PuerkitoBio/goquery */
/*First we will get pages than we visit each page ,
 *and than we extract job from page and put it in to excel*/

var baseURL string = "https://cse.cau.ac.kr/sub05/sub0501.php"

//"https://kr.indeed.com/jobs?q=로봇엔지니어&l=&ts=1598092708172&rq=1&rsIdx=0&fromage=last&newcount=26"

func main() {
	pages := getPages()
	fmt.Println(pages)
}

/*fucking error.. I think I just did go mod init github.com/dlwjddms/scrappers .. and do it again ** */
/* return how many page is it*/
func getPages() int {
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close() // We need to close IO... after function is done prevent memory leak
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body) //Body is basically Byte IO
	checkErr(err)

	//class name of the page : pagination-list
	doc.Find("paging")
	fmt.Println(doc)

	return 0
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
