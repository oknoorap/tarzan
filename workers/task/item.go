package task

import (
	"./../request"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"strings"
	"strconv"
	"reflect"
	"log"
	"time"
)

func GetItem (queue string, args ...interface{}) error {
	// Parse Arguments
	arg := reflect.ValueOf(args[0]).Interface().(map[string]interface{})

	// Parse value
	uri := arg["url"].(string)

	// Open page, without referrer
	log.Println("Scraping Item", uri)
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

	// Default data
	data := map[string]interface{}{
		"item_id": 0,
		"title": "",
		"sales": 0,
		"price": 0,
		"created": "",
		"author": "",
		"category": "",
		"url": "",
		"tags": []string{},
	}

	/**
	 * Parsing data
	 */
	
	// Item Id
	item_id := getItemId(uri)
	data["item_id"] = item_id
	
	// Title
	title := doc.Find("title").Text()
	data["title"] = strings.Trim(title, " ")

	// Sales
	sales, sales_ok := doc.Find(`meta[itemprop="interactionCount"]`).Attr("content")
	if sales_ok {
		sales_count := strings.Replace(sales, "UserDownloads:", "", 1)
		if sales_int, err := strconv.Atoi(sales_count); err == nil {
			data["sales"] = sales_int
		}
	}

	// Date uploaded
	created, created_ok := doc.Find(`time[itemprop="dateCreated"]`).Attr("datetime")
	if created_ok {
		data["created"] = created
	}

	// Tags
	var tags []string
	doc.Find(".meta-attributes__attr-tags a").Each(func(i int, s *goquery.Selection) {
		tags = append(tags, s.Text())
	})
	data["tags"] = tags

	// Author
	author := doc.Find(`a[rel="author"]`).Text()
	data["author"] = strings.Trim(author, " ")

	// Category
	category := doc.Find(`a[itemprop="genre"]`).Text()
	data["category"] = strings.Trim(category, " ")

	// Get Item price
	price, price_ok := doc.Find(`meta[itemprop="price"]`).Attr("content")
	if price_ok {
		price_float, err := strconv.ParseFloat(price, 64)
		if err == nil {
			data["price"] = price_float
		} else {
			log.Println(err)
		}
	}

	if data["author"] != "" {

		// Connect to database
		db, err := connectDb()
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		} else {

			// Set Mgo Session
			db.SetMode(mgo.Monotonic, true)
			dbSession := db.Copy()
			defer dbSession.Close()

			// Pick MongoDB collection
			collection := dbSession.DB("tarzan").C("item")

			// MongoDB data
			now := int32(time.Now().Unix())
			fields := bson.M{
				"$set": bson.M{
					"item_id": data["item_id"],
					"title": data["title"],
					"created": data["created"],
					"author": data["author"],
					"category": data["category"],
					"price": data["price"],
					"url": uri,
					"tags": tags,
					"time": now,
				},

				"$push": bson.M{
					"sales": bson.M{
						"date": now,
						"value": data["sales"],
					},
				},
			}

			// Upsert `page` collection
			_, err := collection.Upsert(bson.M{"item_id": item_id}, fields)
			if err != nil {
				log.Println("Error:", err)
			} else {
				if data["item_id"] != 0 {
					log.Println("Done, ItemId: ", data["item_id"])
				}
			}
		}
	} else {
		// Wait 10 seconds, after that push existing error task as new task
		timer := time.NewTimer(time.Second * 10)
    	<- timer.C

		// Connect redis
		redis, err := dialRedis()
		if err != nil {
			log.Println(err)
		}

		// Since we're have a connection trouble
		// Add again to task
		rpush := redis.RPush("resque:queue:" + queue, `{"class":"GetItem","args":[{"url":"`+ uri +`"}]}`).Err()
		if rpush != nil {
			log.Println(rpush)
		}

		log.Println("Item ID not found, recrawling")
		return nil
	}

	return nil
}


func getItemId (uri string) (output int) {
	output = 0

	if parsedUri, err := url.Parse(uri); err == nil {
		pathSlice := strings.Split(parsedUri.Path, "/")
		lastSlice := pathSlice[len(pathSlice)-1]
		if itemId, err := strconv.Atoi(lastSlice); err == nil {
			output = itemId
		}
	}

	return output
}