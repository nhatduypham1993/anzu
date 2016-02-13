package components

import (
	"github.com/fernandez14/spartangeek-blacker/modules/exceptions"
	"gopkg.in/mgo.v2/bson"
	"math"
	"time"
)

// Set DI instance
func (component *ComponentModel) SetDI(di *Module) {
	component.di = di
}

// Set generic data for the component model
func (component *ComponentModel) SetGeneric(data []byte) {
	component.generic = data
}

// Get generic data
func (component *ComponentModel) GetData() map[string]interface{} {

	var data map[string]interface{}

	err := bson.Unmarshal(component.generic, &data)

	if err != nil {
		panic(err)
	}

	return data
}

// Get Store price for vendor
func (component *ComponentModel) GetVendorPrice(vendor string) (float64, error) {

	if v, exists := component.Store.Vendors[vendor]; exists {

		return v.Price, nil
	}

	return float64(math.NaN()), exceptions.NotFound{"Invalid vendor. Not found."}
}

// Get Vendor object
func (component *ComponentModel) GetVendor(vendor string) (ComponentStoreItemModel, error) {

	if v, exists := component.Store.Vendors[vendor]; exists {

		return v, nil
	}

	return ComponentStoreItemModel{}, exceptions.NotFound{"Invalid vendor. Not found."}
}

// Update component's price
func (component *ComponentModel) UpdatePrice(vendors map[string]map[string]interface{}) {

	database := component.di.Mongo.Database

	// Record price history
	if len(component.Store.Vendors) > 0 {

		historic := &ComponentHistoricModel{
			ComponentId: component.Id,
			Store:       component.Store,
			Created:     time.Now(),
		}

		err := database.C("components_historic").Insert(historic)

		if err != nil {
			panic(err)
		}
	} else {

		// Allocate the map
		component.Store.Vendors = make(map[string]ComponentStoreItemModel)
	}

	set := bson.M{"store.updated_at": time.Now(), "activated": true}

	for key, item := range vendors {

		price, with_price := item["price"]
		stock, with_stock := item["stock"]
		prio, with_prio := item["priority"]

		if with_price && with_stock && with_prio {

			u := ComponentStoreItemModel{
				Price:    price.(float64),
				Stock:    int(stock.(float64)),
				Priority: int(prio.(float64)),
			}

			set["store.vendors."+key] = u

			// Runtime update
			component.Store.Vendors[key] = u
		}
	}

	err := database.C("components").Update(bson.M{"_id": component.Id}, bson.M{"$set": set})

	if err != nil {
		panic(err)
	}

	component.Activated = true

	go component.UpdateAlgolia()
}

func (component *ComponentModel) DeletePrice() {

	database := component.di.Mongo.Database

	// Record price history
	if len(component.Store.Vendors) > 0 {

		historic := &ComponentHistoricModel{
			ComponentId: component.Id,
			Store:       component.Store,
			Created:     time.Now(),
		}

		err := database.C("components_historic").Insert(historic)

		if err != nil {
			panic(err)
		}
	}

	component.Activated = false
	component.Store.Vendors = make(map[string]ComponentStoreItemModel)

	err := database.C("components").Update(bson.M{"_id": component.Id}, bson.M{"$unset": bson.M{"store": 1}, "$set": bson.M{"activated": false}})	

	if err != nil {
		panic(err)
	}

	go component.UpdateAlgolia()
}

func (component *ComponentModel) GetAggregatedUsrVotes(kind string) map[string]int {

	var votes []CommentVotesModel

	aggregated := make(map[string]int)
	database := component.di.Mongo.Database

	pipe := database.C("user_owns").Pipe([]bson.M{
		{"$match": bson.M{"related": kind, "related_id": component.Id, "removed": bson.M{"$exists": false}}},
		{"$group": bson.M{"_id": "$type", "count": bson.M{"$sum": 1}}},
	})

	err := pipe.All(&votes)

	if err != nil {
		panic(err)
	}

	total := 0

	if len(votes) > 0 {
		for _, vote := range votes {
			aggregated[vote.Id] = vote.Count
			total = total + vote.Count
		}
	}

	if kind == "component" {
		for _, k := range []string{"have-it", "had-it", "want-it"} {
			if _, exists := aggregated[k]; !exists {
				aggregated[k] = 0
			}
		}
	}

	if kind == "component-buy" {
		for _, k := range []string{"no", "yes", "maybe", "wow"} {
			if _, exists := aggregated[k]; !exists {
				aggregated[k] = 0
			}
		}
	}

	aggregated["total"] = total

	return aggregated
}

// Perform component's algolia's record sync
func (component *ComponentModel) UpdateAlgolia() {

	index := component.di.Search.Get("components")

	// Compose algolia item
	full_name := component.FullName

	if full_name == "" {
		full_name = component.Name
	}

	var image string

	if len(component.Images) > 0 {
		image = component.Images[0]
	}

	object := make(map[string]interface{})
	object["objectID"] = component.Id.Hex()
	object["name"] = component.Name
	object["full_name"] = full_name
	object["part_number"] = component.PartNumber
	object["slug"] = component.Slug
	object["image"] = image
	object["type"] = component.Type
	object["activated"] = component.Activated

	v, err := component.GetVendor("spartangeek")

	if err == nil {

		object["price"] = v.Price
		object["priority"] = v.Priority
		object["stock"] = v.Stock
	}

	_, err = index.UpdateObject(object)

	if err != nil {
		panic(err)
	}
}
