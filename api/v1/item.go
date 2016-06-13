package v1

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"log"
	"time"
	"reflect"
	"net/http"
)


// Get item by given id
func ItemGet (c echo.Context) error {
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

		// Get ID
		id := bson.ObjectIdHex(c.Param("id"))

		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("item")

		// Iterate all list
		var result ItemView
		err := collection.FindId(id).Select(bson.M{
			"author": true,
			"category": true,
			"created": true,
			"item_id": true,
			"price": true,
			"sales": bson.M{"$slice": -1},
			"subscribed": true,
			"tags": true,
			"title": true,
			"url": true,
		}).One(&result)

		// Send response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
				"data": result,
			})
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}


func ItemUpdate (c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

// Show all items
func ItemList (c echo.Context) error {
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

		// Count all item
		total, _ := collection.Find(nil).Count()

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
			bson.M{ "$unwind": "$sales" },
			bson.M{	"$sort": sort },
			bson.M{	"$skip": offset },
			bson.M{ "$limit": limit },
		})
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
