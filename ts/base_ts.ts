var SERVER = new Object()
var BULK_EDIT:Boolean = false
var COLUMN_TO_SORT:number
var SEARCH_MAPPING = new Object()
var UNDO = new Object()
var SERVER_CONNECTION = false
var WS_AVAILABLE = false


// Menu
var menuItems = new Array()
menuItems.push(new MainMenuItem("playlist", "{{.mainMenu.item.playlist}}", "m3u.png", "{{.mainMenu.headline.playlist}}"))
//menuItems.push(new MainMenuItem("pmsID", "{{.mainMenu.item.pmsID}}", "number.png", "{{.mainMenu.headline.pmsID}}"))
menuItems.push(new MainMenuItem("filter", "{{.mainMenu.item.filter}}", "filter.png", "{{.mainMenu.headline.filter}}"))
menuItems.push(new MainMenuItem("xmltv", "{{.mainMenu.item.xmltv}}", "xmltv.png", "{{.mainMenu.headline.xmltv}}"))
menuItems.push(new MainMenuItem("mapping", "{{.mainMenu.item.mapping}}", "mapping.png", "{{.mainMenu.headline.mapping}}"))
menuItems.push(new MainMenuItem("users", "{{.mainMenu.item.users}}", "users.png", "{{.mainMenu.headline.users}}"))
menuItems.push(new MainMenuItem("settings", "{{.mainMenu.item.settings}}", "settings.png", "{{.mainMenu.headline.settings}}"))
menuItems.push(new MainMenuItem("log", "{{.mainMenu.item.log}}", "log.png", "{{.mainMenu.headline.log}}"))
menuItems.push(new MainMenuItem("logout", "{{.mainMenu.item.logout}}", "logout.png", "{{.mainMenu.headline.logout}}"))

// Settings categories
var settingsCategory = new Array()
settingsCategory.push(new SettingsCategoryItem("{{.settings.category.general}}", "tlsMode,xteveAutoUpdate,hostIP,tuner,epgSource,disallowURLDuplicates,api"))
settingsCategory.push(new SettingsCategoryItem("{{.settings.category.files}}", "update,files.update,temp.path,cache.images,xepg.replace.missing.images"))
settingsCategory.push(new SettingsCategoryItem("{{.settings.category.streaming}}", "buffer,udpxy,buffer.size.kb,storeBufferInRAM,buffer.timeout,user.agent,ffmpeg.path,ffmpeg.options,vlc.path,vlc.options"))
settingsCategory.push(new SettingsCategoryItem("{{.settings.category.backup}}", "backup.path,backup.keep"))
settingsCategory.push(new SettingsCategoryItem("{{.settings.category.authentication}}", "authentication.web,authentication.pms,authentication.m3u,authentication.xml,authentication.api"))

function showPopUpElement(elm) {

  var allElements = new Array("popup-custom");

  for (var i = 0; i < allElements.length; i++) {
    showElement(allElements[i], false)
  }

  showElement(elm, true)

  setTimeout(function(){
    showElement("popup", true);
  }, 10);

  return
}

function showElement(elmID, type) {

  var cssClass:string
  switch(type) {
    case true:  cssClass = "block"; break;
    case false: cssClass = "none"; break;
  }

  document.getElementById(elmID).className = cssClass;
}

function changeButtonAction(element, buttonID, attribute) {
  var value = element.options[element.selectedIndex].value;
  document.getElementById(buttonID).setAttribute(attribute, value)
}

