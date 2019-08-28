package src

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

func makeInteraceFromHDHR(content []byte, playlistName, id string) (channels []interface{}, err error) {

	var hdhrData []interface{}

	err = json.Unmarshal(content, &hdhrData)
	if err == nil {

		for _, d := range hdhrData {

			var channel = make(map[string]string)
			var data = d.(map[string]interface{})

			channel["group-title"] = playlistName
			channel["name"] = data["GuideName"].(string)
			channel["tvg-id"] = data["GuideName"].(string)
			channel["url"] = data["URL"].(string)
			channel["ID-"+id] = data["GuideNumber"].(string)
			channel["_uuid.key"] = "ID-" + id
			channel["_values"] = playlistName + " " + channel["name"]

			channels = append(channels, channel)

		}

	}

	return
}

func getCapability() (xmlContent []byte, err error) {

	var capability Capability
	var buffer bytes.Buffer

	capability.Xmlns = "urn:schemas-upnp-org:device-1-0"
	capability.URLBase = System.ServerProtocol.WEB + "://" + System.Domain

	capability.SpecVersion.Major = 1
	capability.SpecVersion.Minor = 0

	capability.Device.DeviceType = "urn:schemas-upnp-org:device:MediaServer:1"
	capability.Device.FriendlyName = System.Name
	capability.Device.Manufacturer = "Silicondust"
	capability.Device.ModelName = "HDTC-2US"
	capability.Device.ModelNumber = "HDTC-2US"
	capability.Device.SerialNumber = ""
	capability.Device.UDN = "uuid:" + System.DeviceID

	output, err := xml.MarshalIndent(capability, " ", "  ")
	if err != nil {
		ShowError(err, 1003)
	}

	buffer.Write([]byte(xml.Header))
	buffer.Write([]byte(output))
	xmlContent = buffer.Bytes()

	return
}

func getDiscover() (jsonContent []byte, err error) {

	var discover Discover

	discover.BaseURL = System.ServerProtocol.WEB + "://" + System.Domain
	discover.DeviceAuth = System.AppName
	discover.DeviceID = System.DeviceID
	discover.FirmwareName = "bin_" + System.Version
	discover.FirmwareVersion = System.Version
	discover.FriendlyName = System.Name

	discover.LineupURL = fmt.Sprintf("%s://%s/lineup.json", System.ServerProtocol.DVR, System.Domain)
	discover.Manufacturer = "Golang"
	discover.ModelNumber = System.Version
	discover.TunerCount = Settings.Tuner

	jsonContent, err = json.MarshalIndent(discover, "", "  ")

	return
}

func getLineupStatus() (jsonContent []byte, err error) {

	var lineupStatus LineupStatus

	lineupStatus.ScanInProgress = System.ScanInProgress
	lineupStatus.ScanPossible = 0
	lineupStatus.Source = "Cable"
	lineupStatus.SourceList = []string{"Cable"}

	jsonContent, err = json.MarshalIndent(lineupStatus, "", "  ")

	return
}

func getLineup() (jsonContent []byte, err error) {

	var lineup Lineup

	switch Settings.EpgSource {

	case "PMS":
		for i, dsa := range Data.Streams.Active {

			var m3uChannel M3UChannelStructXEPG

			err = json.Unmarshal([]byte(mapToJSON(dsa)), &m3uChannel)
			if err != nil {
				return
			}

			var stream LineupStream
			stream.GuideName = m3uChannel.Name
			switch len(m3uChannel.UUIDValue) {

			case 0:
				stream.GuideNumber = fmt.Sprintf("%d", i+1000)
				guideNumber, err := getGuideNumberPMS(stream.GuideName)
				if err != nil {
					ShowError(err, 0)
				}

				stream.GuideNumber = guideNumber

			default:
				stream.GuideNumber = m3uChannel.UUIDValue

			}

			stream.URL, err = createStreamingURL("DVR", m3uChannel.FileM3UID, stream.GuideNumber, m3uChannel.Name, m3uChannel.URL)
			if err == nil {
				lineup = append(lineup, stream)
			} else {
				ShowError(err, 1202)
			}

		}

	case "XEPG":
		for _, dxc := range Data.XEPG.Channels {

			var xepgChannel XEPGChannelStruct
			err = json.Unmarshal([]byte(mapToJSON(dxc)), &xepgChannel)
			if err != nil {
				return
			}

			if xepgChannel.XActive == true {
				var stream LineupStream
				stream.GuideName = xepgChannel.XName
				stream.GuideNumber = xepgChannel.XChannelID
				//stream.URL = fmt.Sprintf("%s://%s/stream/%s-%s", System.ServerProtocol.DVR, System.Domain, xepgChannel.FileM3UID, base64.StdEncoding.EncodeToString([]byte(xepgChannel.URL)))
				stream.URL, err = createStreamingURL("DVR", xepgChannel.FileM3UID, xepgChannel.XChannelID, xepgChannel.XName, xepgChannel.URL)
				if err == nil {
					lineup = append(lineup, stream)
				} else {
					ShowError(err, 1202)
				}

			}

		}

	}

	jsonContent, err = json.MarshalIndent(lineup, "", "  ")

	Data.Cache.PMS = nil

	saveMapToJSONFile(System.File.URLS, Data.Cache.StreamingURLS)

	return
}

func getGuideNumberPMS(channelName string) (pmsID string, err error) {

	if len(Data.Cache.PMS) == 0 {

		Data.Cache.PMS = make(map[string]string)

		pms, err := loadJSONFileToMap(System.File.PMS)

		if err != nil {
			return "", err
		}

		for key, value := range pms {
			Data.Cache.PMS[key] = value.(string)
		}

	}

	var getNewID = func(channelName string) (id string) {

		var i int

	newID:

		var ids []string
		id = fmt.Sprintf("id-%d", i)

		for _, v := range Data.Cache.PMS {
			ids = append(ids, v)
		}

		if indexOfString(id, ids) != -1 {
			i++
			goto newID
		}

		return
	}

	if value, ok := Data.Cache.PMS[channelName]; ok {

		pmsID = value

	} else {

		pmsID = getNewID(channelName)
		Data.Cache.PMS[channelName] = pmsID
		saveMapToJSONFile(System.File.PMS, Data.Cache.PMS)

	}

	return
}
