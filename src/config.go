package src

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

// System : Beinhaltet alle Systeminformationen
var System SystemStruct

// WebScreenLog : Logs werden im RAM gespeichert und für das Webinterface bereitgestellt
var WebScreenLog WebScreenLogStruct

// Settings : Inhalt der settings.json
var Settings SettingsStrcut

// Data : Alle Daten werden hier abgelegt. (Lineup, XMLTV)
var Data DataStruct

// SystemFiles : Alle Systemdateien
var SystemFiles = []string{"authentication.json", "pms.json", "settings.json", "xepg.json", "urls.json"}

// BufferInformation : Informationen über den Buffer (aktive Streams, maximale Streams)
var BufferInformation sync.Map

// BufferClients : Anzahl der Clients die einen Stream über den Buffer abspielen
var BufferClients sync.Map

// Init : Systeminitialisierung
func Init() (err error) {

	var debug string

	// System Einstellungen
	System.AppName = strings.ToLower(System.Name)
	System.ARCH = runtime.GOARCH
	System.OS = runtime.GOOS
	System.ServerProtocol.API = "http"
	System.ServerProtocol.DVR = "http"
	System.ServerProtocol.M3U = "http"
	System.ServerProtocol.WEB = "http"
	System.ServerProtocol.XML = "http"
	System.DVRLimit = 480
	System.Compatibility = "1.4.4"

	// Default Logeinträge, wird später von denen aus der settings.json überschrieben. Muss gemacht werden, damit die ersten Einträge auch im Log (webUI aangezeigt werden)
	Settings.LogEntriesRAM = 500

	// Variablen für den Update Prozess
	//System.Update.Git = "https://github.com/xteve-project/xTeVe-Downloads/blob"
	System.Update.Git = fmt.Sprintf("https://github.com/%s/%s/blob", System.GitHub.User, System.GitHub.Repo)
	System.Update.Name = "xteve_2"

	// Ordnerpfade festlegen
	var tempFolder = os.TempDir() + string(os.PathSeparator) + System.AppName + string(os.PathSeparator)
	tempFolder = getPlatformPath(strings.Replace(tempFolder, "//", "/", -1))

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
	System.Folder.Temp = tempFolder

	// Dev Info
	showDevInfo()

	// System Ordner erstellen
	err = createSystemFolders()
	if err != nil {
		ShowError(err, 1070)
		return
	}

	if len(System.Flag.Restore) > 0 {
		// Einstellungen werden über CLI wiederhergestellt. Weitere Initialisierung ist nicht notwendig.
		return
	}

	System.File.XML = getPlatformFile(fmt.Sprintf("%s%s.xml", System.Folder.Data, System.AppName))
	System.File.M3U = getPlatformFile(fmt.Sprintf("%s%s.m3u", System.Folder.Data, System.AppName))

	err = activatedSystemAuthentication()
	if err != nil {
		return
	}

	err = resolveHostIP()
	if err != nil {
		ShowError(err, 1002)
	}

	// Menü für das Webinterface
	System.WEB.Menu = []string{"playlist", "filter", "xmltv", "mapping", "users", "settings", "log", "logout"}

	fmt.Println("For help run: " + getPlatformFile(os.Args[0]) + " -h")
	fmt.Println()

	// Überprüfen ob xTeVe als root läuft
	if os.Geteuid() == 0 {
		showWarning(2010)
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

	// Systemdateien erstellen (Falls nicht vorhanden)
	err = createSystemFiles()
	if err != nil {
		ShowError(err, 1071)
		return
	}

	// Bedingte Update Änderungen durchführen
	err = conditionalUpdateChanges()
	if err != nil {
		ShowError(err, 0)
		return
	}

	// Einstellungen laden (settings.json)
	showInfo(fmt.Sprintf("Load Settings:%s", System.File.Settings))

	_, err = loadSettings()
	if err != nil {
		ShowError(err, 0)
		return
	}

	// Berechtigung aller Ordner überprüfen
	err = checkFilePermission(System.Folder.Config)
	if err == nil {
		err = checkFilePermission(System.Folder.Temp)
	}

	// Separaten tmp Ordner für jede Instanz
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

	// Branch festlegen
	System.Branch = Settings.Branch

	if System.Dev == true {
		System.Branch = "Development"
	}

	if len(System.Branch) == 0 {
		System.Branch = "master"
	}

	showInfo(fmt.Sprintf("GitHub:https://github.com/%s", System.GitHub.User))
	showInfo(fmt.Sprintf("Git Branch:%s [%s]", System.Branch, System.GitHub.User))

	// Domainnamen setzten
	setGlobalDomain(fmt.Sprintf("%s:%s", System.IPAddress, Settings.Port))

	System.URLBase = fmt.Sprintf("%s://%s:%s", System.ServerProtocol.WEB, System.IPAddress, Settings.Port)

	// HTML Dateien erstellen, mit dev == true werden die lokalen HTML Dateien verwendet
	if System.Dev == true {

		HTMLInit("webUI", "src", "html"+string(os.PathSeparator), "src"+string(os.PathSeparator)+"webUI.go")
		err = BuildGoFile()
		if err != nil {
			return
		}

	}

	// DLNA Server starten
	err = SSDP()
	if err != nil {
		return
	}

	// HTML Datein laden
	loadHTMLMap()

	return
}

// StartSystem : System wird gestartet
func StartSystem(updateProviderFiles bool) (err error) {

	setDeviceID()

	if System.ScanInProgress == 1 {
		return
	}

	// Systeminformationen in der Konsole ausgeben
	showInfo(fmt.Sprintf("UUID:%s", Settings.UUID))
	showInfo(fmt.Sprintf("Tuner (Plex / Emby):%d", Settings.Tuner))
	showInfo(fmt.Sprintf("EPG Source:%s", Settings.EpgSource))
	showInfo(fmt.Sprintf("Plex Channel Limit:%d", System.DVRLimit))

	// Providerdaten aktualisieren
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
