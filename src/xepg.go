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

// Check provider XMLTV File
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

// Create XEPG Data
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

				// Clearing the Cache
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

				// Clearing the Cache
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

// Create Mapping Menu for the XMLTV Files
func createXEPGMapping() {

	Data.XMLTV.Files = getLocalProviderFiles("xmltv")
	Data.XMLTV.Mapping = make(map[string]interface{})

	var tmpMap = make(map[string]interface{})

	var getFriendlyName = func(channel Channel) (friendlyName string) {
		switch len(channel.DisplayNames) {
		case 1:
			friendlyName = channel.DisplayNames[0].Value
		default:
			friendlyName = fmt.Sprintf("%s (%s)", channel.DisplayNames[1].Value, channel.DisplayNames[0].Value)
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

			// XML Parsing (Provider File)
			if err == nil {

				// Write Data from the XML File to a temporary Map
				var xmltvMap = make(map[string]interface{})

				for _, c := range xmltv.Channel {
					var channel = make(map[string]interface{})

					channel["id"] = c.ID
					channel["display-names"] = c.DisplayNames
					channel["friendly-name"] = getFriendlyName(*c)
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

	// Create selection for the Dummy
	var dummy = make(map[string]interface{})
	var times = []string{"30", "60", "90", "120", "180", "240", "360"}

	for _, i := range times {

		var dummyChannel = make(map[string]interface{})
		dummyChannel["friendly-name"] = i + " Minutes"
		dummyChannel["display-names"] = []DisplayName{{Value: i + " Minutes"}}
		dummyChannel["id"] = i + "_Minutes"
		dummyChannel["icon"] = ""

		dummy[dummyChannel["id"].(string)] = dummyChannel

	}

	Data.XMLTV.Mapping["xTeVe Dummy"] = dummy

	return
}

// Create / update XEPG Database
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

	var getFreeChannelNumber = func(startingChannel ...string) (xChannelID string) {

		sort.Float64s(allChannelNumbers)

		var firstFreeNumber float64 = Settings.MappingFirstChannel
		if startingChannel != nil {
			var startingChannel, _ = strconv.ParseFloat(startingChannel[0], 64)
			if startingChannel > 0 {
				firstFreeNumber = startingChannel
			}
		}

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

	var generateHashForChannel = func(m3uID string, name string, groupTitle string, tvgID string, tvgName string, uuidKey string, uuidValue string) string {
		hash := md5.Sum([]byte(m3uID + name + groupTitle + tvgID + tvgName + uuidKey + uuidValue))
		return hex.EncodeToString(hash[:])
	}

	showInfo("XEPG:" + "Update database")

	// Delete Channel with missing Channel Numbers.
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
		channelHash := generateHashForChannel(channel.FileM3UID, channel.Name, channel.GroupTitle, channel.TvgID, channel.TvgName, channel.UUIDKey, channel.UUIDValue)
		xepgChannelsValuesMap[channelHash] = channel
	}

	for _, dsa := range Data.Streams.Active {

		var channelExists = false  // Decides whether a Channel should be added to the Database
		var channelHasUUID = false // Checks whether the Channel (Stream) has Unique IDs
		var currentXEPGID string   // Current Database ID (XEPG) Used to update the Channel in the Database with the Stream of the M3U

		var m3uChannel M3UChannelStructXEPG

		err = json.Unmarshal([]byte(mapToJSON(dsa)), &m3uChannel)
		if err != nil {
			return
		}

		Data.Cache.Streams.Active = append(Data.Cache.Streams.Active, m3uChannel.Name+m3uChannel.FileM3UID)

		// Try to find the channel based on matching all known values.  If that fails, then move to full channel scan
		m3uChannelHash := generateHashForChannel(m3uChannel.FileM3UID, m3uChannel.Name, m3uChannel.GroupTitle, m3uChannel.TvgID, m3uChannel.TvgName, m3uChannel.UUIDKey, m3uChannel.UUIDValue)
		if val, ok := xepgChannelsValuesMap[m3uChannelHash]; ok {
			channelExists = true
			currentXEPGID = val.XEPG
			if len(m3uChannel.UUIDValue) > 0 {
				channelHasUUID = true
			}
		} else {

			// Run through the XEPG Database to search for the Channel (full scan)
			for _, dxc := range xepgChannelsValuesMap {

				if m3uChannel.FileM3UID == dxc.FileM3UID {

					dxc.FileM3UID = m3uChannel.FileM3UID
					dxc.FileM3UName = m3uChannel.FileM3UName

					// Compare the Stream using a UUID in the M3U with the Channel in the Database
					if len(dxc.UUIDValue) > 0 && len(m3uChannel.UUIDValue) > 0 {

						if dxc.UUIDValue == m3uChannel.UUIDValue && dxc.UUIDKey == m3uChannel.UUIDKey {

							channelExists = true
							channelHasUUID = true
							currentXEPGID = dxc.XEPG
							break

						}

					} else {
						// Compare the Stream to the Channel in the Database using the Channel Name
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
			// Existing Channel
			var xepgChannel XEPGChannelStruct
			err = json.Unmarshal([]byte(mapToJSON(Data.XEPG.Channels[currentXEPGID])), &xepgChannel)
			if err != nil {
				return
			}

			// Update Streaming URL
			xepgChannel.URL = m3uChannel.URL

			// Update Name, the Name is used to check whether the Channel is still available in a Playlist. Function: cleanupXEPG
			xepgChannel.Name = m3uChannel.Name

			// Update Channel Name, only possible with Channel ID's
			if channelHasUUID == true {
				if xepgChannel.XUpdateChannelName == true {
					xepgChannel.XName = m3uChannel.Name
				}
			}

			// Update Channel Logo. Will be overwritten again if the Logo is present in the XMLTV file
			if xepgChannel.XUpdateChannelIcon == true {
				xepgChannel.TvgLogo = m3uChannel.TvgLogo
			}

			Data.XEPG.Channels[currentXEPGID] = xepgChannel

		case false:
			// New Channel
			var xepg = createNewID()
			xChannelID := func() string {
				if m3uChannel.PreserveMapping == "true" {
					return getFreeChannelNumber(m3uChannel.UUIDValue)
				} else {
					return getFreeChannelNumber(m3uChannel.StartingChannel)
				}
			}()
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
			if m3uChannel.TvgShift == "" {
				newChannel.TvgShift = "0"
			} else {
				newChannel.TvgShift = m3uChannel.TvgShift
			}
			newChannel.URL = m3uChannel.URL
			newChannel.XmltvFile = ""
			newChannel.XMapping = ""
			newChannel.DefaultMissingEPG = m3uChannel.DefaultMissingEPG

			if len(m3uChannel.UUIDKey) > 0 {
				newChannel.UUIDKey = m3uChannel.UUIDKey
				newChannel.UUIDValue = m3uChannel.UUIDValue
			}

			newChannel.XName = m3uChannel.Name
			newChannel.XGroupTitle = m3uChannel.GroupTitle
			newChannel.XEPG = xepg
			newChannel.XChannelID = xChannelID
			newChannel.XTimeshift = newChannel.TvgShift

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

// Automatically assign Channels and check the Mapping
func mapping() (err error) {
	showInfo("XEPG:" + "Map channels")

	for xepg, dxc := range Data.XEPG.Channels {

		var xepgChannel XEPGChannelStruct
		err = json.Unmarshal([]byte(mapToJSON(dxc)), &xepgChannel)
		if err != nil {
			return
		}

		// Automatic mapping for new Channels. Is only executed if the Channel is deactivated and no XMLTV file and no XMLTV Channel is assigned.
		if xepgChannel.XActive == false {

			// Values can be "-", therefore len <= 1
			// If either XmltvFile (XMLTV file / EPG source) or XMapping (XMLTV Channel / EPG program) is "-" or null, then look for a matching EPG program.
			// If nothing matches, look for DefaultMissingEPG (a default Dummy xTeVe preference) and set it
			if len(xepgChannel.XmltvFile) <= 1 || len(xepgChannel.XMapping) <= 1 {

				var tvgID = xepgChannel.TvgID

				// Set default for new Channel
				if len(xepgChannel.DefaultMissingEPG) > 1 {
					xepgChannel.XmltvFile = "xTeVe Dummy"
					xepgChannel.XMapping = xepgChannel.DefaultMissingEPG
					xepgChannel.XActive = true
				} else {
					xepgChannel.XmltvFile = "-"
					xepgChannel.XMapping = "-"
				}

				Data.XEPG.Channels[xepg] = xepgChannel

			xmltvMapLoop:
				for file, xmltvChannels := range Data.XMLTV.Mapping {

					if channel, ok := xmltvChannels.(map[string]interface{})[tvgID]; ok {

						if channelID, ok := channel.(map[string]interface{})["id"].(string); ok {

							xepgChannel.XmltvFile = file
							xepgChannel.XMapping = channelID
							xepgChannel.XActive = true

							// If there is a Logo in the XMLTV file, this will be used. If not, then the Logo from the M3U file
							if icon, ok := channel.(map[string]interface{})["icon"].(string); ok {
								if len(icon) > 0 {
									xepgChannel.TvgLogo = icon
								}
							}

							Data.XEPG.Channels[xepg] = xepgChannel
							break

						}

					} else {

						// Search for the proper XEPG channel ID by comparing it's name with every alias in XML file
						for _, xmltvChannel := range xmltvChannels.(map[string]interface{}) {
							xmltvNames := xmltvChannel.(map[string]interface{})["display-names"].([]DisplayName)

							for _, xmltvName := range xmltvNames {
								xmltvNameSolid := strings.ReplaceAll(xmltvName.Value, " ", "")
								xepgNameSolid := strings.ReplaceAll(xepgChannel.Name, " ", "")

								if strings.EqualFold(xmltvNameSolid, xepgNameSolid) {
									xepgChannel.XmltvFile = file
									xepgChannel.XMapping = xmltvChannel.(map[string]interface{})["id"].(string)
									// xepgChannel.XActive = true

									// If there is a Logo in the XMLTV file, this will be used.
									// If not, then the Logo from the M3U file.
									if icon, ok := xmltvChannel.(map[string]interface{})["icon"].(string); ok {
										if len(icon) > 0 {
											xepgChannel.TvgLogo = icon
										}
									}

									break xmltvMapLoop
								}

							}

						}

					}

				}

				Data.XEPG.Channels[xepg] = xepgChannel

			}

		}

		// Check whether the assigned XMLTV Files and Channels still exist.
		if xepgChannel.XActive == true {

			var mapping = xepgChannel.XMapping
			var file = xepgChannel.XmltvFile

			if file != "xTeVe Dummy" {

				if value, ok := Data.XMLTV.Mapping[file].(map[string]interface{}); ok {

					if channel, ok := value[mapping].(map[string]interface{}); ok {

						// Update Channel Logo
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

			if len(xepgChannel.DefaultMissingEPG) > 1 && xepgChannel.XActive == false {
				xepgChannel.XmltvFile = "xTeVe Dummy"
				xepgChannel.XMapping = xepgChannel.DefaultMissingEPG
				xepgChannel.XActive = true
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

// Create XMLTV File
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

				// Channels
				var channel Channel
				channel.ID = xepgChannel.XChannelID
				channel.Icon = Icon{Src: imgc.Image.GetURL(xepgChannel.TvgLogo)}
				channel.DisplayNames = append(channel.DisplayNames, DisplayName{Value: xepgChannel.XName})

				xepgXML.Channel = append(xepgXML.Channel, &channel)

				// Programs

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

// Create Program Data (createXMLTVFile)
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
			timeshift, _ := strconv.Atoi(xepgChannel.XTimeshift)
			progStart := strings.Split(xmltvProgram.Start, " ")
			progStop := strings.Split(xmltvProgram.Stop, " ")
			tzStart, _ := strconv.Atoi(progStart[1])
			tzStop, _ := strconv.Atoi(progStop[1])
			progStart[1] = fmt.Sprintf("%+05d", tzStart+timeshift*100)
			progStop[1] = fmt.Sprintf("%+05d", tzStop+timeshift*100)
			program.Start = strings.Join(progStart, " ")
			program.Stop = strings.Join(progStop, " ")

			// Title
			program.Title = xmltvProgram.Title

			// Subtitle
			program.SubTitle = xmltvProgram.SubTitle

			// Description
			program.Desc = xmltvProgram.Desc

			// Category
			getCategory(program, xmltvProgram, xepgChannel)

			// Credits
			program.Credits = xmltvProgram.Credits

			// Rating
			program.Rating = xmltvProgram.Rating

			// StarRating
			program.StarRating = xmltvProgram.StarRating

			// Country
			program.Country = xmltvProgram.Country

			// Program icon
			getPoster(program, xmltvProgram, xepgChannel)

			// Language
			program.Language = xmltvProgram.Language

			// Episodes numbers
			getEpisodeNum(program, xmltvProgram, xepgChannel)

			// Video
			getVideo(program, xmltvProgram, xepgChannel)

			// Date
			program.Date = xmltvProgram.Date

			// Previously shown
			program.PreviouslyShown = xmltvProgram.PreviouslyShown

			// New
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

// Create Dummy Data (createXMLTVFile)
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

// Expand Categories (createXMLTVFile)
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

// Load the Poster Cover Program from the XMLTV File
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

// Apply Episode system, if none is available and a Category has been set in the mapping, an Episode is created
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

// Create Video Parameters (createXMLTVFile)
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

// Load Local Provider XMLTV file
func getLocalXMLTV(file string, xmltv *XMLTV) (err error) {

	if _, ok := Data.Cache.XMLTV[file]; !ok {

		// Initialize Cache
		if len(Data.Cache.XMLTV) == 0 {
			Data.Cache.XMLTV = make(map[string]XMLTV)
		}

		// Read XML Data
		content, err := readByteFromFile(file)

		// Local XML File does not exist in the folder: Data
		if err != nil {
			ShowError(err, 1004)
			err = errors.New("Local copy of the file no longer exists")
			return err
		}

		// Parse XML File
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

// Create M3U File
func createM3UFile() {

	showInfo("XEPG:" + fmt.Sprintf("Create M3U file (%s)", System.File.M3U))
	_, err := buildM3U([]string{})
	if err != nil {
		ShowError(err, 000)
	}

	saveMapToJSONFile(System.File.URLS, Data.Cache.StreamingURLS)

	return
}

// Clean up the XEPG Database
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
