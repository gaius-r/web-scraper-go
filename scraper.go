package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	// importing colly for web-scraping/crawling
	"github.com/gocolly/colly"
)

// defining a data structure to store the scraped data
type Product struct {
	url, image, name, price string
}

func main() {

	// initializing the list of products to scrape with an empty slice
	var products []Product

	// initializing the list of pages to scrape with an empty slice
	var pagesToScrape []string

	// the first pagination URL to scrape
	pageToScrape := "https://www.scrapingcourse.com/ecommerce/page/1/"

	// initializing the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	// current iteration
	i := 1
	// max pages to scrape
	limit := 5

	// initializing a Colly instance
	fmt.Println("Creating Collector")
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	// iterating over the list of pagination links to implement the crawling logic
	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		// discovering a new page
		newPaginationLink := e.Attr("href")

		// if the discovered page is new
		if !contains(pagesToScrape, newPaginationLink) {
			// if the page discovered should be scraped
			if !contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})

	// iterating over the list of HTML product elements
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		// initializing a new product instance
		product := Product{}

		// scraping the data of interest
		product.url = e.ChildAttr("a", "href")
		product.image = e.ChildAttr("img", "src")
		product.name = e.ChildText("h2")
		product.price = e.ChildText(".price")

		// adding the product instance with scraped data to the list of products
		products = append(products, product)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")

		// until there is still a page to scrape
		if len(pagesToScrape) != 0 && i < limit {
			// getting the current page to scrape and removing it from the list
			pageToScrape := pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			// incrementing the iteration count
			i++

			// visiting a new page
			c.Visit(pageToScrape)
		}

		// opening the CSV file
		file, err := os.Create("Products.csv")
		if err != nil {
			log.Fatalln("Failed to create output CSV file!", err)
		}
		defer file.Close()

		// initializing a file writer
		writer := csv.NewWriter(file)

		// writing the CSV headers
		headers := []string{
			"url",
			"image",
			"name",
			"price",
		}

		writer.Write(headers)

		// writing each product as a CSV row
		for _, product := range products {
			// converting a Product to an array of strings
			record := []string{
				product.url,
				product.image,
				product.name,
				product.price,
			}

			// adding a CSV record to the output file
			writer.Write(record)
		}
		defer writer.Flush()
	})

	c.Visit("https://www.scrapingcourse.com/ecommerce/")
	fmt.Println("Hello World!")
}

// contains checks if a slice contains a specific element
func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
