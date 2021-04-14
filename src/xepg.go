package src

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"sort"

	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"xteve/src/internal/imgcache"
)

// Provider XMLTV Datei überprüfen
func checkXMLCompatibility(id string, body []byte) (err error) {

	var xmltv XMLTV
	var compatibility = make(map[string]int)

	err = xml.Unmarshal(body, &xmltv)
	if err != nil {
		return
	}

	compatibility["xmltv.channels"] = len(xmltv.Channel)
	compatibility["xmltv.programs"] = len(xmltv.Program)

	setProviderCompatibility(id, "xmltv", compatibility)

	return
}

// XEPG Daten erstellen
func buildXEPG(background bool) {

	if System.ScanInProgress == 1 {
		return
	}

	System.ScanInProgress = 1

	var err error

	Data.Cache.Images, err = imgcache.New(System.Folder.ImagesCache, fmt.Sprintf("%s://%s/images/", System.ServerProtocol.WEB, System.Domain), Settings.CacheImages)
	if err != nil {
		ShowError(err, 0)
	}

	if Settings.EpgSource == "XEPG" {

		switch background {

		case true:

			go func() {

				createXEPGMapping()
				createXEPGDatabase()
				mapping()
				cleanupXEPG()
				createXMLTVFile()
				createM3UFile()

				showInfo("XEPG:" + fmt.Sprintf("Ready to use"))

				if Settings.CacheImages == true && System.ImageCachingInProgress == 0 {

					go func() {

						System.ImageCachingInProgress = 1
						showInfo(fmt.Sprintf("Image Caching:Images are cached (%d)", len(Data.Cache.Images.Queue)))

						Data.Cache.Images.Image.Caching()
						Data.Cache.Images.Image.Remove()
						showInfo("Image Caching:Done")

						createXMLTVFile()
						createM3UFile()

						System.ImageCachingInProgress = 0

					}()

				}

				System.ScanInProgress = 0

				// Cache löschen
				/*
					Data.Cache.XMLTV = make(map[string]XMLTV)
					Data.Cache.XMLTV = nil
				*/
				runtime.GC()

			}()

		case false:

			createXEPGMapping()
			createXEPGDatabase()
			mapping()
			cleanupXEPG()

			go func() {

				createXMLTVFile()
				createM3UFile()

				if Settings.CacheImages == true && System.ImageCachingInProgress == 0 {

					go func() {

						System.ImageCachingInProgress = 1
						showInfo(fmt.Sprintf("Image Caching:Images are cached (%d)", len(Data.Cache.Images.Queue)))

						Data.Cache.Images.Image.Caching()
						Data.Cache.Images.Image.Remove()
						showInfo("Image Caching:Done")

						createXMLTVFile()
						createM3UFile()

						System.ImageCachingInProgress = 0

					}()

				}

				showInfo("XEPG:" + fmt.Sprintf("Ready to use"))

				System.ScanInProgress = 0

				// Cache löschen
				//Data.Cache.XMLTV = make(map[string]XMLTV)
				//Data.Cache.XMLTV = nil
				runtime.GC()

			}()

		}

	} else {

		getLineup()
		System.ScanInProgress = 0

	}

}

// XEPG Daten aktualisieren
func updateXEPG(background bool) {

	if System.ScanInProgress == 1 {
		return
	}

	System.ScanInProgress = 1

	if Settings.EpgSource == "XEPG" {

		switch background {

		case false:

			createXEPGDatabase()
			mapping()
			cleanupXEPG()

			go func() {

				createXMLTVFile()
				createM3UFile()
				showInfo("XEPG:" + fmt.Sprintf("Ready to use"))

				System.ScanInProgress = 0

			}()

		case true:
			System.ScanInProgress = 0

		}

	} else {

		System.ScanInProgress = 0

	}

	// Cache löschen
	//Data.Cache.XMLTV = nil //make(map[string]XMLTV)
	//Data.Cache.XMLTV = make(map[string]XMLTV)

	return
}

