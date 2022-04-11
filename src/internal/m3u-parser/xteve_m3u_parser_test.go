package m3u

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
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

	// Read playlist
	file := "test_playlist_1.m3u"
	content, err := ioutil.ReadFile(file)
	assert.NoError(t, err, "Should read playlist")

	// Parse playlist into []interface{}
	rawStreams, err := MakeInterfaceFromM3U(content)
	assert.NoError(t, err, "Should parse playlist")

	// Build []M3UStream from []interface{}
	streams := []M3UStream{}

	for _, rawStream := range rawStreams {
		jsonString, err := json.MarshalIndent(rawStream, "", "  ")
		assert.NoError(t, err, "Should convert from interface")

		stream := M3UStream{}
		err = json.Unmarshal(jsonString, &stream)
		assert.NoError(t, err, "Should convert from interface")

		streams = append(streams, stream)
	}

	assert.Len(t, streams, 4, "Should be 4 streams in total")

	// Test stream 1
	assert.Equal(t, "Channel 1", streams[0].Name, "Names should match")
	assert.Equal(t, "Group 1", streams[0].GroupTitle, "Groups should match")
	assert.Equal(t, "http://example.com/stream/1", streams[0].URL, "URL's should match")
	assert.Equal(t, "Channel.1", streams[0].TvgName, "TVG names should match")
	assert.Equal(t, "tvg.id.1", streams[0].TvgID, "TVG ID's should match")
	assert.Equal(t, "https://example/logo.png", streams[0].TvgLogo, "TVG logos should match")
	assert.Empty(t, streams[0].TvgShift, "Should not have tvg-shift tag")

	// Test stream 2
	assert.Equal(t, "Channel 2", streams[1].Name, "Names should match")
	assert.Equal(t, "Group 2", streams[1].GroupTitle, "Should have a GroupTitle set from EXTGRP")
	assert.Equal(t, "http://example.com/stream/2", streams[1].URL, "URL's should match")
	assert.Equal(t, "Channel.2", streams[1].TvgName, "TVG names should match")
	assert.Equal(t, "tvg.id.2", streams[1].TvgID, "TVG ID's should match")
	assert.Equal(t, "https://example/logo/2.png", streams[1].TvgLogo, "TVG logos should match")
	assert.Empty(t, streams[1].TvgShift, "Should not have tvg-shift tag")

	// Test stream 3
	assert.Equal(t, ",:It's - a difficult name |", streams[2].Name, "Names should match")
	assert.Equal(t, "Group 2", streams[2].GroupTitle, "Should have a GroupTitle set from previous EXTGRP")
	assert.Equal(t, "http://example.com/stream/3", streams[2].URL, "URL's should match")
	assert.Empty(t, streams[2].TvgName, "Should not have tvg-name tag")
	assert.Empty(t, streams[2].TvgID, "Should not have tvg-id tag")
	assert.Empty(t, streams[2].TvgLogo, "Should not have tvg-logo tag")
	assert.Empty(t, streams[2].TvgShift, "Should not have tvg-shift tag")

	// Test stream 4
	assert.Equal(t, "Channel 4", streams[3].Name, "Names should match")
	assert.Equal(t, "Group 4", streams[3].GroupTitle, "Should have a GroupTitle set from group-title, over EXTGRP")
	assert.Equal(t, "http://example.com/stream/4", streams[3].URL, "URL's should match")
	assert.Equal(t, "Channel.4", streams[3].TvgName, "TVG names should match")
	assert.Equal(t, "tvg.id.4", streams[3].TvgID, "TVG ID's should match")
	assert.Equal(t, "https://example/logo/4.png", streams[3].TvgLogo, "TVG logos should match")
	assert.Equal(t, "-5", streams[3].TvgShift, "TVG shifts should match")
}
