class MainMenu {
    constructor() {
        this.DocumentID = "main-menu";
        this.HTMLTag = "LI";
        this.ImagePath = "img/";
    }
    createIMG(src) {
        let element = document.createElement("IMG");
        element.setAttribute("src", this.ImagePath + src);
        return element;
    }
    createValue(value) {
        let element = document.createElement("P");
        element.innerHTML = value;
        return element;
    }
}
class MainMenuItem extends MainMenu {
    constructor(menuKey, value, image, headline) {
        super();
        this.menuKey = menuKey;
        this.value = value;
        this.imgSrc = image;
        this.headline = headline;
    }
    createItem() {
        let item = document.createElement("LI");
        item.setAttribute("onclick", "javascript: openThisMenu(this)");
        item.setAttribute("id", this.id);
        let img = this.createIMG(this.imgSrc);
        let value = this.createValue(this.value);
        item.appendChild(img);
        item.appendChild(value);
        let doc = document.getElementById(this.DocumentID);
        doc.appendChild(item);
        switch (this.menuKey) {
            case "playlist":
                this.tableHeader = ["{{.playlist.table.playlist}}", "{{.playlist.table.tuner}}", "{{.playlist.table.lastUpdate}}", "{{.playlist.table.availability}} %", "{{.playlist.table.type}}", "{{.playlist.table.streams}}", "{{.playlist.table.groupTitle}} %", "{{.playlist.table.tvgID}} %", "{{.playlist.table.uniqueID}} %"];
                break;
            case "xmltv":
                this.tableHeader = ["{{.xmltv.table.guide}}", "{{.xmltv.table.lastUpdate}}", "{{.xmltv.table.availability}} %", "{{.xmltv.table.channels}}", "{{.xmltv.table.programs}}"];
                break;
            case "filter":
                this.tableHeader = ["{{.filter.table.startingChannel}}", "{{.filter.table.name}}", "{{.filter.table.type}}", "{{.filter.table.filter}}"];
                break;
            case "users":
                this.tableHeader = ["{{.users.table.username}}", "{{.users.table.password}}", "{{.users.table.web}}", "{{.users.table.pms}}", "{{.users.table.m3u}}", "{{.users.table.xml}}", "{{.users.table.api}}"];
                break;
            case "mapping":
                this.tableHeader = ["BULK", "{{.mapping.table.chNo}}", "{{.mapping.table.logo}}", "{{.mapping.table.channelName}}", "{{.mapping.table.updateChannelNameRegex}}", "{{.mapping.table.playlist}}", "{{.mapping.table.groupTitle}}", "{{.mapping.table.xmltvFile}}", "{{.mapping.table.xmltvID}}", "{{.mapping.table.timeshift}}"];
                break;
        }
    }
}
class Content {
    constructor() {
        this.DocumentID = "content";
        this.TableID = "content_table";
        this.headerClass = "content_table_header";
        this.interactionID = "content-interaction";
    }
    createHeadline(value) {
        let element = document.createElement("H3");
        element.innerHTML = value;
        return element;
    }
    createHR() {
        return document.createElement("HR");
    }
    createInteraction() {
        let element = document.createElement("DIV");
        element.setAttribute("id", this.interactionID);
        return element;
    }
    createDIV() {
        let element = document.createElement("DIV");
        element.id = this.DivID;
        return element;
    }
    createTABLE() {
        let element = document.createElement("TABLE");
        element.id = this.TableID;
        return element;
    }
    createTableRow() {
        let element = document.createElement("TR");
        element.className = this.headerClass;
        return element;
    }
    createTableContent(menuKey) {
        let data = {};
        let rows = [];
        let fileTypes = [];
        switch (menuKey) {
            case "playlist":
                fileTypes = ["m3u", "hdhr"];
                fileTypes.forEach(fileType => {
                    data = SERVER["settings"]["files"][fileType];
                    var keys = getOwnObjProps(data);
                    keys.forEach(key => {
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
                var keys = getOwnObjProps(data);
                keys.forEach(key => {
                    var tr = document.createElement("TR");
                    tr.id = key;
                    tr.setAttribute('onclick', 'javascript: openPopUp("' + data[key]["type"] + '", this)');
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    cell.value = data[key]["startingChannel"];
                    tr.appendChild(cell.createCell());
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
                fileTypes = new Array("xmltv");
                fileTypes.forEach(fileType => {
                    data = SERVER["settings"]["files"][fileType];
                    var keys = getOwnObjProps(data);
                    keys.forEach(key => {
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
                fileTypes = new Array("users");
                fileTypes.forEach(fileType => {
                    data = SERVER[fileType];
                    var keys = getOwnObjProps(data);
                    keys.forEach(key => {
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
                data = SERVER["xepg"]["epgMapping"];
                var keys = getOwnObjProps(data);
                keys.forEach(key => {
                    var tr = document.createElement("TR");
                    tr.id = key;
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
                    // Channel number
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "INPUTCHANNEL";
                    cell.value = data[key]["x-channelID"];
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
                    // Channel name
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    cell.className = data[key]["x-category"];
                    cell.value = data[key]["x-name"];
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    // Update channel name regex
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    cell.value = data[key]["update-channel-name-regex"];
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    // Playlist
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    cell.value = getValueFromProviderFile(data[key]["_file.m3u.id"], "m3u", "name");
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    // Group (group-title)
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    cell.value = data[key]["x-group-title"];
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    // XMLTV file
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
                    // XMLTV Channel
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    var value = data[key]["x-mapping"];
                    if (value.length > 20) {
                        value = data[key]["x-mapping"].substring(0, 20) + "...";
                    }
                    cell.value = value;
                    var td = cell.createCell();
                    td.setAttribute('onclick', 'javascript: openPopUp("mapping", this)');
                    td.id = key;
                    tr.appendChild(td);
                    // TimeShift
                    var cell = new Cell();
                    cell.child = true;
                    cell.childType = "P";
                    cell.value = data[key]["x-timeshift"];
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
                break;
        }
        return rows;
    }
}
class Cell {
    createCell() {
        let td = document.createElement("TD");
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
            td.setAttribute("onclick", this.onclickFunction);
            td.className = "pointer";
        }
        if (this.tdClassName != undefined) {
            td.className = this.tdClassName;
        }
        return td;
    }
}
class ShowContent extends Content {
    constructor(menuID) {
        super();
        this.menuID = menuID;
    }
    createInput(type, name, value) {
        let input = document.createElement("INPUT");
        input.setAttribute("type", type);
        input.setAttribute("name", name);
        input.setAttribute("value", value);
        return input;
    }
    show() {
        COLUMN_TO_SORT = -1;
        // Delete old content
        let doc = document.getElementById(this.DocumentID);
        doc.innerHTML = "";
        showPreview(false);
        // Headline
        let headline = menuItems[this.menuID].headline;
        let menuKey = menuItems[this.menuID].menuKey;
        let h = this.createHeadline(headline);
        doc.appendChild(h);
        let hr = this.createHR();
        doc.appendChild(hr);
        // Interaction
        let div = this.createInteraction();
        doc.appendChild(div);
        let interaction = document.getElementById(this.interactionID);
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
            case "logout":
                location.reload();
                document.cookie = "Token= ; expires = Thu, 01 Jan 1970 00:00:00 GMT";
                break;
            default:
                break;
        }
        // Create table (if needed)
        var tableHeader = menuItems[this.menuID].tableHeader;
        if (tableHeader.length > 0) {
            var wrapper = document.createElement("DIV");
            doc.appendChild(wrapper);
            wrapper.setAttribute("id", "box-wrapper");
            var table = this.createTABLE();
            wrapper.appendChild(table);
            var header = this.createTableRow();
            table.appendChild(header);
            // Table header
            tableHeader.forEach(element => {
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
                        cell.onclickFunction = "javascript: sortTable(1);";
                        cell.tdClassName = "sortThis";
                    }
                    if (element == "{{.mapping.table.channelName}}") {
                        cell.onclick = true;
                        cell.onclickFunction = "javascript: sortTable(3);";
                    }
                    if (element == "{{.mapping.table.updateChannelNameRegex}}") {
                        cell.onclick = true;
                        cell.onclickFunction = "javascript: sortTable(4);";
                    }
                    if (element == "{{.mapping.table.playlist}}") {
                        cell.onclick = true;
                        cell.onclickFunction = "javascript: sortTable(5);";
                    }
                    if (element == "{{.mapping.table.groupTitle}}") {
                        cell.onclick = true;
                        cell.onclickFunction = "javascript: sortTable(6);";
                    }
                    if (element == "{{.mapping.table.timeshift}}") {
                        cell.onclick = true;
                        cell.onclickFunction = "javascript: sortTable(9);";
                    }
                }
                if (menuKey == "filter") {
                    if (element == "{{.filter.table.startingChannel}}") {
                        cell.onclick = true;
                        cell.onclickFunction = "javascript: sortTable(0);";
                        cell.tdClassName = "sortThis";
                    }
                    if (element == "{{.filter.table.name}}") {
                        cell.onclick = true;
                        cell.onclickFunction = "javascript: sortTable(1);";
                    }
                    if (element == "{{.filter.table.type}}") {
                        cell.onclick = true;
                        cell.onclickFunction = "javascript: sortTable(2);";
                    }
                    if (element == "{{.filter.table.filter}}") {
                        cell.onclick = true;
                        cell.onclickFunction = "javascript: sortTable(3);";
                    }
                }
                header.appendChild(cell.createCell());
            });
            table.appendChild(header);
            // Content of the table
            var rows = this.createTableContent(menuKey);
            rows.forEach(tr => {
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
    }
}
function PageReady() {
    let server = new Server("getServerConfig");
    server.request({});
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
    let obj = SERVER["clientInfo"];
    let keys = getOwnObjProps(obj);
    for (var i = 0; i < keys.length; i++) {
        if (document.getElementById(keys[i])) {
            document.getElementById(keys[i]).innerHTML = obj[keys[i]];
            if (location.protocol === 'https:') {
                if (keys[i] === "xepg-url" || keys[i] === "m3u-url" || keys[i] === "DVR") {
                    document.getElementById(keys[i]).addEventListener('click', function (event) {
                        const target = event.target;
                        navigator.clipboard.writeText(target.innerText.split(" ")[0]).then(() => { });
                    }, false);
                }
            }
        }
    }
    if (!document.getElementById("main-menu")) {
        return;
    }
    // Create menu
    document.getElementById("main-menu").innerHTML = "";
    for (let i = 0; i < menuItems.length; i++) {
        menuItems[i].id = i;
        switch (menuItems[i]["menuKey"]) {
            case "users":
            case "logout":
                if (SERVER["settings"]["authentication.web"] == true) {
                    menuItems[i].createItem();
                }
                break;
            case "mapping":
            case "xmltv":
                if (SERVER["clientInfo"]["epgSource"] == "XEPG") {
                    menuItems[i].createItem();
                }
                break;
            default:
                menuItems[i].createItem();
                break;
        }
    }
    return;
}
function openThisMenu(element) {
    let id = element.id;
    let content = new ShowContent(id);
    content.show();
    calculateWrapperHeight();
    return;
}
class PopupWindow {
    constructor() {
        this.DocumentID = "popup-custom";
        this.InteractionID = "interaction";
        this.doc = document.getElementById(this.DocumentID);
    }
    createTitle(title) {
        let td = document.createElement("TD");
        td.className = "left";
        td.innerHTML = title + ":";
        return td;
    }
    createContent(element) {
        let td = document.createElement("TD");
        td.appendChild(element);
        return td;
    }
    createInteraction() {
        let div = document.createElement("div");
        div.setAttribute("id", "popup-interaction");
        div.className = "interaction";
        this.doc.appendChild(div);
    }
}
class PopupContent extends PopupWindow {
    constructor() {
        super(...arguments);
        this.table = document.createElement("TABLE");
    }
    createHeadline(headline) {
        this.doc.innerHTML = "";
        var element = document.createElement("H3");
        element.innerHTML = headline.toUpperCase();
        this.doc.appendChild(element);
        // Create table
        this.table = document.createElement("TABLE");
        this.doc.appendChild(this.table);
    }
    appendRow(title, element) {
        let tr = document.createElement("TR");
        // Title
        if (title.length != 0) {
            tr.appendChild(this.createTitle(title));
        }
        // Content
        tr.appendChild(this.createContent(element));
        this.table.appendChild(tr);
    }
    createInput(type, name, value) {
        let input = document.createElement("INPUT");
        if (value == undefined) {
            value = "";
        }
        input.setAttribute("type", type);
        input.setAttribute("name", name);
        input.setAttribute("value", value);
        return input;
    }
    createCheckbox(name) {
        let input = document.createElement("INPUT");
        input.setAttribute("type", "checkbox");
        input.setAttribute("name", name);
        return input;
    }
    // Creates a selection of multiple options values with text descriptions
    createSelect(text, values, set, dbKey) {
        let select = document.createElement("SELECT");
        select.setAttribute("name", dbKey);
        for (let i = 0; i < text.length; i++) {
            let option = document.createElement("OPTION");
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
    }
    selectOption(select, value) {
        //select.selectedOptions = value
        let s = select;
        s.options[s.selectedIndex].value = value;
        return select;
    }
    description(value) {
        let tr = document.createElement("TR");
        let td = document.createElement("TD");
        let span = document.createElement("PRE");
        span.innerHTML = value;
        tr.appendChild(td);
        tr.appendChild(this.createContent(span));
        this.table.appendChild(tr);
    }
    // Interaction
    addInteraction(element) {
        let interaction = document.getElementById("popup-interaction");
        interaction.appendChild(element);
    }
}
function openPopUp(dataType, element) {
    let data = {};
    let id;
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
    let content = new PopupContent();
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
            // Interaction
            content.createInteraction();
            // Abort
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Next
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
            // Description
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
                var text = [];
                var values = [];
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
            // Interation
            content.createInteraction();
            // Delete
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
            // Abort
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Update
            if (data["id.provider"] != "-") {
                var input = content.createInput("button", "update", "{{.button.update}}");
                input.setAttribute('onclick', 'javascript: savePopupData("m3u", "' + id + '", false, 1)');
                content.addInteraction(input);
            }
            // Save
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
            // Description
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
                var text = [];
                var values = [];
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
            // Interaction
            content.createInteraction();
            // Delete
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
            // Abort
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Update
            if (data["id.provider"] != "-") {
                var input = content.createInput("button", "update", "{{.button.update}}");
                input.setAttribute('onclick', 'javascript: savePopupData("hdhr", "' + id + '", false, 1)');
                content.addInteraction(input);
            }
            // Save
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
            // Interaction
            content.createInteraction();
            // Abort
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Next
            var input = content.createInput("button", "next", "{{.button.next}}");
            input.setAttribute("onclick", 'javascript: openPopUp("group-title")');
            input.setAttribute("id", 'next');
            content.addInteraction(input);
            break;
        case "custom-filter":
        case "group-title":
            switch (dataType) {
                case "custom-filter":
                    content.createHeadline("{{.filter.custom}}" + " Filter");
                    break;
                case "group-title":
                    content.createHeadline("{{.filter.group}}" + " Filter");
                    break;
            }
            // Name      
            var dbKey = "name";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.filter.name.placeholder}}");
            content.appendRow("{{.filter.name.title}}", input);
            // Description
            var dbKey = "description";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.filter.description.placeholder}}");
            content.appendRow("{{.filter.description.title}}", input);
            // Type
            var dbKey = "type";
            var input = content.createInput("hidden", dbKey, data[dbKey]);
            content.appendRow("", input);
            var filterType = data[dbKey];
            switch (filterType) {
                case "custom-filter":
                    // Case sensitive
                    var dbKey = "caseSensitive";
                    var input = content.createCheckbox(dbKey);
                    input.checked = data[dbKey];
                    content.appendRow("{{.filter.caseSensitive.title}}", input);
                    // Filter Rule (Custom)
                    var dbKey = "filter";
                    var input = content.createInput("text", dbKey, data[dbKey]);
                    input.setAttribute("placeholder", "{{.filter.filterRule.placeholder}}");
                    content.appendRow("{{.filter.filterRule.title}}", input);
                    // Starting Channel Number Mapping
                    var dbKey = "startingChannel";
                    var input = content.createInput("text", dbKey, data[dbKey]);
                    input.setAttribute("placeholder", "{{.filter.startingChannel.placeholder}}");
                    content.appendRow("{{.filter.startingChannel.title}}", input);
                    break;
                case "group-title":
                    //alert(dbKey + " " + filterType)
                    // Filter based on the groups in the M3U
                    var dbKey = "filter";
                    var groupsM3U = getLocalData("m3uGroups", "");
                    var text = groupsM3U["text"];
                    var values = groupsM3U["value"];
                    var select = content.createSelect(text, values, data[dbKey], dbKey);
                    select.setAttribute("onchange", "javascript: this.className = 'changed'");
                    content.appendRow("{{.filter.filterGroup.title}}", select);
                    content.description("{{.filter.filterGroup.description}}");
                    // Case sensetive
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
                    // Preserve M3U Playlist Channel Mapping
                    var dbKey = "preserveMapping";
                    var input = content.createCheckbox(dbKey);
                    input.checked = data[dbKey];
                    content.appendRow("{{.filter.preserveMapping.title}}", input);
                    // Starting Channel Number Mapping
                    var dbKey = "startingChannel";
                    var input = content.createInput("text", dbKey, data[dbKey]);
                    input.setAttribute("placeholder", "{{.filter.startingChannel.placeholder}}");
                    content.appendRow("{{.filter.startingChannel.title}}", input);
                    break;
                default:
                    break;
            }
            // Interaction
            content.createInteraction();
            // Delete
            var input = content.createInput("button", "delete", "{{.button.delete}}");
            input.setAttribute('onclick', 'javascript: savePopupData("filter", "' + id + '", true, 0)');
            input.className = "delete";
            content.addInteraction(input);
            // Abort
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Save
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
            // Description
            var dbKey = "description";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.xmltv.description.placeholder}}");
            content.appendRow("{{.xmltv.description.title}}", input);
            // URL
            var dbKey = "file.source";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.xmltv.fileXMLTV.placeholder}}");
            content.appendRow("{{.xmltv.fileXMLTV.title}}", input);
            // Interaction
            content.createInteraction();
            // Delete
            if (data["id.provider"] != "-") {
                var input = content.createInput("button", "delete", "{{.button.delete}}");
                input.setAttribute('onclick', 'javascript: savePopupData("xmltv", "' + id + '", true, 0)');
                input.className = "delete";
                content.addInteraction(input);
            }
            // Abort
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Update
            if (data["id.provider"] != "-") {
                var input = content.createInput("button", "update", "{{.button.update}}");
                input.setAttribute('onclick', 'javascript: savePopupData("xmltv", "' + id + '", false, 1)');
                content.addInteraction(input);
            }
            // Save
            var input = content.createInput("button", "save", "{{.button.save}}");
            input.setAttribute('onclick', 'javascript: savePopupData("xmltv", "' + id + '", false, 0)');
            content.addInteraction(input);
            break;
        case "users":
            content.createHeadline("{{.mainMenu.item.users}}");
            // User name
            var dbKey = "username";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.users.username.placeholder}}");
            content.appendRow("{{.users.username.title}}", input);
            // New Parssword
            var dbKey = "password";
            var input = content.createInput("password", dbKey, "");
            input.setAttribute("placeholder", "{{.users.password.placeholder}}");
            content.appendRow("{{.users.password.title}}", input);
            // Confirmation
            var dbKey = "confirm";
            var input = content.createInput("password", dbKey, "");
            input.setAttribute("placeholder", "{{.users.confirm.placeholder}}");
            content.appendRow("{{.users.confirm.title}}", input);
            // Authentication WEB
            var dbKey = "authentication.web";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            if (data["defaultUser"] == true) {
                input.setAttribute("onclick", "javascript: return false");
            }
            content.appendRow("{{.users.web.title}}", input);
            // Authentication PMS
            var dbKey = "authentication.pms";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            content.appendRow("{{.users.pms.title}}", input);
            // Authentication M3U
            var dbKey = "authentication.m3u";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            content.appendRow("{{.users.m3u.title}}", input);
            // Authentication XML
            var dbKey = "authentication.xml";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            content.appendRow("{{.users.xml.title}}", input);
            // Authentication API
            var dbKey = "authentication.api";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            content.appendRow("{{.users.api.title}}", input);
            // Interaction
            content.createInteraction();
            // Delete
            if (data["defaultUser"] != true && id != "-") {
                var input = content.createInput("button", "delete", "{{.button.delete}}");
                input.className = "delete";
                input.setAttribute('onclick', 'javascript: savePopupData("' + dataType + '", "' + id + '", true, 0)');
                content.addInteraction(input);
            }
            // Abort
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Save
            var input = content.createInput("button", "save", "{{.button.save}}");
            input.setAttribute("onclick", 'javascript: savePopupData("' + dataType + '", "' + id + '", "false");');
            content.addInteraction(input);
            break;
        case "mapping":
            content.createHeadline("{{.mainMenu.item.mapping}}");
            // Active
            var dbKey = "x-active";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            input.id = "active";
            //input.setAttribute("onchange", "javascript: this.className = 'changed'")
            input.setAttribute("onchange", "javascript: toggleChannelStatus('" + id + "', this)");
            content.appendRow("{{.mapping.active.title}}", input);
            // Channel name
            var dbKey = "x-name";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            if (BULK_EDIT == true) {
                input.style.border = "solid 1px red";
                input.setAttribute("readonly", "true");
            }
            content.appendRow("{{.mapping.channelName.title}}", input);
            content.description(data["name"]);
            // Description
            var dbKey = "x-description";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.mapping.description.placeholder}}");
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            content.appendRow("{{.mapping.description.title}}", input);
            // Update the channel x-name
            if (data.hasOwnProperty("_uuid.key")) {
                if (data["_uuid.key"] != "") {
                    var dbKey = "x-update-channel-name";
                    var input = content.createCheckbox(dbKey);
                    input.setAttribute("onchange", "javascript: this.className = 'changed'");
                    input.checked = data[dbKey];
                    content.appendRow("{{.mapping.updateChannelName.title}}", input);
                }
            }
            // Channel name regex for updating the channel name
            var dbKey = "update-channel-name-regex";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.mapping.updateChannelNameRegex.placeholder}}");
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            content.appendRow("{{.mapping.updateChannelNameRegex.title}}", input);
            content.description("{{.mapping.updateChannelNameRegex.description}}");
            // Channel group regex for updating the channel name
            var dbKey = "update-channel-name-by-group-regex";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("placeholder", "{{.mapping.updateChannelNameByGroupRegex.placeholder}}");
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            content.appendRow("{{.mapping.updateChannelNameByGroupRegex.title}}", input);
            content.description("{{.mapping.updateChannelNameByGroupRegex.description}}");
            // Logo URL (Channel) 
            var dbKey = "tvg-logo";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            input.setAttribute("id", "channel-icon");
            content.appendRow("{{.mapping.channelLogo.title}}", input);
            // Channel logo update
            var dbKey = "x-update-channel-icon";
            var input = content.createCheckbox(dbKey);
            input.checked = data[dbKey];
            input.setAttribute("id", "update-icon");
            input.setAttribute("onchange", "javascript: this.className = 'changed'; changeChannelLogo('" + id + "');");
            content.appendRow("{{.mapping.updateChannelLogo.title}}", input);
            // Expand EPG category
            var dbKey = "x-category";
            var text = ["-", "Kids (Emby only)", "News", "Movie", "Series", "Sports"];
            var values = ["", "Kids", "News", "Movie", "Series", "Sports"];
            var select = content.createSelect(text, values, data[dbKey], dbKey);
            select.setAttribute("onchange", "javascript: this.className = 'changed'");
            content.appendRow("{{.mapping.epgCategory.title}}", select);
            // M3U group title
            var dbKey = "x-group-title";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            input.dataset.oldValue = data[dbKey];
            content.appendRow("{{.mapping.m3uGroupTitle.title}}", input);
            if (data["group-title"] != undefined) {
                content.description(data["group-title"]);
            }
            if (data["x-update-channel-group"] == true) {
                input.disabled = true;
            }
            // Update channel group checkbox
            var dbKey = "x-update-channel-group";
            var input = content.createCheckbox(dbKey);
            input.setAttribute("onchange", "javascript: toggleGroupUpdateCb('" + id + "', this);");
            input.checked = data[dbKey];
            content.appendRow("{{.mapping.updateChannelGroup.title}}", input);
            content.description("{{.mapping.updateChannelGroup.description}}");
            // XMLTV file
            var dbKey = 'x-xmltv-file';
            const xmlTvFile = data[dbKey];
            var xmlTv = new XMLTVFile();
            const xmlTvFileSelect = xmlTv.getFiles(data[dbKey]);
            xmlTvFileSelect.setAttribute('name', dbKey);
            xmlTvFileSelect.setAttribute('id', 'popup-xmltv');
            xmlTvFileSelect.setAttribute('onchange', `javascript: this.className = 'changed'; setXmltvChannel('${id}', this);`);
            content.appendRow('{{.mapping.xmltvFile.title}}', xmlTvFileSelect);
            // XMLTV Mapping
            var dbKey = 'x-mapping';
            var xmlTv = new XMLTVFile();
            const currentXmlTvId = data[dbKey];
            const [xmlTvIdContainer, xmlTvIdInput, xmlTvIdDatalist] = xmlTv.newXmlTvIdPicker(xmlTvFile, currentXmlTvId);
            xmlTvIdContainer.setAttribute('id', 'xmltv-id-picker-container');
            xmlTvIdInput.setAttribute('list', 'xmltv-id-picker-datalist');
            xmlTvIdInput.setAttribute('name', 'x-mapping'); // Should stay x-mapping as it will be used in donePopupData to make a server request
            xmlTvIdInput.setAttribute('id', 'xmltv-id-picker-input');
            xmlTvIdInput.setAttribute('onchange', `javascript: this.className = 'changed'; checkXmltvChannel('${id}', this.value, '${xmlTvFile}');`);
            xmlTvIdDatalist.setAttribute('id', 'xmltv-id-picker-datalist');
            content.appendRow('{{.mapping.xmltvChannel.title}}', xmlTvIdContainer);
            // Timeshift
            var dbKey = "x-timeshift";
            var input = content.createInput("text", dbKey, data[dbKey]);
            input.setAttribute("onchange", "javascript: this.className = 'changed'");
            input.setAttribute("placeholder", "{{.mapping.timeshift.placeholder}}");
            input.setAttribute("id", "timeshift");
            content.appendRow("{{.mapping.timeshift.title}}", input);
            // Interaction
            content.createInteraction();
            // Upload logo
            var input = content.createInput("button", "cancel", "{{.button.uploadLogo}}");
            input.setAttribute("onclick", 'javascript: uploadLogo();');
            content.addInteraction(input);
            // Abort
            var input = content.createInput("button", "cancel", "{{.button.cancel}}");
            input.setAttribute("onclick", 'javascript: showElement("popup", false);');
            content.addInteraction(input);
            // Finished
            var ids = [];
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
class XMLTVFile {
    getFiles(set) {
        let fileIDs = getOwnObjProps(SERVER["xepg"]["xmltvMap"]);
        let values = new Array("-");
        let text = new Array("-");
        for (let i = 0; i < fileIDs.length; i++) {
            if (fileIDs[i] != "xTeVe Dummy") {
                values.push(getValueFromProviderFile(fileIDs[i], "xmltv", "file.xteve"));
                text.push(getValueFromProviderFile(fileIDs[i], "xmltv", "name"));
            }
            else {
                values.push(fileIDs[i]);
                text.push(fileIDs[i]);
            }
        }
        let select = document.createElement("SELECT");
        for (let i = 0; i < text.length; i++) {
            var option = document.createElement("OPTION");
            option.setAttribute("value", values[i]);
            option.innerText = text[i];
            select.appendChild(option);
        }
        if (set != "") {
            select.value = set;
        }
        return select;
    }
    /**
     * @param xmlTvFile XML file path to get EPG from.
     * @param currentXmlTvId Current XMLTV ID to set initial input value to.
     * @returns Array of, sequentially:
     * 1) Container of the picker.
     * 2) Input field to type at and get choice from.
     * 3) Datalist containing every option.
     */
    newXmlTvIdPicker(xmlTvFile, currentXmlTvId) {
        const container = document.createElement('div');
        const input = document.createElement('input');
        input.setAttribute('type', 'text');
        // Initially, set value to '-' if input is empty
        input.value = (currentXmlTvId) ? currentXmlTvId : '-';
        // When input is focused, remove '-' from it
        input.addEventListener('focus', (evt) => {
            const target = evt.target;
            target.value = (target.value === '-') ? '' : target.value;
        });
        // When input lose focus or take a value, if it's empty, set value to '-'
        input.addEventListener('blur', setFallbackValue);
        input.addEventListener('change', setFallbackValue);
        function setFallbackValue(evt) {
            const target = evt.target;
            target.value = (target.value) ? target.value : '-';
        }
        container.appendChild(input);
        const datalist = document.createElement('datalist');
        const option = document.createElement('option');
        option.setAttribute('value', '-');
        option.innerText = '-';
        datalist.appendChild(option);
        const epg = SERVER['xepg']['xmltvMap'][xmlTvFile];
        if (epg) {
            const programIds = getOwnObjProps(epg);
            programIds.forEach((programId) => {
                const program = epg[programId];
                if (program.hasOwnProperty('display-names')) {
                    program['display-names'].forEach((displayName) => {
                        const option = document.createElement('option');
                        option.setAttribute('value', programId);
                        option.innerText = displayName['Value'];
                        datalist.appendChild(option);
                    });
                }
                else {
                    const option = document.createElement('option');
                    option.setAttribute('value', programId);
                    option.innerText = '-';
                    datalist.appendChild(option);
                }
            });
        }
        container.appendChild(datalist);
        return [container, input, datalist];
    }
}
function getValueFromProviderFile(file, fileType, key) {
    if (file == "xTeVe Dummy") {
        return file;
    }
    let fileID;
    let indicator = file.charAt(0);
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
function setXmltvChannel(epgMapId, xmlTvFileSelect) {
    const xmlTv = new XMLTVFile();
    const newXmlTvFile = xmlTvFileSelect.value;
    // Remove old XMLTV ID selection box
    const xmlTvIdPickerParent = document.getElementById('xmltv-id-picker-container').parentElement;
    xmlTvIdPickerParent.innerHTML = '';
    // Create new XMLTV ID selection box
    const tvgId = SERVER['xepg']['epgMapping'][epgMapId]['tvg-id'];
    const [xmlTvIdContainer, xmlTvIdInput, xmlTvIdDatalist] = xmlTv.newXmlTvIdPicker(newXmlTvFile, tvgId);
    xmlTvIdContainer.setAttribute('id', 'xmltv-id-picker-container');
    xmlTvIdInput.setAttribute('list', 'xmltv-id-picker-datalist');
    xmlTvIdInput.setAttribute('name', 'x-mapping'); // Should stay x-mapping as it will be used in donePopupData to make a server request
    xmlTvIdInput.setAttribute('id', 'xmltv-id-picker-input');
    xmlTvIdInput.setAttribute('onchange', `javascript: this.className = 'changed'; checkXmltvChannel('${epgMapId}', this.value, '${newXmlTvFile}');`);
    xmlTvIdInput.classList.add('changed');
    xmlTvIdDatalist.setAttribute('id', 'xmltv-id-picker-datalist');
    // Add new XMLTV ID selection box to it's parent
    xmlTvIdPickerParent.appendChild(xmlTvIdContainer);
    checkXmltvChannel(epgMapId, xmlTvIdInput.value, newXmlTvFile);
}
function checkXmltvChannel(epgMapId, newXmlTvId, xmlTvFile) {
    const channelActiveCb = document.getElementById('active');
    const channelActive = newXmlTvId != '-';
    channelActiveCb.checked = channelActive;
    channelActiveCb.className = 'changed';
    if (xmlTvFile != 'xTeVe Dummy' && channelActive == true) {
        changeChannelLogo(epgMapId);
        return;
    }
    if (xmlTvFile == 'xTeVe Dummy') {
        document.getElementById('update-icon').checked = false;
        document.getElementById('update-icon').className = 'changed';
    }
}
function changeChannelLogo(epgMapId) {
    const channel = SERVER['xepg']['epgMapping'][epgMapId];
    const xmlTvFileSelect = document.getElementById('popup-xmltv');
    const xmlTvFile = xmlTvFileSelect.options[xmlTvFileSelect.selectedIndex].value;
    const xmlTvIdInput = document.getElementById('xmltv-id-picker-input');
    const newXmlTvId = xmlTvIdInput.value;
    const updateLogo = document.getElementById('update-icon').checked;
    let logo;
    if (updateLogo == true && xmlTvFile != 'xTeVe Dummy') {
        if (SERVER['xepg']['xmltvMap'][xmlTvFile].hasOwnProperty(newXmlTvId)) {
            logo = SERVER['xepg']['xmltvMap'][xmlTvFile][newXmlTvId]['icon'];
        }
        else {
            logo = channel['tvg-logo'];
        }
        var logoInput = document.getElementById('channel-icon');
        logoInput.value = logo;
        if (BULK_EDIT == false) {
            logoInput.className = 'changed';
        }
    }
}
function savePopupData(dataType, id, remove, option) {
    if (dataType == "mapping") {
        let data = {};
        let cmd = "saveEpgMapping";
        data["epgMapping"] = SERVER["xepg"]["epgMapping"];
        let server = new Server(cmd);
        server.request(data);
        delete UNDO["epgMapping"];
        return;
    }
    let div = document.getElementById("popup-custom");
    let inputs = div.getElementsByTagName("TABLE")[0].getElementsByTagName("INPUT");
    let selects = div.getElementsByTagName("TABLE")[0].getElementsByTagName("SELECT");
    let input = {};
    let confirmMsg;
    for (let i = 0; i < selects.length; i++) {
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
    for (let i = 0; i < inputs.length; i++) {
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
    let data = {};
    let cmd;
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
                let d = {};
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
            data["files"] = {};
            data["files"][dataType] = {};
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
            data["files"] = {};
            data["files"][dataType] = {};
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
            data["files"] = {};
            data["files"][dataType] = {};
            data["files"][dataType][id] = input;
            break;
        case "filter":
            confirmMsg = "Delete this filter?";
            cmd = "saveFilter";
            data["filter"] = {};
            data["filter"][id] = input;
            break;
        default:
            return;
    }
    if (remove == true) {
        if (!confirm(confirmMsg)) {
            showElement("popup", false);
            return;
        }
    }
    let server = new Server(cmd);
    server.request(data);
}
function donePopupData(dataType, idsStr) {
    let ids = idsStr.split(',');
    let div = document.getElementById("popup-custom");
    let inputs = div.getElementsByClassName("changed");
    ids.forEach(id => {
        let input;
        input = SERVER["xepg"]["epgMapping"][id];
        for (let i = 0; i < inputs.length; i++) {
            let name;
            let value;
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
                case "update-channel-name-regex":
                    document.getElementById(id).childNodes[4].firstChild.innerHTML = value;
                    break;
                case "x-group-title":
                    document.getElementById(id).childNodes[6].firstChild.innerHTML = value;
                    break;
                case "x-xmltv-file":
                    if (value != "xTeVe Dummy" && value != "-") {
                        value = getValueFromProviderFile(value, "xmltv", "name");
                    }
                    if (value == "-") {
                        input["x-active"] = false;
                    }
                    document.getElementById(id).childNodes[7].firstChild.innerHTML = value;
                    break;
                case "x-mapping":
                    if (value == "-") {
                        input["x-active"] = false;
                    }
                    document.getElementById(id).childNodes[8].firstChild.innerHTML = value;
                    break;
                case "x-timeshift":
                    document.getElementById(id).childNodes[9].firstChild.innerHTML = value;
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
        document.getElementById(id).childNodes[2].firstChild.setAttribute("src", input["tvg-logo"]);
    });
    showElement("popup", false);
    return;
}
function showPreview(element) {
    let div = document.getElementById("myStreamsBox");
    switch (element) {
        case false:
            div.className = "notVisible";
            return;
    }
    let streams = ["activeStreams", "inactiveStreams"];
    streams.forEach(preview => {
        let table = document.getElementById(preview);
        table.innerHTML = "";
        let obj = SERVER["data"]["StreamPreviewUI"][preview];
        obj.forEach(channel => {
            let tr = document.createElement("TR");
            let tdKey = document.createElement("TD");
            let tdVal = document.createElement("TD");
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
