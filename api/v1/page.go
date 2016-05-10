package v1

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
	"log"
)


type (
	Error struct {
		Error bool `json:"error"`
		Message string `json:"message"`
	}

	Page struct {
		Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Url string `json:"url"`
		Title string `json:"title"`
		Desc string `json:"desc"`
	}
)

// Create new page
func PageNew (c echo.Context) error {

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
		url := c.FormValue("url")
		desc := c.FormValue("desc")
		title := c.FormValue("title")

		if url == "" || desc == "" || title == "" {
			defaultResponse.Message = "Unknown URL / Desc / Title"
			return c.JSON(http.StatusOK, defaultResponse)
		}

		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("page")

		// MongoDB data
		fields := bson.M{
			"$set": bson.M{
				"url": url,
				"title": title,
				"desc": desc,
				"time": int32(time.Now().Unix()),
			},
		}

		// Upsert `page` collection
		_, err := collection.Upsert(bson.M{"url": url}, fields)

		// Send Response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
			})
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}

// Get page by given id
func PageGet (c echo.Context) error {

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
		collection := dbSession.DB("tarzan").C("page")

		// Iterate all list
		var result Page
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



// Update page by given id
func PageUpdate (c echo.Context) error {
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
		title := c.FormValue("title")

		// Set Mgo Session
		db.SetMode(mgo.Monotonic, true)
		dbSession := db.Copy()
		defer dbSession.Close()

		// Pick MongoDB collection
		collection := dbSession.DB("tarzan").C("page")

		fields := bson.M{
			"$set": bson.M{
				"title": title,
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
			defaultResponse.Message = "Can't Update " + id
		}
	}

	return c.JSON(http.StatusOK, defaultResponse)
}


// Delete pageby given id
func PageDelete (c echo.Context) error {
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
		collection := dbSession.DB("tarzan").C("page")

		// Remove from db
		err := collection.Remove(bson.M{"_id": id})

		// Send Response
		if err == nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": http.StatusOK,
			})
		} else {
			defaultResponse.Message = "Id not exists"
		}

	}

	return c.JSON(http.StatusOK, defaultResponse)
}

// Show all pages
func PageList (c echo.Context) error {

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
		collection := dbSession.DB("tarzan").C("page")

		// Iterate all list
		iterate := collection.Find(nil).Select(bson.M{
			"_id": true,
			"url": true,
			"title": true,
			"desc": true,
		}).Limit(10000).Sort("-time").Iter()

		var result []Page
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