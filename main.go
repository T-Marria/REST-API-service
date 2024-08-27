package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TODO Add PATCH method

// furniture represents data about a record furniture
type furniture struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	ManufacturedBy string  `json:"manufacturedBy"`
	Height         float64 `json:"height"`
	Length         float64 `json:"length"`
	Width          float64 `json:"width"`
}

type furnitureDTO struct {
	Name           string  `json:"name"`
	ManufacturedBy string  `json:"manufacturedBy"`
	Height         float64 `json:"height"`
	Length         float64 `json:"length"`
	Width          float64 `json:"width"`
}

var infoFilename = "furniture.json"

// // slice with records about furniture.
//
//	var furnitureSlice = []furniture{
//		{ID: "1", Name: "Table", ManufacturedBy: "IKEA", Height: 100.5, Length: 50, Width: 50},
//		{ID: "2", Name: "Chair", ManufacturedBy: "Hoff", Height: 50, Length: 33.3, Width: 25.7},
//		{ID: "3", Name: "Bed", ManufacturedBy: "IKEA", Height: 50.8, Length: 200, Width: 150},
//	}

func main() {

	router := gin.Default()
	router.GET("/furniture", getFurniture)
	router.GET("/furniture/:id", getFurnitureByID)
	router.POST("/furniture", postFurniture)
	router.PUT("furniture/:id", putFurnitureByID)
	router.PATCH("furniture/:id", patchFurnitureByID)
	router.DELETE("furniture/:id", deleteFurnitureByID)

	router.Run("localhost:8080")
}

func readFurnitureInfo(filename string) []furniture {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	var furnitureSlice []furniture
	jsonErr := json.Unmarshal(data, &furnitureSlice)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return furnitureSlice
}

func updateFurnitureInfo(slice []furniture, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	jsonData, err := json.Marshal(slice)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(infoFilename, jsonData, 0666)
}

func mapFromDTO(dto furnitureDTO, id string) furniture {
	return furniture{
		ID:             id,
		Name:           dto.Name,
		ManufacturedBy: dto.ManufacturedBy,
		Height:         dto.Height,
		Length:         dto.Length,
		Width:          dto.Width,
	}
}

// getFurniture responds with the list of all furniture as JSON.
func getFurniture(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, readFurnitureInfo(infoFilename))
}

// postFurniture adds an furniture unit from JSON received in the request body.
func postFurniture(c *gin.Context) {
	var newFurnitureDTO furnitureDTO

	if err := c.BindJSON(&newFurnitureDTO); err != nil {
		return
	}

	furnitureSlice := readFurnitureInfo(infoFilename)
	id, err := strconv.ParseInt((furnitureSlice[len(furnitureSlice)-1].ID), 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	id += 1

	newFurniture := mapFromDTO(newFurnitureDTO, fmt.Sprint(id))
	furnitureSlice = append(furnitureSlice, newFurniture)
	updateFurnitureInfo(furnitureSlice, infoFilename)

	c.IndentedJSON(http.StatusCreated, newFurniture)
}

// getFurnitureByID locates the furniture unit whose ID value matches the id
// parameter sent by the client, then returns that furniture unit as a response.
func getFurnitureByID(c *gin.Context) {
	var id = c.Param("id")

	furnitureSlice := readFurnitureInfo(infoFilename)
	for _, furnitureUnit := range furnitureSlice {
		if furnitureUnit.ID == id {
			c.IndentedJSON(http.StatusOK, furnitureUnit)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "furniture not found"})
}

// deleteFurnitureByID locates the furniture unit whose ID value matches the id
// parameter sent by the client, then deletes that furniture unit.
func deleteFurnitureByID(c *gin.Context) {
	var id = c.Param("id")

	furnitureSlice := readFurnitureInfo(infoFilename)
	for i := 0; i < len(furnitureSlice); i++ {
		if furnitureSlice[i].ID == id {
			furnitureSlice = append(furnitureSlice[:i], furnitureSlice[i+1:]...)
			updateFurnitureInfo(furnitureSlice, infoFilename)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "furniture not found"})
}

func putFurnitureByID(c *gin.Context) {
	var id = c.Param("id")
	var newFurnitureDTO furnitureDTO

	if err := c.BindJSON(&newFurnitureDTO); err != nil {
		return
	}
	newFurniture := mapFromDTO(newFurnitureDTO, id)

	furnitureSlice := readFurnitureInfo(infoFilename)
	for i := 0; i < len(furnitureSlice); i++ {
		if furnitureSlice[i].ID == id {
			furnitureSlice[i] = newFurniture
			updateFurnitureInfo(furnitureSlice, infoFilename)
			c.IndentedJSON(http.StatusOK, newFurnitureDTO)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "furniture not found"})
}

func patchFurnitureByID(c *gin.Context) {
	var id = c.Param("id")
	// var newFurnitureDTO furnitureDTO
	// newFurniture := mapFromDTO(newFurnitureDTO, id)

	furnitureSlice := readFurnitureInfo(infoFilename)
	for i := 0; i < len(furnitureSlice); i++ {
		if furnitureSlice[i].ID == id {
			if err := c.BindJSON(&furnitureSlice[i]); err != nil {
				fmt.Println("boobs:", err)
				return
			}
			updateFurnitureInfo(furnitureSlice, infoFilename)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "furniture not found"})
}
