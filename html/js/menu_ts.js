var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
var MainMenu = /** @class */ (function () {
    function MainMenu() {
        this.DocumentID = "main-menu";
        this.HTMLTag = "LI";
        this.ImagePath = "img/";
    }
    MainMenu.prototype.createIMG = function (src) {
        var element = document.createElement("IMG");
        element.setAttribute("src", this.ImagePath + src);
        return element;
    };
    MainMenu.prototype.createValue = function (value) {
        var element = document.createElement("P");
        element.innerHTML = value;
        return element;
    };
    return MainMenu;
}());
var MainMenuItem = /** @class */ (function (_super) {
    __extends(MainMenuItem, _super);
    function MainMenuItem(menuKey, value, image, headline) {
        var _this = _super.call(this) || this;
        _this.menuKey = menuKey;
        _this.value = value;
        _this.imgSrc = image;
        _this.headline = headline;
        return _this;
    }
    MainMenuItem.prototype.createItem = function () {
        var item = document.createElement("LI");
        item.setAttribute("onclick", "javascript: openThisMenu(this)");
        item.setAttribute("id", this.id);
        var img = this.createIMG(this.imgSrc);
        var value = this.createValue(this.value);
        item.appendChild(img);
        item.appendChild(value);
        var doc = document.getElementById(this.DocumentID);
        doc.appendChild(item);
        switch (this.menuKey) {
            case "playlist":
                this.tableHeader = ["{{.playlist.table.playlist}}", "{{.playlist.table.tuner}}", "{{.playlist.table.lastUpdate}}", "{{.playlist.table.availability}} %", "{{.playlist.table.type}}", "{{.playlist.table.streams}}", "{{.playlist.table.groupTitle}} %", "{{.playlist.table.tvgID}} %", "{{.playlist.table.uniqueID}} %"];
                break;
            case "xmltv":
                this.tableHeader = ["{{.xmltv.table.guide}}", "{{.xmltv.table.lastUpdate}}", "{{.xmltv.table.availability}} %", "{{.xmltv.table.channels}}", "{{.xmltv.table.programs}}"];
                break;
            case "filter":
                this.tableHeader = ["{{.filter.table.name}}", "{{.filter.table.type}}", "{{.filter.table.filter}}"];
                break;
            case "users":
                this.tableHeader = ["{{.users.table.username}}", "{{.users.table.password}}", "{{.users.table.web}}", "{{.users.table.pms}}", "{{.users.table.m3u}}", "{{.users.table.xml}}", "{{.users.table.api}}"];
                break;
            case "mapping":
                this.tableHeader = ["BULK", "{{.mapping.table.chNo}}", "{{.mapping.table.logo}}", "{{.mapping.table.channelName}}", "{{.mapping.table.playlist}}", "{{.mapping.table.groupTitle}}", "{{.mapping.table.xmltvFile}}", "{{.mapping.table.xmltvID}}"];
                break;
        }
        //console.log(this.menuKey, this.tableHeader);
    };
    return MainMenuItem;
}(MainMenu));
var Content = /** @class */ (function () {
    function Content() {
        this.DocumentID = "content";
        this.TableID = "content_table";
        this.headerClass = "content_table_header";
        this.interactionID = "content-interaction";
    }
    Content.prototype.createHeadline = function (value) {
        var element = document.createElement("H3");
        element.innerHTML = value;
        return element;
    };
    Content.prototype.createHR = function () {
        var element = document.createElement("HR");
        return element;
    };
    Content.prototype.createInteraction = function () {
        var element = document.createElement("DIV");
        element.setAttribute("id", this.interactionID);
        return element;
    };
    Content.prototype.createDIV = function () {
        var element = document.createElement("DIV");
        element.id = this.DivID;
        return element;
    };
    Content.prototype.createTABLE = function () {
        var element = document.createElement("TABLE");
        element.id = this.TableID;
        return element;
    };
    Content.prototype.createTableRow = function () {
        var element = document.createElement("TR");
        element.className = this.headerClass;
        return element;
    };
    Content.prototype.createTableContent = function (menuKey) {
        var data = new Object();
        var rows = new Array();
        switch (menuKey) {
            case "playlist":
                var fileTypes = new Array("m3u", "hdhr");
                fileTypes.forEach(function (fileType) {
                    data = SERVER["settings"]["files"][fileType];
                    var keys = getObjKeys(data);
                    keys.forEach(function (key) {
                        var tr = document.createElement("TR");
                        tr.id = key;
                        tr.setAttribute('onclick', 'javascript: openPopUp("' + fileType + '", this)');
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["name"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        if (SERVER["settings"]["buffer"] != "-") {
                            cell.value = data[key]["tuner"];
                        }
                        else {
                            cell.value = "-";
                        }
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["last.update"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["provider.availability"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["type"].toUpperCase();
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["compatibility"]["streams"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["compatibility"]["group.title"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["compatibility"]["tvg.id"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["compatibility"]["stream.id"];
                        tr.appendChild(cell.createCell());
                        rows.push(tr);
                    });
                });
                break;
            case "filter":
                delete SERVER["settings"]["filter"][-1];
                data = SERVER["settings"]["filter"];
                var keys = getObjKeys(data);
                keys.forEach(function (key) {
                    var tr = document.createElement("TR");
                    tr.id = key;
                    tr.setAttribute('onclick', 'javascript: openPopUp("' + data[key]["type"] + '", this)');
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    cell.value = data[key]["name"];
                    tr.appendChild(cell.createCell());
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    switch (data[key]["type"]) {
                        case "custom-filter":
                            cell.value = "{{.filter.custom}}";
                            break;
                        case "group-title":
                            cell.value = "{{.filter.group}}";
                            break;
                        default:
                            break;
                    }
                    tr.appendChild(cell.createCell());
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    cell.value = data[key]["filter"];
                    tr.appendChild(cell.createCell());
                    rows.push(tr);
                });
                break;
            case "xmltv":
                var fileTypes = new Array("xmltv");
                fileTypes.forEach(function (fileType) {
                    data = SERVER["settings"]["files"][fileType];
                    var keys = getObjKeys(data);
                    keys.forEach(function (key) {
                        var tr = document.createElement("TR");
                        tr.id = key;
                        tr.setAttribute('onclick', 'javascript: openPopUp("' + fileType + '", this)');
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["name"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["last.update"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["provider.availability"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["compatibility"]["xmltv.channels"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["compatibility"]["xmltv.programs"];
                        tr.appendChild(cell.createCell());
                        rows.push(tr);
                    });
                });
                break;
            case "users":
                var fileTypes = new Array("users");
                fileTypes.forEach(function (fileType) {
                    data = SERVER[fileType];
                    var keys = getObjKeys(data);
                    keys.forEach(function (key) {
                        var tr = document.createElement("TR");
                        tr.id = key;
                        tr.setAttribute('onclick', 'javascript: openPopUp("' + fileType + '", this)');
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = data[key]["data"]["username"];
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        cell.value = "******";
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        if (data[key]["data"]["authentication.web"] == true) {
                            cell.value = "✓";
                        }
                        else {
                            cell.value = "-";
                        }
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        if (data[key]["data"]["authentication.pms"] == true) {
                            cell.value = "✓";
                        }
                        else {
                            cell.value = "-";
                        }
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        if (data[key]["data"]["authentication.m3u"] == true) {
                            cell.value = "✓";
                        }
                        else {
                            cell.value = "-";
                        }
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        if (data[key]["data"]["authentication.xml"] == true) {
                            cell.value = "✓";
                        }
                        else {
                            cell.value = "-";
                        }
                        tr.appendChild(cell.createCell());
                        var cell = new Cell();
                        cell.child = true;
                        cell.childType = "P";
                        if (data[key]["data"]["authentication.api"] == true) {
                            cell.value = "✓";
                        }
                        else {
                            cell.value = "-";
                        }
                        tr.appendChild(cell.createCell());
                        rows.push(tr);
                    });
                });
                break;
            case "mapping":
                BULK_EDIT = false;
                createSearchObj();
                checkUndo("epgMapping");
                console.log("MAPPING");
                data = SERVER["xepg"]["epgMapping"];
                var keys = getObjKeys(data);
                keys.forEach(function (key) {
                    var tr = document.createElement("TR");
                    tr.id = key;
                    //tr.setAttribute('oncontextmenu', 'javascript: rightClick(this)')
                    switch (data[key]["x-active"]) {
                        case true:
                            tr.className = "activeEPG";
                            break;
                        case false:
                            tr.className = "notActiveEPG";
                            break;
                    }
                    // Bulk
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "BULK";
                    cell.value = false;
                    tr.appendChild(cell.createCell());
                    // Kanalnummer
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "INPUTCHANNEL";
                    cell.value = data[key]["x-channelID"];
                    //td.setAttribute('onclick', 'javascript: changeChannelNumber("' + key + '", this)')
                    tr.appendChild(cell.createCell());
                    // Logo
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "IMG";
                    cell.imageURL = data[key]["tvg-logo"];
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    // Kanalname
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    cell.className = data[key]["x-category"];
                    cell.value = data[key]["x-name"];
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    // Playlist
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    //cell.value = data[key]["_file.m3u.name"] 
                    cell.value = getValueFromProviderFile(data[key]["_file.m3u.id"], "m3u", "name");
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    // Gruppe (group-title)
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    cell.value = data[key]["x-group-title"];
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    // XMLTV Datei
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    if (data[key]["x-xmltv-file"] != "-") {
                        cell.value = getValueFromProviderFile(data[key]["x-xmltv-file"], "xmltv", "name");
                    }
                    else {
                        cell.value = data[key]["x-xmltv-file"];
                    }
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    // XMLTV Kanal
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    //var value = str.substring(1, 4);
                    var value = data[key]["x-mapping"];
                    if (value.length > 20) {
                        value = data[key]["x-mapping"].substring(0, 20) + "...";
                    }
                    cell.value = value;
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    rows.push(tr);
                });
                break;
            case "settings":
                alert();
                break;
            default:
                console.log("Table content (menuKey):", menuKey);
                break;
        }
        return rows;
    };
    return Content;
}());
var Cell = /** @class */ (function () {
    function Cell() {
    }
    Cell.prototype.createCell = function () {
        var td = document.createElement("TD");
        if (this.child == true) {
            var element;
            switch (this.childType) {
                case "P":
                    element = document.createElement(this.childType);
                    element.innerHTML = this.value;
                    element.className = this.className;
                    break;
                case "INPUT":
                    element = document.createElement(this.childType);
                    element.value = this.value;
                    element.type = "text";
                    break;
                case "INPUTCHANNEL":
                    element = document.createElement("INPUT");
                    element.setAttribute("onchange", "javscript: changeChannelNumber(this)");
                    element.value = this.value;
                    element.type = "text";
                    break;
                case "BULK":
                    element = document.createElement("INPUT");
                    element.checked = this.value;
                    element.type = "checkbox";
                    element.className = "bulk hideBulk";
                    break;
                case "BULK_HEAD":
                    element = document.createElement("INPUT");
                    element.checked = this.value;
                    element.type = "checkbox";
                    element.className = "bulk hideBulk";
                    element.setAttribute("onclick", "javascript: selectAllChannels()");
                    break;
                case "IMG":
                    element = document.createElement(this.childType);
                    element.setAttribute("src", this.imageURL);
                    if (this.imageURL != "") {
                        element.setAttribute("onerror", "javascript: this.onerror=null;this.src=''");
                        //onerror="this.onerror=null;this.src='missing.gif';"
                    }
            }
            td.appendChild(element);
        }
        else {
            td.innerHTML = this.value;
        }
        if (this.onclick == true) {
            td.setAttribute("onclick", this.onclickFunktion);
            td.className = "pointer";
        }
        if (this.tdClassName != undefined) {
            td.className = this.tdClassName;
        }
        return td;
    };
    return Cell;
}());
var ShowContent = /** @class */ (function (_super) {
    __extends(ShowContent, _super);
    function ShowContent(menuID) {
        var _this = _super.call(this) || this;
        _this.menuID = menuID;
        return _this;
    }
    ShowContent.prototype.createInput = function (type, name, value) {
        var input = document.createElement("INPUT");
        input.setAttribute("type", type);
        input.setAttribute("name", name);
        input.setAttribute("value", value);
        return input;
    };
    ShowContent.prototype.show = function () {
        COLUMN_TO_SORT = -1;
        // Alten Inhalt löschen
        var doc = document.getElementById(this.DocumentID);
        doc.innerHTML = "";
        showPreview(false);
        // Überschrift
        var headline = menuItems[this.menuID].headline;
        var menuKey = menuItems[this.menuID].menuKey;
        var h = this.createHeadline(headline);
        doc.appendChild(h);
        var hr = this.createHR();
        doc.appendChild(hr);
        // Interaktion
        var div = this.createInteraction();
        doc.appendChild(div);
        var interaction = document.getElementById(this.interactionID);
        switch (menuKey) {
            case "playlist":
                var input = this.createInput("button", menuKey, "{{.button.new}}");
                input.setAttribute("id", "-");
                input.setAttribute("onclick", 'javascript: openPopUp("playlist")');
                interaction.appendChild(input);
                break;
            case "filter":
                var input = this.createInput("button", menuKey, "{{.button.new}}");
                input.setAttribute("id", -1);
                input.setAttribute("onclick", 'javascript: openPopUp("filter", this)');
                interaction.appendChild(input);
                break;
            case "xmltv":
                var input = this.createInput("button", menuKey, "{{.button.new}}");
                input.setAttribute("id", "xmltv");
                input.setAttribute("onclick", 'javascript: openPopUp("xmltv")');
                interaction.appendChild(input);
                break;
            case "users":
                var input = this.createInput("button", menuKey, "{{.button.new}}");
                input.setAttribute("id", "users");
                input.setAttribute("onclick", 'javascript: openPopUp("users")');
                interaction.appendChild(input);
                break;
            case "mapping":
                showElement("loading", true);
                var input = this.createInput("button", menuKey, "{{.button.save}}");
                input.setAttribute("onclick", 'javascript: savePopupData("mapping", "", "")');
                interaction.appendChild(input);
                var input = this.createInput("button", menuKey, "{{.button.bulkEdit}}");
                input.setAttribute("onclick", 'javascript: bulkEdit()');
                interaction.appendChild(input);
                var input = this.createInput("search", "search", "");
                input.setAttribute("id", "searchMapping");
                input.setAttribute("placeholder", "{{.button.search}}");
                input.className = "search";
                input.setAttribute("onchange", 'javascript: searchInMapping()');
                interaction.appendChild(input);
                break;
            case "settings":
                var input = this.createInput("button", menuKey, "{{.button.save}}");
                input.setAttribute("onclick", 'javascript: saveSettings();');
                interaction.appendChild(input);
                var input = this.createInput("button", menuKey, "{{.button.backup}}");
                input.setAttribute("onclick", 'javascript: backup();');
                interaction.appendChild(input);
                var input = this.createInput("button", menuKey, "{{.button.restore}}");
                input.setAttribute("onclick", 'javascript: restore();');
                interaction.appendChild(input);
                var wrapper = document.createElement("DIV");
                wrapper.setAttribute("id", "box-wrapper");
                doc.appendChild(wrapper);
                this.DivID = "content_settings";
                var settings = this.createDIV();
                wrapper.appendChild(settings);
                showSettings();
                return;
                break;
            case "log":
                var input = this.createInput("button", menuKey, "{{.button.resetLogs}}");
                input.setAttribute("onclick", 'javascript: resetLogs();');
                interaction.appendChild(input);
                var wrapper = document.createElement("DIV");
                wrapper.setAttribute("id", "box-wrapper");
                doc.appendChild(wrapper);
                this.DivID = "content_log";
                var logs = this.createDIV();
                wrapper.appendChild(logs);
                showLogs(true);
                return;
                break;
            case "logout":
                location.reload();
                document.cookie = "Token= ; expires = Thu, 01 Jan 1970 00:00:00 GMT";
                break;
            default:
                console.log("Show content (menuKey):", menuKey);
                break;
        }
        // Tabelle erstellen (falls benötigt)
        var tableHeader = menuItems[this.menuID].tableHeader;
        if (tableHeader.length > 0) {
            var wrapper = document.createElement("DIV");
            doc.appendChild(wrapper);
            wrapper.setAttribute("id", "box-wrapper");
            var table = this.createTABLE();
            wrapper.appendChild(table);
            var header = this.createTableRow();
            table.appendChild(header);
            // Kopfzeile der Tablle
            tableHeader.forEach(function (element) {
                var cell = new Cell();
                cell.child = true;
                cell.childType = "P";
                cell.value = element;
                if (element == "BULK") {
                    cell.childType = "BULK_HEAD";
                    cell.value = false;
                }
                if (menuKey == "mapping") {
                    if (element == "{{.mapping.table.chNo}}") {
                        cell.onclick = true;
                        cell.onclickFunktion = "javascript: sortTable(1);";
                        cell.tdClassName = "sortThis";
                    }
                    if (element == "{{.mapping.table.channelName}}") {
                        cell.onclick = true;
                        cell.onclickFunktion = "javascript: sortTable(3);";
                    }
                    if (element == "{{.mapping.table.playlist}}") {
                        cell.onclick = true;
                        cell.onclickFunktion = "javascript: sortTable(4);";
                    }
                    if (element == "{{.mapping.table.groupTitle}}") {
                        cell.onclick = true;
                        cell.onclickFunktion = "javascript: sortTable(5);";
                    }
                }
                header.appendChild(cell.createCell());
            });
            table.appendChild(header);
            // Inhalt der Tabelle
            var rows = this.createTableContent(menuKey);
            rows.forEach(function (tr) {
                table.appendChild(tr);
            });
        }
        switch (menuKey) {
            case "mapping":
                sortTable(1);
                break;
            case "filter":
                showPreview(true);
                sortTable(0);
                break;
            default:
                COLUMN_TO_SORT = -1;
                sortTable(0);
                break;
        }
        showElement("loading", false);
    };
    return ShowContent;
}(Content));
function PageReady() {
    var server = new Server("getServerConfig");
    server.request(new Object());
    window.addEventListener("resize", function () {
        calculateWrapperHeight();
    }, true);
    setInterval(function () {
        updateLog();
    }, 10000);
    return;
}
function createLayout() {
    // Client Info
    var obj = SERVER["clientInfo"];
    var keys = getObjKeys(obj);
    for (var i = 0; i < keys.length; i++) {
        if (document.getElementById(keys[i])) {
            document.getElementById(keys[i]).innerHTML = obj[keys[i]];
        }
    }
    if (!document.getElementById("main-menu")) {
        return;
    }
    // Menü erstellen
    document.getElementById("main-menu").innerHTML = "";
    for (var i_1 = 0; i_1 < menuItems.length; i_1++) {
        menuItems[i_1].id = i_1;
        switch (menuItems[i_1]["menuKey"]) {
            case "users":
            case "logout":
                if (SERVER["settings"]["authentication.web"] == true) {
                    menuItems[i_1].createItem();
                }
                break;
            case "mapping":
            case "xmltv":
                if (SERVER["clientInfo"]["epgSource"] == "XEPG") {
                    menuItems[i_1].createItem();
                }
                break;
            default:
                menuItems[i_1].createItem();
                break;
        }
    }
    return;
}
function openThisMenu(element) {
    var id = element.id;
    var content = new ShowContent(id);
    content.show();
    calculateWrapperHeight();
    return;
}
var PopupWindow = /** @class */ (function () {
    function PopupWindow() {
        this.DocumentID = "popup-custom";
        this.InteractionID = "interaction";
        this.doc = document.getElementById(this.DocumentID);
    }
    PopupWindow.prototype.createTitle = function (title) {
        var td = document.createElement("TD");
        td.className = "left";
        td.innerHTML = title + ":";
        return td;
    };
    PopupWindow.prototype.createContent = function (element) {
        var td = document.createElement("TD");
        td.appendChild(element);
        return td;
    };
    PopupWindow.prototype.createInteraction = function () {
        var div = document.createElement("div");
        div.setAttribute("id", "popup-interaction");
        div.className = "interaction";
        this.doc.appendChild(div);
    };
    return PopupWindow;
}());
var PopupContent = /** @class */ (function (_super) {
    __extends(PopupContent, _super);
    function PopupContent() {
        var _this = _super !== null && _super.apply(this, arguments) || this;
        _this.table = document.createElement("TABLE");
        return _this;
    }
    PopupContent.prototype.createHeadline = function (headline) {
        this.doc.innerHTML = "";
        var element = document.createElement("H3");
        element.innerHTML = headline.toUpperCase();
        this.doc.appendChild(element);
        // Tabelle erstellen
        this.table = document.createElement("TABLE");
        this.doc.appendChild(this.table);
    };
    PopupContent.prototype.appendRow = function (title, element) {
        var tr = document.createElement("TR");
        // Bezeichnung
        if (title.length != 0) {
            tr.appendChild(this.createTitle(title));
        }
        // Content
        tr.appendChild(this.createContent(element));
        this.table.appendChild(tr);
    };
    PopupContent.prototype.createInput = function (type, name, value) {
        var input = document.createElement("INPUT");
        if (value == undefined) {
            value = "";
        }
        input.setAttribute("type", type);
        input.setAttribute("name", name);
        input.setAttribute("value", value);
        return input;
    };
    PopupContent.prototype.createCheckbox = function (name) {
        var input = document.createElement("INPUT");
        input.setAttribute("type", "checkbox");
        input.setAttribute("name", name);
        return input;
    };
    PopupContent.prototype.createSelect = function (text, values, set, dbKey) {
        var select = document.createElement("SELECT");
        select.setAttribute("name", dbKey);
        for (var i = 0; i < text.length; i++) {
            var option = document.createElement("OPTION");
            option.setAttribute("value", values[i]);
            option.innerText = text[i];
            select.appendChild(option);
        }
        if (set != "") {
            select.value = set;
        }
        if (set == undefined) {
            select.value = values[0];
        }
        return select;
    };
    PopupContent.prototype.selectOption = function (select, value) {
        //select.selectedOptions = value
        var s = select;
        s.options[s.selectedIndex].value = value;
        return select;
    };
    PopupContent.prototype.description = function (value) {
        var tr = document.createElement("TR");
        var td = document.createElement("TD");
        var span = document.createElement("PRE");
        span.innerHTML = value;
        tr.appendChild(td);
        tr.appendChild(this.createContent(span));
        this.table.appendChild(tr);
    };
    // Interaktion
    PopupContent.prototype.addInteraction = function (element) {
        var interaction = document.getElementById("popup-interaction");
        interaction.appendChild(element);
    };
    return PopupContent;
}(PopupWindow));
function openPopUp(dataType, element) {
    var data = new Object();
    var id;
    switch (element) {
        case undefined:
            switch (dataType) {
                case "group-title":
                    if (id == undefined) {
                        id = -1;
                    }
                    data = getLocalData("filter", id);
                    data["type"] = "group-title";
                    break;
                case "custom-filter":
                    if (id == undefined) {
                        id = -1;
                    }
                    data = getLocalData("filter", id);
                    data["type"] = "custom-filter";
                    break;
                default:
                    data["id.provider"] = "-";
                    data["type"] = dataType;
                    id = "-";
                    break;
            }
            break;
        default:
            id = element.id;
            data = getLocalData(dataType, id);
            break;
    }
    var content = new PopupContent();
    switch (dataType) {
        case "playlist":
            content.createHeadline("{{.playlist.playlistType.title}}");
            // Type
            var text = ["M3U", "HDHomeRun"];
            var values = ["javascript: openPopUp('m3u')", "javascript: openPopUp('hdhr')"];
            var select = content.createSelect(text, values, "", "type");
            select.setAttribute("id", "type");
            select.setAttribute("onchange", 'javascript: changeButtonAction(this, "next", "onclick")'); // changeButtonAction
            content.appendRow("{{.playlist.type.title}}", select);
            // Interaktion
            content.createInteraction();
            // Abbrechen
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Weiter
            var input = content.createInput("button", "next", "{{.button.next}}");
            input.setAttribute("onclick", 'javascript: openPopUp("m3u")');
            input.setAttribute("id", 'next');
            content.addInteraction(input);
            break;
        case "m3u":
            content.createHeadline(dataType);
            // Name
            var dbKey = "name";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.playlist.name.placeholder}}");
            content.appendRow("{{.playlist.name.title}}", input);
            // Beschreibung
            var dbKey = "description";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.playlist.description.placeholder}}");
            content.appendRow("{{.playlist.description.title}}", input);
            // URL
            var dbKey = "file.source";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.playlist.fileM3U.placeholder}}");
            content.appendRow("{{.playlist.fileM3U.title}}", input);
            // Tuner
            if (SERVER["settings"]["buffer"] != "-") {
                var text = new Array();
                var values = new Array();
                for (var i = 1; i <= 100; i++) {
                    text.push(i.toString());
                    values.push(i.toString());
                }
                var dbKey = "tuner";
                var select = content.createSelect(text, values, data[dbKey], dbKey);
                select.setAttribute("onfocus", "javascript: return;");
                content.appendRow("{{.playlist.tuner.title}}", select);
            }
            else {
                var dbKey = "tuner";
                if (data[dbKey] == undefined) {
                    data[dbKey] = 1;
                }
                var input = content.createInput("text", dbKey, data[dbKey]);
                input.setAttribute("readonly", "true");
                input.className = "notAvailable";
                content.appendRow("{{.playlist.tuner.title}}", input);
            }
            content.description("{{.playlist.tuner.description}}");
            // Interaktion
            content.createInteraction();
            // Löschen
            if (data["id.provider"] != "-") {
                var input = content.createInput("button", "delete", "{{.button.delete}}");
                input.className = "delete";
                input.setAttribute('onclick', 'javascript: savePopupData("m3u", "' + id + '", true, 0)');
                content.addInteraction(input);
            }
            else {
                var input = content.createInput("button", "back", "{{.button.back}}");
                input.setAttribute("onclick", 'javascript: openPopUp("playlist")');
                content.addInteraction(input);
            }
            // Abbrechen
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Aktualisieren
            if (data["id.provider"] != "-") {
                var input = content.createInput("button", "update", "{{.button.update}}");
                input.setAttribute('onclick', 'javascript: savePopupData("m3u", "' + id + '", false, 1)');
                content.addInteraction(input);
            }
            // Speichern
            var input = content.createInput("button", "save", "{{.button.save}}");
            input.setAttribute('onclick', 'javascript: savePopupData("m3u", "' + id + '", false, 0)');
            content.addInteraction(input);
            break;
        case "hdhr":
            content.createHeadline(dataType);
            // Name
            var dbKey = "name";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.playlist.name.placeholder}}");
            content.appendRow("{{.playlist.name.title}}", input);
            // Beschreibung
            var dbKey = "description";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.playlist.description.placeholder}}");
            content.appendRow("{{.playlist.description.placeholder}}", input);
            // URL
            var dbKey = "file.source";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.playlist.fileHDHR.placeholder}}");
            content.appendRow("{{.playlist.fileHDHR.title}}", input);
            // Tuner
            if (SERVER["settings"]["buffer"] != "-") {
                var text = new Array();
                var values = new Array();
                for (var i = 1; i <= 100; i++) {
                    text.push(i.toString());
                    values.push(i.toString());
                }
                var dbKey = "tuner";
                var select = content.createSelect(text, values, data[dbKey], dbKey);
                select.setAttribute("onfocus", "javascript: return;");
                content.appendRow("{{.playlist.tuner.title}}", select);
            }
            else {
                var dbKey = "tuner";
                if (data[dbKey] == undefined) {
                    data[dbKey] = 1;
                }
                var input = content.createInput("text", dbKey, data[dbKey]);
                input.setAttribute("readonly", "true");
                input.className = "notAvailable";
                content.appendRow("{{.playlist.tuner.title}}", input);
            }
            content.description("{{.playlist.tuner.description}}");
            // Interaktion
            content.createInteraction();
            // Löschen
            if (data["id.provider"] != "-") {
                var input = content.createInput("button", "delete", "{{.button.delete}}");
                input.setAttribute('onclick', 'javascript: savePopupData("hdhr", "' + id + '", true, 0)');
                input.className = "delete";
                content.addInteraction(input);
            }
            else {
                var input = content.createInput("button", "back", "{{.button.back}}");
                input.setAttribute("onclick", 'javascript: openPopUp("playlist")');
                content.addInteraction(input);
            }
            // Abbrechen
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Aktualisieren
            if (data["id.provider"] != "-") {
                var input = content.createInput("button", "update", "{{.button.update}}");
                input.setAttribute('onclick', 'javascript: savePopupData("hdhr", "' + id + '", false, 1)');
                content.addInteraction(input);
            }
            // Speichern
            var input = content.createInput("button", "save", "{{.button.save}}");
            input.setAttribute('onclick', 'javascript: savePopupData("hdhr", "' + id + '", false, 0)');
            content.addInteraction(input);
            break;
        case "filter":
            content.createHeadline(dataType);
            // Type
            var dbKey = "type";
            var text = ["M3U: " + "{{.filter.type.groupTitle}}", "xTeVe: " + "{{.filter.type.customFilter}}"];
            var values = ["javascript: openPopUp('group-title')", "javascript: openPopUp('custom-filter')"];
            var select = content.createSelect(text, values, "javascript: openPopUp('group-title')", dbKey);
            select.setAttribute("id", id);
            select.setAttribute("onchange", 'javascript: changeButtonAction(this, "next", "onclick");'); // changeButtonAction
            content.appendRow("{{.filter.type.title}}", select);
            // Interaktion
            content.createInteraction();
            // Abbrechen
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Weiter
            var input = content.createInput("button", "next", "{{.button.next}}");
            input.setAttribute("onclick", 'javascript: openPopUp("group-title")');
            input.setAttribute("id", 'next');
            content.addInteraction(input);
            break;
        case "custom-filter":
        case "group-title":
            switch (dataType) {
                case "custom-filter":
                    content.createHeadline("{{.filter.custom}}");
                    break;
                case "group-title":
                    content.createHeadline("{{.filter.group}}");
                    break;
            }
            // Name      
            var dbKey = "name";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.filter.name.placeholder}}");
            content.appendRow("{{.filter.name.title}}", input);
            // Beschreibung
            var dbKey = "description";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.filter.description.placeholder}}");
            content.appendRow("{{.filter.description.title}}", input);
            // Typ
            var dbKey = "type";
            var input = content.createInput("hidden", dbKey, data[dbKey]);
            content.appendRow("", input);
            var filterType = data[dbKey];
            switch (filterType) {
                case "custom-filter":
                    // Groß- Kleinschreibung beachten
                    var dbKey = "caseSensitive";
                    var input = content.createCheckbox(dbKey);
                    input.checked = data[dbKey];
                    content.appendRow("{{.filter.caseSensitive.title}}", input);
                    // Filterregel (Benutzerdefiniert)
                    var dbKey = "filter";
                    var input = content.createInput("text", dbKey, data[dbKey]);
                    input.setAttribute("placeholder", "{{.filter.filterRule.placeholder}}");
                    content.appendRow("{{.filter.filterRule.title}}", input);
                    break;
                case "group-title":
                    //alert(dbKey + " " + filterType)
                    // Filter basierend auf den Gruppen in der M3U
                    var dbKey = "filter";
                    var groupsM3U = getLocalData("m3uGroups", "");
                    var text = groupsM3U["text"];
                    var values = groupsM3U["value"];
                    var select = content.createSelect(text, values, data[dbKey], dbKey);
                    select.setAttribute("onchange", "javascript: this.className = 'changed'");
                    content.appendRow("{{.filter.filterGroup.title}}", select);
                    content.description("{{.filter.filterGroup.description}}");
                    // Groß- Kleinschreibung beachten
                    var dbKey = "caseSensitive";
                    var input = content.createCheckbox(dbKey);
                    input.checked = data[dbKey];
                    content.appendRow("{{.filter.caseSensitive.title}}", input);
                    var dbKey = "include";
                    var input = content.createInput("text", dbKey, data[dbKey]);
                    input.setAttribute("placeholder", "{{.filter.include.placeholder}}");
                    content.appendRow("{{.filter.include.title}}", input);
                    content.description("{{.filter.include.description}}");
                    var dbKey = "exclude";
                    var input = content.createInput("text", dbKey, data[dbKey]);
                    input.setAttribute("placeholder", "{{.filter.exclude.placeholder}}");
                    content.appendRow("{{.filter.exclude.title}}", input);
                    content.description("{{.filter.exclude.description}}");
                    break;
                default:
                    break;
            }
            // Interaktion
            content.createInteraction();
            // Löschen
            var input = content.createInput("button", "delete", "{{.button.delete}}");
            input.setAttribute('onclick', 'javascript: savePopupData("filter", "' + id + '", true, 0)');
            input.className = "delete";
            content.addInteraction(input);
            // Abbrechen
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Speichern
            var input = content.createInput("button", "save", "{{.button.save}}");
            input.setAttribute('onclick', 'javascript: savePopupData("filter", "' + id + '", false, 0)');
            content.addInteraction(input);
            break;
        case "xmltv":
            content.createHeadline(dataType);
            // Name
            var dbKey = "name";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.xmltv.name.placeholder}}");
            content.appendRow("{{.xmltv.name.title}}", input);
            // Beschreibung
            var dbKey = "description";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.xmltv.description.placeholder}}");
            content.appendRow("{{.xmltv.description.title}}", input);
            // URL
            var dbKey = "file.source";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.xmltv.fileXMLTV.placeholder}}");
            content.appendRow("{{.xmltv.fileXMLTV.title}}", input);
            // Interaktion
            content.createInteraction();
            // Löschen
            if (data["id.provider"] != "-") {
                var input = content.createInput("button", "delete", "{{.button.delete}}");
                input.setAttribute('onclick', 'javascript: savePopupData("xmltv", "' + id + '", true, 0)');
                input.className = "delete";
                content.addInteraction(input);
            }
            // Abbrechen
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Aktualisieren
            if (data["id.provider"] != "-") {
                var input = content.createInput("button", "update", "{{.button.update}}");
                input.setAttribute('onclick', 'javascript: savePopupData("xmltv", "' + id + '", false, 1)');
                content.addInteraction(input);
            }
            // Speichern
            var input = content.createInput("button", "save", "{{.button.save}}");
            input.setAttribute('onclick', 'javascript: savePopupData("xmltv", "' + id + '", false, 0)');
            content.addInteraction(input);
            break;
        case "users":
            content.createHeadline("{{.mainMenu.item.users}}");
            // Benutzername 
            var dbKey = "username";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.users.username.placeholder}}");
            content.appendRow("{{.users.username.title}}", input);
            // Neues Passwort 
            var dbKey = "password";
            var input = content.createInput("password", dbKey, "");
            input.setAttribute("placeholder", "{{.users.password.placeholder}}");
            content.appendRow("{{.users.password.title}}", input);
            // Bestätigung 
            var dbKey = "confirm";
            var input = content.createInput("password", dbKey, "");
            input.setAttribute("placeholder", "{{.users.confirm.placeholder}}");
            content.appendRow("{{.users.confirm.title}}", input);
            // Berechtigung WEB
            var dbKey = "authentication.web";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            if (data["defaultUser"] == true) {
                input.setAttribute("onclick", "javascript: return false");
            }
            content.appendRow("{{.users.web.title}}", input);
            // Berechtigung PMS
            var dbKey = "authentication.pms";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            content.appendRow("{{.users.pms.title}}", input);
            // Berechtigung M3U
            var dbKey = "authentication.m3u";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            content.appendRow("{{.users.m3u.title}}", input);
            // Berechtigung XML
            var dbKey = "authentication.xml";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            content.appendRow("{{.users.xml.title}}", input);
            // Berechtigung API
            var dbKey = "authentication.api";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            content.appendRow("{{.users.api.title}}", input);
            // Interaktion
            content.createInteraction();
            // Löschen
            if (data["defaultUser"] != true && id != "-") {
                var input = content.createInput("button", "delete", "{{.button.delete}}");
                input.className = "delete";
                input.setAttribute('onclick', 'javascript: savePopupData("' + dataType + '", "' + id + '", true, 0)');
                content.addInteraction(input);
            }
            // Abbrechen
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Speichern
            var input = content.createInput("button", "save", "{{.button.save}}");
            input.setAttribute("onclick", 'javascript: savePopupData("' + dataType + '", "' + id + '", "false");');
            content.addInteraction(input);
            break;
        case "mapping":
            content.createHeadline("{{.mainMenu.item.mapping}}");
            // Aktiv 
            var dbKey = "x-active";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            input.id = "active";
            //input.setAttribute("onchange", "javascript: this.className = 'changed'")
            input.setAttribute("onchange", "javascript: toggleChannelStatus('" + id + "', this)");
            content.appendRow("{{.mapping.active.title}}", input);
            // Kanalname 
            var dbKey = "x-name";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            if (BULK_EDIT == true) {
                input.style.border = "solid 1px red";
                input.setAttribute("readonly", "true");
            }
            content.appendRow("{{.mapping.channelName.title}}", input);
            content.description(data["name"]);
            // Beschreibung 
            var dbKey = "x-description";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.mapping.description.placeholder}}");
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            content.appendRow("{{.mapping.description.title}}", input);
            // Aktualisierung des Kanalnamens
            if (data.hasOwnProperty("_uuid.key")) {
                if (data["_uuid.key"] != "") {
                    var dbKey = "x-update-channel-name";
                    var input = content.createCheckbox(dbKey);
                    input.setAttribute("onchange", "javascript: this.className = 'changed'");
                    input.checked = data[dbKey];
                    content.appendRow("{{.mapping.updateChannelName.title}}", input);
                }
            }
            // Logo URL (Kanal) 
            var dbKey = "tvg-logo";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            input.setAttribute("id", "channel-icon");
            content.appendRow("{{.mapping.channelLogo.title}}", input);
            // Aktualisierung des Kanallogos
            var dbKey = "x-update-channel-icon";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            input.setAttribute("id", "update-icon");
            input.setAttribute("onchange", "javascript: this.className = 'changed'; changeChannelLogo('" + id + "');");
            content.appendRow("{{.mapping.updateChannelLogo.title}}", input);
            // Erweitern der EPG Kategorie
            var dbKey = "x-category";
            var text = ["-", "Kids (Emby only)", "News", "Movie", "Series", "Sports"];
            var values = ["", "Kids", "News", "Movie", "Series", "Sports"];
            var select = content.createSelect(text, values, data[dbKey], dbKey);
            select.setAttribute("onchange", "javascript: this.className = 'changed'");
            content.appendRow("{{.mapping.epgCategory.title}}", select);
            // M3U Gruppentitel
            var dbKey = "x-group-title";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            content.appendRow("{{.mapping.m3uGroupTitle.title}}", input);
            if (data["group-title"] != undefined) {
                content.description(data["group-title"]);
            }
            // XMLTV Datei
            var dbKey = "x-xmltv-file";
            var xmlFile = data[dbKey];
            var xmltv = new XMLTVFile();
            var select = xmltv.getFiles(data[dbKey]);
            select.setAttribute("name", dbKey);
            select.setAttribute("id", "popup-xmltv");
            select.setAttribute("onchange", "javascript: this.className = 'changed'; setXmltvChannel('" + id + "',this);");
            content.appendRow("{{.mapping.xmltvFile.title}}", select);
            var file = data[dbKey];
            // XMLTV Mapping
            var dbKey = "x-mapping";
            var xmltv = new XMLTVFile();
            var select = xmltv.getPrograms(file, data[dbKey]);
            select.setAttribute("name", dbKey);
            select.setAttribute("id", "popup-mapping");
            select.setAttribute("onchange", "javascript: this.className = 'changed'; checkXmltvChannel('" + id + "',this,'" + xmlFile + "');");
            sortSelect(select);
            content.appendRow("{{.mapping.xmltvChannel.title}}", select);
            // Interaktion
            content.createInteraction();
            // Logo hochladen
            var input = content.createInput("button", "cancel", "{{.button.uploadLogo}}");
            input.setAttribute("onclick", 'javascript: uploadLogo();');
            content.addInteraction(input);
            // Abbrechen
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Fertig
            var ids = new Array();
            ids = getAllSelectedChannels();
            if (ids.length == 0) {
                ids.push(id);
            }
            var input = content.createInput("button", "save", "{{.button.done}}");
            input.setAttribute("onclick", 'javascript: donePopupData("' + dataType + '", "' + ids + '", "false");');
            content.addInteraction(input);
            break;
        default:
            break;
    }
    showPopUpElement('popup-custom');
}
var XMLTVFile = /** @class */ (function () {
    function XMLTVFile() {
    }
    XMLTVFile.prototype.getFiles = function (set) {
        var fileIDs = getObjKeys(SERVER["xepg"]["xmltvMap"]);
        var values = new Array("-");
        var text = new Array("-");
        for (var i = 0; i < fileIDs.length; i++) {
            if (fileIDs[i] != "xTeVe Dummy") {
                values.push(getValueFromProviderFile(fileIDs[i], "xmltv", "file.xteve"));
                text.push(getValueFromProviderFile(fileIDs[i], "xmltv", "name"));
            }
            else {
                values.push(fileIDs[i]);
                text.push(fileIDs[i]);
            }
        }
        var select = document.createElement("SELECT");
        for (var i = 0; i < text.length; i++) {
            var option = document.createElement("OPTION");
            option.setAttribute("value", values[i]);
            option.innerText = text[i];
            select.appendChild(option);
        }
        if (set != "") {
            select.value = set;
        }
        return select;
    };
    XMLTVFile.prototype.getPrograms = function (file, set) {
        //var fileIDs:string[] = getObjKeys(SERVER["xepg"]["xmltvMap"])
        var values = getObjKeys(SERVER["xepg"]["xmltvMap"][file]);
        var text = new Array();
        var displayName;
        for (var i = 0; i < values.length; i++) {
            if (SERVER["xepg"]["xmltvMap"][file][values[i]].hasOwnProperty('display-name') == true) {
                displayName = SERVER["xepg"]["xmltvMap"][file][values[i]]["display-name"];
            }
            else {
                displayName = "-";
            }
            text[i] = displayName + " (" + values[i] + ")";
        }
        text.unshift("-");
        values.unshift("-");
        var select = document.createElement("SELECT");
        for (var i = 0; i < text.length; i++) {
            var option = document.createElement("OPTION");
            option.setAttribute("value", values[i]);
            option.innerText = text[i];
            select.appendChild(option);
        }
        if (set != "") {
            select.value = set;
        }
        if (select.value != set) {
            select.value = "-";
        }
        return select;
    };
    return XMLTVFile;
}());
function getValueFromProviderFile(file, fileType, key) {
    if (file == "xTeVe Dummy") {
        return file;
    }
    var fileID;
    var indicator = file.charAt(0);
    switch (indicator) {
        case "M":
            fileType = "m3u";
            fileID = file;
            break;
        case "H":
            fileType = "hdhr";
            fileID = file;
            break;
        case "X":
            fileType = "xmltv";
            fileID = file.substring(0, file.lastIndexOf('.'));
            break;
    }
    if (SERVER["settings"]["files"][fileType].hasOwnProperty(fileID) == true) {
        var data = SERVER["settings"]["files"][fileType][fileID];
        return data[key];
    }
    return;
}
function setXmltvChannel(id, element) {
    var xmltv = new XMLTVFile();
    var xmlFile = element.value;
    var tvgId = SERVER["xepg"]["epgMapping"][id]["tvg-id"];
    var td = document.getElementById("popup-mapping").parentElement;
    td.innerHTML = "";
    var select = xmltv.getPrograms(element.value, tvgId);
    select.setAttribute("name", "x-mapping");
    select.setAttribute("id", "popup-mapping");
    select.setAttribute("onchange", "javascript: this.className = 'changed'; checkXmltvChannel('" + id + "',this,'" + xmlFile + "');");
    select.className = "changed";
    sortSelect(select);
    td.appendChild(select);
    checkXmltvChannel(id, select, xmlFile);
}
function checkXmltvChannel(id, element, xmlFile) {
    var value = element.value;
    var bool;
    var checkbox = document.getElementById('active');
    var channel = SERVER["xepg"]["epgMapping"][id];
    var updateLogo;
    if (value == "-") {
        bool = false;
    }
    else {
        bool = true;
    }
    checkbox.checked = bool;
    checkbox.className = "changed";
    console.log(xmlFile);
    // Kanallogo aktualisieren
    /*
    updateLogo = (document.getElementById("update-icon") as HTMLInputElement).checked
    console.log(updateLogo);
    */
    if (xmlFile != "xTeVe Dummy" && bool == true) {
        //(document.getElementById("update-icon") as HTMLInputElement).checked = true;
        //(document.getElementById("update-icon") as HTMLInputElement).className = "changed";
        console.log("ID", id);
        changeChannelLogo(id);
        return;
    }
    if (xmlFile == "xTeVe Dummy") {
        document.getElementById("update-icon").checked = false;
        document.getElementById("update-icon").className = "changed";
    }
    return;
}
function changeChannelLogo(id) {
    var updateLogo;
    var channel = SERVER["xepg"]["epgMapping"][id];
    var f = document.getElementById("popup-xmltv");
    var xmltvFile = f.options[f.selectedIndex].value;
    var m = document.getElementById("popup-mapping");
    var xMapping = m.options[m.selectedIndex].value;
    var xmltvLogo = SERVER["xepg"]["xmltvMap"][xmltvFile][xMapping]["icon"];
    updateLogo = document.getElementById("update-icon").checked;
    if (updateLogo == true && xmltvFile != "xTeVe Dummy") {
        if (SERVER["xepg"]["xmltvMap"][xmltvFile].hasOwnProperty(xMapping)) {
            var logo = xmltvLogo;
        }
        else {
            logo = channel["tvg-logo"];
        }
        var logoInput = document.getElementById("channel-icon");
        logoInput.value = logo;
        if (BULK_EDIT == false) {
            logoInput.className = "changed";
        }
    }
}
function savePopupData(dataType, id, remove, option) {
    if (dataType == "mapping") {
        var data = new Object();
        console.log("Save mapping data");
        cmd = "saveEpgMapping";
        data["epgMapping"] = SERVER["xepg"]["epgMapping"];
        console.log("SEND TO SERVER");
        var server = new Server(cmd);
        server.request(data);
        delete UNDO["epgMapping"];
        return;
    }
    console.log("Save popup data");
    var div = document.getElementById("popup-custom");
    var inputs = div.getElementsByTagName("TABLE")[0].getElementsByTagName("INPUT");
    var selects = div.getElementsByTagName("TABLE")[0].getElementsByTagName("SELECT");
    var input = new Object();
    var confirmMsg;
    for (var i = 0; i < selects.length; i++) {
        var name;
        name = selects[i].name;
        var value = selects[i].value;
        switch (name) {
            case "tuner":
                input[name] = parseInt(value);
                break;
            default:
                input[name] = value;
                break;
        }
    }
    for (var i = 0; i < inputs.length; i++) {
        switch (inputs[i].type) {
            case "checkbox":
                name = inputs[i].name;
                input[name] = inputs[i].checked;
                break;
            case "text":
            case "hidden":
            case "password":
                name = inputs[i].name;
                switch (name) {
                    case "tuner":
                        input[name] = parseInt(inputs[i].value);
                        break;
                    default:
                        input[name] = inputs[i].value;
                        break;
                }
                break;
        }
    }
    var data = new Object();
    var cmd;
    if (remove == true) {
        input["delete"] = true;
    }
    switch (dataType) {
        case "users":
            confirmMsg = "Delete this user?";
            if (id == "-") {
                cmd = "saveNewUser";
                data["userData"] = input;
            }
            else {
                cmd = "saveUserData";
                var d = new Object();
                d[id] = input;
                data["userData"] = d;
            }
            break;
        case "m3u":
            confirmMsg = "Delete this playlist?";
            switch (option) {
                // Popup: Save
                case 0:
                    cmd = "saveFilesM3U";
                    break;
                // Popup: Update
                case 1:
                    cmd = "updateFileM3U";
                    break;
            }
            data["files"] = new Object;
            data["files"][dataType] = new Object;
            data["files"][dataType][id] = input;
            break;
        case "hdhr":
            confirmMsg = "Delete this HDHomeRun tuner?";
            switch (option) {
                // Popup: Save
                case 0:
                    cmd = "saveFilesHDHR";
                    break;
                // Popup: Update
                case 1:
                    cmd = "updateFileHDHR";
                    break;
            }
            data["files"] = new Object;
            data["files"][dataType] = new Object;
            data["files"][dataType][id] = input;
            break;
        case "xmltv":
            confirmMsg = "Delete this XMLTV file?";
            switch (option) {
                // Popup: Save
                case 0:
                    cmd = "saveFilesXMLTV";
                    break;
                // Popup: Update
                case 1:
                    cmd = "updateFileXMLTV";
                    break;
            }
            data["files"] = new Object;
            data["files"][dataType] = new Object;
            data["files"][dataType][id] = input;
            break;
        case "filter":
            confirmMsg = "Delete this filter?";
            cmd = "saveFilter";
            data["filter"] = new Object;
            data["filter"][id] = input;
            break;
        default:
            console.log(dataType, id);
            return;
            break;
    }
    if (remove == true) {
        if (!confirm(confirmMsg)) {
            showElement("popup", false);
            return;
        }
    }
    console.log("SEND TO SERVER");
    console.log(data);
    var server = new Server(cmd);
    server.request(data);
}
function donePopupData(dataType, idsStr) {
    var ids = idsStr.split(',');
    var div = document.getElementById("popup-custom");
    var inputs = div.getElementsByClassName("changed");
    ids.forEach(function (id) {
        var input = new Object();
        input = SERVER["xepg"]["epgMapping"][id];
        console.log(input);
        for (var i = 0; i < inputs.length; i++) {
            var name;
            var value;
            switch (inputs[i].tagName) {
                case "INPUT":
                    switch (inputs[i].type) {
                        case "checkbox":
                            name = inputs[i].name;
                            value = inputs[i].checked;
                            input[name] = value;
                            break;
                        case "text":
                            name = inputs[i].name;
                            value = inputs[i].value;
                            input[name] = value;
                            break;
                    }
                    break;
                case "SELECT":
                    name = inputs[i].name;
                    value = inputs[i].value;
                    input[name] = value;
                    break;
            }
            switch (name) {
                case "tvg-logo":
                    //(document.getElementById(id).childNodes[2].firstChild as HTMLElement).setAttribute("src", value)
                    break;
                case "x-name":
                    document.getElementById(id).childNodes[3].firstChild.innerHTML = value;
                    break;
                case "x-category":
                    document.getElementById(id).childNodes[3].firstChild.className = value;
                    break;
                case "x-group-title":
                    document.getElementById(id).childNodes[5].firstChild.innerHTML = value;
                    break;
                case "x-xmltv-file":
                    if (value != "xTeVe Dummy" && value != "-") {
                        value = getValueFromProviderFile(value, "xmltv", "name");
                    }
                    if (value == "-") {
                        input["x-active"] = false;
                    }
                    document.getElementById(id).childNodes[6].firstChild.innerHTML = value;
                    break;
                case "x-mapping":
                    if (value == "-") {
                        input["x-active"] = false;
                    }
                    document.getElementById(id).childNodes[7].firstChild.innerHTML = value;
                    break;
                default:
            }
            createSearchObj();
            searchInMapping();
        }
        if (input["x-active"] == false) {
            document.getElementById(id).className = "notActiveEPG";
        }
        else {
            document.getElementById(id).className = "activeEPG";
        }
        console.log(input["tvg-logo"]);
        document.getElementById(id).childNodes[2].firstChild.setAttribute("src", input["tvg-logo"]);
    });
    showElement("popup", false);
    return;
}
function showPreview(element) {
    var div = document.getElementById("myStreamsBox");
    switch (element) {
        case false:
            div.className = "notVisible";
            return;
            break;
    }
    var streams = ["activeStreams", "inactiveStreams"];
    streams.forEach(function (preview) {
        var table = document.getElementById(preview);
        table.innerHTML = "";
        var obj = SERVER["data"]["StreamPreviewUI"][preview];
        obj.forEach(function (channel) {
            var tr = document.createElement("TR");
            var tdKey = document.createElement("TD");
            var tdVal = document.createElement("TD");
            tdKey.className = "tdKey";
            tdVal.className = "tdVal";
            switch (preview) {
                case "activeStreams":
                    tdKey.innerText = "Channel: (+)";
                    break;
                case "inactiveStreams":
                    tdKey.innerText = "Channel: (-)";
                    break;
            }
            tdVal.innerText = channel;
            tr.appendChild(tdKey);
            tr.appendChild(tdVal);
            table.appendChild(tr);
        });
    });
    showElement("loading", false);
    div.className = "visible";
    return;
}
