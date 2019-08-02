function showConfig(obj) {
  config = obj;
  //setMenuItem();
  createMenu();
  //document.getElementById("page").className = "";
}

showMyStreams

function showMyStreams(allStreamsObj) {

  var streamTypeKeys = getObjKeys(allStreamsObj)

  for (var s = 0; s < streamTypeKeys.length; s++) {
    var streamType = streamTypeKeys[s];
    var obj = new Object();
    obj = allStreamsObj[streamType];
    switch(streamType) {
      case "activeStreams": activeStreams     = obj; break;
    }

    document.getElementById(streamType).innerHTML = "";

    var streamsObj    = new Object();
    var streamsNames  = new Array();
    
    var keys = getObjKeys(obj)
    
    // Create Object (streamsObj) for the streams and sort by name (streamsNames)
    for (var i = 0; i < keys.length; i++) {
      var name    = obj[keys[i]]["name"];
      var tmp = new Object();
      var streamKey = getObjKeys(obj[keys[i]]);

      for (var j = 0; j < streamKey.length; j++) {
        tmp[streamKey[j]] = obj[keys[i]][streamKey[j]];      
      }

      streamsObj[name] = tmp;
      streamsNames.push(name)
    }

    streamsNames.sort();
    
    // Create Table for activeStreams
    var table = document.getElementById(streamType);

    for (var i = 0; i < streamsNames.length; i++) {
      var newEntry = new Object();
      newEntry["_element"] = "TR";
      table.appendChild(createElement(newEntry));
      var line = table.lastChild;

      var tmp = streamsObj[streamsNames[i]]
      var keys = getObjKeys(tmp)

      var newKey = new Object()
      newKey["_element"]  = "TD";
      //newKey["_text"]     = streamsNames[i];
      switch(streamType) {
        case "activeStreams": newKey["_text"]     = "Channel (+):"; break;
        case "inactiveStreams": newKey["_text"]   = "Channel (-):"; break;
      }
      
      newKey["class"]     = "tdKey";
      console.log();


      var newVal = new Object()
      newVal["_element"]  = "TD";
      newVal["_text"]     = streamsNames[i];
      newVal["class"]     = "tdVal";
      //newVal["_text"]     = value;
      
      line.appendChild(createElement(newKey));
      line.appendChild(createElement(newVal));
      
    }

  }

  return
}

function showActiveStreams(obj) {
  document.getElementById("activeStreams").innerHTML = "";
  activeStreams     = obj;
  var streamsObj    = new Object();
  var streamsNames  = new Array();
  
  var keys = getObjKeys(obj)
  
  // Create Object (streamsObj) for the streams and sort by name (streamsNames)
  for (var i = 0; i < keys.length; i++) {
    var name    = obj[keys[i]]["name"];
    var tmp = new Object();
    var streamKey = getObjKeys(obj[keys[i]]);

    for (var j = 0; j < streamKey.length; j++) {
      tmp[streamKey[j]] = obj[keys[i]][streamKey[j]];      
    }

    streamsObj[name] = tmp;
    streamsNames.push(name)
  }

  streamsNames.sort();
  
  // Create Table for activeStreams
  var table = document.getElementById("activeStreams");

  for (var i = 0; i < streamsNames.length; i++) {
    var newEntry = new Object();
    newEntry["_element"] = "TR";
    table.appendChild(createElement(newEntry));
    var line = table.lastChild;

    var tmp = streamsObj[streamsNames[i]]
    var keys = getObjKeys(tmp)

    var newKey = new Object()
    newKey["_element"]  = "TD";
    //newKey["_text"]     = streamsNames[i];
    newKey["_text"]     = "Channel:";
    newKey["class"]     = "tdKey";
    console.log();


    var newVal = new Object()
    newVal["_element"]  = "TD";
    newVal["_text"]     = streamsNames[i];
    newVal["class"]     = "tdVal";
    //newVal["_text"]     = value;
    
    line.appendChild(createElement(newKey));
    line.appendChild(createElement(newVal));
    
  }

}

function parseLogs(obj) {
  log = obj
  var keys = getObjKeys(obj)

  var msgType;
  for (var i = 0; i < keys.length; i++) {
    switch(keys[i]) {
      case "warnings":  msgType = "warningMsg"; break;
      case "errors":    msgType = "errorMsg"; break;
    }

    switch(obj[keys[i]]) {
      case 0: msgType = "tdVal"; break;
      default: break;
    }

    if(document.getElementById(keys[i])){
      document.getElementById(keys[i]).className = msgType;
    }
    
    
  }
  return
}


