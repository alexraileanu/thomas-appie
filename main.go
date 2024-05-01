package main

import (
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

	s := gocron.NewScheduler(time.Local)
	dbConnection, err := db.New(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	if err != nil {
		panic(err)
	}
	dbService := db.NewDBService(dbConnection)
	t, err := thomas.New(dbService)
	if err != nil {
		panic(err)
	}
	// scheduler that runs every monday at 10AM
	s.Every(1).Week().Monday().At("10:30").Do(func() {
		productsToWatch, err := utl.ParseProductsJson()
		if err != nil {
			panic(err)
		}
		t.Go()
		err = dbService.SaveProduct(productsToWatch)
		if err != nil {
			panic(err)
		}
	})
	s.StartAsync()

	go func() {
		h := http.NewServer(dbService)
		h.Start()
	}()

	t.Close()
}
