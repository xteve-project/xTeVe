package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"../src/internal/authentication"

	"github.com/gorilla/websocket"
)

// StartWebserver : Startet den Webserver
func StartWebserver() (err error) {

	var port = Settings.Port

	http.HandleFunc("/", Index)
	http.HandleFunc("/stream/", Stream)
	http.HandleFunc("/xmltv/", xTeVe)
	http.HandleFunc("/m3u/", xTeVe)
	http.HandleFunc("/data/", WS)
	http.HandleFunc("/web/", Web)
	http.HandleFunc("/download/", Download)
	http.HandleFunc("/api/", API)
	http.HandleFunc("/images/", Images)
	http.HandleFunc("/data_images/", DataImages)
	//http.HandleFunc("/auto/", Auto)

	showInfo("DVR IP:" + System.IPAddress + ":" + Settings.Port)

	var ips = len(System.IPAddressesV4) + len(System.IPAddressesV6) - 1
	switch ips {

	case 0:
		showHighlight(fmt.Sprintf("Web Interface:%s://%s:%s/web/", System.ServerProtocol.WEB, System.IPAddress, Settings.Port))

	case 1:
		showHighlight(fmt.Sprintf("Web Interface:%s://%s:%s/web/ | xTeVe is also available via the other %d IP", System.ServerProtocol.WEB, System.IPAddress, Settings.Port, ips))

	default:
		showHighlight(fmt.Sprintf("Web Interface:%s://%s:%s/web/ | xTeVe is also available via the other %d IP's", System.ServerProtocol.WEB, System.IPAddress, Settings.Port, len(System.IPAddressesV4)+len(System.IPAddressesV6)-1))

	}

	if err = http.ListenAndServe(":"+port, nil); err != nil {
		ShowError(err, 1001)
		return
	}

	return
}

// Index : Web Server /
func Index(w http.ResponseWriter, r *http.Request) {

	var err error
	var response []byte
	var path = r.URL.Path
	var debug string

	setGlobalDomain(r.Host)

	debug = fmt.Sprintf("Web Server Request:Path: %s", path)
	showDebug(debug, 2)

	switch path {

	case "/discover.json":
		response, err = getDiscover()
		w.Header().Set("Content-Type", "application/json")

	case "/lineup_status.json":
		response, err = getLineupStatus()
		w.Header().Set("Content-Type", "application/json")

	case "/lineup.json":
		if Settings.AuthenticationPMS == true {

			_, err := basicAuth(r, "authentication.pms")
			if err != nil {
				ShowError(err, 000)
				httpStatusError(w, r, 403)
				return
			}

		}

		response, err = getLineup()
		w.Header().Set("Content-Type", "application/json")

	case "/device.xml", "/capability":
		response, err = getCapability()
		w.Header().Set("Content-Type", "application/xml")

	default:
		response, err = getCapability()
		w.Header().Set("Content-Type", "application/xml")
	}

	if err == nil {

		w.WriteHeader(200)
		w.Write(response)
		return

	}

	httpStatusError(w, r, 500)

	return
}

// Stream : Web Server /stream/
func Stream(w http.ResponseWriter, r *http.Request) {

	var path = strings.Replace(r.RequestURI, "/stream/", "", 1)
	//var stream = strings.SplitN(path, "-", 2)

	streamInfo, err := getStreamInfo(path)
	if err != nil {
		ShowError(err, 1203)
		httpStatusError(w, r, 404)
		return
	}

	if strings.Index(streamInfo.URL, "rtsp://") != -1 || strings.Index(streamInfo.URL, "rtp://") != -1 {
		err = errors.New("RTSP and RTP streams are not supported")
		ShowError(err, 2004)

		showInfo("Streaming URL:" + streamInfo.URL)
		http.Redirect(w, r, streamInfo.URL, 302)

		showInfo("Streaming Info:URL was passed to the client")
		return
	}

	showInfo(fmt.Sprintf("Buffer:%t", Settings.Buffer))

	if Settings.Buffer == true {
		showInfo(fmt.Sprintf("Buffer Size:%d KB", Settings.BufferSize))
	}

	showInfo(fmt.Sprintf("Channel Name:%s", streamInfo.Name))
	showInfo(fmt.Sprintf("Client User-Agent:%s", r.Header.Get("User-Agent")))

	// Prüfen ob der Buffer verwendet werden soll
	switch Settings.Buffer {

	case true:
		bufferingStream(streamInfo.PlaylistID, streamInfo.URL, streamInfo.Name, w, r)

	case false:
		showInfo("Streaming URL:" + streamInfo.URL)
		http.Redirect(w, r, streamInfo.URL, 302)

		showInfo("Streaming Info:URL was passed to the client.")
		showInfo("Streaming Info:xTeVe is no longer involved, the client connects directly to the streaming server.")

	}

	return
}

