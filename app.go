package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

import "github.com/gin-gonic/gin"

import "github.com/Jeffail/gabs"

const AppVersion = "0.2.0"

var dataFile string

func main() {

	listenPort := flag.Int("p", 8000, "The listening port - defaults to 8080.")
	file := flag.String("f", "./data.json", "The JSON file to load- defaults to ./data.json.")
	version := flag.Bool("v", false, "Display the current version")
	flag.Parse()

	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	} else {
		setupGlobalOptions(*file)
		router := setupAPI()

		p := strconv.Itoa(*listenPort)
		addr := ":" + p
		router.Run(addr)
	}

}

// function to set the global variables for listening port and file
func setupGlobalOptions(file string) {
	dataFile = file
}

//  set up the API - define routes and return router
func setupAPI() *gin.Engine {
	// turn off that noisy notice from Gin
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// attach Handler functions
	r.Use(Cors())

	// define API routes
	r.GET("/", getJSON)
	r.GET("/:name", getAll)
	r.GET("/:name/:id", getById)
	r.POST("/:name", create)
	r.DELETE("/:name/:id", deleteById)

	// routes for js and cors
	r.OPTIONS("/:name", accessControlHeaders)
	r.OPTIONS("/:name/:id", accessControlHeaders)

	return r
}

// CORS support handler
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

// ---- handlers

// GET /
// This displays the serialized JSON from the file. It's a sanity check
// so the user knows all the data imported properly from the file. We
// parse it and then print it back out.
func getJSON(c *gin.Context) {
	code := 200
	jsonParsed, err := getData()
	if err != nil {
		code = 404
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(code, jsonParsed.String())
}

// GET /:name
// Displays the json records for the resource
func getAll(c *gin.Context) {
	code := 200
	jsonParsed, err := getData()
	if err != nil {
		code = 404
	}
	name := c.Params.ByName("name")
	result := jsonParsed.S(name).String()
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(code, result)
}

// GET /:name/:id
// Displays the json record for the given id
func getById(c *gin.Context) {
	code := 200
	name := c.Params.ByName("name")
	id := c.Params.ByName("id")
	result := "{}"

	c.Writer.Header().Set("Content-Type", "application/json")

	record, err := findById(name, id)

	if err != nil {
		code = 404
	}

	result = record.String()

	c.String(code, result)
}

// DELETE /:name/:id
// Displays the json record for the given id
func deleteById(c *gin.Context) {
	name := c.Params.ByName("name")
	id := c.Params.ByName("id")
	result := "{}"
	code := 200

	record, err := removeById(name, id)

	if err != nil {
		code = 422
	} else {
		result = record.String()
		code = 200
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(code, result)
}

// POST /:name
// Create a new resource and persist it
func create(c *gin.Context) {
	code := 201

	collectionName := c.Params.ByName("name")

	// Read data from body and parse it out
	data, _ := ioutil.ReadAll(c.Request.Body)

	newRecord, err := createNewRecord(collectionName, data)

	if err != nil {
		code = 422
	}

	c.String(code, newRecord.String())

}

// OPTIONS /:name
// OPTIONS /:name/:id
// Set OPTIONS header response for javascript + CORS via .fetch or .xHTTPrequest
func accessControlHeaders(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, PUT, PATCH")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

// --- JSON logic

// getNextID - simulates the auto-incrementing key feature an API would have.
func getNextId(collection string) (float64, error) {
	// I recognize this implementation is terrible -
	// Read the json, find all ids, find max one, adds one to result.
	// Not really something you can rely on for scaling, and TOTALLY
	// not concurrent. There will probably be a race condition between the time we
	// get this id and the time we add the new record.
	max := 0.0
	jsonParsed, err := getData()

	if err != nil {
		return -1, fmt.Errorf("Unable to read data")
	}

	records, _ := jsonParsed.S(collection).Children()
	for _, record := range records {
		value := record.S("id").Data().(float64)
		if value > max {
			max = value
		}
	}

	return (max + 1), nil
}

// Takes a collection and finds the record with the id. collection and id are
// strings
func findById(name string, idString string) (*gabs.Container, error) {
	var result *gabs.Container

	id, err := toFloat(idString)

	if err != nil {
		return result, err
	} else {

		jsonParsed, err := getData()
		if err != nil {
			return result, err
		}

		collection := jsonParsed.S(name)
		children, _ := collection.Children()
		for _, child := range children {
			if child.S("id").Data().(float64) == id {
				result = child
				break
			}
		}

		// if the result is an empty json response then we didn't find anything.
		// I bet there's a better way. But the gabs.Container always
		// returns an empty JSON string if it doesn't match. So this works
		// for now.
		if result.String() == "{}" {
			return result, fmt.Errorf("json: no record found for id %g", id)
		}

		return result, nil
	}

}

func removeById(name string, idString string) (*gabs.Container, error) {
	var result *gabs.Container
	indexToDelete := -1

	id, err := toFloat(idString)

	if err != nil {
		return result, err
	} else {

		jsonParsed, err := getData()

		if err != nil {
			return result, err
		}

		collection := jsonParsed.S(name)
		children, _ := collection.Children()

		// find the index of the record we have to delete
		for index, child := range children {

			// if we find it....
			if child.S("id").Data().(float64) == id {
				// save the record we found as the result along with the index
				result = child
				indexToDelete = index
				break
			}
		}

		// if we didn't find the record....
		if indexToDelete == -1 {
			return result, fmt.Errorf("json: no record found for id %g", id)
		} else {

			// remove the index we found from the array
			err = jsonParsed.ArrayRemove(indexToDelete, name)

			if err == nil {
				saveData(jsonParsed)
			}
		}

		return result, err
	}

}

// creates a new JSON record and adds it to the collection.
func createNewRecord(collectionName string, data []byte) (*gabs.Container, error) {
	// parse the body to a new record.
	newRecord, err := gabs.ParseJSON([]byte(data))
	if err != nil {
		return newRecord, err
	}

	// give data from user an id
	id, err := getNextId(collectionName)

	if err != nil {
		return newRecord, err
	}
	newRecord.Set(id, "id")

	// update our overall collection
	jsonParsed, err := getData()
	if err != nil {
		return newRecord, err
	} else {
		jsonParsed.ArrayAppend(newRecord.Data(), collectionName)
		err := saveData(jsonParsed)
		return newRecord, err
	}
}

// Load the json from the file
func getData() (*gabs.Container, error) {

	// TODO: probably not super efficient to read the file all the time but
	// I don't know another way to make this global right now.

	// TODO: I bet gabs can read json right from a file. I should explore that
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		fmt.Print("No data file found. Exiting.\n")
		os.Exit(1)
	}

	return gabs.ParseJSON([]byte(data))
}

// save json back to file.
func saveData(json *gabs.Container) error {
	data := []byte(json.String())
	return ioutil.WriteFile(dataFile, data, 0644)
}

// wrapper for float
func toFloat(id string) (float64, error) {
	return strconv.ParseFloat(id, 64)
}
