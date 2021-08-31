package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Email string `json:"email"`
	BloodType string `json:"blood_type"`
}

func main() {
	// create a new gin router
	router := gin.Default()

	// define a single endpoint
	router.GET("/", helloWorldhandler)

	// CRUD enpoints for data

	// create
	router.POST("/createUser", createUserHandler)

	// retrieve
	router.GET("/getUser", getSingleUserHandler)

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
	//create user
	//......
	var user User
	//user := User{}

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return
	}

	// save user somewhere ...

	c.JSON(200, gin.H{
		"message": "succesfully created user",
		"data": user,
	})
}

func getSingleUserHandler(c *gin.Context) {
	var user User
	user = User{
		Name: "victor",
		Age: 1243,
		Email: "my@email.com",
	}
	c.JSON(200, gin.H{
		"message": "hello world",
		"data": user,
	})
}

func getAllUserHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
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





