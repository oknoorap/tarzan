package task

import (
	"./../request"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"reflect"
	"log"
)

func GetPageItems (queue string, args ...interface{}) error {
	// Parse Arguments
	arg := reflect.ValueOf(args[0]).Interface().(map[string]interface{})

	// Parse value
	uri := arg["url"].(string)

	// Open page, without referrer
	log.Println("Scraping Category", uri)
	tor := new(request.Tor)
	body, err := tor.Open(uri, "")

	// Tell error to console
	if err != nil {
		log.Println("Error:", err)
		return nil
	}

	// Parse html
	html := strings.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}

	// Connect redis
	redis, err := dialRedis()
	if err != nil {
		log.Println(err)
		return nil
	}

	// Get title link
	doc.Find(".item-thumbnail__image a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			if !strings.Contains(href, "full_screen_preview") {
				err := redis.RPush("resque:queue:" + queue, `{"class":"GetItem","args":[{"url":"http://themeforest.net`+ href +`"}]}`).Err()

				if err != nil {
					log.Println(err)
				}
			}
		}
	})

	// Get next page
	next_page, np_ok := doc.Find(".pagination__next").Attr("href")
	if np_ok {

		err := redis.RPush("resque:queue:" + queue, `{"class":"GetPageItems","args":[{"url":"http://themeforest.net`+ next_page +`"}]}`).Err()

		if err != nil {
			log.Println(err)
		}
	}
	
	return nil
}