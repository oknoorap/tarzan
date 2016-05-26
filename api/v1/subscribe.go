package v1

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
	"log"
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

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("group")

		// Iterate all list
		iterate := collection.Find(nil).Select(bson.M{
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
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
	})
}