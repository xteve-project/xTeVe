class SettingsCategory {
  DocumentID:string = "content_settings"
  createCategoryHeadline(value:string):any {
    var element = document.createElement("H4")
    element.innerHTML = value
    return element
  }

  createHR():any {
    var element = document.createElement("HR")
    return element
  }

  createSettings(settingsKey:string):any {
    var setting = document.createElement("TR")
    var content:PopupContent = new PopupContent()
    var data = SERVER["settings"][settingsKey]

    switch (settingsKey) {

      // Text inputs
      case "update":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.update.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createInput("text", "update", data.toString())
        input.setAttribute("placeholder", "{{.settings.update.placeholder}}")
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "backup.path":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.backupPath.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createInput("text", "backup.path", data)
        input.setAttribute("placeholder", "{{.settings.backupPath.placeholder}}")
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "temp.path":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.tempPath.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createInput("text", "temp.path", data)
        input.setAttribute("placeholder", "{{.settings.tmpPath.placeholder}}")
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "user.agent":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.userAgent.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createInput("text", "user.agent", data)
        input.setAttribute("placeholder", "{{.settings.userAgent.placeholder}}")
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "buffer.timeout":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.bufferTimeout.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createInput("text", "buffer.timeout", data)
        input.setAttribute("placeholder", "{{.settings.bufferTimeout.placeholder}}")
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "ffmpeg.path":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.ffmpegPath.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createInput("text", "ffmpeg.path", data)
        input.setAttribute("placeholder", "{{.settings.ffmpegPath.placeholder}}")
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "ffmpeg.options":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.ffmpegOptions.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createInput("text", "ffmpeg.options", data)
        input.setAttribute("placeholder", "{{.settings.ffmpegOptions.placeholder}}")
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "vlc.path":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.vlcPath.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createInput("text", "vlc.path", data)
        input.setAttribute("placeholder", "{{.settings.vlcPath.placeholder}}")
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "vlc.options":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.vlcOptions.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createInput("text", "vlc.options", data)
        input.setAttribute("placeholder", "{{.settings.vlcOptions.placeholder}}")
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      // Checkboxes
      case "tlsMode":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.tlsMode.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "disallowURLDuplicates":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.disallowURLDuplicates.title}}" + ":"
  
        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)
  
        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "authentication.web":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.authenticationWEB.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "authentication.pms":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.authenticationPMS.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "authentication.m3u":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.authenticationM3U.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "authentication.xml":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.authenticationXML.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "authentication.api":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.authenticationAPI.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "files.update":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.filesUpdate.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "cache.images":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.cacheImages.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "xepg.replace.missing.images":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.replaceEmptyImages.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "storeBufferInRAM":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.storeBufferInRAM.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "xteveAutoUpdate":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.xteveAutoUpdate.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "api":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.api.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "enableMappedChannels":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.enableMappedChannels.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createCheckbox(settingsKey)
        input.checked = data
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      // Select
      case "hostIP":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.hostIP.title}}" + ":"

        var tdRight = document.createElement("TD")
        var text: any[] = SERVER["ipAddressesV4Host"]
        var values: any[] = SERVER["ipAddressesV4Host"]

        var select = content.createSelect(text, values, data, settingsKey)
        select.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(select)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "tuner":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.tuner.title}}" + ":"

        var tdRight = document.createElement("TD")
        var text = new Array()
        var values = new Array()

        for (var i = 1; i <= 100; i++) {
          text.push(i)
          values.push(i)
        }

        var select = content.createSelect(text, values, data, settingsKey)
        select.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(select)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "epgSource":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.epgSource.title}}" + ":"

        var tdRight = document.createElement("TD")
        var text:any[] = ["PMS", "XEPG"]
        var values:any[] = ["PMS", "XEPG"]

        var select = content.createSelect(text, values, data, settingsKey)
        select.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(select)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "defaultMissingEPG":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.defaultMissingEPG.title}}" + ":"

        var tdRight = document.createElement("TD")
        var text:any[] = [
          "-", "30 Minutes (30_Minutes)", "60 Minutes (60_Minutes)", "90 Minutes (90_Minutes)",
          "120 Minutes (120_Minutes)", "180 Minutes (180_Minutes)", "240 Minutes (240_Minutes)", "360 Minutes (360_Minutes)"
        ]
        var values:any[] = [
          "-", "30_Minutes", "60_Minutes", "90_Minutes", "120_Minutes", "180_Minutes", "240_Minutes", "360_Minutes"
        ]

        var select = content.createSelect(text, values, data, settingsKey)
        select.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(select)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "backup.keep":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.backupKeep.title}}" + ":"

        var tdRight = document.createElement("TD")
        var text:any[] = ["5", "10", "20", "30", "40", "50"]
        var values:any[] = ["5", "10", "20", "30", "40", "50"]

        var select = content.createSelect(text, values, data, settingsKey)
        select.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(select)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

      case "buffer.size.kb":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.bufferSize.title}}" + ":"

        var tdRight = document.createElement("TD")
        var text:any[] = ["0.5 MB", "1 MB", "2 MB", "3 MB", "4 MB", "5 MB", "6 MB", "7 MB", "8 MB"]
        var values:any[] = ["512", "1024", "2048", "3072", "4096", "5120", "6144", "7168", "8192"]

        var select = content.createSelect(text, values, data, settingsKey)
        select.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(select)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

       case "buffer":
        var tdLeft = document.createElement("TD")
        tdLeft.innerHTML = "{{.settings.streamBuffering.title}}" + ":"

        var tdRight = document.createElement("TD")
        var text:any[] = ["{{.settings.streamBuffering.info_false}}", "xTeVe: ({{.settings.streamBuffering.info_xteve}})", "FFmpeg: ({{.settings.streamBuffering.info_ffmpeg}})", "VLC: ({{.settings.streamBuffering.info_vlc}})"]
        var values:any[] = ["-", "xteve", "ffmpeg", "vlc"]

        var select = content.createSelect(text, values, data, settingsKey)
        select.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(select)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

    case "udpxy":

        var tdLeft = document.createElement("TD");
        tdLeft.innerHTML = "{{.settings.udpxy.title}}" + ":"

        var tdRight = document.createElement("TD")
        var input = content.createInput("text", "udpxy", data)
        input.setAttribute("placeholder", "{{.settings.udpxy.placeholder}}")
        input.setAttribute("onchange", "javascript: this.className = 'changed'")
        tdRight.appendChild(input)

        setting.appendChild(tdLeft)
        setting.appendChild(tdRight)
        break

    }

    return setting

  }


  createDescription(settingsKey:string):any {

    var description = document.createElement("TR")
    var text:string
    switch (settingsKey) {

      case "tlsMode":
        text = "{{.settings.tlsMode.description}}"
        break

      case "disallowURLDuplicates":
        text = "{{.settings.disallowURLDuplicates.description}}"
        break
  
      case "authentication.web":
        text = "{{.settings.authenticationWEB.description}}"
        break

      case "authentication.m3u":
        text = "{{.settings.authenticationM3U.description}}"
        break

      case "authentication.pms":
        text = "{{.settings.authenticationPMS.description}}"
        break

      case "authentication.xml":
        text = "{{.settings.authenticationXML.description}}"
        break

      case "authentication.api":
        if (SERVER["settings"]["authentication.web"] == true) {
          text = "{{.settings.authenticationAPI.description}}"
        }
        break

      case "xteveAutoUpdate":
        text = "{{.settings.xteveAutoUpdate.description}}"
        break

      case "backup.keep":
        text = "{{.settings.backupKeep.description}}"
        break

      case "backup.path":
        text = "{{.settings.backupPath.description}}"
        break

      case "temp.path":
        text = "{{.settings.tempPath.description}}"
        break

      case "buffer":
        text = "{{.settings.streamBuffering.description}}"
        break

      case "buffer.size.kb":
        text = "{{.settings.bufferSize.description}}"
        break

      case "storeBufferInRAM":
        text = "{{.settings.storeBufferInRAM.description}}"
        break

      case "buffer.timeout":
        text = "{{.settings.bufferTimeout.description}}"
        break

      case "user.agent":
        text = "{{.settings.userAgent.description}}"
        break

      case "ffmpeg.path":
        text = "{{.settings.ffmpegPath.description}}"
        break

      case "ffmpeg.options":
        text = "{{.settings.ffmpegOptions.description}}"
        break

      case "vlc.path":
        text = "{{.settings.vlcPath.description}}"
        break

      case "vlc.options":
        text = "{{.settings.vlcOptions.description}}"
        break

      case "epgSource":
        text = "{{.settings.epgSource.description}}"
        break

      case "hostIP":
        text = "{{.settings.hostIP.description}}"
        break

      case "tuner":
        text = "{{.settings.tuner.description}}"
        break

      case "update":
        text = "{{.settings.update.description}}"
        break

      case "api":
        text = "{{.settings.api.description}}"
        break

      case "defaultMissingEPG":
        text = "{{.settings.defaultMissingEPG.description}}"
        break

      case "enableMappedChannels":
        text = "{{.settings.enableMappedChannels.description}}"
        break

      case "files.update":
        text = "{{.settings.filesUpdate.description}}"
        break

      case "cache.images":
        text = "{{.settings.cacheImages.description}}"
        break

      case "xepg.replace.missing.images":
        text = "{{.settings.replaceEmptyImages.description}}"
        break

      case "udpxy":
        text = "{{.settings.udpxy.description}}"
        break

      default:
        text = ""
        break

    }

    var tdLeft = document.createElement("TD")
    tdLeft.innerHTML = ""

    var tdRight = document.createElement("TD")
    var pre = document.createElement("PRE")
    pre.innerHTML = text
    tdRight.appendChild(pre)

    description.appendChild(tdLeft)
    description.appendChild(tdRight)

    return description

  }

}

