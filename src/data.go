package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"xteve/src/internal/authentication"
	"xteve/src/internal/imgcache"
)

// Einstellungen ändern (WebUI)
func updateServerSettings(request RequestStruct) (settings SettingsStruct, err error) {

	var oldSettings = jsonToMap(mapToJSON(Settings))
	var newSettings = jsonToMap(mapToJSON(request.Settings))
	var reloadData = false
	var cacheImages = false
	var createXEPGFiles = false
	var debug string

	// -vvv [URL] --sout '#transcode{vcodec=mp4v, acodec=mpga} :standard{access=http, mux=ogg}'

	for key, value := range newSettings {

		if _, ok := oldSettings[key]; ok {

			switch key {

			case "tuner":
				showWarning(2105)

			case "epgSource":
				reloadData = true

			case "update":
				// Leerzeichen aus den Werten entfernen und Formatierung der Uhrzeit überprüfen (0000 - 2359)
				var newUpdateTimes = make([]string, 0)

				for _, v := range value.([]interface{}) {

					v = strings.Replace(v.(string), " ", "", -1)

					_, err := time.Parse("1504", v.(string))
					if err != nil {
						ShowError(err, 1012)
						return Settings, err
					}

					newUpdateTimes = append(newUpdateTimes, v.(string))

				}

				if len(newUpdateTimes) == 0 {
					//newUpdateTimes = append(newUpdateTimes, "0000")
				}

				value = newUpdateTimes

			case "cache.images":
				cacheImages = true

			case "xepg.replace.missing.images":
				createXEPGFiles = true

			case "backup.path":
				value = strings.TrimRight(value.(string), string(os.PathSeparator)) + string(os.PathSeparator)
				err = checkFolder(value.(string))
				if err == nil {

					err = checkFilePermission(value.(string))
					if err != nil {
						return
					}

				}

				if err != nil {
					return
				}

			case "temp.path":
				value = strings.TrimRight(value.(string), string(os.PathSeparator)) + string(os.PathSeparator)
				err = checkFolder(value.(string))
				if err == nil {

					err = checkFilePermission(value.(string))
					if err != nil {
						return
					}

				}

				if err != nil {
					return
				}

			case "ffmpeg.path", "vlc.path":
				var path = value.(string)
				if len(path) > 0 {

					err = checkFile(path)
					if err != nil {
						return
					}

				}

			case "scheme.m3u", "scheme.xml":
				createXEPGFiles = true

			}

			oldSettings[key] = value

			switch fmt.Sprintf("%T", value) {

			case "bool":
				debug = fmt.Sprintf("Save Setting:Key: %s | Value: %t (%T)", key, value, value)

			case "string":
				debug = fmt.Sprintf("Save Setting:Key: %s | Value: %s (%T)", key, value, value)

			case "[]interface {}":
				debug = fmt.Sprintf("Save Setting:Key: %s | Value: %v (%T)", key, value, value)

			case "float64":
				debug = fmt.Sprintf("Save Setting:Key: %s | Value: %d (%T)", key, int(value.(float64)), value)

			default:
				debug = fmt.Sprintf("%T", value)
			}

			showDebug(debug, 1)

		}

	}

	// Einstellungen aktualisieren
	err = json.Unmarshal([]byte(mapToJSON(oldSettings)), &Settings)
	if err != nil {
		return
	}

	if Settings.AuthenticationWEB == false {

		Settings.AuthenticationAPI = false
		Settings.AuthenticationM3U = false
		Settings.AuthenticationPMS = false
		Settings.AuthenticationWEB = false
		Settings.AuthenticationXML = false

	}

	// Buffer Einstellungen überprüfen
	if len(Settings.FFmpegOptions) == 0 {
		Settings.FFmpegOptions = System.FFmpeg.DefaultOptions
	}

	if len(Settings.VLCOptions) == 0 {
		Settings.VLCOptions = System.VLC.DefaultOptions
	}

	switch Settings.Buffer {

	case "ffmpeg":

		if len(Settings.FFmpegPath) == 0 {
			err = errors.New(getErrMsg(2020))
			return
		}

	case "vlc":

		if len(Settings.VLCPath) == 0 {
			err = errors.New(getErrMsg(2021))
			return
		}

	}

	err = saveSettings(Settings)
	if err == nil {

		settings = Settings

		if reloadData == true {

			err = buildDatabaseDVR()
			if err != nil {
				return
			}

			buildXEPG(false)

		}

		if cacheImages == true {

			if Settings.EpgSource == "XEPG" && System.ImageCachingInProgress == 0 {

				Data.Cache.Images, err = imgcache.New(System.Folder.ImagesCache, fmt.Sprintf("%s://%s/images/", System.ServerProtocol.WEB, System.Domain), Settings.CacheImages)
				if err != nil {
					ShowError(err, 0)
				}

				switch Settings.CacheImages {

				case false:
					createXMLTVFile()
					createM3UFile()

				case true:
					go func() {

						createXMLTVFile()
						createM3UFile()

						System.ImageCachingInProgress = 1
						showInfo("Image Caching:Images are cached")

						Data.Cache.Images.Image.Caching()
						showInfo("Image Caching:Done")

						System.ImageCachingInProgress = 0

						buildXEPG(false)

					}()

				}

			}

		}

		if createXEPGFiles == true {

			go func() {
				createXMLTVFile()
				createM3UFile()
			}()

		}

	}

	return
}

