package src

import (
	"net"
	"xteve/src/internal/imgcache"
)

// SystemStruct : Contains all System Information
type SystemStruct struct {
	Addresses struct {
		DVR string
		M3U string
		XML string
	}

	APIVersion          string
	AppName             string
	ARCH                string
	BackgroundProcess   bool
	Branch              string
	Build               string
	Compatibility       string
	ConfigurationWizard bool
	DBVersion           string
	Dev                 bool
	DeviceID            string
	Domain              string
	PlexChannelLimit    int

	FFmpeg struct {
		DefaultOptions string
		Path           string
	}

	VLC struct {
		DefaultOptions string
		Path           string
	}

	File struct {
		Authentication    string
		M3U               string
		PMS               string
		ServerCert        string
		ServerCertPrivKey string
		Settings          string
		URLS              string
		XEPG              string
		XML               string
	}

	Compressed struct {
		GZxml string
	}

	Flag struct {
		Branch  string
		Debug   int
		Info    bool
		Port    string
		Restore string
		SSDP    bool
	}

	Folder struct {
		Backup       string
		Cache        string
		Certificates string
		Config       string
		Data         string
		ImagesCache  string
		ImagesUpload string
		Temp         string
	}

	Hostname               string
	ImageCachingInProgress int
	IPAddressesList        []string // Every IP address available (IPv4 + IPv6)
	IPAddressesV4          []string // Every IPv4 address available in string format
	IPAddressesV4Host      []string // Every IPv4 address available except loopback and link-local
	IPAddressesV4Raw       []net.IP // Every IPv4 address available in net.IP format
	IPAddressesV6          []string // Every IPv6 address available
	Name                   string
	OS                     string
	ScanInProgress         int
	TimeForAutoUpdate      string

	Notification map[string]Notification

	ServerProtocol struct {
		API string
		DVR string
		M3U string
		WEB string
		XML string
	}

	GitHub struct {
		Branch string
		Repo   string
		Update bool
		User   string
	}

	Update struct {
		Git  string
		Name string
	}

	URLBase string
	UDPxy   string
	Version string
	WEB     struct {
		Menu []string
	}
}

// GitStruct : Update information from GitHub
type GitStruct struct {
	Filename string `json:"filename"`
	Version  string `json:"version"`
}

// DataStruct : All Data is stored here. (Lineup, XMLTV)
type DataStruct struct {
	Cache struct {
		Images      *imgcache.Cache
		ImagesCache []string
		ImagesFiles []string
		ImagesURLS  []string
		PMS         map[string]string

		StreamingURLS map[string]StreamInfo
		XMLTV         map[string]XMLTV

		Streams struct {
			Active []string
		}
	}

	Filter []Filter

	Playlist struct {
		M3U struct {
			Groups struct {
				Text  []string
				Value []string
			}
		}
	}

	StreamPreviewUI struct {
		Active   []string
		Inactive []string
	}

	Streams struct {
		Active   []interface{}
		All      []interface{}
		Inactive []interface{}
	}

	XMLTV struct {
		Files   []string
		Mapping map[string]interface{}
	}

	XEPG struct {
		Channels  map[string]interface{}
		XEPGCount int64
	}
}

// Filter : Used for the Filter Rules
type Filter struct {
	CaseSensitive   bool
	PreserveMapping bool
	Rule            string
	Type            string
	StartingChannel string
}

// XEPGChannelStruct : XEPG Structure
type XEPGChannelStruct struct {
	FileM3UID                     string `json:"_file.m3u.id"`
	FileM3UName                   string `json:"_file.m3u.name"`
	FileM3UPath                   string `json:"_file.m3u.path"`
	GroupTitle                    string `json:"group-title"`
	Name                          string `json:"name"`
	TvgID                         string `json:"tvg-id"`
	TvgLogo                       string `json:"tvg-logo"`
	TvgName                       string `json:"tvg-name"`
	TvgShift                      string `json:"tvg-shift"`
	UpdateChannelNameRegex        string `json:"update-channel-name-regex"`
	UpdateChannelNameByGroupRegex string `json:"update-channel-name-by-group-regex"`
	URL                           string `json:"url"`
	UUIDKey                       string `json:"_uuid.key"`
	UUIDValue                     string `json:"_uuid.value,omitempty"`
	Values                        string `json:"_values"`
	XActive                       bool   `json:"x-active"`
	XCategory                     string `json:"x-category"`
	XChannelID                    string `json:"x-channelID"`
	XEPG                          string `json:"x-epg"`
	XGroupTitle                   string `json:"x-group-title"`
	XMapping                      string `json:"x-mapping"`
	XmltvFile                     string `json:"x-xmltv-file"`
	XName                         string `json:"x-name"`
	XUpdateChannelIcon            bool   `json:"x-update-channel-icon"`
	XUpdateChannelName            bool   `json:"x-update-channel-name"`
	XUpdateChannelGroup           bool   `json:"x-update-channel-group"`
	XDescription                  string `json:"x-description"`
	XTimeshift                    string `json:"x-timeshift"`
}

