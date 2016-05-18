package v1

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"log"
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
		err := collection.FindId(id).One(&result)

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

		// Iterate all list
		iterate := collection.Find(nil).Select(bson.M{
			"_id": true,
			"item_id": true,
			"url": true,
			"author": true,
			"title": true,
			"sales": bson.M{"$slice": -1},
		}).Limit(limit).Sort("-item_id").Skip(offset).Iter()

		var result []Item
		err := iterate.All(&result)

		// Send response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
				"total": total,
				"list": result,
			})
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}
