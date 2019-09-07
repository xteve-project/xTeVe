package src

/*
	Tuner-Limit Bild als Video rendern [ffmpeg]
	-loop 1 -i stream-limit.jpg -c:v libx264 -t 1 -pix_fmt yuv420p -vf scale=1920:1080  stream-limit.ts
*/

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

func createStreamID(stream map[int]ThisStream) (streamID int) {

	var debug string

	streamID = 0
	for i := 0; i <= len(stream); i++ {

		if _, ok := stream[i]; !ok {
			streamID = i
			break
		}

	}

	debug = fmt.Sprintf("Streaming Status:Stream ID = %d", streamID)
	showDebug(debug, 1)

	return
}

func bufferingStream(playlistID, streamingURL, channelName string, w http.ResponseWriter, r *http.Request) {

	time.Sleep(time.Duration(Settings.BufferTimeout) * time.Millisecond)

	var playlist Playlist
	var client ThisClient
	var stream ThisStream
	var streaming = false
	var streamID int
	var debug string
	var timeOut = 0
	var newStream = true

	//w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Connection", "close")

	// Überprüfen ob die Playlist schon verwendet wird
	if p, ok := BufferInformation.Load(playlistID); !ok {

		var playlistType string
		// Playlist wird noch nicht verwendet, Default-Werte für die Playlist erstellen
		playlist.Folder = System.Folder.Temp + playlistID + string(os.PathSeparator)
		playlist.PlaylistID = playlistID
		playlist.Streams = make(map[int]ThisStream)
		playlist.Clients = make(map[int]ThisClient)

		err := checkFolder(playlist.Folder)
		if err != nil {
			ShowError(err, 000)
			httpStatusError(w, r, 404)
			return
		}

		switch playlist.PlaylistID[0:1] {

		case "M":
			playlistType = "m3u"

		case "H":
			playlistType = "hdhr"

		}

		playlist.Tuner = getTuner(playlistID, playlistType)

		playlist.PlaylistName = getProviderParameter(playlist.PlaylistID, playlistType, "name")

		// Default-Werte für den Stream erstellen
		streamID = createStreamID(playlist.Streams)

		client.Connection = 1
		stream.URL = streamingURL
		stream.ChannelName = channelName
		stream.Status = false

		playlist.Streams[streamID] = stream
		playlist.Clients[streamID] = client

		BufferInformation.Store(playlistID, playlist)

	} else {

		// Playlist wird bereits zum streamen verwendet
		// Überprüfen ob die URL bereit von einem anderen Client gestreamt wird.

		playlist = p.(Playlist)

		for id := range playlist.Streams {

			stream = playlist.Streams[id]
			client = playlist.Clients[id]

			if streamingURL == stream.URL {

				streamID = id
				newStream = false
				client.Connection++

				//playlist.Streams[streamID] = stream
				playlist.Clients[streamID] = client

				BufferInformation.Store(playlistID, playlist)

				debug = fmt.Sprintf("Restream Status:Playlist: %s - Channel: %s - Connections: %d", playlist.PlaylistName, stream.ChannelName, client.Connection)

				showDebug(debug, 1)

				if c, ok := BufferClients.Load(playlistID + stream.MD5); ok {

					var clients = c.(ClientConnection)
					clients.Connection = clients.Connection + 1
					showInfo(fmt.Sprintf("Streaming Status:Channel: %s (Clients: %d)", stream.ChannelName, clients.Connection))

					BufferClients.Store(playlistID+stream.MD5, clients)

				}

				break
			}

		}

		// Neuer Stream bei einer bereits aktiven Playlist
		if newStream == true {

			// Prüfen ob die Playlist noch einen weiteren Stream erlaubt (Tuner)
			if len(playlist.Streams) >= playlist.Tuner {

				showInfo(fmt.Sprintf("Streaming Status:Playlist: %s - No new connections available. Tuner = %d", playlist.PlaylistName, playlist.Tuner))

				if value, ok := webUI["html/video/stream-limit.ts"]; ok {

					var content string
					content = GetHTMLString(value.(string))

					w.WriteHeader(200)
					w.Header().Set("Content-type", "video/mpeg")
					w.Header().Set("Content-Length:", "0")

					for i := 1; i < 60; i++ {
						_ = i
						w.Write([]byte(content))
						time.Sleep(time.Duration(500) * time.Millisecond)
					}

					return
				}

				return
			}

			// Playlist erlaubt einen weiterern Stream (Das Limit des Tuners ist noch nicht erreicht)
			// Default-Werte für den Stream erstellen
			stream = ThisStream{}
			client = ThisClient{}

			streamID = createStreamID(playlist.Streams)

			client.Connection = 1
			stream.URL = streamingURL
			stream.ChannelName = channelName
			stream.Status = false

			playlist.Streams[streamID] = stream
			playlist.Clients[streamID] = client

			BufferInformation.Store(playlistID, playlist)

		}

	}

	// Überprüfen ob der Stream breits von einem anderen Client abgespielt wird
	if playlist.Streams[streamID].Status == false && newStream == true {

		// Neuer Buffer wird benötigt
		stream = playlist.Streams[streamID]
		stream.MD5 = getMD5(streamingURL)
		stream.Folder = playlist.Folder + stream.MD5 + string(os.PathSeparator)
		stream.PlaylistID = playlistID
		stream.PlaylistName = playlist.PlaylistName

		playlist.Streams[streamID] = stream
		BufferInformation.Store(playlistID, playlist)

		go connectToStreamingServer(streamID, playlistID)

		showInfo(fmt.Sprintf("Streaming Status:Playlist: %s - Tuner: %d / %d", playlist.PlaylistName, len(playlist.Streams), playlist.Tuner))

		var clients ClientConnection
		clients.Connection = 1
		BufferClients.Store(playlistID+stream.MD5, clients)

	}

	w.WriteHeader(200)

	for { // Loop 1: Warten bis das erste Segment durch den Buffer heruntergeladen wurde

		if p, ok := BufferInformation.Load(playlistID); ok {

			var playlist = p.(Playlist)

			if stream, ok := playlist.Streams[streamID]; ok {

				if stream.Status == false {

					timeOut++

					time.Sleep(time.Duration(100) * time.Millisecond)

					if c, ok := BufferClients.Load(playlistID + stream.MD5); ok {

						var clients = c.(ClientConnection)

						if clients.Error != nil || timeOut > 200 {
							killClientConnection(streamID, stream.PlaylistID, false)
							return
						}

					}

					continue
				}

				var oldSegments []string

				for { // Loop 2: Temporäre Datein sind vorhanden, Daten können zum Client gesendet werden
					// HTTP Clientverbindung überwachen
					cn, ok := w.(http.CloseNotifier)
					if ok {

						select {

						case <-cn.CloseNotify():
							killClientConnection(streamID, playlistID, false)
							return

						default:
							if c, ok := BufferClients.Load(playlistID + stream.MD5); ok {

								var clients = c.(ClientConnection)
								if clients.Error != nil {
									ShowError(clients.Error, 0)
									killClientConnection(streamID, playlistID, false)
									return
								}

							} else {

								return

							}

						}

					}

					if _, err := os.Stat(stream.Folder); os.IsNotExist(err) {
						killClientConnection(streamID, playlistID, false)
						return
					}

					var tmpFiles = getTmpFiles(&stream)
					//fmt.Println("Buffer Loop:", stream.Connection)

					for _, f := range tmpFiles {

						if _, err := os.Stat(stream.Folder); os.IsNotExist(err) {
							killClientConnection(streamID, playlistID, false)
							return
						}

						oldSegments = append(oldSegments, f)

						var fileName = stream.Folder + f

						file, err := os.Open(fileName)
						defer file.Close()

						if err == nil {

							l, err := file.Stat()
							if err == nil {

								debug = fmt.Sprintf("Buffer Status:Send to client (%s)", fileName)
								showDebug(debug, 2)

								var buffer = make([]byte, int(l.Size()))
								_, err = file.Read(buffer)

								if err == nil {

									file.Seek(0, 0)

									if streaming == false {

										contentType := http.DetectContentType(buffer)
										_ = contentType
										//w.Header().Set("Content-type", "video/mpeg")
										w.Header().Set("Content-type", contentType)
										w.Header().Set("Content-Length", "0")
										w.Header().Set("Connection", "close")

									}

									/*
										// HDHR Header
										w.Header().Set("Cache-Control", "no-cache")
										w.Header().Set("Pragma", "no-cache")
										w.Header().Set("transferMode.dlna.org", "Streaming")
									*/

									_, err := w.Write(buffer)

									if err != nil {
										file.Close()
										killClientConnection(streamID, playlistID, false)
										return
									}

									file.Close()
									streaming = true

								}

								file.Close()

							}

							var n = indexOfString(f, oldSegments)

							if n > 20 {

								var fileToRemove = stream.Folder + oldSegments[0]
								os.RemoveAll(getPlatformFile(fileToRemove))
								oldSegments = append(oldSegments[:0], oldSegments[0+1:]...)

							}

						}

						file.Close()

					}

					if len(tmpFiles) == 0 {
						time.Sleep(time.Duration(100) * time.Millisecond)
					}

				} // Ende Loop 2

			} else {

				// Stream nicht vorhanden
				killClientConnection(streamID, stream.PlaylistID, false)
				showInfo(fmt.Sprintf("Streaming Status:Playlist: %s - Tuner: %d / %d", playlist.PlaylistName, len(playlist.Streams), playlist.Tuner))
				return

			}

		} // Ende BufferInformation

	} // Ende Loop 1

}

