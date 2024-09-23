package main

import (
	"github.com/alexraileanu/thomas-appie/pkg/logger"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"

	"github.com/alexraileanu/thomas-appie/pkg/db"
	"github.com/alexraileanu/thomas-appie/pkg/http"
	"github.com/alexraileanu/thomas-appie/pkg/thomas"
	"github.com/alexraileanu/thomas-appie/pkg/utl"
)

func main() {
	godotenv.Load()

	loggerService := logger.New()

	loggerService.Info("Starting thomas", nil)
	s := gocron.NewScheduler(time.Local)

	loggerService.Info("Connecting to the db", nil)
	dbConnection, err := db.New(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	if err != nil {
		loggerService.Error("Error connecting to the db", map[string]interface{}{"error": err.Error()})
		panic(err)
	}
	dbService := db.NewDBService(dbConnection, loggerService)

	t, err := thomas.New(dbService, loggerService)
	if err != nil {
		panic(err)
	}
	// scheduler that runs every monday at 10AM
	s.Every(1).Week().Monday().At("10:30").Do(func() {
		loggerService.Info("Fetching products from the Appie", nil)
		productsToWatch, err := utl.ParseProductsJson()
		if err != nil {
			panic(err)
		}
		t.Go()

		loggerService.Info("Saving products to the db", map[string]interface{}{"products": productsToWatch})
		err = dbService.SaveProduct(productsToWatch)
		if err != nil {
			panic(err)
		}
	})
	s.StartAsync()

	go func() {
		loggerService.Info("Starting http server", nil)
		h := http.NewServer(dbService)
		h.Start()
	}()

	t.Close()
}
