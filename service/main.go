package main

import (
	"service/api"
	db "service/database"

	"github.com/labstack/echo/v4"
)

func main() {

	// Connect to database
	db.DB = db.ConnectToDatabase()
	// Start server
	e := echo.New()

	// List POST routes
	e.POST("/api/countries", api.CreateCountry)
	e.POST("/api/cities", api.CreateCity)
	e.POST("/api/temperatures", api.CreateTemp)

	// List GET routes
	e.GET("/api/countries", api.GetCountries)
	e.GET("/api/cities", api.GetCities)
	e.GET("/api/cities/country/:id_Tara", api.GetCitiesByCountry)
	e.GET("/api/temperatures", api.GetTempsParams)
	e.GET("/api/temperatures/cities/:id_oras", api.GetTempsParamsIdCity)
	e.GET("/api/temperatures/countries/:id_tara", api.GetTempsParamsIdCountry)

	// List PUT routes\
	e.PUT("/api/countries/:id", api.UpdateCountry)
	e.PUT("/api/cities/:id", api.UpdateCity)
	e.PUT("/api/temperatures/:id", api.UpdateTemp)

	// List DELETE routes
	e.DELETE("/api/countries/:id", api.DeleteCountry)
	e.DELETE("/api/cities/:id", api.DeleteCity)
	e.DELETE("/api/temperatures/:id", api.DeleteTemp)

	e.Logger.Fatal(e.Start(":6000"))
}
