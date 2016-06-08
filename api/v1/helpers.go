package v1

import(
	"gopkg.in/mgo.v2"
	"reflect"
	"sort"
	"time"
	"log"
)

type TimeSlice []time.Time
func (p TimeSlice) Len() int {
	return len(p)
}

func (p TimeSlice) Less(i, j int) bool {
	return p[i].Before(p[j])
}

func (p TimeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

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

func sum_sales(sales_series []map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	for _, value := range sales_series {
		for date, item := range value {
			if data[date] == nil {
				data[date] = map[string]interface{}{
					"sales": int32(0),
					"price": float32(0),
				}
			}

			data_in_date := reflect.ValueOf(data[date])
			prev_sales := data_in_date.MapIndex(reflect.ValueOf("sales")).Elem().Int()
			prev_price := data_in_date.MapIndex(reflect.ValueOf("price")).Elem().Float()

			current_item := reflect.ValueOf(item)
			current_sales := current_item.MapIndex(reflect.ValueOf("sales")).Elem().Int()
			current_price := current_item.MapIndex(reflect.ValueOf("price")).Elem().Float()

			sum_sales := prev_sales + current_sales
			data_in_date.SetMapIndex(reflect.ValueOf("sales"), reflect.ValueOf(sum_sales))

			if current_sales > 0 {
				sum_price := prev_price + current_price
				data_in_date.SetMapIndex(reflect.ValueOf("price"), reflect.ValueOf(sum_price))
			}
		}
	}

	return data
}


// https://play.golang.org/p/zdbg0NvTun
func get_sales_from_series(item SalesSeries, time_location string, format_date string) map[string]interface{} {
	tf_timezone, _ := time.LoadLocation(time_location)
	data := make(map[string]interface{})

	for _, sales := range item.Sales {
		date_index := time.Unix(int64(sales.Date), 0).In(tf_timezone).Format(format_date)
		salesValue := sales.Value
		data[date_index] = map[string]interface{}{
			"sales": salesValue,
			"price": float32(item.Price),
		}
	}

	data_keys := []time.Time{}
	for key, _ := range data {
		a, _ := time.Parse(format_date, key)
		data_keys = append(data_keys, a)
	}
	sort.Sort(TimeSlice(data_keys))

	counted_sales := map[string]int32{}
	for i := 0; i < len(data_keys); i++ {
		key := data_keys[i].Format(format_date)
		item := reflect.ValueOf(data[key])
		sales := int32(item.MapIndex(reflect.ValueOf("sales")).Elem().Int())

		if i > 0 {
			prev_key := data_keys[i-1].Format(format_date)
			prev_item := reflect.ValueOf(data[prev_key])
			prev_sales := int32(prev_item.MapIndex(reflect.ValueOf("sales")).Elem().Int())
			count_sales := sales - prev_sales
			counted_sales[key] = count_sales
			
		}
	}

	for date, value := range counted_sales {
		item := reflect.ValueOf(data[date])
		if item.IsValid() {
			price_reflector := item.MapIndex(reflect.ValueOf("price")).Elem()
			price := float32(0)
			if price_reflector.Kind().String() == "float32" {
				price = float32(price_reflector.Float())
			} else {
				price = float32(price_reflector.Int())
			}
			total_sales := price * float32(value)
			item.SetMapIndex(reflect.ValueOf("sales"), reflect.ValueOf(value))
			item.SetMapIndex(reflect.ValueOf("price"), reflect.ValueOf(total_sales))
		}
	}
	
	if len(data_keys) > 0 {
		first_data := reflect.ValueOf(data[data_keys[0].Format(format_date)])
		first_data.SetMapIndex(reflect.ValueOf("sales"), reflect.ValueOf(0))
	}

	return data
}