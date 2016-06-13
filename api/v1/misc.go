package v1

import (
	"../request"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/url"
	"net/http"
	"encoding/json"
	"strconv"
	"reflect"
	"strings"
	"time"
)

func CategoryList (c echo.Context) error {
	// Set default response to error
	defaultResponse := Error{
		Error: true,
		Message: "Unknown Error",
	}

	// Connect to mongodb
	db, err := connectDb()

	if err != nil {

		log.Fatalf("CreateSession: %s\n", err)
		defaultResponse.Message = "Can't connect db"

	} else {


		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("item")

		// Iterate all list
		var result jsondata
		err := collection.Find(nil).Distinct("category", &result)
		result.delete("")

		// Send response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
				"list": result,
			})
		} else {
			defaultResponse.Message = err.Error()
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}

func GetPreview (c echo.Context) error {
	// Set default response to error
	defaultResponse := Error{
		Error: true,
		Message: "Unknown Error",
	}

	uri := c.QueryParam("uri")
	item_id := getItemId(uri)


	// Connect DB
	db, err := connectDb()
	if err != nil {
		defaultResponse.Message = err.Error()
	} else {
		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()
		
		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("item")

		// Find if there is imge preview
		var result struct {
			Img_preview string `json:"url"`
		}

		err := collection.Find(bson.M{
			"item_id": item_id,
			"img_preview": bson.M{"$exists": 1},
		}).One(&result)

		if err != nil {
			// Open page, without referrer
			log.Println("Get Preview Item", uri)
			tor := new(request.Tor)
			body, err := tor.Open(uri, "")

			// Tell error to console
			if err != nil {
				defaultResponse.Message = err.Error()
			}

			// Parse html
			html := strings.NewReader(body)
			doc, err := goquery.NewDocumentFromReader(html)
			if err != nil {
				defaultResponse.Message = err.Error()
			}

			img, img_ok := doc.Find(`[property="og:image"]`).Attr("content")
			if img_ok && item_id != 0 {

			}

			// Update `image preview` collection
			err = collection.Update(bson.M{"item_id": item_id}, bson.M{
				"$set": bson.M{"img_preview": img},
			})

			// Send Response
			if err == nil {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"status": http.StatusOK,
					"img": img,
				})
			} else {
				defaultResponse.Message = err.Error()
			}
		} else {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
				"img": result,
			})
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
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

