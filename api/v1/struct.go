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
		Created string `json:"created"`
		Sales []ItemSales `json:"sales"`
		Subscribed bool `json:"subscribed"`
	}

	ItemView struct {
		Item_id int32 `json:"item_id"`
		Url string `json:"url"`
		Author string `json:"author"`
		Title string `json:"title"`
		Price float32 `json:"price"`
		Created string `json:"created"`
		Category string `json:"category"`
		Tags []string `json:"tags"`
		Sales []ItemSales `json:"sales"`
		Subscribed bool `json:"subscribed"`
	}

	ItemViewSearch struct {
		Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Item_id int32 `json:"item_id"`
		Url string `json:"url"`
		Author string `json:"author"`
		Title string `json:"title"`
		Price float32 `json:"price"`
		Created string `json:"created"`
		Category string `json:"category"`
		Sales []ItemSales `json:"sales"`
		Subscribed bool `json:"subscribed"`
	}

	ItemSubscribe struct {
		Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Item_id int32 `json:"item_id"`
		Url string `json:"url"`
		Author string `json:"author"`
		Title string `json:"title"`
		Created string `json:"created"`
		Sales []ItemSales `json:"sales"`
		Subscribe_group_id []string `json:"subscribe_group_id"`
		Subscribed bool `json:"subscribed"`
	}

	ItemSales struct {
		Date int32 `json:"date"`
		Value int32 `json:"value"`
	}

	Sales struct {
		Sales int32   `json:"sales"`
		Price float32 `json:"price"`
		Date  string  `json:"date"`
	}

	SalesSeries struct {
		Sales []ItemSales `json:"sales"`
		Price float32     `json:"price"`
	}

	Group struct {
		Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Name string `json:"name"`
		Desc string `json:"desc"`
	}

	TagsCount struct {
		Label string `json:"label" bson:"_id"`
		Count int32 `json:"count"`
	}
)