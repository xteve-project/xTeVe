
function setMenuItem() {

  menu = new Object();
  subMenu = new Object();

  var menu_m3u = new Object();
  menu_m3u["_menuType"]       = "inputArray";
  menu_m3u["_element"]        = "LI";
  menu_m3u["_configKey"]      = "files.m3u";
  menu_m3u["_text"]           = "Playlist";
  menu_m3u["_icon"]           = "img/m3u.png";
  menu_m3u["_headline"]       = "Playlists: Local or remote";
  menu_m3u["_usage"]          = "<b>Info</b><br>Availability: File availability in percent<br>Streams:      Number of streams in the file.<br>group-title:  Streams that are assigned to a group. Simplifies filtering streams<br>tvg-id:       This ID is used for automatic mapping, must match with the channel ID in the XMLTV file.<br>Unique ID:    Streams with a unique ID to identify them. Allows channel name changes in the M3U without losing the XMLTV mapping (PPV / live events).<br><br><b>Usage M3U:</b><br>Remote playlist: http://your.iptv.provider.com/file.m3u<br>Local  playlist: /path/to/file.m3u<br><br><b>Usage HDHomeRun:</b><br>IP: 192.168.1.10:5004<br>"
  menu_m3u["name"]            = "file";
  menu_m3u["id"]              = "file";
  menu_m3u["value"]           = menu_m3u["name"];
  menu_m3u["placeholder"]     = "Playlist: local or remote";
  menu_m3u["onclick"]         = "javascript: toggleMenu(this);";
  menu_m3u["class"]           = "menu-notActive";


  var menu_filter = new Object();
  menu_filter["_menuType"]    = "inputArray";
  menu_filter["_element"]     = "LI";
  menu_filter["_configKey"]   = "filter";
  menu_filter["_text"]        = "Filter";
  menu_filter["_icon"]        = "img/filter.png";
  menu_filter["_headline"]    = "Filter by M3U parameters, e.g. group-title";
  menu_filter["_usage"]       = "<b>Usage:</b><br>Sport - All sports channels<br>Sport {HD} - All HD sports channels<br>Sport {HD} !{ES,DE} - All HD sports channels, but no Spanish and German<br><br>To filter the streams of a HDHomeRun, the playlist name can be entered:<br>My tuner {HD}"
  //menu_filter["_usage"]       = "<b>Usage:</b><br>All sports channels: Sport<br>All HD sports channels: Sport {HD}<br>All HD sports channels, but no Spanish and German: Sport {HD} !{ES,DE}"
  menu_filter["name"]         = "filter";
  menu_filter["id"]           = "M3U";
  menu_filter["value"]        = menu_filter["name"];
  menu_filter["placeholder"]  = "Filter streams: Sport";
  menu_filter["onclick"]      = "javascript: toggleMenu(this);";
  menu_filter["class"]        = "menu-notActive";

  var menu_id = new Object();
  menu_id["_menuType"]        = "inputArray";
  menu_id["_element"]         = "LI";
  menu_id["_configKey"]       = "id";
  menu_id["_text"]            = "PMS ID";
  menu_id["_icon"]            = "img/number.png";
  menu_id["_headline"]        = "Setup PMS guide number";
  menu_id["_usage"]           = 'Some playlists have unique channel IDs.<br>Enter the keyword of the ID. The channel assignment in PMS will change as a result.<br><br>e.g. channelID<br>#EXTINF:0 type="stream" <b>channelId</b>="81", My Streaming Channel HD<br><br>Only enter here if you know what you are doing!'
  menu_id["name"]             = "id";
  menu_id["id"]               = "id";
  menu_id["value"]            = menu_id["name"];
  menu_id["placeholder"]      = "Unique ID from the M3U file";
  menu_id["onclick"]          = "javascript: toggleMenu(this);";
  menu_id["class"]            = "menu-notActive";


  var menu_xmltv = new Object();
  menu_xmltv["_menuType"]     = "inputArray";
  menu_xmltv["_element"]      = "LI";
  menu_xmltv["_configKey"]    = "files.xmltv";
  menu_xmltv["_text"]         = "XMLTV";
  menu_xmltv["_icon"]         = "img/xmltv.png";
  menu_xmltv["_headline"]     = "XMLTV files: Local or remote";
  menu_xmltv["_usage"]        = "<b>Info:</b><br>Availability: File availability in percent<br>Channels:     Number of channels in the file<br>Programs:     Number of EPG data<br><br><b>Usage:</b><br>Remote XMLTV file: http://your.epg.provider.com/guide.xml<br>Local  XMLTV file: /path/to/guide.xml"
  menu_xmltv["name"]          = "xmltv";
  menu_xmltv["id"]            = "xmltv";
  menu_xmltv["value"]         = menu_xmltv["name"];
  menu_xmltv["placeholder"]   = "XMLTV File: local or remote";
  menu_xmltv["onclick"]       = "javascript: toggleMenu(this);";
  menu_xmltv["class"]         = "menu-notActive";

  menu_mapping = new Object();
  menu_mapping["_element"]   = "LI";
  menu_mapping["_text"]      = "Mapping";
  menu_mapping["_icon"]      = "img/mapping.png";
  menu_mapping["_configKey"] = "mapping";
  menu_mapping["_headline"]  = "XMLTV assignment and sorting of channels";
  menu_mapping["id"]         = "mapping";
  menu_mapping["onclick"]    = "javascript: toggleMenu(this);";
  menu_mapping["class"]      = "menu-notActive phone";

  menu_users = new Object();
  menu_users["_element"]   = "LI";
  menu_users["_text"]      = "Users";
  menu_users["_icon"]      = "img/users.png";
  menu_users["_configKey"] = "users";
  menu_users["_headline"]  = "Administration of users and permissions";
  menu_users["id"]         = "users";
  menu_users["onclick"]    = "javascript: toggleMenu(this);";
  menu_users["class"]      = "menu-notActive";
  menu_users["_usage"]     = "<b>Authorization groups:</b><br>WEB: Users can log in to the web interface<br>PMS: Programs like Plex can access the channel list. Login via DVR IP: username:password@xteve.ip:port<br>M3U: Allows clients to download the M3U playlist.<br>XML: Allows clients to download the XMLTV file.<br>API: Allows clients to use the API interface.<br><br>!!! For PMS authentication, only the following special characters are valid: !$()=.,-:;<br><br>The individual authentication groups can be activated / deactivated in the settings menu."
  
  menu_settings = new Object();
  menu_settings["_element"]   = "LI";
  menu_settings["_text"]      = "Settings";
  menu_settings["_icon"]      = "img/settings.png";
  menu_settings["_configKey"] = "settings";
  menu_settings["_headline"]  = "Settings";
  menu_settings["_subMenu"]   = "701,702,703,704,705,706,707,708,799,710,711,712,713,714";
  menu_settings["id"]         = "settings";
  menu_settings["onclick"]    = "javascript: toggleMenu(this);";
  menu_settings["class"]      = "menu-notActive";

  menu_log = new Object();
  menu_log["_element"]        = "LI";
  menu_log["_text"]           = "Log";
  menu_log["_icon"]           = "img/log.png";
  menu_log["_headline"]       = "Log";
  menu_log["_configKey"]      = "log";
  menu_log["id"]              = "log";
  menu_log["onclick"]         = "javascript: toggleMenu(this);";
  menu_log["class"]           = "menu-notActive";

  menu_logout = new Object();
  menu_logout["_element"]   = "LI";
  menu_logout["_text"]      = "Logout";
  menu_logout["_icon"]      = "img/logout.png";
  menu_logout["id"]         = "logout";
  menu_logout["onclick"]    = "javascript: logout();";
  menu_logout["class"]      = "menu-notActive";

  var menu_schedule = new Object();
  menu_schedule["_menuType"]  = "inputArray";
  menu_schedule["_element"]   = "LI";
  menu_schedule["_configKey"] = "update";
  menu_schedule["_text"]      = "Schedule";
  menu_schedule["_icon"]      = "img/schedule.png";
  menu_schedule["_headline"]  = "Schedule for updating M3U, XMLTV files and creating a local backup";
  menu_schedule["_usage"]     = "<b>Usage:</b><br>0815 = 8:15 am<br>1930 = 7:30 pm"
  menu_schedule["name"]       = "update";
  menu_schedule["id"]         = "update";
  menu_schedule["value"]      = menu_id["name"];
  menu_schedule["placeholder"]= "time of day (24-hour clock)";
  menu_schedule["onclick"]    = "javascript: toggleMenu(this);";
  menu_schedule["class"]      = "menu-notActive";

  var menu_filesUpdate = new Object();
  menu_filesUpdate["_element"] = "LI";
  menu_filesUpdate["_menuType"]     = "checkbox";
  menu_filesUpdate["_configKey"]    = "files.update";
  menu_filesUpdate["_label"]        = "Update the provider files at system startup";
  menu_filesUpdate["_headline"]     = "Update the provider files at system startup";
  menu_filesUpdate["_usage"]        = "Playlists and XMLTV files are updated by xTeVe at system startup."
  menu_filesUpdate["name"]          = "files.update";
  menu_filesUpdate["id"]            = "files.update";
  menu_filesUpdate["value"]         = menu_filesUpdate["name"];
  menu_filesUpdate["onclick"]       = "javascript: toggleMenu(this);";
  menu_filesUpdate["class"]         = "menu-notActive";
  
  var menu_tuner = new Object();
  menu_tuner["_element"]      = "LI";
  menu_tuner["_menuType"]     = "select";
  menu_tuner["_configKey"]    = "tuner";
  menu_tuner["_label"]        = "Available tuners";
  menu_tuner["_text"]         = "Tuner";
  menu_tuner["_icon"]         = "img/tuner.png";
  menu_tuner["_headline"]     = "Number of tuners";
  menu_tuner["_usage"]        = "This setting is only used by Plex and Emby.<br>The number of concurrent streams allowed by the IPTV provider.<br>After a change, xTeVe must be delete in the PMS DVR settings and set up again."
  menu_tuner["name"]          = "tuner";
  menu_tuner["id"]            = "tuner";
  menu_tuner["value"]         = menu_tuner["name"];
  menu_tuner["placeholder"]   = "Number of tuners";
  menu_tuner["onclick"]       = "javascript: toggleMenu(this);";
  menu_tuner["class"]         = "menu-notActive";

  var optionValues = new Array();
  for (var i = 1; i <= 100; i++) {
    optionValues.push(i)
  }
  menu_tuner["_optionValues"] = optionValues;
  
  var menu_epg = new Object();
  menu_epg["_element"] = "LI";
  menu_epg["_menuType"]     = "select";
  menu_epg["_configKey"]    = "epgSource";
  menu_epg["_label"]        = "Selection of the EPG source";
  menu_epg["_text"]         = "EPG source";
  menu_epg["_headline"]     = "Selection of the EPG source";
  menu_epg["_usage"]        = "PMS:   Use EPG data from Plex or Emby.<br>XEPG:  Use of external EPG data (XMLTV).<br>       Several XMLTV sources possible.<br>       Allows editing and order channels.<br>       M3U / XMLTV export (HTTP link for IPTV apps)."
  menu_epg["name"]          = "epgSource";
  menu_epg["id"]            = "epgSource";
  menu_epg["value"]         = menu_epg["name"];
  menu_epg["placeholder"]   = "EPG source";
  menu_epg["onclick"]       = "javascript: toggleMenu(this);";
  menu_epg["class"]         = "menu-notActive";
  menu_epg["_optionValues"] = new Array("PMS", "XEPG");

  var menu_xepg = new Object();
  menu_xepg["_element"] = "LI";
  menu_xepg["_menuType"]     = "checkbox";
  menu_xepg["_configKey"]    = "xteveAutoUpdate";
  menu_xepg["_label"]        = "Automatic update of xTeVe";
  menu_xepg["_headline"]     = "Automatic update of xTeVe";
  menu_xepg["_usage"]        = "If a new version of xTeVe is available, it will be automatically installed."
  menu_xepg["name"]          = "xteveAutoUpdate";
  menu_xepg["id"]            = "xteveAutoUpdate";
  menu_xepg["value"]         = menu_xepg["name"];
  menu_xepg["onclick"]       = "javascript: toggleMenu(this);";
  menu_xepg["class"]         = "menu-notActive";

  var menu_autoBackupPath = new Object();
  menu_autoBackupPath["_element"]   = "LI";
  menu_autoBackupPath["_menuType"]  = "singleInput";
  menu_autoBackupPath["_configKey"] = "backup.path";
  menu_autoBackupPath["_label"]     = "Location for automatic backups";
  menu_autoBackupPath["_headline"]  = "Location for automatic backups";
  menu_autoBackupPath["_usage"]     = "Before any update of the provider data by the schedule, xTeVe creates a backup. The path for the automatic backups can be changed. xTeVe requires write permission for this folder."
  menu_autoBackupPath["name"]       = "backup.path";
  menu_autoBackupPath["id"]         = "backup.path";
  menu_autoBackupPath["value"]      = menu_autoBackupPath["name"];
  menu_autoBackupPath["onclick"]    = "javascript: toggleMenu(this);";
  menu_autoBackupPath["class"]      = "menu-notActive";

  var menu_autoBackupKeep = new Object();
  menu_autoBackupKeep["_element"]   = "LI";
  menu_autoBackupKeep["_menuType"]  = "select";
  menu_autoBackupKeep["_configKey"] = "backup.keep";
  menu_autoBackupKeep["_text"]      = "Keep";
  menu_autoBackupKeep["_label"]     = "Number of backups to keep";
  menu_autoBackupKeep["_headline"]  = "Number of backups to keep";
  menu_autoBackupKeep["_usage"]     = ""
  menu_autoBackupKeep["name"]       = "backup.keep";
  menu_autoBackupKeep["id"]         = "backup.keep";
  menu_autoBackupKeep["value"]      = menu_autoBackupKeep["name"];
  menu_autoBackupKeep["onclick"]    = "javascript: toggleMenu(this);";
  menu_autoBackupKeep["class"]      = "menu-notActive";

  var optionValues = new Array(5, 10, 20, 30, 40, 50);
  menu_autoBackupKeep["_optionValues"] = optionValues;


  var menu_buffer = new Object();
  menu_buffer["_element"] = "LI";
  menu_buffer["_menuType"]     = "checkbox";
  menu_buffer["_configKey"]    = "buffer";
  menu_buffer["_label"]        = "Stream buffering [Experimental]";
  menu_buffer["_headline"]     = "Stream buffering [Experimental]";
  menu_buffer["_usage"]        = "With activated buffer, streams can be played and recorded more fluently.<br>The stream is passed from xTeVe to Plex / Emby"
  menu_buffer["name"]          = "buffer";
  menu_buffer["id"]            = "buffer";
  menu_buffer["value"]         = menu_buffer["name"];
  menu_buffer["onclick"]       = "javascript: toggleMenu(this);";
  menu_buffer["class"]         = "menu-notActive";

  var menu_api = new Object();
  menu_api["_element"] = "LI";
  menu_api["_menuType"]     = "checkbox";
  menu_api["_configKey"]    = "api";
  menu_api["_label"]        = "API interface";
  menu_api["_headline"]     = "API interface";
  menu_api["_usage"]        = 'Via API interface it is possible to send commands to xTeVe. API documentation is available <a href="https://xteve.de?scroll=api">here</a> '
  //menu_api["_usage"]        = 'Via API interface it is possible to send commands to xTeVe. API documentation is available <a href="http://localhost:1313?scroll=api">here</a> '
  menu_api["name"]          = "api";
  menu_api["id"]            = "api";
  menu_api["value"]         = menu_api["name"];
  menu_api["onclick"]       = "javascript: toggleMenu(this);";
  menu_api["class"]         = "menu-notActive";

  var menu_authenticationWeb = new Object();
  menu_authenticationWeb["_element"] = "LI";
  menu_authenticationWeb["_menuType"]     = "checkbox";
  menu_authenticationWeb["_configKey"]    = "authentication.web";
  menu_authenticationWeb["_label"]        = "User authentication";
  menu_authenticationWeb["_headline"]     = "User authentication";
  menu_authenticationWeb["_usage"]        = "Access to xTeVe requires authentication."
  menu_authenticationWeb["name"]          = "authentication.web";
  menu_authenticationWeb["id"]            = "authentication.web";
  menu_authenticationWeb["value"]         = menu_authenticationWeb["name"];
  menu_authenticationWeb["onclick"]       = "javascript: toggleMenu(this);";
  menu_authenticationWeb["class"]         = "menu-notActive";
  
  var menu_authenticationPms = new Object();
  menu_authenticationPms["_element"] = "LI";
  menu_authenticationPms["_menuType"]     = "checkbox";
  menu_authenticationPms["_configKey"]    = "authentication.pms";
  menu_authenticationPms["_label"]        = "Plex authentication.";
  menu_authenticationPms["_headline"]     = "Plex authentication.";
  menu_authenticationPms["_usage"]        = "Plex requests are only possible with authentication.<br>Warning!!! After activating this function xTeVe must be delete in the PMS DVR settings and set up again."
  menu_authenticationPms["name"]          = "authentication.pms";
  menu_authenticationPms["id"]            = "authentication.pms";
  menu_authenticationPms["value"]         = menu_authenticationPms["name"];
  menu_authenticationPms["onclick"]       = "javascript: toggleMenu(this);";
  menu_authenticationPms["class"]         = "menu-notActive";
  
  var menu_authenticationM3u = new Object();
  menu_authenticationM3u["_element"] = "LI";
  menu_authenticationM3u["_menuType"]     = "checkbox";
  menu_authenticationM3u["_configKey"]    = "authentication.m3u";
  menu_authenticationM3u["_label"]        = "M3U authentication.";
  menu_authenticationM3u["_headline"]     = "M3U authentication.";
  menu_authenticationM3u["_usage"]        = "Downloading the M3U file via an HTTP request is only possible with authentication."
  menu_authenticationM3u["name"]          = "authentication.m3u";
  menu_authenticationM3u["id"]            = "authentication.m3u";
  menu_authenticationM3u["value"]         = menu_authenticationM3u["name"];
  menu_authenticationM3u["onclick"]       = "javascript: toggleMenu(this);";
  menu_authenticationM3u["class"]         = "menu-notActive";
  

  var menu_authenticationXml = new Object();
  menu_authenticationXml["_element"] = "LI";
  menu_authenticationXml["_menuType"]     = "checkbox";
  menu_authenticationXml["_configKey"]    = "authentication.xml";
  menu_authenticationXml["_label"]        = "XEPG authentication";
  menu_authenticationXml["_headline"]     = "XEPG authentication";
  menu_authenticationXml["_usage"]        = "Downloading the XEPG (XMLTV) file via an HTTP request is only possible with authentication."
  menu_authenticationXml["name"]          = "authentication.xml";
  menu_authenticationXml["id"]            = "authentication.xml";
  menu_authenticationXml["value"]         = menu_authenticationXml["name"];
  menu_authenticationXml["onclick"]       = "javascript: toggleMenu(this);";
  menu_authenticationXml["class"]         = "menu-notActive";

  var menu_authenticationApi = new Object();
  menu_authenticationApi["_element"] = "LI";
  menu_authenticationApi["_menuType"]     = "checkbox";
  menu_authenticationApi["_configKey"]    = "authentication.api";
  menu_authenticationApi["_label"]        = "API authentication";
  menu_authenticationApi["_headline"]     = "API authentication";
  menu_authenticationApi["_usage"]        = "Access to the API interface is only possible with authentication."
  menu_authenticationApi["name"]          = "authentication.api";
  menu_authenticationApi["id"]            = "authentication.api";
  menu_authenticationApi["value"]         = menu_authenticationApi["name"];
  menu_authenticationApi["onclick"]       = "javascript: toggleMenu(this);";
  menu_authenticationApi["class"]         = "menu-notActive";
  
  
  // Main menu
  menu[10] = menu_m3u;

  switch(config["epgSource"]) {
    case "PMS":
      menu[20] = menu_id;
      break;
    
    case "XMLTV":
      menu[40] = menu_xmltv;
      break;

    case "XEPG":
      menu[40] = menu_xmltv;
      menu[50] = menu_mapping;
      break;
  }
  
  menu[30] = menu_filter;
  
  if (config["authentication.web"] == true) {
    menu[60] = menu_users;
  }

  menu[70] = menu_settings;
  menu[80] = menu_log;
  if (config["authentication.web"] == true) {
    menu[100] = menu_logout;
  } 
  

  // Sub-Menu

  subMenu[701] = menu_schedule;
  subMenu[702] = menu_filesUpdate;
  subMenu[703] = menu_tuner;
  subMenu[704] = menu_epg;
  subMenu[705] = menu_xepg;
  subMenu[706] = menu_autoBackupPath;
  subMenu[707] = menu_autoBackupKeep;
  subMenu[708] = menu_buffer;
  
  subMenu[710] = menu_authenticationWeb;
  
  if (config["authentication.web"] == true) {
    subMenu[711] = menu_authenticationPms;
    subMenu[712] = menu_authenticationM3u;
    subMenu[713] = menu_authenticationXml;
    subMenu[714] = menu_authenticationApi;
  }

  subMenu[799] = menu_api;

  

  return
}

