// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//		Schemes: http
//	 Host: localhost:8080
//		BasePath: /
//		Version: 1.0.0
//		Contact: Vishal Patel <vishster@gmail.com>
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	handlers "Vishal151/recipes-api/handlers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var authHandler *handlers.AuthHandler
var recipesHandler *handlers.RecipesHandler

func init() {
	/* recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)

	ctx := context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	var listOfRecipies []interface{}
	for _, recipe := range recipes {
		listOfRecipies = append(listOfRecipies, recipe)
	}
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	insertMany, err := collection.InsertMany(ctx, listOfRecipies)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted recipes: ", len(insertMany.InsertedIDs))*/
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("localhost:6379"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func main() {
	router := gin.Default()

	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("recipes_api", store))

	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.POST("/signout", authHandler.SignOutHandler)

	authorized := router.Group("/")

	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
		authorized.GET("recipe/:id", recipesHandler.GetOneRecipeHandler)
		// authorized.GET("/recipes/search", SearchRecipesHandler)
	}
	router.Run()
}
