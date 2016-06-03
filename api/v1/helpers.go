package v1

import(
	"gopkg.in/mgo.v2"
	"reflect"
	"time"
	"log"
)

func connectDb () (*mgo.Session, error) {
	return mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: []string{"localhost:27017"},
		Timeout: 60 * time.Second,
	})
}

func logPanic (err error) {
	log.Println(err.Error())
	panic(err)
}

func get_sales_from_series(sales_series []MarketValueSeries) {
	for item := range sales_series {
		item_data := []map[string]interface{}{}
		price_market := make(map[string]interface{})
		item_data_reflect := reflect.ValueOf(item_data)

		for sales := range item.Sales {
			date_index := time.Unix(int64(sales.Date), 0).In(tf_timezone).Format(format_date)
			if data[date_index] == nil {
				data[date_index] = map[string]interface{}{
					"sales": int32(0),
					"price": float32(0),
				}
			}

			var item_data_found reflect.Value
			item_data_reflect = reflect.ValueOf(item_data)
			item_exists := false

			for i := 0; i < item_data_reflect.Len(); i++ {
				current_item := item_data_reflect.Index(i)
				item_date := current_item.MapIndex(reflect.ValueOf("date")).Elem().String()
				if item_date == date_index {
					item_exists = true
					item_data_found = current_item
				}
			}

			if item_exists {
				item_data_found.SetMapIndex(reflect.ValueOf("sales"), reflect.ValueOf(sales.Value))
			} else {
				item_data = append(item_data, map[string]interface{}{
					"sales": sales.Value,
					"date": date_index,
				})
			}

			price_market[date_index] = item.Price
		}

		for i := 0; i < item_data_reflect.Len(); i++ {
			if i > 0 {
				current_data := item_data_reflect.Index(i)
				current_sales := current_data.MapIndex(reflect.ValueOf("sales")).Elem().Int()
				current_date := current_data.MapIndex(reflect.ValueOf("date")).Elem().String()

				previous_item := item_data_reflect.Index(i - 1)
				previous_sales := previous_item.MapIndex(reflect.ValueOf("sales")).Elem().Int()
				sum_sales := current_sales - previous_sales

				abc := reflect.ValueOf(data[current_date])
				data_sales := abc.MapIndex(reflect.ValueOf("sales")).Elem().Int()
				sum := data_sales + sum_sales
				abc.SetMapIndex(reflect.ValueOf("sales"), reflect.ValueOf(sum))
			}
		}
	}
}