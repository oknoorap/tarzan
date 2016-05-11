package v1

import (
	"gopkg.in/mgo.v2/bson"
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

	Item struct {
		Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Url string `json:"url"`
		Author string `json:"author"`
		Title string `json:"title"`
		Sales []ItemSales `json:"sales"`
	}

	ItemSales struct {
		Date int32 `json:"date"`
		Value int32 `json:"value"`
	}
)