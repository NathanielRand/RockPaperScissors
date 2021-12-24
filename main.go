package main

import (
	"fmt"
	"log"
	// "math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var version string = "1.0.0"

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
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	guildID := m.Message.GuildID

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Grab message content from guild.
	content := m.Content

	if strings.Contains(content, "!rpshelp") {
		// Build help message
		author := m.Author.Username

		// Title
		commandHelpTitle := "Looks like you need a hand. Check out my goodies below... \n \n"

		// Notes
		note0 := "- Commands are case-sensitive. They must be in lower-case :) . \n"
		note1 := "- Dev: Narsiq#5638. DM me for requests/questions/love. \n"

		// Commands
		commandHelp := "❔  !rpshelp : Provides a list of my commands. \n"
		commandChallenge := "🦶🏽  !rps @User : Challenge tagged user. \n"
		commandSite := "🔗  !rpssite : Link to the RPS website \n"
		commandSupport := "✨  !rpssupport : Link to the RPS Patreon. \n"
		commandVersion := "🤖  !rpsversion : Current RPS version. \n"

		// Build help message
		message := "Whats up " + author + "\n \n" + commandHelpTitle + "NOTES: \n \n" + note0 + note1 + "\n" + "COMMANDS: \n \n" + commandHelp + commandChallenge + "\n" + "OTHER: \n \n" + commandSite + commandSupport + commandVersion + "\n \n" + "https://www.patreon.com/BotVoteTo"

		// Reply to help request with build message above.
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!rpssite") {
		// Build start vote message
		author := m.Author.Username
		message := "Here ya go " + author + "..." + "\n" + "https://discordbots.dev/"

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!rpssupport") {
		// Build start vote message
		author := m.Author.Username
		message := "Thanks for thinking of me " + author + " 💖." + "\n" + "https://www.patreon.com/BotVoteTo"

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!rpsversion") {
		// Build start vote message
		message := "VoteBot is currently running version " + version

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!rps") {
		// Trim bot command from string to grab User tagged
		trimmed := strings.TrimPrefix(content, "!rps ")
		trimmedUser := strings.Trim(trimmed, "<@!>")

		// GuildMember(guildID, userID string) (st *Member, err error)
		targetMember, err := s.GuildMember(guildID, trimmedUser)
		if err != nil {
			fmt.Println(err)
		}

		// Create memberNickname for target identification in challenge message
		targetNickname := targetMember.Nick
		if targetNickname == "" {
			targetNickname = targetMember.User.Username
		}

		// Build start challenge message
		author := m.Author.Username
		authorID := m.Author.ID
		startMessage := author + " is challenging " + targetNickname + " to a game of rock paper scissors... \n \n" + "Accept or decline by reacting below in the next 30 seconds. Only " + targetNickname + "'s reaction will trigger a DM to play."

		// Send start challenge message
		challengeMessage, err := s.ChannelMessageSendReply(m.ChannelID, startMessage, m.Reference())
		if err != nil {
			fmt.Println(err)
		}

		// Add yes reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, challengeMessage.ID, "✔️")
		if err != nil {
			fmt.Println(err)
		}

		// Add no reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, challengeMessage.ID, "❌")
		if err != nil {
			fmt.Println(err)
		}

		// challengerMove, targetMove := checksendReactions(s, m, trimmed, trimmedUser, authorID, challengeMessage.ID, author)

		// fmt.Println("challengerMove: ", challengerMove)
		// fmt.Println("targetMove: ", targetMove)

		yes, err2 := s.MessageReactions(m.ChannelID, challengeMessage.ID, "✔️", 100, "", "")
		if err2 != nil {
			fmt.Println(err2)
		}

		// Count no reactions from vote message
		no, err3 := s.MessageReactions(m.ChannelID, challengeMessage.ID, "❌", 100, "", "")
		if err3 != nil {
			fmt.Println(err3)
		}

		fmt.Println(yes)
		fmt.Println(no)

		timeout := time.After(30 * time.Second)
		ticker := time.Tick(500 * time.Millisecond)

		// Keep trying until we're timed out or get a result/error
		for {
			select {
			// Got a timeout! fail with a timeout error
			case <-timeout:
				authorID := m.Author.ID
				authorAt := "<@!" + authorID + ">"
				timeoutMessage := "Bummer, " + authorAt + ". Challenge period ended without a reaction from " + targetNickname + " 😞"
				_, err := s.ChannelMessageSendReply(challengeMessage.ID, timeoutMessage, m.Reference())
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Challenge period ended without a reaction from the target ID: ", targetNickname)
				return
			// Got a tick, we should check on reactions
			case <-ticker:
				yes, err2 = s.MessageReactions(m.ChannelID, challengeMessage.ID, "✔️", 100, "", "")
				if err2 != nil {
					fmt.Println(err2)
				}

				for _, yes := range yes {

					fmt.Printf("Name: %s ID: %s\n", yes.Username, yes.ID)

					if yes.ID == trimmedUser {

						// If accepeted by target return with DM to author and target with reaction options
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

						// Add rock reaction to request message
						err := s.MessageReactionAdd(authorChannel.ID, authorChannelMessage.ID, "⛰️")
						if err != nil {
							fmt.Println(err)
						}

						// Add paper reaction to request message
						err = s.MessageReactionAdd(authorChannel.ID, authorChannelMessage.ID, "🧻")
						if err != nil {
							fmt.Println(err)
						}

						// Add scissors reaction to request message
						err = s.MessageReactionAdd(authorChannel.ID, authorChannelMessage.ID, "✂️")
						if err != nil {
							fmt.Println(err)
						}

						// If accepeted by target return with DM to author and target with reaction options
						targetChannel, errTargetChannel := s.UserChannelCreate(trimmedUser)
						if errTargetChannel != nil {
							fmt.Println(err)
						}

						// Send request move message to author
						targetMessage := "Time to clutch " + targetNickname + ". Select your move below."
						targetChannelMessage, errTargetChannelMessage := s.ChannelMessageSend(targetChannel.ID, targetMessage)
						if errTargetChannelMessage != nil {
							fmt.Println(err)
						}

						// Add rock reaction to request message
						err = s.MessageReactionAdd(targetChannel.ID, targetChannelMessage.ID, "⛰️")
						if err != nil {
							fmt.Println(err)
						}

						// Add paper reaction to request message
						err = s.MessageReactionAdd(targetChannel.ID, targetChannelMessage.ID, "🧻")
						if err != nil {
							fmt.Println(err)
						}

						// Add scissors reaction to request message
						err = s.MessageReactionAdd(targetChannel.ID, targetChannelMessage.ID, "✂️")
						if err != nil {
							fmt.Println(err)
						}

						timeout2 := time.After(30 * time.Second)
						ticker2 := time.Tick(500 * time.Millisecond)

						for {

							select {
							// Got a timeout! fail with a timeout error

							case <-timeout2:
								timeoutMessage := "Bummer, challenge period ended without a move from both participants 😞"
								_, err := s.ChannelMessageSend(authorChannel.ID, timeoutMessage)
								if err != nil {
									fmt.Println(err)
								}
								_, err = s.ChannelMessageSend(targetChannel.ID, timeoutMessage)
								if err != nil {
									fmt.Println(err)
								}
								_, err = s.ChannelMessageSendReply(m.ChannelID, timeoutMessage, m.Reference())
								if err != nil {
									fmt.Println(err)
								}

								return
							// Got a tick, we should check on reactions
							case <-ticker2:
								// Check reactions from the challenger
								challengerRock, err2 := s.MessageReactions(authorChannel.ID, authorChannelMessage.ID, "⛰️", 100, "", "")
								if err2 != nil {
									fmt.Println(err)
								}
								challengerPaper, err2 := s.MessageReactions(authorChannel.ID, authorChannelMessage.ID, "🧻", 100, "", "")
								if err2 != nil {
									fmt.Println(err)
								}
								challengerScissors, err2 := s.MessageReactions(authorChannel.ID, authorChannelMessage.ID, "✂️", 100, "", "")
								if err2 != nil {
									fmt.Println(err)
								}

								// Check reactions from target
								targetRock, err3 := s.MessageReactions(targetChannel.ID, targetChannelMessage.ID, "⛰️", 100, "", "")
								if err2 != nil {
									fmt.Println(err3)
								}
								targetPaper, err4 := s.MessageReactions(targetChannel.ID, targetChannelMessage.ID, "🧻", 100, "", "")
								if err2 != nil {
									fmt.Println(err4)
								}
								targetScissors, err5 := s.MessageReactions(targetChannel.ID, targetChannelMessage.ID, "✂️", 100, "", "")
								if err2 != nil {
									fmt.Println(err5)
								}

								// Moves to be assigned based on the compare statements
								var challengerMove string
								var targetMove string

								// Check if reactions have adde value.
								if (len(challengerRock) > 1 || len(challengerPaper) > 1 || len(challengerScissors) > 1) && (len(targetRock) > 1 || len(targetPaper) > 1 || len(targetScissors) > 1) {
									if len(challengerRock) > 1 {
										if challengerRock[0].ID == authorID {
											challengerMove = "Rock"
										}
									}
									if len(challengerPaper) > 1 {
										if challengerPaper[0].ID == authorID {
											challengerMove = "Paper"
										}
									}
									if len(challengerScissors) > 1 {
										if challengerScissors[0].ID == authorID {
											challengerMove = "Scissors"
										}
									}
									if len(targetRock) > 1 {
										if targetRock[0].ID == trimmedUser {
											targetMove = "Rock"
										}
									}
									if len(targetPaper) > 1 {
										if targetPaper[0].ID == trimmedUser {
											targetMove = "Paper"
										}
									}
									if len(targetScissors) > 1 {
										if targetScissors[0].ID == trimmedUser {
											targetMove = "Scissors"
										}
									}
								}

								if challengerMove != "" && targetMove != "" {
									checkMovesAndNotify(s, m, challengerMove, targetMove, author, trimmed, targetNickname)
									return
								}
							}
						}

					}
				}
				// checkSomething() isn't done yet, but it didn't fail either, let's try again
			}
		}
	}

}