function createMenu() {

  showElement("popup", false);

  //console.log(config);
  setMenuItem();
  var menuItems = getObjKeys(menu)
  var nav = document.getElementsByTagName("NAV")[0];
  nav.innerHTML = "";
  var newItem = new Object();

  for (var i = 0; i < menuItems.length; i++) {

    
    var newItem = menu[menuItems[i]];
    newItem["id"] = menuItems[i];

    
    switch(newItem.hasOwnProperty("_icon")) {
      case true: 
        var itemText = newItem["_text"];
        delete newItem["_text"]
        nav.appendChild(createElement(newItem));
        newItem["_text"] = itemText;
        var newIcon = new Object();
        newIcon["_element"] = "IMG";
        newIcon["src"] = newItem["_icon"];

        var currentElement = document.getElementById(menuItems[i]);
        currentElement.appendChild(createElement(newIcon));


        var text = new Object();
        text["_element"] = "P"
        text["_text"] = itemText;
        text["class"] = "nav-text"
        currentElement.appendChild(createElement(text));
        break;

      default:
        nav.appendChild(createElement(newIcon));
        break;
    }

  }
  if (activeMenu != undefined) {
    //console.log(activeMenu);
    toggleMenu(activeMenu);
  }

  return
}

function toggleMenu(elm) {
  //showStreams(false);
  clearInterval(logInterval)
  activeMenu = elm;
  var item = menu[elm.id]
  var div = document.getElementById("settings");
  div.innerHTML = "";
  
  // Set Headline
  var headline = new Object();
  headline["_element"] = "H4";
  headline["_text"] = item["_headline"];
  div.appendChild(createElement(headline));

  // Sub-Menu
  if (item.hasOwnProperty("_subMenu") == true) {
    openSubMenu(item);
    return
  }

  // Mapping, Users, Log, Files
  switch(item["_configKey"]) {
    case "mapping":     openMappingEditor(item); return; break;
    case "users":       openUsers(item); return; break;
    case "log":         showLog(item); return; break;
    case "files.m3u":   openFiles(item, "m3u"); return; break;
    case "files.xmltv": openFiles(item, "xmltv"); return; break;

    case "filter":      showStreams(true); break;
  }

 

  var newHR = new Object();
  newHR["_element"] = "HR"
  div.appendChild(createElement(newHR));
  
  var newEntry = new Object();
  newEntry["_element"]  = "INPUT";
  newEntry["type"] = "button";
  //newEntry["class"] = "save";
  newEntry["value"] = "Save";
  newEntry["onclick"] = "saveData2('settings')"
  div.appendChild(createElement(newEntry));


  var newWrapper = new Object();
  newWrapper["_element"]  = "DIV";
  newWrapper["id"]        = "box-wrapper";
  div.appendChild(createElement(newWrapper));

  div = div.lastChild;
  
  div.appendChild(createMenuItem(item))

  // usage Info  
  switch(menu[activeMenu.id].hasOwnProperty("_usage")) {
    case true: 
      var usageItem = new Object();
      usageItem["_element"] = "PRE"
      usageItem["_text"]    = menu[activeMenu.id]["_usage"];
      div.appendChild(createElement(usageItem));
  }
  
  calculateWrapperHeight();

}

