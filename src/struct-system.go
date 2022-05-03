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
	CaseSensitive     bool
	PreserveMapping   bool
	Rule              string
	Type              string
	StartingChannel   string
	DefaultMissingEPG string
}

// XEPGChannelStruct : XEPG Structure
type XEPGChannelStruct struct {
	FileM3UID                     string `json:"_file.m3u.id,required"`
	FileM3UName                   string `json:"_file.m3u.name,required"`
	FileM3UPath                   string `json:"_file.m3u.path,required"`
	GroupTitle                    string `json:"group-title,required"`
	Name                          string `json:"name,required"`
	TvgID                         string `json:"tvg-id,required"`
	TvgLogo                       string `json:"tvg-logo,required"`
	TvgName                       string `json:"tvg-name,required"`
	TvgShift                      string `json:"tvg-shift,required"`
	UpdateChannelNameRegex        string `json:"update-channel-name-regex,required"`
	UpdateChannelNameByGroupRegex string `json:"update-channel-name-by-group-regex,required"`
	URL                           string `json:"url,required"`
	UUIDKey                       string `json:"_uuid.key,required"`
	UUIDValue                     string `json:"_uuid.value,omitempty"`
	Values                        string `json:"_values,required"`
	XActive                       bool   `json:"x-active,required"`
	XCategory                     string `json:"x-category,required"`
	XChannelID                    string `json:"x-channelID,required"`
	XEPG                          string `json:"x-epg,required"`
	XGroupTitle                   string `json:"x-group-title,required"`
	XMapping                      string `json:"x-mapping,required"`
	XmltvFile                     string `json:"x-xmltv-file,required"`
	XName                         string `json:"x-name,required"`
	XUpdateChannelIcon            bool   `json:"x-update-channel-icon,required"`
	XUpdateChannelName            bool   `json:"x-update-channel-name,required"`
	XUpdateChannelGroup           bool   `json:"x-update-channel-group,required"`
	XDescription                  string `json:"x-description,required"`
	XTimeshift                    string `json:"x-timeshift,required"`
	DefaultMissingEPG             string `json:"x-default-missing-epg,required"`
}

// M3UChannelStructXEPG : M3U Structure for XEPG
type M3UChannelStructXEPG struct {
	FileM3UID         string `json:"_file.m3u.id,required"`
	FileM3UName       string `json:"_file.m3u.name,required"`
	FileM3UPath       string `json:"_file.m3u.path,required"`
	GroupTitle        string `json:"group-title,required"`
	Name              string `json:"name,required"`
	TvgID             string `json:"tvg-id,required"`
	TvgLogo           string `json:"tvg-logo,required"`
	TvgName           string `json:"tvg-name,required"`
	TvgShift          string `json:"tvg-shift,required"`
	URL               string `json:"url,required"`
	UUIDKey           string `json:"_uuid.key,required"`
	UUIDValue         string `json:"_uuid.value,required"`
	Values            string `json:"_values,required"`
	PreserveMapping   string `json:"_preserve-mapping,required"`
	StartingChannel   string `json:"_starting-channel,required"`
	DefaultMissingEPG string `json:"_default-missing-epg,required"`
}

// FilterStruct : Filter Structure
type FilterStruct struct {
	Active            bool   `json:"active,required"`
	CaseSensitive     bool   `json:"caseSensitive,required"`
	PreserveMapping   bool   `json:"preserveMapping,required"`
	Description       string `json:"description,required"`
	Exclude           string `json:"exclude,required"`
	Filter            string `json:"filter,required"`
	Include           string `json:"include,required"`
	Name              string `json:"name,required"`
	Rule              string `json:"rule,omitempty"`
	Type              string `json:"type,required"`
	StartingChannel   string `json:"startingChannel,required"`
	DefaultMissingEPG string `json:"defaultMissingEPG,required"`
}

// StreamingURLS : Information on all Streaming URL's
type StreamingURLS struct {
	Streams map[string]StreamInfo `json:"channels,required"`
}

// StreamInfo : Information about the Channel for the Streaming URL
type StreamInfo struct {
	ChannelNumber string `json:"channelNumber,required"`
	Name          string `json:"name,required"`
	PlaylistID    string `json:"playlistID,required"`
	URL           string `json:"url,required"`
	URLid         string `json:"urlID,required"`
}

// Notification : Notifications in the Web Interface
type Notification struct {
	Headline string `json:"headline,required"`
	Message  string `json:"message,required"`
	New      bool   `json:"new,required"`
	Time     string `json:"time,required"`
	Type     string `json:"type,required"`
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
	DisallowURLDuplicates bool     `json:"disallowURLDuplicates"`
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
	HostIP                    string                `json:"hostIP"` // IP chosen in web client. Used to form m3u and xml files.
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
