var Log = /** @class */ (function () {
    function Log() {
    }
    Log.prototype.createLog = function (entry) {
        var element = document.createElement("PRE");
        if (entry.indexOf("WARNING") != -1) {
            element.className = "warningMsg";
        }
        if (entry.indexOf("ERROR") != -1) {
            element.className = "errorMsg";
        }
        if (entry.indexOf("DEBUG") != -1) {
            element.className = "debugMsg";
        }
        element.innerHTML = entry;
        return element;
    };
    return Log;
}());
function showLogs(bottom) {
    var log = new Log();
    var logs = SERVER["log"]["log"];
    var div = document.getElementById("content_log");
    div.innerHTML = "";
    var keys = getObjKeys(logs);
    keys.forEach(function (logID) {
        var entry = log.createLog(logs[logID]);
        div.append(entry);
    });
    setTimeout(function () {
        if (bottom == true) {
            var wrapper = document.getElementById("box-wrapper");
            wrapper.scrollTop = wrapper.scrollHeight;
        }
    }, 10);
}
function resetLogs() {
    var cmd = "resetLogs";
    var data = new Object();
    var server = new Server(cmd);
    server.request(data);
}