func Search (c echo.Context) error {

	// Set default response to error
	defaultResponse := Error{
		Error: true,
		Message: "Unknown Error",
	}

	// Connect to mongodb
	db, err := connectDb()

	if err != nil {

		log.Fatalf("CreateSession: %s\n", err)
		defaultResponse.Message = "Can't connect db"

	} else {
		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()

		// Decode search criteria into bson
		var search bson.M
		searchParams := c.FormValue("search")
		if err := json.Unmarshal([]byte(searchParams), &search); err != nil {
			logPanic(err)
		}

		// Convert string to bson.RegEx
		// And then delete it
		// Author
		operator := "$regex"
		author := reflect.ValueOf(search["author"])
		if author.IsValid() {
			regex := author.MapIndex(reflect.ValueOf("regex")).Elem().String()
			options := author.MapIndex(reflect.ValueOf("options")).Elem().String()
			if author.MapIndex(reflect.ValueOf("not")).Elem().Bool() {
				operator = "$not"
			}
			author.SetMapIndex(reflect.ValueOf(operator), reflect.ValueOf(bson.RegEx{regex, options}))
			author.SetMapIndex(reflect.ValueOf("regex"), reflect.Value{})
			author.SetMapIndex(reflect.ValueOf("options"), reflect.Value{})
			author.SetMapIndex(reflect.ValueOf("not"), reflect.Value{})
		}

		// tags
		tags := reflect.ValueOf(search["tags"])
		if tags.IsValid() {
			regex := tags.MapIndex(reflect.ValueOf("regex")).Elem().String()
			options := tags.MapIndex(reflect.ValueOf("options")).Elem().String()
			if tags.MapIndex(reflect.ValueOf("not")).Elem().Bool() {
				operator = "$not"
			}
			tags.SetMapIndex(reflect.ValueOf(operator), reflect.ValueOf(bson.RegEx{regex, options}))
			tags.SetMapIndex(reflect.ValueOf("regex"), reflect.Value{})
			tags.SetMapIndex(reflect.ValueOf("options"), reflect.Value{})
			tags.SetMapIndex(reflect.ValueOf("not"), reflect.Value{})
		}

		// Category
		category := reflect.ValueOf(search["category"])
		if category.IsValid() {
			regex := category.MapIndex(reflect.ValueOf("regex")).Elem().String()
			options := category.MapIndex(reflect.ValueOf("options")).Elem().String()
			if category.MapIndex(reflect.ValueOf("not")).Elem().Bool() {
				operator = "$not"
			}
			category.SetMapIndex(reflect.ValueOf(operator), reflect.ValueOf(bson.RegEx{regex, options}))
			category.SetMapIndex(reflect.ValueOf("regex"), reflect.Value{})
			category.SetMapIndex(reflect.ValueOf("options"), reflect.Value{})
			category.SetMapIndex(reflect.ValueOf("not"), reflect.Value{})
		}

		// title
		title := reflect.ValueOf(search["title"])
		if title.IsValid() {
			regex := title.MapIndex(reflect.ValueOf("regex")).Elem().String()
			options := title.MapIndex(reflect.ValueOf("options")).Elem().String()
			if title.MapIndex(reflect.ValueOf("not")).Elem().Bool() {
				operator = "$not"
			}
			title.SetMapIndex(reflect.ValueOf(operator), reflect.ValueOf(bson.RegEx{regex, options}))
			title.SetMapIndex(reflect.ValueOf("regex"), reflect.Value{})
			title.SetMapIndex(reflect.ValueOf("options"), reflect.Value{})
			title.SetMapIndex(reflect.ValueOf("not"), reflect.Value{})
		}

		// Set collection as item
		collection := dbSession.DB("tarzan").C("item")

		// Get limit, default limit is 100
		limit := 100
		if limit_int, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
			limit = limit_int
		}

		// Get parameters, default offset is 0
		offset := 0
		if offset_int, err := strconv.Atoi(c.QueryParam("offset")); err == nil {
			if offset_int > 0 {
				offset = offset_int * limit
			}
		}


		// Get total of results
		total, _ := collection.Find(search).Count()

		// Iterate all list
		/*find := collection.Find(search).Select(bson.M{
			"_id": true,
			"item_id": true,
			"url": true,
			"author": true,
			"title": true,
			"created": true,
			"category": true,
			"subscribed": true,
			"price": true,
			"sales": bson.M{"$slice": -1},
		}).Limit(limit).Skip(offset)*/

		// Sort
		sortby := c.QueryParam("sort")
		sortorder := c.QueryParam("order")
		sortorder_int := 1
		if sortorder == "desc" {
			sortorder_int = -1
		}

		sort := bson.M{}
		if sortby != "" {
			if sortby == "sales" {
				sortby = "sales.value"
			}
			sort[sortby] = sortorder_int
		} else {
			sort["created"] = -1
		}
		
		// Iterate all list
		now := time.Now().In(time.Local)
		format_date := "02/01/2006"
		start_date := int32(now.AddDate(0, 0, -8).Unix())
		end_date := int32(now.Unix())

		aggregate := collection.Pipe([]bson.M{
			bson.M{
				"$project": bson.M{
					"_id": 1,
					"item_id": 1,
					"url": 1,
					"author": 1,
					"title": 1,
					"category": 1,
					"price": 1,
					"created": 1,
					"subscribed": 1,
					"weeksales": bson.M{
						"$filter": bson.M{
							"input": "$sales",
							"as": "sales",
							"cond": bson.M{
								"$and": []bson.M{
									bson.M{"$gt": []interface{}{"$$sales.date", start_date}},
									bson.M{"$lte": []interface{}{"$$sales.date", end_date}},
								},
							},
						},
					},
					"sales": bson.M{"$slice": []interface{}{"$sales", -1}},
				},
			},
			bson.M{ "$match": search },
			bson.M{ "$unwind": "$sales" },
			bson.M{	"$sort": sort },
			bson.M{	"$skip": offset },
			bson.M{ "$limit": limit },
		})
		
		// Get items
		var result []Item
		err := aggregate.Iter().All(&result)

		// Count week sales
		for i := range result {
			sales_series_item := SalesSeries{result[i].Weeksales, result[i].Price}
			sales_series := get_sales_from_series(sales_series_item, tf_timezone_str, format_date)

			count := int32(0)
			for _, v := range sales_series {
				count += int32(reflect.ValueOf(v).MapIndex(reflect.ValueOf("sales")).Elem().Int())
			}
			result[i].WeeksalesData = count
		}

		// Send response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
				"total": total,
				"list": result,
			})
		} else {
			defaultResponse.Message = err.Error()
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}