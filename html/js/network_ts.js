var Server = /** @class */ (function () {
    function Server(cmd) {
        this.cmd = cmd;
    }
    Server.prototype.request = function (data) {
        if (SERVER_CONNECTION == true) {
            return;
        }
        SERVER_CONNECTION = true;
        console.log(data);
        if (this.cmd != "updateLog") {
            showElement("loading", true);
            UNDO = new Object();
        }
        switch (window.location.protocol) {
            case "http:":
                this.protocol = "ws://";
                break;
            case "https:":
                this.protocol = "wss://";
                break;
        }
        var url = this.protocol + window.location.hostname + ":" + window.location.port + "/data/" + "?Token=" + getCookie("Token");
        data["cmd"] = this.cmd;
        var ws = new WebSocket(url);
        ws.onopen = function () {
            WS_AVAILABLE = true;
            console.log("REQUEST (JS):");
            console.log(data);
            console.log("REQUEST: (JSON)");
            console.log(JSON.stringify(data));
            this.send(JSON.stringify(data));
        };
        ws.onerror = function (e) {
            console.log("No websocket connection to xTeVe could be established. Check your network configuration.");
            SERVER_CONNECTION = false;
            if (WS_AVAILABLE == false) {
                alert("No websocket connection to xTeVe could be established. Check your network configuration.");
            }
        };
        ws.onmessage = function (e) {
            SERVER_CONNECTION = false;
            showElement("loading", false);
            console.log("RESPONSE:");
            var response = JSON.parse(e.data);
            console.log(response);
            if (response.hasOwnProperty("token")) {
                document.cookie = "Token=" + response["token"];
            }
            if (response["status"] == false) {
                alert(response["err"]);
                if (response.hasOwnProperty("reload")) {
                    location.reload();
                }
                return;
            }
            if (response.hasOwnProperty("logoURL")) {
                var div = document.getElementById("channel-icon");
                div.value = response["logoURL"];
                div.className = "changed";
                return;
            }
            switch (data["cmd"]) {
                case "updateLog":
                    SERVER["log"] = response["log"];
                    if (document.getElementById("content_log")) {
                        showLogs(false);
                    }
                    return;
                    break;
                default:
                    SERVER = new Object();
                    SERVER = response;
                    break;
            }
            if (response.hasOwnProperty("openMenu")) {
                var menu = document.getElementById(response["openMenu"]);
                menu.click();
                showElement("popup", false);
            }
            if (response.hasOwnProperty("openLink")) {
                window.location = response["openLink"];
            }
            if (response.hasOwnProperty("alert")) {
                alert(response["alert"]);
            }
            if (response.hasOwnProperty("reload")) {
                location.reload();
            }
            if (response.hasOwnProperty("wizard")) {
                createLayout();
                configurationWizard[response["wizard"]].createWizard();
                return;
            }
            createLayout();
        };
    };
    return Server;
}());
function getCookie(name) {
    var value = "; " + document.cookie;
    var parts = value.split("; " + name + "=");
    if (parts.length == 2)
        return parts.pop().split(";").shift();
}
