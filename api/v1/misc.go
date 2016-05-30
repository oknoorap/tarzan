package v1

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
	"reflect"
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
		var result []string
		err := collection.Find(nil).Distinct("category", &result)

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
		find := collection.Find(search).Select(bson.M{
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
		}).Limit(limit).Skip(offset)
		
		// Get items
		var result []ItemViewSearch
		iterate := find.Iter()
		err := iterate.All(&result)

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