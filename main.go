package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO Add DELETE method
// TODO Change IDs to int
// TODO Add auto-generated IDs

// ! Change comments !
// furniture represents data about a record furniture
type furniture struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Fabricator string  `json:"fabricator"`
	Height     float64 `json:"height"`
	Length     float64 `json:"length"`
	Width      float64 `json:"width"`
}

// furniture slice to seed record furniture data.
var furnitureSlice = []furniture{
	{ID: "1", Name: "Table", Fabricator: "IKEA", Height: 100.5, Length: 50, Width: 50},
	{ID: "2", Name: "Chair", Fabricator: "Hoff", Height: 50, Length: 33.3, Width: 25.7},
	{ID: "3", Name: "Bed", Fabricator: "IKEA", Height: 50.8, Length: 200, Width: 150},
}

func main() {
	router := gin.Default()
	router.GET("/furniture", getFurniture)
	router.GET("/furniture/:id", getFurnitureByID)
	router.POST("/furniture", postFurniture)

	router.Run("localhost:8080")
}

// getFurniture responds with the list of all furniture as JSON.
func getFurniture(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, furnitureSlice)
}

// postFurniture adds an furniture from JSON received in the request body.
func postFurniture(c *gin.Context) {
	var newFurniture furniture

	if err := c.BindJSON(&newFurniture); err != nil {
		return
	}

	furnitureSlice = append(furnitureSlice, newFurniture)
	c.IndentedJSON(http.StatusCreated, newFurniture)
}

// getFurnitureByID locates the furniture whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getFurnitureByID(c *gin.Context) {
	var id = c.Param("id")

	for _, furnitureItem := range furnitureSlice {
		if furnitureItem.ID == id {
			c.IndentedJSON(http.StatusOK, furnitureItem)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "furniture not found"})
}
