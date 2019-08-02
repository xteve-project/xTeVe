var config        = new Object();
var menu          = new Object();
var subMenu       = new Object();
var activeStreams = new Object();
var xEPG          = new Object();
var users         = new Object();
var log           = new Object();
var undo          = new Object();
var webSockets    = true;
var closeLog, version, activeMenu;
var columnToSort  = 0


if (window.WebSocket === undefined) {
  alert("Your browser does not support WebSockets");
  webSockets = false;
} 

function pageReady() {
  var data = new Object();
  data["cmd"] = "getServerConfig";
  xTeVe(data);
  //showLoadingScreen(false);

  var resizeHandle = document.getElementById("openStreams");
  var box = document.getElementById("myStreamsBox");
  resizeHandle.addEventListener("mousedown", initialiseResize, false);

  function initialiseResize(e) {
    window.addEventListener("mousemove", startResizing, false);
    window.addEventListener("mouseup", stopResizing, false);
  }

  function startResizing(e) {
    box.style.height = (e.clientY - box.offsetTop) + "px";
      
    var elm = document.getElementById("allStreams");
    if (e.clientY > 120) {
      elm.className = "visible"; 
    } else {
      elm.className = "notVisible"; 
    }

    calculateWrapperHeight();

  }
  function stopResizing(e) {
      window.removeEventListener('mousemove', startResizing, false);
      window.removeEventListener('mouseup', stopResizing, false);
      calculateWrapperHeight();
  }

  window.addEventListener("resize", function(){
    calculateWrapperHeight();
  }, true);
}


function getObjKeys(obj) {
  var keys = new Array();

  for (var i in obj) {
    if (obj.hasOwnProperty(i)) {
      keys.push(i);
    }
  }

  return keys;
}


function createElement(item) {
  //console.log(item);
  var element = document.createElement(item["_element"]);
  if (item.hasOwnProperty("_text")) {
    //element.innerHTML = "<p>" + item["_text"] + "</p>";
    element.innerHTML = item["_text"];
  }

  var keys = getObjKeys(item);
  for (var i = 0; i < keys.length; i++) {
    if (keys[i].charAt(0) != "_") {
      //console.log(keys[i], item[keys[i]]);
      element.setAttribute(keys[i], item[keys[i]]);
    }
  }

  //console.log(element);
  return element;
}

function modifyOption(id, options, values) {
  var select = document.getElementById(id);
  select.innerHTML = "";

  for (var i = 0; i < options.length; i++) {
  
    var element = document.createElement("OPTION")
  
    element.value = values[i];
    element.innerHTML = options[i];

    document.getElementById(id).appendChild(element);
  
  }

}


function startWebSocket() {
  if (webSockets == false) {
    return;
  }

  //ws.send('{"cmd": "getServerConfig1"}');

}

function checkErr(obj) {
  //alert(obj["err"])
  //screenLog(obj["err"], "error")
  console.log(obj);
  var newObj = new Object();
  var newErr = new Object();
  newErr["key"]   = "Error";
  newErr["value"] = obj["err"];
  newErr["type"]  = "error";

  newObj[0] = newErr
  showLog(newObj);
  return
}

function screenLog(msg, msgType, show) {
  return
  clearTimeout(closeLog)
  var div = document.getElementById("screenLog");
  var newMsg = new Object();
  
  newMsg["_element"] = "P";
  
  switch(msgType) {
    case "error":   newMsg["class"] = "errorMsg"; break;
    case "warning": newMsg["class"] = "warningMsg"; break;
    //default:      newMsg["class"] = "infoMsg"
  }

  newMsg["_text"] = msg;

  div.appendChild(createElement(newMsg));

  div.scrollTop = div.scrollHeight;

  if (show == false) {
    return;
  }

  div.className = ""
  closeLog = setTimeout(closeScreenLog, 10000);
}


function closeScreenLog() {
  var div = document.getElementById("screenLog");
  div.className = "screenLogHidden"
}

