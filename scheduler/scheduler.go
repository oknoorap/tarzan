package main

import (
	//"fmt"
	"flag"
	"time"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/redis.v3"
)

var (
	interval = flag.Int("interval", 5, "Cron Interval in Minute. Default is 5 hour.")
)

func main() {
	// Parsing Country Code
	flag.Parse()

	// Check first time
	go check()

	// Check daily
	go cron(time.Duration(*interval) * time.Hour)

	// Tick forever
	select {}
}


func cron (t time.Duration) {
	for _ = range time.Tick(t) {
		go check()
	}
}


func connectDb () (*mgo.Session, error) {
	return mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: []string{"localhost:27017"},
		Timeout: 60 * time.Second,
	})
}

func dialRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := client.Ping().Result()
	return client, err
}


type PageUrl struct {
	Id bson.ObjectId
	Url string
}

func check () {
	// Connect to mongodb
	db, err := connectDb()

	if err != nil {

		log.Fatalf("CreateSession: %s\n", err)

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
		}).Limit(10000).Sort("-time").Iter()

		var result []PageUrl
		err := iterate.All(&result)

		if err == nil {
			// Connect redis
			redis, err := dialRedis()
			if err != nil {
				log.Println(err)
			}

			for _, item := range result {
		    	err := redis.RPush("resque:queue:main", `{"class":"GetPageItems","args":[{"url":"`+ item.Url +`"}]}`).Err()

				if err != nil {
					log.Println(err)
				}
			}
		} else {
			log.Println(err)
		}
	}
}