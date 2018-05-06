package main

// lunchGetter rss.
//  Handles the rss reads for lunchgetter.

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Rss2 struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	// Required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`
	// Optional
	PubDate  string `xml:"channel>pubDate"`
	ItemList []Item `xml:"channel>item"`
}

// GetRSSResturants, gets the RSS menu from the uris,
// in the global rssfeeds map.
func getRSSresturants(rssfeed map[string]string) {
	var wg sync.WaitGroup
	for resturant, uri := range rssfeed {
		wg.Add(1)
		go getRSSFeed(uri, resturant, &wg)
	}
	wg.Wait()
}

// getRSSFeed, populates the sent in week map with the
// resturants menu in the format week[weekday] = "menustring"
func getRSSFeed(uri string, resturant string, wg *sync.WaitGroup) {
	defer wg.Done()
	week := make(map[string]string)
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	r := Rss2{}
	err = xml.Unmarshal([]byte(data), &r)
	if err != nil {
		log.Fatal(err)
	}
	var description string
	for _, day := range r.ItemList {
		title := strings.Fields(day.Title)
		description = day.Description
		description = strings.Replace(description, "<p>", "", -1)
		description = strings.Replace(description, "</p>", "", -1)
		description = strings.Replace(description, "<br />", "\n", -1)
		description = strings.Replace(description, "<strong>", "\n", 1)
		description = strings.Replace(description, "</strong>", "\n", 1)
		description = strings.Replace(description, "<div>", "\n", 1)
		description = strings.Replace(description, "</div>", "\n", 1)
		curday := swedishWeekdaysToEnglish(strings.Title(title[0]))
		week[curday] = description
	}
	menus[resturant] = week
}
