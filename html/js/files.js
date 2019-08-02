function openFiles(elm, fileType) {
  //document.getElementById("settings").innerHTML = "Test";
  
  columnToSort = 0; 
  var newDiv = document.getElementById("settings");
  
  var newEntry = new Object();
  newEntry["_element"]  = "HR";
  newDiv.appendChild(createElement(newEntry));

  var newEntry = new Object();
  newEntry["_element"]  = "INPUT";
  newEntry["type"] = "button";
  newEntry["class"] = "button";
  newEntry["value"] = "New";
  newEntry["onclick"] = 'fileDetail("-", "' + fileType + '")';
  newDiv.appendChild(createElement(newEntry));

  var newEntry = new Object();
  newEntry["_element"]  = "INPUT";
  newEntry["type"] = "button";
  newEntry["class"] = "button";
  newEntry["value"] = "Update";
  newEntry["onclick"] = "fileDetail(0)";
  //newDiv.appendChild(createElement(newEntry));

  var div = document.getElementById("settings");

   // Build table
  var newTable = new Object();
  newTable["_element"]  = "TABLE";
  newTable["id"]        = "id_mapping";
  newTable["class"]     = "table-mapping";
  div.appendChild(createElement(newTable));

  setTimeout(function(){ 
    createFilesTable(fileType); 
  }, 10);

}

function createFilesTable(fileType) {
  var table = document.getElementById("id_mapping");
  var availableFileTypes = new Array();
  
  table.innerHTML = "";
  var newTR = new Object();
  newTR["_element"] = "TR";
  newTR["class"]    = "table-mapping-header";
  table.appendChild(createElement(newTR));

  var tr = table.lastChild;

  switch(fileType) {
    case "xmltv": 
      availableFileTypes = new Array("xmltv"); 
      var trHeadlines = new Array("Guide", "Last Update", "Availability %", "Channels", "Programs")
      var compatibilityKeys = new Array("xmltv.channels", "xmltv.programs")
      break;

    case "m3u":
      availableFileTypes = new Array("m3u", "hdhr"); 
      var trHeadlines = new Array("Playlist", "Last Update", "Availability %", "Type", "Streams", "group-title %", "tvg-id %", "Unique ID %");
      var compatibilityKeys = new Array("streams", "group.title", "tvg.id", "stream.id");
      break;
  }

  for (var i = 0; i < trHeadlines.length; i++) {
    var newTD = new Object();
    newTD["_element"] = "TD";
    newTD["_text"]    = trHeadlines[i];
    tr.appendChild(createElement(newTD));
  }
  
  for (var i = 0; i < availableFileTypes.length; i++) {
    
    var fileType = availableFileTypes[i]

    var data = config["files"][fileType];
    
    var allFiles = getObjKeys(data)
  
    for (var f = 0; f < allFiles.length; f++) {
      var elm           = data[allFiles[f]];
      var table         = document.getElementById("id_mapping");
      var fileID        = elm["id.provider"];
      var name          = elm["name"];
      var lastUpdate    = elm["last.update"];
      var availability  = elm["provider.availability"];
      var type          = elm["type"].toUpperCase();
      var compatibility = elm["compatibility"];

      // Create TR
      var newTR = new Object();
      newTR["_element"]       = "TR";
      newTR["class"]          = "";
      newTR["id"]             = fileID;
      newTR["onclick"]        = 'javascript: fileDetail("' + fileID + '","' + fileType + '");';
      table.appendChild(createElement(newTR));

      var tr = table.lastChild;

      // Create file name TD
      var newTD = new Object();
      newTD["_element"] = "P";
      newTD["_text"]    = name;
      createNewTD(newTD, tr);

      // Create last update TD
      var newTD = new Object();
      newTD["_element"] = "P";
      newTD["_text"]    = lastUpdate;
      createNewTD(newTD, tr);

      // Create availability TD
      var newTD = new Object();
      newTD["_element"] = "P";
      newTD["_text"]    = availability;
      createNewTD(newTD, tr);

      if (fileType == "m3u" || fileType == "hdhr") {

        // Create Type TD
        var newTD = new Object();
        newTD["_element"] = "P";
        newTD["_text"]    = type;
        createNewTD(newTD, tr);
  
      }
      
      // Create all compatibility TDs

      for (var j = 0; j < compatibilityKeys.length; j++) {
        var newTD = new Object();
        newTD["_element"] = "P";
        newTD["_text"]    = compatibility[compatibilityKeys[j]];
        createNewTD(newTD, tr);
      }

    }

  }
  
  
  sortTable(0)

  // usage Info  
  var div = document.getElementById("settings");
  switch(menu[activeMenu.id].hasOwnProperty("_usage")) {
    case true: 
      var usageItem = new Object();
      usageItem["_element"] = "PRE"
      usageItem["_text"]    = menu[activeMenu.id]["_usage"];

      var newHR = new Object();
      newHR["_element"] = "HR"
      div.appendChild(createElement(newHR));
      div.appendChild(createElement(usageItem));
      break;
  }

  calculateWrapperHeight();
  return;
}


