package m3u

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
)

// MakeInterfaceFromM3U :
func MakeInterfaceFromM3U(byteStream []byte) (allChannels []interface{}, err error) {

	var content = string(byteStream)
	var channelName string
	var uuids []string

	var parseMetaData = func(channel string) (stream map[string]string) {

		stream = make(map[string]string)
		var exceptForParameter = `[a-z-A-Z=]*(".*?")`
		var exceptForChannelName = `,([^\n]*|,[^\r]*)`

		var lines = strings.Split(strings.Replace(channel, "\r\n", "\n", -1), "\n")

		// Remove lines with # and blank lines
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
					// Parse all parameters
					var p = regexp.MustCompile(exceptForParameter)
					var streamParameter = p.FindAllString(line, -1)

					for _, p := range streamParameter {

						line = strings.Replace(line, p, "", 1)

						p = strings.Replace(p, `"`, "", -1)
						var parameter = strings.SplitN(p, "=", 2)

						if len(parameter) == 2 {

							// Set TVG Key as lowercase
							switch strings.Contains(parameter[0], "tvg") {

							case true:
								stream[strings.ToLower(parameter[0])] = parameter[1]
							case false:
								stream[parameter[0]] = parameter[1]

							}

							// URL's are not passed to the filter function
							if !strings.Contains(parameter[1], "://") && len(parameter[1]) > 0 {
								value = value + parameter[1] + " "
							}

						}

					}

					// Parse channel names
					n := regexp.MustCompile(exceptForChannelName)
					var name = n.FindAllString(line, 1)

					if len(name) > 0 {
						channelName = name[0]
						channelName = strings.Replace(channelName, `,`, "", 1)
						channelName = strings.TrimRight(channelName, "\r\n")
						channelName = strings.Trim(channelName, " ")
					}

					if len(channelName) == 0 {

						if v, ok := stream["tvg-name"]; ok {
							channelName = v
						}

					}

					channelName = strings.Trim(channelName, " ")

					// Channels without names are skipped
					if len(channelName) == 0 {
						return
					}

					stream["name"] = channelName
					value = value + channelName

					stream["_values"] = value

				}

			}

		}

		// Search for a unique ID in the stream
		for key, value := range stream {

			if !strings.Contains(strings.ToLower(key), "tvg-id") {

				if strings.Contains(strings.ToLower(key), "id") {

					if indexOfString(value, uuids) != -1 {
						log.Println(fmt.Sprintf("Channel: %s - %s = %s ", stream["name"], key, value))
						break
					}

					uuids = append(uuids, value)

					stream["_uuid.key"] = key
					stream["_uuid.value"] = value
					break

				}

			}

		}

		return
	}

	if strings.Contains(content, "#EXT-X-TARGETDURATION") || strings.Contains(content, "#EXT-X-MEDIA-SEQUENCE") {
		err = errors.New("Invalid M3U file, an extended M3U file is required.")
		return
	}

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

		err = errors.New("Invalid M3U file, an extended M3U file is required.")

	}

	return
}

func indexOfString(element string, data []string) int {

	for k, v := range data {
		if element == v {
			return k
		}
	}

	return -1
}
