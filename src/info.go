package src

import (
	"fmt"
	"strings"
)

// ShowSystemInfo : View System Information
func ShowSystemInfo() {

	fmt.Print("Creating the information takes a moment...")
	err := buildDatabaseDVR()
	if err != nil {
		ShowError(err, 0)
		return
	}

	buildXEPG(false)

	fmt.Println("OK")
	println()

	fmt.Println(fmt.Sprintf("Version:             %s %s.%s", System.Name, System.Version, System.Build))
	fmt.Println(fmt.Sprintf("Branch:              %s", System.Branch))
	fmt.Println(fmt.Sprintf("GitHub:              %s/%s | Git update = %t", System.GitHub.User, System.GitHub.Repo, System.GitHub.Update))
	fmt.Println(fmt.Sprintf("Folder (config):     %s", System.Folder.Config))

	fmt.Println(fmt.Sprintf("Streams:             %d / %d", len(Data.Streams.Active), len(Data.Streams.All)))
	fmt.Println(fmt.Sprintf("Filter:              %d", len(Data.Filter)))
	fmt.Println(fmt.Sprintf("XEPG Chanels:        %d", int(Data.XEPG.XEPGCount)))

	println()
	fmt.Println(fmt.Sprintf("IPv4 Addresses:"))

	for i, ipv4 := range System.IPAddressesV4 {

		switch count := i; {

		case count < 10:
			fmt.Println(fmt.Sprintf("  %d.                 %s", count, ipv4))
			break
		case count < 100:
			fmt.Println(fmt.Sprintf("  %d.                %s", count, ipv4))
			break

		}

	}

	println()
	fmt.Println(fmt.Sprintf("IPv6 Addresses:"))

	for i, ipv4 := range System.IPAddressesV6 {

		switch count := i; {

		case count < 10:
			fmt.Println(fmt.Sprintf("  %d.                 %s", count, ipv4))
			break
		case count < 100:
			fmt.Println(fmt.Sprintf("  %d.                %s", count, ipv4))
			break

		}

	}

	println("---")

	fmt.Println("Settings [General]")
	fmt.Println(fmt.Sprintf("xTeVe Update:        %t", Settings.XteveAutoUpdate))
	fmt.Println(fmt.Sprintf("UUID:                %s", Settings.UUID))
	fmt.Println(fmt.Sprintf("Tuner (Plex / Emby): %d", Settings.Tuner))
	fmt.Println(fmt.Sprintf("EPG Source:          %s", Settings.EpgSource))

	println("---")

	fmt.Println("Settings [Files]")
	fmt.Println(fmt.Sprintf("Schedule:            %s", strings.Join(Settings.Update, ",")))
	fmt.Println(fmt.Sprintf("Files Update:        %t", Settings.FilesUpdate))
	fmt.Println(fmt.Sprintf("Folder (tmp):        %s", Settings.TempPath))
	fmt.Println(fmt.Sprintf("Image Chaching:      %t", Settings.CacheImages))
	fmt.Println(fmt.Sprintf("Replace EPG Image:   %t", Settings.XepgReplaceMissingImages))

	println("---")

	fmt.Println("Settings [Streaming]")
	fmt.Println(fmt.Sprintf("Buffer:              %s", Settings.Buffer))
	fmt.Println(fmt.Sprintf("UDPxy:               %s", Settings.UDPxy))
	fmt.Println(fmt.Sprintf("Buffer Size:         %d KB", Settings.BufferSize))
	fmt.Println(fmt.Sprintf("Timeout:             %d ms", int(Settings.BufferTimeout)))
	fmt.Println(fmt.Sprintf("User Agent:          %s", Settings.UserAgent))

	println("---")

	fmt.Println("Settings [Backup]")
	fmt.Println(fmt.Sprintf("Folder (backup):     %s", Settings.BackupPath))
	fmt.Println(fmt.Sprintf("Backup Keep:         %d", Settings.BackupKeep))

}