func getTmpFiles(stream *ThisStream) (tmpFiles []string) {

	var tmpFolder = stream.Folder
	var fileIDs []float64

	if _, err := os.Stat(tmpFolder); !os.IsNotExist(err) {

		files, err := ioutil.ReadDir(getPlatformPath(tmpFolder))
		if err != nil {
			ShowError(err, 000)
			return
		}

		if len(files) > 1 {

			for _, file := range files {

				var fileID = strings.Replace(file.Name(), ".ts", "", -1)
				var f, err = strconv.ParseFloat(fileID, 64)

				if err == nil {
					fileIDs = append(fileIDs, f)
				}

			}

			sort.Float64s(fileIDs)
			fileIDs = fileIDs[:len(fileIDs)-1]

			for _, file := range fileIDs {

				var fileName = fmt.Sprintf("%d.ts", int64(file))

				if indexOfString(fileName, stream.OldSegments) == -1 {
					tmpFiles = append(tmpFiles, fileName)
					stream.OldSegments = append(stream.OldSegments, fileName)
				}

			}

		}

	}

	return
}

func killClientConnection(streamID int, playlistID string, force bool) {

	if p, ok := BufferInformation.Load(playlistID); ok {

		var playlist = p.(Playlist)

		if force == true {
			delete(playlist.Streams, streamID)
			showInfo(fmt.Sprintf("Streaming Status:Playlist: %s - Tuner: %d / %d", playlist.PlaylistName, len(playlist.Streams), playlist.Tuner))
			return
		}

		if stream, ok := playlist.Streams[streamID]; ok {

			if c, ok := BufferClients.Load(playlistID + stream.MD5); ok {

				var clients = c.(ClientConnection)
				clients.Connection = clients.Connection - 1
				BufferClients.Store(playlistID+stream.MD5, clients)

				showInfo("Streaming Status:Client has terminated the connection")
				showInfo(fmt.Sprintf("Streaming Status:Channel: %s (Clients: %d)", stream.ChannelName, clients.Connection))

				if clients.Connection <= 0 {
					BufferClients.Delete(playlistID + stream.MD5)
					delete(playlist.Streams, streamID)
				}

			}

			BufferInformation.Store(playlistID, playlist)

			if len(playlist.Streams) > 0 {
				showInfo(fmt.Sprintf("Streaming Status:Playlist: %s - Tuner: %d / %d", playlist.PlaylistName, len(playlist.Streams), playlist.Tuner))
			}

		}

	}

}

