// package main

// import (
// 	"github.com/Surdy-A/course-api/models"
// 	"github.com/gin-gonic/gin"

// 	"github.com/astaxie/beego/orm"
// 	//"./models"
// 	"net/http"
// )

// var ORM orm.Ormer

// func init() {
// 	models.ConnectToDb()
// 	ORM = models.GetOrmObject()
// }

// func main() {
// 	// Creates a gin router with default middleware:
// 	// logger and recovery (crash-free) middleware
// 	router := gin.Default()

// 	router.POST("/createUser", createUser)
// 	router.GET("/readUsers", readUsers)
// 	router.PUT("/updateUser", updateUser)
// 	router.DELETE("/deleteUser", deleteUser)

// 	// By default it serves on :8080 unless a
// 	// PORT environment variable was defined.
// 	router.Run(":3000")
// 	// router.Run(":3000") for a hard coded port
// }

// func createUser(c *gin.Context) {
// 	var newUser models.Users
// 	c.BindJSON(&newUser)
// 	_, err := ORM.Insert(&newUser)
// 	if err == nil {
// 		c.JSON(http.StatusOK, gin.H{
// 			"status":    http.StatusOK,
// 			"email":     newUser.Email,
// 			"user_name": newUser.UserName,
// 			"user_id":   newUser.UserId})
// 	} else {
// 		c.JSON(http.StatusInternalServerError,
// 			gin.H{"status": http.StatusInternalServerError, "error": "Failed to create the user"})
// 	}
// }

// func readUsers(c *gin.Context) {
// 	var user []models.Users
// 	_, err := ORM.QueryTable("users").All(&user)
// 	if err == nil {
// 		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "users": &user})
// 	} else {
// 		c.JSON(http.StatusInternalServerError,
// 			gin.H{"status": http.StatusInternalServerError, "error": "Failed to read the users"})
// 	}
// }

// func updateUser(c *gin.Context) {
// 	var updateUser models.Users
// 	c.BindJSON(&updateUser)
// 	_, err := ORM.QueryTable("users").Filter("email", updateUser.Email).Update(
// 		orm.Params{"user_name": updateUser.UserName,
// 			"password": updateUser.Password})
// 	if err == nil {
// 		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
// 	} else {
// 		c.JSON(http.StatusInternalServerError,
// 			gin.H{"status": http.StatusInternalServerError, "error": "Failed to update the users"})
// 	}
// }

// func deleteUser(c *gin.Context) {
// 	var delUser models.Users
// 	c.BindJSON(&delUser)
// 	_, err := ORM.QueryTable("users").Filter("email", delUser.Email).Delete()
// 	if err == nil {
// 		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
// 	} else {
// 		c.JSON(http.StatusInternalServerError,
// 			gin.H{"status": http.StatusInternalServerError, "error": "Failed to delete the users"})
// 	}
// }

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	handlers "github.com/Surdy-A/course-api/handlers"

)

var recipesHandler *handlers.recipesHandler

func init()  {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/recipedb?readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	//collection := client.Database("recipedb").Collection("recipes")

	
}
func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	router.Run()
}
