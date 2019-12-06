package src

// SystemStruct : Beinhaltet alle Systeminformationen
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
	DVRLimit            int

	FFmpeg struct {
		DefaultOptions string
		Path           string
	}

	VLC struct {
		DefaultOptions string
		Path           string
	}

	File struct {
		Authentication string
		M3U            string
		PMS            string
		Settings       string
		URLS           string
		XEPG           string
		XML            string
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
		Config       string
		Data         string
		ImagesCache  string
		ImagesUpload string
		Temp         string
	}

	Hostname               string
	ImageCachingInProgress int
	IPAddress              string
	IPAddressesList        []string
	IPAddressesV4          []string
	IPAddressesV6          []string
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
	Version string
	WEB     struct {
		Menu []string
	}
}

// GitStruct : Updateinformationen von GitHub
type GitStruct struct {
	Filename string `json:"filename"`
	Version  string `json:"version"`
}

// DataStruct : Alle Daten werden hier abgelegt. (Lineup, XMLTV)
type DataStruct struct {
	Cache struct {
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

// Filter : Wird für die Filterregeln verwendet
type Filter struct {
	CaseSensitive bool
	Rule          string
	Type          string
}

// XEPGChannelStruct : XEPG Struktur
type XEPGChannelStruct struct {
	FileM3UID          string `json:"_file.m3u.id,required"`
	FileM3UName        string `json:"_file.m3u.name,required"`
	FileM3UPath        string `json:"_file.m3u.path,required"`
	GroupTitle         string `json:"group-title,required"`
	Name               string `json:"name,required"`
	TvgID              string `json:"tvg-id,required"`
	TvgLogo            string `json:"tvg-logo,required"`
	TvgName            string `json:"tvg-name,required"`
	URL                string `json:"url,required"`
	UUIDKey            string `json:"_uuid.key,required"`
	UUIDValue          string `json:"_uuid.value,omitempty"`
	Values             string `json:"_values,required"`
	XActive            bool   `json:"x-active,required"`
	XCategory          string `json:"x-category,required"`
	XChannelID         string `json:"x-channelID,required"`
	XEPG               string `json:"x-epg,required"`
	XGroupTitle        string `json:"x-group-title,required"`
	XMapping           string `json:"x-mapping,required"`
	XmltvFile          string `json:"x-xmltv-file,required"`
	XName              string `json:"x-name,required"`
	XUpdateChannelIcon bool   `json:"x-update-channel-icon,required"`
	XUpdateChannelName bool   `json:"x-update-channel-name,required"`
}

// M3UChannelStructXEPG : M3U Struktur für XEPG
type M3UChannelStructXEPG struct {
	FileM3UID   string `json:"_file.m3u.id,required"`
	FileM3UName string `json:"_file.m3u.name,required"`
	FileM3UPath string `json:"_file.m3u.path,required"`
	GroupTitle  string `json:"group-title,required"`
	Name        string `json:"name,required"`
	TvgID       string `json:"tvg-id,required"`
	TvgLogo     string `json:"tvg-logo,required"`
	TvgName     string `json:"tvg-name,required"`
	URL         string `json:"url,required"`
	UUIDKey     string `json:"_uuid.key,required"`
	UUIDValue   string `json:"_uuid.value,required"`
	Values      string `json:"_values,required"`
}

// FilterStruct : Filter Struktur
type FilterStruct struct {
	Active        bool   `json:"active,required"`
	CaseSensitive bool   `json:"caseSensitive,required"`
	Description   string `json:"description,required"`
	Exclude       string `json:"exclude,required"`
	Filter        string `json:"filter,required"`
	Include       string `json:"include,required"`
	Name          string `json:"name,required"`
	Rule          string `json:"rule,omitempty"`
	Type          string `json:"type,required"`
}

// StreamingURLS : Informationen zu allen streaming URL's
type StreamingURLS struct {
	Streams map[string]StreamInfo `json:"channels,required"`
}

// StreamInfo : Informationen zum Kanal für die streaming URL
type StreamInfo struct {
	ChannelNumber string `json:"channelNumber,required"`
	Name          string `json:"name,required"`
	PlaylistID    string `json:"playlistID,required"`
	URL           string `json:"url,required"`
	URLid         string `json:"urlID,required"`
}

// Notification : Notifikationen im Webinterface
type Notification struct {
	Headline string `json:"headline,required"`
	Message  string `json:"message,required"`
	New      bool   `json:"new,required"`
	Time     string `json:"time,required"`
	Type     string `json:"type,required"`
}

// SettingsStrcut : Inhalt der settings.json
type SettingsStrcut struct {
	API               bool     `json:"api"`
	AuthenticationAPI bool     `json:"authentication.api"`
	AuthenticationM3U bool     `json:"authentication.m3u"`
	AuthenticationPMS bool     `json:"authentication.pms"`
	AuthenticationWEB bool     `json:"authentication.web"`
	AuthenticationXML bool     `json:"authentication.xml"`
	BackupKeep        int      `json:"backup.keep"`
	BackupPath        string   `json:"backup.path"`
	Branch            string   `json:"git.branch,omitempty"`
	Buffer            string   `json:"buffer"`
	BufferSize        int      `json:"buffer.size.kb"`
	BufferTimeout     float64  `json:"buffer.timeout"`
	CacheImages       bool     `json:"cache.images"`
	EpgSource         string   `json:"epgSource"`
	FFmpegOptions     string   `json:"ffmpeg.options"`
	FFmpegPath        string   `json:"ffmpeg.path"`
	VLCOptions        string   `json:"vlc.options"`
	VLCPath           string   `json:"vlc.path"`
	FileM3U           []string `json:"file,omitempty"`  // Beim Wizard wird die M3U in ein Slice gespeichert
	FileXMLTV         []string `json:"xmltv,omitempty"` // Altes Speichersystem der Provider XML Datei Slice (Wird für die Umwandlung auf das neue benötigt)

	Files struct {
		HDHR  map[string]interface{} `json:"hdhr"`
		M3U   map[string]interface{} `json:"m3u"`
		XMLTV map[string]interface{} `json:"xmltv"`
	} `json:"files"`

	FilesUpdate               bool                  `json:"files.update"`
	Filter                    map[int64]interface{} `json:"filter"`
	Key                       string                `json:"key,omitempty"`
	Language                  string                `json:"language"`
	LogEntriesRAM             int                   `json:"log.entries.ram"`
	M3U8AdaptiveBandwidthMBPS int                   `json:"m3u8.adaptive.bandwidth.mbps"`
	MappingFirstChannel       float64               `json:"mapping.first.channel"`
	Port                      string                `json:"port"`
	SSDP                      bool                  `json:"ssdp"`
	SchemeM3U                 string                `json:"scheme.m3u"`
	SchemeXML                 string                `json:"scheme.xml"`
	TempPath                  string                `json:"temp.path"`
	Tuner                     int                   `json:"tuner"`
	Update                    []string              `json:"update"`
	UpdateURL                 string                `json:"update.url,omitempty"`
	UserAgent                 string                `json:"user.agent"`
	UUID                      string                `json:"uuid"`
	Version                   string                `json:"version"`
	XepgReplaceMissingImages  bool                  `json:"xepg.replace.missing.images"`
	XteveAutoUpdate           bool                  `json:"xteveAutoUpdate"`
}

// LanguageUI : Sprache für das WebUI
type LanguageUI struct {
	Login struct {
		Failed string
	}
}
