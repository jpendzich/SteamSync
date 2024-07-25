package internal

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetSaves(name string) (string, string, error) {
	windows := ""
	linux := ""
	doc, err := goquery.NewDocument("https://www.pcgamingwiki.com/wiki/" + name)
	if err != nil {
		return windows, linux, err
	}

	doc.Find("h3").EachWithBreak(func(i int, header *goquery.Selection) bool {
		if header.Text() == "Save game data location" {
			table := header.Next().Find("table")
			table.Find("tr").Each(func(i int, row *goquery.Selection) {
				if strings.TrimSpace(row.Find("th").Text()) == "Windows" {
					windows = row.Find("td").Find("span").Text()
				} else if strings.TrimSpace(row.Find("th").Text()) == "Steam Play (Linux)" {
					linux = row.Find("td").Find("span").Text()
				}
			})
			return false
		}
		return true
	})
	return windows, linux, nil
}
