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
