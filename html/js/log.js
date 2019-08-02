var logInterval


function updateLog() {
  var data = new Object();
  data["cmd"] = "getLog";
  xTeVe(data);
  writeLogInDiv();
  return
}

function writeLogInDiv() {
  var logs = log["log"];
  var div = document.getElementById("settings").lastChild.lastChild;
  div.innerHTML = "";

  var max = 50;
  
  
  for (var i = 0; i < logs.length; i++) {
    var newEntry = new Object();
    newEntry["_element"]  = "P";

    if (logs[i].includes("ERROR")) {
//      case "warnings":  msgType = "warningMsg"; break;
      newEntry["class"]   = "errorMsg";
    }

    if (logs[i].includes("WARNING")) {
//      case "warnings":  msgType = "warningMsg"; break;
      newEntry["class"]   = "warningMsg";
    }

    newEntry["_text"]     = logs[i];
    
    div.appendChild(createElement(newEntry));
  }

  calculateWrapperHeight();
  var scrollDiv = document.getElementById("box-wrapper");
  scrollDiv.scrollTop = scrollDiv.scrollHeight;
}

function showLog(obj) {
  //logInterval = setInterval(updateLog, 5000);

  var logs = log["log"];

  var div = document.getElementById("settings");

  var newEntry = new Object();
  newEntry["_element"]  = "HR";
  div.appendChild(createElement(newEntry));
  //div = div.lastChild;

  var newEntry = new Object();
  newEntry["_element"]    = "INPUT";
  newEntry["type"]        = "button";
  newEntry["class"]       = "button";
  newEntry["value"]       = "Empty Log";
  newEntry["onclick"]     = "emptyLog()";
  div.appendChild(createElement(newEntry));

  var newEntry = new Object();
  newEntry["_element"]    = "code";
  newEntry["_text"]        = "Update Log: ";
  div.appendChild(createElement(newEntry));

  var newEntry = new Object();
  newEntry["_element"]    = "INPUT";
  newEntry["type"]        = "checkbox";
  //newEntry["checked"]     = "checkbox";
  newEntry["onclick"]     = "logUpdates(this)";
  div.appendChild(createElement(newEntry));

  var newEntry = new Object();
  newEntry["_element"]  = "HR";
  div.appendChild(createElement(newEntry));

  


  
  var newWrapper = new Object();
  newWrapper["_element"]  = "DIV";
  newWrapper["id"]        = "box-wrapper";
  div.appendChild(createElement(newWrapper));
  div = div.lastChild;

  var newPre = new Object();
  newPre["_element"]  = "PRE";
  newPre["id"]        = "logScreen";
  div.appendChild(createElement(newPre));

  div = div.lastChild;

  writeLogInDiv()
  return
}

function emptyLog() {
  var data = new Object();
  data["cmd"] = "emptyLog";
  xTeVe(data);
  return
}

function logUpdates(elm) {
  switch(elm.checked) {
    case false: clearInterval(logInterval); break;
    case true: logInterval = setInterval(updateLog, 5000); break;
    
  }
}