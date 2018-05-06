# goLunchAggregator

## Introduction
I started to write this when I started learning Go, it proved to be a very useful tool to avoid
looking at multiple different websites to see if there is a good lunch dailys that matches
all demands in a lunch group. RSSfeed and Gastrogate is all I needed to aggregate
so far so thats why that is the options.

## Configuration
 This uses a settings-file caleld .lunchsettings in $HOME dir (this program have only been
   tested on Linux and MacOs there might be issues on windows because of this).
This file follow the classic ini settings syntax. The below is the one I tested with
(so it is based on Kista, Stockholm):

```
[RSSfeed]
Tastory = http://www.tastory.com/modules/MenuRss/MenuRss/CurrentWeek?costNumber=6209&language=sv

[Gastrogate]
Uppereast = "https://uppereast.gastrogate.com/lunch/"
Nordicforum = "https://nordicforum.gastrogate.com/lunchmeny/"
```

## Usage:
gla - Will give you all the weeks entry to stdout.
gla -today will give you todays menu options
gla -weekday=$WEEKDAY will give you $WEEKDAYs lunchoption, weekday can be in Swedish or English currently.
gla -settings=$PATH will change where to look for settings file from the default of $HOME/.lunchsettings to $PATH