// Providerdaten speichern (WebUI)
func saveFiles(request RequestStruct, fileType string) (err error) {

	var filesMap = make(map[string]interface{})
	var newData = make(map[string]interface{})
	var indicator string
	var reloadData = false

	switch fileType {
	case "m3u":
		filesMap = Settings.Files.M3U
		newData = request.Files.M3U
		indicator = "M"

	case "hdhr":
		filesMap = Settings.Files.HDHR
		newData = request.Files.HDHR
		indicator = "H"

	case "xmltv":
		filesMap = Settings.Files.XMLTV
		newData = request.Files.XMLTV
		indicator = "X"
	}

	if len(filesMap) == 0 {
		filesMap = make(map[string]interface{})
	}

	for dataID, data := range newData {

		if dataID == "-" {

			// Neue Providerdatei
			dataID = indicator + randomString(19)
			data.(map[string]interface{})["new"] = true
			filesMap[dataID] = data

		} else {

			// Bereits vorhandene Providerdatei
			for key, value := range data.(map[string]interface{}) {

				var oldData = filesMap[dataID].(map[string]interface{})
				oldData[key] = value

			}

		}

		switch fileType {

		case "m3u":
			Settings.Files.M3U = filesMap

		case "hdhr":
			Settings.Files.HDHR = filesMap

		case "xmltv":
			Settings.Files.XMLTV = filesMap

		}

		// Neue Providerdatei
		if _, ok := data.(map[string]interface{})["new"]; ok {

			reloadData = true
			err = getProviderData(fileType, dataID)
			delete(data.(map[string]interface{}), "new")

			if err != nil {
				delete(filesMap, dataID)
				return
			}

		}

		if _, ok := data.(map[string]interface{})["delete"]; ok {

			deleteLocalProviderFiles(dataID, fileType)
			reloadData = true

		}

		err = saveSettings(Settings)
		if err != nil {
			return
		}

		if reloadData == true {

			err = buildDatabaseDVR()
			if err != nil {
				return err
			}

			buildXEPG(false)

		}

		Settings, _ = loadSettings()

	}

	return
}

// Providerdaten manuell aktualisieren (WebUI)
func updateFile(request RequestStruct, fileType string) (err error) {

	var updateData = make(map[string]interface{})

	switch fileType {

	case "m3u":
		updateData = request.Files.M3U

	case "hdhr":
		updateData = request.Files.HDHR

	case "xmltv":
		updateData = request.Files.XMLTV
	}

	for dataID := range updateData {

		err = getProviderData(fileType, dataID)
		if err == nil {
			err = buildDatabaseDVR()
			buildXEPG(false)
		}

	}

	return
}

