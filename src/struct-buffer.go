package src

import "time"

// Playlist : Enthält allen Playlistinformationen, die der Buffer benötigr
type Playlist struct {
	Folder       string
	PlaylistID   string
	PlaylistName string
	Tuner        int

	Clients map[int]ThisClient
	Streams map[int]ThisStream
}

// ThisClient : Clientinfos
type ThisClient struct {
	Connection int
}

// ThisStream : Enthält Informationen zu dem abzuspielenden Stream einer Playlist
type ThisStream struct {
	ChannelName      string
	Error            string
	Folder           string
	MD5              string
	NetworkBandwidth int
	PlaylistID       string
	PlaylistName     string
	Status           bool
	URL              string

	Segment []Segment

	// Serverinformationen
	Location           string
	URLFile            string
	URLHost            string
	URLPath            string
	URLRedirect        string
	URLScheme          string
	URLStreamingServer string

	// Wird nur für HLS / M3U8 verwendet
	Body             string
	Difference       float64
	Duration         float64
	DynamicBandwidth bool
	FirstSequence    int64
	HLS              bool
	LastSequence     int64
	M3U8URL          string
	NewSegCount      int
	OldSegCount      int
	Sequence         int64
	TimeDiff         float64
	TimeEnd          time.Time
	TimeSegDuration  float64
	TimeStart        time.Time
	Version          int
	Wait             float64

	DynamicStream map[int]DynamicStream

	// Lokale Temp Datein
	OldSegments []string
}

// Segment : URL Segmente (HLS / M3U8)
type Segment struct {
	Duration     float64
	Info         bool
	PlaylistType string
	Sequence     int64
	URL          string
	Version      int
	Wait         float64

	StreamInf struct {
		AverageBandwidth int
		Bandwidth        int
		Framerate        float64
		Resolution       string
		SegmentURL       string
	}
}

// DynamicStream : Streaminformationen bei dynamischer Bandbreite
type DynamicStream struct {
	AverageBandwidth int
	Bandwidth        int
	Framerate        float64
	Resolution       string
	URL              string
}

// ClientConnection : Client Verbindungen
type ClientConnection struct {
	Connection int
	Error      error
}

// BandwidthCalculation : Bandbreitenberechnung für den Stream
type BandwidthCalculation struct {
	NetworkBandwidth int
	Size             int
	Start            time.Time
	Stop             time.Time
	TimeDiff         float64
}
