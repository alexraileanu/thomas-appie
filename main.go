package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
	_ "github.com/joho/godotenv/autoload"
)

type shield struct {
	Theme string `json:"theme"`
	Text  string `json:"text"`
}

type discount struct {
	Theme string
	Start string `json:"startDate"`
	End   string `json:"endDate"`
}

type productResponse struct {
	Cards []struct {
		Products []struct {
			Shield   shield   `json:"shield"`
			Discount discount `json:"discount"`
		} `json:"products"`
	} `json:"cards"`
}

type bonusData struct {
	Shield   shield
	Discount discount
}

type productToCheck struct {
	ApiName      string `json:"apiName"`
	FriendlyName string `json:"friendlyName"`
	TaxonomyId   int    `json:"taxonomyId"`

	BonusData bonusData
}

const AppieURL = "https://www.ah.nl/zoeken/api/products/taxonomy-brand?brand=%s&taxonomyId=%d"

func main() {
	thomas := initThomas()

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
func checkProduct(product productToCheck) (*http.Response, error) {
	productUrl := fmt.Sprintf(AppieURL, url.QueryEscape(product.ApiName), product.TaxonomyId)
	req, err := http.NewRequest("GET", productUrl, bytes.NewBuffer(nil))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// initThomas makes the initial connection to discord
func initThomas() *discordgo.Session {
	thomas, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}

	err = thomas.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	return thomas
}

// parseBody parses the response body into the given target
func parseBody(resp *http.Response, target any) error {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, target)
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
		resp, err := checkProduct(product)
		if err != nil {
			panic(err)
		}

		response := new(productResponse)
		err = parseBody(resp, &response)
		if err != nil {
			panic(err)
		}

		p := response.Cards[0].Products[0]

		// shield and discount hold the actual discount info
		// the properties aren't actually present if the p is not discounted
		// so we check if the shield/discount values are equal to their respective empty structs
		hasShield := p.Shield != shield{}
		hasDiscount := p.Discount != discount{}

		if hasShield && hasDiscount {
			productsInBonus = append(productsInBonus, productToCheck{
				FriendlyName: product.FriendlyName,
				ApiName:      product.ApiName,
				BonusData: bonusData{
					Shield: shield{
						Text: p.Shield.Text,
					},
					Discount: discount{
						Start: p.Discount.Start,
						End:   p.Discount.End,
					},
				},
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
				Value: fmt.Sprintf("%s; starts: %s, ends: %s", prod.BonusData.Shield.Text, prod.BonusData.Discount.Start, prod.BonusData.Discount.End),
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
