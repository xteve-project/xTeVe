<div align="center" style="background-color: #111; padding: 100;">
    <a href="https://github.com/senexcrenshaw/xTeVe"><img width="880" height="200" src="html/img/logo_b_880x200.jpg" alt="xTeVe" /></a>
</div>
<br>

# xTeVe

## M3U Proxy and EPG aggregator for Plex DVR and Emby Live TV

### This is a fork of <https://github.com/xteve-project/xTeVe>, all credit goes to the original author

Documentation for setup and configuration is [here](https://github.com/xteve-project/xTeVe-Documentation/blob/main/en/configuration.md).

---

## Features

### Files

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

* See [releases page](https://github.com/senexcrenshaw/xTeVe/releases)

---

## TLS mode

This mode can be enabled by ticking the checkbox in `Settings -> General`.

Unless the server's certificate and it's private key already exists in xTeVe config directory, xTeVe will generate a self-signed automatically.

Self-signed certificate will only allow TLS mode to start up but not to actually establish a secure connections.
For truly working HTTPS, you should [generate](https://gist.github.com/fntlnz/cf14feb5a46b2eda428e000157447309) a certificate by yourself and **also** add the CA certificate to the client-side certificate storage (where the web browser, Plex etc. is).

Certificate and it's private key should be placed in xTeVe config directory like so:

```text
/home/username/.xteve/certificates/xteve.crt
/home/username/.xteve/certificates/xteve.key
```

If the certificate is signed by a certificate authority (CA), it should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.

This will also enable copy to clipboad by clicking the green links at the header. (DVR IP,M3U URL,XEPG URL)

---

## Docker

### Get an image

Pull from dockerhub:

```sh
docker pull senexcrenshaw/xteve:latest
```

**OR** build your own image based on Dockerfile from this repository:

```sh
git clone https://github.com/SenexCrenshaw/xTeVe.git
cd xTeVe
docker build --tag senexcrenshaw/xteve .
```

### Create a container

```sh
docker create \
    --tty \
    --publish 34400:34400 \
    --name xteve \
    senexcrenshaw/xteve
```

With the specific timezone, ip and port:

```sh
docker create \
    --tty \
    --env TZ=Europe/Amsterdam \
    --env XTEVE_PORT=12345 \
    --publish 192.168.88.218:12345:12345 \
    --name xteve \
    senexcrenshaw/xteve
```

### Start a container

```sh
docker start xteve
```

#### Attach to a started container

```sh
docker attach xteve
```

To detach from a container, press `Ctrl + C`.

#### Access web UI

Open `http(s)://<ip>:<port>/web/` in browser, for example:
`http://192.168.88.218:34400/web/`

#### Stop a running container

```sh
docker stop xteve
```

---

### xTeVe Beta branch

New features and bug fixes are only available in beta branch. Only after successful testing are they are merged into the main branch.

**It is not recommended to use the beta version in a production system.**  

With the command line argument `branch` the Git Branch can be changed. xTeVe must be started via the terminal.  

#### Switch from main to beta branch

```text
xteve -branch beta

...
[xTeVe] GitHub:                https://github.com/senexcrenshaw
[xTeVe] Git Branch:            beta [senexcrenshaw]
...
```

#### Switch from beta to main branch

```text
xteve -branch main

...
[xTeVe] GitHub:                https://github.com/senexcrenshaw
[xTeVe] Git Branch:            main [senexcrenshaw]
...
```

When the branch is changed, an update is only performed if there is a new version and the update function is activated in the settings.  

---

## Build from source code [Go / Golang]

### Requirements

* [Go](https://golang.org) (go1.18 or newer)

### Dependencies

* [avfs](https://github.com/avfs/avfs)
* [go-ssdp](https://github.com/koron/go-ssdp)
* [lo](https://github.com/samber/lo)
* [osext](https://github.com/kardianos/osext)
* [testify](https://github.com/stretchr/testify)
* [websocket](https://github.com/gorilla/websocket)

### Build

#### 1. Download source code

```sh
git clone https://github.com/senexcrenshaw/xTeVe.git
```

#### 2. Install dependencies

```sh
go mod tidy
```

Or

```sh
go get github.com/avfs/avfs@latest 
go get github.com/gorilla/websocket
go get github.com/kardianos/osext
go get github.com/koron/go-ssdp
go get github.com/samber/lo
go get github.com/stretchr/testify
```

#### 3. Update dependencies (optional)

```sh
go get -u ./...
```

#### 5. Update web files (optional)

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

`xteve -branch main` or `xteve -branch beta`

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

## Forks

When creating a fork, the xTeVe GitHub account must be changed from the source code or the update function disabled.

xteve.go - Line: 29

```go
var GitHub = GitHubStruct{Branch: "main", User: "senexcrenshaw", Repo: "xTeVe", Update: true}

// Branch: GitHub Branch
// User:   GitHub Username
// Repo:   GitHub Repository
// Update: Automatic updates from the GitHub repository [true|false]
```
