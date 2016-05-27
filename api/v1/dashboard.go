package v1

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
	"net/http"
	"reflect"
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

		date := c.QueryParam("date")
		now := time.Now()
		end_date := int32(now.Unix())
		year, _, _ := now.Date()

		if date == "today" {
			format_date = "02/01/2006 15:00 PM"
			start_date = int32(now.AddDate(0, 0, -1).Unix())
		} else if date == "week" {
			format_date = "02/01/2006"
			start_date = int32(now.AddDate(0, 0, -8).Unix())
		} else if date == "month" {
			format_date = "02/01/2006"
			_, month, _ := now.AddDate(0, -1, 0).Date()
			start_date = int32(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Unix())
		} else if date == "lastmonth" {
			format_date = "02/01/2006"
			_, month, _ := now.AddDate(0, -1, 0).Date()
			start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
			start_date = int32(start.Unix())
			end_date = int32(time.Date(year, month, start.Day(), 0, 0, 0, 0, time.UTC).Unix())
		} else if date == "year" {
			format_date = "01/2006"
			start_date = int32(time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC).Unix())
		}

		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("item")

		// Iterate all list
		// db.item.aggregate([])
		aggregate := collection.Pipe([]bson.M{
			bson.M{
				"$project": bson.M{
					"sales": bson.M{
						"$filter": bson.M{
							"input": "$sales",
							"as": "sales",
							"cond": bson.M{
								"$and": []bson.M{
									bson.M{"$gte": []interface{}{"$$sales.date", start_date}},
									bson.M{"$lte": []interface{}{"$$sales.date", end_date}},
								},
							},
						},
					},
					"price": 1,
				},
			},
			bson.M{"$limit": 50000},
		})

		var result []MarketValueSeries
		err := aggregate.Iter().All(&result)

		// Send response
		if err == nil {

			data := make(map[string]interface{})

			for _, item := range result {
				itemSales := item.Sales
				itemPrice := item.Price
				priceMarket := make(map[string]interface{})

				for index, sales := range itemSales {
					dateIndex := time.Unix(int64(sales.Date), 0).Format(format_date)
					salesValue := sales.Value
					
					if index > 0 {
						salesValue = salesValue - itemSales[index - 1].Value
					}

					if data[dateIndex] == nil {
						data[dateIndex] = map[string]interface{}{
							"sales": int32(0),
							"price": float32(0),
						}
					}

					date := reflect.ValueOf(data[dateIndex])
					if date.IsValid() {
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

			/*_.each(data, function (item) {

				// Get the same time of document,
				// Prevent duplicate unixtime of timeline when scraping
				_.each(item.sales, function (sales, index) {
					var date = moment.unix(sales.date).format(dateFormat), salesValue = sales.value

					if (index > 0) salesValue = salesValue - item.sales[index - 1].value
					if (!marketValue[date]) marketValue[date] = {sales: 0, price: 0}
					marketValue[date].sales += parseInt(salesValue)
					marketValue[date].price += parseInt(item.price)
				})
			})*/

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