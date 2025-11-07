package main

import (
	"os"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/joho/godotenv"

	"github.com/alexraileanu/thomas-appie/pkg/config"
	"github.com/alexraileanu/thomas-appie/pkg/logger"
	"github.com/alexraileanu/thomas-appie/pkg/thomas"

	"github.com/alexraileanu/thomas-appie/pkg/db"
	"github.com/alexraileanu/thomas-appie/pkg/http"
)

func main() {
	godotenv.Load()

	enableLogs := os.Getenv("ENABLE_LOGS") == "true"
	loggerService := logger.New(enableLogs)

	conf := config.New()
	err := conf.ParseConfig()
	if err != nil {
		loggerService.Error("Error parsing config", map[string]interface{}{"error": err.Error()})
		panic(err)
	}

	loggerService.Info("Starting thomas", nil)
	s := gocron.NewScheduler(time.Local)

	loggerService.Info("Connecting to the db", nil)
	dbConnection, err := db.New(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	if err != nil {
		//loggerService.Error("Error connecting to the db", map[string]interface{}{"error": err.Error()})
		panic(err)
	}
	dbService := db.NewDBService(dbConnection, loggerService)

	t, err := thomas.New(dbService, loggerService, conf)
	if err != nil {
		panic(err)
	}

	loggerService.Info("Starting cron job", map[string]interface{}{"cron": conf.Thomas.Cron})
	s.Cron(conf.Thomas.Cron).Do(func() {
		loggerService.Info("Fetching products from the Appie", nil)
		productsToWatch, err := dbService.GetProducts()
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
		h := http.NewServer(dbService, conf.Appie, loggerService)
		h.Start()
	}()

	t.Close()
}
