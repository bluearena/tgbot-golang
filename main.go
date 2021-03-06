package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/floodcode/tbf"
)

const (
	configPath = "config.json"
)

var (
	bot       tbf.TelegramBotFramework
	botConfig BotConfig
)

// BotConfig contains bot's environment variables
type BotConfig struct {
	Token string `json:"token"`
	Delay int    `json:"delay"`
}

func main() {
	configData, err := ioutil.ReadFile(configPath)
	checkError(err)

	err = json.Unmarshal(configData, &botConfig)
	checkError(err)

	bot, err = tbf.New(botConfig.Token)
	checkError(err)

	addRoutes()

	err = bot.Poll(tbf.PollConfig{
		Delay: botConfig.Delay,
	})

	checkError(err)
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func addRoutes() {
	bot.AddRoute("help", helpAction)
	bot.AddRoute("start", helpAction)
	bot.AddRoute("ping", pingAction)
	bot.AddRoute("compile", compileAction)
	bot.AddRoute("main", mainAction)
	bot.AddRoute("fmt", fmtAction)
}

func helpAction(req tbf.Request) {
	req.QuickMessageMD(fmt.Sprintf(strings.Join([]string{
		"Available commads:",
		"/help - Get this message",
		"/compile - Compile code",
		"/main - Compile code in main function",
		"/fmt - Format code",
	}, "\n")))
}

func pingAction(req tbf.Request) {
	req.QuickMessage("Pong!")
}

func compileAction(req tbf.Request) {
	code := getCode(req)
	req.SendTyping()

	output, err := runCode(code)
	if err != nil {
		req.QuickMessage("Error: " + err.Error())
		return
	}

	output = strings.TrimSpace(output)
	if len(output) == 0 {
		req.QuickMessage("No output")
	} else {
		req.QuickMessageMD(fmt.Sprintf("```\n%s\n```", output))
	}
}

func mainAction(req tbf.Request) {
	codeTemplate := `
	package main

	func main() {
		%s
	}`

	req.Args = fmt.Sprintf(codeTemplate, getCode(req))
	compileAction(req)
}

func fmtAction(req tbf.Request) {
	formatedCode := strings.TrimSpace(formatCode(getCode(req)))
	if len(formatedCode) == 0 {
		req.QuickMessage("No output")
	} else {
		req.QuickMessageMD(fmt.Sprintf("```\n%s\n```", formatedCode))
	}
}

func getCode(req tbf.Request) string {
	if len(req.Args) != 0 {
		return req.Args
	}

	req.QuickMessage("Now send me the code")

	for {
		newReq := req.WaitNext()
		if newReq.Message.Text != "" {
			return newReq.Message.Text
		}

		req.QuickMessage("You should send code as a text message")
	}
}
