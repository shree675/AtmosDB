package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/mattn/go-shellwords"

	t "atmosdb/types"
	"atmosdb/util"
)

func Run(conn string) {
	session := createSession(conn)
	fmt.Println("Initialized an interactive terminal session: " + session.SId)
	fmt.Println()

	rl, err := readline.New(util.GetMagentaStr("atmos-cli> "))
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		os.Exit(0)
	}()

	for {
		line, err := rl.Readline()
		if err == io.EOF {
			fmt.Println("[EOF] detected")
			break
		}
		if err != nil {
			util.PrintError(err.Error())
			return
		}

		line = strings.TrimSpace(line)
		if line == "exit" {
			break
		} else if line == "db.version" {
			printVersion(session)
			continue
		}

		args, _ := shellwords.Parse(line)
		process(args, session)
	}

	fmt.Println("Graceful exit")
}

func createSession(url string) t.SessionConfig {
	return t.SessionConfig{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		SId:  uuid.NewString(),
		Conn: url}
}

func printVersion(session t.SessionConfig) {
	res, err := session.Client.Get(session.Conn + "/version")
	if err != nil {
		util.PrintRed("Error while fetching value: " + err.Error())
		return
	}

	defer res.Body.Close()
	defer io.Copy(io.Discard, res.Body)

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		util.PrintError(err.Error())
		return
	}

	version := strings.TrimSpace(string(respBody))
	fmt.Println(strings.Trim(version, "\""))
}

func process(args []string, session t.SessionConfig) {
	switch args[0] {
	case string(util.GET):
		HandleGetKey(session, args)
	case string(util.SETINT):
		HandleSetValue(session, args, util.INT, false)
	case string(util.SETFLOAT):
		HandleSetValue(session, args, util.FLOAT, false)
	case string(util.SETSTR):
		HandleSetValue(session, args, util.STRING, false)
	case string(util.SETINT_TTL):
		HandleSetValue(session, args, util.INT, true)
	case string(util.SETFLOAT_TTL):
		HandleSetValue(session, args, util.FLOAT, true)
	case string(util.SETSTR_TTL):
		HandleSetValue(session, args, util.STRING, true)
	case string(util.DEL):
		HandleDeleteValue(session, args)
	case string(util.INCR):
		HandleIncrementValue(session, args, 1)
	case string(util.DECR):
		HandleIncrementValue(session, args, -1)
	case string(util.EXISTS):
		HandleGetKey(session, args)
	default:
		util.PrintGray("Unknown command '" + args[0] + "'")
	}
}
