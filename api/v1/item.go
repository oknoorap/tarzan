package v1

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)


// Get item by given id
func ItemGet (c echo.Context) error {
	return c.JSON(http.StatusOK, "")
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

		// Iterate all list
		iterate := collection.Find(nil).Select(bson.M{
			"_id": true,
			"url": true,
			"author": true,
			"title": true,
			"sales": bson.M{"$slice": -1},
		}).Limit(1000).Sort("-time").Iter()

		var result []Item
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
