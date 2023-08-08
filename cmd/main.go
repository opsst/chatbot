package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	route := echo.New()

	bot, err := linebot.New("61d3f3790119fe63c4eabd5ed60d3dc6", "K12cxfkZ8VKnXizjw6TMeh0uDsSzwjlAaxf75H9RzE57JeonPZBTDOJ8J1s/GfN/ndCcMwicQcKRrxJ2YLDKEFD6Kl/IBbnxOkRZ82BADKlrlgTWxRaTfEL3ZmWi66mzU++jiVdOYX7Gzkule+5RrwdB04t89/1O/w1cDnyilFU=")

	if err != nil {
		log.Fatal(err)
	}

	route.POST("/callback", func(c echo.Context) error {

		// body, err := ioutil.ReadAll(c.Request().Body)
		// if err != nil {
		// 	// ...
		// }
		// decoded, err := base64.StdEncoding.DecodeString(c.Request().Header.Get("x-line-signature"))
		// if err != nil {
		// 	// ...
		// }
		// hash := hmac.New(sha256.New, []byte("<channel secret>"))
		// hash.Write(body)
		// // Compare decoded signature and `hash.Sum(nil)` by using `hmac.Equal`

		events, err := bot.ParseRequest(c.Request())
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				return c.JSON(http.StatusBadRequest, err)
			} else {
				return c.JSON(http.StatusInternalServerError, err)
			}

		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
		return c.NoContent(http.StatusOK)
	})

	go RunServer(route, 3020)
	Shutdown(route)
	// router := route.New

}
func Shutdown(router *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	if err := router.Shutdown(context.Background()); err != nil {

		log.Fatal(err)
	}
}

func RunServer(router *echo.Echo, port int) {
	startPort := fmt.Sprintf(":%d", port)
	router.Logger.Fatal(router.Start(startPort))
}
