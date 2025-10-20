package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	t "atmosdb/types"
	"atmosdb/util"
)

func HandleGetKey(session t.SessionConfig, args []string) {
	if len(args) != 2 {
		util.PrintGray("Incorrect arguments for " + args[0] + ", expecting 'key'")
		return
	}

	payload := t.InputPayload{
		SId: session.SId,
		Key: args[1],
	}

	body, _ := json.Marshal(&payload)
	reqBody := bytes.NewBuffer(body)
	res, err := session.Client.Post(session.Conn+"/get", "application/json", reqBody)
	if err != nil {
		util.PrintRed("Error while fetching value: " + err.Error())
		return
	}

	defer res.Body.Close()
	defer io.Copy(io.Discard, res.Body)

	if res.StatusCode == http.StatusNotFound {
		if args[0] == string(util.EXISTS) {
			fmt.Println(false)
		} else {
			fmt.Println(nil)
		}
		return
	}
	if res.StatusCode != http.StatusOK {
		util.PrintRed("Error while fetching value: Received status code " + strconv.Itoa(res.StatusCode))
		return
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		util.PrintError(err.Error())
		return
	}

	if args[0] == string(util.EXISTS) {
		fmt.Println(true)
		return
	}

	var op t.OutputPayload
	if err := json.Unmarshal(respBody, &op); err != nil {
		util.PrintRed("Incompatible response obtained from server: " + err.Error())
		return
	}

	util.PrintBlue("(" + util.Types[op.Type].GetType().String() + ")")
	fmt.Println(op.Val)
}

func HandleSetValue(session t.SessionConfig, args []string, code util.BaseType, ttl bool) {
	if (!ttl && len(args) != 3) || (ttl && len(args) != 4) {
		msg := "Incorrect arguments for " + args[0] + ", expecting 'key' 'value'"
		if ttl {
			msg += " 'ttl'"
		}
		util.PrintGray(msg)
		return
	}

	payload := t.InputPayload{
		SId:  session.SId,
		Key:  args[1],
		Val:  args[2],
		Type: int8(code),
		Op:   int8(util.PUT),
	}
	if ttl {
		v, err := strconv.Atoi(args[3])
		if err != nil {
			util.PrintGray("Incorrect type provided for TTL")
			return
		}
		payload.Ttl = uint32(v)
	}

	body, _ := json.Marshal(&payload)
	respBody := bytes.NewBuffer(body)
	res, err := session.Client.Post(session.Conn+"/set", "application/json", respBody)
	if err != nil {
		util.PrintRed("Error while storing value: " + err.Error())
		return
	}

	defer res.Body.Close()
	defer io.Copy(io.Discard, res.Body)

	if res.StatusCode == http.StatusOK {
		util.PrintGreen("[OK]")
	} else {
		util.PrintYellow("[FAILED]")
	}
}

func HandleDeleteValue(session t.SessionConfig, args []string) {
	if len(args) != 2 {
		util.PrintGray("Incorrect arguments for " + args[0] + ", expecting 'key'")
		return
	}

	payload := t.InputPayload{
		SId: session.SId,
		Key: args[1],
		Op:  int8(util.DELETE),
	}

	body, _ := json.Marshal(&payload)
	respBody := bytes.NewBuffer(body)
	res, err := session.Client.Post(session.Conn+"/set", "application/json", respBody)
	if err != nil {
		util.PrintRed("Error while deleting key: " + err.Error())
		return
	}

	defer res.Body.Close()
	defer io.Copy(io.Discard, res.Body)

	if res.StatusCode == http.StatusOK {
		util.PrintGreen("[OK]")
	} else {
		util.PrintYellow("[FAILED]")
	}
}

func HandleIncrementValue(session t.SessionConfig, args []string, order int) {
	if len(args) != 2 {
		util.PrintGray("Incorrect arguments for " + args[0] + ", expecting 'key'")
		return
	}

	payload := t.InputPayload{
		SId: session.SId,
		Key: args[1],
		Val: strconv.Itoa(order),
		Op:  int8(util.DELTA),
	}

	body, _ := json.Marshal(&payload)
	respBody := bytes.NewBuffer(body)
	res, err := session.Client.Post(session.Conn+"/set", "application/json", respBody)
	if err != nil {
		util.PrintRed("Error while deleting key: " + err.Error())
		return
	}

	defer res.Body.Close()
	defer io.Copy(io.Discard, res.Body)

	if res.StatusCode == http.StatusOK {
		util.PrintGreen("[OK]")
	} else {
		util.PrintYellow("[FAILED]")
	}
}

func HandleSubscribeKey(session t.SessionConfig, args []string) {
	if len(args) != 2 {
		util.PrintGray("Incorrect arguments for " + args[0] + ", expecting 'key'")
		return
	}

	payload := t.InputSubscriptionPayload{
		SId: session.SId,
		Key: args[1],
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	body, _ := json.Marshal(&payload)
	req, err := http.NewRequestWithContext(ctx, "POST", session.Conn+"/subscribe", bytes.NewBuffer(body))
	if err != nil {
		util.PrintRed("Error while subscribing to key: " + err.Error())
		return
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")

	res, err := session.SSEClient.Do(req)
	if err != nil {
		util.PrintRed("Error while subscribing to key: " + err.Error())
		return
	}
	defer res.Body.Close()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func(cancel context.CancelFunc) {
		<-ch
		cancel()
	}(cancel)

	for {
		buf := make([]byte, util.StreamBufSize)

		n, err := res.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				util.PrintRed("[FAILED] Stream closed by server, the key might not exist")
			} else if errors.Is(err, context.Canceled) {
				return
			} else {
				util.PrintRed("[FAILED] " + err.Error())
			}
			return
		}

		value := string(buf[:n])
		if value == util.StreamDeleteId {
			util.PrintYellow("[TERM] Key has been removed")
			return
		}

		fmt.Printf("%s\n", value)
	}
}
