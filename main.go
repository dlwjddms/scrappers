package main

import (
	"os"
	"strings"

	"github.com/dlwjddms/scrappers/github.com/dlwjddms/scrappers/scrap"
	"github.com/labstack/echo"
)

const fileName = "cau_notice.csv"

func handleHome(c echo.Context) error {
	// respond to this html
	return c.File("home.html")
}
func handleScrape(c echo.Context) error {
	defer os.Remove(fileName)
	term := strings.ToLower(scrap.CleanString(c.FormValue("term")))
	scrap.Scrape(term)
	// returns a file
	return c.Attachment(fileName, fileName)
}

// we will make go sever with echo
func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))

}
