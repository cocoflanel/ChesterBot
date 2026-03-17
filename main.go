package main

import (
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	//restTest := restGet("https://api.restful-api.dev/objects")
	//log.Println(restTest)
	runBot()
	/*link, err := restPutImg("https://media.discordapp.net/attachments/541836238471299114/1473448726696558743/HArtXgCbYAAtan7.png?ex=69963f8f&is=6994ee0f&hm=b8edb7d0ba5408ad5f3e3e52ab622b435546f79995e6f6b536ee69e2ba119fe4&")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully uploaded img at: " + link)*/
}

func runBot() {

	//import a discord auth token from bot.key and instantiate a new bot with it
	authToken := loadFile("bot.key")
	sess, err := discordgo.New("Bot " + authToken)
	errorCheck(err)
	//if needed, permission int is 116736

	//add message handler, intents, and open the session
	sess.AddHandler(helloMesages)
	sess.Identify.Intents = discordgo.IntentsGuildMessages
	err = sess.Open()
	errorCheck(err)
	log.Println("The bot is listening")

	//code to hold the thread until interrupt is sent
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

// Primary message handler
func helloMesages(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	//Chester Commands
	if strings.HasPrefix(m.Content, "chester") || strings.HasPrefix(m.Content, "Chester") {
		log.Println(m.Content)

		//check if there's any attatchments and loop over them
		if len(m.Attachments) > 0 {
			for _, attatchment := range m.Attachments {
				//for any image attatchements, log and print out the image proxy url.
				if strings.HasPrefix(attatchment.ContentType, "image") {
					log.Printf("image found: %+v\n", attatchment.ProxyURL)
					s.ChannelMessageSend(m.ChannelID, "Oh boy yummy image for me! "+attatchment.ProxyURL)
					restPutImg(attatchment.ProxyURL)
				} else {
					//for non-image attatchments
					log.Println(attatchment.ContentType)
					s.ChannelMessageSend(m.ChannelID, "*gorps your files* Oh that's not quite as good as an image.")
				}
			}
		} else {
			//for messages with no attatchments
			s.ChannelMessageSend(m.ChannelID, "Hey that's me! Gimme an image I've been so good :D")
		}
	} else if strings.Contains(strings.ToLower(m.Content), "good boy") {

		s.ChannelMessageSend(m.ChannelID, getAGoodBoy())

	}
}

func getAGoodBoy() string {
	//Good Boy GIFs, replies with any of the following
	goodBoyGIF := ""
	switch rand.IntN(8) {
	case 1:
		goodBoyGIF = "https://c.tenor.com/_4xCiEhhoZsAAAAd/tenor.gif"
	case 2:
		goodBoyGIF = "https://c.tenor.com/4jSFH4ktsHsAAAAC/tenor.gif"
	case 3:
		goodBoyGIF = "https://c.tenor.com/vXG7hQc33IoAAAAC/tenor.gif"
	case 4:
		goodBoyGIF = "https://c.tenor.com/cDJj3LEw0UIAAAAd/tenor.gif"
	case 5:
		goodBoyGIF = "https://c.tenor.com/ts6NhvDRU1AAAAAd/tenor.gif"
	case 6:
		goodBoyGIF = "https://c.tenor.com/qIxZz0K4yVQAAAAC/tenor.gif"
	case 7:
		goodBoyGIF = "https://tenor.com/view/straight-face-inuyasha-inuyahsa-gif-21752958"
	default:
		goodBoyGIF = "https://c.tenor.com/KY0JkwS42xEAAAAd/tenor.gif"

	}
	return goodBoyGIF
}

func loadFile(filename string) string {
	data, err := os.ReadFile(filename)
	errorCheck(err)
	return (string(data))
}

// basic error check function to throw panic when nothing more is required
func errorCheck(e error) {
	if e != nil {
		log.Panic(e)
	}
}
