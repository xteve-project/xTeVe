package up2date

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

// ClientInfo : Information about the key (NAME OS, ARCH, UUID, KEY)
type ClientInfo struct {
	Arch   string `json:"arch,required"`
	Branch string `json:"branch,required"`
	CMD    string `json:"cmd,omitempty"`
	Name   string `json:"name,required"`
	OS     string `json:"os,required"`
	URL    string `json:"url,required"`

	Response ServerResponse `json:"response,omitempty"`
}

//ServerResponse : Response from server after client request
type ServerResponse struct {
	Status    bool   `json:"status,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Version   string `json:"version,omitempty"`
	UpdateBIN string `json:"update.url.bin,omitempty"`
	UpdateZIP string `json:"update.url.zip,omitempty"`
	Filename  string `json:"filename.bin,omitempty"`
}

// Updater : Client infos
var Updater ClientInfo

// UpdateURL : URL for the new binary
var UpdateURL string

// Init : Init
func Init() {
	Updater.OS = runtime.GOOS
	Updater.Arch = runtime.GOARCH
}

// GetVersion : Information about the latest version
func GetVersion() (err error) {

	Updater.CMD = "getVersion"
	err = serverRequest()
	return
}

func serverRequest() (err error) {

	var serverResponse ServerResponse
	jsonByte, err := json.MarshalIndent(Updater, "", "  ")
	if err == nil {

		// Serververbindung prüfen
		u, err := url.Parse(Updater.URL)
		if err != nil {
			return err
		}
		var server = u.Host

		timeout := time.Duration(1 * time.Second)
		_, err = net.DialTimeout("tcp", server, timeout)
		if err != nil {
			return err
		}

		// Check redirect 301 <---> 308
		redirect, err := http.NewRequest("POST", Updater.URL, nil)

		client := &http.Client{}
		client.CheckRedirect = func(redirect *http.Request, via []*http.Request) error {
			return errors.New("Redirect")
		}

		resp, err := client.Do(redirect)

		if err != nil {
			// Redirect
			if resp.StatusCode >= 301 && resp.StatusCode <= 308 { //status code 301 <---> 308
				Updater.URL = resp.Header.Get("Location")
			} else {
				return err
			}
		}
		// ---

		req, err := http.NewRequest("POST", Updater.URL, bytes.NewBuffer(jsonByte))
		req.Header.Set("Content-Type", "application/json")

		client = &http.Client{}
		resp, err = client.Do(req)

		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			//fmt.Println(resp.StatusCode, Updater.URL, Updater.CMD)
			err = fmt.Errorf(fmt.Sprintf("%d: %s (%s)", resp.StatusCode, http.StatusText(resp.StatusCode), Updater.URL))
			return err
		}

		Updater.CMD = ""
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &serverResponse)

		if err != nil {
			return err
		}

		Updater.Response = serverResponse

	}

	return
}
