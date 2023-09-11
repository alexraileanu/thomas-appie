package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

type productToCheck struct {
	ID           int
	ApiName      string
	FriendlyName string
	RefererUrl   string

	BonusData productInfoResponse
}

const AppieURL = "https://www.ah.nl/gql"

func main() {
	godotenv.Load()
	thomas, err := initThomas()
	if err != nil {
		panic(err)
	}

	s := gocron.NewScheduler(time.Local)

	// scheduler that runs every monday at 10AM
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

func readQueryFile() (string, error) {
	data, err := os.ReadFile("queryFormat.json")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// checkProduct gets the product info from the API
func checkProduct(product *productToCheck) error {
	gqlQuery, err := readQueryFile()
	if err != nil {
		return err
	}

	preparedRequest := fmt.Sprintf(gqlQuery, product.ID, time.Now().Format("2006-01-02"))

	client := resty.New()
	r := client.R()

	// we pretend we're a valid browser request
	r.Header.Add("client-name", "ah-products")
	r.Header.Add("client-version", "6.500.0")
	r.Header.Add("Referer", product.RefererUrl)
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0")
	r.Header.Add("Content-Type", "application/json")

	resp, err := r.SetBody(preparedRequest).Post(AppieURL)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp.Body(), &product.BonusData)
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
		err := checkProduct(&product)
		if err != nil {
			panic(err)
		}
		hasDiscount := product.BonusData.Data.Product.Price.Discount.SegmentId != 0
		if hasDiscount {
			productsInBonus = append(productsInBonus, productToCheck{
				FriendlyName: product.FriendlyName,
				ApiName:      product.ApiName,
				BonusData:    product.BonusData,
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
				Value: fmt.Sprintf("%s %s", prod.BonusData.Data.Product.Price.Discount.Description, prod.BonusData.Data.Product.SmartLabel),
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
