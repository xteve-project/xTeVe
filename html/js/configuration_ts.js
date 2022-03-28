class WizardCategory {
    constructor() {
        this.DocumentID = "content";
    }
    createCategoryHeadline(value) {
        var element = document.createElement("H4");
        element.innerHTML = value;
        return element;
    }
}
class WizardItem extends WizardCategory {
    constructor(key, headline) {
        super();
        this.headline = headline;
        this.key = key;
    }
    createWizard() {
        var headline = this.createCategoryHeadline(this.headline);
        var key = this.key;
        var content = new PopupContent();
        var description;
        var doc = document.getElementById(this.DocumentID);
        doc.innerHTML = "";
        doc.appendChild(headline);
        switch (key) {
            case "tuner":
                var text = new Array();
                var values = new Array();
                for (var i = 1; i <= 100; i++) {
                    text.push(i);
                    values.push(i);
                }
                var select = content.createSelect(text, values, "1", key);
                select.setAttribute("class", "wizard");
                select.id = key;
                doc.appendChild(select);
                description = "{{.wizard.tuner.description}}";
                break;
            case "epgSource":
                var text = ["PMS", "XEPG"];
                var values = ["PMS", "XEPG"];
                var select = content.createSelect(text, values, "XEPG", key);
                select.setAttribute("class", "wizard");
                select.id = key;
                doc.appendChild(select);
                description = "{{.wizard.epgSource.description}}";
                break;
            case "m3u":
                var input = content.createInput("text", key, "");
                input.setAttribute("placeholder", "{{.wizard.m3u.placeholder}}");
                input.setAttribute("class", "wizard");
                input.id = key;
                doc.appendChild(input);
                description = "{{.wizard.m3u.description}}";
                break;
            case "xmltv":
                var input = content.createInput("text", key, "");
                input.setAttribute("placeholder", "{{.wizard.xmltv.placeholder}}");
                input.setAttribute("class", "wizard");
                input.id = key;
                doc.appendChild(input);
                description = "{{.wizard.xmltv.description}}";
                break;
            default:
                console.log(key);
                break;
        }
        var pre = document.createElement("PRE");
        pre.innerHTML = description;
        doc.appendChild(pre);
        console.log(headline, key);
    }
}
function readyForConfiguration(wizard) {
    var server = new Server("getServerConfig");
    server.request(new Object());
    showElement("loading", false);
    configurationWizard[wizard].createWizard();
}
function saveWizard() {
    var cmd = "saveWizard";
    var div = document.getElementById("content");
    var config = div.getElementsByClassName("wizard");
    var wizard = new Object();
    for (var i = 0; i < config.length; i++) {
        var name;
        var value;
        switch (config[i].tagName) {
            case "SELECT":
                name = config[i].name;
                value = config[i].value;
                // If the value is a number, store it as a number
                if (isNaN(value)) {
                    wizard[name] = value;
                }
                else {
                    wizard[name] = parseInt(value);
                }
                break;
            case "INPUT":
                switch (config[i].type) {
                    case "text":
                        name = config[i].name;
                        value = config[i].value;
                        if (value.length == 0) {
                            var msg = name.toUpperCase() + ": " + "{{.alert.missingInput}}";
                            alert(msg);
                            return;
                        }
                        wizard[name] = value;
                        break;
                }
                break;
            default:
                // code...
                break;
        }
    }
    var data = new Object();
    data["wizard"] = wizard;
    var server = new Server(cmd);
    server.request(data);
    console.log(data);
}
// Wizard
var configurationWizard = new Array();
configurationWizard.push(new WizardItem("tuner", "{{.wizard.tuner.title}}"));
configurationWizard.push(new WizardItem("epgSource", "{{.wizard.epgSource.title}}"));
configurationWizard.push(new WizardItem("m3u", "{{.wizard.m3u.title}}"));
configurationWizard.push(new WizardItem("xmltv", "{{.wizard.xmltv.title}}"));