func clientConnection(stream ThisStream) (status bool) {

	status = true

	if _, ok := BufferClients.Load(stream.PlaylistID + stream.MD5); !ok {

		var debug = fmt.Sprintf("Streaming Status:Remove temporary files (%s)", stream.Folder)
		showDebug(debug, 1)

		status = false

		debug = fmt.Sprintf("Remove tmp folder:%s", stream.Folder)
		showDebug(debug, 1)

		os.RemoveAll(stream.Folder)

		if p, ok := BufferInformation.Load(stream.PlaylistID); ok {

			showInfo(fmt.Sprintf("Streaming Status:Channel: %s - No client is using this channel anymore. Streaming Server connection has ended", stream.ChannelName))

			var playlist = p.(Playlist)

			showInfo(fmt.Sprintf("Streaming Status:Playlist: %s - Tuner: %d / %d", playlist.PlaylistName, len(playlist.Streams), playlist.Tuner))

			if len(playlist.Streams) <= 0 {
				BufferInformation.Delete(stream.PlaylistID)
			}

		}

		status = false

	}

	return
}

func connectToStreamingServer(streamID int, playlistID string) {

	if p, ok := BufferInformation.Load(playlistID); ok {

		var playlist = p.(Playlist)

		var timeOut = 0
		var debug string
		var tmpSegment = 1
		var tmpFolder = playlist.Streams[streamID].Folder
		var m3u8Segments []string
		var bandwidth BandwidthCalculation
		var networkBandwidth = Settings.M3U8AdaptiveBandwidthMBPS * 1e+6

		var defaultSegment = func() {

			var segment Segment

			if len(playlist.Streams[streamID].Location) > 0 {
				segment.URL = playlist.Streams[streamID].Location
			} else {
				segment.URL = playlist.Streams[streamID].URL
			}

			segment.Duration = 0

			var stream = playlist.Streams[streamID]
			stream.Segment = []Segment{}
			stream.Segment = append(stream.Segment, segment)

			stream.HLS = false
			stream.Sequence = 0
			stream.Wait = 0
			stream.NetworkBandwidth = networkBandwidth

			playlist.Streams[streamID] = stream

			timeOut++

		}

		var addErrorToStream = func(err error) {

			var stream = playlist.Streams[streamID]

			if c, ok := BufferClients.Load(playlistID + stream.MD5); ok {

				var clients = c.(ClientConnection)
				clients.Error = err
				BufferClients.Store(playlistID+stream.MD5, clients)

			}

		}

		os.RemoveAll(getPlatformPath(tmpFolder))

		err := checkFolder(tmpFolder)
		if err != nil {
			ShowError(err, 0)
			addErrorToStream(err)
			return
		}

		// M3U8 Segmente
	InitBuffer:
		defaultSegment()

		if len(m3u8Segments) > 30 {
			m3u8Segments = m3u8Segments[15:]
		}
		if timeOut >= 10 {
			return
		}

		var stream ThisStream = playlist.Streams[streamID]

		if stream.Status == false {

			if strings.Index(stream.URL, ".m3u8") != -1 {
				showInfo("Streaming Type:" + "[HLS / M3U8]")
			} else {
				showInfo("Streaming Type:" + "[TS]")
			}

			showInfo("Streaming URL:" + stream.URL)

		}

		var s = 0

		stream.TimeStart = time.Now()
		bandwidth.Start = stream.TimeStart
		bandwidth.Size = 0

		for {

			if clientConnection(stream) == false {
				return
			}

			if len(stream.Segment) == 0 || len(stream.URL) == 0 {
				goto InitBuffer
			}

			var segment = stream.Segment[0]

			var currentURL = strings.Trim(segment.URL, "\r\n")

			if len(currentURL) == 0 {
				goto InitBuffer
			}

			debug = fmt.Sprintf("Connection to:%s", currentURL)
			showDebug(debug, 2)

			// Sprung für Redirect (301 <---> 308)
		Redirect:

			req, err := http.NewRequest("GET", currentURL, nil)
			req.Header.Set("User-Agent", Settings.UserAgent)
			req.Header.Set("Connection", "close")
			//req.Header.Set("Range", "bytes=0-")
			req.Header.Set("Accept", "*/*")
			debugRequest(req)

			client := &http.Client{}
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return errors.New("Redirect")
			}

			resp, err := client.Do(req)

			if resp != nil && err != nil {
				debugResponse(resp)
			}

			if err != nil {

				if resp == nil {

					err = errors.New("No response from streaming server")
					fmt.Println("Current URL:", currentURL)
					ShowError(err, 0)

					addErrorToStream(err)

					killClientConnection(streamID, playlistID, true)
					clientConnection(stream)

					return
				}

				// Redirect
				if resp.StatusCode >= 301 && resp.StatusCode <= 308 {

					debug = fmt.Sprintf("Streaming Status:HTTP response status [%d] %s", resp.StatusCode, http.StatusText(resp.StatusCode))
					showDebug(debug, 2)

					currentURL = strings.Trim(resp.Header.Get("Location"), "\r\n")

					stream.Location = currentURL

					if len(currentURL) > 0 {

						debug = fmt.Sprintf("HTTP Redirect:%s", stream.Location)
						showDebug(debug, 2)
						defer resp.Body.Close()
						goto Redirect

					} else {

						err = errors.New("Streaming server")
						ShowError(err, 4002)
						addErrorToStream(err)

						defer resp.Body.Close()

						return

					}

				} else {

					ShowError(err, 0)
					addErrorToStream(err)

					defer resp.Body.Close()

					return

				}

				defer resp.Body.Close()

			}

			defer resp.Body.Close()

			// HTTP Status überprüfen, bei Fehlern wird der Stream beendet
			var contentType = resp.Header.Get("Content-Type")
			var httpStatusCode = resp.StatusCode
			var httpStatusInfo = fmt.Sprintf("HTTP Response Status [%d] %s", httpStatusCode, http.StatusText(resp.StatusCode))

			if resp.StatusCode != http.StatusOK {

				showInfo("Content Type:" + contentType)
				showInfo("Streaming Status:" + httpStatusInfo)
				showInfo("Error with this URL:" + currentURL)

				var err = errors.New(http.StatusText(resp.StatusCode))
				ShowError(err, 4004)

				debug = fmt.Sprintf("Streaming Status:Playlist: %s - Tuner: %d / %d", playlist.PlaylistName, len(playlist.Streams), playlist.Tuner)
				showDebug(debug, 1)

				BufferInformation.Store(playlist.PlaylistID, playlist)
				addErrorToStream(err)

				killClientConnection(streamID, playlistID, true)
				clientConnection(stream)
				resp.Body.Close()

				return
			}

			// Informationen über den Streamingserver auslesen
			if stream.Status == false {

				if len(stream.URLStreamingServer) == 0 {

					u, _ := url.Parse(currentURL)
					p, _ := url.Parse(currentURL)

					stream.URLScheme = u.Scheme
					stream.URLHost = req.Host
					stream.URLPath = p.Path
					stream.URLFile = path.Base(p.Path)

					stream.URLRedirect = fmt.Sprintf("%s://%s%s", stream.URLScheme, stream.URLHost, stream.URLPath)
					stream.URLStreamingServer = fmt.Sprintf("%s://%s", stream.URLScheme, stream.URLHost)

				}

				debug = fmt.Sprintf("Server URL:%s", stream.URLStreamingServer)
				showDebug(debug, 1)

				debug = fmt.Sprintf("Temp Folder:%s", tmpFolder)
				showDebug(debug, 1)

				showInfo("Streaming Status:" + "HTTP Response Status [" + strconv.Itoa(resp.StatusCode) + "] " + http.StatusText(resp.StatusCode))
				showInfo("Content Type:" + contentType)

			} else {

				debug = fmt.Sprintf("Content Type:%s", contentType)
				showDebug(debug, 2)

			}

			// Content Type bereinigen
			if len(contentType) > 0 {
				var ct = strings.SplitN(contentType, ";", 2)
				contentType = strings.ToLower(ct[0])
			}

			switch contentType {

			// M3U8 Playlist
			case "application/x-mpegurl", "application/vnd.apple.mpegurl", "audio/mpegurl":
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					ShowError(err, 0)
					addErrorToStream(err)
				}

				stream.Body = string(body)
				stream.HLS = true
				stream.M3U8URL = currentURL

				err = parseM3U8(&stream)
				if err != nil {
					ShowError(err, 4050)
					addErrorToStream(err)
				}

			// Video Stream (TS)
			case "video/mpeg", "video/mp4", "video/mp2t", "video/m2ts", "application/octet-stream", "binary/octet-stream", "application/mp2t":

				var fileSize int

				// Größe des Buffers
				buffer := make([]byte, 1024*Settings.BufferSize*2)
				var tmpFileSize = 1024 * Settings.BufferSize * 1

				debug = fmt.Sprintf("Buffer Size:%d KB [SERVER CONNECTION]", len(buffer)/1024)
				showDebug(debug, 3)

				debug = fmt.Sprintf("Buffer Size:%d KB [CLIENT CONNECTION]", tmpFileSize/1024)
				showDebug(debug, 3)

				var tmpFile = fmt.Sprintf("%s%d.ts", tmpFolder, tmpSegment)

				if clientConnection(stream) == false {
					resp.Body.Close()
					return
				}

				bufferFile, err := os.Create(tmpFile)
				if err != nil {

					addErrorToStream(err)
					bufferFile.Close()
					resp.Body.Close()
					return

				}

				for {

					if fileSize == 0 {

						debug = fmt.Sprintf("Buffer Status:Buffering (%s)", tmpFile)
						showDebug(debug, 2)

					}

					timeOut = 0
					// Buffer mit Daten vom Server füllen
					n, err := resp.Body.Read(buffer)

					if err != nil && err != io.EOF {

						ShowError(err, 0)
						addErrorToStream(err)
						resp.Body.Close()
						return

					}

					defer resp.Body.Close()

					if _, err := bufferFile.Write(buffer[:n]); err != nil {

						ShowError(err, 0)
						addErrorToStream(err)
						resp.Body.Close()
						return

					}

					defer bufferFile.Close()

					fileSize = fileSize + n

					if clientConnection(stream) == false {

						resp.Body.Close()
						bufferFile.Close()

						err = os.RemoveAll(stream.Folder)
						if err != nil {
							ShowError(err, 4005)
						}
						return

					}

					// Buffer auf die Festplatte speichern
					if fileSize >= tmpFileSize || n == 0 {

						bandwidth.Stop = time.Now()
						bandwidth.Size += fileSize

						bandwidth.TimeDiff = bandwidth.Stop.Sub(bandwidth.Start).Seconds()

						networkBandwidth = int(float64(bandwidth.Size) / bandwidth.TimeDiff * 1000)

						stream.NetworkBandwidth = networkBandwidth
						bandwidth.NetworkBandwidth = stream.NetworkBandwidth

						debug = fmt.Sprintf("Buffer Status:Done (%s)", tmpFile)
						showDebug(debug, 2)

						bufferFile.Close()

						stream.Status = true
						playlist.Streams[streamID] = stream
						BufferInformation.Store(playlistID, playlist)

						tmpSegment++

						tmpFile = fmt.Sprintf("%s%d.ts", tmpFolder, tmpSegment)

						if clientConnection(stream) == false {

							bufferFile.Close()
							resp.Body.Close()

							err = os.RemoveAll(stream.Folder)
							if err != nil {
								ShowError(err, 4005)
							}

							return
						}

						bufferFile, err = os.Create(tmpFile)
						if err != nil {
							addErrorToStream(err)
							resp.Body.Close()
							return
						}

						fileSize = 0

						if n == 0 {
							bufferFile.Close()
							resp.Body.Close()
							break
						}

					}

				}

				//--

			// Umbekanntes Format
			default:
				showInfo("Content Type:" + resp.Header.Get("Content-Type"))
				err = errors.New("Streaming error")
				ShowError(err, 4003)

				addErrorToStream(err)
				resp.Body.Close()
				return
			}

			s++

			// Wartezeit für den Download das nächste Segments berechnen
			if stream.HLS == true {

				var sleep float64

				if segment.Duration > 0 {

					stream.TimeEnd = time.Now()
					stream.TimeDiff = stream.TimeEnd.Sub(stream.TimeStart).Seconds()

					sleep = (segment.Duration - stream.TimeDiff) - (segment.Duration * 0.25)

					if sleep < 0 {
						sleep = 0
					}

					debug = fmt.Sprintf("HLS Status:Download time: %f s | Segment duration: %f s | Sleep: %f s Sequence: %d", stream.TimeDiff, segment.Duration, sleep, segment.Sequence)
					showDebug(debug, 1)

					if sleep > 0 {

						for i := 0.0; i < sleep*1000; i = i + 100 {

							_ = i
							time.Sleep(time.Duration(100) * time.Millisecond)

							if _, err := os.Stat(stream.Folder); os.IsNotExist(err) {
								break
							}

						}

					}

				}

			}

			stream.Segment = stream.Segment[1:len(stream.Segment)]

			resp.Body.Close()

		} // Ende for loop

	} // Ende BufferInformation

}

