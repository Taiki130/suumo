package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type property struct {
	Name     string
	Cost     int
	Style    string
	Age      string
	Nearests []string
	URL      string
}

func main() {
	urls := []string{
		"https://suumo.jp/chintai/bc_100389341666/",
		"https://suumo.jp/chintai/bc_100391519369/",
		"https://suumo.jp/chintai/bc_100401120582/",
		"https://suumo.jp/chintai/bc_100401129371/",
		"https://suumo.jp/chintai/bc_100401976254/",
		"https://suumo.jp/chintai/bc_100394102781/",
		"https://suumo.jp/chintai/bc_100395160410/",
		"https://suumo.jp/chintai/bc_100402562270/",
		"https://suumo.jp/chintai/bc_100395653896/",
		"https://suumo.jp/chintai/bc_100361824432/",
	}

	var properties []property

	// 物件名
	nameSelector := ".section_h1-header > h1"
	// 家賃
	rentSelector := ".property_view_main-emphasis"
	// 管理費
	manageRentSelector := ".property_view_main-data > div > div.property_data-body"
	// 間取り
	styleSelector := ".property_view_detail-body > ul > li:nth-child(1) > div > div.property_data-body"
	// 築年数
	ageSelector := ".property_view_detail-body > ul > li:nth-child(5) > div > div.property_data-body"
	// 最寄り駅
	nearestSelector := ".l-property_view_detail-sub > div:nth-child(1) > div > div.property_view_detail-body"
	nearestInnerSelector := ".property_view_detail-text"

	for _, url := range urls {
		c := colly.NewCollector()

		var name, rent, manageRent, style, age string
		c.OnHTML(nameSelector, func(e *colly.HTMLElement) {
			name = strings.TrimSpace(e.Text)
		})

		c.OnHTML(rentSelector, func(e *colly.HTMLElement) {
			rent = strings.TrimSpace(e.Text)
		})
		c.OnHTML(manageRentSelector, func(e *colly.HTMLElement) {
			manageRent = strings.TrimSpace(e.Text)
		})
		c.OnHTML(styleSelector, func(e *colly.HTMLElement) {
			style = strings.TrimSpace(e.Text)
		})
		c.OnHTML(ageSelector, func(e *colly.HTMLElement) {
			age = strings.TrimSpace(e.Text)
		})

		var nearests []string
		c.OnHTML(nearestSelector, func(e *colly.HTMLElement) {
			e.ForEach(nearestInnerSelector, func(_ int, s *colly.HTMLElement) {
				nearests = append(nearests, s.Text)
			})
		})
		c.Visit(url)

		var cost, manageRentInt int
		var err error
		if strings.Contains(manageRent, "-") {
			manageRentInt = 0
		} else {
			manageRentInt, err = strconv.Atoi(strings.Replace(manageRent, "円", "", -1))
			if err != nil {
				log.Fatal(err)
			}
		}

		trimedRent := strings.Replace(rent, "万円", "", -1)
		if strings.Contains(trimedRent, ".") {
			rentInt, err := strconv.ParseFloat(trimedRent, 64)
			if err != nil {
				log.Fatal(err)
			}
			cost = int(rentInt*10000) + manageRentInt
		} else {
			rentInt, err := strconv.Atoi(trimedRent)
			if err != nil {
				log.Fatal(err)
			}
			cost = rentInt*10000 + manageRentInt
		}

		properties = append(properties, property{
			Name:     name,
			Cost:     cost,
			Style:    style,
			Age:      age,
			Nearests: nearests,
			URL:      url,
		})
	}

	fmt.Println("| 物件名 | 家賃（管理費込み） | 間取り | 築年数 | 最寄り駅 | URL | 備考 | 評価① | 評価② |")
	fmt.Println("| ----- | --------------- | ----- | ----- | ------ | ---- | ---- | ----- | ------ |")
	for _, prop := range properties {
		fmt.Printf("| %s | %d | %s | %s | %v | %s | | | |\n", prop.Name, prop.Cost, prop.Style, prop.Age, prop.Nearests, prop.URL)
	}
}
