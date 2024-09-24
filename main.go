package main

import (
	"food-app/infrastructure/auth"
	"food-app/infrastructure/persistence"
	middleware2 "food-app/interfaces/adapter/middleware"
	"food-app/interfaces/common/file_upload"
	"food-app/interfaces/controller"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	// create database connection
	db, err := persistence.NewDatabase()
	if err != nil {
		panic(err)
	}
	defer func(db *persistence.Repositories) {
		err := db.Close()
		if err != nil {
			log.Println("error closing the database: ", err)
			return
		}
	}(db)
	// create redis connection
	redisService, err := auth.NewRedisDB()
	if err != nil {
		log.Fatal(err)
		return
	}

	tk := auth.NewToken()
	fd := file_upload.NewFileUpload()

	users := controller.NewUsers(db.User, redisService.Auth, tk)
	foods := controller.NewFood(db.Food, db.User, fd, redisService.Auth, tk)
	authenticate := controller.NewAuthenticate(db.User, redisService.Auth, tk)

	r := gin.Default()
	r.Use(middleware2.CORSMiddleware()) //For CORS

	//user rout
	userRoutes := r.Group("/v1/users")
	{
		userRoutes.POST("", users.SaveUser)
		userRoutes.GET("", users.GetUsers)
		userRoutes.GET("/:user_id", users.GetUser)
	}

	//post routes
	foodRoutes := r.Group("/v1/food")
	{
		foodRoutes.POST("", middleware2.AuthMiddleware(), middleware2.MaxSizeAllowed(8192000), foods.SaveFood)
		foodRoutes.PUT("/:food_id", middleware2.AuthMiddleware(), middleware2.MaxSizeAllowed(8192000), foods.UpdateFood)
		foodRoutes.GET("/:food_id", foods.GetFoodAndCreator)
		foodRoutes.DELETE("/:food_id", middleware2.AuthMiddleware(), foods.DeleteFood)
		foodRoutes.GET("", foods.GetAllFood)
	}

	//authentication routes
	authentication := r.Group("/v1/auth")
	{
		authentication.POST("/login", authenticate.Login)
		authentication.POST("/logout", authenticate.Logout)
		authentication.POST("/refresh", authenticate.Refresh)
	}

	//Starting the application
	appPort := os.Getenv("APP_PORT") //using heroku host
	if appPort == "" {
		appPort = "8888" //localhost
	}
	log.Fatal(r.Run(":" + appPort))
}
