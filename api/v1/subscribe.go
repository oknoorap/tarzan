package v1

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
	"log"
	"reflect"
	"net/http"
)

func subunsub (subscribe bool, c echo.Context) error {

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

		// Get group ID
		group_id := c.FormValue("group_id")

		// Convert item_id to integer
		item_id := 0
		if item_id_int, err := strconv.Atoi(c.FormValue("item_id")); err == nil {
			item_id = item_id_int
		} else {
			defaultResponse.Message = err.Error()
		}

		// MongoDB data
		fields := bson.M{}

		if subscribe {
			fields = bson.M{
				"$addToSet": bson.M{
					"subscribe_group_id": group_id,
				},
				"$set": bson.M{
					"subscribed": subscribe,
				},
			}
		} else {
			fields = bson.M{
				"$set": bson.M{
					"subscribe_group_id": []string{},
					"subscribed": subscribe,
				},
			}
		}

		// Update `item` collection
		info, err := collection.Upsert(bson.M{"item_id": item_id}, fields)

		// Send Response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
				"info": info,
			})
		} else {
			defaultResponse.Message = err.Error()
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}


func Subscribe (c echo.Context) error {
	return subunsub(true, c)
}


func Unsubscribe (c echo.Context) error {
	return subunsub(false, c)
}


func GroupNew (c echo.Context) error {
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

		// Get POST data
		name := c.FormValue("name")
		desc := c.FormValue("desc")

		if name == "" || desc == "" {
			defaultResponse.Message = "Unknown Name and Desc"
			return c.JSON(http.StatusOK, defaultResponse)
		}

		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("group")

		// MongoDB data
		fields := bson.M{
			"$set": bson.M{
				"name": name,
				"desc": desc,
				"time": int32(time.Now().Unix()),
			},
		}

		// Upsert `page` collection
		_, err := collection.Upsert(bson.M{"name": name}, fields)

		// Send Response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
			})
		} else {
			defaultResponse.Message = err.Error()
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}

func GroupGet (c echo.Context) error {
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
		collection := dbSession.DB("tarzan").C("group")

		// Iterate all list
		var result Group
		err := collection.FindId(id).One(&result)

		// Send response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
				"data": result,
			})
		} else {
			defaultResponse.Message = err.Error()
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}

func GroupUpdate (c echo.Context) error {
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
		id := c.Param("id")
		desc := c.FormValue("desc")

		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("group")

		fields := bson.M{
			"$set": bson.M{
				"desc": desc,
			},
		}
		_, err := collection.Upsert(bson.M{"_id": bson.ObjectIdHex(id)}, fields)

		// Send response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
			})
		} else {
			defaultResponse.Message = err.Error()
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}

func GroupDelete (c echo.Context) error {
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

		// Get DELETE data
		id := bson.ObjectIdHex(c.Param("id"))

		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("group")

		// Remove from db
		err := collection.Remove(bson.M{"_id": id})

		// Send Response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
			})
		} else {
			defaultResponse.Message = err.Error()
		}

	}

	return c.JSON(http.StatusOK, defaultResponse)
}


func GroupList (c echo.Context) error {
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

		// Query Fields
		query := bson.M{}

		// Get Category ID
		cat_id := c.QueryParam("cat_id")
		if cat_id != "" {
			query["_id"] = bson.ObjectIdHex(cat_id)
		}

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("group")

		// Iterate all list
		iterate := collection.Find(query).Select(bson.M{
			"_id": true,
			"name": true,
			"desc": true,
		}).Limit(10000).Sort("-time").Iter()

		var result []Group
		err := iterate.All(&result)

		// Send response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
				"list": result,
			})
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}


func SubscribeList (c echo.Context) error {
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

		// Iterate all list
		in := c.QueryParam("in")
		in_array := []interface{}{in}
		aggregate := collection.Pipe([]bson.M{
			bson.M{
				"$project": bson.M{
					"_id": 1,
					"item_id": 1,
					"url": 1,
					"author": 1,
					"price": 1,
					"category": 1,
					"title": 1,
					"created": 1,
					"subscribed": 1,
					"subscribe_group_id": 1,
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
					"sales": bson.M{
						"$slice": []interface{}{"$sales", -1},
					},
				},
			},
			bson.M{
				"$match": bson.M{
					"subscribed": true,
					"subscribe_group_id": bson.M{"$in": in_array},
				},
			},
			bson.M{"$sort": sort},
			bson.M{"$unwind": "$sales"},
			bson.M{"$limit": 500},
		})

		var result []ItemSubscribe
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
				"list": result,
			})
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}