// M3UChannelStructXEPG : M3U Structure for XEPG
type M3UChannelStructXEPG struct {
	FileM3UID       string `json:"_file.m3u.id"`
	FileM3UName     string `json:"_file.m3u.name"`
	FileM3UPath     string `json:"_file.m3u.path"`
	GroupTitle      string `json:"group-title"`
	Name            string `json:"name"`
	TvgID           string `json:"tvg-id"`
	TvgLogo         string `json:"tvg-logo"`
	TvgName         string `json:"tvg-name"`
	TvgShift        string `json:"tvg-shift"`
	URL             string `json:"url"`
	UUIDKey         string `json:"_uuid.key"`
	UUIDValue       string `json:"_uuid.value"`
	Values          string `json:"_values"`
	PreserveMapping string `json:"_preserve-mapping"`
	StartingChannel string `json:"_starting-channel"`
}

// FilterStruct : Filter Structure
type FilterStruct struct {
	Active          bool   `json:"active"`
	CaseSensitive   bool   `json:"caseSensitive"`
	PreserveMapping bool   `json:"preserveMapping"`
	Description     string `json:"description"`
	Exclude         string `json:"exclude"`
	Filter          string `json:"filter"`
	Include         string `json:"include"`
	Name            string `json:"name"`
	Rule            string `json:"rule,omitempty"`
	Type            string `json:"type"`
	StartingChannel string `json:"startingChannel"`
}

// StreamingURLS : Information on all Streaming URL's
type StreamingURLS struct {
	Streams map[string]StreamInfo `json:"channels"`
}

// StreamInfo : Information about the Channel for the Streaming URL
type StreamInfo struct {
	ChannelNumber string `json:"channelNumber"`
	Name          string `json:"name"`
	PlaylistID    string `json:"playlistID"`
	URL           string `json:"url"`
	URLid         string `json:"urlID"`
}

// Notification : Notifications in the Web Interface
type Notification struct {
	Headline string `json:"headline"`
	Message  string `json:"message"`
	New      bool   `json:"new"`
	Time     string `json:"time"`
	Type     string `json:"type"`
}

// SettingsStruct : Content of settings.json
type SettingsStruct struct {
	API                   bool     `json:"api"`
	AuthenticationAPI     bool     `json:"authentication.api"`
	AuthenticationM3U     bool     `json:"authentication.m3u"`
	AuthenticationPMS     bool     `json:"authentication.pms"`
	AuthenticationWEB     bool     `json:"authentication.web"`
	AuthenticationXML     bool     `json:"authentication.xml"`
	BackupKeep            int      `json:"backup.keep"`
	BackupPath            string   `json:"backup.path"`
	Branch                string   `json:"git.branch,omitempty"`
	Buffer                string   `json:"buffer"`
	BufferSize            int      `json:"buffer.size.kb"`
	BufferTimeout         float64  `json:"buffer.timeout"`
	CacheImages           bool     `json:"cache.images"`
	ClearXMLTVCache       bool     `json:"clearXMLTVCache"`
	DefaultMissingEPG     string   `json:"defaultMissingEPG"`
	DisallowURLDuplicates bool     `json:"disallowURLDuplicates"`
	EnableMappedChannels  bool     `json:"enableMappedChannels"`
	EpgSource             string   `json:"epgSource"`
	FFmpegOptions         string   `json:"ffmpeg.options"`
	FFmpegPath            string   `json:"ffmpeg.path"`
	VLCOptions            string   `json:"vlc.options"`
	VLCPath               string   `json:"vlc.path"`
	FileM3U               []string `json:"file,omitempty"`  // In the Wizard, the M3U is saved in a Slice
	FileXMLTV             []string `json:"xmltv,omitempty"` // Old Storage System of the provider XML File Slice (Required for the conversion to the new one)

	Files struct {
		HDHR  map[string]interface{} `json:"hdhr"`
		M3U   map[string]interface{} `json:"m3u"`
		XMLTV map[string]interface{} `json:"xmltv"`
	} `json:"files"`

	FilesUpdate               bool                  `json:"files.update"`
	Filter                    map[int64]interface{} `json:"filter"`
	HostIP                    string                `json:"hostIP"`   // IP chosen in web client. Used to form m3u and xml files.
	HostName                  string                `json:"hostName"` // Hostname chosen in web client. Used to form m3u and xml files.
	Key                       string                `json:"key,omitempty"`
	Language                  string                `json:"language"`
	LogEntriesRAM             int                   `json:"log.entries.ram"`
	M3U8AdaptiveBandwidthMBPS int                   `json:"m3u8.adaptive.bandwidth.mbps"`
	MappingFirstChannel       float64               `json:"mapping.first.channel"`
	Port                      string                `json:"port"`
	SSDP                      bool                  `json:"ssdp"`
	StoreBufferInRAM          bool                  `json:"storeBufferInRAM"`
	TempPath                  string                `json:"temp.path"`
	TLSMode                   bool                  `json:"tlsMode"`
	Tuner                     int                   `json:"tuner"`
	Update                    []string              `json:"update"`
	UpdateURL                 string                `json:"update.url,omitempty"`
	UserAgent                 string                `json:"user.agent"`
	UUID                      string                `json:"uuid"`
	UDPxy                     string                `json:"udpxy"`
	Version                   string                `json:"version"`
	XepgReplaceMissingImages  bool                  `json:"xepg.replace.missing.images"`
	XteveAutoUpdate           bool                  `json:"xteveAutoUpdate"`
}

// LanguageUI : Language for the WebUI
type LanguageUI struct {
	Login struct {
		Failed string
	}
}
