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
		Item_id int32 `json:"item_id"`
		Url string `json:"url"`
		Author string `json:"author"`
		Title string `json:"title"`
		Sales []ItemSales `json:"sales"`
	}

	ItemView struct {
		Item_id int32 `json:"item_id"`
		Url string `json:"url"`
		Author string `json:"author"`
		Title string `json:"title"`
		Price string `json:"price"`
		Created string `json:"created"`
		Category string `json:"category"`
		Tags []string `json:"tags"`
		Sales []ItemSales `json:"sales"`
	}

	ItemSales struct {
		Date int32 `json:"date"`
		Value int32 `json:"value"`
	}
)