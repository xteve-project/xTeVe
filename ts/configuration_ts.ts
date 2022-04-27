class WizardCategory {
  DocumentID = "content"

  createCategoryHeadline(value:string):any {
    var element = document.createElement("H4")
    element.innerHTML = value
    return element
  }
}

class WizardItem extends WizardCategory {
  key:string
  headline:string

  constructor(key:string, headline:string) {
    super()
    this.headline = headline
    this.key = key
  }

  createWizard():void {
    var headline = this.createCategoryHeadline(this.headline)
    var key = this.key
    var content:PopupContent = new PopupContent()
    var description:string

    var doc = document.getElementById(this.DocumentID)
    doc.innerHTML = ""
    doc.appendChild(headline)

    switch (key) {
      case "tuner":
        var text = new Array()
        var values = new Array()

        for (var i = 1; i <= 100; i++) {
          text.push(i)
          values.push(i)
        }

        var select = content.createSelect(text, values, "1", key)
        select.setAttribute("class", "wizard")
        select.id = key
        doc.appendChild(select)

        description = "{{.wizard.tuner.description}}"

        break;
      
      case "epgSource":
        var text:any[] = ["PMS", "XEPG"]
        var values:any[] = ["PMS", "XEPG"]

        var select = content.createSelect(text, values, "XEPG", key)
        select.setAttribute("class", "wizard")
        select.id = key
        doc.appendChild(select)

        description = "{{.wizard.epgSource.description}}"

        break

      case "m3u":
        var input = content.createInput("text", key, "")
        input.setAttribute("placeholder", "{{.wizard.m3u.placeholder}}")
        input.setAttribute("class", "wizard")
        input.id = key
        doc.appendChild(input)

        description = "{{.wizard.m3u.description}}"

        break

      case "xmltv":
        var input = content.createInput("text", key, "")
        input.setAttribute("placeholder", "{{.wizard.xmltv.placeholder}}")
        input.setAttribute("class", "wizard")
        input.id = key
        doc.appendChild(input)

        description = "{{.wizard.xmltv.description}}"

      break

      default:
        break;
    }

    var pre = document.createElement("PRE")
    pre.innerHTML = description
    doc.appendChild(pre)
  }


}


function readyForConfiguration(wizard:number) {

  var server:Server = new Server("getServerConfig")
  server.request(new Object())

  showElement("loading", false)

  configurationWizard[wizard].createWizard()

}

function saveWizard() {

  var cmd = "saveWizard"
  var div = document.getElementById("content")
  var config = div.getElementsByClassName("wizard")

  var wizard = new Object()

  for (var i = 0; i < config.length; i++) {

    var name:string
    var value:any
    
    switch (config[i].tagName) {
      case "SELECT":
        name = (config[i] as HTMLSelectElement).name
        value = (config[i] as HTMLSelectElement).value

        // If the value is a number, store it as a number
        if(isNaN(value)){
          wizard[name] = value
        } else {
          wizard[name] = parseInt(value)
        }

        break

      case "INPUT":
        switch ((config[i] as HTMLInputElement).type) {
          case "text":
            name = (config[i] as HTMLInputElement).name
            value = (config[i] as HTMLInputElement).value

            if (value.length == 0) {
              var msg = name.toUpperCase() + ": " + "{{.alert.missingInput}}"
              alert(msg)
              return
            }

            wizard[name] = value
            break
        }
        break
      
      default:
        // code...
        break;
    }

  }

  var data = new Object()
  data["wizard"] = wizard

  var server:Server = new Server(cmd)
  server.request(data)

}

// Wizard
var configurationWizard = new Array()
configurationWizard.push(new WizardItem("tuner", "{{.wizard.tuner.title}}"))
configurationWizard.push(new WizardItem("epgSource", "{{.wizard.epgSource.title}}"))
configurationWizard.push(new WizardItem("m3u", "{{.wizard.m3u.title}}"))
configurationWizard.push(new WizardItem("xmltv", "{{.wizard.xmltv.title}}"))