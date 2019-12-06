package src

import (
  b64 "encoding/base64"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "path/filepath"
  "strings"
)

func getCacheImageURL(imageURL string) (cacheImageURL string) {

  if Settings.CacheImages == false {
    return imageURL
  }

  imageURL = strings.Trim(imageURL, "\r\n")

  p, err := url.Parse(imageURL)
  if err != nil {
    // URL konnte nicht geparst werden, die ursprüngliche image url wird zurückgegeben
    showInfo(fmt.Sprintf("Image Caching:Image URL: %s", imageURL))
    showWarning(4101)
    return imageURL
  }
  var urlMD5 = getMD5(imageURL)
  var fileExtension = filepath.Ext(p.Path)

  if len(fileExtension) == 0 {
    // Keine Dateierweiterung vorhanden, die ursprüngliche image url wird zurückgegeben
    return imageURL
  }

  if indexOfString(urlMD5+fileExtension, Data.Cache.ImagesFiles) == -1 {
    Data.Cache.ImagesFiles = append(Data.Cache.ImagesFiles, urlMD5+fileExtension)
  }

  if System.ImageCachingInProgress == 1 {
    return imageURL
  }

  if indexOfString(urlMD5+fileExtension, Data.Cache.ImagesCache) != -1 {

    cacheImageURL = fmt.Sprintf("%s://%s/images/%s%s", System.ServerProtocol.XML, System.Domain, urlMD5, fileExtension)

  } else {

    if strings.Contains(imageURL, System.Domain+"/images/") == false {

      if indexOfString(imageURL, Data.Cache.ImagesURLS) == -1 {
        Data.Cache.ImagesURLS = append(Data.Cache.ImagesURLS, imageURL)
      }

    }

    cacheImageURL = imageURL

  }

  return
}

func cachingImages() {

  if Settings.CacheImages == false || System.ImageCachingInProgress == 1 {
    return
  }

  System.ImageCachingInProgress = 1

  showInfo("Image Caching:Images are cached")

  for _, imageURL := range Data.Cache.ImagesURLS {

    if len(imageURL) > 0 {
      cacheImage(imageURL)
    }

  }

  showInfo("Image Caching:Done")

  // Bilder die nicht mehr verwendet werden, werden gelöscht
  files, err := ioutil.ReadDir(System.Folder.ImagesCache)
  if err != nil {
    ShowError(err, 0)
    return
  }

  for _, file := range files {

    if indexOfString(file.Name(), Data.Cache.ImagesFiles) == -1 {

      var debug = fmt.Sprintf("Image Caching:Remove file: %s %s %d", System.Folder.ImagesCache+file.Name(), file.Name(), len(file.Name()))
      showDebug(debug, 1)
      err := os.RemoveAll(System.Folder.ImagesCache + file.Name())
      if err != nil {
        ShowError(err, 0)
      }

    }

  }

  System.ImageCachingInProgress = 0

  return
}

func cacheImage(imageURL string) {

  var debug string
  var urlMD5 = getMD5(imageURL)
  var fileExtension = filepath.Ext(imageURL)

  debug = fmt.Sprintf("Image Caching:File: %s Download: %s", urlMD5+fileExtension, imageURL)
  showDebug(debug, 1)

  resp, err := http.Get(imageURL)
  if err != nil {
    return
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return
  }

  var filePath = System.Folder.ImagesCache + urlMD5 + fileExtension

  // Datei speichern
  file, err := os.Create(filePath)
  if err != nil {
    return
  }

  defer file.Close()

  _, err = io.Copy(file, resp.Body)

  return
}

func uploadLogo(input, filename string) (logoURL string, err error) {

  b64data := input[strings.IndexByte(input, ',')+1:]

  // BAse64 in bytes umwandeln un speichern
  sDec, err := b64.StdEncoding.DecodeString(b64data)

  if err != nil {
    return
  }

  var file = fmt.Sprintf("%s%s", System.Folder.ImagesUpload, filename)

  err = writeByteToFile(file, sDec)
  if err != nil {
    return
  }

  logoURL = fmt.Sprintf("%s://%s/data_images/%s", System.ServerProtocol.XML, System.Domain, filename)

  return

}