// Providerdaten löschen (WebUI)
func deleteLocalProviderFiles(dataID, fileType string) {

	var removeData = make(map[string]interface{})
	var fileExtension string

	switch fileType {

	case "m3u":
		removeData = Settings.Files.M3U
		fileExtension = ".m3u"

	case "hdhr":
		removeData = Settings.Files.HDHR
		fileExtension = ".json"

	case "xmltv":
		removeData = Settings.Files.XMLTV
		fileExtension = ".xml"
	}

	if _, ok := removeData[dataID]; ok {
		delete(removeData, dataID)
		os.RemoveAll(System.Folder.Data + dataID + fileExtension)
	}

	return
}

// Filtereinstellungen speichern (WebUI)
func saveFilter(request RequestStruct) (settings SettingsStruct, err error) {

	var filterMap = make(map[int64]interface{})
	var newData = make(map[int64]interface{})
	var defaultFilter FilterStruct
	var newFilter = false

	defaultFilter.Active = true
	defaultFilter.CaseSensitive = false

	filterMap = Settings.Filter
	newData = request.Filter

	var createNewID = func() (id int64) {

	newID:
		if _, ok := filterMap[id]; ok {
			id++
			goto newID
		}

		return id
	}

	for dataID, data := range newData {

		if dataID == -1 {

			// Neuer Filter
			newFilter = true
			dataID = createNewID()
			filterMap[dataID] = jsonToMap(mapToJSON(defaultFilter))

		}

		// Filter aktualisieren / löschen
		for key, value := range data.(map[string]interface{}) {

			// Filter löschen
			if _, ok := data.(map[string]interface{})["delete"]; ok {
				delete(filterMap, dataID)
				break
			}

			if filter, ok := data.(map[string]interface{})["filter"].(string); ok {

				if len(filter) == 0 {

					err = errors.New(getErrMsg(1014))
					if newFilter == true {
						delete(filterMap, dataID)
					}

					return
				}

			}

			if oldData, ok := filterMap[dataID].(map[string]interface{}); ok {
				oldData[key] = value
			}

		}

	}

	err = saveSettings(Settings)
	if err != nil {
		return
	}

	settings = Settings

	err = buildDatabaseDVR()
	if err != nil {
		return
	}

	buildXEPG(false)

	return
}

// XEPG Mapping speichern
func saveXEpgMapping(request RequestStruct) (err error) {

	var tmp = Data.XEPG

	Data.Cache.Images, err = imgcache.New(System.Folder.ImagesCache, fmt.Sprintf("%s://%s/images/", System.ServerProtocol.WEB, System.Domain), Settings.CacheImages)
	if err != nil {
		ShowError(err, 0)
	}

	err = json.Unmarshal([]byte(mapToJSON(request.EpgMapping)), &tmp)
	if err != nil {
		return
	}

	err = saveMapToJSONFile(System.File.XEPG, request.EpgMapping)
	if err != nil {
		return err
	}

	Data.XEPG.Channels = request.EpgMapping

	if System.ScanInProgress == 0 {

		System.ScanInProgress = 1
		cleanupXEPG()
		System.ScanInProgress = 0
		buildXEPG(true)

	} else {

		// Wenn während des erstellen der Datanbank das Mapping erneut gespeichert wird, wird die Datenbank erst später erneut aktualisiert.
		go func() {

			if System.BackgroundProcess == true {
				return
			}

			System.BackgroundProcess = true

			for {
				time.Sleep(time.Duration(1) * time.Second)
				if System.ScanInProgress == 0 {
					break
				}

			}

			System.ScanInProgress = 1
			cleanupXEPG()
			System.ScanInProgress = 0
			buildXEPG(false)
			showInfo("XEPG:" + fmt.Sprintf("Ready to use"))

			System.BackgroundProcess = false

		}()

	}

	return
}

