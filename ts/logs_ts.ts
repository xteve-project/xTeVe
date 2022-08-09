class Log {

  createLog(entry:string):any {

    var element = document.createElement("PRE");
    entry = String(entry);

    if (entry.indexOf("WARNING") != -1) {
      element.className = "warningMsg"
    }

    if (entry.indexOf("ERROR") != -1) {
      element.className = "errorMsg"
    }

    if (entry.indexOf("DEBUG") != -1) {
      element.className = "debugMsg"
    }

    element.innerHTML = entry

    return element
  }

}

function showLogs(bottom:boolean) {

  var log = new Log()

  var logs = SERVER["log"]["log"]
  var div = document.getElementById("content_log")

  div.innerHTML = ""

  var keys = getOwnObjProps(logs)

  keys.forEach(logID => {

    var entry = log.createLog(logs[logID])

    div.append(entry)
  
  });

  setTimeout(function(){ 

    if (bottom == true) {
  
      var wrapper = document.getElementById("box-wrapper");
      wrapper.scrollTop = wrapper.scrollHeight;

    }

  }, 10);

}

function resetLogs() {

  var cmd = "resetLogs"
  var data = new Object()
  var server:Server = new Server(cmd)
  server.request(data)

}