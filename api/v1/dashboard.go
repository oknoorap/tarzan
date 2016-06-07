package v1

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
	sorting "sort"
	"reflect"
	"net/http"
	"strconv"
)


// Get item by given id
func MarketValue (c echo.Context) error {
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

		// Get default limit by date
		var (
			start_date int32
			format_date string
		)

		tf_timezone_str := "Australia/Melbourne"
		local_str := "America/New_York"
		local, _ := time.LoadLocation(local_str)
		date := c.QueryParam("date")
		now := time.Now().In(time.Local)
		end_date := int32(now.Unix())
		year, month, day := now.Date()

		if date == "today" {
			format_date = "02/01/2006 15:00 PM"
			start_date = int32(time.Date(year, month, day, 0, 0, 0, 0, local).Unix())
		} else if date == "week" {
			format_date = "02/01/2006"
			start_date = int32(now.AddDate(0, 0, -8).Unix())
		} else if date == "month" {
			format_date = "02/01/2006"
			_, m, _ := now.AddDate(0, -1, 0).Date()
			start_date = int32(time.Date(year, m, 1, -1, 0, 0, 0, local).Unix())
		} else if date == "lastmonth" {
			format_date = "02/01/2006"
			_, m, _ := now.AddDate(0, -1, 0).Date()
			start := time.Date(year, m, 1, 0, 0, 0, 0, local)
			start_date = int32(start.Unix())
			end_date = int32(time.Date(year, m, start.Day(), 0, 0, 0, 0, local).Unix())
		} else if date == "year" {
			format_date = "01/2006"
			start_date = int32(time.Date(year, time.January, 1, 0, 0, 0, 0, local).Unix())
		}

		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("item")

		// Build query
		project_query := bson.M{}
		match_query := bson.M{}
		limit := 50000

		// Select basic field
		project_query["price"] = 1
		project_query["sales"] = bson.M{
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
		}

		// If request contains category
		category := c.QueryParam("category")
		if category != "" {
			match_query["category"] = bson.M{
				"$regex": bson.RegEx{category, ""},
			}
			project_query["category"] = 1
		}

		// If request contains group
		group_id := c.QueryParam("group")
		if group_id != "" {
			match_query["subscribe_group_id"] = bson.M{
				"$in": []string{group_id},
			}
			project_query["subscribe_group_id"] = 1
		}

		// If request contains item_id
		item_id, err := strconv.Atoi(c.QueryParam("item_id"))
		if err == nil {
			match_query["item_id"] = item_id
			project_query["item_id"] = 1
		}


		// If request contains category
		bestselling := c.QueryParam("bestselling")
		sort := bson.M{"$sort": bson.M{"created": 1}}
		if bestselling != "" {
			project_query["_id"] = 1
			project_query["title"] = 1
			project_query["item_id"] = 1
			project_query["subscribed"] = 1
			project_query["url"] = 1
			project_query["img_preview"] = 1
			/*sort = bson.M{
				"$sort": bson.M{"sales": -1},
			}

			//db.item.aggregate([{$project: {sales: {$arrayElemAt: [{$filter: {input: "$sales", as: "sales", cond: {$gte: ["$$sales.value", 20]} } }, -1]}}  }, {$sort: {sales: -1} } ])
			project_query["sales"] = bson.M{
				"$arrayElemAt": []interface{}{
					bson.M{
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
					-1,
				},
			}*/
		}


		match := bson.M{"$match": match_query}
		project := bson.M{"$project": project_query}

		// Iterate all list
		// db.item.aggregate([])
		aggregate := collection.Pipe([]bson.M{project, match, sort, bson.M{"$limit": limit}})

		if bestselling == "" {
			var result []SalesSeries
			var data []map[string]interface{}
			err = aggregate.Iter().All(&result)

			if err == nil {

				// Get sales series from items
				for _, item := range result {
					series := get_sales_from_series(item, tf_timezone_str, format_date)
					data = append(data, series)
				}

				// Sum all sales
				counted_data := sum_sales(data)

				return c.JSON(http.StatusOK, ResponseOK{http.StatusOK, counted_data})
			} else {
				defaultResponse.Message = err.Error()
			}

		} else {
			var result []map[string]interface{}
			data := make(map[string]interface{})
			err = aggregate.Iter().All(&result)

			if err == nil {

				for _, item := range result {
					var price_value float32
					price_reflector := reflect.ValueOf(item["price"])
					if price_reflector.Kind().String() == "float64" {
						price_value = float32(price_reflector.Float())
					} else {
						price_value = float32(price_reflector.Int())
					}

					var item_sales []ItemSales
					sales_collection := reflect.ValueOf(item["sales"])
					for i := 0; i < sales_collection.Len(); i++ {
						var sales_date int32
						sales_date_reflector := sales_collection.Index(i).Elem().MapIndex(reflect.ValueOf("date"))
						if sales_date_reflector.Elem().Kind().String() == "int" {
							sales_date = int32(sales_date_reflector.Elem().Int())
						} else {
							sales_date = int32(sales_date_reflector.Elem().Float())
						}

						var sales_value int32
						sales_value_reflector := sales_collection.Index(i).Elem().MapIndex(reflect.ValueOf("value"))
						if sales_value_reflector.Elem().Kind().String() == "int" {
							sales_value = int32(sales_value_reflector.Elem().Int())
						} else {
							sales_value = int32(sales_value_reflector.Elem().Float())
						}
						item_sales = append(item_sales, ItemSales{sales_date, sales_value})
					}


					sales_series := get_sales_from_series(SalesSeries{item_sales, price_value}, tf_timezone_str, format_date)
					sales_total := int64(0)
					for _, sales := range sales_series {
						sales_reflect := reflect.ValueOf(sales)
						sales_value := sales_reflect.MapIndex(reflect.ValueOf("sales")).Elem().Int()
						sales_total += sales_value
					}
					
					if sales_total > 0 {
						key := strconv.Itoa(item["item_id"].(int))
						data[key] = map[string]interface{}{
							"id": item["_id"],
							"item_id": item["item_id"],
							"price": item["price"],
							"url": item["url"],
							"img_preview": item["img_preview"],
							"title": item["title"],
							"total": sales_total,
						}
					}
				}

				// Sorting by total value
				sort_total := make(map[int][]string)
				for key, v := range data {
					val := int(reflect.ValueOf(v).MapIndex(reflect.ValueOf("total")).Elem().Int())
					sort_total[val] = append(sort_total[val], key)
				}

				var sort_keys []int
				for key, _ := range sort_total {
					sort_keys = append(sort_keys, key)
				}
				sorting.Sort(sorting.Reverse(sorting.IntSlice(sort_keys)))

				new_data := make(map[string]interface{})
				var sort_slice []int
				if len(sort_keys) > 9 {
					sort_slice = sort_keys[0:10]
				} else {
					sort_slice = sort_keys
				}
				for _, key := range sort_slice {
					for _, val := range sort_total[key] {
						if reflect.ValueOf(new_data).Len() < 15 {
							new_data[val] = data[val]
						}
					}
				}

				return c.JSON(http.StatusOK, map[string]interface{}{
					"status": http.StatusOK,
					"data": new_data,
				})
			} else {
				defaultResponse.Message = err.Error()
			}
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}

// Get item by given id
func Tags (c echo.Context) error {
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
		aggregate := collection.Pipe([]bson.M{
			bson.M{ "$unwind": "$tags"},
			bson.M{
				"$group": bson.M{
					"_id": "$tags",
					"count": bson.M{"$sum": 1},
				},
			},
			bson.M{ "$sort": bson.M{"count": -1}},
			bson.M{ "$limit": 10},
		})
		var result []TagsCount
		err := aggregate.Iter().All(&result)

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