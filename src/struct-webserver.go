package src

// RequestStruct : Anfragen über die Websocket Schnittstelle
type RequestStruct struct {
	// Befehle an xTeVe
	Cmd string `json:"cmd,required"`

	// Benutzer
	DeleteUser bool                   `json:"deleteUser,omitempty"`
	UserData   map[string]interface{} `json:"userData,omitempty"`

	// Mapping
	EpgMapping map[string]interface{} `json:"epgMapping,omitempty"`

	// Restore
	Base64 string `json:"base64,omitempty"`

	// Neue Werte für die Einstellungen (settings.json)
	Settings struct {
		API                      *bool     `json:"api,omitempty"`
		AuthenticationAPI        *bool     `json:"authentication.api,omitempty"`
		AuthenticationM3U        *bool     `json:"authentication.m3u,omitempty"`
		AuthenticationPMS        *bool     `json:"authentication.pms,omitempty"`
		AuthenticationWEP        *bool     `json:"authentication.web,omitempty"`
		AuthenticationXML        *bool     `json:"authentication.xml,omitempty"`
		BackupKeep               *int      `json:"backup.keep,omitempty"`
		BackupPath               *string   `json:"backup.path,omitempty"`
		Buffer                   *string   `json:"buffer,omitempty"`
		BufferSize               *int      `json:"buffer.size.kb, omitempty"`
		BufferTimeout            *float64  `json:"buffer.timeout,omitempty"`
		CacheImages              *bool     `json:"cache.images,omitempty"`
		EpgSource                *string   `json:"epgSource,omitempty"`
		FFmpegOptions            *string   `json:"ffmpeg.options,omitempty"`
		FFmpegPath               *string   `json:"ffmpeg.path,omitempty"`
		VLCOptions               *string   `json:"vlc.options,omitempty"`
		VLCPath                  *string   `json:"vlc.path,omitempty"`
		FilesUpdate              *bool     `json:"files.update,omitempty"`
		TempPath                 *string   `json:"temp.path,omitempty"`
		Tuner                    *int      `json:"tuner,omitempty"`
		Update                   *[]string `json:"update,omitempty"`
		UserAgent                *string   `json:"user.agent,omitempty"`
		XepgReplaceMissingImages *bool     `json:"xepg.replace.missing.images,omitempty"`
		XteveAutoUpdate          *bool     `json:"xteveAutoUpdate,omitempty"`
		SchemeM3U                *string   `json:"scheme.m3u,omitempty"`
		SchemeXML                *string   `json:"scheme.xml,omitempty"`
	} `json:"settings,omitempty"`

	// Upload Logo
	Filename string `json:"filename,omitempty"`

	// Filter
	Filter map[int64]interface{} `json:"filter,omitempty"`

	// Dateien (M3U, HDHR, XMLTV)
	Files struct {
		HDHR  map[string]interface{} `json:"hdhr,omitempty"`
		M3U   map[string]interface{} `json:"m3u,omitempty"`
		XMLTV map[string]interface{} `json:"xmltv,omitempty"`
	} `json:"files,omitempty"`

	// Wizard
	Wizard struct {
		EpgSource *string `json:"epgSource,omitempty"`
		M3U       *string `json:"m3u,omitempty"`
		Tuner     *int    `json:"tuner,omitempty"`
		XMLTV     *string `json:"xmltv,omitempty"`
	} `json:"wizard,omitempty"`
}

// ResponseStruct : Antworten an den Client (WEB)
type ResponseStruct struct {
	ClientInfo struct {
		ARCH      string `json:"arch"`
		Branch    string `json:"branch,omitempty"`
		DVR       string `json:"DVR"`
		EpgSource string `json:"epgSource"`
		Errors    int    `json:"errors"`
		M3U       string `json:"m3u-url,required"`
		OS        string `json:"os"`
		Streams   string `json:"streams"`
		UUID      string `json:"uuid"`
		Version   string `json:"version"`
		Warnings  int    `json:"warnings"`
		XEPGCount int64  `json:"xepg"`
		XML       string `json:"xepg-url,required"`
	} `json:"clientInfo,omitempty"`

	Data struct {
		Playlist struct {
			M3U struct {
				Groups struct {
					Text  []string `json:"text,required"`
					Value []string `json:"value,required"`
				} `json:"groups,required"`
			} `json:"m3u,required"`
		} `json:"playlist,required"`

		StreamPreviewUI struct {
			Active   []string `json:"activeStreams,required"`
			Inactive []string `json:"inactiveStreams,required"`
		}
	} `json:"data,required"`

	Alert               string                 `json:"alert,omitempty"`
	ConfigurationWizard bool                   `json:"configurationWizard,required"`
	Error               string                 `json:"err,omitempty"`
	Log                 WebScreenLogStruct     `json:"log,required"`
	LogoURL             string                 `json:"logoURL,omitempty"`
	OpenLink            string                 `json:"openLink,omitempty"`
	OpenMenu            string                 `json:"openMenu,omitempty"`
	Reload              bool                   `json:"reload,omitempty"`
	Settings            SettingsStrcut         `json:"settings,required"`
	Status              bool                   `json:"status,required"`
	Token               string                 `json:"token,omitempty"`
	Users               map[string]interface{} `json:"users,omitempty"`
	Wizard              int                    `json:"wizard,omitempty"`
	XEPG                map[string]interface{} `json:"xepg,required"`

	Notification map[string]Notification `json:"notification,omitempty"`
}

// APIRequestStruct : Anfrage über die API Schnittstelle
type APIRequestStruct struct {
	Cmd      string `json:"cmd"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Username string `json:"username"`
}

// APIResponseStruct : Antwort an den Client (API)
type APIResponseStruct struct {
	EpgSource     string `json:"epg.source,omitempty"`
	Error         string `json:"err,omitempty"`
	Status        bool   `json:"status,required"`
	StreamsActive int64  `json:"streams.active,omitempty"`
	StreamsAll    int64  `json:"streams.all,omitempty"`
	StreamsXepg   int64  `json:"streams.xepg,omitempty"`
	Token         string `json:"token,omitempty"`
	URLDvr        string `json:"url.dvr,omitempty"`
	URLM3U        string `json:"url.m3u,omitempty"`
	URLXepg       string `json:"url.xepg,omitempty"`
	VersionAPI    string `json:"version.api,omitempty"`
	VersionXteve  string `json:"version.xteve,omitempty"`
}

// WebScreenLogStruct : Logs werden im RAM gespeichert und für das Webinterface bereitgestellt
type WebScreenLogStruct struct {
	Errors   int      `json:"errors,required"`
	Log      []string `json:"log,required"`
	Warnings int      `json:"warnings,required"`
}
