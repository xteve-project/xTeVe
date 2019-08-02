var configMenu = new Object();
var wizard = new Array("key", "tuner", "epgSource", "m3u", "complete");
var activeWizard;
var dvrIP

var configMenu_tuner = new Object();
configMenu_tuner["_element"]      = "SELECT";
configMenu_tuner["_menuType"]     = "singleInput";
configMenu_tuner["_configKey"]    = "tuner";
configMenu_tuner["_label"]        = "Available tuners";
configMenu_tuner["name"]          = "tuner";
configMenu_tuner["id"]            = "Tuner";
configMenu_tuner["placeholder"]   = "Tuner";
configMenu_tuner["_usage"]        = "This setting is only used by Plex and Emby.<br>The number of concurrent streams allowed by the IPTV provider."


var optionValues = new Array();
for (var i = 1; i <= 100; i++) {
  optionValues.push(i)
}
configMenu_tuner["_optionValues"] = optionValues;

var configMenu_epg = new Object();
configMenu_epg["_element"]      = "SELECT";
configMenu_epg["_menuType"]     = "singleInput";
configMenu_epg["_configKey"]    = "epgSource";
configMenu_epg["_label"]        = "Selection of the EPG source";
configMenu_epg["name"]          = "epgSource";
configMenu_epg["id"]            = "EPG source";
configMenu_epg["placeholder"]   = "EPG source";
configMenu_epg["_optionValues"] = new Array("PMS", "XEPG");
configMenu_epg["_usage"]        = "PMS:   Use EPG data from Plex or Emby<br>XEPG:  Use of external EPG data (XMLTV)<br>       Several XMLTV sources possible<br>       Allows editing and order channels<br>       M3U / XMLTV export (HTTP link for IPTV apps)"

var configMenu_m3u = new Object();
configMenu_m3u["_element"]        = "INPUT";
configMenu_m3u["_menuType"]       = "inputArray";
configMenu_m3u["_configKey"]      = "file";
configMenu_m3u["_label"]          = "M3U File: local or remote";
configMenu_m3u["name"]            = "file";
configMenu_m3u["id"]              = "m3u";
configMenu_m3u["type"]            = "text";
configMenu_m3u["placeholder"]     = "M3U File";
configMenu_m3u["_usage"]          = "Remote playlist: http://your.provider.com/file.m3u<br>Local  playlist: /path/to/file.m3u"


configMenu_m3u["value"]           = "http://websrv.local:8080/kabel.m3u";

var configMenu_complete = new Object();
configMenu_complete["_element"]        = "H2";
configMenu_complete["_menuType"]       = "inputArray";
configMenu_complete["_configKey"]      = "file";
configMenu_complete["_text"]           = "xTeVe was successfully set up";
configMenu_complete["name"]            = "complete";
configMenu_complete["id"]              = "complete";
configMenu_complete["type"]            = "text";
configMenu_complete["class"]           = "center";

configMenu["tuner"]     = configMenu_tuner;
configMenu["epgSource"] = configMenu_epg;
configMenu["m3u"]       = configMenu_m3u;
configMenu["complete"]  = configMenu_complete;

function readyForConfiguration() {
  var data = new Object();
  data["cmd"] = "getServerConfig";
  xTeVe(data);
  showLoadingScreen(false);
}

function createConfiguration(elm) {

  activeWizard = elm;
  var item  = configMenu[elm];

  var div   = document.getElementById("content");
  div.innerHTML = "";
  div.setAttribute("data-configKey", item["_configKey"]);
  div.setAttribute("data-menuType", item["_menuType"]);
  
  switch(item.hasOwnProperty("_label")) {
    case true:
      var newItem = new Object();
      newItem["_element"] = "LABEL";
      newItem["_text"]    = item["_label"]; 
      newItem["for"]      = item["id"];
      div.appendChild(createElement(newItem));
      break
  }

  switch(item["_element"]) {
    case "SELECT":
      div.appendChild(createElement(item));
      var selectElement = div.getElementsByTagName("SELECT")[0];
      var values = item["_optionValues"];
      for (var i = 0; i < values.length; i++) {
        var newEntry = new Object;
        newEntry["_element"]  = "OPTION";
        newEntry["_text"]     = item["id"] + ": " + values[i];
        newEntry["value"]     = values[i];
        selectElement.appendChild(createElement(newEntry));
      }
      //return
      break;

    default: 
      div.appendChild(createElement(item));
      break;


  }
  //alert()

  switch(item.hasOwnProperty("_usage")) {
    case true: 
      var usageItem = new Object();
      usageItem["_element"] = "PRE"
      usageItem["_text"]    = item["_usage"];
      div.appendChild(createElement(usageItem));
  }

  if (activeWizard == "complete") {
    document.getElementById("next").value = "Finished"
    //document.getElementById("next").setAttribute("onclick", "javascript: location.reload();")
  }

  //div.appendChild(createElement(item));
}

