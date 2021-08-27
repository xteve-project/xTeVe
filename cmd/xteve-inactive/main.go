package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	xteve "xteve/src"
)

var port = flag.String("port", "", ": Server port          [34400] (default: 34400)")
var host = flag.String("host", "", ": Server host                  (default: localhost)")

func main() {
	flag.Parse()

	portNum := 34400
	if port != nil {
		var err error
		portNum, err = strconv.Atoi(*port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable parse port: %v\n", err)
			os.Exit(-1)
		}
	}

	hostname := "localhost"
	if host != nil {
		hostname = *host
	}

	requestBody, err := json.Marshal(&xteve.APIRequestStruct{
		Cmd: "status",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshall request: %v\n", err)
		os.Exit(-1)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s:%d/api/", hostname, portNum), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get API: %v\n", err)
		os.Exit(-1)
	}

	defer resp.Body.Close()

	respStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable read response: %v\n", err)
		os.Exit(-1)
	}

	var apiresp xteve.APIResponseStruct
	err = json.Unmarshal(respStr, &apiresp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable parse response: %v\n", err)
		fmt.Fprintf(os.Stderr, "%s\n", respStr)
		os.Exit(-1)
	}

	os.Exit(int(apiresp.TunerActive))
}
