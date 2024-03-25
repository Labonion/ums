package main

import (
	"context"
	"log"
	"markie-backend/database"
	"markie-backend/routes"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func handleShutdown(server *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	database.CloseMongoClient()

	log.Println("Server gracefully stopped")
}

func handlePrefight(c *gin.Context) {

}

func main() {
	godotenv.Load()
	database.Redis()
	router := gin.Default()
	allowedOrigins := os.Getenv("CORS_ORIGINS")
	allowedMethods := os.Getenv("CORS_METHODS")
	origins := strings.Split(allowedOrigins, ",")
	methods := strings.Split(allowedMethods, ",")
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     methods,
		AllowHeaders:     []string{"X-API-KEY", "Content-Type"},
		AllowCredentials: true,
	}))

	routes.GetRoutes(router)

	port := os.Getenv("PORT")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	log.Printf("Server is running on port %s", port)

	handleShutdown(server)
}
