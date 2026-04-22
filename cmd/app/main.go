package main

import (
	"crud-go/internal/repository"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	db := repository.Connect()
	pingErr := db.Ping()
	if pingErr != nil {
		log.Panic(pingErr)
	} else {
		log.Println("Подключение к БД установлено успешно")
	}

	defer repository.Close(db)

	r := gin.Default()
	r.Run(":" + os.Getenv("APP_PORT"))

}