function cancelData(element) {
  createMenu();
}

function saveData(element) {
  var data = new Object();
  var div = element.parentNode.parentNode;
  var inputs = div.getElementsByTagName("INPUT");
  
  var configKey = div.getAttribute("data-configkey");
  var menuType = div.getAttribute("data-menutype");
  var value;
  var valueArr = new Array();
  
  for (var i = 0; i < inputs.length; i++) {
    if (inputs[i].type == "text" && inputs[i].value != undefined && inputs[i].value != "" ) {
      console.log(inputs[i].value, menuType)
      switch(menuType) {
        case "inputArray": valueArr.push(inputs[i].value); break;
        case "singleInput": value = inputs[i].value; break;
      }
    }
  }

  switch(menuType) {
    case "inputArray":  data[configKey] = valueArr; break;
    case "singleInput":  
      if (isNaN(value) == false) {
        value = parseInt(value);
        data[configKey] = value;
        break;
      }

      if (value == undefined) {
        data["delete"] = configKey; 
      } else {
        data[configKey] = value
      }

      break;
  }

  data["cmd"] = "saveConfig";
  console.log(data);
  xTeVe(data)
}

function xTeVe(data) {
  if (webSockets == false) {
    alert("Your browser does not support WebSockets");
    return;
  } else {
    if (data["cmd"] != "getLog") {
      showLoadingScreen(true)
    }
  }
  delete undo["epgMapping"];
  
  var protocolWS
  switch(window.location.protocol) {
    case "http:":   protocolWS = "ws://"; break;
    case "https:":  protocolWS = "wss://"; break;
  }


  var ws = new WebSocket(protocolWS + window.location.hostname + ":" + window.location.port + "/data/" + "?Token=" + getCookie("Token"));
  ws.onopen = function() {
    console.log(data)
    ws.send(JSON.stringify(data));
  }

  ws.onmessage = function (e) {
    var response = JSON.parse(e.data);
    console.log(response);
    
    if (response.hasOwnProperty("clientInfo")) {
      createClintInfo(response["clientInfo"]);
    }

    if (response.hasOwnProperty("log")) {
      createClintInfo(response["log"]);
    }

    if (response.hasOwnProperty("status")) {
      if (response["status"] == false) {
        alert(response["err"])
        if(response.hasOwnProperty("reload")) {
          location.reload();
        }
        //checkErr(response)
        console.log(response);
        updateXteveStatus(response);
        setTimeout(function(){ showLoadingScreen(false); }, 300);
        
        return
      }

      updateXteveStatus(response)

      //console.log(data["cmd"]);
      switch(data["cmd"]) {
        case "saveUserData":    createMenu(); break;
        case "saveNewUser":     createMenu(); break;
        case "saveFilesXMLTV":  //createMenu(); break;
        case "saveFilesM3U":    //createMenu(); return; break;
        case "saveConfig":    
          data = new Object();
          data["cmd"] = "checkToken";
          xTeVe(data);
          break;

        case "emptyLog": writeLogInDiv(); break;
        case "getLog":  return; break;
      }


    }

    if (config["files"] == undefined || config["files"].length == 0) {
      createMenu();
      document.getElementById(10).click()
    }

    setTimeout(function(){ showLoadingScreen(false); }, 0);
  }
  
}

function updateXteveStatus(response) {
  var keys = getObjKeys(response);
  //console.log(keys);

  for (var i = 0; i < keys.length; i++) {
    switch(keys[i]) {
      case "alert":         alert(response[keys[i]]); break;
      case "config":        showConfig(response[keys[i]]); break;
      case "log":           parseLogs(response[keys[i]]); break;
      case "myStreams":     showMyStreams(response[keys[i]]); break;
      case "xEPG":          xEPG = response[keys[i]]; break;
      case "users":         users = response[keys[i]]; break;
      case "token":         document.cookie = "Token=" + response[keys[i]]; break;
      case "reload":        location.reload(); break;
      case "openLink":      window.location = response["openLink"]; break;
      //case "version": version = response[keys[i]]; break;
    }
  }
}

function getValueFromProviderFile(xXmltvFile, fileType, key) {

  var fileID = xXmltvFile.substring(0, xXmltvFile.lastIndexOf('.'))

  if (config["files"][fileType].hasOwnProperty(fileID) == true) {
    var data = config["files"][fileType][fileID];
    return data[key]
  }

}




