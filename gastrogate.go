package main

// lunchHandler, gastrogate
// Handles the scrapping of gastrogate data from lunchmneus.

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Item struct {
	// Required
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	// Optional
	Content  string `xml:"encoded"`
	PubDate  string `xml:"pubDate"`
	Comments string `xml:"comments"`
}

// getGastrogateResturants, gets the menus from gastrogate homepage.
func getGastrogateResturants() {
	var wg sync.WaitGroup
	for resturant, uri := range gastrogate {
		wg.Add(1)
		go getGastrogateMenu(uri, resturant, &wg)
	}
	wg.Wait()
}

// getGastrogateMenu, get the menu html into a map.
// This is where the html is fetched.
// The map is setup with weekday as key and that weekday
// menu items as value.
func getGastrogateMenu(uri string, resturant string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, _ := html.Parse(resp.Body)
	table := getTableFromHTML(doc)
	week := getMenuItemsFromTable(table)
	menus[resturant] = week
}

// getMenuItemsFromTable, takes a html table and returns a
// map of the weeks menu. This is based on each weekday beeing an
// h3 on gastrogate.
func getMenuItemsFromTable(table *html.Node) map[string]string {
	var f func(*html.Node)
	week := make(map[string]string)
	var curday string
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h3" {
			daymap := strings.Fields(n.FirstChild.Data)
			curday = swedishWeekdaysToEnglish(strings.Title(daymap[0]))
		}
		if n.Type == html.ElementNode && n.Data == "td" && n.Attr[0].Val == "td_title" {
			week[curday] = week[curday] + strings.TrimSpace(n.FirstChild.Data) + "\n"
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(table)
	return week
}

// getTableFromHTML, returns a table from an *html.Node
func getTableFromHTML(document *html.Node) *html.Node {
	var table *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			table = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(document)
	return table
}
