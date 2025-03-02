package api

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mwdev22/CarRental/docs"
	"github.com/mwdev22/CarRental/internal/handlers"
	"github.com/mwdev22/CarRental/internal/services"
	"github.com/mwdev22/CarRental/internal/store"
	"github.com/mwdev22/CarRental/internal/store/postgres"
	"github.com/mwdev22/CarRental/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type api struct {
	addr string
	db   *sqlx.DB
}

func New(addr string, db *sqlx.DB) *api {
	return &api{
		addr: addr,
		db:   db,
	}
}

func (a *api) Start() error {

	// --- APPLICATION SETUP ---
	mux := http.NewServeMux()
	// log files to inspect
	logDir, err := filepath.Abs("./log")
	if err != nil {
		log.Fatalf("failed to resolve log directory: %v", err)
	}
	fs := http.FileServer(http.Dir(logDir))
	mux.Handle("/log/", http.StripPrefix("/log/", fs))

	// api docs
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// --- STORAGE AND SERVICES ---
	userStore := postgres.NewUserRepo(a.db)
	userService := services.NewUserService(userStore)
	carStore := postgres.NewCarRepository(a.db)
	carService := services.NewCarService(carStore)
	companyStore := postgres.NewCompanyRepository(a.db)
	companyService := services.NewCompanyService(companyStore)
	bookingStore := postgres.NewBookingRepository(a.db)
	bookingService := services.NewBookingService(bookingStore, carStore, userStore)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	userCompCache := store.NewRedisCache(redisClient)

	// --- MAIN ROUTES ---

	_ = handlers.NewUserHandler(mux, userService, utils.MakeLogger("user"))
	_ = handlers.NewCarHandler(mux, carService, utils.MakeLogger("car"))
	_ = handlers.NewCompanyHandler(mux, companyService, utils.MakeLogger("company"))
	_ = handlers.NewBookingHandler(mux, bookingService, carService, userCompCache, utils.MakeLogger("booking"))

	c := cors.New(cors.Options{
		AllowedOrigins:      []string{"*"},
		AllowCredentials:    true,
		AllowedMethods:      []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowPrivateNetwork: true,
		AllowedHeaders:      []string{"*"},
	})

	server := &http.Server{
		Addr:         a.addr,
		Handler:      c.Handler(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  15 * time.Second,
	}
	log.Printf("Starting server on %s", a.addr)

	return server.ListenAndServe()
}
