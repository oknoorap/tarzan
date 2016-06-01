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

		loc, _ := time.LoadLocation("Australia/Sydney")
		date := c.QueryParam("date")
		now := time.Now().In(time.Local)
		end_date := int32(now.Unix())
		year, month, day := now.Date()

		if date == "today" {
			format_date = "02/01/2006 15:00 PM"
			start_date = int32(time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix())
		} else if date == "week" {
			format_date = "02/01/2006"
			start_date = int32(now.AddDate(0, 0, -7).Unix())
		} else if date == "month" {
			format_date = "02/01/2006"
			_, m, _ := now.AddDate(0, -1, 0).Date()
			start_date = int32(time.Date(year, m, 1, 0, 0, 0, 0, time.Local).Unix())
		} else if date == "lastmonth" {
			format_date = "02/01/2006"
			_, m, _ := now.AddDate(0, -1, 0).Date()
			start := time.Date(year, m, 1, 0, 0, 0, 0, time.UTC)
			start_date = int32(start.Unix())
			end_date = int32(time.Date(year, m, start.Day(), 0, 0, 0, 0, time.Local).Unix())
		} else if date == "year" {
			format_date = "01/2006"
			start_date = int32(time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Unix())
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

		match := bson.M{"$match": match_query}
		project := bson.M{"$project": project_query}

		// Iterate all list
		// db.item.aggregate([])
		aggregate := collection.Pipe([]bson.M{project, match, bson.M{"$limit": 50000}})


		var result []MarketValueSeries
		err = aggregate.Iter().All(&result)

		// Send response
		if err == nil {

			data := make(map[string]interface{})

			for _, item := range result {
				itemSales := item.Sales
				itemPrice := item.Price
				priceMarket := make(map[string]interface{})
				lastSales := int32(0)

				if len(itemSales) > 0 {
					lastSales = itemSales[len(itemSales)-1].Value
				}

				for index, sales := range itemSales {
					dateIndex := time.Unix(int64(sales.Date), 0).In(loc).Format(format_date)
					dateIndex2 := time.Unix(int64(sales.Date), 0).In(time.Local).Format(format_date)
					salesValue := sales.Value
					log.Println(sales.Date, dateIndex, dateIndex2, sales.Value)

					if data[dateIndex] == nil {
						data[dateIndex] = map[string]interface{}{
							"sales": int32(0),
							"price": float32(0),
						}
					}

					date := reflect.ValueOf(data[dateIndex])
					if date.IsValid() {
						if index > 0 {
							salesValue = salesValue - itemSales[index - 1].Value
						} else {
							salesValue = lastSales - salesValue
						}
						sumSales := salesValue + int32(date.MapIndex(reflect.ValueOf("sales")).Elem().Int())
						date.SetMapIndex(reflect.ValueOf("sales"), reflect.ValueOf(sumSales))
					}

					if priceMarket[dateIndex] == nil {
						priceMarket[dateIndex] = 0
					}

					priceMarket[dateIndex] = itemPrice
				}

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