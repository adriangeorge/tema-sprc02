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

	var operations []bson.M
	
	res, err := db.TaraCollection.Find(db.Context, bson.D{{}})
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	var tara_list []models.Tara
	err = res.All(db.Context, &tara_list)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	return ctx.JSON(200, tara_list)
}

// GET /api/temperatures/cities/:id_oras?from=Date&until=Date
func GetTempsParamsIdCity(ctx echo.Context) error {
	res, err := db.TaraCollection.Find(db.Context, bson.D{{}})
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	var tara_list []models.Tara
	err = res.All(db.Context, &tara_list)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	return ctx.JSON(200, tara_list)
}

// GET /api/temperatures/countries/:id_tara?from=Date&until=Date
func GetTempsParamsIdCountry(ctx echo.Context) error {
	res, err := db.TaraCollection.Find(db.Context, bson.D{{}})
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	var tara_list []models.Tara
	err = res.All(db.Context, &tara_list)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	return ctx.JSON(200, tara_list)
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
		return ctx.JSON(400, "Bad update")
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
