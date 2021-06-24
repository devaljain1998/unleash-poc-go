package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/Unleash/unleash-client-go/v3/context"
	"github.com/gin-gonic/gin"
)

// Globals:
var userIds [10]string
var variantInfo map[string]string

func init() {
	unleashError := unleash.Initialize(
		unleash.WithListener(&unleash.DebugListener{}),
		unleash.WithAppName("my-unleash-poc-app"),
		unleash.WithUrl("http://localhost:4242/api/"),
		unleash.WithCustomHeaders(http.Header{"Authorization": {"4119a4b596fdc7e13dabd29601dd2a63b76da4cec263187d5ffb9514a1d5303b"}}),
		unleash.WithRefreshInterval(10),
	)
	fmt.Println(unleashError)
	if unleashError != nil {
		fmt.Println("Problem in initializing Unleash", unleashError)
		os.Exit(1)
	}
	fmt.Println("Unleash initialized!")

	// Initializing User IDs:
	userIds = [10]string{"U1", "U2", "U3", "U4", "U5", "U6", "U7", "U8", "U9", "U10"}

	// Initializing Map:
	variantInfo = make(map[string]string)
	for _, uid := range userIds {
		variantInfo[uid] = ""
	}

	fmt.Println("map: ", variantInfo)
}

func getRandomUID() string {
	num := rand.Intn(10)
	fmt.Println("got random num: ", num)

	if num >= len(userIds) {
		panic("random number greater then index.")
	}

	return userIds[num]
}

func main() {
	// Setting up logs:
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.Println("Started unleash-poc!")

	server := gin.Default()
	// Base URL:
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Simple GET Url for User ID:
	server.GET("/greet", func(c *gin.Context) {
		// get feature
		featureName := "poc.greet"
		uid := getRandomUID()
		ctx := context.Context{
			UserId: uid,
		}

		// checking if the feature is enabled:
		if unleash.IsEnabled(featureName, unleash.WithContext(ctx)) {
			variant := unleash.GetVariant(featureName)

			var jsonMap map[string]string

			err := json.Unmarshal([]byte(variant.Payload.Value), &jsonMap)
			if err != nil {
				fmt.Println(err)
			}

			variantInfo[uid] = jsonMap["greeting"]

			// response:
			c.JSON(200, jsonMap)
		} else {
			c.JSON(404, gin.H{
				"message": "feature " + featureName + " not found!",
			})
		}

		fmt.Println("variant map: ", variantInfo)
		log.Println("variant map: ", variantInfo)
	})

	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	fmt.Println("Server started on port: ", 8080)
}