function getLocalData(dataType, id):object {
  let data = {}
  switch(dataType) {
    case "m3u":
      data = SERVER["settings"]["files"][dataType][id]
      break

    case "hdhr":
      data = SERVER["settings"]["files"][dataType][id]
      break

    case "filter":
    case "custom-filter":
    case "group-title":
      if (id == -1) {
        data["active"] = true
        data["caseSensitive"] = false
        data["description"] = ""
        data["exclude"] = ""
        data["filter"] = ""
        data["include"] = ""
        data["name"] = ""
        data["type"] = "group-title"
        data["preserveMapping"] = true
        data["startingChannel"] = SERVER["settings"]["mapping.first.channel"]
        data["defaultMissingEPG"] = "-"
        SERVER["settings"]["filter"][id] = data
      }
      data = SERVER["settings"]["filter"][id]
      break

    case "xmltv":
      data = SERVER["settings"]["files"][dataType][id]
      break

    case "users":
      data = SERVER["users"][id]["data"]
      break

    case "mapping":
      data = SERVER["xepg"]["epgMapping"][id]
      break

    case "m3uGroups":
      data = SERVER["data"]["playlist"]["m3u"]["groups"]
      break
  }

  return data
}

function getOwnObjProps(object: Object): string[] {
  return object ? Object.getOwnPropertyNames(object) : [];
}

function getAllSelectedChannels():string[] {

  var channels:string[] = new Array()

  if (BULK_EDIT == false) {
    return channels
  }

  var trs = document.getElementById("content_table").getElementsByTagName("TR")

  for (var i = 1; i < trs.length; i++) {

    if ((trs[i] as HTMLElement).style.display != "none") {

      if ((trs[i].firstChild.firstChild as HTMLInputElement).checked == true) {
        channels.push(trs[i].id)
      }

    }

  }

  return channels
}

function selectAllChannels() {

  var bulk:Boolean = false
  var trs = document.getElementById("content_table").getElementsByTagName("TR")

  if ((trs[0].firstChild.firstChild as HTMLInputElement).checked == true) {
    bulk = true
  }

  for (var i = 1; i < trs.length; i++) {

    if ((trs[i] as HTMLElement).style.display != "none") {

      switch (bulk) {

        case true:
          (trs[i].firstChild.firstChild as HTMLInputElement).checked = true
          break

        case false:
          (trs[i].firstChild.firstChild as HTMLInputElement).checked = false
          break

      }

    }

  }

  return
}

function bulkEdit() {

  BULK_EDIT = !BULK_EDIT
  var className:string
  var rows = document.getElementsByClassName("bulk");

  switch (BULK_EDIT) {
    case true:
      className = "bulk showBulk"
      break;

    case false:
      className = "bulk hideBulk"
      break;
  }

  for (var i = 0; i < rows.length; i++) {
    rows[i].className = className;
    (rows[i] as HTMLInputElement).checked = false
  }

  return
}

function sortTable(column) {

  if (column == COLUMN_TO_SORT) {
    return;
  }

  const table       = document.getElementById("content_table");
  const tableHead   = table.getElementsByTagName("TR")[0];
  const tableItems  = tableHead.getElementsByTagName("TD");

  type SortEntry = {
    key: string | number;
    row: HTMLTableRowElement;
  }

  const sortArr: SortEntry[] = [];
  let xValue: string | number;

  if (column >= 0 && COLUMN_TO_SORT >= 0)  {
    tableItems[COLUMN_TO_SORT].className = "pointer";
    tableItems[column].className = "sortThis";
  }

  COLUMN_TO_SORT = column;

  const rows = (table as HTMLTableElement).rows;

  if (rows[1] != undefined) {
    const tableHeader = rows[0];

    let x: any = rows[1].getElementsByTagName("TD")[column];

    for (let i = 1; i < rows.length; i++) {

      x = rows[i].getElementsByTagName("TD")[column];

      switch(x.childNodes[0].tagName.toLowerCase()) {
        case "input":
          xValue = x.getElementsByTagName("INPUT")[0].value.toLowerCase();
          break;

        case "p":
          xValue = x.getElementsByTagName("P")[0].innerText.toLowerCase();
          break;

        default:
          break;
      }

      sortArr.push({key: xValue ? xValue : i, row: rows[i]});

    }

    while (table.firstChild) {
      table.removeChild(table.firstChild);
    }

    sortArr.sort((se1: SortEntry, se2: SortEntry): number => {
      const se1KeyNum = parseFloat(String(se1.key));
      const se2KeyNum = parseFloat(String(se2.key));

      if (!isNaN(se1KeyNum) && !isNaN(se2KeyNum)) {
        return se1KeyNum - se2KeyNum;
      }

      if (se1.key < se2.key) {
        return -1;
      }

      if (se1.key > se2.key) {
        return 1;
      }

      return 0;
    });

    table.appendChild(tableHeader);

    sortArr.forEach((se: SortEntry) => {
      table.appendChild(se.row);
    });

  }

  return
}

