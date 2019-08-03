package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

// Entwicklerinfos anzeigen
func showDevInfo() {

	if System.Dev == true {

		fmt.Print("\033[31m")
		fmt.Println("* * * * * D E V   M O D E * * * * *")
		fmt.Println("Version: ", System.Version)
		fmt.Println("Build:   ", System.Build)
		fmt.Println("* * * * * * * * * * * * * * * * * *")
		fmt.Print("\033[0m")
		fmt.Println()

	}

	return
}

// Alle Systemordner erstellen
func createSystemFolders() (err error) {

	e := reflect.ValueOf(&System.Folder).Elem()

	for i := 0; i < e.NumField(); i++ {

		var folder = e.Field(i).Interface().(string)

		err = checkFolder(folder)

		if err != nil {
			return
		}

	}

	return
}

// Alle Systemdateien erstellen
func createSystemFiles() (err error) {

	var debug string
	for _, file := range SystemFiles {

		var filename = getPlatformFile(System.Folder.Config + file)

		err = checkFile(filename)
		if err != nil {
			// Datei existiert nicht, wird jetzt erstellt
			err = saveMapToJSONFile(filename, make(map[string]interface{}))
			if err != nil {
				return
			}

			debug = fmt.Sprintf("Create File:%s", filename)
			showDebug(debug, 1)

		}

		switch file {

		case "authentication.json":
			System.File.Authentication = filename
		case "pms.json":
			System.File.PMS = filename
		case "settings.json":
			System.File.Settings = filename
		case "xepg.json":
			System.File.XEPG = filename
		case "urls.json":
			System.File.URLS = filename

		}

	}

	return
}

// Einstellungen laden und default Werte setzen (xTeVe)
func loadSettings() (settings SettingsStrcut, err error) {

	settingsMap, err := loadJSONFileToMap(System.File.Settings)
	if err != nil {
		return
	}

	// Deafult Werte setzten
	var defaults = make(map[string]interface{})
	var dataMap = make(map[string]interface{})

	dataMap["xmltv"] = make(map[string]interface{})
	dataMap["m3u"] = make(map[string]interface{})
	dataMap["hdhr"] = make(map[string]interface{})

	defaults["api"] = false
	defaults["authentication.api"] = false
	defaults["authentication.m3u"] = false
	defaults["authentication.pms"] = false
	defaults["authentication.web"] = false
	defaults["authentication.xml"] = false
	defaults["backup.keep"] = 10
	defaults["backup.path"] = System.Folder.Backup
	defaults["buffer"] = false
	defaults["buffer.size.kb"] = 1024
	defaults["buffer.timeout"] = 500
	defaults["cache.images"] = false
	defaults["epgSource"] = "XEPG"
	defaults["files"] = dataMap
	defaults["files.update"] = true
	defaults["filter"] = make(map[string]interface{})
	defaults["git.branch"] = System.Branch
	defaults["language"] = "en"
	defaults["log.entries.ram"] = 500
	defaults["mapping.first.channel"] = 1000
	defaults["xepg.replace.missing.images"] = true
	defaults["m3u8.adaptive.bandwidth.mbps"] = 10
	defaults["port"] = "34400"
	defaults["ssdp"] = true
	defaults["tuner"] = 1
	defaults["update"] = []string{"0000"}
	defaults["user.agent"] = System.Name
	defaults["uuid"] = createUUID()
	defaults["version"] = System.DBVersion
	defaults["xteveAutoUpdate"] = true
	defaults["temp.path"] = System.Folder.Temp

	// Default Werte setzen
	for key, value := range defaults {
		if _, ok := settingsMap[key]; !ok {
			settingsMap[key] = value
		}
	}

	err = json.Unmarshal([]byte(mapToJSON(settingsMap)), &settings)
	if err != nil {
		return
	}

	// Einstellungen von den Flags übernehmen
	if len(System.Flag.Port) > 0 {
		settings.Port = System.Flag.Port
	}

	if len(System.Flag.Branch) > 0 {
		settings.Branch = System.Flag.Branch
		showInfo(fmt.Sprintf("Git Branch:Switching Git Branch to -> %s", settings.Branch))
	}

	settings.Version = System.DBVersion

	err = saveSettings(settings)

	return
}

