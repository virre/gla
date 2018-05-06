/*
goLunchAggregator v.01
This will fetch menus from all configured lunch resturants.
Currently works with rssfeeds and gastrogate menus.

@TODO:
	- Regexp replace tags instead of exact.
	- Add option to serve as website.
	- Try to sort on diet options. (When info available.)
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

var menus = map[string]map[string]string{}
var rssfeeds = map[string]string{}
var gastrogate = map[string]string{}

func main() {
	usr, _ := user.Current()
	home := usr.HomeDir
	if home == "" {
		home = "/tmp"
	}
	settingsPtr := flag.String("settings", home+"/.lunchsettings", "The path to the config file that contains the settings.")
	weekdayPtr := flag.String("weekday", "", "Look for a certain weekday")
	todayPtr := flag.Bool("today", false, "Enable for todays menu")
	flag.Parse()
	weekday := *weekdayPtr
	if *todayPtr != false {
		now := time.Now()
		weekday = now.Weekday().String()
	}
	setUrisFromSettings(*settingsPtr)
	getRSSresturants(rssfeeds)
	getGastrogateResturants()
	for resturant, menu := range menus {
		fmt.Printf("%s\n------\n", resturant)
		if weekday != "" {
			weekday = swedishWeekdaysToEnglish(strings.Title(weekday))
			fmt.Printf("%s\n%s\n", weekday, menu[weekday])
		} else {
			for weekday, menu := range menu {
				fmt.Printf("Menu for %s\n%s\n\n", weekday, menu)
			}
		}
		fmt.Printf("\n\n")
	}
}

//setUrisFromSettings, load gastrogate and rssfeeds from an ini formated config file.
func setUrisFromSettings(settings_path string) {
	cfg, err := ini.Load(settings_path)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Could not load settings file. Standard is ~/.lunchsettings")
	}
	rss, _ := cfg.GetSection("RSSfeed")
	gastrogates, _ := cfg.GetSection("Gastrogate")
	rssfeeds = rss.KeysHash()
	gastrogate = gastrogates.KeysHash()
}

// swedishWeekdaysToEnglish, Returns weekday in English from a Swedish string.
func swedishWeekdaysToEnglish(weekday string) string {
	weekday = strings.Replace(weekday, ",", "", -1)
	switch weekday {
	case "Måndag":
		return "Monday"
	case "Tisdag":
		return "Tuesday"
	case "Onsdag":
		return "Wednesday"
	case "Torsdag":
		return "Thursday"
	case "Fredag":
		return "Friday"
	case "Lördag":
		return "Saturday"
	case "Söndag":
		return "Sunday"
	}
	return weekday
}