function createMenuItem(item) {
  var element = document.createElement("DIV");
  switch(item["_menuType"]) {
    case "inputArray":
      if (config.hasOwnProperty(item["_configKey"]) == true) {
        var value = config[item["_configKey"]];
      } else {
        var value = new Array();
      }
      
      for (var i = 0; i < value.length; i++) {
        var newEntry = new Object();
        newEntry = item
        delete newEntry["onclick"];
        newEntry["_element"]  = "INPUT";
        newEntry["value"]     = value[i];
        newEntry["type"]      = "search";
        newEntry["data-menutype"] = item["_menuType"];
        newEntry["data-menukey"] = item["_configKey"];
        element.appendChild(createElement(newEntry));

      }
      // New entry for array
      var newEntry = new Object();
      newEntry["_element"]      = "INPUT";
      newEntry["type"]          = "search";
      newEntry["name"]          = item["name"];
      newEntry["placeholder"]   = item["placeholder"];
      newEntry["value"]         = "";
      newEntry["data-menutype"] = item["_menuType"];
      newEntry["data-menukey"]  = item["_configKey"];
      element.appendChild(createElement(newEntry));
      break;
  
    case "singleInput":
      var value = config[item["_configKey"]];
      if (value == undefined) {
        value = "";
      }
      var newEntry = new Object();
      newEntry = item;
      delete newEntry["onclick"];
      newEntry["_element"]  = "INPUT";
      newEntry["value"]     = value;
      newEntry["type"]      = "search";
      newEntry["data-menutype"] = item["_menuType"];
      newEntry["data-menukey"] = item["_configKey"];
      element.appendChild(createElement(newEntry));
      break;

    case "checkbox":
      var value = config[item["_configKey"]];
      if (value == undefined) {
        value = false;
      }
      var newEntry = new Object();
      newEntry = item;
      delete newEntry["onclick"];
      newEntry["_element"]  = "INPUT";
      newEntry["value"]     = value;
      newEntry["type"]      = "checkbox";
      newEntry["data-menutype"] = item["_menuType"];
      newEntry["data-menukey"] = item["_configKey"];
      element.appendChild(createElement(newEntry));
      element.getElementsByTagName("INPUT")[0].checked = value;
      break;
    
    case "select":
      var value = config[item["_configKey"]];
      var newEntry = new Object();
      newEntry = item;
      delete newEntry["onclick"]
      newEntry["_element"]  = "SELECT";
      element.appendChild(createElement(newEntry));
      var selectElement = element.getElementsByTagName("SELECT")[0];
      var values = item["_optionValues"];
      for (var i = 0; i < values.length; i++) {
        var newEntry = new Object;
        newEntry["_element"]  = "OPTION";
        newEntry["_text"]     = item["_text"] + ": " + values[i];
        newEntry["value"]     = values[i];
        selectElement.appendChild(createElement(newEntry));
      }
      selectElement.value = value;
      break;
    
  }
  return element;
}

