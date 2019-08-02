package src

import (
  b64 "encoding/base64"
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

  // Alte Backups löschen
  files, err := ioutil.ReadDir(System.Folder.Backup)

  if err == nil {

    for _, file := range files {

      if filepath.Ext(file.Name()) == ".zip" && strings.Contains(file.Name(), "xteve_auto_backup") {
        oldBackupFiles = append(oldBackupFiles, file.Name())
      }

    }

    // Alle Backups löschen
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

  // Backup erstellen
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

func xteveRestore(input string) (newWebURL string, err error) {

  var newPort, oldPort string

  // Base64 Json String in base64 umwandeln
  b64data := input[strings.IndexByte(input, ',')+1:]

  // Base64 in bytes umwandeln und speichern
  sDec, err := b64.StdEncoding.DecodeString(b64data)

  if err != nil {
    return
  }

  var archive = System.Folder.Temp + "restore.zip"

  err = writeByteToFile(archive, sDec)
  if err != nil {
    return
  }

  // Zip Archiv entpacken
  err = extractZIP(archive, System.Folder.Config)
  if err != nil {
    return
  }

  // Neue Config laden um den Port zu überprüfen
  newConfig, err := loadJSONFileToMap(System.Folder.Config + "settings.json")
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

  return
}
