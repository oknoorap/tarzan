package v1

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
	"net/http"
	"reflect"
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

		tf_timezone, _ := time.LoadLocation("Australia/Melbourne")
		local, _ := time.LoadLocation("America/New_York")
		date := c.QueryParam("date")
		now := time.Now().In(time.Local)
		end_date := int32(now.Unix())
		year, month, day := now.Date()

		if date == "today" {
			format_date = "02/01/2006 15:00 PM"
			start_date = int32(time.Date(year, month, day, 0, 0, 0, 0, local).Unix())
		} else if date == "week" {
			format_date = "02/01/2006"
			start_date = int32(now.AddDate(0, 0, -7).Unix())
		} else if date == "month" {
			format_date = "02/01/2006"
			_, m, _ := now.AddDate(0, -1, 0).Date()
			start_date = int32(time.Date(year, m, 1, 0, 0, 0, 0, local).Unix())
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
			limit = 10
			project_query["title"] = 1
			project_query["item_id"] = 1
			project_query["subscribed"] = 1
			sort = bson.M{
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
			}
		}


		match := bson.M{"$match": match_query}
		project := bson.M{"$project": project_query}

		// Iterate all list
		// db.item.aggregate([])
		aggregate := collection.Pipe([]bson.M{project, match, sort, bson.M{"$limit": limit}})
		data := make(map[string]interface{})

		if bestselling == "" {
			var result []SalesSeries
			err = aggregate.Iter().All(&result)

			if err == nil {
				for _, item := range result {
					itemData := []map[string]interface{}{}
					itemDataOf := reflect.ValueOf(itemData)
					itemSales := item.Sales
					itemPrice := item.Price
					priceMarket := make(map[string]interface{})

					for _, sales := range itemSales {
						dateIndex := time.Unix(int64(sales.Date), 0).In(tf_timezone).Format(format_date)
						if data[dateIndex] == nil {
							data[dateIndex] = map[string]interface{}{
								"sales": int32(0),
								"price": float32(0),
							}
						}

						var itemDataFound reflect.Value
						itemDataOf = reflect.ValueOf(itemData)
						itemExists := false
						for i := 0; i < itemDataOf.Len(); i++ {
							currentItem := itemDataOf.Index(i)
							itemDate := currentItem.MapIndex(reflect.ValueOf("date")).Elem().String()
							if itemDate == dateIndex {
								itemExists = true
								itemDataFound = currentItem
							}
						}

						if itemExists {
							itemDataFound.SetMapIndex(reflect.ValueOf("sales"), reflect.ValueOf(sales.Value))
						} else {
							itemData = append(itemData, map[string]interface{}{
								"sales": sales.Value,
								"date":  dateIndex,
							})
						}


						priceMarket[dateIndex] = itemPrice
					}

					for i := 0; i < itemDataOf.Len(); i++ {
						if i > 0 {
							currentItem := itemDataOf.Index(i)
							currentSales := currentItem.MapIndex(reflect.ValueOf("sales")).Elem().Int()
							currentDate := currentItem.MapIndex(reflect.ValueOf("date")).Elem().String()

							previousItem := itemDataOf.Index(i - 1)
							previousSales := previousItem.MapIndex(reflect.ValueOf("sales")).Elem().Int()
							
							sumSales := currentSales - previousSales

							abc := reflect.ValueOf(data[currentDate])
							datasales := abc.MapIndex(reflect.ValueOf("sales")).Elem().Int()
							sum := datasales + sumSales
							abc.SetMapIndex(reflect.ValueOf("sales"), reflect.ValueOf(sum))
						}
					}

					/*itemValue := reflect.ValueOf(itemData)
					for i := 0; i < itemValue.Len(); i++ {
						if i > 0 {
							currentSales := itemValue.Index(i).MapIndex(reflect.ValueOf("sales")).Int()
							previousSales := itemValue.Index(i - 1).MapIndex(reflect.ValueOf("sales")).Int()
							sumSales := currentSales - previousSales
							itemValue.SetMapIndex(reflect.ValueOf("sales"), reflect.ValueOf(sumSales))
						}
					}*/

					//log.Println(itemData)

					for date, price := range priceMarket {
						marketInDate := reflect.ValueOf(data[date])
						if marketInDate.IsValid() {
							priceVal := float32(reflect.ValueOf(price).Float())
							es := marketInDate.MapIndex(reflect.ValueOf("price"))
							if es.IsValid() {
								var abc float32
								if es.Elem().Kind().String() == "int" {
									abc = float32(es.Elem().Int())
								} else {
									abc = float32(es.Elem().Float())
								}

								sumPrice := priceVal + abc
								marketInDate.SetMapIndex(reflect.ValueOf("price"), reflect.ValueOf(sumPrice))
							}
						}
					}
				}

				return c.JSON(http.StatusOK, map[string]interface{}{
					"status": http.StatusOK,
					"data": data,
				})
			} else {
				defaultResponse.Message = err.Error()
			}
		} else {
			var result []map[string]interface{}
			err = aggregate.Iter().All(&result)
			log.Println(result)

			if err == nil {
				for _, value := range result {
					valueOf := reflect.ValueOf(value)
					item_id := strconv.Itoa(int(valueOf.MapIndex(reflect.ValueOf("item_id")).Elem().Int()))
					title := valueOf.MapIndex(reflect.ValueOf("title")).Elem().String()

					var priceval float32
					price := valueOf.MapIndex(reflect.ValueOf("price"))
					if price.IsValid() {
						if price.Elem().Kind().String() == "int" {
							priceval = float32(price.Elem().Int())
						} else {
							priceval = float32(price.Elem().Float())
						}
					}
					data[item_id] = map[string]interface{}{
						"title": title,
						"price": priceval,
					}
				}

				return c.JSON(http.StatusOK, map[string]interface{}{
					"status": http.StatusOK,
					"data": data,
				})
			} else {
				defaultResponse.Message = err.Error()
			}
			defaultResponse.Message = "hahahhahaha"
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