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

	router := gin.Default()
	router.Use(middleware2.CORSMiddleware()) //For CORS
	router.Use(middleware2.LoggerMiddleware())
	//For error handling
	router.Use(middleware2.ErrorHandler())
	//routes

	controller.RouterUser(router, db.User, redisService.Auth, tk)
	controller.RouterFood(router, db.Food, db.User, fd, redisService.Auth, tk)
	controller.RouterAuthenticate(router, db.User, redisService.Auth, tk)
	
	//Starting the application
	appPort := os.Getenv("APP_PORT") //using heroku host
	if appPort == "" {
		appPort = "8888" //localhost
	}
	log.Fatal(router.Run(":" + appPort))
}