function showScreenLog() {
  clearTimeout(closeLog)
  var div = document.getElementById("screenLog");
  var currentClass = div.className;
  div.className = "screenLogHidden"

  switch(currentClass) {
    case "screenLogHidden":  div.className = ""; break; 
    case "": div.className = "screenLogHidden"; break;
  }
}

function showLoadingScreen(elm) {
  var div = document.getElementById("loading");
  switch(elm) {
    case true: div.className = "block"; break;
    case false: div.className = "none"; break;

    /*
    case true: div.style.display = "block"; break;
    case false: div.style.display = "none"; break;
    */
  }
}

function createClintInfo(obj) {
  //console.log(obj);
  var keys = getObjKeys(obj);
  for (var i = 0; i < keys.length; i++) {
    if(document.getElementById(keys[i])){
      document.getElementById(keys[i]).innerHTML = obj[keys[i]];
    }
  }
  //document.getElementById("clientInfo").className = "visible";
}

function showElement(elmID, type) {
  switch(type) {
    case true:  cssClass = "block"; break;
    case false: cssClass = "none"; break;
  }

  document.getElementById(elmID).className = cssClass;
} 

function showPopUpElement(elm) {
  var allElements = new Array("deleteUserDetail", "mapping-detail", "user-detail", "file-detail");

  for (var i = 0; i < allElements.length; i++) {
    showElement(allElements[i], false)
  }

  showElement(elm, true)

  setTimeout(function(){ 
    showElement("popup", true);
  }, 10);
}

  // body...

function showStreams(force) {

  var elmBox = document.getElementById("myStreamsBox");
  var elm = document.getElementById("allStreams");
  //console.log(elm);
  show = elm.className;

  switch(force) {
    case true: show = "notVisible"; break;
    case false: show = "visible"; break;
  }

  switch(show) {
    case "notVisible": 
      elm.className = "visible"; 
      elmBox.style.height = "100px";
      break;

    default: 
      elm.className = "notVisible"; 
      elmBox.style.height = "20px";
      break;
  }

  var show = elm.style.display; {
    //console.log(elm.style.display);
  }
  
  calculateWrapperHeight();
}

function xteveBackup() {
  console.log("xteveBackup");
  var data = new Object();
  data["cmd"] = "xteveBackup";

  xTeVe(data);
}

function xteveRestore(elm) {
  var restore = document.createElement("INPUT");
  restore.setAttribute("type", "file");
  restore.setAttribute("class", "notVisible");
  restore.setAttribute("name", "");
  restore.id = "upload";

  document.body.appendChild(restore);
  restore.click();

  restore.onchange = function() {
    var filename = restore.files[0].name
    //console.log(restore.srcElement.files[0]);
    var check = confirm("File: " + filename + "\nAll data will be replaced with those from the backup.\nShould the files be restored?"); 
    if (check == true) {
      var reader  = new FileReader();
      var file = document.querySelector('input[type=file]').files[0];
      if (file) {
        reader.readAsDataURL(file);
        reader.onload = function() {
          console.log(reader.result);
          var data = new Object();
          data["cmd"]     = "xteveRestore"
          data["base64"]  = reader.result

          xTeVe(data);
          return
        };
      } else {
        alert("File could not be loaded")
      }
    }
  };
}

function getBase64(file) {
   var reader = new FileReader();
   reader.readAsDataURL(file);
   reader.onload = function() {
     console.log(reader.result);
   };
   reader.onerror = function(error) {
     console.log('Error: ', error);
   };
}

function logout() {
  document.cookie.split(';').forEach(function(c) {
    document.cookie = c.trim().split('=')[0] + '=;' + 'expires=Thu, 01 Jan 1970 00:00:00 UTC;';
  });
  location.reload();
}

function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}

function setCookie(token) {
  //console.log(token);
  document.cookie = "Token=" + token
}