// Mapping Menü für die XMLTV Dateien erstellen
func createXEPGMapping() {

	Data.XMLTV.Files = getLocalProviderFiles("xmltv")
	Data.XMLTV.Mapping = make(map[string]interface{})

	var tmpMap = make(map[string]interface{})

	var friendlyDisplayName = func(channel Channel) (displayName string) {
		var dn = channel.DisplayName
		displayName = dn[0].Value

		switch len(dn) {
		case 1:
			displayName = dn[0].Value
		default:
			displayName = fmt.Sprintf("%s (%s)", dn[1].Value, dn[0].Value)
		}

		return
	}

	if len(Data.XMLTV.Files) > 0 {

		for i := len(Data.XMLTV.Files) - 1; i >= 0; i-- {

			var file = Data.XMLTV.Files[i]

			var err error
			var fileID = strings.TrimSuffix(getFilenameFromPath(file), path.Ext(getFilenameFromPath(file)))
			showInfo("XEPG:" + "Parse XMLTV file: " + getProviderParameter(fileID, "xmltv", "name"))

			//xmltv, err = getLocalXMLTV(file)
			var xmltv XMLTV

			err = getLocalXMLTV(file, &xmltv)
			if err != nil {
				Data.XMLTV.Files = append(Data.XMLTV.Files, Data.XMLTV.Files[i+1:]...)
				var errMsg = err.Error()
				err = errors.New(getProviderParameter(fileID, "xmltv", "name") + ": " + errMsg)
				ShowError(err, 000)
			}

			// XML Parsen (Provider Datei)
			if err == nil {

				// Daten aus der XML Datei in eine temporäre Map schreiben
				var xmltvMap = make(map[string]interface{})

				for _, c := range xmltv.Channel {
					var channel = make(map[string]interface{})

					channel["id"] = c.ID
					channel["display-name"] = friendlyDisplayName(*c)
					channel["icon"] = c.Icon.Src

					xmltvMap[c.ID] = channel

				}

				tmpMap[getFilenameFromPath(file)] = xmltvMap
				Data.XMLTV.Mapping[getFilenameFromPath(file)] = xmltvMap

			}

		}

		Data.XMLTV.Mapping = tmpMap
		tmpMap = make(map[string]interface{})

	} else {

		if System.ConfigurationWizard == false {
			showWarning(1007)
		}

	}

	// Auswahl für den Dummy erstellen
	var dummy = make(map[string]interface{})
	var times = []string{"30", "60", "90", "120", "180", "240", "360"}

	for _, i := range times {

		var dummyChannel = make(map[string]string)
		dummyChannel["display-name"] = i + " Minutes"
		dummyChannel["id"] = i + "_Minutes"
		dummyChannel["icon"] = ""

		dummy[dummyChannel["id"]] = dummyChannel

	}

	Data.XMLTV.Mapping["xTeVe Dummy"] = dummy

	return
}

