package testcase

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/token"
	"gopkg.in/yaml.v2"
)

var (
	testGuildID   = "3326534247441079828" // replace your guild id
	testChannelID = "116482"              // replace your channel id
	testMessageID = `08e092eeb983afef9e0110f9bb5d1a1231343431313532313836373838333234303420801e
28003091c4bb02380c400c48d8a7928d06`  // replace your channel id
	ctx context.Context
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

	os.Exit(m.Run())
}