function createSearchObj() {

  SEARCH_MAPPING = new Object()
  var data = SERVER["xepg"]["epgMapping"]
  var channels = getOwnObjProps(data)

  var channelKeys:string[] = ["x-active", "x-channelID", "x-name", "updateChannelNameRegex", "_file.m3u.name", "x-group-title", "x-xmltv-file"]

  channels.forEach(id => {

    channelKeys.forEach(key => {

      if (key == "x-active") {

        switch (data[id][key]) {
          case true:
            SEARCH_MAPPING[id] = "online "
            break;

          case false:
            SEARCH_MAPPING[id] = "offline "
            break;

        }

      } else {

        if (key == "x-xmltv-file") {
          var xmltvFile = getValueFromProviderFile(data[id][key], "xmltv", "name")

          if (xmltvFile != undefined) {
            SEARCH_MAPPING[id] = SEARCH_MAPPING[id] + xmltvFile + " "
          }

        } else {
          SEARCH_MAPPING[id] = SEARCH_MAPPING[id] + data[id][key] + " "
        }


      }

    })

  })

  return
}

function searchInMapping() {

  var searchValue = (document.getElementById("searchMapping") as HTMLInputElement).value;
  var trs = document.getElementById("content_table").getElementsByTagName("TR")

  for (var i = 1; i < trs.length; ++i) {

    var id = trs[i].getAttribute("id")
    var element = SEARCH_MAPPING[id]

    switch (element.toLowerCase().includes(searchValue.toLowerCase())) {
      case true:
        document.getElementById(id).style.display = ""
        break;

      case false:
        document.getElementById(id).style.display = "none"
        break;
    }


  }

  return
}

function calculateWrapperHeight() {

  if (document.getElementById("box-wrapper")){

    var elm = document.getElementById("box-wrapper");

    var divs = new Array("myStreamsBox", "clientInfo", "content");
    var elementsHeight = 0 - elm.offsetHeight;
    for (var i = 0; i < divs.length; i++) {
      elementsHeight = elementsHeight + document.getElementById(divs[i]).offsetHeight;
    }

    elm.style.height = window.innerHeight - elementsHeight + "px";

  }

  return
}

function changeChannelNumber(element) {

  var dbID = element.parentNode.parentNode.id

  var newNumber:number = parseFloat(element.value)
  var channelNumbers:number[] = []
  var data = SERVER["xepg"]["epgMapping"]
  var channels = getOwnObjProps(data)

  if (isNaN(newNumber)) {
    alert("{{.alert.invalidChannelNumber}}")
    return
  }

  channels.forEach(id => {

    var channelNumber = parseFloat(data[id]["x-channelID"])
    channelNumbers.push(channelNumber)

  })

  for (var i = 0; i < channelNumbers.length; i++) {

    if (channelNumbers.indexOf(newNumber) == -1) {
      break
    }

    if (Math.floor(newNumber) == newNumber) {
      newNumber = newNumber + 1
    } else {
      newNumber = newNumber + 0.1;
      newNumber.toFixed(1)
      newNumber = Math.round(newNumber * 10) / 10
    }

  }

  data[dbID]["x-channelID"] = newNumber.toString()
  element.value = newNumber

  if (COLUMN_TO_SORT == 1) {
    COLUMN_TO_SORT = -1
    sortTable(1)
  }

  return
}

function backup() {

  var data = new Object()
  var cmd = "xteveBackup"
  var server:Server = new Server(cmd)
  server.request(data)

  return

}