func parseM3U8(stream *ThisStream) (err error) {

	var debug string
	var noNewSegment = false
	var lastSegmentDuration float64
	var segment Segment
	var m3u8Segments []Segment
	var sequence int64

	stream.DynamicBandwidth = false

	debug = fmt.Sprintf(`M3U8 Playlist:`+"\n"+`%s`, stream.Body)
	showDebug(debug, 3)

	var getBandwidth = func(line string) int {

		var infos = strings.Split(line, ",")

		for _, info := range infos {

			if strings.Contains(info, "BANDWIDTH=") {

				var bandwidth = strings.Replace(info, "BANDWIDTH=", "", -1)
				n, err := strconv.Atoi(bandwidth)
				if err == nil {
					return n
				}

			}

		}

		return 0
	}

	var parseParameter = func(line string, segment *Segment) (err error) {

		line = strings.Trim(line, "\r\n")

		var parameters = []string{"#EXT-X-VERSION:", "#EXT-X-PLAYLIST-TYPE:", "#EXT-X-MEDIA-SEQUENCE:", "#EXT-X-STREAM-INF:", "#EXTINF:"}

		for _, parameter := range parameters {

			if strings.Contains(line, parameter) {

				var value = strings.Replace(line, parameter, "", -1)

				switch parameter {

				case "#EXT-X-VERSION:":
					version, err := strconv.Atoi(value)
					if err == nil {
						segment.Version = version
					}

				case "#EXT-X-PLAYLIST-TYPE:":
					segment.PlaylistType = value

				case "#EXT-X-MEDIA-SEQUENCE:":
					n, err := strconv.ParseInt(value, 10, 64)
					if err == nil {
						stream.Sequence = n
						sequence = n
					}

				case "#EXT-X-STREAM-INF:":
					segment.Info = true
					segment.StreamInf.Bandwidth = getBandwidth(value)

				case "#EXTINF:":
					var d = strings.Split(value, ",")
					if len(d) > 0 {

						value = strings.Replace(d[0], ",", "", -1)
						duration, err := strconv.ParseFloat(value, 64)
						if err == nil {
							segment.Duration = duration
						} else {
							ShowError(err, 1050)
							return err
						}

					}

				}

			}

		}

		return
	}

	var parseURL = func(line string, segment *Segment) {

		// Prüfen ob die Adresse eine gültige URL ist (http://... oder /path/to/stream)
		_, err := url.ParseRequestURI(line)
		if err == nil {

			// Prüfen ob die Domain in der Adresse enhalten ist
			u, _ := url.Parse(line)

			if len(u.Host) == 0 {
				// Adresse enthällt nicht die Domain, Redirect wird der Adresse hinzugefügt
				segment.URL = stream.URLStreamingServer + line
			} else {
				// Domain in der Adresse enthalten
				segment.URL = line
			}

		} else {

			// keine URL, sondern ein Dateipfad (media/file-01.ts)
			var serverURLPath = strings.Replace(stream.M3U8URL, path.Base(stream.M3U8URL), line, -1)
			segment.URL = serverURLPath

		}

		return
	}

	if strings.Contains(stream.Body, "#EXTM3U") {

		var lines = strings.Split(strings.Replace(stream.Body, "\r\n", "\n", -1), "\n")

		if stream.DynamicBandwidth == false {
			stream.DynamicStream = make(map[int]DynamicStream)
		}

		// Parameter parsen
		for i, line := range lines {

			_ = i

			if len(line) > 0 {

				if line[0:1] == "#" {

					err := parseParameter(line, &segment)
					if err != nil {
						return err
					}

					lastSegmentDuration = segment.Duration

				}

				// M3U8 enthällt mehrere Links zu weiteren M3U8 Wiedergabelisten (Bandbreitenoption)
				if segment.Info == true && len(line) > 0 && line[0:1] != "#" {

					var dynamicStream DynamicStream

					segment.Duration = 0
					noNewSegment = false

					stream.DynamicBandwidth = true
					parseURL(line, &segment)

					dynamicStream.Bandwidth = segment.StreamInf.Bandwidth
					dynamicStream.URL = segment.URL

					stream.DynamicStream[dynamicStream.Bandwidth] = dynamicStream

				}

				// Segment mit TS Stream
				if segment.Duration > 0 && line[0:1] != "#" {

					parseURL(line, &segment)

					if len(segment.URL) > 0 {
						segment.Sequence = sequence
						m3u8Segments = append(m3u8Segments, segment)
						sequence++
					}

				}

			}

		}

	} else {

		err = errors.New(getErrMsg(4051))
		return
	}

	if len(m3u8Segments) > 0 {

		noNewSegment = true

		if stream.Status == false {

			if len(m3u8Segments) >= 2 {
				m3u8Segments = m3u8Segments[0 : len(m3u8Segments)-1]
			}

		}

		for _, s := range m3u8Segments {

			segment = s

			if stream.Status == false {

				noNewSegment = false
				stream.LastSequence = segment.Sequence

				// Stream ist vom Typ VOD. Es muss das erste Segment der M3U8 Playlist verwendet werden.
				if strings.ToUpper(segment.PlaylistType) == "VOD" {
					break
				}

			} else {

				if segment.Sequence > stream.LastSequence {

					stream.LastSequence = segment.Sequence
					noNewSegment = false
					break

				}

			}

		}

	}

	if noNewSegment == false {

		if stream.DynamicBandwidth == true {
			switchBandwidth(stream)
		} else {
			stream.Segment = append(stream.Segment, segment)
		}

	}

	if noNewSegment == true {

		var sleep = lastSegmentDuration * 0.5

		for i := 0.0; i < sleep*1000; i = i + 100 {

			_ = i
			time.Sleep(time.Duration(100) * time.Millisecond)

			if _, err := os.Stat(stream.Folder); os.IsNotExist(err) {
				break
			}

			err := checkFile(stream.Folder + "remove")
			if err == nil {
				os.RemoveAll(stream.Folder)
				break
			}

		}

	}

	return
}

