package main

import (
	"fmt"
	"log"

	// importing colly for web-scraping
	"github.com/gocolly/colly"
)

func main() {
	// scraping logic
	fmt.Println("Creating Collector")
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Somehting went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		// printing all URLs associated with the a links in the page
		fmt.Printf("%v", e.Attr("href"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	c.Visit("https://en.wikipedia.org/wiki/Main_Page")
	fmt.Println("Hello World!")
}