// Einstellungen speichern (xTeVe)
func saveSettings(settings SettingsStrcut) (err error) {

	if settings.BackupKeep == 0 {
		settings.BackupKeep = 10
	}

	if len(settings.BackupPath) == 0 {
		settings.BackupPath = System.Folder.Backup
	}

	if settings.BufferTimeout < 0 {
		settings.BufferTimeout = 0
	}

	System.Folder.Temp = settings.TempPath + settings.UUID + string(os.PathSeparator)

	err = writeByteToFile(System.File.Settings, []byte(mapToJSON(settings)))
	if err != nil {
		return
	}

	Settings = settings

	if System.Dev == true {
		Settings.UUID = "2019-01-DEV-xTeVe!"
	}

	setDeviceID()

	return
}

// Zugriff über die Domain ermöglichen
func setGlobalDomain(domain string) {

	System.Domain = domain

	switch Settings.AuthenticationPMS {
	case true:
		System.Addresses.DVR = "username:password@" + System.Domain
	case false:
		System.Addresses.DVR = System.Domain
	}

	switch Settings.AuthenticationM3U {
	case true:
		System.Addresses.M3U = System.ServerProtocol.M3U + "://" + System.Domain + "/m3u/xteve.m3u?username=xxx&password=yyy<br>(Specific groups: [http://...&group-title=foo,bar])"
	case false:
		System.Addresses.M3U = System.ServerProtocol.M3U + "://" + System.Domain + "/m3u/xteve.m3u     (Specific groups: [http://...?group-title=foo,bar])"
	}

	switch Settings.AuthenticationXML {
	case true:
		System.Addresses.XML = System.ServerProtocol.XML + "://" + System.Domain + "/xmltv/xteve.xml?username=xxx&password=yyy"
	case false:
		System.Addresses.XML = System.ServerProtocol.XML + "://" + System.Domain + "/xmltv/xteve.xml"
	}

	if Settings.EpgSource != "XEPG" {
		System.Addresses.M3U = getErrMsg(2106)
		System.Addresses.XML = getErrMsg(2106)
	}

	return
}

// UUID generieren
func createUUID() (uuid string) {
	uuid = time.Now().Format("2006-01") + "-" + randomString(4) + "-" + randomString(6)
	return
}

// Eindeutige Geräte ID für Plex generieren
func setDeviceID() {

	var id = Settings.UUID

	switch Settings.Tuner {
	case 1:
		System.DeviceID = id

	default:
		System.DeviceID = fmt.Sprintf("%s:%d", id, Settings.Tuner)
	}

	return
}

// Provider Streaming-URL zu xTeVe Streaming-URL konvertieren
func createStreamingURL(streamingType, playlistID, channelNumber, channelName, url string) (streamingURL string, err error) {

	var streamInfo StreamInfo
	var serverProtocol string

	if len(Data.Cache.StreamingURLS) == 0 {
		Data.Cache.StreamingURLS = make(map[string]StreamInfo)
	}

	var urlID = getMD5(fmt.Sprintf("%s-%s", playlistID, url))

	if s, ok := Data.Cache.StreamingURLS[urlID]; ok {

		streamInfo = s

	} else {

		streamInfo.URL = url
		streamInfo.Name = channelName
		streamInfo.PlaylistID = playlistID
		streamInfo.ChannelNumber = channelNumber
		streamInfo.URLid = urlID

		Data.Cache.StreamingURLS[urlID] = streamInfo

	}

	switch streamingType {

	case "DVR":
		serverProtocol = System.ServerProtocol.DVR

	case "M3U":
		serverProtocol = System.ServerProtocol.M3U

	}

	streamingURL = fmt.Sprintf("%s://%s/stream/%s", serverProtocol, System.Domain, streamInfo.URLid)

	return
}

func getStreamInfo(urlID string) (streamInfo StreamInfo, err error) {

	if len(Data.Cache.StreamingURLS) == 0 {

		tmp, err := loadJSONFileToMap(System.File.URLS)
		if err != nil {
			return streamInfo, err
		}

		err = json.Unmarshal([]byte(mapToJSON(tmp)), &Data.Cache.StreamingURLS)
		if err != nil {
			return streamInfo, err
		}

	}

	if s, ok := Data.Cache.StreamingURLS[urlID]; ok {
		streamInfo = s
		streamInfo.URL = strings.Trim(streamInfo.URL, "\r\n")
	} else {
		err = errors.New("Streaming error")
	}

	return
}
