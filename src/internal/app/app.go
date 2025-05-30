package app

import (
	"bank_api/src/cmd/docs"
	"bank_api/src/internal/handler"
	"bank_api/src/internal/repository"
	"bank_api/src/internal/service"
	"bank_api/src/pkg/db"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func RunServiceInstance() {
	err := godotenv.Load()
	if err != nil {
		log.Println("failed loading .env: ", err)
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	//if err := migrations.Run(context.Background(), connString); err != nil {
	//	log.Fatal("failed migrations:", err)
	//}

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Println("failed create pool:", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		log.Println("failed ping postgres:", err)
	}

	defer pool.Close()

	dataBase := db.NewDataBase(pool)
	operationRepository := repository.NewOperation(dataBase)
	userRepository := repository.NewUser(dataBase)

	account := service.NewBankAccountService(dataBase, userRepository, operationRepository)
	accountHandler := handler.NewHandler(context.Background(), account)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	gin.ErrorLogger()

	apiV1 := router.Group("/api/v1")
	{
		account := apiV1.Group("/bank_account")

		account.PUT("/update", accountHandler.UpdateBalance)
		account.POST("/transfer", accountHandler.Transfer)
		account.GET("/:id/operations", accountHandler.GetLastOperations)
	}
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
