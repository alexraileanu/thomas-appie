package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
	"github.com/hasura/go-graphql-client"
)

type productToCheck struct {
	ID           int
	ApiName      string
	FriendlyName string
	TaxonomyId   int

	BonusData GQLQuery
}

const AppieURL = "https://www.ah.nl/gql"

func main() {
	thomas, err := initThomas()
	if err != nil {
		panic(err)
	}

	s := gocron.NewScheduler(time.Local)

	// scheduler that runs every day at 10AM (for now for debug purposes only)
	// eventually it will run every monday
	s.Every(1).Week().Monday().At("10:30").Do(func() {
		productsToWatch, err := parseProductsJson()
		if err != nil {
			panic(err)
		}
		goThomasGo(thomas, productsToWatch)
	})
	s.StartBlocking()

	handleClose(thomas)
}

// checkProduct gets the product info from the API
func checkProduct(product productToCheck) (GQLQuery, error) {
	gql := new(Query)
	client := graphql.NewClient(AppieURL, nil).WithRequestModifier(func(request *http.Request) {
		request.Header.Add("client-name", "ah-products")
		request.Header.Add("client-version", "6.462.0")
	})
	vars := map[string]interface{}{
		"id":   product.ID,
		"date": time.Now().Format("2006-01-02"),
	}

	err := client.Query(context.Background(), &gql.Query, vars)
	if err != nil {
		panic(err)
	}

	return gql.Query, nil
}

// initThomas makes the initial connection to discord
func initThomas() (*discordgo.Session, error) {
	thomas, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		return nil, err
	}

	err = thomas.Open()
	if err != nil {
		return nil, err
	}

	return thomas, nil
}

// handleClose cleanly disconnects and shuts thomas down
func handleClose(thomas *discordgo.Session) {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Thomas is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	thomas.Close()
}

// goThomasGo fetches the info from the API and instructs thomas to send the discord message
func goThomasGo(thomas *discordgo.Session, products []productToCheck) {
	var productsInBonus []productToCheck
	var productsNotInBonus []productToCheck

	for _, product := range products {
		response, err := checkProduct(product)
		if err != nil {
			panic(err)
		}
		hasDiscount := response.Product.Price.Discount.SegmentId != 0
		if hasDiscount {
			productsInBonus = append(productsInBonus, productToCheck{
				FriendlyName: product.FriendlyName,
				ApiName:      product.ApiName,
				BonusData:    response,
			})
		} else {
			productsNotInBonus = append(productsNotInBonus, productToCheck{
				FriendlyName: product.FriendlyName,
				ApiName:      product.ApiName,
			})
		}
	}

	var inBonusFields []*discordgo.MessageEmbedField
	if len(productsInBonus) != 0 {
		for _, prod := range productsInBonus {
			inBonusFields = append(inBonusFields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("%s (%s)", prod.FriendlyName, prod.ApiName),
				Value: fmt.Sprintf("%s %s", prod.BonusData.Product.Price.Discount.Description, prod.BonusData.Product.SmartLabel),
			})
		}
	}

	var notInBonusFields []*discordgo.MessageEmbedField
	if len(productsNotInBonus) != 0 {
		for _, prod := range productsNotInBonus {
			notInBonusFields = append(notInBonusFields, &discordgo.MessageEmbedField{
				Name: fmt.Sprintf("%s (%s)", prod.FriendlyName, prod.ApiName),
			})
		}
	}

	embeds := []*discordgo.MessageEmbed{
		{
			Color:  0xff7900,
			Title:  "Products that are in bonus this week at the Appie",
			Fields: inBonusFields,
		},
		{
			Title:  "Products that aren't in bonus this week at the Appie",
			Color:  0xff0000,
			Fields: notInBonusFields,
		},
	}
	thomas.ChannelMessageSendEmbeds(os.Getenv("DISCORD_CHANNEL_ID"), embeds)
}

func parseProductsJson() ([]productToCheck, error) {
	fileContents, err := os.ReadFile(os.Getenv("PRODUCTS_JSON_FILE_PATH"))
	if err != nil {
		return nil, err
	}

	products := new([]productToCheck)
	err = json.Unmarshal(fileContents, products)
	if err != nil {
		return nil, err
	}

	return *products, nil
}