// XEPG Datenbank erstellen / aktualisieren
func createXEPGDatabase() (err error) {

	var allChannelNumbers = make([]float64, 0, System.UnfilteredChannelLimit)
	Data.Cache.Streams.Active = make([]string, 0, System.UnfilteredChannelLimit)
	Data.XEPG.Channels = make(map[string]interface{}, System.UnfilteredChannelLimit)

	Data.XEPG.Channels, err = loadJSONFileToMap(System.File.XEPG)
	if err != nil {
		ShowError(err, 1004)
		return err
	}

	var createNewID = func() (xepg string) {

		var firstID = 0 //len(Data.XEPG.Channels)

	newXEPGID:

		if _, ok := Data.XEPG.Channels["x-ID."+strconv.FormatInt(int64(firstID), 10)]; ok {
			firstID++
			goto newXEPGID
		}

		xepg = "x-ID." + strconv.FormatInt(int64(firstID), 10)
		return
	}

	var getFreeChannelNumber = func() (xChannelID string) {

		sort.Float64s(allChannelNumbers)

		var firstFreeNumber float64 = Settings.MappingFirstChannel

		for {

			if indexOfFloat64(firstFreeNumber, allChannelNumbers) == -1 {
				xChannelID = fmt.Sprintf("%g", firstFreeNumber)
				allChannelNumbers = append(allChannelNumbers, firstFreeNumber)
				return
			}

			firstFreeNumber++

		}

		return
	}

	var generateHashForChannel = func(m3uID string, groupTitle string, tvgID string, tvgName string, uuidKey string, uuidValue string) string {
		hash := md5.Sum([]byte(m3uID + groupTitle + tvgID + tvgName + uuidKey + uuidValue))
		return hex.EncodeToString(hash[:])
	}

	showInfo("XEPG:" + "Update database")

	// Kanal mit fehlenden Kanalnummern löschen.  Delete channel with missing channel numbers
	for id, dxc := range Data.XEPG.Channels {

		var xepgChannel XEPGChannelStruct
		err = json.Unmarshal([]byte(mapToJSON(dxc)), &xepgChannel)
		if err != nil {
			return
		}

		if len(xepgChannel.XChannelID) == 0 {
			delete(Data.XEPG.Channels, id)
		}

		if xChannelID, err := strconv.ParseFloat(xepgChannel.XChannelID, 64); err == nil {
			allChannelNumbers = append(allChannelNumbers, xChannelID)
		}

	}

	// Make a map of the db channels based on their previously downloaded attributes -- filename, group, title, etc
	var xepgChannelsValuesMap = make(map[string]XEPGChannelStruct, System.UnfilteredChannelLimit)
	for _, v := range Data.XEPG.Channels {
		var channel XEPGChannelStruct
		err = json.Unmarshal([]byte(mapToJSON(v)), &channel)
		if err != nil {
			return
		}
		channelHash := generateHashForChannel(channel.FileM3UID, channel.GroupTitle, channel.TvgID, channel.TvgName, channel.UUIDKey, channel.UUIDValue)
		xepgChannelsValuesMap[channelHash] = channel
	}

	for _, dsa := range Data.Streams.Active {

		var channelExists = false  // Entscheidet ob ein Kanal neu zu Datenbank hinzugefügt werden soll.  Decides whether a channel should be added to the database
		var channelHasUUID = false // Überprüft, ob der Kanal (Stream) eindeutige ID's besitzt.  Checks whether the channel (stream) has unique IDs
		var currentXEPGID string   // Aktuelle Datenbank ID (XEPG). Wird verwendet, um den Kanal in der Datenbank mit dem Stream der M3u zu aktualisieren. Current database ID (XEPG) Used to update the channel in the database with the stream of the M3u

		var m3uChannel M3UChannelStructXEPG

		err = json.Unmarshal([]byte(mapToJSON(dsa)), &m3uChannel)
		if err != nil {
			return
		}

		Data.Cache.Streams.Active = append(Data.Cache.Streams.Active, m3uChannel.Name+m3uChannel.FileM3UID)

		// Try to find the channel based on matching all known values.  If that fails, then move to full channel scan
		m3uChannelHash := generateHashForChannel(m3uChannel.FileM3UID, m3uChannel.GroupTitle, m3uChannel.TvgID, m3uChannel.TvgName, m3uChannel.UUIDKey, m3uChannel.UUIDValue)
		if val, ok := xepgChannelsValuesMap[m3uChannelHash]; ok {
			channelExists = true
			currentXEPGID = val.XEPG
			if len(m3uChannel.UUIDValue) > 0 {
				channelHasUUID = true
			}
		} else {

			// XEPG Datenbank durchlaufen um nach dem Kanal zu suchen.  Run through the XEPG database to search for the channel (full scan)
			for _, dxc := range xepgChannelsValuesMap {

				if m3uChannel.FileM3UID == dxc.FileM3UID {

					dxc.FileM3UID = m3uChannel.FileM3UID
					dxc.FileM3UName = m3uChannel.FileM3UName

					// Vergleichen des Streams anhand einer UUID in der M3U mit dem Kanal in der Databank.  Compare the stream using a UUID in the M3U with the channel in the database
					if len(dxc.UUIDValue) > 0 && len(m3uChannel.UUIDValue) > 0 {

						if dxc.UUIDValue == m3uChannel.UUIDValue && dxc.UUIDKey == m3uChannel.UUIDKey {

							channelExists = true
							channelHasUUID = true
							currentXEPGID = dxc.XEPG
							break

						}

					} else {
						// Vergleichen des Streams mit dem Kanal in der Databank anhand des Kanalnamens.  Compare the stream to the channel in the database using the channel name
						if dxc.Name == m3uChannel.Name {
							channelExists = true
							currentXEPGID = dxc.XEPG
							break
						}

					}

				}

			}
		}

		switch channelExists {

		case true:
			// Bereits vorhandener Kanal
			var xepgChannel XEPGChannelStruct
			err = json.Unmarshal([]byte(mapToJSON(Data.XEPG.Channels[currentXEPGID])), &xepgChannel)
			if err != nil {
				return
			}

			// Streaming URL aktualisieren
			xepgChannel.URL = m3uChannel.URL

			// Name aktualisieren, anhand des Names wird überprüft ob der Kanal noch in einer Playlist verhanden. Funktion: cleanupXEPG
			xepgChannel.Name = m3uChannel.Name

			// Kanalname aktualisieren, nur mit Kanal ID's möglich
			if channelHasUUID == true {
				if xepgChannel.XUpdateChannelName == true {
					xepgChannel.XName = m3uChannel.Name
				}
			}

			// Kanallogo aktualisieren. Wird bei vorhandenem Logo in der XMLTV Datei wieder überschrieben
			if xepgChannel.XUpdateChannelIcon == true {
				xepgChannel.TvgLogo = m3uChannel.TvgLogo
			}

			Data.XEPG.Channels[currentXEPGID] = xepgChannel

		case false:
			// Neuer Kanal
			var xepg = createNewID()
			var xChannelID = getFreeChannelNumber()

			var newChannel XEPGChannelStruct
			newChannel.FileM3UID = m3uChannel.FileM3UID
			newChannel.FileM3UName = m3uChannel.FileM3UName
			newChannel.FileM3UPath = m3uChannel.FileM3UPath
			newChannel.Values = m3uChannel.Values
			newChannel.GroupTitle = m3uChannel.GroupTitle
			newChannel.Name = m3uChannel.Name
			newChannel.TvgID = m3uChannel.TvgID
			newChannel.TvgLogo = m3uChannel.TvgLogo
			newChannel.TvgName = m3uChannel.TvgName
			newChannel.URL = m3uChannel.URL
			newChannel.XmltvFile = ""
			newChannel.XMapping = ""

			if len(m3uChannel.UUIDKey) > 0 {
				newChannel.UUIDKey = m3uChannel.UUIDKey
				newChannel.UUIDValue = m3uChannel.UUIDValue
			}

			newChannel.XName = m3uChannel.Name
			newChannel.XGroupTitle = m3uChannel.GroupTitle
			newChannel.XEPG = xepg
			newChannel.XChannelID = xChannelID

			newChannel.XUpdateChannelIcon = Settings.ChannelDefaults.XUpdateChannelIcon
			newChannel.XUpdateChannelName = Settings.ChannelDefaults.XUpdateChannelName

			Data.XEPG.Channels[xepg] = newChannel

		}

	}
	showInfo("XEPG:" + "Save DB file")
	err = saveMapToJSONFile(System.File.XEPG, Data.XEPG.Channels)
	if err != nil {
		return
	}

	return
}

