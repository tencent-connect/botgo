package apitest

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
)

var conf struct {
	AppID uint64 `yaml:"appid"`
	Token string `yaml:"token"`
}
var botToken *token.Token
var api openapi.OpenAPI

var (
	testGuildID   = "3326534247441079828" // replace your guild id
	testChannelID = "1595028"             // replace your channel id
	testMessageID = `08e092eeb983afef9e0110f9bb5d1a1231343431313532313836373838333234303420801e
28003091c4bb02380c400c48d8a7928d06` // replace your channel id
	testRolesID            = `10054557`                             // replace your roles id
	testMemberID           = `1201318637970874066`                  // replace your member id
	testMarkdownTemplateID = 1231231231231231                       // replace your markdown template id
	testInteractionD       = "e924431f-aaed-4e78-8375-9295b1f1e0ef" // replace your interaction id
	ctx                    context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Println("read conf failed")
		os.Exit(1)
	}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(conf)

	botToken = token.BotToken(conf.AppID, conf.Token)
	api = botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	if err := botToken.InitToken(context.Background()); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}
