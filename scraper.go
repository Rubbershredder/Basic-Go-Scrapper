package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Book struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Image string `json:"image"`
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
	)

	var books []Book

	c.OnHTML("article.product_pod", func(h *colly.HTMLElement) {
		book := Book{
			Name:  h.ChildAttr("h3 a", "title"),
			Image: h.ChildAttr("div.image_container img", "src"),
			Price: h.ChildText("p.price_color"),
		}

		// Clean up the image URL
		book.Image = "https://books.toscrape.com/" + strings.TrimPrefix(book.Image, "../")

		books = append(books, book)
	})

	c.OnScraped(func(r *colly.Response) {
		// for _, book := range books {
		// 	fmt.Printf("Name: %s\nImage: %s\nPrice: %s\n\n", book.Name, book.Image, book.Price)
		// }
		jsonData, err := json.MarshalIndent(books, "", " ")
		if err != nil {
			log.Fatal("Error marshalling JSON:", err)
		}

		err = ioutil.WriteFile("books.json", jsonData, 0664)
		if err != nil {
			log.Fatal("Error writing JSON file:", err)
		}
		fmt.Printf("Scraped %d books and saved to books.json\n", len(books))
	})

	err := c.Visit("https://books.toscrape.com/")
	if err != nil {
		log.Fatal(err)
	}
}