// Benutzerdaten speichern (WebUI)
func saveUserData(request RequestStruct) (err error) {

	var userData = request.UserData

	var newCredentials = func(userID string, newUserData map[string]interface{}) (err error) {

		var newUsername, newPassword string
		if username, ok := newUserData["username"].(string); ok {
			newUsername = username
		}

		if password, ok := newUserData["password"].(string); ok {
			newPassword = password
		}

		if len(newUsername) > 0 {
			err = authentication.ChangeCredentials(userID, newUsername, newPassword)
		}

		return
	}

	for userID, newUserData := range userData {

		err = newCredentials(userID, newUserData.(map[string]interface{}))
		if err != nil {
			return
		}

		if request.DeleteUser == true {
			err = authentication.RemoveUser(userID)
			return
		}

		delete(newUserData.(map[string]interface{}), "password")
		delete(newUserData.(map[string]interface{}), "confirm")

		if _, ok := newUserData.(map[string]interface{})["delete"]; ok {

			authentication.RemoveUser(userID)

		} else {

			err = authentication.WriteUserData(userID, newUserData.(map[string]interface{}))
			if err != nil {
				return
			}

		}

	}

	return
}

// Neuen Benutzer anlegen (WebUI)
func saveNewUser(request RequestStruct) (err error) {

	var data = request.UserData
	var username = data["username"].(string)
	var password = data["password"].(string)

	delete(data, "password")
	delete(data, "confirm")

	userID, err := authentication.CreateNewUser(username, password)
	if err != nil {
		return
	}

	err = authentication.WriteUserData(userID, data)
	return
}

// Wizard (WebUI)
func saveWizard(request RequestStruct) (nextStep int, err error) {

	var wizard = jsonToMap(mapToJSON(request.Wizard))

	for key, value := range wizard {

		switch key {

		case "tuner":
			Settings.Tuner = int(value.(float64))
			nextStep = 1

		case "epgSource":
			Settings.EpgSource = value.(string)
			nextStep = 2

		case "m3u", "xmltv":

			var filesMap = make(map[string]interface{})
			var data = make(map[string]interface{})
			var indicator, dataID string

			filesMap = make(map[string]interface{})

			data["type"] = key
			data["new"] = true

			switch key {

			case "m3u":
				filesMap = Settings.Files.M3U
				data["name"] = "M3U"
				indicator = "M"

			case "xmltv":
				filesMap = Settings.Files.XMLTV
				data["name"] = "XMLTV"
				indicator = "X"

			}

			dataID = indicator + randomString(19)
			data["file.source"] = value.(string)

			filesMap[dataID] = data

			switch key {
			case "m3u":
				Settings.Files.M3U = filesMap
				nextStep = 3

				err = getProviderData(key, dataID)

				if err != nil {
					ShowError(err, 000)
					delete(filesMap, dataID)
					return
				}

				err = buildDatabaseDVR()
				if err != nil {
					ShowError(err, 000)
					delete(filesMap, dataID)
					return
				}

				if Settings.EpgSource == "PMS" {
					nextStep = 10
				}

			case "xmltv":
				Settings.Files.XMLTV = filesMap
				nextStep = 10

				err = getProviderData(key, dataID)

				if err != nil {

					ShowError(err, 000)
					delete(filesMap, dataID)
					return

				}

				buildXEPG(false)
				System.ScanInProgress = 0

			}

		}

	}

	err = saveSettings(Settings)
	if err != nil {
		return
	}

	return
}

// Filterregeln erstellen
func createFilterRules() (err error) {

	Data.Filter = nil
	var dataFilter Filter

	for _, f := range Settings.Filter {

		var filter FilterStruct

		var exclude, include string

		err = json.Unmarshal([]byte(mapToJSON(f)), &filter)
		if err != nil {
			return
		}

		switch filter.Type {

		case "custom-filter":
			dataFilter.CaseSensitive = filter.CaseSensitive
			dataFilter.Rule = filter.Filter
			dataFilter.Type = filter.Type

			Data.Filter = append(Data.Filter, dataFilter)

		case "group-title":
			if len(filter.Include) > 0 {
				include = fmt.Sprintf(" {%s}", filter.Include)
			}

			if len(filter.Exclude) > 0 {
				exclude = fmt.Sprintf(" !{%s}", filter.Exclude)
			}

			dataFilter.CaseSensitive = filter.CaseSensitive
			dataFilter.Rule = fmt.Sprintf("%s%s%s", filter.Filter, include, exclude)
			dataFilter.Type = filter.Type

			Data.Filter = append(Data.Filter, dataFilter)
		}

	}

	return
}

