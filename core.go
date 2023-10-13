package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
)

type Core struct {
	server         		*fiber.App
	mainFolderPath 		string
	port 		   		int
	disableStartMessage bool
	logs 				bool
}

func newCore(folder string, port int) *Core {
	s, _ := strings.CutSuffix(folder, "\\")

	return &Core{
		server: fiber.New(fiber.Config{
			DisableStartupMessage: true},
		),
		mainFolderPath: s,
		port: port,
		disableStartMessage: false,
	}
}

func (c *Core) reload() {
	time := time.Now()
	fmt.Println(time.Format("02.01.2006 15:04:05") + " [ Started reloading files ] 🌐")

	err := c.server.Shutdown()
	if err != nil {
		log.Fatal(err)
	}
	c.server = fiber.New(fiber.Config{
		DisableStartupMessage: true},
	)
	c.searchFolder(c.mainFolderPath, "")
	
	go c.server.Listen(":" + fmt.Sprint(c.port))
}

func (c *Core) startServer() {
	if (!c.disableStartMessage) {
		color.HiMagenta("\nListening on adres -> 127.0.0.1:" + fmt.Sprint(c.port))
		color.Red(serverOnlineText)
		color.HiMagenta(guitarArt)
		color.HiMagenta("- Press `r` to reload")
	}
	go c.server.Listen(":" + fmt.Sprint(c.port))
}

func (c *Core) searchFolder(path string, endpoint string) {
	div, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err.Error() + "\nError: Could not read given folder")
	}

	_, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range div {
		fileName := file.Name()
		nameLen := len(fileName)

		if file.IsDir() {
			route := fileName
			if fileName[0] == '[' &&
				fileName[nameLen - 1] == ']' &&
				nameLen > 2 {

				runes := []rune(fileName[:nameLen - 1])
				runes[0] = ':'
				route = string(runes)
			}
			c.searchFolder(
				path + "\\" + fileName,
				endpoint + "/" + route,
			)

		} else if fileName[0] == '+' && nameLen > 1 {
			var name = fileName
			var addition = ""

			if name != "+page.html" {
				addition, _ = strings.CutPrefix(name, "+")
				addition = "/" + addition
			}
			var fullEndpoint = endpoint + addition
			if (fullEndpoint == "") {
				fullEndpoint = "/"
			}
			var logBegin = "- " + fullEndpoint + " - "
			var filePath = path + "/" + name
			
			c.server.Get(
				fullEndpoint,
				func(ctx *fiber.Ctx) error {
					_, err := os.Stat(filePath)
					if err != nil {
						if (c.logs) {
							fmt.Println(logBegin + "[ 404 ] NOT_FOUND")
						}
						return ctx.SendStatus(404)
					}
					if (c.logs) {
						fmt.Println(logBegin + "[ 200 ] OK")
					}
					return ctx.SendFile(filePath, false)
				},
			)
		}
	}
}