class SettingsCategoryItem extends SettingsCategory {
  headline:string
  settingsKeys:string

  constructor(headline:string, settingsKeys:string) {
    super()
    this.headline = headline
    this.settingsKeys = settingsKeys
  }

  createCategory():void {
    var headline = this.createCategoryHeadline(this.headline)
    var settingsKeys = this.settingsKeys

    var doc = document.getElementById(this.DocumentID)
    doc.appendChild(headline)

    // Create a table for the category

    var table = document.createElement("TABLE")

    var keys = settingsKeys.split(",")

    keys.forEach(settingsKey => {

      switch (settingsKey) {

        case "authentication.pms":
        case "authentication.m3u":
        case "authentication.xml":
        case "authentication.api":
          if (SERVER["settings"]["authentication.web"] == false) {
            break
          }

        default:
          var item = this.createSettings(settingsKey)
          var description = this.createDescription(settingsKey)

          table.appendChild(item)
          table.appendChild(description)
          break

      }

    });

    doc.appendChild(table)
    doc.appendChild(this.createHR())
  }

}

function showSettings() {

  for (let i = 0; i < settingsCategory.length; i++) {
    settingsCategory[i].createCategory()
  }

}

function saveSettings() {

  var cmd = "saveSettings"
  var div = document.getElementById("content_settings")
  var settings = div.getElementsByClassName("changed")

  var newSettings = new Object();

  for (let i = 0; i < settings.length; i++) {

    var name:string
    var value:any

    switch (settings[i].tagName) {
      case "INPUT":

        switch ((settings[i] as HTMLInputElement).type) {
          case "checkbox":
            name = (settings[i] as HTMLInputElement).name
            value = (settings[i] as HTMLInputElement).checked
            newSettings[name] = value
            break

          case "text":
            name = (settings[i] as HTMLInputElement).name
            value = (settings[i] as HTMLInputElement).value

            switch (name) {
              case "update":
                value = value.split(",")
                value = value.filter(function(e:any) { return e})
                break

              case "buffer.timeout":
                value = parseFloat(value)

            }

            newSettings[name] = value
            break
        }

        break

      case "SELECT":
        name = (settings[i] as HTMLSelectElement).name
        value = (settings[i] as HTMLSelectElement).value

        // If the value is a number, store it as a number
        if(isNaN(value)){
          newSettings[name] = value
        } else {
          newSettings[name] = parseInt(value)
        }

        break

    }

  }

  var data = new Object()
  data["settings"] = newSettings

  var server:Server = new Server(cmd)
  server.request(data)

}
