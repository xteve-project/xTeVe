<div align="center" style="background-color: #111; padding: 100;">
    <a href="https://github.com/SCP002/xTeVe"><img width="880" height="200" src="html/img/logo_b_880x200.jpg" alt="xTeVe" /></a>
</div>
<br>

# xTeVe
## M3U Proxy and EPG aggregator for Plex DVR and Emby Live TV.

#### This is a fork of <https://github.com/xteve-project/xTeVe>, all credit goes to the original author.

Documentation for setup and configuration is [here](https://github.com/xteve-project/xTeVe-Documentation/blob/master/en/configuration.md).

#### Donation
* **Bitcoin:** 1c1iCe4CJPfNUXtqxKBbW2Qd2EtqRPWme  
![Bitcoin](html/img/BC-QR.jpg "Bitcoin - xTeVe")

--- 

## Features

#### Files
* Merge external M3U files
* Merge external XMLTV files (EPG aggregation)
* Automatic M3U and XMLTV update
* M3U and XMLTV export

#### Channel management
* Filtering streams
* Teleguide timeshift
* Channel mapping
* Channel order
* Channel logos
* Channel categories

#### Streaming
* Buffer with HLS / M3U8 support
* Re-streaming
* Number of tuners adjustable
* Compatible with Plex / Emby EPG

---

## Downloads
* See [releases page](https://github.com/SCP002/xTeVe/releases)

---

### xTeVe Beta branch
New features and bug fixes are only available in beta branch. Only after successful testing are they are merged into the master branch.

**It is not recommended to use the beta version in a production system.**  

With the command line argument `branch` the Git Branch can be changed. xTeVe must be started via the terminal.  

#### Switch from master to beta branch:
```
xteve -branch beta

...
[xTeVe] GitHub:                https://github.com/SCP002
[xTeVe] Git Branch:            beta [SCP002]
...
```

#### Switch from beta to master branch:
```
xteve -branch master

...
[xTeVe] GitHub:                https://github.com/SCP002
[xTeVe] Git Branch:            master [SCP002]
...
```

When the branch is changed, an update is only performed if there is a new version and the update function is activated in the settings.  

---

## Run

#### Requirements

---

## Build from source code [Go / Golang]

#### Requirements
* [Go](https://golang.org) (go1.18 or newer)

#### Dependencies
* [go-ssdp](https://github.com/koron/go-ssdp)
* [websocket](https://github.com/gorilla/websocket)
* [osext](https://github.com/kardianos/osext)

#### Build

#### 1. Download source code

#### 2. Install dependencies

```sh
go mod tidy
```

Or

```sh
go get github.com/koron/go-ssdp
go get github.com/gorilla/websocket
go get github.com/kardianos/osext
```

#### 3. Update dependencies (optional)

```sh
go get -u ./...
```

#### 4. Update web files (optional)

If TypeScript files were changed, run:

```sh
tsc -p ./ts/tsconfig.json
```

Then, to embed updated JavaScript files into the source code (src/webUI.go), run it in development mode at least once:

```sh
go build xteve.go
xteve -dev
```

:exclamation: To not to get CreateFile error, do not forget to switch your binary to "regular" mode after runnning with `-dev` flag:

`xteve -branch master` or `xteve -branch beta`

#### 4. Build xTeVe

```sh
go build xteve.go
```

Or use convenient cross-compile tool. To build binaries for every OS / architecture pair into `./xteve-build/` folder:

```sh
go get github.com/mitchellh/gox
go install github.com/mitchellh/gox
gox -output="./xteve-build/{{.Dir}}_{{.OS}}_{{.Arch}}" ./
```

---

## Fork without pull request :mega:
When creating a fork, the xTeVe GitHub account must be changed from the source code or the update function disabled.
Future updates of the xteve-project would update your fork. :wink:

xteve.go - Line: 29
```Go
var GitHub = GitHubStruct{Branch: "master", User: "xteve-project", Repo: "xTeVe-Downloads", Update: true}

/*
  Branch: GitHub Branch
  User:   GitHub Username
  Repo:   GitHub Repository
  Update: Automatic updates from the GitHub repository [true|false]
*/

```
