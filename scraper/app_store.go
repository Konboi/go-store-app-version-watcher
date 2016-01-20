package scraper

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const (
	APP_STORE_BASE_URI           = "https://itunes.apple.com/jp/app"
	APP_STORE_RULE               = "div#left-stack .lockup ul.list li span"
	APP_STORE_VERSION_ATTR_NAME  = "itemprop"
	APP_STORE_VERSION_ATTR_VALUE = "softwareVersion"
)

func AppStore(appId string) (version string, err error) {
	uri := fmt.Sprintf("%s/%s", APP_STORE_BASE_URI, appId)
	doc, err := goquery.NewDocument(uri)

	if err != nil {
		return version, err
	}

	doc.Find(APP_STORE_RULE).Each(func(i int, s *goquery.Selection) {
		attr, _ := s.Attr(APP_STORE_VERSION_ATTR_NAME)

		if attr == APP_STORE_VERSION_ATTR_VALUE {
			version = s.Text()
		}
	})

	if version == "" {
		return version, fmt.Errorf("Can not find version at %s", uri)
	}

	return version, nil
}
