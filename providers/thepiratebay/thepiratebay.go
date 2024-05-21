package thepiratebay

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"

	"github.com/tnychn/torrodle/models"
	"github.com/tnychn/torrodle/request"
)

const (
	Name = "ThePirateBay"
	Site = "https://thepiratebay.org"
)

type provider struct {
	models.Provider
}

func New() models.ProviderInterface {
	provider := &provider{}
	provider.Name = Name
	provider.Site = Site
	provider.Categories = models.Categories{
		All:   "/search.php?q=%v&all=on&search=Pirate+Search&page=%d&orderby=",
		Movie: "/search.php?q=%v&video=on&search=Pirate+Search&page=%d&orderby=",
		TV:    "/search.php?q=%v&video=on&search=Pirate+Search&page=%d&orderby=",
		Porn:  "/search.php?q=%v&porn=on&search=Pirate+Search&page=%d&orderby=",
		// All:   "/search/%v/%d/99/0",
		// Movie: "/search/%v/%d/99/200",
		// TV:    "/search/%v/%d/99/200",
		// Porn:  "/search/%v/%d/99/500",
	}
	return provider
}

func (provider *provider) Search(query string, count int, categoryURL models.CategoryURL) ([]models.Source, error) {
	results, err := provider.Query(query, categoryURL, count, 30, 0, extractor)
	return results, err
}

func extractor(surl string, page int, results *[]models.Source, wg *sync.WaitGroup) {
	logrus.Infof("ThePirateBay: [%d] Extracting results...\n", page)
	_, html, err := request.Get(nil, surl, nil)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("ThePirateBay: [%d]", page), err)
		wg.Done()
		return
	}
	var sources []models.Source
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	ol := doc.Find("ol#torrents")
	ol.Find("li.list-entry").Each(func(i int, li *goquery.Selection) {
		// tds := li.Find("td")
		// a := tds.Eq(1).Find("a.detLink")
		// title
		title := li.Find("span.list-item.item-name.item-title a").Text()
		// seeders
		// s := tds.Eq(2).Text()
		s := li.Find("span.list-item.item-seed").Text()
		seeders, _ := strconv.Atoi(strings.TrimSpace(s))
		// leechers
		// l := tds.Eq(3).Text()
		l := li.Find("span.list-item.item-leech").Text()
		leechers, _ := strconv.Atoi(strings.TrimSpace(l))
		// filesize
		// re := regexp.MustCompile(`Size\s(.*?),`)
		// text := tds.Eq(1).Find("font").Text()
		// fs := re.FindStringSubmatch(text)[1]
		fs := li.Find("span.list-item.item-size").Text()
		filesize, _ := humanize.ParseBytes(strings.TrimSpace(fs)) // convert human words to bytes number
		// url
		URL, _ := li.Find("span.list-item.item-name.item-title a").Eq(0).Attr("href")
		// magnet
		magnet, _ := li.Find(`span.item-icons a`).Eq(0).Attr("href")
		// ---
		source := models.Source{
			From:     "ThePirateBay",
			Title:    strings.TrimSpace(title),
			URL:      Site + URL,
			Seeders:  seeders,
			Leechers: leechers,
			FileSize: int64(filesize),
			Magnet:   magnet,
		}
		sources = append(sources, source)
	})

	logrus.Debugf("ThePirateBay: [%d] Amount of results: %d", page, len(sources))
	*results = append(*results, sources...)
	wg.Done()
}