// Kanäle automatisch zuordnen und das Mapping überprüfen
func mapping() (err error) {
	showInfo("XEPG:" + "Map channels")

	for xepg, dxc := range Data.XEPG.Channels {

		var xepgChannel XEPGChannelStruct
		err = json.Unmarshal([]byte(mapToJSON(dxc)), &xepgChannel)
		if err != nil {
			return
		}

		// Automatische Mapping für neue Kanäle. Wird nur ausgeführt, wenn der Kanal deaktiviert ist und keine XMLTV Datei und kein XMLTV Kanal zugeordnet ist.
		if xepgChannel.XActive == false {

			// Werte kann "-" sein, deswegen len < 1
			if len(xepgChannel.XmltvFile) < 1 && len(xepgChannel.XmltvFile) < 1 {

				var tvgID = xepgChannel.TvgID

				// Default für neuen Kanal setzen
				xepgChannel.XmltvFile = "-"
				xepgChannel.XMapping = "-"

				Data.XEPG.Channels[xepg] = xepgChannel

				for file, xmltvChannels := range Data.XMLTV.Mapping {

					if channel, ok := xmltvChannels.(map[string]interface{})[tvgID]; ok {

						if channelID, ok := channel.(map[string]interface{})["id"].(string); ok {

							xepgChannel.XmltvFile = file
							xepgChannel.XMapping = channelID
							xepgChannel.XActive = true

							// Falls in der XMLTV Datei ein Logo existiert, wird dieses verwendet. Falls nicht, dann das Logo aus der M3U Datei
							if icon, ok := channel.(map[string]interface{})["icon"].(string); ok {
								if len(icon) > 0 {
									xepgChannel.TvgLogo = icon
								}
							}

							Data.XEPG.Channels[xepg] = xepgChannel
							break

						}

					}

				}

			}

		}

		// Überprüfen, ob die zugeordneten XMLTV Dateien und Kanäle noch existieren.
		if xepgChannel.XActive == true {

			var mapping = xepgChannel.XMapping
			var file = xepgChannel.XmltvFile

			if file != "xTeVe Dummy" {

				if value, ok := Data.XMLTV.Mapping[file].(map[string]interface{}); ok {

					if channel, ok := value[mapping].(map[string]interface{}); ok {

						// Kanallogo aktualisieren
						if logo, ok := channel["icon"].(string); ok {

							if xepgChannel.XUpdateChannelIcon == true && len(logo) > 0 {
								xepgChannel.TvgLogo = logo
							}

						}

					} else {

						ShowError(fmt.Errorf(fmt.Sprintf("Missing EPG data: %s", xepgChannel.Name)), 0)
						showWarning(2302)
						xepgChannel.XActive = false

					}

				} else {

					var fileID = strings.TrimSuffix(getFilenameFromPath(file), path.Ext(getFilenameFromPath(file)))

					ShowError(fmt.Errorf("Missing XMLTV file: %s", getProviderParameter(fileID, "xmltv", "name")), 0)
					showWarning(2301)
					xepgChannel.XActive = false

				}

			}

			if len(xepgChannel.XmltvFile) == 0 {
				xepgChannel.XmltvFile = "-"
				xepgChannel.XActive = false
			}

			if len(xepgChannel.XMapping) == 0 {
				xepgChannel.XMapping = "-"
				xepgChannel.XActive = false
			}

			Data.XEPG.Channels[xepg] = xepgChannel

		}

	}

	err = saveMapToJSONFile(System.File.XEPG, Data.XEPG.Channels)
	if err != nil {
		return
	}

	return
}