// Datenbank für das DVR System erstellen
func buildDatabaseDVR() (err error) {

	System.ScanInProgress = 1

	Data.Streams.All = make([]interface{}, 0, System.UnfilteredChannelLimit)
	Data.Streams.Active = make([]interface{}, 0, System.UnfilteredChannelLimit)
	Data.Streams.Inactive = make([]interface{}, 0, System.UnfilteredChannelLimit)
	Data.Playlist.M3U.Groups.Text = []string{}
	Data.Playlist.M3U.Groups.Value = []string{}
	Data.StreamPreviewUI.Active = []string{}
	Data.StreamPreviewUI.Inactive = []string{}

	var availableFileTypes = []string{"m3u", "hdhr"}

	var tmpGroupsM3U = make(map[string]int64)

	err = createFilterRules()
	if err != nil {
		return
	}

	for _, fileType := range availableFileTypes {

		var playlistFile = getLocalProviderFiles(fileType)

		for n, i := range playlistFile {

			var channels []interface{}
			var groupTitle, tvgID, uuid int = 0, 0, 0
			var keys = []string{"group-title", "tvg-id", "uuid"}
			var compatibility = make(map[string]int)

			var id = strings.TrimSuffix(getFilenameFromPath(i), path.Ext(getFilenameFromPath(i)))
			var playlistName = getProviderParameter(id, fileType, "name")

			switch fileType {

			case "m3u":
				channels, err = parsePlaylist(i, fileType)
			case "hdhr":
				channels, err = parsePlaylist(i, fileType)

			}

			if err != nil {
				ShowError(err, 1005)
				err = errors.New(playlistName + ": Local copy of the file no longer exists")
				ShowError(err, 0)
				playlistFile = append(playlistFile[:n], playlistFile[n+1:]...)
			}

			// Streams analysieren
			for _, stream := range channels {

				var s = stream.(map[string]string)
				s["_file.m3u.path"] = i
				s["_file.m3u.name"] = playlistName
				s["_file.m3u.id"] = id

				// Kompatibilität berechnen
				for _, key := range keys {

					switch key {
					case "uuid":
						if value, ok := s["_uuid.key"]; ok {
							if len(value) > 0 {
								uuid++
							}
						}

					case "group-title":
						if value, ok := s[key]; ok {
							if len(value) > 0 {

								if _, ok := tmpGroupsM3U[value]; ok {
									tmpGroupsM3U[value]++
								} else {
									tmpGroupsM3U[value] = 1
								}

								groupTitle++
							}
						}

					case "tvg-id":
						if value, ok := s[key]; ok {
							if len(value) > 0 {
								tvgID++
							}
						}

					}

				}

				Data.Streams.All = append(Data.Streams.All, stream)

				// Neuer Filter ab Version 1.3.0
				var preview string
				var status = filterThisStream(stream)

				if name, ok := s["name"]; ok {
					var group string

					if v, ok := s["group-title"]; ok {
						group = v
					}

					preview = fmt.Sprintf("%s [%s]", name, group)

				}

				switch status {

				case true:
					Data.StreamPreviewUI.Active = append(Data.StreamPreviewUI.Active, preview)
					Data.Streams.Active = append(Data.Streams.Active, stream)

				case false:
					Data.StreamPreviewUI.Inactive = append(Data.StreamPreviewUI.Inactive, preview)
					Data.Streams.Inactive = append(Data.Streams.Inactive, stream)

				}

			}

			if tvgID == 0 {
				compatibility["tvg.id"] = 0
			} else {
				compatibility["tvg.id"] = int(tvgID * 100 / len(channels))
			}

			if groupTitle == 0 {
				compatibility["group.title"] = 0
			} else {
				compatibility["group.title"] = int(groupTitle * 100 / len(channels))
			}

			if uuid == 0 {
				compatibility["stream.id"] = 0
			} else {
				compatibility["stream.id"] = int(uuid * 100 / len(channels))
			}

			compatibility["streams"] = len(channels)

			setProviderCompatibility(id, fileType, compatibility)

		}

	}

	for group, count := range tmpGroupsM3U {
		var text = fmt.Sprintf("%s (%d)", group, count)
		var value = fmt.Sprintf("%s", group)
		Data.Playlist.M3U.Groups.Text = append(Data.Playlist.M3U.Groups.Text, text)
		Data.Playlist.M3U.Groups.Value = append(Data.Playlist.M3U.Groups.Value, value)
	}

	sort.Strings(Data.Playlist.M3U.Groups.Text)
	sort.Strings(Data.Playlist.M3U.Groups.Value)

	if len(Data.Streams.Active) == 0 && len(Data.Streams.All) <= System.UnfilteredChannelLimit && len(Settings.Filter) == 0 {
		Data.Streams.Active = Data.Streams.All
		Data.Streams.Inactive = make([]interface{}, 0)

		Data.StreamPreviewUI.Active = Data.StreamPreviewUI.Inactive
		Data.StreamPreviewUI.Inactive = []string{}

	}

	if len(Data.Streams.Active) > System.PlexChannelLimit {
		showWarning(2000)
	}

	if len(Settings.Filter) == 0 && len(Data.Streams.All) > System.UnfilteredChannelLimit {
		showWarning(2001)
	}

	System.ScanInProgress = 0
	showInfo(fmt.Sprintf("All streams:%d", len(Data.Streams.All)))
	showInfo(fmt.Sprintf("Active streams:%d", len(Data.Streams.Active)))
	showInfo(fmt.Sprintf("Filter:%d", len(Data.Filter)))

	sort.Strings(Data.StreamPreviewUI.Active)
	sort.Strings(Data.StreamPreviewUI.Inactive)

	return
}

