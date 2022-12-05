package api

import (
	"encoding/json"
	"fmt"
	db "service/database"
	"service/models"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// POST /api/cities
func CreateCity(ctx echo.Context) error {
	fmt.Println("City Create Req")
	// Marshall body
	var new_oras models.Orase
	err := json.NewDecoder(ctx.Request().Body).Decode(&new_oras)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, datatypes")
	}

	// Validate step
	if err := db.Validator.Struct(new_oras); err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, missing fields")
	}

	// Check if country exists
	res := db.TaraCollection.FindOne(db.Context, bson.D{{Key: "id", Value: new_oras.Id_tara}})
	if res.Err() == mongo.ErrNoDocuments {
		return ctx.JSON(404, "No such country")
	}
	// Reserve index before adding to db
	id, err := db.IndexQuery("id_orase", true)
	if err != nil {
		return ctx.JSON(500, "Something went wrong")
	}
	new_oras.Id = id

	_, err = db.OraseCollection.InsertOne(db.Context, new_oras)
	if err != nil {
		db.IndexQuery("id_orase", false)
		return ctx.JSON(409, "Invalid request, (city name,id_tara) conflict")
	}
	return ctx.JSON(201, bson.M{"id": new_oras.Id})
}

// GET /api/cities
func GetCities(ctx echo.Context) error {
	res, err := db.OraseCollection.Find(db.Context, bson.D{{}})
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	var orase_list []models.Orase
	err = res.All(db.Context, &orase_list)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	return ctx.JSON(200, orase_list)
}

// GET /api/cities/country/:id_Tara
func GetCitiesByCountry(ctx echo.Context) error {
	id_tara, err := strconv.ParseInt(ctx.Param("id_Tara"), 10, 32)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid id")
	}
	match := bson.D{{Key: "$match", Value: bson.D{{Key: "id_tara", Value: id_tara}}}}
	res, err := db.OraseCollection.Aggregate(db.Context, mongo.Pipeline{match})
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	var orase_list []models.Orase
	err = res.All(db.Context, &orase_list)
	fmt.Printf("%v\n", orase_list)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(500, "Something went wrong")
	}

	return ctx.JSON(200, orase_list)
}

// PUT /api/cities/:id
func UpdateCity(ctx echo.Context) error {
	fmt.Println("UPDATING")
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid id")
	}

	// Marshall body
	var upd_oras models.Orase
	err = json.NewDecoder(ctx.Request().Body).Decode(&upd_oras)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, datatypes")
	}

	// Validate step
	if err := db.Validator.Struct(upd_oras); err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, missing fields")
	}

	upd_oras.Id = int(id)
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: upd_oras}}

	res, err := db.OraseCollection.UpdateOne(db.Context, filter, update)
	if err != nil {
		return ctx.JSON(400, "Bad update")
	}
	if res.ModifiedCount == 0 {
		fmt.Println(err)
		return ctx.JSON(404, "Invalid id")
	}

	return ctx.JSON(200, res)
}

// DELETE /api/cities/:id
func DeleteCity(ctx echo.Context) error {
	fmt.Println("DELETE")
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid id")
	}
	res, err := db.OraseCollection.DeleteOne(db.Context, bson.D{{Key: "id", Value: id}})

	if err != nil || res.DeletedCount == 0 {
		fmt.Println(err)
		return ctx.JSON(404, "Invalid request")
	} else {
		return ctx.JSON(200, "Delete successful")
	}

}
