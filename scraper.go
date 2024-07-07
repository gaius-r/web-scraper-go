package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	// importing colly for web-scraping
	"github.com/gocolly/colly"
)

// defining a data structure to store the scraped data
type Product struct {
	url, image, name, price string
}

func main() {
	// scraping logic

	var products []Product

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