// Auto : HDHR routing (wird derzeit nicht benutzt)
func Auto(w http.ResponseWriter, r *http.Request) {

	var channelID = strings.Replace(r.RequestURI, "/auto/v", "", 1)
	fmt.Println(channelID)

	/*
		switch Settings.Buffer {

		case true:
			var playlistID, streamURL, err = getStreamByChannelID(channelID)
			if err == nil {
				bufferingStream(playlistID, streamURL, w, r)
			} else {
				httpStatusError(w, r, 404)
			}

		case false:
			httpStatusError(w, r, 423)
		}
	*/

	return
}

// xTeVe : Web Server /xmltv/ und /m3u/
func xTeVe(w http.ResponseWriter, r *http.Request) {

	var requestType, groupTitle, file, content string
	var err error
	var path = strings.TrimPrefix(r.URL.Path, "/")
	var groups = []string{}

	setGlobalDomain(r.Host)

	// XMLTV Datei
	if strings.Contains(path, "xmltv/") {

		requestType = "xml"

		file = System.Folder.Data + getFilenameFromPath(path)

		content, err = readStringFromFile(file)
		if err != nil {
			httpStatusError(w, r, 404)
			return
		}

	}

	// M3U Datei
	if strings.Contains(path, "m3u/") {

		requestType = "m3u"
		groupTitle = r.URL.Query().Get("group-title")

		if System.Dev == false {
			// false: Dateiname wird im Header gesetzt
			// true: M3U wird direkt im Browser angezeigt
			w.Header().Set("Content-Disposition", "attachment; filename="+getFilenameFromPath(path))
		}

		if len(groupTitle) > 0 {
			groups = strings.Split(groupTitle, ",")
		}

		content, err = buildM3U(groups)
		if err != nil {
			ShowError(err, 000)
		}

	}

	// Authentifizierung überprüfen
	err = urlAuth(r, requestType)
	if err != nil {
		ShowError(err, 000)
		httpStatusError(w, r, 403)
		return
	}

	if err == nil {
		w.Write([]byte(content))
	}

	return
}

// Images : Image Cache /images/
func Images(w http.ResponseWriter, r *http.Request) {

	var path = strings.TrimPrefix(r.URL.Path, "/")
	var filePath = System.Folder.ImagesCache + getFilenameFromPath(path)

	content, err := readByteFromFile(filePath)
	if err != nil {
		httpStatusError(w, r, 404)
		return
	}

	w.Header().Add("Content-Type", getContentType(filePath))
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(content)))
	w.WriteHeader(200)
	w.Write(content)

	return
}

// DataImages : Image Pfad für Logos / Bilder die hochgeladen wurden /data_images/
func DataImages(w http.ResponseWriter, r *http.Request) {

	var path = strings.TrimPrefix(r.URL.Path, "/")
	var filePath = System.Folder.ImagesUpload + getFilenameFromPath(path)

	content, err := readByteFromFile(filePath)
	if err != nil {
		httpStatusError(w, r, 404)
		return
	}

	w.Header().Add("Content-Type", getContentType(filePath))
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(content)))
	w.WriteHeader(200)
	w.Write(content)

	return
}

