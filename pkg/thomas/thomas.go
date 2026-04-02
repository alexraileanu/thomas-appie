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

const (
	colorBonus    = 0xff7900
	colorNotBonus = 0xff0000
	maxEmbedFields = 25
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

	if err = session.Open(); err != nil {
		return nil, err
	}

	return &Thomas{session: session, dbService: dbService, loggerService: loggerService, config: config}, nil
}

func (t *Thomas) Go() {
	products, err := t.dbService.GetProducts()
	if err != nil {
		t.loggerService.Error("Error fetching products from the db", map[string]interface{}{"error": err.Error()})
		panic(err)
	}

	a := appie.New(t.loggerService, t.config.Appie)

	t.loggerService.Info("Checking products", map[string]interface{}{"count": len(products)})
	inBonus, notInBonus, err := a.PerformProductsCheck(products)
	if err != nil {
		panic(err)
	}

	_ = t.dbService.SaveDiscountedProducts(append(inBonus, notInBonus...))

	inBonusFields := productFields(inBonus, true)
	notInBonusFields := productFields(notInBonus, false)

	var embeds []*discordgo.MessageEmbed
	embeds = append(embeds, buildEmbeds(inBonusFields, "Products that are in bonus this week at the Appie", colorBonus)...)
	embeds = append(embeds, buildEmbeds(notInBonusFields, "Products that aren't in bonus this week at the Appie", colorNotBonus)...)

	t.loggerService.Info("Sending Discord message", map[string]interface{}{"in_bonus": len(inBonus), "not_in_bonus": len(notInBonus)})
	if _, err = t.session.ChannelMessageSendEmbeds(os.Getenv("DISCORD_CHANNEL_ID"), embeds); err != nil {
		t.loggerService.Error("Error sending Discord message", map[string]interface{}{"error": err.Error()})
		return
	}
	t.loggerService.Info("Discord message sent", nil)
}

func (t *Thomas) Close() {
	t.loggerService.Info("Thomas is running, press CTRL-C to exit", nil)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	t.session.Close()
}

// productFields builds Discord embed fields for a list of products.
// For bonus products it includes the discount description; for non-bonus only the name.
func productFields(products []appie.Product, withDescription bool) []*discordgo.MessageEmbedField {
	fields := make([]*discordgo.MessageEmbedField, 0, len(products))
	for _, p := range products {
		field := &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%s (%s)", p.FriendlyName, p.ApiName),
		}
		if withDescription && len(p.DiscountedProducts) > 0 {
			field.Value = fmt.Sprintf("%s %s", p.DiscountedProducts[0].Description, p.DiscountedProducts[0].Label)
		}
		fields = append(fields, field)
	}
	return fields
}

// buildEmbeds splits fields into chunks of maxEmbedFields (Discord's per-embed limit)
// and wraps each chunk in a MessageEmbed with the given title and color.
func buildEmbeds(fields []*discordgo.MessageEmbedField, title string, color int) []*discordgo.MessageEmbed {
	var embeds []*discordgo.MessageEmbed
	for i := 0; i < len(fields); i += maxEmbedFields {
		end := i + maxEmbedFields
		if end > len(fields) {
			end = len(fields)
		}
		embeds = append(embeds, &discordgo.MessageEmbed{
			Color:  color,
			Title:  title,
			Fields: fields[i:end],
		})
	}
	return embeds
}