// XMLTV Datei erstellen
func createXMLTVFile() (err error) {

	// Image Cache
	// 4edd81ab7c368208cc6448b615051b37.jpg
	var imgc = Data.Cache.Images

	Data.Cache.ImagesFiles = []string{}
	Data.Cache.ImagesURLS = []string{}
	Data.Cache.ImagesCache = []string{}

	files, err := ioutil.ReadDir(System.Folder.ImagesCache)
	if err == nil {

		for _, file := range files {

			if indexOfString(file.Name(), Data.Cache.ImagesCache) == -1 {
				Data.Cache.ImagesCache = append(Data.Cache.ImagesCache, file.Name())
			}

		}

	}

	if len(Data.XMLTV.Files) == 0 && len(Data.Streams.Active) == 0 {
		Data.XEPG.Channels = make(map[string]interface{})
		return
	}

	showInfo("XEPG:" + fmt.Sprintf("Create XMLTV file (%s)", System.File.XML))

	var xepgXML XMLTV

	xepgXML.Generator = System.Name

	if System.Branch == "master" {
		xepgXML.Source = fmt.Sprintf("%s - %s", System.Name, System.Version)
	} else {
		xepgXML.Source = fmt.Sprintf("%s - %s.%s", System.Name, System.Version, System.Build)
	}

	var tmpProgram = &XMLTV{}

	for _, dxc := range Data.XEPG.Channels {

		var xepgChannel XEPGChannelStruct
		err := json.Unmarshal([]byte(mapToJSON(dxc)), &xepgChannel)
		if err == nil {

			if xepgChannel.XActive == true {

				// Kanäle
				var channel Channel
				channel.ID = xepgChannel.XChannelID
				channel.Icon = Icon{Src: imgc.Image.GetURL(xepgChannel.TvgLogo)}
				channel.DisplayName = append(channel.DisplayName, DisplayName{Value: xepgChannel.XName})

				xepgXML.Channel = append(xepgXML.Channel, &channel)

				// Programme

				*tmpProgram, err = getProgramData(xepgChannel)
				if err == nil {

					for _, program := range tmpProgram.Program {
						xepgXML.Program = append(xepgXML.Program, program)
					}

				}

			}

		}

	}

	var content, _ = xml.MarshalIndent(xepgXML, "  ", "    ")
	var xmlOutput = []byte(xml.Header + string(content))
	writeByteToFile(System.File.XML, xmlOutput)

	showInfo("XEPG:" + fmt.Sprintf("Compress XMLTV file (%s)", System.Compressed.GZxml))
	err = compressGZIP(&xmlOutput, System.Compressed.GZxml)

	xepgXML = XMLTV{}

	return
}