function fileDetail(fileID, fileType) {

  optionsText  = new Array("M3U", "HDHomeRun - [Experimental]")
  optionsValue = new Array("m3u", "hdhr")

  switch (fileType) {
    
    case "m3u": 
      document.getElementById("name").setAttribute("placeholder", "Playlist name"); 
      document.getElementById("description").setAttribute("placeholder", "Description of this playlist"); 
      document.getElementById("file-detail-headline").innerHTML = "M3U Playlist"; 
      document.getElementById("file-path").innerHTML = "M3U File:"; 
      document.getElementById("file.source").setAttribute("placeholder", "Local or remote");
      break;

    case "hdhr": 
      document.getElementById("name").setAttribute("placeholder", "HDHomeRun name"); 
      document.getElementById("description").setAttribute("placeholder", "Description of this HDHomeRun tuner"); 
      document.getElementById("file-detail-headline").innerHTML = "HDHomeRun"; 
      document.getElementById("file-path").innerHTML = "HDHomeRun IP:"; 
      document.getElementById("file.source").setAttribute("placeholder", "IP address and port of the tuner (192.168.1.10:5004)");
      break;
    
    case "xmltv": 
      document.getElementById("name").setAttribute("placeholder", "XMLTV name"); 
      document.getElementById("description").setAttribute("placeholder", "Description of this XMLTV file"); 
      document.getElementById("file-detail-headline").innerHTML = "XMLTV File"; 
      document.getElementById("file-path").innerHTML = "XMLTV File:";
      document.getElementById("file.source").setAttribute("placeholder", "Local or remote");

      optionsText  = new Array("XMLTV")
      optionsValue = new Array("xmltv")
      break;
  }

  modifyOption("type", optionsText, optionsValue)
  
  showPopUpElement('file-detail');

  document.getElementById("saveFileDetail").setAttribute("onclick", 'javascript: saveFileDetail("' + fileID + '","' + fileType + '", false)');
  document.getElementById("updateFileDetail").setAttribute("onclick", 'javascript: updateFile("' + fileID + '","' + fileType + '", false)');
  document.getElementById("deleteFileDetail").setAttribute("onclick", 'javascript: saveFileDetail("' + fileID + '","' + fileType + '", true)');

  var data = new Object();

  switch(fileID) {

    case "-": // New file
      data["name"]        = "";
      data["description"] = "";
      data["file.source"] = "";
      data["type"] = fileType;
      
      document.getElementById("deleteFileDetail").className = "delete";
      document.getElementById("type").setAttribute("onchange", "changeFileType(this);")
      document.getElementById("type").setAttribute("data-id", fileID)
      
      showElement("deleteFileDetail", false);
      showElement("updateFileDetail", false);
      
      if (fileType == "xmltv") {
        showElement("type", false);
        showElement("file-type", false);
      } else {
        showElement("type", true);
        showElement("file-type", true);
      }
      
      break;

    default: 
      data = config["files"][fileType][fileID];
      document.getElementById("deleteFileDetail").className = "delete";
      
      showElement("updateFileDetail", true);
      showElement("type", false);
      showElement("file-type", false);
      
      break;

  }

  var keys = getObjKeys(data);
  
  for (var i = 0; i < keys.length; i++) {

    if(document.getElementById(keys[i])){
      document.getElementById(keys[i]).value = data[keys[i]];
    } 


  }

}

function changeFileType(elm) {

  var fileID = elm.getAttribute("data-id");
  var fileType = elm.options[elm.selectedIndex].value;
  
  fileDetail(fileID, fileType)

}


function saveFileDetail(fileID, fileType, deleteFile) {

  if (fileID == undefined) {
    alert("ID is missing!!!");
    return 
  }

  var inputs      = document.getElementById("file-detail").getElementsByTagName("INPUT");
  var selects     = document.getElementById("file-detail").getElementsByTagName("SELECT");
  var newFileData = new Object();
  var data        = new Object();

  for (var i = 0; i < inputs.length; i++) {
    switch(inputs[i].type) {
      case "text": newFileData[inputs[i].name] = inputs[i].value; break;
    }
  }

  for (var i = 0; i < selects.length; i++) {
    newFileData[selects[i].id] = selects[i].options[selects[i].selectedIndex].value;
  }

  if (deleteFile == true) {
    switch(fileType) {
      case "m3u":   var alertText = "Delete this playlist?"; break;
      case "hdhr": var alertText = "Delete this HDHomeRun tuner?"; break;
      case "xmltv": var alertText = "Delete this XMLTV file?"; break;
    }

    if (confirm(alertText)) {
      newFileData["delete"] = true
      data = buildFilesObj(fileType, fileID, newFileData);
      console.log(data);
      
    } else {
      showElement("popup", false);
      return
    
    }

  } else {
  
    switch(config["files"][fileType].hasOwnProperty(fileID)) {

      case true: 
        data = config["files"][fileType][fileID]; 
        if (data["file.source"] != newFileData["file.source"]) {
          data["update"] = true
        } else {
          data["updatePlaylistName"] = true;
        }
        break;
      
      case false: 
        newFileData["new"] = true;
        data = buildFilesObj(fileType, fileID, newFileData);
        break

    }
  
  }  
  
  switch(fileType) {

    case "m3u":   data["cmd"] = "saveFilesM3U"; break;
    case "hdhr":  data["cmd"] = "saveFilesHDHR"; break;
    case "xmltv": data["cmd"] = "saveFilesXMLTV"; break;

  }
  //console.log(data);
  xTeVe(data);
  return
}

function updateFile(fileID, fileType, allFiles) {
  
  switch(config["files"][fileType].hasOwnProperty(fileID)) {

    case true: 
    
      var data = new Object();
      var data = buildFilesObj(fileType, fileID, config["files"][fileType][fileID])
      data["new"] = true

      switch(fileType) {

        case "m3u":   data["cmd"] = "updateFileM3U"; break;
        case "hdhr":  data["cmd"] = "updateFileHDHR"; break;
        case "xmltv": data["cmd"] = "updateFileXMLTV"; break;

      }
      
      xTeVe(data);
      
      break;
  }

}

function buildFilesObj(fileType, fileID, obj) {

  var data = new Object();
  data["files"] = new Object();
  data["files"][fileType] = new Object();
  data["files"][fileType][fileID] = obj
  return data

}