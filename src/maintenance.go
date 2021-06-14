package src

import (
	"fmt"
	"math/rand"
	"time"
)

// InitMaintenance : Initialize maintenance process
func InitMaintenance() (err error) {

	rand.Seed(time.Now().Unix())
	System.TimeForAutoUpdate = fmt.Sprintf("0%d%d", randomTime(0, 2), randomTime(10, 59))

	go maintenance()

	return
}

func maintenance() {

	for {

		var t = time.Now()

		// Update the playlist and XMLTV files
		if System.ScanInProgress == 0 {

			for _, schedule := range Settings.Update {

				if schedule == t.Format("1504") {

					showInfo("Update:" + schedule)

					// Create a backup
					err := xTeVeAutoBackup()
					if err != nil {
						ShowError(err, 000)
					}

					// Update Playlist and XMLTV Files
					getProviderData("m3u", "")
					getProviderData("hdhr", "")

					if Settings.EpgSource == "XEPG" {
						getProviderData("xmltv", "")
					}

					// Create database for DVR
					err = buildDatabaseDVR()
					if err != nil {
						ShowError(err, 000)
					}

					if Settings.CacheImages == false && System.ImageCachingInProgress == 0 {
						removeChildItems(System.Folder.ImagesCache)
					}

					// Create XEPG Files
					Data.Cache.XMLTV = make(map[string]XMLTV)
					buildXEPG(false)

				}

			}

			// Update xTeVe (Binary)
			if System.TimeForAutoUpdate == t.Format("1504") {
				BinaryUpdate()
			}

		}

		time.Sleep(60 * time.Second)

	}

	return
}

func randomTime(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
