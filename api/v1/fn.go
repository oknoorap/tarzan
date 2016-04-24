package v1

import(
	"time"
	"gopkg.in/mgo.v2"
)

func connectDb () (*mgo.Session, error) {
	return mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: []string{"localhost:27017"},
		Timeout: 60 * time.Second,
	})
}