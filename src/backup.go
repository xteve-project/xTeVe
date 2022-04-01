package src

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func xTeVeAutoBackup() (err error) {

	var archiv = "xteve_auto_backup_" + time.Now().Format("20060102_1504") + ".zip"
	var target string
	var sourceFiles = make([]string, 0)
	var oldBackupFiles = make([]string, 0)
	var debug string

	if len(Settings.BackupPath) > 0 {
		System.Folder.Backup = Settings.BackupPath
	}

	showInfo("Backup Path:" + System.Folder.Backup)

	err = checkFolder(System.Folder.Backup)
	if err != nil {
		ShowError(err, 1070)
		return
	}

	// Delete Old Backups
	files, err := ioutil.ReadDir(System.Folder.Backup)

	if err == nil {

		for _, file := range files {

			if filepath.Ext(file.Name()) == ".zip" && strings.Contains(file.Name(), "xteve_auto_backup") {
				oldBackupFiles = append(oldBackupFiles, file.Name())
			}

		}

		// Delete All Backups
		var end int
		switch Settings.BackupKeep {
		case 0:
			end = 0
		default:
			end = Settings.BackupKeep - 1
		}

		for i := 0; i < len(oldBackupFiles)-end; i++ {

			os.RemoveAll(System.Folder.Backup + oldBackupFiles[i])
			debug = fmt.Sprintf("Delete backup file:%s", oldBackupFiles[i])
			showDebug(debug, 1)

		}

		if Settings.BackupKeep == 0 {
			return
		}

	} else {

		return

	}

	// Create a Backup
	if err == nil {

		target = System.Folder.Backup + archiv

		for _, i := range SystemFiles {
			sourceFiles = append(sourceFiles, System.Folder.Config+i)
		}

		sourceFiles = append(sourceFiles, System.Folder.ImagesUpload)

		err = zipFiles(sourceFiles, target)

		if err == nil {

			debug = fmt.Sprintf("Create backup file:%s", target)
			showDebug(debug, 1)

			showInfo("Backup file:" + target)

		}

	}

	return
}

func xteveBackup() (archiv string, err error) {

	err = checkFolder(System.Folder.Temp)
	if err != nil {
		return
	}

	archiv = "xteve_backup_" + time.Now().Format("20060102_1504") + ".zip"

	var target = System.Folder.Temp + archiv
	var sourceFiles = make([]string, 0)

	for _, i := range SystemFiles {
		sourceFiles = append(sourceFiles, System.Folder.Config+i)
	}

	sourceFiles = append(sourceFiles, System.Folder.Data)

	err = zipFiles(sourceFiles, target)
	if err != nil {
		ShowError(err, 0)
		return
	}

	return
}

func xteveRestore(archive string) (newWebURL string, err error) {

	var newPort, oldPort, backupVersion, tmpRestore string

	tmpRestore = System.Folder.Temp + "restore" + string(os.PathSeparator)

	err = checkFolder(tmpRestore)
	if err != nil {
		return
	}

	// Unpack the ZIP Archive in tmp
	err = extractZIP(archive, tmpRestore)
	if err != nil {
		return
	}

	// Load a new Config to check the Port and Version
	newConfig, err := loadJSONFileToMap(tmpRestore + "settings.json")
	if err != nil {
		ShowError(err, 0)
		return
	}

	backupVersion = newConfig["version"].(string)
	if backupVersion < System.Compatibility {
		err = errors.New(getErrMsg(1013))
		return
	}

	if err = os.RemoveAll(System.Folder.Config); err != nil {
		ShowError(err, 1073)
	}

	// Extract the ZIP Archive into the Config Folder
	err = extractZIP(archive, System.Folder.Config)
	if err != nil {
		return
	}

	// Load a new Config to check the Port and Version
	newConfig, err = loadJSONFileToMap(System.Folder.Config + "settings.json")
	if err != nil {
		ShowError(err, 0)
		return
	}

	newPort = newConfig["port"].(string)
	oldPort = Settings.Port

	if newPort == oldPort {

		if err != nil {
			ShowError(err, 0)
		}

		loadSettings()

		err := Init()
		if err != nil {
			ShowError(err, 0)
			return "", err
		}

		err = StartSystem(true)
		if err != nil {
			ShowError(err, 0)
			return "", err
		}

		return "", err
	}

	var url = System.URLBase + "/web/"
	newWebURL = strings.Replace(url, ":"+oldPort, ":"+newPort, 1)

	os.RemoveAll(tmpRestore)

	return
}

func xteveRestoreFromWeb(input string) (newWebURL string, err error) {

	// Convert base64 JSON string to base64
	b64data := input[strings.IndexByte(input, ',')+1:]

	// Convert Base64 into bytes and save
	sDec, err := b64.StdEncoding.DecodeString(b64data)

	if err != nil {
		return
	}

	var archive = System.Folder.Temp + "restore.zip"

	err = writeByteToFile(archive, sDec)
	if err != nil {
		return
	}

	newWebURL, err = xteveRestore(archive)

	return
}

// XteveRestoreFromCLI : Recovery from the Command Line
func XteveRestoreFromCLI(archive string) (err error) {

	var confirm string

	println()
	showInfo(fmt.Sprintf("Version:%s Build: %s", System.Version, System.Build))
	showInfo(fmt.Sprintf("Backup File:%s", archive))
	showInfo(fmt.Sprintf("System Folder:%s", getPlatformPath(System.Folder.Config)))
	println()

	fmt.Print("All data will be replaced with those from the backup. Should the files be restored? [yes|no]:")

	fmt.Scanln(&confirm)

	switch strings.ToLower(confirm) {

	case "yes":
		break

	case "no":
		return

	default:
		fmt.Println("Invalid input")
		return

	}

	if len(System.Folder.Config) > 0 {

		err = checkFilePermission(System.Folder.Config)
		if err != nil {
			return
		}

		_, err = xteveRestore(archive)
		if err != nil {
			return
		}

		showHighlight(fmt.Sprintf("Restor:Backup was successfully restored. %s can now be started normally", System.Name))

	}
	return
}
