package main

import (
	"github.com/joho/godotenv"

	"github.com/alexraileanu/thomas-appie/pkg/db"
	"github.com/alexraileanu/thomas-appie/pkg/thomas"
	"github.com/alexraileanu/thomas-appie/pkg/utl"
)

func main() {
	godotenv.Load()
	t, err := thomas.New()
	if err != nil {
		panic(err)
	}
	//s := gocron.NewScheduler(time.Local)
	dbConnection, err := db.New()
	if err != nil {
		panic(err)
	}

	// scheduler that runs every monday at 10AM
	//s.Every(1).Week().Monday().At("10:30").Do(func() {
	productsToWatch, err := utl.ParseProductsJson()
	if err != nil {
		panic(err)
	}
	//goThomasGo(thomas, productsToWatch)
	t.Go(productsToWatch)
	err = dbConnection.Save(productsToWatch)
	if err != nil {
		panic(err)
	}
	//})
	//s.StartBlocking()

	t.Close()
}
