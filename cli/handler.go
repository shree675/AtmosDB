package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