func switchBandwidth(stream *ThisStream) (err error) {

	var bandwidth []int
	var dynamicStream DynamicStream
	var segment Segment

	for key := range stream.DynamicStream {
		bandwidth = append(bandwidth, key)
	}

	sort.Ints(bandwidth)

	if len(bandwidth) > 0 {

		for i := range bandwidth {

			segment.StreamInf.Bandwidth = stream.DynamicStream[bandwidth[i]].Bandwidth

			dynamicStream = stream.DynamicStream[bandwidth[0]]

			if stream.NetworkBandwidth == 0 {

				dynamicStream = stream.DynamicStream[bandwidth[0]]
				break

			} else {

				if bandwidth[i] > stream.NetworkBandwidth {
					break
				}

				dynamicStream = stream.DynamicStream[bandwidth[i]]

			}

		}

	} else {

		err = errors.New("M3U8 does not contain streaming URLs")
		return

	}

	segment.URL = dynamicStream.URL
	segment.Duration = 0
	stream.Segment = append(stream.Segment, segment)

	return
}

func getTuner(id, playlistType string) (tuner int) {

	switch Settings.Buffer {

	case false:
		tuner = Settings.Tuner

	case true:

		i, err := strconv.Atoi(getProviderParameter(id, playlistType, "tuner"))
		if err == nil {
			tuner = i
		} else {
			ShowError(err, 0)
			tuner = 1
		}

	}

	return
}

