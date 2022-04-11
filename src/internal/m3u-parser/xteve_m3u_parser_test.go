package m3u

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

type M3UStream struct {
	GroupTitle string `json:"group-title,required"`
	Name       string `json:"name,required"`
	TvgID      string `json:"tvg-id,required"`
	TvgLogo    string `json:"tvg-logo,required"`
	TvgName    string `json:"tvg-name,required"`
	TvgShift   string `json:"tvg-shift,omitempty"`
	URL        string `json:"url,required"`
	UUIDKey    string `json:"_uuid.key,omitempty"`
	UUIDValue  string `json:"_uuid.value,omitempty"`
}

func TestMakeInterfaceFromM3U(t *testing.T) {

	var file = "test_playlist_1.m3u"
	var content, err = ioutil.ReadFile(file)
	if err != nil {
		t.Error(err)
		return
	}

	streams, err := MakeInterfaceFromM3U(content)

	if err != nil {
		t.Error(err)
	}

	err = checkStream(streams)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Streams:", len(streams))
	t.Log(streams)

}

func checkStream(streamInterface []interface{}) (err error) {

	for i, s := range streamInterface {

		var stream = s.(map[string]string)
		var m3uStream M3UStream

		jsonString, err := json.MarshalIndent(stream, "", "  ")

		if err == nil {

			err = json.Unmarshal(jsonString, &m3uStream)
			if err == nil {

				log.Print(fmt.Sprintf("Stream:        %d", i))
				log.Print(fmt.Sprintf("Name*:         %s", m3uStream.Name))
				log.Print(fmt.Sprintf("URL*:          %s", m3uStream.URL))
				log.Print(fmt.Sprintf("tvg-name:      %s", m3uStream.TvgName))
				log.Print(fmt.Sprintf("tvg-id**:      %s", m3uStream.TvgID))
				log.Print(fmt.Sprintf("tvg-logo:      %s", m3uStream.TvgLogo))
				log.Print(fmt.Sprintf("tvg-shift:     %s", m3uStream.TvgShift))
				log.Print(fmt.Sprintf("group-title**: %s", m3uStream.GroupTitle))

				if len(m3uStream.UUIDKey) > 0 {
					log.Print(fmt.Sprintf("UUID key***:   %s", m3uStream.UUIDKey))
					log.Print(fmt.Sprintf("UUID value:    %s", m3uStream.UUIDValue))
				} else {
					log.Print(fmt.Sprintf("UUID key:    false"))
				}

			}

		}

		log.Println(fmt.Sprintf("- - - - - (*: Required) | (**: Nice to have) | (***: Love it) - - - - -"))
	}

	return
}