function openSubMenu(item) {
  var entrys = item["_subMenu"].split(",");
  var div = document.getElementById("settings");

  var newHR = new Object();
  newHR["_element"] = "HR"
  div.appendChild(createElement(newHR));
  
  var newEntry = new Object();
  newEntry["_element"]  = "INPUT";
  newEntry["type"] = "button";
  //newEntry["class"] = "save";
  newEntry["value"] = "Save";
  newEntry["onclick"] = "saveData2('settings')"
  div.appendChild(createElement(newEntry));

  if (item["_configKey"] == "settings") {
    var newEntry = new Object();
    newEntry["_element"]  = "INPUT";
    newEntry["type"] = "button";
    //newEntry["class"] = "save";
    newEntry["value"] = "Backup";
    newEntry["onclick"] = "xteveBackup()"
    div.appendChild(createElement(newEntry));
  }

  if (item["_configKey"] == "settings") {
    var newEntry = new Object();
    newEntry["_element"]  = "INPUT";
    newEntry["type"] = "button";
    //newEntry["class"] = "save";
    newEntry["value"] = "Restore";
    newEntry["onclick"] = "xteveRestore(this)"
    div.appendChild(createElement(newEntry));
  }


  var newWrapper = new Object();
  newWrapper["_element"]  = "DIV";
  newWrapper["id"]        = "box-wrapper";
  div.appendChild(createElement(newWrapper));
  
  div = div.lastChild;
  

  for (var i = 0; i < entrys.length; i++) {
    var item = subMenu[entrys[i]];
    if (item == undefined) {
      break;
    }
    
    var container = new Object();
    container["_element"] = "DIV";
    div.appendChild(createElement(container));

    var divContainer = div.lastChild;
    
    var headline = new Object();
    headline["_element"] = "H5";
    headline["_text"] = item["_headline"];
    divContainer.appendChild(createElement(headline));

    divContainer.appendChild(createMenuItem(item))

    switch(item.hasOwnProperty("_usage")) {
      case true: 
        var usageItem = new Object();
        usageItem["_element"] = "PRE"
        usageItem["_text"]    = item["_usage"];
        divContainer.appendChild(createElement(usageItem));
    }

    var hr = new Object();
    hr["_element"] = "HR";
    divContainer.appendChild(createElement(hr));
  
  }

  calculateWrapperHeight();
  return
}