// WS : Web Sockets /ws/
func WS(w http.ResponseWriter, r *http.Request) {

	var request RequestStruct
	var response ResponseStruct
	response.Status = true

	var newToken string

	if r.Header.Get("Origin") != "http://"+r.Host {
		httpStatusError(w, r, 403)
		return
	}

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		ShowError(err, 0)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	setGlobalDomain(r.Host)

	for {

		err = conn.ReadJSON(&request)

		if err != nil {
			return
		}

		if System.ConfigurationWizard == false {

			switch Settings.AuthenticationWEB {

			// Token Authentication
			case true:

				var token string
				tokens, ok := r.URL.Query()["Token"]

				if !ok || len(tokens[0]) < 1 {
					token = "-"
				} else {
					token = tokens[0]
				}

				newToken, err = tokenAuthentication(token)
				if err != nil {

					response.Status = false
					response.Reload = true
					response.Error = err.Error()
					request.Cmd = "-"

					if err = conn.WriteJSON(response); err != nil {
						ShowError(err, 1102)
					}

					return
				}

				response.Token = newToken
				response.Users, _ = authentication.GetAllUserData()

			}

		}

		switch request.Cmd {
		// Daten lesen
		case "getServerConfig":
			//response.Config = Settings

		case "updateLog":
			response = setDefaultResponseData(response, false)
			if err = conn.WriteJSON(response); err != nil {
				ShowError(err, 1022)
			} else {
				return
				break
			}
			return

		case "loadFiles":
			//response.Response = Settings.Files

		// Daten schreiben
		case "saveSettings":
			var authenticationUpdate = Settings.AuthenticationWEB
			response.Settings, err = updateServerSettings(request)
			if err == nil {

				response.OpenMenu = strconv.Itoa(indexOfString("settings", System.WEB.Menu))

				if Settings.AuthenticationWEB == true && authenticationUpdate == false {
					response.Reload = true
				}

			}

		case "saveFilesM3U":
			err = saveFiles(request, "m3u")
			if err == nil {
				response.OpenMenu = strconv.Itoa(indexOfString("playlist", System.WEB.Menu))
			}

		case "updateFileM3U":
			err = updateFile(request, "m3u")
			if err == nil {
				response.OpenMenu = strconv.Itoa(indexOfString("playlist", System.WEB.Menu))
			}

		case "saveFilesHDHR":
			err = saveFiles(request, "hdhr")
			if err == nil {
				response.OpenMenu = strconv.Itoa(indexOfString("playlist", System.WEB.Menu))
			}

		case "updateFileHDHR":
			err = updateFile(request, "hdhr")
			if err == nil {
				response.OpenMenu = strconv.Itoa(indexOfString("playlist", System.WEB.Menu))
			}

		case "saveFilesXMLTV":
			err = saveFiles(request, "xmltv")
			if err == nil {
				response.OpenMenu = strconv.Itoa(indexOfString("xmltv", System.WEB.Menu))
			}

		case "updateFileXMLTV":
			err = updateFile(request, "xmltv")
			if err == nil {
				response.OpenMenu = strconv.Itoa(indexOfString("xmltv", System.WEB.Menu))
			}

		case "saveFilter":
			response.Settings, err = saveFilter(request)
			if err == nil {
				response.OpenMenu = strconv.Itoa(indexOfString("filter", System.WEB.Menu))
			}

		case "saveEpgMapping":
			err = saveXEpgMapping(request)

		case "saveUserData":
			err = saveUserData(request)
			if err == nil {
				response.OpenMenu = strconv.Itoa(indexOfString("users", System.WEB.Menu))
			}

		case "saveNewUser":
			err = saveNewUser(request)
			if err == nil {
				response.OpenMenu = strconv.Itoa(indexOfString("users", System.WEB.Menu))
			}

		case "resetLogs":
			WebScreenLog.Log = make([]string, 0)
			WebScreenLog.Errors = 0
			WebScreenLog.Warnings = 0
			response.OpenMenu = strconv.Itoa(indexOfString("log", System.WEB.Menu))

		case "xteveBackup":
			file, errNew := xteveBackup()
			err = errNew
			if err == nil {
				response.OpenLink = fmt.Sprintf("%s://%s/download/%s", System.ServerProtocol.WEB, System.Domain, file)
			}

		case "xteveRestore":
			WebScreenLog.Log = make([]string, 0)
			WebScreenLog.Errors = 0
			WebScreenLog.Warnings = 0

			if len(request.Base64) > 0 {

				newWebURL, err := xteveRestoreFromWeb(request.Base64)
				if err != nil {
					ShowError(err, 000)
					response.Alert = err.Error()
				}

				if err == nil {

					if len(newWebURL) > 0 {
						response.Alert = "Backup was successfully restored.\nThe port of the sTeVe URL has changed, you have to restart xTeVe.\nAfter a restart, xTeVe can be reached again at the following URL:\n" + newWebURL
					} else {
						response.Alert = "Backup was successfully restored."
						response.Reload = true
					}
					showInfo("xTeVe:" + "Backup successfully restored.")
				}

			}

		case "uploadLogo":
			if len(request.Base64) > 0 {
				response.LogoURL, err = uploadLogo(request.Base64, request.Filename)

				if err == nil {

					if err = conn.WriteJSON(response); err != nil {
						ShowError(err, 1022)
					} else {
						return
					}

				}

			}

		case "saveWizard":
			nextStep, errNew := saveWizard(request)

			err = errNew
			if err == nil {

				if nextStep == 10 {
					System.ConfigurationWizard = false
					response.Reload = true
				} else {
					response.Wizard = nextStep
				}

			}

			/*
				case "wizardCompleted":
					System.ConfigurationWizard = false
					response.Reload = true
			*/
		default:
			fmt.Println("+ + + + + + + + + + +", request.Cmd)

			var requestMap = make(map[string]interface{}) // Debug
			_ = requestMap
			if System.Dev == true {
				fmt.Println(mapToJSON(requestMap))
			}

		}

		if err != nil {
			response.Status = false
			response.Error = err.Error()
			response.Settings = Settings
		}

		response = setDefaultResponseData(response, true)
		if System.ConfigurationWizard == true {
			response.ConfigurationWizard = System.ConfigurationWizard
		}

		if err = conn.WriteJSON(response); err != nil {
			ShowError(err, 1022)
		} else {
			break
		}

	}

	return
}

