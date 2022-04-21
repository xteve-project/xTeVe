// Copyright 2019 marmei. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the
// LICENSE file.
// GitHub: https://github.com/xteve-project/xTeVe

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"xteve/src"
)

// GitHubStruct : GitHub Account. The Updates are published via this Account
type GitHubStruct struct {
	Branch string
	Repo   string
	Update bool
	User   string
}

// GitHub : GitHub Account
// If you want to fork this project, enter your Github account here. This prevents a newer version of xTeVe from updating your version.
var GitHub = GitHubStruct{Branch: "master", User: "SCP002", Repo: "xTeVe", Update: false}

// Branch:	GitHub Branch
// User: 	GitHub Username
// Repo: 	GitHub Repository
// Update:	Automatic updates from the GitHub repository [true|false]

// Name : Program Name
const Name = "xTeVe"

// Version : Version, the Build Number is parsed in the main func
const Version = "2.3.0.0000"

// DBVersion : Database Version
const DBVersion = "2.2.2"

// APIVersion : API Version
const APIVersion = "1.1.0"

var homeDirectory = fmt.Sprintf("%s%s.%s%s", src.GetUserHomeDirectory(), string(os.PathSeparator), strings.ToLower(Name), string(os.PathSeparator))
var samplePath = fmt.Sprintf("%spath%sto%sxteve%s", string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator))
var sampleRestore = fmt.Sprintf("%spath%sto%sfile%s", string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator))

var configFolder = flag.String("config", "", ": Config Folder        ["+samplePath+"] (default: "+homeDirectory+")")
var port = flag.String("port", "", ": Server port          [34400] (default: 34400)")
var restore = flag.String("restore", "", ": Restore from backup  ["+sampleRestore+"xteve_backup.zip]")

var gitBranch = flag.String("branch", "", ": Git Branch           [master|beta] (default: master)")
var debug = flag.Int("debug", 0, ": Debug level          [0 - 3] (default: 0)")
var info = flag.Bool("info", false, ": Show system info")
var h = flag.Bool("h", false, ": Show help")

// Activates Development Mode. The local Files are then used for the Webserver.
var dev = flag.Bool("dev", false, ": Activates the developer mode, the source code must be available. The local files for the web interface are used.")

func main() {

	// Separate Build Number from Version Number
	var build = strings.Split(Version, ".")

	var system = &src.System
	system.APIVersion = APIVersion
	system.Branch = GitHub.Branch
	system.Build = build[len(build)-1:][0]
	system.DBVersion = DBVersion
	system.GitHub = GitHub
	system.Name = Name
	system.Version = strings.Join(build[0:len(build)-1], ".")

	// Panic
	defer func() {

		if r := recover(); r != nil {

			fmt.Println()
			fmt.Println("* * * * * FATAL ERROR * * * * *")
			fmt.Println("OS:  ", runtime.GOOS)
			fmt.Println("Arch:", runtime.GOARCH)
			fmt.Println("Err: ", r)
			fmt.Println()

			pc := make([]uintptr, 20)
			runtime.Callers(2, pc)

			for i := range pc {

				if runtime.FuncForPC(pc[i]) != nil {

					f := runtime.FuncForPC(pc[i])
					file, line := f.FileLine(pc[i])

					if string(file)[0:1] != "?" {
						fmt.Printf("%s:%d %s\n", filepath.Base(file), line, f.Name())
					}

				}

			}

			fmt.Println()
			fmt.Println("* * * * * * * * * * * * * * * *")

		}

	}()

	flag.Parse()

	if *h {
		flag.Usage()
		return
	}

	system.Dev = *dev

	// Display System Information
	if *info {

		system.Flag.Info = true

		err := src.Init()
		if err != nil {
			src.ShowError(err, 0)
			os.Exit(0)
		}

		src.ShowSystemInfo()
		return

	}

	// Webserver Port
	if len(*port) > 0 {
		system.Flag.Port = *port
	}

	// Branch
	system.Flag.Branch = *gitBranch
	if len(system.Flag.Branch) > 0 {
		fmt.Println("Git Branch is now:", system.Flag.Branch)
	}

	// Debug Level
	system.Flag.Debug = *debug
	if system.Flag.Debug > 3 {
		flag.Usage()
		return
	}

	// Storage location for the Configuration Files
	if len(*configFolder) > 0 {
		system.Folder.Config = *configFolder
	}

	// Restore Backup
	if len(*restore) > 0 {

		system.Flag.Restore = *restore

		err := src.Init()
		if err != nil {
			src.ShowError(err, 0)
			os.Exit(0)
		}

		err = src.XteveRestoreFromCLI(*restore)
		if err != nil {
			src.ShowError(err, 0)
		}

		os.Exit(0)
	}

	err := src.Init()
	if err != nil {
		src.ShowError(err, 0)
		os.Exit(0)
	}

	err = src.BinaryUpdate()
	if err != nil {
		src.ShowError(err, 0)
	}

	err = src.StartSystem(false)
	if err != nil {
		src.ShowError(err, 0)
		os.Exit(0)
	}

	err = src.InitMaintenance()
	if err != nil {
		src.ShowError(err, 0)
		os.Exit(0)
	}

	err = src.StartWebserver()
	if err != nil {
		src.ShowError(err, 0)
		os.Exit(0)
	}

}
