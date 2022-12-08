package api

import (
	"encoding/json"
	"fmt"
	db "service/database"
	"service/models"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// POST /api/temperatures
func CreateTemp(ctx echo.Context) error {
	// Marshall body
	var new_temp models.Temperaturi
	err := json.NewDecoder(ctx.Request().Body).Decode(&new_temp)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, datatypes")
	}

	// Validate step
	if err := db.Validator.Struct(new_temp); err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, missing fields")
	}

	// Check if city exists
	res := db.OraseCollection.FindOne(db.Context, bson.D{{Key: "id", Value: *(new_temp.Id_oras)}})
	if res.Err() == mongo.ErrNoDocuments {
		fmt.Println(*(new_temp.Id_oras))
		return ctx.JSON(404, "No such city")
	}

	// Reserve index before adding to db
	id, err := db.IndexQuery("id_temperaturi", true)
	if err != nil {
		return ctx.JSON(500, "Something went wrong")
	}
	new_temp.Id = id
	new_temp.Timestamp = time.Now().UTC()
	_, err = db.TemperaturiCollection.InsertOne(db.Context, new_temp)
	if err != nil {
		db.IndexQuery("id_temperaturi", false)
		return ctx.JSON(409, "Invalid request, country name conflict")
	}
	return ctx.JSON(201, bson.M{"id": new_temp.Id})
}

// GET /api/temperatures?lat=Double&lon=Double&from=Date&until=Date
func GetTempsParams(ctx echo.Context) error {
	lat := ctx.FormValue("lat")
	lon := ctx.FormValue("lon")
	from := ctx.FormValue("from")
	until := ctx.FormValue("until")
	fmt.Println(lat, lon, from, until)

	var operations []bson.D

	// Join
	city_lookup := bson.D{{
		Key: "$lookup",
		Value: bson.D{
			{Key: "from", Value: "Orase"},
			{Key: "localField", Value: "id_oras"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "o"},
		},
	}}
	operations = append(operations, city_lookup)

	// Location matchers
	if len(lat) > 0 {
		float_val, _ := strconv.ParseFloat(lat, 64)
		fmt.Println(float_val)
		latitude_match := bson.D{{
			Key:   "$match",
			Value: bson.D{{Key: "o.latitudine", Value: float_val}},
		}}
		operations = append(operations, latitude_match)
	}

	if len(lon) > 0 {
		float_val, _ := strconv.ParseFloat(lon, 64)
		longitudine_match := bson.D{{
			Key:   "$match",
			Value: bson.D{{Key: "o.longitudine", Value: float_val}},
		}}
		operations = append(operations, longitudine_match)
	}

	// Time matchers
	if len(from) > 0 {
		date, _ := time.Parse("2006-01-02", from)
		fmt.Println(date)
		from_match := bson.D{{
			Key: "$match",
			Value: bson.D{{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: date},
			}}},
		}}
		operations = append(operations, from_match)
	}

	if len(until) > 0 {
		date, _ := time.Parse("2006-01-02", until)
		fmt.Println(date)
		from_match := bson.D{{
			Key: "$match",
			Value: bson.D{{Key: "timestamp", Value: bson.D{
				{Key: "$lte", Value: date},
			}}},
		}}
		operations = append(operations, from_match)
	}

	res, err := db.TemperaturiCollection.Aggregate(db.Context, operations)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	var temperaturi_list []models.Temperaturi
	fmt.Println(temperaturi_list)

	var return_temp_array []models.Temperaturi_Response = make([]models.Temperaturi_Response, 0)

	err = res.All(db.Context, &temperaturi_list)
	for i, _ := range temperaturi_list {
		var new_temp models.Temperaturi_Response
		new_temp.Id = temperaturi_list[i].Id
		new_temp.Timestamp = temperaturi_list[i].Timestamp
		new_temp.Valoare = temperaturi_list[i].Valoare

		return_temp_array = append(return_temp_array, new_temp)
	}

	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}
	fmt.Println(return_temp_array)
	return ctx.JSON(200, return_temp_array)
}