// Speicherort aller lokalen Providerdateien laden, immer für eine Dateityp (M3U, XMLTV usw.)
func getLocalProviderFiles(fileType string) (localFiles []string) {

	var fileExtension string
	var dataMap = make(map[string]interface{})

	switch fileType {

	case "m3u":
		fileExtension = ".m3u"
		dataMap = Settings.Files.M3U

	case "hdhr":
		fileExtension = ".json"
		dataMap = Settings.Files.HDHR

	case "xmltv":
		fileExtension = ".xml"
		dataMap = Settings.Files.XMLTV

	}

	for dataID := range dataMap {
		localFiles = append(localFiles, System.Folder.Data+dataID+fileExtension)
	}

	return
}

// Providerparameter anhand von dem Key ausgeben
func getProviderParameter(id, fileType, key string) (s string) {

	var dataMap = make(map[string]interface{})

	switch fileType {
	case "m3u":
		dataMap = Settings.Files.M3U

	case "hdhr":
		dataMap = Settings.Files.HDHR

	case "xmltv":
		dataMap = Settings.Files.XMLTV
	}

	if data, ok := dataMap[id].(map[string]interface{}); ok {

		if v, ok := data[key].(string); ok {
			s = v
		}

		if v, ok := data[key].(float64); ok {
			s = strconv.Itoa(int(v))
		}

	}

	return
}

// Provider Statistiken Kompatibilität aktualisieren
func setProviderCompatibility(id, fileType string, compatibility map[string]int) {

	var dataMap = make(map[string]interface{})

	switch fileType {
	case "m3u":
		dataMap = Settings.Files.M3U

	case "hdhr":
		dataMap = Settings.Files.HDHR

	case "xmltv":
		dataMap = Settings.Files.XMLTV
	}

	if data, ok := dataMap[id].(map[string]interface{}); ok {

		data["compatibility"] = compatibility

		switch fileType {
		case "m3u":
			Settings.Files.M3U = dataMap
		case "hdhr":
			Settings.Files.HDHR = dataMap
		case "xmltv":
			Settings.Files.XMLTV = dataMap
		}

		saveSettings(Settings)

	}

}
