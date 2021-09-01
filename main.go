package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	router.PATCH("/updateUser/:name", updateUserHandler)

	// delete
	router.DELETE("/deleteUser/:name", deleteUserHandler)



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

	// create an empty user
	var user User
	query := bson.M{
		"name": name,
	}
	err := dbClient.Database("usersdb").Collection("users").FindOne(context.Background(), query).Decode(&user)

	// if no match was found
	// err would not be nil
	// so we return a user not found error
	if err != nil {
		fmt.Println("user not found", err)
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
	// create an empty array of users to store the result
	var users []User

	cursor, err := dbClient.Database("usersdb").Collection("users").Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not process request, could get users",
		})
		return
	}

	err = cursor.All(context.Background(), &users)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not process request, could get users",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": users,
	})
}

func updateUserHandler(c *gin.Context) {
	// get the value passed from the client
	name := c.Param("name")

	// creating an empty object to store request data
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
	filterQuery := bson.M{
		"name": name,
	}

	updateQuery := bson.M{
		"$set": bson.M{
			"name": user.Name,
			"age": user.Age,
			"email": user.Email,
		},
	}

	_, err = dbClient.Database("usersdb").Collection("users").UpdateOne(context.Background(), filterQuery, updateQuery)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not process request, could not update user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User updated!",
	})
}

func deleteUserHandler(c *gin.Context) {
	// get the value passed from the client
	name := c.Param("name")

	query := bson.M{
		"name": name,
	}
	_, err := dbClient.Database("usersdb").Collection("users").DeleteOne(context.Background(), query)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not process request, could not delete user",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "user deleted!",
	})
}





