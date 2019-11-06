<div align="center" style="background-color: #111; padding: 100;">
    <a href="https://github.com/xteve-project/xTeVe"><img width="880" height="200" src="html/img/logo_b_880x200.jpg" alt="xTeVe" /></a>
</div>
<br>

# xTeVe
## M3U Proxy for Plex DVR and Emby Live TV.  

Documentation for setup and configuration is [here](https://github.com/xteve-project/xTeVe-Documentation/blob/master/en/configuration.md).

#### Donation
* **Bitcoin:** 1c1iCe4CJPfNUXtqxKBbW2Qd2EtqRPWme  
![Bitcoin](html/img/BC-QR.jpg "Bitcoin - xTeVe")

## Requirements
### Plex
* Plex Media Server (1.11.1.4730 or newer)
* Plex Client with DVR support
* Plex Pass

### Emby
* Emby Server (3.5.3.0 or newer)
* Emby Client with Live-TV support
* Emby Premiere

--- 

## Features

#### Files
* Merge external M3U files
* Merge external XMLTV files
* Automatic M3U and XMLTV update
* M3U and XMLTV export

#### Channel management
* Filtering streams
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

## Downloads v2 | 64 Bit only
#### 64 Bit Intel / AMD

* [Windows](https://github.com/xteve-project/xTeVe-Downloads/blob/master/xteve_windows_amd64.zip?raw=true)
* [OS X](https://github.com/xteve-project/xTeVe-Downloads/blob/master/xteve_darwin_amd64.zip?raw=true)
* [Linux](https://github.com/xteve-project/xTeVe-Downloads/blob/master/xteve_linux_amd64.zip?raw=true)
* [FreeBSD](https://github.com/xteve-project/xTeVe-Downloads/blob/master/xteve_freebsd_amd64.zip?raw=true)

#### 64 Bit ARM
* [Linux](https://github.com/xteve-project/xTeVe-Downloads/blob/master/xteve_linux_arm64.zip?raw=true)

#### Recommended Docker Image (Linux 64 Bit)
Thanks to @alturismo and @LeeD for creating the Docker Images.

**Created by alturismo:**  
[xTeVe](https://hub.docker.com/r/alturismo/xteve)  
[xTeVe / Guide2go](https://hub.docker.com/r/alturismo/xteve_guide2go)  
[xTeVe / Guide2go / owi2plex](https://hub.docker.com/r/alturismo/xteve_g2g_owi)

Including:  
- Guide2go: XMLTV grabber for Schedules Direct  
- owi2plex: XMLTV file grabber for Enigma receivers

**Created by LeeD:**  
[xTeVe / Guide2go / Zap2XML](https://hub.docker.com/r/dnsforge/xteve)  

Including:  
- Guide2go: XMLTV grabber for Schedules Direct  
- Zap2XML: Perl based zap2it XMLTV grabber  
- Bash: A Unix / Linux shell  
- Crond: Daemon to execute scheduled commands  
- Perl: Programming language   

---

### xTeVe Beta branch
New features and bug fixes are only available in beta branch. Only after successful testing are they are merged into the master branch.

**It is not recommended to use the beta version in a production system.**  

With the command line argument `branch` the Git Branch can be changed. xTeVe must be started via the terminal.  

#### Switch from master to beta branch:
```
xteve -branch beta

...
[xTeVe] GitHub:                https://github.com/xteve-project
[xTeVe] Git Branch:            beta [xteve-project]
...
```

#### Switch from beta to master branch:
```
xteve -branch master

...
[xTeVe] GitHub:                https://github.com/xteve-project
[xTeVe] Git Branch:            master [xteve-project]
...
```

When the branch is changed, an update is only performed if there is a new version and the update function is activated in the settings.  

---

## Build from source code [Go / Golang]

#### Requirements
* [Go](https://golang.org) (go1.12.4 or newer)

#### Dependencies
* [go-ssdp](https://github.com/koron/go-ssdp)
* [websocket](https://github.com/gorilla/websocket)
* [osext](https://github.com/kardianos/osext)

#### Build
1. Download source code
2. Install dependencies
```
go get github.com/koron/go-ssdp
go get github.com/gorilla/websocket
go get github.com/kardianos/osext
```
3. Build xTeVe
```
go build xteve.go
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


