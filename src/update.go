package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	up2date "../src/internal/up2date/client"

	"reflect"
)

// BinaryUpdate : Binary Update Prozess. Git Branch master und beta wird von GitHub geladen.
func BinaryUpdate() (err error) {

	if System.GitHub.Update == false {
		showWarning(2099)
		return
	}

	var debug string

	var updater = &up2date.Updater
	updater.Name = System.Update.Name
	updater.Branch = System.Branch

	up2date.Init()

	switch System.Branch {

	// Update von GitHub
	case "master", "beta":

		var gitInfo = fmt.Sprintf("%s/%s/info.json?raw=true", System.Update.Git, System.Branch)
		var zipFile = fmt.Sprintf("%s/%s/%s_%s_%s.zip?raw=true", System.Update.Git, System.Branch, System.AppName, System.OS, System.ARCH)
		var body []byte

		var git GitStruct

		resp, err := http.Get(gitInfo)
		if err != nil {
			ShowError(err, 0)
			return err
		}

		if resp.StatusCode != http.StatusOK {

			if resp.StatusCode == 404 {
				err = fmt.Errorf(fmt.Sprintf("Update Server: %s (%s)", http.StatusText(resp.StatusCode), gitInfo))
				ShowError(err, 6003)
				return nil
			}

			err = fmt.Errorf(fmt.Sprintf("%d: %s (%s)", resp.StatusCode, http.StatusText(resp.StatusCode), gitInfo))

			return err
		}

		body, err = ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &git)
		if err != nil {
			return err
		}

		updater.Response.Status = true
		updater.Response.UpdateZIP = zipFile
		updater.Response.Version = git.Version
		updater.Response.Filename = git.Filename

	// Update vom eigenen Server
	default:

		updater.URL = Settings.UpdateURL

		if len(updater.URL) == 0 {
			showInfo(fmt.Sprintf("Update URL:No server URL specified, update will not be performed. Branch: %s", System.Branch))
			return
		}

		showInfo("Update URL:" + updater.URL)
		fmt.Println("-----------------")

		// Versionsinformationen vom Server laden
		err = up2date.GetVersion()
		if err != nil {

			debug = fmt.Sprintf(err.Error())
			showDebug(debug, 1)

			return nil
		}

		if len(updater.Response.Reason) > 0 {

			err = fmt.Errorf(fmt.Sprintf("Update Server: %s", updater.Response.Reason))
			ShowError(err, 6002)

			return nil
		}

	}

	var currentVersion = System.Version + "." + System.Build

	// Versionsnummer überprüfen
	if updater.Response.Version > currentVersion && updater.Response.Status == true {

		if Settings.XteveAutoUpdate == true {
			// Update durchführen
			var fileType, url string

			showInfo(fmt.Sprintf("Update Available:Version: %s", updater.Response.Version))

			switch System.Branch {

			// Update von GitHub
			case "master", "beta":
				showInfo(fmt.Sprintf("Update Server:GitHub"))

			// Update vom eigenen Server
			default:
				showInfo(fmt.Sprintf("Update Server:%s", Settings.UpdateURL))

			}

			showInfo(fmt.Sprintf("Start Update:Branch: %s", updater.Branch))

			// Neue Version als BIN Datei herunterladen
			if len(updater.Response.UpdateBIN) > 0 {
				url = updater.Response.UpdateBIN
				fileType = "bin"
			}

			// Neue Version als ZIP Datei herunterladen
			if len(updater.Response.UpdateZIP) > 0 {
				url = updater.Response.UpdateZIP
				fileType = "zip"
			}

			if len(url) > 0 {

				err = up2date.DoUpdate(fileType, updater.Response.Filename)
				if err != nil {
					ShowError(err, 6002)
				}

			}

		} else {
			// Hinweis ausgeben
			showWarning(6004)
		}

	}

	return nil
}

func conditionalUpdateChanges() (err error) {

checkVersion:
	settingsMap, err := loadJSONFileToMap(System.File.Settings)
	if err != nil || len(settingsMap) == 0 {
		return
	}

	if settingsVersion, ok := settingsMap["version"].(string); ok {

		// Letzte Kompatible Version (1.4.4)
		if settingsVersion < System.Compatibility {
			err = errors.New(getErrMsg(1013))
			return
		}

		switch settingsVersion {

		case "1.4.4":
			// UUID Wert in xepg.json setzen
			err = setValueForUUID()
			if err != nil {
				return
			}

			// Neuer Filter (WebUI). Alte Filtereinstellungen werden konvertiert
			if oldFilter, ok := settingsMap["filter"].([]interface{}); ok {
				var newFilterMap = convertToNewFilter(oldFilter)
				settingsMap["filter"] = newFilterMap

				settingsMap["version"] = "1.9.0"

				err = saveMapToJSONFile(System.File.Settings, settingsMap)
				if err != nil {
					return
				}

				goto checkVersion

			} else {
				err = errors.New(getErrMsg(1030))
				return
			}

		case "1.9.0":
			// Falls es in einem späteren Update Änderungen an der Datenbank gibt, geht es hier weiter

			break

		}

	} else {
		// settings.json ist zu alt (älter als Version 1.4.4)
		err = errors.New(getErrMsg(1013))
	}

	return
}

func convertToNewFilter(oldFilter []interface{}) (newFilterMap map[int]interface{}) {

	newFilterMap = make(map[int]interface{})

	switch reflect.TypeOf(oldFilter).Kind() {

	case reflect.Slice:
		s := reflect.ValueOf(oldFilter)

		for i := 0; i < s.Len(); i++ {

			var newFilter FilterStruct
			newFilter.Active = true
			newFilter.Name = fmt.Sprintf("Custom filter %d", i+1)
			newFilter.Filter = s.Index(i).Interface().(string)
			newFilter.Type = "custom-filter"
			newFilter.CaseSensitive = false

			newFilterMap[i] = newFilter

		}

	}

	return
}

func setValueForUUID() (err error) {

	xepg, err := loadJSONFileToMap(System.File.XEPG)

	for _, c := range xepg {

		var xepgChannel = c.(map[string]interface{})

		if uuidKey, ok := xepgChannel["_uuid.key"].(string); ok {

			if value, ok := xepgChannel[uuidKey].(string); ok {

				if len(value) > 0 {
					xepgChannel["_uuid.value"] = value
				}

			}

		}

	}

	err = saveMapToJSONFile(System.File.XEPG, xepg)

	return
}