func debugRequest(req *http.Request) {

	var debugLevel = 3

	if System.Flag.Debug < debugLevel {
		return
	}

	var debug string

	fmt.Println()
	debug = "Request:* * * * * * BEGIN HTTP(S) REQUEST * * * * * * "
	showDebug(debug, debugLevel)

	debug = fmt.Sprintf("Method:%s", req.Method)
	showDebug(debug, debugLevel)

	debug = fmt.Sprintf("Proto:%s", req.Proto)
	showDebug(debug, debugLevel)

	debug = fmt.Sprintf("URL:%s", req.URL)
	showDebug(debug, debugLevel)

	for name, headers := range req.Header {

		name = strings.ToLower(name)

		for _, h := range headers {
			debug = fmt.Sprintf("Header:%v: %v", name, h)
			showDebug(debug, debugLevel)
		}

	}

	debug = "Request:* * * * * * END HTTP(S) REQUEST * * * * * *"
	showDebug(debug, debugLevel)

	return
}

func debugResponse(resp *http.Response) {

	var debugLevel = 3

	if System.Flag.Debug < debugLevel {
		return
	}

	var debug string

	fmt.Println()

	debug = "Response:* * * * * * BEGIN RESPONSE * * * * * * "
	showDebug(debug, debugLevel)

	debug = fmt.Sprintf("Proto:%s", resp.Proto)
	showDebug(debug, debugLevel)

	debug = fmt.Sprintf("Status Code:%d", resp.StatusCode)
	showDebug(debug, debugLevel)

	debug = fmt.Sprintf("Status Text:%s", http.StatusText(resp.StatusCode))
	showDebug(debug, debugLevel)

	for key, value := range resp.Header {

		switch fmt.Sprintf("%T", value) {

		case "[]string":
			debug = fmt.Sprintf("Header:%v: %s", key, strings.Join(value, " "))

		default:
			debug = fmt.Sprintf("Header:%v: %v", key, value)
		}

		showDebug(debug, debugLevel)

	}

	debug = "Pesponse:* * * * * * END RESPONSE * * * * * * "
	showDebug(debug, debugLevel)

	return
}
