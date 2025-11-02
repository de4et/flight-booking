package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"flight-booking/internal/adapters/protobuf"
	"flight-booking/internal/adapters/redis"
	"flight-booking/internal/database"
	"flight-booking/internal/service"
	"flight-booking/internal/service/providers"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db: database.New(),
	}

	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	slog.Debug("Connecting to redis", "addr", redisAddr)
	c, err := redis.NewRedisSROCache(redisAddr, os.Getenv("REDIS_PASSWORD"), protobuf.NewTripsSerializer())
	if err != nil {
		panic("couldn't start redis")
	}

	svc := service.NewMultipleSearchService(c)
	svc.AddProviderService(providers.NewStubGDS(5))
	svc.AddProviderService(providers.NewStubGDS(1))

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(svc),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
