package src

import "time"

// Playlist : Enthält allen Playlistinformationen, die der Buffer benötigr
type Playlist struct {
	Folder       string
	PlaylistID   string
	PlaylistName string
	Tuner        int

	Clients map[int]ThisClient
	Streams map[int]ThisStream
}

// ThisClient : Clientinfos
type ThisClient struct {
	Connection int
}

// ThisStream : Enthält Informationen zu dem abzuspielenden Stream einer Playlist
type ThisStream struct {
	ChannelName      string
	Error            string
	Folder           string
	MD5              string
	NetworkBandwidth int
	PlaylistID       string
	PlaylistName     string
	Status           bool
	URL              string

	Segment []Segment

	// Serverinformationen
	Location           string
	URLFile            string
	URLHost            string
	URLPath            string
	URLRedirect        string
	URLScheme          string
	URLStreamingServer string

	// Wird nur für HLS / M3U8 verwendet
	Body             string
	Difference       float64
	Duration         float64
	DynamicBandwidth bool
	FirstSequence    int64
	HLS              bool
	LastSequence     int64
	M3U8URL          string
	NewSegCount      int
	OldSegCount      int
	Sequence         int64
	TimeDiff         float64
	TimeEnd          time.Time
	TimeSegDuration  float64
	TimeStart        time.Time
	Version          int
	Wait             float64

	DynamicStream map[int]DynamicStream

	// Lokale Temp Datein
	OldSegments []string
}

// Segment : URL Segmente (HLS / M3U8)
type Segment struct {
	Duration     float64
	Info         bool
	PlaylistType string
	Sequence     int64
	URL          string
	Version      int
	Wait         float64

	StreamInf struct {
		AverageBandwidth int
		Bandwidth        int
		Framerate        float64
		Resolution       string
		SegmentURL       string
	}
}

// DynamicStream : Streaminformationen bei dynamischer Bandbreite
type DynamicStream struct {
	AverageBandwidth int
	Bandwidth        int
	Framerate        float64
	Resolution       string
	URL              string
}

// ClientConnection : Client Verbindungen
type ClientConnection struct {
	Connection int
	Error      error
}

// BandwidthCalculation : Bandbreitenberechnung für den Stream
type BandwidthCalculation struct {
	NetworkBandwidth int
	Size             int
	Start            time.Time
	Stop             time.Time
	TimeDiff         float64
}

/*
var args = "-hide_banner -loglevel panic -re -i " + url + " -codec copy -f mpegts pipe:1"
		//var args = "-re -i " + url + " -codec copy -f mpegts pipe:1"
		cmd := exec.Command("/usr/local/bin/ffmpeg", strings.Split(args, " ")...)

		//run := exec.Command("/usr/local/bin/ffmpeg", "-hide_banner", "-loglevel", "panic", "-re", "-i", url, "-codec", "copy", "-f", "mpegts", "pipe:1")
		//run := exec.Command("/usr/local/bin/ffmpeg", "-re", "-i", url, "-codec", "copy", "-f", "mpegts", "pipe:1")

		stderr, _ := cmd.StderrPipe()
		cmd.Start()

		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}
		cmd.Wait()

		os.Exit(0)
*/

/*

ffmpegOut, _ := run.StderrPipe()
		//run.Start()

		scanner = bufio.NewScanner(ffmpegOut)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}

		ffmpegOut, err = run.StdoutPipe()
		if err != nil {
			ShowError(err, 0)
			return
		}

		stderr, stderrErr := run.StderrPipe()
		if stderrErr != nil {
			fmt.Println(stderrErr)
		}

		_ = stderr

		if startErr := run.Start(); startErr != nil {
			fmt.Println(startErr)

			return
		}

		n, err := ffmpegOut.Read(buffer)
		_ = n
		_ = stream
		_ = fileSize

		if err != nil && err != io.EOF {

			ShowError(err, 0)
			addErrorToStream(err)
			return

		}

		defer bufferFile.Close()

		scanner = bufio.NewScanner(ffmpegOut)

		for scanner.Scan() {
			//fmt.Printf("%s\n", scanner.Text())
			//fmt.Println(scanner)
			thisLine := scanner.Text()
			line := make([]byte, len(thisLine))

			buffer = append(buffer, line...)

			fmt.Println(len(buffer))

			if len(buffer) > tmpFileSize {

				if _, err := bufferFile.Write(buffer[:]); err != nil {

					ShowError(err, 0)
					addErrorToStream(err)
					run.Process.Kill()
					return

				}

				buffer = make([]byte, 1024*Settings.BufferSize*2)

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
					run.Process.Kill()

					err = os.RemoveAll(stream.Folder)
					if err != nil {
						ShowError(err, 4005)
					}

					return
				}

				bufferFile, err = os.Create(tmpFile)
				if err != nil {
					addErrorToStream(err)
					run.Process.Kill()
					return
				}

				fileSize = 0

				if n == 0 {
					bufferFile.Close()
					run.Process.Kill()
					break
				}

				os.Exit(0)

			}



		}

*/
