package v1

import(
	"gopkg.in/mgo.v2"
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