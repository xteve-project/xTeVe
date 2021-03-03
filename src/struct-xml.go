package src

import "encoding/xml"

// XMLTV : XMLTV Datei
type XMLTV struct {
	Generator string   `xml:"generator-info-name,attr"`
	Source    string   `xml:"source-info-name,attr"`
	XMLName   xml.Name `xml:"tv"`

	Channel []*Channel `xml:"channel"`
	Program []*Program `xml:"programme"`
}

// Channel : Kanäle
type Channel struct {
	ID          string        `xml:"id,attr"`
	DisplayName []DisplayName `xml:"display-name"`
	Icon        Icon          `xml:"icon"`
}

// DisplayName : Kanalname
type DisplayName struct {
	Value string `xml:",chardata"`
}

// Icon : Senderlogo
type Icon struct {
	Src string `xml:"src,attr"`
}

// Program : Programme
type Program struct {
	Channel string `xml:"channel,attr"`
	Start   string `xml:"start,attr"`
	Stop    string `xml:"stop,attr"`

	Title           []*Title         `xml:"title"`
	SubTitle        []*SubTitle      `xml:"sub-title"`
	Desc            []*Desc          `xml:"desc"`
	Category        []*Category      `xml:"category"`
	Country         []*Country       `xml:"country"`
	EpisodeNum      []*EpisodeNum    `xml:"episode-num"`
	Poster          []Poster         `xml:"icon"`
	Language        []*Language      `xml:"language"`
	Video           Video            `xml:"video"`
	Date            string           `xml:"date"`
	PreviouslyShown *PreviouslyShown `xml:"previously-shown"`
	New             *New             `xml:"new"`
	Live            *Live            `xml:"live"`
	Premiere        *Live            `xml:"premiere"`
	StarRating      []*StarRating    `xml:"star-rating"`
}

// Title : Programmtitel
type Title struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// SubTitle : Kurzbeschreibung
type SubTitle struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

//Desc : Programmbeschreibung
type Desc struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// Category : Kategorien
type Category struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// Language : Sprachen
type Language struct {
	Value string `xml:",chardata"`
}

// Country : Länder
type Country struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// EpisodeNum : Episodennummerierung
type EpisodeNum struct {
	System string `xml:"system,attr"`
	Value  string `xml:",chardata"`
}

// StarRating : star-rating
type StarRating struct {
	System string `xml:"system,attr"`
	Value  string `xml:"value,omitempty"`
}

// Poster : Programmposter / Cover
type Poster struct {
	Height string `xml:"height,attr"`
	Src    string `xml:"src,attr"`
	Value  string `xml:",chardata"`
	Width  string `xml:"width,attr"`
}

// Video : Video Metadaten
type Video struct {
	Aspect  string `xml:"aspect,omitempty"`
	Colour  string `xml:"colour,omitempty"`
	Present string `xml:"present,omitempty"`
	Quality string `xml:"quality,omitempty"`
}

// PreviouslyShown : Widerholung bzw. Erstausstrahlung
type PreviouslyShown struct {
	Start string `xml:"start,attr"`
}

// New : Sendung als neu deklarieren
type New struct {
	Value string `xml:",chardata"`
}

// Live : Sendung als Liveübertragung deklarieren
type Live struct {
	Value string `xml:",chardata"`
}
