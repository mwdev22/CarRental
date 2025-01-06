package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mwdev22/CarRental/internal/config"
	"github.com/mwdev22/CarRental/internal/database"
	"github.com/mwdev22/CarRental/internal/services"
	"github.com/mwdev22/CarRental/internal/store"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

var (
	testServer   *httptest.Server
	authHeader   string
	testUsername = utils.GenerateUniqueString("testuser")
	testPassword = "testpassword"
)

func TestMain(m *testing.M) {
	_ = config.New()
	// changes the working directory to the project root, important for running migrations
	_, err := initializeTests()
	if err != nil {
		fmt.Printf("failed to initialize tests: %v", err)
		os.Exit(1)
	}

	code := m.Run()

	os.Remove("./test.db")

	os.Exit(code)
}

func initializeTests() (*httptest.Server, error) {
	// change working directory to the project root
	_ = config.New()

	// setup the test DB
	testDB, err := database.OpenTestSqlDB()
	if err != nil {
		return nil, err
	}

	// stores and services
	userStore := store.NewUserRepo(testDB)
	userService := services.NewUserService(userStore)

	companyStore := store.NewCompanyRepository(testDB)
	companyService := services.NewCompanyService(companyStore)

	carStore := store.NewCarRepo(testDB)
	carService := services.NewCarService(carStore)

	// handlers
	mux := http.NewServeMux()
	handlers := []types.Handler{
		NewUserHandler(mux, userService),
		NewCompanyHandler(mux, companyService),
		NewCarHandler(mux, carService),
	}
	for _, h := range handlers {
		h.RegisterRoutes()
	}

	// setup the test server
	testServer = httptest.NewServer(mux)
	return testServer, nil
}
