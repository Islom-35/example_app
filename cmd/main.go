package main

import (
	"example/internal/db"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	productAdap "example/internal/product/adapters"
	productApp "example/internal/product/app"
	productHandler "example/internal/product/controller/handler"

	userAdap "example/internal/users/adapters"
	userApp "example/internal/users/app"
	userHandler "example/internal/users/controller/handler"
)

func main() {
	db, err := db.New(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
	)

	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	// defer db.Close()

	router := gin.Default()

	// Define a route and its handler
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})

	productRepo := productAdap.NewProductRepository(db.DB)
	productService := productApp.NewProductService(productRepo)
	productHandler := productHandler.NewProductHandler(productService)

	productGroup := router.Group("/product")

	{
		productGroup.POST("", productHandler.CreateProduct)
		productGroup.GET("/:id", productHandler.GetProductByID)
		productGroup.GET("/page", productHandler.GetPagesProduct)
		productGroup.DELETE("/:id", productHandler.DeleteProductByID)
		productGroup.PUT("/:id", productHandler.UpdateProductByID)
	}

	userRepo := userAdap.NewUserRepository(db.DB)
	userService := userApp.NewUserService(userRepo)
	userHandler := userHandler.NewUserHandler(userService)

	userGroup := router.Group("/users")

	{
		userGroup.POST("/login", userHandler.LoginUserHandler)
		userGroup.POST("/sign-up", userHandler.SignUpUserHandler)
		userGroup.GET("/page", userHandler.GetPagesUser)
	}

	server := &http.Server{Addr: os.Getenv("HTTP_PORT"), Handler: router}
	log.Println("Starting server on port...", os.Getenv("HTTP_PORT"))
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Println(">>>>", err)
	}
	log.Println(server.Addr)
	log.Println(">>>")
	defer server.Close()

}
