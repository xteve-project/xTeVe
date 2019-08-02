package src

import "encoding/xml"

// Capability : HDHR Capability XML
type Capability struct {
	URLBase string   `xml:"URLBase"`
	XMLName xml.Name `xml:"root"`
	Xmlns   string   `xml:"xmlns,attr"`

	SpecVersion struct {
		Major int `xml:"major"`
		Minor int `xml:"minor"`
	} `xml:"specVersion"`

	Device struct {
		DeviceType   string `xml:"deviceType"`
		FriendlyName string `xml:"friendlyName"`
		Manufacturer string `xml:"manufacturer"`
		ModelName    string `xml:"modelName"`
		ModelNumber  string `xml:"modelNumber"`
		SerialNumber string `xml:"serialNumber"`
		UDN          string `xml:"UDN"`
	} `xml:"device"`
}

// Discover : HDHR Discover /discover.json
type Discover struct {
	BaseURL         string `json:"BaseURL"`
	DeviceAuth      string `json:"DeviceAuth"`
	DeviceID        string `json:"DeviceID"`
	FirmwareName    string `json:"FirmwareName"`
	FirmwareVersion string `json:"FirmwareVersion"`
	FriendlyName    string `json:"FriendlyName"`
	LineupURL       string `json:"LineupURL"`
	Manufacturer    string `json:"Manufacturer"`
	ModelNumber     string `json:"ModelNumber"`
	TunerCount      int    `json:"TunerCount"`
}

// LineupStatus : HDHR Lineup status /lineup_status.json
type LineupStatus struct {
	ScanInProgress int      `json:"ScanInProgress"`
	ScanPossible   int      `json:"ScanPossible"`
	Source         string   `json:"Source"`
	SourceList     []string `json:"SourceList"`
}

// Lineup : HDHR Lineup /lineup.json
type Lineup []interface {
	//GuideName string `json:"GuideName"`
	//GuideNumber string `json:"GuideNumber"`
	//URL         string `json:"URL"`
}

// LineupStream : HDHR einzelner Stream im Lineup
type LineupStream struct {
	GuideName   string `json:"GuideName"`
	GuideNumber string `json:"GuideNumber"`
	URL         string `json:"URL"`
}
