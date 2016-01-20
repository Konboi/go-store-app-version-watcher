package scraper

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const (
	GOOGLE_PLAY_BASE_URI           = "https://play.google.com"
	GOOGLE_PLAY_RULE               = "div.meta-info div.content"
	GOOGLE_PLAY_VERSION_ATTR_NAME  = "itemprop"
	GOOGLE_PLAY_VERSION_ATTR_VALUE = "softwareVersion"
)

func GooglePlay(appId string) (version string, err error) {
	uri := fmt.Sprintf("%s/store/apps/details?id=%s", GOOGLE_PLAY_BASE_URI, appId)
	doc, err := goquery.NewDocument(uri)

	if err != nil {
		return version, err
	}

	doc.Find(GOOGLE_PLAY_RULE).Each(func(i int, s *goquery.Selection) {
		attr, _ := s.Attr(GOOGLE_PLAY_VERSION_ATTR_NAME)

		if attr == GOOGLE_PLAY_VERSION_ATTR_VALUE {
			version = s.Text()
		}
	})

	if version == "" {
		return version, fmt.Errorf("Can not find version at %s", uri)
	}

	return version, nil
}