function toggleChannelStatus(id:string) {

  var element:any
  var status:boolean

  if(document.getElementById("active")) {
    var checkbox = (document.getElementById("active") as HTMLInputElement)
    status = (checkbox).checked
  }


  var ids:string[] = getAllSelectedChannels()
  if (ids.length == 0) {
    ids.push(id)
  }

  ids.forEach(id => {

    var channel = SERVER["xepg"]["epgMapping"][id]

    channel["x-active"] = status

    switch (channel["x-active"]) {
      case true:
        if (channel["x-xmltv-file"] == "-" || channel["x-mapping"] == "-") {

          if (BULK_EDIT == false) {
            alert(channel["x-name"] + ": Missing XMLTV file / channel")
            checkbox.checked = false
          }

          channel["x-active"] = false

        }

        break

      case false:
        // code...
        break;
    }

    if (channel["x-active"] == false) {
      document.getElementById(id).className = "notActiveEPG"
    } else {
      document.getElementById(id).className = "activeEPG"
    }

  });

}

function toggleGroupUpdateCb(xepgId: string, target: HTMLInputElement) {
  target.className = 'changed';

  const groupInput: HTMLInputElement = document.querySelector('input[name="x-group-title"]');
  const mapping = getLocalData('mapping', xepgId);

  if (target.checked) {
    groupInput.dataset.oldValue = groupInput.value;
    groupInput.value = mapping['group-title'];
    groupInput.disabled = true;
  } else {
    groupInput.value = groupInput.dataset.oldValue;
    groupInput.disabled = false;
  }

  groupInput.className = 'changed';
}

function restore() {

  if (document.getElementById('upload')) {
    document.getElementById('upload').remove()
  }

  var restore = document.createElement("INPUT");
  restore.setAttribute("type", "file");
  restore.setAttribute("class", "notVisible");
  restore.setAttribute("name", "");
  restore.id = "upload";

  document.body.appendChild(restore);
  restore.click();

  restore.onchange = function() {

    var filename = (restore as HTMLInputElement).files[0].name
    var check = confirm("File: " + filename + "\n{{.confirm.restore}}");

    if (check == true) {

      var reader  = new FileReader();
      var file = (document.querySelector('input[type=file]') as HTMLInputElement).files[0];

      if (file) {

        reader.readAsDataURL(file);
        reader.onload = function() {
          var data = new Object();
          var cmd = "xteveRestore"
          data["base64"]  = reader.result

          var server:Server = new Server(cmd)
          server.request(data)

        };

      } else {
        alert("File could not be loaded")
      }

      restore.remove()
      return
    }

  }

  return
}

function uploadLogo() {

  if (document.getElementById('upload')) {
    document.getElementById('upload').remove()
  }

  var upload = document.createElement("INPUT");
  upload.setAttribute("type", "file");
  upload.setAttribute("class", "notVisible");
  upload.setAttribute("name", "");
  upload.id = "upload";

  document.body.appendChild(upload);
  upload.click();

  upload.onblur = function() {
    alert()
  }

  upload.onchange = function() {

    var filename = (upload as HTMLInputElement).files[0].name

    var reader  = new FileReader();
    var file = (document.querySelector('input[type=file]') as HTMLInputElement).files[0];

    if (file) {

      reader.readAsDataURL(file);
      reader.onload = function() {
        var data = new Object();
        var cmd = "uploadLogo"
        data["base64"]  = reader.result
        data["filename"]  = file.name

        var server:Server = new Server(cmd)
        server.request(data)

        var updateLogo = (document.getElementById('update-icon') as HTMLInputElement)
        updateLogo.checked = false
        updateLogo.className = "changed"

      };

    } else {
      alert("File could not be loaded")
    }

    upload.remove()
    return
  }

}

function checkUndo(key:string) {

  switch (key) {
    case "epgMapping":
      if (UNDO.hasOwnProperty(key)) {
        SERVER["xepg"][key] = JSON.parse(JSON.stringify(UNDO[key]))
      } else {
        UNDO[key] = JSON.parse(JSON.stringify(SERVER["xepg"][key]));
      }
      break;

    default:

      break;
  }

  return
}

function updateLog() {

  var server:Server = new Server("updateLog")
  server.request(new Object())

}