func checkMovesAndNotify(s *discordgo.Session, m *discordgo.MessageCreate, challengerMove string, targetMove string, author string, trimmed string, targetNickname string) {
	authorID := m.Author.ID
	authorAt := "<@!" + authorID + ">"

	// Rock vs Rock
	if challengerMove == "Rock" && targetMove == "Rock" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Rock" + " 🤡" + "\n" + trimmed + ": Rock" + " 🤡" + "\n \n"
		resultsMessage := "Welp, this is akward... " + author + " and " + targetNickname + " actually tied."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Rock vs Paper
	if challengerMove == "Rock" && targetMove == "Paper" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Rock" + "\n" + trimmed + ": Paper" + " 👑" + "\n \n"
		resultsMessage := author + " got rekt by " + targetNickname + ". Imagine..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Rock vs Scissors
	if challengerMove == "Rock" && targetMove == "Scissors" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Rock" + " 👑" + "\n" + trimmed + ": Scissors" + "\n \n"
		resultsMessage := author + " kinda clapped " + targetNickname + " ngl..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Paper vs Paper
	if challengerMove == "Paper" && targetMove == "Paper" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Paper" + " 🤡" + "\n" + trimmed + ": Paper" + " 🤡" + "\n \n"
		resultsMessage := "Welp, this is akward... " + author + " and " + targetNickname + " actually tied."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Paper vs Rock
	if challengerMove == "Paper" && targetMove == "Rock" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Paper" + " 👑" + "\n" + trimmed + ": Rock" + "\n \n"
		resultsMessage := author + " kinda clapped " + targetNickname + " ngl..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Paper vs Scissors
	if challengerMove == "Paper" && targetMove == "Scissors" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Paper" + "\n" + trimmed + ": Scissors" + " 👑" + "\n \n"
		resultsMessage := author + " got rekt by " + targetNickname + ". Imagine..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Scissors vs Scissors
	if challengerMove == "Scissors" && targetMove == "Scissors" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Scissors" + " 🤡" + "\n" + trimmed + ": Scissors" + " 🤡" + "\n \n"
		resultsMessage := "Welp, this is akward... " + author + " and " + targetNickname + " actually tied."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Scissors vs Rock
	if challengerMove == "Scissors" && targetMove == "Rock" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Scissors" + "\n" + trimmed + ": Rock" + " 👑" + "\n \n"
		resultsMessage := author + " got rekt by " + targetNickname + ". Imagine..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	// Scissors vs Paper
	if challengerMove == "Scissors" && targetMove == "Paper" {
		resultsTitle := author + " ⚔️ " + targetNickname + "\n \n"
		resultsOverview := authorAt + ": Scissors" + " 👑" + "\n" + trimmed + ": Paper" + "\n \n"
		resultsMessage := author + " kinda clapped " + trimmed + " ngl..."
		resultsFull := resultsTitle + resultsOverview + resultsMessage
		_, err := s.ChannelMessageSendReply(m.ChannelID, resultsFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}
}