// GET /api/temperatures/cities/:id_oras?from=Date&until=Date
func GetTempsParamsIdCity(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id_oras"), 10, 32)

	from := ctx.FormValue("from")
	until := ctx.FormValue("until")
	fmt.Println(from, until)

	var operations []bson.D

	// Join
	city_lookup := bson.D{{
		Key: "$lookup",
		Value: bson.D{
			{Key: "from", Value: "Orase"},
			{Key: "localField", Value: "id_oras"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "o"},
		},
	}}
	operations = append(operations, city_lookup)

	city_match := bson.D{{
		Key:   "$match",
		Value: bson.D{{Key: "o.id", Value: id}},
	}}
	operations = append(operations, city_match)

	// Time matchers
	if len(from) > 0 {
		date, _ := time.Parse("2006-01-02", from)
		fmt.Println(date)
		from_match := bson.D{{
			Key: "$match",
			Value: bson.D{{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: date},
			}}},
		}}
		operations = append(operations, from_match)
	}

	if len(until) > 0 {
		date, _ := time.Parse("2006-01-02", until)
		fmt.Println(date)
		from_match := bson.D{{
			Key: "$match",
			Value: bson.D{{Key: "timestamp", Value: bson.D{
				{Key: "$lte", Value: date},
			}}},
		}}
		operations = append(operations, from_match)
	}

	res, err := db.TemperaturiCollection.Aggregate(db.Context, operations)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	var return_temp_array []models.Temperaturi_Response = make([]models.Temperaturi_Response, 0)
	err = res.All(db.Context, &return_temp_array)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}
	fmt.Println(return_temp_array)
	return ctx.JSON(200, return_temp_array)
}

// GET /api/temperatures/countries/:id_tara?from=Date&until=Date
func GetTempsParamsIdCountry(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id_tara"), 10, 32)

	from := ctx.FormValue("from")
	until := ctx.FormValue("until")
	fmt.Println(from, until)

	var operations []bson.D

	// Join
	city_lookup := bson.D{{
		Key: "$lookup",
		Value: bson.D{
			{Key: "from", Value: "Orase"},
			{Key: "localField", Value: "id_oras"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "o"},
		},
	}}
	operations = append(operations, city_lookup)

	country_lookup := bson.D{{
		Key: "$lookup",
		Value: bson.D{
			{Key: "from", Value: "Tari"},
			{Key: "localField", Value: "o.id_tara"},
			{Key: "foreignField", Value: "id"},
			{Key: "as", Value: "t"},
		},
	}}
	country_match := bson.D{{
		Key:   "$match",
		Value: bson.D{{Key: "t.id", Value: id}},
	}}
	operations = append(operations, country_lookup)
	operations = append(operations, country_match)

	// Time matchers
	if len(from) > 0 {
		date, _ := time.Parse("2006-01-02", from)
		fmt.Println(date)
		from_match := bson.D{{
			Key: "$match",
			Value: bson.D{{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: date},
			}}},
		}}
		operations = append(operations, from_match)
	}

	if len(until) > 0 {
		date, _ := time.Parse("2006-01-02", until)
		fmt.Println(date)
		from_match := bson.D{{
			Key: "$match",
			Value: bson.D{{Key: "timestamp", Value: bson.D{
				{Key: "$lte", Value: date},
			}}},
		}}
		operations = append(operations, from_match)
	}

	res, err := db.TemperaturiCollection.Aggregate(db.Context, operations)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	var temperaturi_list []models.Temperaturi
	fmt.Println(temperaturi_list)

	var return_temp_array []models.Temperaturi_Response = make([]models.Temperaturi_Response, 0)

	err = res.All(db.Context, &return_temp_array)

	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}
	fmt.Println(return_temp_array)
	return ctx.JSON(200, return_temp_array)
}

// PUT /api/countries/:id
func UpdateTemp(ctx echo.Context) error {
	fmt.Println("UPDATING")
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid id")
	}

	// Marshall body
	var new_temp models.Temperaturi
	err = json.NewDecoder(ctx.Request().Body).Decode(&new_temp)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, datatypes")
	}

	// Validate step
	if err := db.Validator.Struct(new_temp); err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, missing fields")
	}

	new_temp.Id = int(id)
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: new_temp}}

	res, err := db.TemperaturiCollection.UpdateOne(db.Context, filter, update)
	if err != nil {
		return ctx.JSON(409, "Bad update")
	}
	if res.ModifiedCount == 0 {
		fmt.Println(err)
		return ctx.JSON(404, "Invalid id")
	}

	return ctx.JSON(200, res)
}

// DELETE /api/countries/:id
func DeleteTemp(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid id")
	}

	res, err := db.TemperaturiCollection.DeleteOne(db.Context, bson.D{{Key: "id", Value: id}})

	if err != nil || res.DeletedCount == 0 {
		fmt.Println(err)
		return ctx.JSON(404, "Invalid request")
	} else {
		return ctx.JSON(200, "Delete successful")
	}

}