// Programmdaten erstellen (createXMLTVFile)
func getProgramData(xepgChannel XEPGChannelStruct) (xepgXML XMLTV, err error) {

	var xmltvFile = System.Folder.Data + xepgChannel.XmltvFile
	var channelID = xepgChannel.XMapping

	var xmltv XMLTV

	if xmltvFile == System.Folder.Data+"xTeVe Dummy" {
		xmltv = createDummyProgram(xepgChannel)
	} else {

		err = getLocalXMLTV(xmltvFile, &xmltv)
		if err != nil {
			return
		}

	}

	for _, xmltvProgram := range xmltv.Program {

		if xmltvProgram.Channel == channelID {
			//fmt.Println(&channelID)
			var program = &Program{}

			// Channel ID
			program.Channel = xepgChannel.XChannelID
			program.Start = xmltvProgram.Start
			program.Stop = xmltvProgram.Stop

			// Title
			program.Title = xmltvProgram.Title

			// Sub title (Untertitel)
			program.SubTitle = xmltvProgram.SubTitle

			// Description (Beschreibung)
			program.Desc = xmltvProgram.Desc

			// Category (Kategorie)
			getCategory(program, xmltvProgram, xepgChannel)

			// Credits : (Credits)
			program.Credits = xmltvProgram.Credits

			// Rating (Bewertung)
			program.Rating = xmltvProgram.Rating

			// StarRating (Bewertung / Kritiken)
			program.StarRating = xmltvProgram.StarRating

			// Country (Länder)
			program.Country = xmltvProgram.Country

			// Program icon (Poster / Cover)
			getPoster(program, xmltvProgram, xepgChannel)

			// Language (Sprache)
			program.Language = xmltvProgram.Language

			// Episodes numbers (Episodennummern)
			getEpisodeNum(program, xmltvProgram, xepgChannel)

			// Video (Videoparameter)
			getVideo(program, xmltvProgram, xepgChannel)

			// Date (Datum)
			program.Date = xmltvProgram.Date

			// Previously shown (Wiederholung)
			program.PreviouslyShown = xmltvProgram.PreviouslyShown

			// New (Neu)
			program.New = xmltvProgram.New

			// Live
			program.Live = xmltvProgram.Live

			// Premiere
			program.Premiere = xmltvProgram.Premiere

			xepgXML.Program = append(xepgXML.Program, program)

		}

	}

	return
}

