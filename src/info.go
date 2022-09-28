package src

import (
	"fmt"
	"strings"
)

// ShowSystemVersion basic version info
func ShowSystemVersion() {
  fmt.Println(fmt.Sprintf("%s.%s", System.Version, System.Build))
}

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

	fmt.Printf("Version:             %s %s.%s\n", System.Name, System.Version, System.Build)
	fmt.Printf("Branch:              %s\n", System.Branch)
	fmt.Printf("GitHub:              %s/%s | Git update = %t\n", System.GitHub.User, System.GitHub.Repo, System.GitHub.Update)
	fmt.Printf("Folder (config):     %s\n", System.Folder.Config)

	fmt.Printf("Streams:             %d / %d\n", len(Data.Streams.Active), len(Data.Streams.All))
	fmt.Printf("Filter:              %d\n", len(Data.Filter))
	fmt.Printf("XEPG Chanels:        %d\n", int(Data.XEPG.XEPGCount))

	println()
	fmt.Println("IPv4 Addresses:")

	for i, ipv4 := range System.IPAddressesV4 {

		switch count := i; {

		case count < 10:
			fmt.Printf("  %d.                 %s\n", count, ipv4)
		case count < 100:
			fmt.Printf("  %d.                %s\n", count, ipv4)

		}

	}

	println()
	fmt.Println("IPv6 Addresses:")

	for i, ipv4 := range System.IPAddressesV6 {

		switch count := i; {

		case count < 10:
			fmt.Printf("  %d.                 %s\n", count, ipv4)
		case count < 100:
			fmt.Printf("  %d.                %s\n", count, ipv4)
		}

	}

	println("---")

	fmt.Println("Settings [General]")
	fmt.Printf("xTeVe Update:        %t\n", Settings.XteveAutoUpdate)
	fmt.Printf("UUID:                %s\n", Settings.UUID)
	fmt.Printf("Tuner (Plex / Emby): %d\n", Settings.Tuner)
	fmt.Printf("EPG Source:          %s\n", Settings.EpgSource)

	println("---")

	fmt.Println("Settings [Files]")
	fmt.Printf("Schedule:            %s\n", strings.Join(Settings.Update, ","))
	fmt.Printf("Files Update:        %t\n", Settings.FilesUpdate)
	fmt.Printf("Folder (tmp):        %s\n", Settings.TempPath)
	fmt.Printf("Image Chaching:      %t\n", Settings.CacheImages)
	fmt.Printf("Replace EPG Image:   %t\n", Settings.XepgReplaceMissingImages)

	println("---")

	fmt.Println("Settings [Streaming]")
	fmt.Printf("Buffer:              %s\n", Settings.Buffer)
	fmt.Printf("UDPxy:               %s\n", Settings.UDPxy)
	fmt.Printf("Buffer Size:         %d KB\n", Settings.BufferSize)
	fmt.Printf("Timeout:             %d ms\n", int(Settings.BufferTimeout))
	fmt.Printf("User Agent:          %s\n", Settings.UserAgent)

	println("---")

	fmt.Println("Settings [Backup]")
	fmt.Printf("Folder (backup):     %s\n", Settings.BackupPath)
	fmt.Printf("Backup Keep:         %d\n", Settings.BackupKeep)

}
