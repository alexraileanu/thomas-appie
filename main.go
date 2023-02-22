package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
	_ "github.com/joho/godotenv/autoload"
)

type shield struct {
	Theme string
	Text  string
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

const AppieURL = "https://www.ah.nl/zoeken/api/products/taxonomy-brand?brand=Perla%20Superiore"

//const AppieURL = "https://www.ah.nl/zoeken/api/products/taxonomy-brand?brand=Perla%20Huisblends"

func main() {

	thomas := initThomas()

	s := gocron.NewScheduler(time.Local)
	s.Every(1).Day().At("10:30").Do(func() {
		goThomasGo(thomas)
	})

	handleClose(thomas)
}

func makeRequest() (*http.Response, error) {
	req, err := http.NewRequest("GET", AppieURL, bytes.NewBuffer(nil))
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

func handleClose(thomas *discordgo.Session) {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Thomas is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	thomas.Close()
}

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

func goThomasGo(thomas *discordgo.Session) {
	resp, err := makeRequest()
	if err != nil {
		panic(err)
	}

	response := new(productResponse)
	err = parseBody(resp, &response)
	if err != nil {
		panic(err)
	}

	product := response.Cards[0].Products[0]
	hasShield := product.Shield != shield{}
	hasDiscount := product.Discount != discount{}

	msg := "Beans not in bonus :("
	if hasShield && hasDiscount {
		msg = fmt.Sprintf(`
Beans are in bonus!

%s (starts: %s; ends: %s)
`, product.Shield.Text, product.Discount.Start, product.Discount.End)
	}

	thomas.ChannelMessageSend(os.Getenv("DISCORD_CHANNEL_ID"), msg)
}