// Web : Web Server /web/
func Web(w http.ResponseWriter, r *http.Request) {

	var lang = make(map[string]interface{})
	var err error

	var requestFile = strings.Replace(r.URL.Path, "/web", "html", -1)
	var content, contentType, file string

	var language LanguageUI

	if System.Dev == true {

		lang, err = loadJSONFileToMap(fmt.Sprintf("html/lang/%s.json", Settings.Language))
		if err != nil {
			ShowError(err, 000)
		}

	} else {

		var languageFile = "html/lang/en.json"

		if value, ok := webUI[languageFile].(string); ok {
			content = GetHTMLString(value)
			lang = jsonToMap(content)
		}

	}

	err = json.Unmarshal([]byte(mapToJSON(lang)), &language)
	if err != nil {
		ShowError(err, 000)
		return
	}

	if getFilenameFromPath(requestFile) == "html" {

		if len(Data.Streams.All) == 0 && System.ScanInProgress == 0 {
			System.ConfigurationWizard = true
		}

		switch System.ConfigurationWizard {

		case true:
			file = requestFile + "configuration.html"
			Settings.AuthenticationWEB = false

		case false:
			file = requestFile + "index.html"

		}

		if System.ScanInProgress == 1 {
			file = requestFile + "maintenance.html"
		}

		switch Settings.AuthenticationWEB {
		case true:

			var username, password, confirm string
			switch r.Method {
			case "POST":
				var allUsers, _ = authentication.GetAllUserData()

				username = r.FormValue("username")
				password = r.FormValue("password")

				if len(allUsers) == 0 {
					confirm = r.FormValue("confirm")
				}

				// Erster Benutzer wird angelegt (Passwortbestätigung ist vorhanden)
				if len(confirm) > 0 {

					var token, err = createFirstUserForAuthentication(username, password)
					if err != nil {
						httpStatusError(w, r, 429)
						return
					}
					// Redirect, damit die Daten aus dem Browser gelöscht werden.
					w = authentication.SetCookieToken(w, token)
					http.Redirect(w, r, "/web", 301)
					return

				}

				// Benutzername und Passwort vorhanden, wird jetzt überprüft
				if len(username) > 0 && len(password) > 0 {

					var token, err = authentication.UserAuthentication(username, password)
					if err != nil {
						file = requestFile + "login.html"
						lang["authenticationErr"] = language.Login.Failed
						break
					}

					w = authentication.SetCookieToken(w, token)
					http.Redirect(w, r, "/web", 301) // Redirect, damit die Daten aus dem Browser gelöscht werden.

				} else {
					w = authentication.SetCookieToken(w, "-")
					http.Redirect(w, r, "/web", 301) // Redirect, damit die Daten aus dem Browser gelöscht werden.
				}

				return

			case "GET":
				lang["authenticationErr"] = ""
				_, token, err := authentication.CheckTheValidityOfTheTokenFromHTTPHeader(w, r)

				if err != nil {
					file = requestFile + "login.html"
					break
				}

				err = checkAuthorizationLevel(token, "authentication.web")
				if err != nil {
					file = requestFile + "login.html"
					break
				}

			}

			allUserData, err := authentication.GetAllUserData()
			if err != nil {
				ShowError(err, 000)
				httpStatusError(w, r, 403)
				return
			}

			if len(allUserData) == 0 && Settings.AuthenticationWEB == true {
				file = requestFile + "create-first-user.html"
			}

		}

		requestFile = file

		if value, ok := webUI[requestFile]; ok {

			content = GetHTMLString(value.(string))

			if contentType == "text/plain" {
				w.Header().Set("Content-Disposition", "attachment; filename="+getFilenameFromPath(requestFile))
			}

		} else {

			httpStatusError(w, r, 404)
			return
		}

	}

	if value, ok := webUI[requestFile].(string); ok {

		content = GetHTMLString(value)
		contentType = getContentType(requestFile)

		if contentType == "text/plain" {
			w.Header().Set("Content-Disposition", "attachment; filename="+getFilenameFromPath(requestFile))
		}

	} else {
		httpStatusError(w, r, 404)
		return
	}

	contentType = getContentType(requestFile)

	if System.Dev == true {
		// Lokale Webserver Dateien werden geladen, nur für die Entwicklung
		content, _ = readStringFromFile(requestFile)
	}

	w.Header().Add("Content-Type", contentType)
	w.WriteHeader(200)

	if contentType == "text/html" || contentType == "application/javascript" {
		content = parseTemplate(content, lang)
	}

	w.Write([]byte(content))
}