function saveData2(elm) {
  var div   = document.getElementById(elm);
  var inputs = div.getElementsByTagName("INPUT");
  var selects = div.getElementsByTagName("SELECT");
  var value, configKey;
  var data = new Object();
  var valueArr = new Array();
  var newData = false;
  
  for (var i = 0; i < inputs.length; i++) {
    if (inputs[i].type != "button") {
      var menuType = inputs[i].getAttribute("data-menutype");
      
      //console.log(menuType);
      switch(menuType) {
        case "singleInput":
          value = inputs[i].value;
          if (value == "" || value == undefined) {
            data = new Object();
            data["delete"] = inputs[i].name
            newData = true;
          } else {
            newData = true;
            data[inputs[i].name] = value;
            console.log(data);
          }
          break;
        case "inputArray": 
          value = inputs[i].value;
          if (value != "" && value != undefined) {
            newData = true;
            valueArr.push(value)
            data[inputs[i].name] = valueArr;
            configKey = inputs[i].name;
          } 
          
          break;

        case "checkbox":
          value = inputs[i].checked
          data[inputs[i].name] = value;
      }
      
    }
    
  }


  // Delete config key
  if (valueArr.length == 0 && newData == false) {
    newData = true;
    data = new Object();
    data["delete"] = configKey;
  } 


  for (var i = 0; i < selects.length; i++) {
    var value = selects[i].options[selects[i].selectedIndex].value;
    switch(isNaN(value)) {
      case false: value = parseInt(value); break;
    }

    data[selects[i].name] = value;
    newData = true;
  }

  //console.log(data, newData);

  if (newData == true) {
    data["cmd"] = "saveConfig";
    if (!data.hasOwnProperty('filter')) {
      data["filter"] = config["filter"]
    }
    var settings = new Object();
    settings["cmd"] = data["cmd"];
    settings["settings"] = data;
    console.log(settings);
    xTeVe(settings);
  }
}
