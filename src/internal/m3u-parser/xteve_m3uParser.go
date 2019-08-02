package m3u

import (
  "errors"
  "net/url"
  "regexp"
  "strings"
)

// MakeInterfaceFromM3U :
func MakeInterfaceFromM3U(byteStream []byte) (allChannels []interface{}, err error) {

  var content = string(byteStream)
  var channelName string

  var parseMetaData = func(channel string) (stream map[string]string) {

    stream = make(map[string]string)
    var exceptForParameter = `[a-z-A-Z=]*(".*?")`
    var exceptForChannelName = `,([^\n]*|,[^\r]*)`

    var lines = strings.Split(strings.Replace(channel, "\r\n", "\n", -1), "\n")

    // Zeilen mit # und leerer Zeilen entfernen
    for i := len(lines) - 1; i >= 0; i-- {

      if len(lines[i]) == 0 || lines[i][0:1] == "#" {
        lines = append(lines[:i], lines[i+1:]...)
      }

    }

    if len(lines) >= 2 {

      for _, line := range lines {

        _, err := url.ParseRequestURI(line)

        switch err {

        case nil:
          stream["url"] = strings.Trim(line, "\r\n")

        default:

          var value string
          // Alle Parameter parsen
          var p = regexp.MustCompile(exceptForParameter)
          var streamParameter = p.FindAllString(line, -1)

          for _, p := range streamParameter {

            line = strings.Replace(line, p, "", 1)

            p = strings.Replace(p, `"`, "", -1)
            var parameter = strings.Split(p, "=")

            if len(parameter) == 2 {

              // TVG Key als Kleinbuchstaben speichern
              switch strings.Contains(parameter[0], "tvg") {

              case true:
                stream[strings.ToLower(parameter[0])] = parameter[1]
              case false:
                stream[parameter[0]] = parameter[1]

              }

              // URL's nicht an die Filterfunktion übergeben
              if !strings.Contains(parameter[1], "://") && len(parameter[1]) > 0 {
                value = value + parameter[1] + " "
              }

            }

          }

          // Kanalnamen parsen
          n := regexp.MustCompile(exceptForChannelName)
          var name = n.FindAllString(line, 1)

          if len(name) > 0 {
            channelName = name[0]
            channelName = strings.Replace(channelName, `,`, "", 1)
            channelName = strings.TrimRight(channelName, "\r\n")
            channelName = strings.TrimRight(channelName, " ")
          }

          if len(channelName) == 0 {

            if v, ok := stream["tvg-name"]; ok {
              channelName = v
            }

          }

          channelName = strings.TrimRight(channelName, " ")

          // Kanäle ohne Namen werden augelassen
          if len(channelName) == 0 {
            return
          }

          stream["name"] = channelName
          value = value + channelName

          stream["_values"] = value

        }

      }

    }

    // Nach eindeutiger ID im Stream suchen
    for key, value := range stream {

      if !strings.Contains(strings.ToLower(key), "tvg-id") {

        if strings.Contains(strings.ToLower(key), "id") {

          stream["_uuid.key"] = key
          stream["_uuid.value"] = value
          //os.Exit(0)
          break

        }

      }

    }

    return
  }

  //fmt.Println(content)

  if strings.Contains(content, "#EXTM3U") {

    var channels = strings.Split(content, "#EXTINF")

    channels = append(channels[:0], channels[1:]...)

    for _, channel := range channels {

      var stream = parseMetaData(channel)

      if len(stream) > 0 && stream != nil {
        allChannels = append(allChannels, stream)
      }

    }

  } else {

    err = errors.New("No valid m3u file")

  }

  return
}

// MakeInterfaceFromM3U2 :
func MakeInterfaceFromM3U2(byteStream []byte) (allChannels []interface{}, err error) {
  var content = string(byteStream)
  //var allChannels = make([]interface{}, 0)

  var channels = strings.Split(content, "#EXTINF")

  var parseMetaData = func(metaData string) map[string]string {
    var values string // Save all values in a key
    var channel = make(map[string]string)

    var exceptForParameter = `[a-z-A-Z=]*(".*?")`
    //var exceptForChannelName  = `(,[^.$\n]*|,[^.$\r]*)`
    var exceptForChannelName = `(,[^\n]*|,[^\r]*)`

    var exceptForStreamingURL = `(\n.*?\n|\r.*?\r|\n.*?\z|\r.*?\z)`
    //var exceptForStreamingURL = `^(([^:/?#]+):)?(//([^/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`

    // Parse all parameters
    p := regexp.MustCompile(exceptForParameter)
    var parameter = p.FindAllString(metaData, -1)
    //fmt.Println(parameter)
    for _, i := range parameter {
      var remove = i
      i = strings.Replace(i, `"`, "", -1)
      if strings.Contains(i, "=") {
        var item = strings.Split(i, "=")
        switch strings.Contains(item[0], "tvg") {
        case true:
          channel[strings.ToLower(item[0])] = item[1]
        case false:
          channel[item[0]] = item[1]
        }

        switch strings.Contains(item[1], "://") {
        case false:
          values = values + item[1] + " "
        }

      }
      metaData = strings.Replace(metaData, remove, "", 1)
    }

    // Parse channel name (after the comma)
    n := regexp.MustCompile(exceptForChannelName)
    var name = n.FindAllString(metaData, 1)
    //name[len(name) - 1] = strings.Replace(name[len(name) - 1], `\r`, "", -1)

    var channelName string
    if len(name) == 0 {
      if v, ok := channel["tvg-name"]; ok {
        channelName = v
      }
    } else {
      channelName = name[len(name)-1][1:len(name[len(name)-1])]
    }

    channelName = strings.Replace(channelName, `"`, "", -1)

    var replacer = strings.NewReplacer("\n", "", "\r", "")
    channel["name"] = replacer.Replace(channelName)

    values = values + channelName + " "

    // Parse streaming URL
    u := regexp.MustCompile(exceptForStreamingURL)
    var streamingURL = u.FindAllString(metaData, -1)
    var url = strings.Replace(streamingURL[0], "\n", "", -1)
    url = strings.Replace(url, "\r", "", -1)
    url = strings.Trim(url, "\r\n")
    channel["url"] = url

    channel["_values"] = values

    // Search for a unique ID

    for key, value := range channel {
      if !strings.Contains(strings.ToLower(key), "tvg-id") {
        if strings.Contains(strings.ToLower(key), "id") {
          channel["_uuid.key"] = key
          channel["_uuid.value"] = value
          break
        }
      }
    }

    return channel
  }

  if strings.Contains(channels[0], "#EXTM3U") {

    for _, thisStream := range channels {
      if !strings.Contains(thisStream, "#EXTM3U") {
        var channel = parseMetaData(thisStream)
        allChannels = append(allChannels, channel)
      }
    }

  } else {
    err = errors.New("No valid m3u file")
  }

  return
}