// API : API request /api/
func API(w http.ResponseWriter, r *http.Request) {

	/*
			API Bedingungen (ohne Authentifizierung):
			- API muss in den Einstellungen aktiviert sein

			Beispiel API Request mit curl
			Status:
			curl -X POST -H "Content-Type: application/json" -d '{"cmd":"status"}' http://localhost:34400/api/

			- - - - -

			API Bedingungen (mit Authentifizierung):
			- API muss in den Einstellungen aktiviert sein
			- API muss bei den Authentifizierungseinstellungen aktiviert sein
			- Benutzer muss die Berechtigung API haben

			Nach jeder API Anfrage wird ein Token generiert, dieser ist einmal in 60 Minuten gültig.
			In jeder Antwort ist ein neuer Token enthalten

			Beispiel API Request mit curl
			Login:
			curl -X POST -H "Content-Type: application/json" -d '{"cmd":"login","username":"plex","password":"123"}' http://localhost:34400/api/

			Antwort:
			{
		  	"status": true,
		  	"token": "U0T-NTSaigh-RlbkqERsHvUpgvaaY2dyRGuwIIvv"
			}

			Status mit Verwendung eines Tokens:
			curl -X POST -H "Content-Type: application/json" -d '{"cmd":"status","token":"U0T-NTSaigh-RlbkqERsHvUpgvaaY2dyRGuwIIvv"}' http://localhost:4400/api/

			Antwort:
			{
			  "epg.source": "XEPG",
			  "status": true,
			  "streams.active": 7,
			  "streams.all": 63,
			  "streams.xepg": 2,
			  "token": "mXiG1NE1MrTXDtyh7PxRHK5z8iPI_LzxsQmY-LFn",
			  "url.dvr": "localhost:34400",
			  "url.m3u": "http://localhost:34400/m3u/xteve.m3u",
			  "url.xepg": "http://localhost:34400/xmltv/xteve.xml",
			  "version.api": "1.1.0",
			  "version.xteve": "1.3.0"
			}
	*/

	setGlobalDomain(r.Host)
	var request APIRequestStruct
	var response APIResponseStruct

	var responseAPIError = func(err error) {

		var response APIResponseStruct

		response.Status = false
		response.Error = err.Error()
		w.Write([]byte(mapToJSON(response)))
		return

	}

	response.Status = true

	if Settings.API == false {
		httpStatusError(w, r, 423)
		return
	}

	if r.Method == "GET" {
		httpStatusError(w, r, 404)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		httpStatusError(w, r, 400)
		return

	}

	err = json.Unmarshal(b, &request)
	if err != nil {
		httpStatusError(w, r, 400)
		return
	}

	w.Header().Set("content-type", "application/json")

	if Settings.AuthenticationAPI == true {
		var token string
		switch len(request.Token) {
		case 0:
			if request.Cmd == "login" {
				token, err = authentication.UserAuthentication(request.Username, request.Password)
				if err != nil {
					responseAPIError(err)
					return
				}

			} else {
				err = errors.New("Login incorrect")
				if err != nil {
					responseAPIError(err)
					return
				}

			}

		default:
			token, err = tokenAuthentication(request.Token)
			fmt.Println(err)
			if err != nil {
				responseAPIError(err)
				return
			}

		}
		err = checkAuthorizationLevel(token, "authentication.api")
		if err != nil {
			responseAPIError(err)
			return
		}

		response.Token = token

	}

	switch request.Cmd {
	case "login": // Muss nichts übergeben werden

	case "status":

		response.VersionXteve = System.Version
		response.VersionAPI = System.APIVersion
		response.StreamsActive = int64(len(Data.Streams.Active))
		response.StreamsAll = int64(len(Data.Streams.All))
		response.StreamsXepg = int64(Data.XEPG.XEPGCount)
		response.EpgSource = Settings.EpgSource
		response.URLDvr = System.Domain
		response.URLM3U = System.ServerProtocol.M3U + "://" + System.Domain + "/m3u/xteve.m3u"
		response.URLXepg = System.ServerProtocol.XML + "://" + System.Domain + "/xmltv/xteve.xml"

	case "update.m3u":
		err = getProviderData("m3u", "")
		if err != nil {
			break
		}

		err = buildDatabaseDVR()
		if err != nil {
			break
		}

	case "update.hdhr":

		err = getProviderData("hdhr", "")
		if err != nil {
			break
		}

		err = buildDatabaseDVR()
		if err != nil {
			break
		}

	case "update.xmltv":
		err = getProviderData("xmltv", "")
		if err != nil {
			break
		}

	case "update.xepg":
		buildXEPG(false)

	default:
		err = errors.New(getErrMsg(5000))

	}

	if err != nil {
		responseAPIError(err)
	}

	w.Write([]byte(mapToJSON(response)))

	return
}

