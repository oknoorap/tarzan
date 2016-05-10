package task

import(
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/redis.v3"
)

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