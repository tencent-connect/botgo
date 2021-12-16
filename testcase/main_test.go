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
	testGuildID   = "7195342413866929087" // replace your guild id
	testChannelID = "1611157"             // replace your channel id
	testMessageID = `08f1908d979791ee521095ab621a12313434313135323138363
73133393436313220801e280030f1fa8d703813400c489cb5e18d06`  // replace your channel id
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
