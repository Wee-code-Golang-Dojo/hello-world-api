package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)


type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Email string `json:"email"`
	BloodType string `json:"blood_type"`
}

var Users []User

var dbClient *mongo.Client

func main() {
	// connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// we are trying to connect to mongodb on a specified URL - mongodb://localhost:27017
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		// if there was an error connecting
		// print out the error and exit using log.Fatal()
		log.Fatalf("Could not connect to the db: %v\n", err)
	}

	dbClient = client
	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		// if there was an issue with the pin
		// print out the error and exit using log.Fatal()
		log.Fatalf("MOngo db not available: %v\n", err)
	}

	// create a new gin router
	router := gin.Default()

	// define a single endpoint
	router.GET("/", helloWorldhandler)

	// CRUD enpoints for data

	// create
	router.POST("/createUser", createUserHandler)

	// retrieve
	// name is a placeholder to represent data the users send
	router.GET("/getUser/:name", getSingleUserHandler)

	router.GET("/getUsers", getAllUserHandler)

	// update
	router.PATCH("/updateUser", updateUserHandler)

	// delete
	router.DELETE("/deleteUser", helloWorldhandler)



	// run the server on the port 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	_ = router.Run(":"+ port)
}



func helloWorldhandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
	})
}

func createUserHandler(c *gin.Context) {
	// create an empty user object
	var user User

	// gets the user data that was sent from the client
	// fills up our empty user object with the sent data
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return
	}

	//// add single user to the list of users
	//Users = append(Users, user)

	// using the dbclient we created above
	// specify the dbname and the collection name we want to add data to using the InsertOne method.
	_, err = dbClient.Database("usersdb").Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println("error saving user", err)
	//	if saving ws not successful
		c.JSON(500, gin.H{
			"error": "Could not process request, could not save user",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "succesfully created user",
		"data": user,
	})
}

func getSingleUserHandler(c *gin.Context) {
	// get the value passed from the client
	name := c.Param("name")

	fmt.Println("name", name)

	// create an empty user
	var user User
	// initialize a boolean variable as false
	userAvailable := false

	// loop through the users array to find a match
	for _, value := range Users {

		// check the current iteration of users
		// check if the name matches the client request
		if value.Name == name {
			// if it matches assign the value to the empty user object we created
			user = value

			// set user available boolean to true since there was a match
			userAvailable = true
		}
	}

	// if no match was found
	// the userAvailable would still be false, if so return a not found error
	// check if user is empty, if so return a not found error
	if !userAvailable {
		c.JSON(404, gin.H{
			"error": "no user with name found: " + name,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": user,
	})
}

func getAllUserHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "success",
		"data": Users,
	})
}

func updateUserHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "User updated!",
	})
}

func deleteUserHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "user deleted!",
	})
}





