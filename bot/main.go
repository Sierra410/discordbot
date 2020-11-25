package bot

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	ownId    string
	logDebug = log.New(os.Stderr, "\033[34m[DEBUG]\033[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	logErr   = log.New(os.Stderr, "\033[31m[Error]\033[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	logInfo  = log.New(os.Stdout, "\033[36m[Info]\033[0m ", log.Ldate|log.Ltime)
)

func main() {
	var err error

	flag.StringVar(&configFilePath, "c", configFilePath, "config file")

	flag.Parse()

	logInfo.Println("Config:", configFilePath)
	logInfo.Println("Starting...")

	cfg = NewConfig(configFilePath)
	cfg.Load()
	cfg.cfg.SetAutosaveTime(time.Second * 15)

	token := interfaceToString(cfg.Get("", "Token"))
	if token == "" {
		logErr.Println("No token!")
		os.Exit(1)
	}

	//Session
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	//Handlers
	session.AddHandler(func(session *discordgo.Session, data *discordgo.Ready) {
		ownId = data.User.ID
		logInfo.Println("Userid: ", ownId)

		for _, x := range readyCallbacks {
			x(session, data)
		}
	})

	// session.AddHandler(rulesReactionAdd)
	// session.AddHandler(rulesReationRemove)

	session.AddHandler(messageCreate)

	//Open connection
	err = session.Open()
	if err != nil {
		panic(err)
	}

	//Wait for termination
	termchan := make(chan os.Signal, 1)
	signal.Notify(termchan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-termchan

	//Buy ice cream
	session.Close()
}
