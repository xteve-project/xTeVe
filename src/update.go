package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	up2date "xteve/src/internal/up2date/client"

	"reflect"
)

// BinaryUpdate : Binary update process. Git Branch master and beta is loaded from GitHub.
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

	// Update from GitHub
	case "master", "beta":

		var gitInfo = fmt.Sprintf("%s/%s/info.json?raw=true", System.Update.Git, System.Branch)
		var zipFile = fmt.Sprintf("%s/%s/%s_%s_%s.zip?raw=true", System.Update.Git, System.Branch, System.AppName, System.OS, System.ARCH)
		var body []byte

		var git GitStruct

		resp, err := http.Get(gitInfo)
		if err != nil {
			ShowError(err, 6003)
			return nil
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

	// Update from your own Server
	default:

		updater.URL = Settings.UpdateURL

		if len(updater.URL) == 0 {
			showInfo(fmt.Sprintf("Update URL:No server URL specified, update will not be performed. Branch: %s", System.Branch))
			return
		}

		showInfo("Update URL:" + updater.URL)
		fmt.Println("-----------------")

		// Load version information from the Server
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

	// Check Version Number
	if updater.Response.Version > currentVersion && updater.Response.Status == true {

		if Settings.XteveAutoUpdate == true {
			// Perform update
			var fileType, url string

			showInfo(fmt.Sprintf("Update Available:Version: %s", updater.Response.Version))

			switch System.Branch {

			// Update from GitHub
			case "master", "beta":
				showInfo(fmt.Sprintf("Update Server:GitHub"))

			// Update from your own Server
			default:
				showInfo(fmt.Sprintf("Update Server:%s", Settings.UpdateURL))

			}

			showInfo(fmt.Sprintf("Start Update:Branch: %s", updater.Branch))

			// Download the new version as a BIN File
			if len(updater.Response.UpdateBIN) > 0 {
				url = updater.Response.UpdateBIN
				fileType = "bin"
			}

			// Download the new version as a ZIP File
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
			// Display update exception
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

		if settingsVersion > System.DBVersion {
			showInfo("Settings DB Version:" + settingsVersion)
			showInfo("System DB Version:" + System.DBVersion)
			err = errors.New(getErrMsg(1031))
			return
		}

		// Latest Compatible Version (1.4.4)
		if settingsVersion < System.Compatibility {
			err = errors.New(getErrMsg(1013))
			return
		}

		switch settingsVersion {

		case "1.4.4":
			// Set UUID Value in xepg.json
			err = setValueForUUID()
			if err != nil {
				return
			}

			// New filter (WebUI). Old Filter Settings are converted
			if oldFilter, ok := settingsMap["filter"].([]interface{}); ok {
				var newFilterMap = convertToNewFilter(oldFilter)
				settingsMap["filter"] = newFilterMap

				settingsMap["version"] = "2.0.0"

				err = saveMapToJSONFile(System.File.Settings, settingsMap)
				if err != nil {
					return
				}

				goto checkVersion

			} else {
				err = errors.New(getErrMsg(1030))
				return
			}

		case "2.0.0":

			if oldBuffer, ok := settingsMap["buffer"].(bool); ok {

				var newBuffer string
				switch oldBuffer {
				case true:
					newBuffer = "xteve"
				case false:
					newBuffer = "-"
				}

				settingsMap["buffer"] = newBuffer

				settingsMap["version"] = "2.1.0"

				err = saveMapToJSONFile(System.File.Settings, settingsMap)
				if err != nil {
					return
				}

				goto checkVersion

			} else {
				err = errors.New(getErrMsg(1030))
				return
			}

		case "2.1.0", "2.1.1":
			// Database verison <= 2.1.1 has broken XEPG mapping

			// Clear XEPG mapping
			Data.XEPG.Channels = make(map[string]interface{})
			Data.XEPG.XEPGCount = 0
			Data.Cache.Streams = struct{ Active []string }{}

			err = saveMapToJSONFile(System.File.XEPG, Data.XEPG.Channels)
			if err != nil {
				ShowError(err, 000)
				return err
			}

			// Notify user
			showWarning(2022)
			sendAlert(getErrMsg(2022))

			// Update database version
			settingsMap["version"] = "2.2.0"

			err = saveMapToJSONFile(System.File.Settings, settingsMap)
			if err != nil {
				return
			}

			goto checkVersion

		case "2.2.0", "2.2.1", "2.2.2", "2.2.3", "2.3.0":
			// If there are changes to the Database in a later update, continue here

			break
		}

	} else {
		// settings.json is too old (older than Version 1.4.4)
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