// Download : Datei Download
func Download(w http.ResponseWriter, r *http.Request) {

	var path = r.URL.Path
	var file = System.Folder.Temp + getFilenameFromPath(path)
	w.Header().Set("Content-Disposition", "attachment; filename="+getFilenameFromPath(file))

	content, err := readStringFromFile(file)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	os.RemoveAll(System.Folder.Temp + getFilenameFromPath(path))
	w.Write([]byte(content))
	return
}

func setDefaultResponseData(response ResponseStruct, data bool) (defaults ResponseStruct) {

	defaults = response

	// Folgende Daten immer an den Client übergeben
	defaults.ClientInfo.ARCH = System.ARCH
	defaults.ClientInfo.EpgSource = Settings.EpgSource
	defaults.ClientInfo.DVR = System.Addresses.DVR
	defaults.ClientInfo.M3U = System.Addresses.M3U
	defaults.ClientInfo.XML = System.Addresses.XML
	defaults.ClientInfo.OS = System.OS
	defaults.ClientInfo.Streams = fmt.Sprintf("%d / %d", len(Data.Streams.Active), len(Data.Streams.All))
	defaults.ClientInfo.UUID = Settings.UUID
	defaults.ClientInfo.Errors = WebScreenLog.Errors
	defaults.ClientInfo.Warnings = WebScreenLog.Warnings
	defaults.Notification = System.Notification
	defaults.Log = WebScreenLog

	switch System.Branch {

	case "master":
		defaults.ClientInfo.Version = fmt.Sprintf("%s", System.Version)

	default:
		defaults.ClientInfo.Version = fmt.Sprintf("%s (%s)", System.Version, System.Build)
		defaults.ClientInfo.Branch = System.Branch

	}

	if data == true {

		defaults.Users, _ = authentication.GetAllUserData()
		//defaults.DVR = System.DVRAddress

		if Settings.EpgSource == "XEPG" {

			defaults.ClientInfo.XEPGCount = Data.XEPG.XEPGCount

			var XEPG = make(map[string]interface{})

			if len(Data.Streams.Active) > 0 {

				XEPG["epgMapping"] = Data.XEPG.Channels
				XEPG["xmltvMap"] = Data.XMLTV.Mapping

			} else {

				XEPG["epgMapping"] = make(map[string]interface{})
				XEPG["xmltvMap"] = make(map[string]interface{})

			}

			defaults.XEPG = XEPG

		}

		defaults.Settings = Settings

		defaults.Data.Playlist.M3U.Groups.Text = Data.Playlist.M3U.Groups.Text
		defaults.Data.Playlist.M3U.Groups.Value = Data.Playlist.M3U.Groups.Value
		defaults.Data.StreamPreviewUI.Active = Data.StreamPreviewUI.Active
		defaults.Data.StreamPreviewUI.Inactive = Data.StreamPreviewUI.Inactive

	}

	return
}

func httpStatusError(w http.ResponseWriter, r *http.Request, httpStatusCode int) {
	http.Error(w, fmt.Sprintf("%s [%d]", http.StatusText(httpStatusCode), httpStatusCode), httpStatusCode)
	return
}

func getContentType(filename string) (contentType string) {

	if strings.HasSuffix(filename, ".html") {
		contentType = "text/html"
	} else if strings.HasSuffix(filename, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(filename, ".js") {
		contentType = "application/javascript"
	} else if strings.HasSuffix(filename, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filename, ".jpg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(filename, ".gif") {
		contentType = "image/gif"
	} else if strings.HasSuffix(filename, ".svg") {
		contentType = "image/svg+xml"
	} else if strings.HasSuffix(filename, ".mp4") {
		contentType = "video/mp4"
	} else if strings.HasSuffix(filename, ".webm") {
		contentType = "video/webm"
	} else if strings.HasSuffix(filename, ".ogg") {
		contentType = "video/ogg"
	} else if strings.HasSuffix(filename, ".mp3") {
		contentType = "audio/mp3"
	} else if strings.HasSuffix(filename, ".wav") {
		contentType = "audio/wav"
	} else {
		contentType = "text/plain"
	}

	return
}
