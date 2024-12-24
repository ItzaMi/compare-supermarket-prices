package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

var baseProductUrls []string = []string{
	"https://www.auchan.pt/on/demandware.store/Sites-AuchanPT-Site/pt_PT/Search-UpdateGrid?srule=price-low-to-high&q=",
	"https://www.continente.pt/pesquisa/?start=0&srule=price-low-to-high&q=",
}
var productList []string = []string{
	"leite soja",
	"pao mistura",
}

func parseAuchan(c *colly.Collector, url string, product string) [4]string {
	var result [4]string
	c.OnHTML("body div.auc-product:first-of-type div.product", func(e *colly.HTMLElement) {
		result[0] = product
		result[1] = e.ChildText("h3 a")
		result[2] = " "
		result[3] = e.ChildText("span.value")
	})
	c.Visit(url)
	return result
}

func parseContinente(c *colly.Collector, url string, product string) [4]string {
	var result [4]string
	c.OnHTML("div.product-grid div.productTile:first-of-type", func(e *colly.HTMLElement) {
		result[0] = product
		result[1] = e.ChildText("h2")
		result[2] = e.ChildText("p.pwc-tile--brand")
		result[3] = e.ChildText("span.ct-price-formatted")
	})
	c.Visit(url)
	return result
}

func visitStore(store string, c *colly.Collector, url string, product string) [4]string {
	var result [4]string
	switch {
	case store == "auchan":
		result = parseAuchan(c, url, product)
	case store == "continente":
		result = parseContinente(c, url, product)
	}
	return result
}

func main() {
	var visitedStoreName string

	c := colly.NewCollector()

	for i := 0; i < len(baseProductUrls); i++ {
		var stringSeparator string

		if strings.Contains(baseProductUrls[i], "auchan") {
			stringSeparator = "%20"
			visitedStoreName = "auchan"
		}

		if strings.Contains(baseProductUrls[i], "continente") {
			stringSeparator = "+"
			visitedStoreName = "continente"
		}

		fName := visitedStoreName + ".csv"
		file, err := os.Create(fName)
		if err != nil {
			log.Fatalf("Cannot create file %q: %s\n", fName, err)
			return
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()

		writer.Write([]string{"Product", "Name", "Brand", "Price"})

		for j := 0; j < len(productList); j++ {
			var product = productList[j]
			var newString = strings.Join(strings.Split(product, " "), stringSeparator)
			var url = baseProductUrls[i] + newString

			var info = visitStore(visitedStoreName, c, url, product)

			fmt.Println(info)

			writer.Write(info[:])
		}
	}
}