// Dummy Daten erstellen (createXMLTVFile)
func createDummyProgram(xepgChannel XEPGChannelStruct) (dummyXMLTV XMLTV) {

	var imgc = Data.Cache.Images
	var currentTime = time.Now()
	var dateArray = strings.Fields(currentTime.String())
	var offset = " " + dateArray[2]
	var currentDay = currentTime.Format("20060102")
	var startTime, _ = time.Parse("20060102150405", currentDay+"000000")

	showInfo("Create Dummy Guide:" + "Time offset" + offset + " - " + xepgChannel.XName)

	var dl = strings.Split(xepgChannel.XMapping, "_")
	dummyLength, err := strconv.Atoi(dl[0])
	if err != nil {
		ShowError(err, 000)
		return
	}

	for d := 0; d < 4; d++ {

		var epgStartTime = startTime.Add(time.Hour * time.Duration(d*24))

		for t := dummyLength; t <= 1440; t = t + dummyLength {

			var epgStopTime = epgStartTime.Add(time.Minute * time.Duration(dummyLength))

			var epg Program
			poster := Poster{}

			epg.Channel = xepgChannel.XMapping
			epg.Start = epgStartTime.Format("20060102150405") + offset
			epg.Stop = epgStopTime.Format("20060102150405") + offset
			epg.Title = append(epg.Title, &Title{Value: xepgChannel.XName + " (" + epgStartTime.Weekday().String()[0:2] + ". " + epgStartTime.Format("15:04") + " - " + epgStopTime.Format("15:04") + ")", Lang: "en"})

			if len(xepgChannel.XDescription) == 0 {
				epg.Desc = append(epg.Desc, &Desc{Value: "xTeVe: (" + strconv.Itoa(dummyLength) + " Minutes) " + epgStartTime.Weekday().String() + " " + epgStartTime.Format("15:04") + " - " + epgStopTime.Format("15:04"), Lang: "en"})
			} else {
				epg.Desc = append(epg.Desc, &Desc{Value: xepgChannel.XDescription, Lang: "en"})
			}

			if Settings.XepgReplaceMissingImages == true {
				poster.Src = imgc.Image.GetURL(xepgChannel.TvgLogo)
				epg.Poster = append(epg.Poster, poster)
			}

			if xepgChannel.XCategory != "Movie" {
				epg.EpisodeNum = append(epg.EpisodeNum, &EpisodeNum{Value: epgStartTime.Format("2006-01-02 15:04:05"), System: "original-air-date"})
			}

			epg.New = &New{Value: ""}

			dummyXMLTV.Program = append(dummyXMLTV.Program, &epg)
			epgStartTime = epgStopTime

		}

	}

	return
}

// Kategorien erweitern (createXMLTVFile)
func getCategory(program *Program, xmltvProgram *Program, xepgChannel XEPGChannelStruct) {

	for _, i := range xmltvProgram.Category {

		category := &Category{}
		category.Value = i.Value
		category.Lang = i.Lang
		program.Category = append(program.Category, category)

	}

	if len(xepgChannel.XCategory) > 0 {

		category := &Category{}
		category.Value = xepgChannel.XCategory
		category.Lang = "en"
		program.Category = append(program.Category, category)

	}

	return
}

// Programm Poster Cover aus der XMLTV Datei laden
func getPoster(program *Program, xmltvProgram *Program, xepgChannel XEPGChannelStruct) {

	var imgc = Data.Cache.Images

	for _, poster := range xmltvProgram.Poster {
		poster.Src = imgc.Image.GetURL(poster.Src)
		program.Poster = append(program.Poster, poster)
	}

	if Settings.XepgReplaceMissingImages == true {

		if len(xmltvProgram.Poster) == 0 {
			var poster Poster
			poster.Src = imgc.Image.GetURL(poster.Src)
			program.Poster = append(program.Poster, poster)
		}

	}

}

// Episodensystem übernehmen, falls keins vorhanden ist und eine Kategorie im Mapping eingestellt wurden, wird eine Episode erstellt
func getEpisodeNum(program *Program, xmltvProgram *Program, xepgChannel XEPGChannelStruct) {

	program.EpisodeNum = xmltvProgram.EpisodeNum

	if len(xepgChannel.XCategory) > 0 && xepgChannel.XCategory != "Movie" {

		if len(xmltvProgram.EpisodeNum) == 0 {

			var timeLayout = "20060102150405"

			t, err := time.Parse(timeLayout, strings.Split(xmltvProgram.Start, " ")[0])
			if err == nil {
				program.EpisodeNum = append(program.EpisodeNum, &EpisodeNum{Value: t.Format("2006-01-02 15:04:05"), System: "original-air-date"})
			} else {
				ShowError(err, 0)
			}

		}

	}

	return
}

