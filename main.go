package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// create a new gin router
	router := gin.Default()

	// define a single endpoint
	router.GET("/", helloWorldhandler)

	// define a single endpoint
	router.GET("/", saveUser)

	// run the server on the port 3000
	_ = router.Run(":3000")
}

func helloWorldhandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
	})
}

func saveUser(c *gin.Context) {
	err := saveUserToDB("tobi")

	if err != nil {
		c.JSON(500, gin.H{
			"error": "user not saved",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "user saved",
	})
}

func saveUserToDB(name string) error {
	// connect to the db
	// save the user to the db
	return nil
}