function saveData() {

  var div   = document.getElementById("content");
  var inputs = div.getElementsByTagName("INPUT");
  var selects = div.getElementsByTagName("SELECT");
  var value;
  var data = new Object();
  var valueArr = new Array();
  var newData = false;

  if (activeWizard == "complete") {
    data["cmd"] = "wizardCompleted";
    showLoadingScreen(true)
    xTeVe(data);
    return
  }
  
  for (var i = 0; i < inputs.length; i++) {
    var menuType = inputs[i].parentElement.getAttribute("data-menutype");
    if (inputs[i].value != undefined && inputs[i].value != "" ) {
      newData = true;

      console.log(inputs[i].id)
      switch(inputs[i].id) {
        case "m3u": 
          var newPlaylist = new Object();
          newPlaylist["file.source"] = inputs[i].value;
          //newPlaylist["name"] = inputs[i].value;
          newPlaylist["type"] = "m3u";
          newPlaylist["new"] = true;

          data["files"] = new Object();
          data["files"]["m3u"] = new Object();
          data["files"]["m3u"]["-"] = newPlaylist;
          
          data["cmd"] = "saveFilesM3U";
          xTeVe(data)
          return
      }
      /*
      switch(menuType) {
        case "singleInput":
          data[inputs[i].name] = inputs[i].value; break;
        case "inputArray": 
          valueArr.push(inputs[i].value);
          data[inputs[i].name] = valueArr; break

      }
      */
    } else {
      inputs[i].style.borderBottomColor = "red";
      return;
    }
  }


  for (var i = 0; i < selects.length; i++) {
    var value = selects[i].options[selects[i].selectedIndex].value;
    if (isNaN(value) == false) {
      value = parseInt(value);
      data[selects[i].name] = value;
      newData = true;
      break;
    }
    data[selects[i].name] = value;
    newData = true;
  }


  //console.log(data, newData);
  if (newData == true) {
    config = data
    data["cmd"] = "saveConfig";
    xTeVe(data);
  }
}

function xTeVe(data) {

  if (webSockets == false) {
    alert("Your browser does not support WebSockets");
    return;
  }

  if (activeWizard == "m3u" || activeWizard == "epgSource") {
    showLoadingScreen(true);
  }

  var protocolWS
  switch(window.location.protocol) {
    case "http:":   protocolWS = "ws://"; break;
    case "https:":  protocolWS = "wss://"; break;
  }

  var ws = new WebSocket(protocolWS + window.location.hostname + ":" + window.location.port + "/data/" + "?Token=" + getCookie("Token"));
  
  ws.onopen = function() {
    ws.send(JSON.stringify(data));
  }

  ws.onmessage = function (e) {
    
    var response = JSON.parse(e.data);
    
    if (response.hasOwnProperty("clientInfo")) {
      createClintInfo(response["clientInfo"]);
    }

    if (response.hasOwnProperty("status")) {
      if (response["status"] == false) {
        document.getElementById("headline").style.borderColor = "red";
        showErr(response["err"]);
        showLoadingScreen(false)
        return
      } else {
        document.getElementById("err").innerHTML = "";
        document.getElementById("headline").style.borderColor = "lawngreen";
      }

      dvrIP = response["DVR"]
      switch(response["configurationWizard"]) {
        case true: 
          if (activeWizard == undefined) {
            activeWizard = wizard[0]
          }
          var n = wizard.indexOf(activeWizard);
          n++;
          activeWizard = wizard[n]

          if (activeWizard == undefined) {
            data["cmd"] = "wizardCompleted";
            xTeVe(data)
          } else {
            //console.log(activeWizard);
            createConfiguration(activeWizard); 
          }
          
        break;
      }

      switch(response["reload"]) {
        
        
        case true: 
    
          setTimeout(function(){ 
            location.reload();
          }, 100);
          
          //location.reload();
          
          break;
        
      }

      
    }

    setTimeout(function(){ showLoadingScreen(false); }, 300);
  }
  
}

function showErr(elm) {
  document.getElementById("err").innerHTML = elm;
}