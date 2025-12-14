package thomas

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexraileanu/thomas-appie/pkg/config"
	"github.com/alexraileanu/thomas-appie/pkg/logger"

	"github.com/bwmarrin/discordgo"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
	"github.com/alexraileanu/thomas-appie/pkg/db"
)

type Thomas struct {
	session       *discordgo.Session
	dbService     *db.Service
	loggerService *logger.Service

	config config.Config
}

func New(dbService *db.Service, loggerService *logger.Service, config config.Config) (*Thomas, error) {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		return nil, err
	}

	err = session.Open()
	if err != nil {
		return nil, err
	}

	return &Thomas{session: session, dbService: dbService, loggerService: loggerService, config: config}, nil
}

func (t *Thomas) Go() {
	t.loggerService.Info("Fetching products from the db", nil)
	products, err := t.dbService.GetProducts()
	if err != nil {
		t.loggerService.Error("Error fetching products from the db", map[string]interface{}{"error": err.Error()})
		panic(err)
	}

	a := appie.New(t.loggerService, t.config.Appie)

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

	inBonus := t.fixFieldsForDiscord(inBonusFields)
	notInBonus := t.fixFieldsForDiscord(notInBonusFields)

	var embeds []*discordgo.MessageEmbed
	embeds = append(embeds, t.fixEmbedsForDiscord(inBonus, "Products that are in bonus this week at the Appie", 0xff7900)...)
	embeds = append(embeds, t.fixEmbedsForDiscord(notInBonus, "Products that aren't in bonus this week at the Appie", 0xff0000)...)

	_, err = t.session.ChannelMessageSendEmbeds(os.Getenv("DISCORD_CHANNEL_ID"), embeds)
	if err != nil {
		t.loggerService.Error("Error sending message", map[string]interface{}{"error": err.Error()})
		return
	}
}

func (t *Thomas) fixFieldsForDiscord(fields []*discordgo.MessageEmbedField) [][]*discordgo.MessageEmbedField {
	var fixedFields [][]*discordgo.MessageEmbedField
	if len(fields) > 25 {
		for i := 0; i < len(fields); i += 25 {
			end := i + 25
			if end > len(fields) {
				end = len(fields)
			}
			fixedFields = append(fixedFields, fields[i:end])
		}
	} else {
		fixedFields = append(fixedFields, fields)
	}

	return fixedFields
}

func (t *Thomas) fixEmbedsForDiscord(embeds [][]*discordgo.MessageEmbedField, title string, color int) []*discordgo.MessageEmbed {
	var fixedEmbeds []*discordgo.MessageEmbed
	if len(embeds) > 1 {
		for _, bonus := range embeds {
			fixedEmbeds = append(fixedEmbeds, &discordgo.MessageEmbed{
				Color:  color,
				Title:  title,
				Fields: bonus,
			})
		}
	} else {
		fixedEmbeds = append(fixedEmbeds, &discordgo.MessageEmbed{
			Color:  color,
			Title:  title,
			Fields: embeds[0],
		})
	}

	return fixedEmbeds
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
