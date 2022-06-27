package main

import (
	"fmt"
	"log"
	"strconv"

	// "math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var version string = "2.0.0"

func goDotEnvVariable(key string) string {
	// Load .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Return value from key provided.
	return os.Getenv(key)
}

func main() {

	// Grab bot token env var.
	botToken := goDotEnvVariable("BOT_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	// dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	author := m.Author.Username
	authorID := m.Author.ID

	guildID := m.Message.GuildID

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Grab message content from guild.
	content := m.Content

	if strings.Contains(content, "!rpshelp") {

		// Title
		commandHelpTitle := "Looks like you need a hand. Check out my goodies below... \n \n"

		// Notes
		// note0 := "- Commands are case-sensitive. They must be in lower-case :) . \n"
		note1 := "- Dev: Narsiq#5638. DM me for requests/questions/love. \n"

		// Commands
		commandHelp := "❔  !rpshelp : Provides a list of my commands. \n"
		commandChallenge := "🦶🏽  !rps @User : Challenge tagged user. \n"
		commandSite := "🔗  !rpssite : Link to the RPS website \n"
		commandSupport := "✨  !rpssupport : Link to the RPS Patreon. \n"
		commandVersion := "🤖  !rpsversion : Current RPS version. \n"

		// Build help message
		message := "Whats up " + author + "\n \n" + commandHelpTitle + "NOTES: \n \n" + note1 + "\n" + "COMMANDS: \n \n" + commandHelp + commandChallenge + "\n" + "OTHER: \n \n" + commandSite + commandSupport + commandVersion + "\n \n" + "https://www.patreon.com/BotVoteTo"

		// Reply to help request with build message above.
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!rpssite") {
		// Create website message
		message := "Here ya go " + author + "..." + "\n" + "https://discordbots.dev/"

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!rpssupport") {
		// Create support message
		message := "Thanks for thinking of me " + author + " 💖." + "\n" + "https://www.patreon.com/discordbotsdev"

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!rpsversion") {
		// Create version message
		message := "RockPaperScissors is currently running version " + version

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!rpsstats") {
		// TODO: This will need to be updated to iterate through
		// all shards once the bot joins 1,000 servers.
		guilds := s.State.Ready.Guilds
		fmt.Println(len(guilds))
		guildCount := len(guilds)

		guildCountStr := strconv.Itoa(guildCount)

		// // Build start vote message
		message := "RockPaperScissors is currently on " + guildCountStr + " servers. Such wow!"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!rps") {
		// Trim bot command from string to grab User tagged
		trimCommand := strings.TrimPrefix(content, "!rps ")
		targetUserID := strings.Trim(trimCommand, "<@!>")

		// GuildMember(guildID, userID string) (st *Member, err error)
		targetUser, err := s.GuildMember(guildID, targetUserID)
		if err != nil {
			fmt.Println(err)
		}

		// Create targetUserNickname for target identification in challenge message
		targetUserNickname := targetUser.Nick
		if targetUserNickname == "" {
			targetUserNickname = targetUser.User.Username
		}

		// Build start challenge message
		startMessage := author + " is challenging " + targetUserNickname + " to a game of rock paper scissors... \n \n" + "Accept or decline by reacting below in the next 30 seconds. Only " + targetUserNickname + "'s reaction will trigger a DM to play."

		// Send start challenge message
		challengeMessage, err := s.ChannelMessageSendReply(m.ChannelID, startMessage, m.Reference())
		if err != nil {
			fmt.Println(err)
		}

		// Add a yes reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, challengeMessage.ID, "✔️")
		if err != nil {
			fmt.Println(err)
		}

		// Add a no reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, challengeMessage.ID, "❌")
		if err != nil {
			fmt.Println(err)
		}

		// BELOW CODE TO NEW FUNC

		// timeout := time.After(30 * time.Second)
		// ticker := time.Tick(500 * time.Millisecond)

		result := checkReactions(s, m, challengeMessage.ID, targetUserID, author, authorID, targetUserNickname)

		if !result {
			authorAt := "<@!" + authorID + ">"
			timeoutMessage := "Bummer, " + authorAt + ". Challenge period ended without a reaction from " + targetUserNickname + " 😞"
			_, err := s.ChannelMessageSendReply(m.ChannelID, timeoutMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		}

	}

}

func checkReactions(s *discordgo.Session, m *discordgo.MessageCreate, challengeMessageID string, targetUserID string, author string, authorID string, targetUserNickname string) bool {
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				reactions, err := s.MessageReactions(m.ChannelID, challengeMessageID, "✔️", 100, "", "")
				if err != nil {
					fmt.Println(err)
				}

				for _, user := range reactions {
					if user.ID == targetUserID {
						fmt.Println("Target User reacted! Yay")
						sendChallengeDMs(s, m, author, authorID, targetUserID, targetUserNickname)
						done <- true
						return
					}
				}

				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(30 * time.Second)
	ticker.Stop()
	done <- true

	// Call challenge expiration message fun
	fmt.Println("Ticker stopped")

	return false
}

func sendChallengeDMs(s *discordgo.Session, m *discordgo.MessageCreate, author string, authorID string, targetUserID string, targetUserNickname string) {
	// Create DM to author with reaction options
	authorChannel, errAuthorChannel := s.UserChannelCreate(authorID)
	if errAuthorChannel != nil {
		fmt.Println(errAuthorChannel)
	}

	// Send request move message to author
	challengerMessage := "Time to clutch " + author + ". Select your move below."
	authorChannelMessage, errAuthorChannelMessage := s.ChannelMessageSend(authorChannel.ID, challengerMessage)
	if errAuthorChannelMessage != nil {
		fmt.Println(errAuthorChannelMessage)
	}

	// Add reactions to Author's DM challenge message
	AddMessageReactions(s, m, authorChannel.ID, authorChannelMessage.ID)

	// Create DM to target with reaction options
	targetChannel, errTargetChannel := s.UserChannelCreate(targetUserID)
	if errTargetChannel != nil {
		fmt.Println(errTargetChannel)
	}

	// Send request move message to author
	targetMessage := "Time to clutch " + targetUserNickname + ". Select your move below."
	targetChannelMessage, errTargetChannelMessage := s.ChannelMessageSend(targetChannel.ID, targetMessage)
	if errTargetChannelMessage != nil {
		fmt.Println(errTargetChannelMessage)
	}

	// Add reactions to Target's DM challenge message
	AddMessageReactions(s, m, targetChannel.ID, targetChannelMessage.ID)

	// Make a channel to recieve a value from the
	// checkMoves func for both author and target moves
	channelTarget := make(chan string)
	channelAuthor := make(chan string)

	// Check the moves from both author and target.
	go checkMoves(s, m, targetChannel.ID, targetChannelMessage.ID, channelTarget)
	go checkMoves(s, m, authorChannel.ID, authorChannelMessage.ID, channelAuthor)

	// Recieve values from channel sent from checkMoves calls
	moveTarget, moveAuthor := <-channelTarget, <-channelAuthor

	// Build the "@" value for author and target.
	authorAt := "<@!" + authorID + ">"
	targetAt := "<@!" + targetUserID + ">"

	// Check to if author and/or target did not respond to send message.
	if moveTarget == "" && moveAuthor == "" {
		timeoutMessage := "Bummer, challenge period ended without a move from both participants 😞"
		_, err := s.ChannelMessageSendReply(m.ChannelID, timeoutMessage, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if moveTarget != "" && moveAuthor == "" {
		timeoutMessage := "Bummer, challenge period ended without a move from" + authorAt + " 😞"
		_, err := s.ChannelMessageSendReply(m.ChannelID, timeoutMessage, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if moveTarget == "" && moveAuthor != "" {
		timeoutMessage := "Bummer, challenge period ended without a move from" + targetAt + " 😞"
		_, err := s.ChannelMessageSendReply(m.ChannelID, timeoutMessage, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if moveTarget != "" && moveAuthor != "" {
		checkMovesAndNotify(s, m, moveAuthor, moveTarget, author, targetUserID, targetUserNickname)
	}
}

func AddMessageReactions(s *discordgo.Session, m *discordgo.MessageCreate, recipientChannelID string, recipientChannelMessageID string) {
	err := s.MessageReactionAdd(recipientChannelID, recipientChannelMessageID, "⛰️")
	if err != nil {
		fmt.Println(err)
	}

	// Add paper reaction to request message
	err = s.MessageReactionAdd(recipientChannelID, recipientChannelMessageID, "🧻")
	if err != nil {
		fmt.Println(err)
	}

	// Add scissors reaction to request message
	err = s.MessageReactionAdd(recipientChannelID, recipientChannelMessageID, "✂️")
	if err != nil {
		fmt.Println(err)
	}
}

func checkMoves(s *discordgo.Session, m *discordgo.MessageCreate, challengeDMChannelID string, challengeDMMessageID string, c chan string) {
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				reactionRock, err := s.MessageReactions(challengeDMChannelID, challengeDMMessageID, "⛰️", 2, "", "")
				if err != nil {
					fmt.Println(err)
				}

				if len(reactionRock) > 1 {
					c <- "Rock"
					done <- true
					return
				}

				reactionPaper, err := s.MessageReactions(challengeDMChannelID, challengeDMMessageID, "🧻", 2, "", "")
				if err != nil {
					fmt.Println(err)
				}

				if len(reactionPaper) > 1 {
					c <- "Paper"
					done <- true
					return
				}

				reactionScissors, err := s.MessageReactions(challengeDMChannelID, challengeDMMessageID, "✂️", 2, "", "")
				if err != nil {
					fmt.Println(err)
				}

				if len(reactionScissors) > 1 {
					c <- "Scissors"
					done <- true
					return
				}

				fmt.Println("Tick: ", t)
			}
		}
	}()

	// Stuff related to waiting
	time.Sleep(15 * time.Second)
	ticker.Stop()
	done <- true

	// Send empty string through the channel
	c <- ""
}

func checkMovesAndNotify(s *discordgo.Session, m *discordgo.MessageCreate, moveAuthor string, moveTarget string, author string, targetUserID string, targetNickname string) {
	authorID := m.Author.ID
	authorAt := "<@!" + authorID + ">"
	targetAt := "<@!" + targetUserID + ">"

	// Rock vs Rock
	if moveAuthor == "Rock" && moveTarget == "Rock" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Rock" + " 🤡" + "\n" + targetAt + ": Rock" + " 🤡" + "\n \n"
		resultsMessage := "Welp, this is akward... " + author + " and " + targetNickname + " actually tied."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Rock vs Paper
	if moveAuthor == "Rock" && moveTarget == "Paper" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Rock" + "\n" + targetAt + ": Paper" + " 👑" + "\n \n"
		resultsMessage := author + " got rekt by " + targetNickname + ". Imagine..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Rock vs Scissors
	if moveAuthor == "Rock" && moveTarget == "Scissors" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Rock" + " 👑" + "\n" + targetAt + ": Scissors" + "\n \n"
		resultsMessage := author + " kinda clapped " + targetNickname + " ngl..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Paper vs Paper
	if moveAuthor == "Paper" && moveTarget == "Paper" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Paper" + " 🤡" + "\n" + targetAt + ": Paper" + " 🤡" + "\n \n"
		resultsMessage := "Welp, this is akward... " + author + " and " + targetNickname + " actually tied."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Paper vs Rock
	if moveAuthor == "Paper" && moveTarget == "Rock" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Paper" + " 👑" + "\n" + targetAt + ": Rock" + "\n \n"
		resultsMessage := author + " kinda clapped " + targetNickname + " ngl..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Paper vs Scissors
	if moveAuthor == "Paper" && moveTarget == "Scissors" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Paper" + "\n" + targetAt + ": Scissors" + " 👑" + "\n \n"
		resultsMessage := author + " got rekt by " + targetNickname + ". Imagine..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Scissors vs Scissors
	if moveAuthor == "Scissors" && moveTarget == "Scissors" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Scissors" + " 🤡" + "\n" + targetAt + ": Scissors" + " 🤡" + "\n \n"
		resultsMessage := "Welp, this is akward... " + author + " and " + targetNickname + " actually tied."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Scissors vs Rock
	if moveAuthor == "Scissors" && moveTarget == "Rock" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Scissors" + "\n" + targetAt + ": Rock" + " 👑" + "\n \n"
		resultsMessage := author + " got rekt by " + targetNickname + ". Imagine..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Scissors vs Paper
	if moveAuthor == "Scissors" && moveTarget == "Paper" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Scissors" + " 👑" + "\n" + targetAt + ": Paper" + "\n \n"
		resultsMessage := author + " kinda clapped " + targetNickname + " ngl..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}
}
