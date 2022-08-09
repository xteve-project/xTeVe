package src

import "encoding/xml"

// XMLTV : XMLTV File
type XMLTV struct {
	Generator string   `xml:"generator-info-name,attr"`
	Source    string   `xml:"source-info-name,attr"`
	XMLName   xml.Name `xml:"tv"`

	Channel []*Channel `xml:"channel"`
	Program []*Program `xml:"programme"`
}

// Channel : Channels
type Channel struct {
	ID           string        `xml:"id,attr"`
	DisplayNames []DisplayName `xml:"display-name"`
	Icon         Icon          `xml:"icon"`
}

// DisplayName : Channel Name
type DisplayName struct {
	Value string `xml:",chardata"`
}

// Icon : Station Logo
type Icon struct {
	Src string `xml:"src,attr"`
}

// Program : Programs
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
	Credits         Credits          `xml:"credits,omitempty"` //`xml:",innerxml,omitempty"`
	Rating          []Rating         `xml:"rating"`
	StarRating      []StarRating     `xml:"star-rating"`
	Language        []*Language      `xml:"language"`
	Video           Video            `xml:"video"`
	Date            string           `xml:"date"`
	PreviouslyShown *PreviouslyShown `xml:"previously-shown"`
	New             *New             `xml:"new"`
	Live            *Live            `xml:"live"`
	Premiere        *Live            `xml:"premiere"`
}

// Title : Program Title
type Title struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// SubTitle : Brief Description
type SubTitle struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

//Desc : Program Description
type Desc struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// Category : Categories
type Category struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// Rating : Rating
type Rating struct {
	System string `xml:"system,attr"`
	Value  string `xml:"value"`
	Icon   []Icon `xml:"icon"`
}

// StarRating : Rating / Reviews
type StarRating struct {
	Value  string `xml:"value"`
	System string `xml:"system,attr"`
}

// Language : Langueages
type Language struct {
	Value string `xml:",chardata"`
}

// Country : Countries
type Country struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// EpisodeNum : Episode Numbering
type EpisodeNum struct {
	System string `xml:"system,attr"`
	Value  string `xml:",chardata"`
}

// Poster : Program Poster / Cover
type Poster struct {
	Height string `xml:"height,attr"`
	Src    string `xml:"src,attr"`
	Value  string `xml:",chardata"`
	Width  string `xml:"width,attr"`
}

// Credits : Credits
type Credits struct {
	Director  []Director  `xml:"director,omitempty"`
	Actor     []Actor     `xml:"actor,omitempty"`
	Writer    []Writer    `xml:"writer,omitempty"`
	Presenter []Presenter `xml:"presenter,omitempty"`
	Producer  []Producer  `xml:"producer,omitempty"`
}

// Director : Director
type Director struct {
	Value string `xml:",chardata"`
}

// Actor : Actor
type Actor struct {
	Value string `xml:",chardata"`
	Role  string `xml:"role,attr,omitempty"`
}

// Writer : Writer
type Writer struct {
	Value string `xml:",chardata"`
}

// Presenter : Presenter
type Presenter struct {
	Value string `xml:",chardata"`
}

// Producer : Producer
type Producer struct {
	Value string `xml:",chardata"`
}

// Video : Video Metadata
type Video struct {
	Aspect  string `xml:"aspect,omitempty"`
	Colour  string `xml:"colour,omitempty"`
	Present string `xml:"present,omitempty"`
	Quality string `xml:"quality,omitempty"`
}

// PreviouslyShown : Repetition or first Broadcast
type PreviouslyShown struct {
	Start string `xml:"start,attr"`
}

// New : Declare the Broadcast as new
type New struct {
	Value string `xml:",chardata"`
}

// Live : Declare the Broadcast as a Live Broadcast
type Live struct {
	Value string `xml:",chardata"`
}
