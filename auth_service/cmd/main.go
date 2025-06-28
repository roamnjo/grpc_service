package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/roamnjo/grpc_service/internal/auth"
	"github.com/roamnjo/grpc_service/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongo(ctx context.Context, uri string) (*mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database("grpc_auth_service"), nil
}

func main() {
	log := logger.New(slog.LevelInfo)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Error("loading .env file:", err)
	}

	db, err := ConnectToMongo(ctx, os.Getenv("DB_URI"))
	if err != nil {
		log.Error("connecting to mongodb:", err)
		return
	}

	repo := auth.NewRepository(db)
	handler := auth.NewHandler(repo, log)

	r := gin.Default()
	r.POST("/signup", handler.SignUp)
	r.POST("/signin", handler.SignIn)

	log.Info("Starting server on port 8080")
	err = r.Run()
	if err != nil {
		log.Error("starting server:", err)
	}
}
