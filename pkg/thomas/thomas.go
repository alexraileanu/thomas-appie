package thomas

import (
	"fmt"
	"github.com/alexraileanu/thomas-appie/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
	"github.com/alexraileanu/thomas-appie/pkg/db"
)

type Thomas struct {
	session       *discordgo.Session
	dbService     *db.Service
	loggerService *logger.Service
}

func New(dbService *db.Service, loggerService *logger.Service) (*Thomas, error) {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		return nil, err
	}

	err = session.Open()
	if err != nil {
		return nil, err
	}

	return &Thomas{session: session, dbService: dbService, loggerService: loggerService}, nil
}

func (t *Thomas) Go() {
	t.loggerService.Info("Fetching products from the db", nil)
	products, err := t.dbService.GetProducts()
	if err != nil {
		t.loggerService.Error("Error fetching products from the db", map[string]interface{}{"error": err.Error()})
		panic(err)
	}

	a := appie.New(t.loggerService)

	t.loggerService.Info("Checking products from the Appie", nil)
	productsInBonus, productsNotInBonus, err := a.PerformProductsCheck(products)
	if err != nil {
		panic(err)
	}
	var inBonusFields []*discordgo.MessageEmbedField
	if len(productsInBonus) != 0 {
		for _, prod := range productsInBonus {
			inBonusFields = append(inBonusFields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("%s (%s)", prod.FriendlyName, prod.ApiName),
				Value: fmt.Sprintf("%s %s", prod.DiscountedProducts[0].Description, prod.DiscountedProducts[0].Label),
			})
		}
	}
	_ = t.dbService.SaveDiscountedProducts(append(productsInBonus, productsNotInBonus...))

	var notInBonusFields []*discordgo.MessageEmbedField
	if len(productsNotInBonus) != 0 {
		for _, prod := range productsNotInBonus {
			notInBonusFields = append(notInBonusFields, &discordgo.MessageEmbedField{
				Name: fmt.Sprintf("%s (%s)", prod.FriendlyName, prod.ApiName),
			})
		}
	}

	t.loggerService.Info("Sending message to Discord", nil)
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
	_, err = t.session.ChannelMessageSendEmbeds(os.Getenv("DISCORD_CHANNEL_ID"), embeds)
	if err != nil {
		t.loggerService.Error("Error sending message", map[string]interface{}{"error": err.Error()})
		return
	}
}

func (t *Thomas) Close() {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Thomas is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Thomas session.
	t.session.Close()
}
