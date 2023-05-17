package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	// Read a file from leveldb
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	chatIDs := []int64{}

	// Retrieve data from leveldb
	bytes, err := db.Get([]byte(ChattIDObjectName), nil)
	if err == nil {
		// Unmarshal data
		var tmp ListChatIDs
		err = json.Unmarshal(bytes, &tmp)
		if err != nil {
			log.Panic(err)
		}
		// Set chatIDs
		chatIDs = tmp.ChatIDs
	} else {
		log.Printf("No chatIDs found in db")
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Get arguments & split to strings
		args := strings.Split(update.Message.CommandArguments(), " ")
		fmt.Printf("args: %v, len: %v\n", args, len(args))

		// Get data from db
		chatID := update.Message.Chat.ID
		scores := map[string]int{}

		if !contains(chatIDs, chatID) {
			chatIDs = append(chatIDs, chatID)
			bytes, _ := json.Marshal(&ListChatIDs{chatIDs})
			err := db.Put([]byte(ChattIDObjectName), bytes, nil)
			if err != nil {
				continue
			}
		} else {
			bytes, err := db.Get([]byte(ScoreboardObjectName+fmt.Sprint(chatID)), nil)
			if err == nil {
				// Unmarshal data
				var tmp Scoreboard
				err = json.Unmarshal(bytes, &tmp)
				if err != nil {
					continue
				}
				// Set chatIDs
				scores = tmp.Scores
			}
		}
		startScores := cloneMap(scores)

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(chatID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "Available commands:\n" +
				"/init <name_1> <name_2> <score> (default: 0) \n" +
				"/remove <name_1> <name_2>\n" +
				"/add <name_1> <name_2> <score> (default: 1) \n" +
				"/sub <name_1> <name_2> <score> (default: 1) \n" +
				"/reset\n" +
				"/show"

		case "init":
			if len(args) < 1 || args[0] == "" {
				msg.Text = "You need to provide a name"
			} else {
				usernames, score, _ := parseInput(args, 0)
				fmt.Printf("usernames, score: %v, %v\n", usernames, score)
				initializeScores(scores, usernames, score)
				fmt.Printf("scores, startScores: %v, %v\n", scores, startScores)
				msg.Text = diffMaps(startScores, scores)
			}
		case "remove":
			if len(args) < 1 || args[0] == "" {
				msg.Text = "You need to provide a name"
			} else {
				usernames, _, _ := parseInput(args, 0)
				removeUsers(scores, usernames)
				msg.Text = diffMaps(startScores, scores)
			}
		case "add":
			if len(args) < 1 || args[0] == "" {
				msg.Text = "You need to provide a name"
			} else {
				usernames, score, _ := parseInput(args, 1)
				addScores(scores, usernames, score)
				msg.Text = diffMaps(startScores, scores)
			}
		case "sub":
			if len(args) < 1 || args[0] == "" {
				msg.Text = "You need to provide a name"
			} else {
				usernames, score, _ := parseInput(args, 1)
				subScores(scores, usernames, score)
				msg.Text = diffMaps(startScores, scores)
			}
		case "reset":
			scores = map[string]int{}
			msg.Text = diffMaps(startScores, scores)
		case "show":
			base := 0
			inputN := false
			if len(args) >= 1 && args[0] != "" {
				_, base, inputN = parseInput(args, 0)
			}
			msg.Text = showScores(scores, base, inputN)
		default:
			msg.Text = "I don't know that command"
		}

		bytes, _ := json.Marshal(&Scoreboard{scores})
		db.Put([]byte(ScoreboardObjectName+fmt.Sprint(chatID)), bytes, nil)

		fmt.Printf("Msg: %s\n", msg.Text)

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
