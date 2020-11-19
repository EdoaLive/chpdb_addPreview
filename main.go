package main

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var (
	bot *tgbotapi.BotAPI
)

func main() {
	var err error
	// verySecretToken is in secret.go, do not commit it!
	bot, err = tgbotapi.NewBotAPI(verySecretToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	re := regexp.MustCompile(`https://ift\.tt/\w+`)

	re2 := regexp.MustCompile(`href="(.+?)".* data-preloader="adp_CometPhotoRootQueryRelayPreloader_&#123;N&#125;"`)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			// parseMessage(update.Message)

			match := re.MatchString(update.Message.Text)
			if match {

				url1 := re.FindString(update.Message.Text)
				sm := getImgUrl(url1, re2)

				msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, sm)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}

}

func getImgUrl(url1 string, re2 *regexp.Regexp) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sm := re2.FindSubmatch(body)
	return html.UnescapeString(string(sm[1]))
}

func getRedirect(s string) string {
	req, err := http.NewRequest("GET", s, nil)
	if err != nil {
		panic(err)
	}

	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	response, err := client.Do(req)
	if err != nil {
		if response.StatusCode == http.StatusMovedPermanently { //status code 302
			redir, _ := response.Location()
			return redir.String()
		}
	}
	return s
}
