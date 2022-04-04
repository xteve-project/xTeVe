package src

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

// System : Contains all System Information
var System SystemStruct

// WebScreenLog : Logs are saved in RAM and made available for the Web interface
var WebScreenLog WebScreenLogStruct

// Settings : Content of settings.json
var Settings SettingsStruct

// Data : All data is stored here. (Lineup, XMLTV)
var Data DataStruct

// SystemFiles : All System Files
var SystemFiles = []string{"authentication.json", "pms.json", "settings.json", "xepg.json", "urls.json"}

// BufferInformation : Information about the Buffer (active Streams, maximum Streams)
var BufferInformation sync.Map

// BufferClients : Number of Clients playing a Stream over the Buffer
var BufferClients sync.Map

// Lock : Lock Map
var Lock = sync.RWMutex{}

// Init : System Initialization
func Init() (err error) {

	var debug string

	// System Settings
	System.AppName = strings.ToLower(System.Name)
	System.ARCH = runtime.GOARCH
	System.OS = runtime.GOOS
	System.ServerProtocol.API = "http"
	System.ServerProtocol.DVR = "http"
	System.ServerProtocol.M3U = "http"
	System.ServerProtocol.WEB = "http"
	System.ServerProtocol.XML = "http"
	System.PlexChannelLimit = 480
	System.UnfilteredChannelLimit = 480
	System.Compatibility = "1.4.4"

	// FFmpeg Default Settings
	System.FFmpeg.DefaultOptions = "-hide_banner -loglevel error -i [URL] -c copy -f mpegts pipe:1"
	System.VLC.DefaultOptions = "-I dummy [URL] --sout #std{mux=ts,access=file,dst=-}"

	// Default Log Entries, which will later be overwritten by those from settings.json. Needed so that the first entries are also displayed in the Log (webUI are displayed)
	Settings.LogEntriesRAM = 500

	// Variables for the Update Process
	//System.Update.Git = "https://github.com/xteve-project/xTeVe-Downloads/blob"
	System.Update.Git = fmt.Sprintf("https://github.com/%s/%s/blob", System.GitHub.User, System.GitHub.Repo)
	System.Update.Name = "xteve_2"

	// Define folder paths
	if len(System.Folder.Config) == 0 {
		System.Folder.Config = GetUserHomeDirectory() + string(os.PathSeparator) + "." + System.AppName + string(os.PathSeparator)
	} else {
		System.Folder.Config = strings.TrimRight(System.Folder.Config, string(os.PathSeparator)) + string(os.PathSeparator)
	}

	System.Folder.Config = getPlatformPath(System.Folder.Config)

	System.Folder.Backup = System.Folder.Config + "backup" + string(os.PathSeparator)
	System.Folder.Data = System.Folder.Config + "data" + string(os.PathSeparator)
	System.Folder.Cache = System.Folder.Config + "cache" + string(os.PathSeparator)
	System.Folder.ImagesCache = System.Folder.Cache + "images" + string(os.PathSeparator)
	System.Folder.ImagesUpload = System.Folder.Data + "images" + string(os.PathSeparator)
	System.Folder.Temp = getDefaultTempDir()

	// Dev Info
	showDevInfo()

	// Create System Folder
	err = createSystemFolders()
	if err != nil {
		ShowError(err, 1070)
		return
	}

	if len(System.Flag.Restore) > 0 {
		// Settings are restored via CLI. No further Initialization is necessary.
		return
	}

	System.File.XML = getPlatformFile(fmt.Sprintf("%s%s.xml", System.Folder.Data, System.AppName))
	System.File.M3U = getPlatformFile(fmt.Sprintf("%s%s.m3u", System.Folder.Data, System.AppName))

	System.Compressed.GZxml = getPlatformFile(fmt.Sprintf("%s%s.xml.gz", System.Folder.Data, System.AppName))

	err = activatedSystemAuthentication()
	if err != nil {
		return
	}

	err = resolveHostIP()
	if err != nil {
		ShowError(err, 1002)
	}

	// Menu for the Web interface
	System.WEB.Menu = []string{"playlist", "filter", "xmltv", "mapping", "users", "settings", "log", "logout"}

	fmt.Println("For help run: " + getPlatformFile(os.Args[0]) + " -h")
	fmt.Println()

	// Check whether xTeVe is running as root
	if os.Geteuid() == 0 {
		showWarning(2110)
	}

	if System.Flag.Debug > 0 {
		debug = fmt.Sprintf("Debug Level:%d", System.Flag.Debug)
		showDebug(debug, 1)
	}

	showInfo(fmt.Sprintf("Version:%s Build: %s", System.Version, System.Build))
	showInfo(fmt.Sprintf("Database Version:%s", System.DBVersion))
	showInfo(fmt.Sprintf("System IP Addresses:IPv4: %d | IPv6: %d", len(System.IPAddressesV4), len(System.IPAddressesV6)))
	showInfo("Hostname:" + System.Hostname)
	showInfo(fmt.Sprintf("System Folder:%s", getPlatformPath(System.Folder.Config)))

	// Create System Files (If not available)
	err = createSystemFiles()
	if err != nil {
		ShowError(err, 1071)
		return
	}

	// Perform conditional Update Changes
	err = conditionalUpdateChanges()
	if err != nil {
		return
	}

	// Load Settings (settings.json)
	showInfo(fmt.Sprintf("Load Settings:%s", System.File.Settings))

	_, err = loadSettings()
	if err != nil {
		ShowError(err, 0)
		return
	}

	// Check the permissions on all Folders
	err = checkFilePermission(System.Folder.Config)
	if err == nil {
		err = checkFilePermission(System.Folder.Temp)
	}

	// Separate tmp Folder for each Instance
	//System.Folder.Temp = System.Folder.Temp + Settings.UUID + string(os.PathSeparator)
	showInfo(fmt.Sprintf("Temporary Folder:%s", getPlatformPath(System.Folder.Temp)))

	err = checkFolder(System.Folder.Temp)
	if err != nil {
		return
	}

	err = removeChildItems(getPlatformPath(System.Folder.Temp))
	if err != nil {
		return
	}

	// Set Branch
	System.Branch = Settings.Branch

	if System.Dev == true {
		System.Branch = "Development"
	}

	if len(System.Branch) == 0 {
		System.Branch = "master"
	}

	showInfo(fmt.Sprintf("GitHub:https://github.com/%s", System.GitHub.User))
	showInfo(fmt.Sprintf("Git Branch:%s [%s]", System.Branch, System.GitHub.User))

	// Set Domain Names
	setGlobalDomain(fmt.Sprintf("%s:%s", System.IPAddress, Settings.Port))

	System.URLBase = fmt.Sprintf("%s://%s:%s", System.ServerProtocol.WEB, System.IPAddress, Settings.Port)

	// Create HTML Files, with dev == true the local HTML Files are used
	if System.Dev == true {

		HTMLInit("webUI", "src", "html"+string(os.PathSeparator), "src"+string(os.PathSeparator)+"webUI.go")
		err = BuildGoFile()
		if err != nil {
			return
		}

	}

	// Start the DLNA Server
	err = SSDP()
	if err != nil {
		return
	}

	// Load HTML Files
	loadHTMLMap()

	return
}

// StartSystem : System is starting up
func StartSystem(updateProviderFiles bool) (err error) {

	setDeviceID()

	if System.ScanInProgress == 1 {
		return
	}

	// Output System Information in the Console
	showInfo(fmt.Sprintf("UUID:%s", Settings.UUID))
	showInfo(fmt.Sprintf("Tuner (Plex / Emby):%d", Settings.Tuner))
	showInfo(fmt.Sprintf("EPG Source:%s", Settings.EpgSource))
	showInfo(fmt.Sprintf("Plex Channel Limit:%d", System.PlexChannelLimit))
	showInfo(fmt.Sprintf("Unfiltered Chan. Limit:%d", System.UnfilteredChannelLimit))

	// Update Provider Data
	if len(Settings.Files.M3U) > 0 && Settings.FilesUpdate == true || updateProviderFiles == true {

		err = xTeVeAutoBackup()
		if err != nil {
			ShowError(err, 1090)
		}

		getProviderData("m3u", "")
		getProviderData("hdhr", "")

		if Settings.EpgSource == "XEPG" {
			getProviderData("xmltv", "")
		}

	}

	err = buildDatabaseDVR()
	if err != nil {
		ShowError(err, 0)
		return
	}

	buildXEPG(false)

	return
}
