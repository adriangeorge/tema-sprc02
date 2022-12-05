package api

import (
	"encoding/json"
	"fmt"
	db "service/database"
	"service/models"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// POST /api/countries
func CreateCountry(ctx echo.Context) error {
	// Marshall body
	var new_tara models.Tara
	err := json.NewDecoder(ctx.Request().Body).Decode(&new_tara)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, datatypes")
	}

	// Validate step
	if err := db.Validator.Struct(new_tara); err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, missing fields")
	}

	// Reserve index before adding to db
	id, err := db.IndexQuery("id_tari", true)
	if err != nil {
		return ctx.JSON(500, "Something went wrong")
	}
	new_tara.Id = id

	_, err = db.TaraCollection.InsertOne(db.Context, new_tara)
	if err != nil {
		db.IndexQuery("id_tari", false)
		return ctx.JSON(409, "Invalid request, country name conflict")
	}
	return ctx.JSON(201, bson.M{"id": new_tara.Id})
}

// GET /api/countries
func GetCountries(ctx echo.Context) error {
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
func UpdateCountry(ctx echo.Context) error {
	fmt.Println("UPDATING")
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid id")
	}

	// Marshall body
	var new_tara models.Tara
	err = json.NewDecoder(ctx.Request().Body).Decode(&new_tara)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, datatypes")
	}

	// Validate step
	if err := db.Validator.Struct(new_tara); err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid request body, missing fields")
	}

	new_tara.Id = int(id)
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: new_tara}}

	res, err := db.TaraCollection.UpdateOne(db.Context, filter, update)
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
func DeleteCountry(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(400, "Invalid id")
	}

	res, err := db.TaraCollection.DeleteOne(db.Context, bson.D{{Key: "id", Value: id}})

	if err != nil || res.DeletedCount == 0 {
		fmt.Println(err)
		return ctx.JSON(404, "Invalid request")
	} else {
		return ctx.JSON(200, "Delete successful")
	}

}