// Videoparameter erstellen (createXMLTVFile)
func getVideo(program *Program, xmltvProgram *Program, xepgChannel XEPGChannelStruct) {

	var video Video
	video.Present = xmltvProgram.Video.Present
	video.Colour = xmltvProgram.Video.Colour
	video.Aspect = xmltvProgram.Video.Aspect
	video.Quality = xmltvProgram.Video.Quality

	if len(xmltvProgram.Video.Quality) == 0 {

		if strings.Contains(strings.ToUpper(xepgChannel.XName), " HD") || strings.Contains(strings.ToUpper(xepgChannel.XName), " FHD") {
			video.Quality = "HDTV"
		}

		if strings.Contains(strings.ToUpper(xepgChannel.XName), " UHD") || strings.Contains(strings.ToUpper(xepgChannel.XName), " 4K") {
			video.Quality = "UHDTV"
		}

	}

	program.Video = video

	return
}

// Lokale Provider XMLTV Datei laden
func getLocalXMLTV(file string, xmltv *XMLTV) (err error) {

	if _, ok := Data.Cache.XMLTV[file]; !ok {

		// Cache initialisieren
		if len(Data.Cache.XMLTV) == 0 {
			Data.Cache.XMLTV = make(map[string]XMLTV)
		}

		// XML Daten lesen
		content, err := readByteFromFile(file)

		// Lokale XML Datei existiert nicht im Ordner: data
		if err != nil {
			ShowError(err, 1004)
			err = errors.New("Local copy of the file no longer exists")
			return err
		}

		// XML Datei parsen
		err = xml.Unmarshal(content, &xmltv)
		if err != nil {
			return err
		}

		Data.Cache.XMLTV[file] = *xmltv

	} else {
		*xmltv = Data.Cache.XMLTV[file]
	}

	return
}

// M3U Datei erstellen
func createM3UFile() {

	showInfo("XEPG:" + fmt.Sprintf("Create M3U file (%s)", System.File.M3U))
	_, err := buildM3U([]string{})
	if err != nil {
		ShowError(err, 000)
	}

	saveMapToJSONFile(System.File.URLS, Data.Cache.StreamingURLS)

	return
}

// XEPG Datenbank bereinigen
func cleanupXEPG() {

	//fmt.Println(Settings.Files.M3U)

	var sourceIDs []string

	for source := range Settings.Files.M3U {
		sourceIDs = append(sourceIDs, source)
	}

	for source := range Settings.Files.HDHR {
		sourceIDs = append(sourceIDs, source)
	}

	showInfo("XEPG:" + fmt.Sprintf("Cleanup database"))
	Data.XEPG.XEPGCount = 0

	for id, dxc := range Data.XEPG.Channels {

		var xepgChannel XEPGChannelStruct
		err := json.Unmarshal([]byte(mapToJSON(dxc)), &xepgChannel)
		if err == nil {

			if indexOfString(xepgChannel.Name+xepgChannel.FileM3UID, Data.Cache.Streams.Active) == -1 {
				delete(Data.XEPG.Channels, id)
			} else {
				if xepgChannel.XActive == true {
					Data.XEPG.XEPGCount++
				}
			}

			if indexOfString(xepgChannel.FileM3UID, sourceIDs) == -1 {
				delete(Data.XEPG.Channels, id)
			}

		}

	}

	err := saveMapToJSONFile(System.File.XEPG, Data.XEPG.Channels)
	if err != nil {
		ShowError(err, 000)
		return
	}

	showInfo("XEPG Channels:" + fmt.Sprintf("%d", Data.XEPG.XEPGCount))

	if len(Data.Streams.Active) > 0 && Data.XEPG.XEPGCount == 0 {
		showWarning(2005)
	}

	return
}

// Streaming URL für die Channels App generieren
func getStreamByChannelID(channelID string) (playlistID, streamURL string, err error) {

	err = errors.New("Channel not found")

	for _, dxc := range Data.XEPG.Channels {

		var xepgChannel XEPGChannelStruct
		err := json.Unmarshal([]byte(mapToJSON(dxc)), &xepgChannel)

		fmt.Println(xepgChannel.XChannelID)

		if err == nil {

			if channelID == xepgChannel.XChannelID {

				playlistID = xepgChannel.FileM3UID
				streamURL = xepgChannel.URL

				return playlistID, streamURL, nil
			}

		}

	}

	return
}
