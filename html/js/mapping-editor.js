var mappingError  = false;
var bulk          = false;
var bulkEditAll   = false; 
var selectObj     = new Object();
var searchObj     = new Object();

var bulkIDs       = new Array();
var bulkChangeObj = new Object();

function checkUndo(key, elm) {
  var tmp = new Object();
  tmp = elm
  console.log("--");
  if (undo.hasOwnProperty("epgMapping")) {
    xEPG["epgMapping"] = JSON.parse(JSON.stringify(undo["epgMapping"]));;
  } else {
    undo["epgMapping"] = JSON.parse(JSON.stringify(elm));
  }
}

//var plexCategories = new Array("-", "Action sports", "Action", "Adults only", "Adventure", "Aerobics", "Animals", "Animated", "Anime", "Anthology", "Archery", "Art", "Arts/crafts", "Auction", "Auto racing", "Auto", "Aviation", "Awards", "Ballet", "Baseball", "Basketball", "Bicycle racing", "Bicycle", "Billiards", "Biography", "Boat racing", "Boat", "Bowling", "Boxing", "Bus./financial", "Children", "Collectibles", "Comedy drama", "Comedy", "Community", "Computers", "Consumer", "Cooking", "Crime drama", "Crime", "Dance", "Dark comedy", "Debate", "Diving", "Docudrama", "Documentary", "Drama", "Educational", "Entertainment", "Environment", "Equestrian", "Erotic", "Event", "Fantasy", "Fashion", "Feature Film", "Fishing", "Football", "Game show", "Gaming", "Gay/lesbian", "Golf", "Handball", "Health", "Historical drama", "History", "Hockey", "Holiday", "Home improvement", "Horror", "Horse", "House/garden", "How-to", "Interview", "Intl soccer", "Law", "Martial arts", "Medical", "Military", "Miniseries", "Mixed martial arts", "Motorcycle racing", "Motorcycle", "Motorsports", "Mountain biking", "Music", "Musical comedy", "Musical", "Mystery", "Nature", "News", "Newsmagazine", "Olympics", "Opera", "Outdoors", "Parade", "Paranormal", "Parenting", "Performing arts", "Playoff sports", "Poker", "Politics", "Pro wrestling", "Public affairs", "Reality", "Religious", "Rodeo", "Roller derby", "Romance", "Romantic comedy", "Rugby", "Running", "Sailing", "Science fiction", "Science", "Self improvement", "Series", "Shooting", "Shopping", "Short Film", "Sitcom", "Skiing", "Snooker", "Soap", "Soccer", "Special", "Sports", "sports", "Sports event", "Sports non-event", "Sports talk", "Standup", "Surfing", "Suspense", "TV Movie", "Talk", "Technology", "Tennis", "Theater", "Thriller", "Track/field", "Travel", "Triathlon", "Variety", "Volleyball", "War", "Watersports", "Weather", "Western", "Wrestling", "Yacht racing", "movie", "series", "sports", "tvshow");
var plexCategoriesValues = new Array("-", "Kids", "News", "Movie", "Series", "Sports")
var plexCategoriesOption = new Array("-", "Kids (Emby only)", "News", "Movie", "Series", "Sports")


function openMappingEditor(elm) {
  var columnToSort  = 1

  checkUndo("epgMapping", xEPG["epgMapping"])

  var newDiv = document.getElementById("settings");
  
  var newEntry = new Object();
  newEntry["_element"]  = "HR";
  newDiv.appendChild(createElement(newEntry));

  var newEntry = new Object();
  newEntry["_element"]    = "INPUT";
  newEntry["type"]        = "button";
  newEntry["class"]       = "button";
  newEntry["value"]       = "Save";
  newEntry["onclick"]     = "saveXEPG()";
  newDiv.appendChild(createElement(newEntry));

  var newEntry = new Object();
  newEntry["_element"]    = "INPUT";
  newEntry["type"]        = "button";
  newEntry["class"]       = "button";
  newEntry["value"]       = "Bulk Edit";
  newEntry["onclick"]     = "bulkEdit()";
  newDiv.appendChild(createElement(newEntry));

  var newEntry = new Object();
  newEntry["_element"]    = "INPUT";
  newEntry["type"]        = "button";
  newEntry["class"]       = "button";
  newEntry["value"]       = "Show XEPG";
  newEntry["onclick"]     = "showXEPG()";
  newDiv.appendChild(createElement(newEntry));

  var newEntry = new Object();
  newEntry["_element"]    = "INPUT";
  newEntry["class"]       = "search";
  newEntry["id"]          = "searchMapping";
  newEntry["type"]        = "search";
  newEntry["placeholder"] = "Search";
  newEntry["onchange"]    = "searchInMapping()";
  newDiv.appendChild(createElement(newEntry));

  var div = document.getElementById("settings");
  //screenLog("Duplicate ID", "error", true)
  

  // Build table

  var newWrapper = new Object();
  newWrapper["_element"]  = "DIV";
  newWrapper["id"]        = "box-wrapper";
  div.appendChild(createElement(newWrapper));


  var newTable = new Object();
  newTable["_element"]  = "TABLE";
  newTable["id"]        = "id_mapping";
  newTable["class"]     = "table-mapping";
  div.lastChild.appendChild(createElement(newTable));
  showLoadingScreen(true);

  setTimeout(function(){ 
    createMappingTable(); 
  }, 10);

}

function createSearchObj() {
  searchObj = new Object();
  var IDs = getObjKeys(xEPG["epgMapping"])
  for (var i = IDs.length - 1; i >= 0; i--) {
    var item = xEPG["epgMapping"][IDs[i]];
    var searchID = item["x-epg"];
    var searchValue = ""; 
    searchValue = searchValue + item["x-channelID"] + " ";
    searchValue = searchValue + item["x-category"] + " ";
    searchValue = searchValue + item["x-name"] + " ";
    searchValue = searchValue + item["x-group-title"] + " ";
    searchValue = searchValue + item["x-xmltv-file"] + " ";
    searchValue = searchValue + item["_file.m3u.name"] + " ";

    switch(item["x-active"]) {
      case true:  searchValue = searchValue + "online"; break;
      case false: searchValue = searchValue + "offline"; break;
    }

    searchObj[searchValue] = searchID;

  }
}


function calculateWrapperHeight() {

  if (document.getElementById("box-wrapper")){

    var elm = document.getElementById("box-wrapper");
    
    var divs = new Array("myStreamsBox", "clientInfo", "settings");
    var elementsHeight = 0 - elm.offsetHeight;
    for (var i = 0; i < divs.length; i++) {
      elementsHeight = elementsHeight + document.getElementById(divs[i]).offsetHeight;
    }

    elm.style.height = window.innerHeight - elementsHeight + "px";

  }

  if (document.getElementById("menu-wrapper")){

    var elm = document.getElementById("menu-wrapper");
    
    var offest = document.getElementById("settings").offsetHeight + document.getElementById("myStreamsBox").offsetHeight + document.getElementById("clientInfo").offsetHeight;
    
    if (window.innerHeight > offest) {
      elm.style.height = window.innerHeight + "px"
    } else {
      elm.style.height = offest + "px"
    }
    

  }


}

function createMappingTable() {
  columnToSort = 1;
  createSearchObj();

  // Create table (Header)
  var table = document.getElementById("id_mapping");
  table.innerHTML = "";
  var newTR = new Object();
  newTR["_element"] = "TR";
  newTR["class"]    = "table-mapping-header";
  table.appendChild(createElement(newTR));

  var tr = document.getElementById("id_mapping").lastChild;
  var trHeadlines = new Array("Bulk", "Ch. No.", "Logo", "Channel Name", "Playlist", "Group Title", "XMLTV File", "XMLTV ID", "Timeshift")

  for (var i = 0; i < trHeadlines.length; i++) {
    var newTD = new Object();

    newTD["_element"] = "TD";
    newTD["_text"]    = trHeadlines[i];

    

    var width = "";
    switch(trHeadlines[i]) {

      case "Bulk":  
        
        maxWidth = "32px"; 
        minWidth = "32px"; 
    
        // Create bulk TD
        var newCheckbox = new Object();
        newCheckbox["_element"] = "INPUT";
        newCheckbox["type"]     = "checkbox";
        newCheckbox["class"]    = "bulk hideBulk";
        newCheckbox["onmouseout"] = "javascript: this.blur()"
        newCheckbox["onclick"]    = "javascript: bulkEditAllChannels()"

    
        //newTD.appendChild(createElement(newCheckbox));

        break;

      case "Ch. No.": 
        maxWidth = "80px"; 
        minWidth = "70px"; 
        newTD["onclick"]  = "javscript: sortTable(" + i + ");";
        newTD["class"]    = "pointer";
        break;
      
      case "Logo":  maxWidth = "120px"; minWidth = "60px"; break;
      
      case "Channel Name":  
        maxWidth = "50%"; 
        minWidth = "200px"; 
        newTD["onclick"]  = "javscript: sortTable(" + i + ");";
        newTD["class"] = "pointer";
        break;

      case "Playlist":      
        maxWidth = "150px"; 
        minWidth = "100px"; 
        newTD["onclick"]  = "javscript: sortTable(" + i + ");";
        newTD["class"]    = "pointer";
        break;
      
      case "Group Title":   
        maxWidth = "150px"; 
        minWidth = "100px"; 
        newTD["onclick"]  = "javscript: sortTable(" + i + ");";
        newTD["class"]    = "pointer";
        break;
      
      case "XMLTV File":    
        maxWidth = "150px"; 
        minWidth = "100px"; 
        //newTD["onclick"]  = "javscript: sortTable(" + i + ");";
        newTD["class"]    = "";
        break;
      

      case "XMLTV ID":      maxWidth = "150px"; minWidth = "100px"; break;

      case "Timeshift":      maxWidth = "50px"; minWidth = "20px"; break;

      default: 
        newTD["class"]    = "";
        break;
    }

    tr.appendChild(createElement(newTD));
    if (trHeadlines[i] == "Bulk") {
      tr.lastChild.innerHTML = "";
      tr.lastChild.appendChild(createElement(newCheckbox));
      
    }
    
    var elm = tr.lastChild;
    elm.style.width = maxWidth;
    elm.style.maxWidth = maxWidth;
    elm.style.minWidth = minWidth;

  }
  calculateWrapperHeight();
  var IDs = getObjKeys(xEPG["epgMapping"])

  var allXmltvFiles = getObjKeys(xEPG["xmltvMap"]);

  if (allXmltvFiles == 0) {
    showLoadingScreen(false);
    return;
  }

  // Sort IDs
  var posObj = new Object();
  for (var i = 0; i < IDs.length; i++) {
    var item  = xEPG["epgMapping"][IDs[i]];
    var pos
    switch(isNaN(xEPG["epgMapping"][IDs[i]]["x-channelID"])) {
      case false: pos = parseFloat(xEPG["epgMapping"][IDs[i]]["x-channelID"]) ; break;
    }
    posObj[pos] = item;
  }
  posFloat = getObjKeys(posObj)
  function sortFloat(a,b) { return a - b; }
  posFloat.sort(sortFloat)

  //console.log(posFloat);

  // ---

  if (IDs.length > 200) {
    setTimeout(function(){ 
      showLoadingScreen(true);
    }, 1);

  }


  // table for int channel ID's
  for (var i = 0; i < posFloat.length; i++) {

    var table = document.getElementById("id_mapping");
    var item  = posObj[posFloat[i]];
    //var item  = xEPG["epgMapping"][IDs[i]];
    //console.log(item);
    var newTR = new Object();
    newTR["_element"]       = "TR";
    newTR["class"]          = "";
    newTR["id"]             = item["x-epg"];
    newTR["oncontextmenu"]  = 'javascript: switchChannelStatus("' + item["x-epg"] + '"); return false;';
    table.appendChild(createElement(newTR));

    var tr = document.getElementById("id_mapping").lastChild;
    
    // Create bulk TD
    var newTD = new Object();
    newTD["_element"]   = "INPUT";
    newTD["type"]       = "checkbox";
    newTD["class"]      = "bulk hideBulk";
    newTD["onmouseout"] = "javascript: this.blur()"
    
    createNewTD(newTD, tr);
    

    // Create ID TD
    var newTD = new Object();
    newTD["_element"] = "INPUT";
    newTD["type"]     = "text"
    newTD["class"]    = "w40px";
    newTD["value"]    = item["x-channelID"];
    newTD["onfocusout"] = "javascript: arrangeTable(this);"
    createNewTD(newTD, tr);

    // Create IMG TD
    var newTD = new Object();
    newTD["_element"] = "IMG";
    newTD["onclick"]  = 'javascript: mappingDetail("' + item["x-epg"] + '");';
    if (item["tvg-logo"] != undefined) {
      newTD["src"]      = item["tvg-logo"];
    } else {
      item["tvg-logo"] = "";
      newTD["src"] = "";
    }
    createNewTD(newTD, tr);
    tr.lastChild.setAttribute("onclick", 'javascript: mappingDetail("' + item["x-epg"] + '");')

    // Create P TD (channel name)
    var newTD = new Object();
    newTD["_element"] = "P";
    newTD["_text"]    = item["x-name"];
    newTD["class"]     = item["x-category"];

    createNewTD(newTD, tr);
    tr.lastChild.setAttribute("onclick", 'javascript: mappingDetail("' + item["x-epg"] + '");')
    tr.lastChild.lastChild.style.padding = "5px 10px";

    // Create P TD (Playlist Name)
    var newTD = new Object();
    newTD["_element"] = "P";
    newTD["_text"]    = item["_file.m3u.name"];
    newTD["class"]     = item["tableEllipsis"];
    
    createNewTD(newTD, tr);
    tr.lastChild.setAttribute("onclick", 'javascript: mappingDetail("' + item["x-epg"] + '");')
    
    // Create P TD (Group Title)
    var newTD = new Object();
    newTD["_element"] = "P";
    newTD["_text"]    = item["x-group-title"];
    newTD["class"]     = item["tableEllipsis"];
    
    createNewTD(newTD, tr);
    tr.lastChild.setAttribute("onclick", 'javascript: mappingDetail("' + item["x-epg"] + '");')

    
    // Create P TD (XMLTV file)
    var newTD = new Object();
    newTD["_element"]       = "P";
    newTD["class"]    = "tableEllipsis";
    newTD["_text"] = "-"

    if (allXmltvFiles.indexOf(item["x-xmltv-file"]) != -1) {
      var xXmltvFile = item["x-xmltv-file"];
      switch(xXmltvFile) {
        case "-":           newTD["_text"]  = xXmltvFile; break;
        case "xTeVe Dummy": newTD["_text"]  = xXmltvFile; break;
        default:            newTD["_text"]  = getValueFromProviderFile(xXmltvFile, "xmltv", "name"); break;
        
      }
      //console.log(newTD);

      //newTD["_text"]    = item["x-xmltv-file"];
    } else {
      //newTD["_text"] = "-"
    }
    createNewTD(newTD, tr);
    tr.lastChild.setAttribute("onclick", 'javascript: mappingDetail("' + item["x-epg"] + '");')

    // Creatr P TD (XMLTV channel ID)
    newTD["_element"] = "P";
    newTD["class"]    = "tableEllipsis";

    if (item["x-mapping"] != undefined) {
      newTD["_text"]    = item["x-mapping"];
    }
    
    createNewTD(newTD, tr);
    tr.lastChild.setAttribute("onclick", 'javascript: mappingDetail("' + item["x-epg"] + '");')


    var xXmltvFile  = item["x-xmltv-file"];
    var xMapping    = item["x-mapping"];
    var tvgID       = item["tvg-id"];
    
    //console.log(item["x-epg"]);
    //console.log(item);

    if (item["x-active"] == true) {
      tr.className = "activeEPG";
    } else {
      tr.className = "notActiveEPG";
    }
    
  }

  sortTable(1);

  setTimeout(function(){ 
    showLoadingScreen(false);
  }, 5);
}

function searchInMapping(elm) {

  var search = document.getElementById("searchMapping").value;
  var values = getObjKeys(searchObj)
  
  for (var i = values.length - 1; i >= 0; i--) {
    var id = searchObj[values[i]];
    var bool = values[i].toLowerCase().includes(search.toLowerCase());
    switch(bool) {
      case true:  document.getElementById(id).style.display = ""; break;
      case false: document.getElementById(id).style.display = "none"; break;
    }
  }

}

function mappingDetail(xepg) {
  
  bulkIDs   = new Array();
  var activeElement = document.activeElement;
  // If input id, return
  if (activeElement.tagName == "INPUT") {
    return
  }

  if (bulk == true) {
    var elm = document.getElementsByClassName("bulk");
    for (var i = 1; i < elm.length; i++) {
      if (elm[i].checked == true) {
        var id = elm[i].parentElement.parentElement.id;
        bulkIDs.push(id)
      }
      
    }

    if (bulkIDs.length == 0) {
      showElement('popup', false)
      alert("No channels selected for editing")
      return
    }

    xepg = bulkIDs[0]
  }


  createSearchObj();
  
  showPopUpElement('mapping-detail');

  var thisChannel = xEPG["epgMapping"][xepg];
  //console.log(thisChannel);
  var xXmltvFile  = thisChannel["x-xmltv-file"];
  var xMapping    = thisChannel["x-mapping"];
  var xCategory   = thisChannel["x-category"];

  if (xXmltvFile == undefined) {
    thisChannel["x-xmltv-file"] = "-";
    xXmltvFile = "-";
  }

  if (xMapping == undefined) {
    thisChannel["x-mapping"] = "-";
    xMapping = "-";
  }

  /*
  console.log("ID:", xepg);
  console.log("XMLTV File:", xXmltvFile);
  console.log("Mapping:", xMapping);
  */

  var keys = getObjKeys(thisChannel);
  for (var i = 0; i < keys.length; i++) {
    if(document.getElementById(keys[i])){
      var td = document.getElementById(keys[i])
    } else {
      var td = undefined;
    }
    
    var newItem = new Object();
    var values, text = new Array();
    switch(keys[i]) {
      case "x-xmltv-file": 
        var fileIDs = getObjKeys(xEPG["xmltvMap"]);
        var value = new Array("-");
        var text  = new Array("-");

        for (var j = fileIDs.length - 1; j >= 0; j--) {
          if (fileIDs[j] != "xTeVe Dummy") {
            value.push(getValueFromProviderFile(fileIDs[j], "xmltv", "file.xteve"))
            text.push(getValueFromProviderFile(fileIDs[j], "xmltv", "name"))
          } else {
            value.push(fileIDs[j])
            text.push(fileIDs[j])
          }
          
        }
        
        newItem["_element"]       = "SELECT";
        newItem["_optionValues"]  = value;
        newItem["_optionText"]    = text
        newItem["value"]          = xXmltvFile;
        newItem["onchange"]       = 'javascript: changeXmltvFile("' + xepg + '",this);';

        break;

      case "x-mapping": 

        var values = getObjKeys(xEPG["xmltvMap"][xXmltvFile]);

        for (var j = 0; j < values.length; j++) {
          
          if (xEPG["xmltvMap"][xXmltvFile][values[j]].hasOwnProperty('display-name') == true) {
            var displayName = xEPG["xmltvMap"][xXmltvFile][values[j]]["display-name"];
          } else {
            var displayName = "-"
          }
          
          //text[j] = values[j] + " (" + displayName + ")";
          text[j] = displayName + " (" + values[j]  + ")";
        }

        text.unshift("-");
        values.unshift("-");
        newItem["_element"]       = "SELECT";
        newItem["_optionValues"]  = values;
        newItem["_optionText"]    = text
        newItem["value"]          = xMapping;
        newItem["onchange"]       = 'javascript: mappingChannel("' + xepg + '",this);';
        break;

      case "x-category":
        //var values = plexCategoriesValues
        newItem["_element"]       = "SELECT";
        newItem["_optionValues"]  = plexCategoriesValues;
        newItem["_optionText"]    = plexCategoriesOption;
        newItem["value"]          = xCategory;
        newItem["onchange"]       = 'saveCategory("' + xepg + '")';
        break;

      case "tvg-logo":
        document.getElementById("channel-logo").setAttribute("src", thisChannel["tvg-logo"]);
        newItem["_element"]       = "INPUT";
        newItem["type"]           = "text";
        newItem["value"]          = thisChannel["tvg-logo"];
        newItem["onfocusout"]     = 'saveChannelLogo("' + xepg + '")';
        newItem["placeholder"]    = 'Image URL';
        break;

      case "x-update-channel-icon":
        newItem["_element"]       = "INPUT";
        newItem["type"]           = "checkbox";
        switch(JSON.parse(thisChannel["x-update-channel-icon"])) {
          case true: newItem["checked"]        = thisChannel["x-update-channel-icon"];
            break
        }
        newItem["onchange"]     = 'saveChannelIconUpdate("' + xepg + '")';
        break;

      case "x-name":
        newItem["_element"]       = "INPUT";
        newItem["type"]           = "text";
        newItem["value"]          = thisChannel["x-name"];
        newItem["onfocusout"]     = 'saveChannelName("' + xepg + '")';
        newItem["placeholder"]    = 'Channel Name';
        break;

      case "x-update-channel-name":
        if (thisChannel.hasOwnProperty("_uuid.key") == true) {
          newItem["_element"]       = "INPUT";
          newItem["type"]           = "checkbox";
          switch(JSON.parse(thisChannel["x-update-channel-name"])) {
            case true: newItem["checked"]        = thisChannel["x-update-channel-name"];
              break
          }
          newItem["onchange"]     = 'saveChannelNameUpdate("' + xepg + '")';
          showElement("streamHasCUID", true)

          break;
        } else {
          //streamHasCUID
          showElement("streamHasCUID", false)
          break;
        }
        
      case "x-active":
        newItem["_element"]       = "INPUT";
        newItem["type"]           = "checkbox";
        switch(JSON.parse(thisChannel["x-active"])) {
          case true: newItem["checked"]        = thisChannel["x-active"];
            break
        }
        newItem["onchange"]     = 'saveChannelStatus("' + xepg + '")';
        break;

      case "x-group-title":
        newItem["_element"]       = "INPUT";
        newItem["type"]           = "text";
        newItem["value"]          = thisChannel["x-group-title"];
        newItem["onfocusout"]     = 'saveGroupTitle("' + xepg + '")';
        newItem["placeholder"]    = 'Group Title';
        break;

      default:
        newItem["_element"]       = "P";
        newItem["_text"]          = thisChannel[keys[i]];
        break;
      
    }
    
    if (td != undefined) {
      td.innerHTML = "";
      var element = createNewElement(newItem)
      //console.log(element);
      td.appendChild(element);
    }

  }

  if (bulk == true) {

    var elm = document.getElementsByClassName("noBulk");
    for (var i = 0; i < elm.length; i++) {
      elm[i].lastChild.setAttribute("readonly", true)
      elm[i].lastChild.style.borderColor = "red";
    }

    xepg = bulkIDs[0]
  }

  sortSelect(document.getElementById("x-xmltv-file").lastChild);
  sortSelect(document.getElementById("x-mapping").lastChild);
  
}

function sortSelect(elem) {

  var tmpAry = [];
  // Retain selected value before sorting
  var selectedValue = elem[elem.selectedIndex].value;
  // Grab all existing entries
  for (var i=0;i<elem.options.length;i++) tmpAry.push(elem.options[i]);
  // Sort array by text attribute
  tmpAry.sort(function(a,b){ return (a.text < b.text)?-1:1; });
  // Wipe out existing elements
  while (elem.options.length > 0) elem.options[0] = null;
  // Restore sorted elements
  var newSelectedIndex = 0;
  for (var i=0;i<tmpAry.length;i++) {
      elem.options[i] = tmpAry[i];
      if(elem.options[i].value == selectedValue) newSelectedIndex = i;
  }
  elem.selectedIndex = newSelectedIndex; // Set new selected index after sorting
  return;
}


function switchChannelStatus(xepg) {
  var thisChannel = xEPG["epgMapping"][xepg];
  var xXmltvFile = thisChannel["x-xmltv-file"];

  if (xEPG["xmltvMap"].hasOwnProperty(xXmltvFile) == true) {
    if (thisChannel["x-mapping"] != "-" && thisChannel["x-mapping"] != undefined) {
      thisChannel["x-active"] = !thisChannel["x-active"];
      var tr = document.getElementById(xepg);
      switch(thisChannel["x-active"]) {
        case true: tr.className = "activeEPG"; break;
        case false: tr.className = "notActiveEPG"; break;
      }
      document.getElementById("logInfo").className = "notVisible";

    } else {
      var err = "XMLTV Channel is not selected"
      alert(err)
      /*
      var newError = new Object();
      newError["err"] = "Channel is not selected";
      checkErr(newError);
      */
    }

  } else {
    var err = "XMLTV File is not selected"
    alert(err)
    /*
    var newError = new Object();
    newError["err"] = "XMLTV file is not selected";
    checkErr(newError);
    */
  }

  searchInMapping();

}

function createNewElement(newItem) {

  var element = createElement(newItem);
  
  switch(newItem["_element"]) {
    case "SELECT":
      //element[]
      var values  = newItem["_optionValues"];
      var text    = newItem["_optionText"];

      for (var i = 0; i < values.length; i++) {
        //console.log(item);
        var newEntry = new Object;
        newEntry["_element"]  = "OPTION";
        newEntry["_text"]     = text[i];
        newEntry["value"]     = values[i];
        element.appendChild(createElement(newEntry));
      }
      element.value = newItem["value"];
      break;
    
    default: 
      
      //element.appendChild(createElement(newItem));
      break;
  }
  
  return element;
}

function saveBulk(key, value) {
  for (var i = 0; i < bulkIDs.length; i++) {
    var id = bulkIDs[i]
    var thisChannel = xEPG["epgMapping"][id];
    thisChannel[key] = value;

    switch(key) {
      case "tvg-logo":      document.getElementById(id).childNodes[2].lastChild.setAttribute("src", value); break;
      
      case "x-category":    document.getElementById(id).childNodes[3].lastChild.className = value; break;

      case "x-xmltv-file":
        var element = document.getElementById(id).childNodes[6].lastChild;
        switch(value) {
          case "-":           element.innerHTML = value; break;
          case "xTeVe Dummy": element.innerHTML = value; break;
          default:            element.innerHTML = getValueFromProviderFile(value, "xmltv", "name"); break;
        }



      //document.getElementById(id).childNodes[5].lastChild.innerHTML = value.replace(/^.*[\\\/]/, ''); break;
      case "x-mapping":   
        document.getElementById(id).childNodes[7].lastChild.innerHTML = value;
        if (value == "-") {
          thisChannel["x-active"] = false;
          document.getElementById(id).className = "notActiveEPG";
        } else {
          thisChannel["x-active"] = true;
          document.getElementById(id).className = "activeEPG";
        }
        break;

      case "x-group-title": document.getElementById(id).childNodes[5].lastChild.innerHTML = value; break;

      case "x-active":
        var tr = document.getElementById(id);
        
        if (thisChannel.hasOwnProperty("x-xmltv-file") == true) {
          if (thisChannel["x-mapping"] != "-" && thisChannel["x-mapping"] != undefined && thisChannel["x-xmltv-file"] != "-" && thisChannel["x-xmltv-file"] != undefined) {
            switch(thisChannel["x-active"]) {
              case true: tr.className = "activeEPG"; break;
              case false: tr.className = "notActiveEPG"; break;
            }
            break;
          }
        }

    }

    updateChannelLogo(id)

  }

}

function updateChannelLogo(xepg) {
  var thisChannel = xEPG["epgMapping"][xepg];
  if (thisChannel["x-update-channel-icon"] == true) {
    var xXmltvFile  = thisChannel["x-xmltv-file"];
    var xMapping    = thisChannel["x-mapping"];

    if (xXmltvFile != "-" && xXmltvFile.length > 0 && xMapping != "-" && xMapping.length > 0) {
      if (xEPG["xmltvMap"][xXmltvFile][xMapping].hasOwnProperty("icon")) {
        var logoURL = xEPG["xmltvMap"][xXmltvFile][xMapping]["icon"];
        thisChannel["tvg-logo"] = logoURL;
        document.getElementById(xepg).childNodes[2].lastChild.setAttribute("src", logoURL);
        document.getElementById("channel-logo").setAttribute("src", logoURL);
      } else {
        alert("No logo URL in the XMLTV file available")
      }
      
    }
    
    /*
    if (xEPG["xmltvMap"][xXmltvFile][xMapping]["icon"] != undefined) {

      
    }
    */
    
  }
}

function saveChannelLogo(xepg) {
  if (bulk == false) {
    var thisChannel = xEPG["epgMapping"][xepg];
    thisChannel["tvg-logo"] = document.getElementById("tvg-logo").lastChild.value;
    document.getElementById(xepg).childNodes[2].lastChild.setAttribute("src", thisChannel["tvg-logo"]);
    mappingDetail(xepg);
    return
  }

  if (bulk == true) {
    var key   = "tvg-logo";
    var value = document.getElementById("tvg-logo").lastChild.value;
    saveBulk(key, value);

    mappingDetail(xepg);
    return
  }
}

function saveChannelIconUpdate(xepg) {

  var key   = "x-update-channel-icon";
  var value = JSON.parse(document.getElementById("x-update-channel-icon").lastChild.checked);
  if (bulk == false) {
    var thisChannel = xEPG["epgMapping"][xepg];
    thisChannel[key] = value
    updateChannelLogo(xepg)
    
    mappingDetail(xepg);
    searchInMapping();
    return
  }

  if (bulk == true) {
    saveBulk(key, value);
    mappingDetail(xepg);
    return
  }
  
}

function saveChannelName(xepg) {
  if (bulk == false) {
    var thisChannel = xEPG["epgMapping"][xepg];
    thisChannel["x-name"] = document.getElementById("x-name").lastChild.value;
    document.getElementById(xepg).childNodes[3].lastChild.innerHTML = thisChannel["x-name"];
    mappingDetail(xepg);
    searchInMapping();
  }
  
}

function saveChannelNameUpdate(xepg) {
  var key   = "x-update-channel-name";
  var value = JSON.parse(document.getElementById("x-update-channel-name").lastChild.checked);

  if (bulk == false) {
    var thisChannel = xEPG["epgMapping"][xepg];
    thisChannel[key] = value
    mappingDetail(xepg);
    searchInMapping();
    return
  }

  if (bulk == true) {
    saveBulk(key, value);
    mappingDetail(xepg);
    return
  }

}

function saveChannelStatus(xepg) {
  var thisChannel = xEPG["epgMapping"][xepg];
  var xXmltvFile = thisChannel["x-xmltv-file"];

  var key   = "x-active";
  var value = JSON.parse(document.getElementById("x-active").lastChild.checked);

  if (xEPG["xmltvMap"].hasOwnProperty(xXmltvFile) == true) {
    if (thisChannel["x-mapping"] != "-" && thisChannel["x-mapping"] != undefined) {
      thisChannel["x-active"] = !thisChannel["x-active"];
      var tr = document.getElementById(xepg);
      switch(thisChannel["x-active"]) {
        case true: tr.className = "activeEPG"; break;
        case false: tr.className = "notActiveEPG"; break;
      }
      
    } else {
      var err = "XMLTV Channel is not selected"
      alert(err)
      value = false
    }

  } else {
    if (value == true) {
      var err = "XMLTV File is not selecte"
      alert(err)
      value = false
    }
  }

  

  if (bulk == false) {
    var thisChannel = xEPG["epgMapping"][xepg];
    thisChannel[key] = value
    mappingDetail(xepg);
    searchInMapping();

    var tr = document.getElementById(xepg);
    switch(thisChannel["x-active"]) {
      case true: tr.className = "activeEPG"; break;
      case false: tr.className = "notActiveEPG"; break;
    }

    return
  }

  if (bulk == true) {
    saveBulk(key, value);
    mappingDetail(xepg);
    return
  }

}

function saveGroupTitle(xepg) {
  var key   = "x-group-title";
  var value = document.getElementById("x-group-title").lastChild.value;

  if (bulk == false) {
    var thisChannel = xEPG["epgMapping"][xepg];
    document.getElementById(xepg).childNodes[5].lastChild.innerHTML = value;
    thisChannel[key] = value;
    mappingDetail(xepg);
    searchInMapping();
  }

  if (bulk == true) {
    saveBulk(key, value);
    mappingDetail(xepg);
    return
  }

}

function saveCategory(xepg) {
  var key   = "x-category";
  var value = document.getElementById("x-category").lastChild.value;

  if (bulk == false) {
    var thisChannel = xEPG["epgMapping"][xepg];
    thisChannel[key] = value
    document.getElementById(xepg).childNodes[3].lastChild.className = value
    mappingDetail(xepg);
    searchInMapping();
  }

  if (bulk == true) {
    saveBulk(key, value);
    mappingDetail(xepg);
    return
  }

}

function arrangeTable(elm) {
  var tr = elm.parentElement.parentElement;
  var newPosition = elm.value;
  var x_channelID = tr.id;

  switch(isNaN(newPosition)) {
    case true: 
      alert("Ch. No. must be a number");
      mappingError = true;
      break;
  }


  //var item = xEPG["epgMapping"][id];
  var keys = getObjKeys(xEPG["epgMapping"])
  for (var i = 0; i < keys.length; i++) {
    var item = xEPG["epgMapping"][keys[i]];
    if (item["x-epg"] == x_channelID) {

      // Check if position exist
      var oldPosition = item["x-channelID"];

      if (oldPosition != newPosition) {

        console.log(newPosition, newPosition.length);
        if (newPosition.length == 0) {
          mappingError = true
          newPosition = oldPosition;
          
        }

        if (mappingError == true) {
          elm.value = oldPosition;
          return;
        }

        for (var j = keys.length - 1; j >= 0; j--) {
          var channel = xEPG["epgMapping"][keys[j]];
          if (keys[j] != x_channelID) {
            if (newPosition == channel["x-channelID"]) { // If position exist, set next free position.
              newPosition++;
              elm.value = newPosition;
              arrangeTable(elm);
              return;
              /*
              var newError = new Object();
              newError["err"] = "Duplicate ID";
              checkErr(newError);
              sortTable();
              mappingError = true;
              document.getElementById(x_channelID).getElementsByTagName("INPUT")[0].focus();
              return;
              */
            }
          }
        }

      }

      //console.log(oldPosition, newPosition);
      if (keys[i] == x_channelID && oldPosition != newPosition) {  
        item["x-channelID"] = newPosition;
      } 

      document.getElementById("logInfo").className = "notVisible";
      if (columnToSort == 1) {
        sortTable(columnToSort);
      }
      mappingError = false;

    }
  }
}

function changeXmltvFile(xepg, elm) {

  var thisChannel = xEPG["epgMapping"][xepg];
  
  var xXmltvFile    = elm.value;
  var channelID     = thisChannel["tvg-id"];
  thisChannel["x-xmltv-file"] = xXmltvFile;

  if (bulk == false) {

    setTimeout(function(){ 

      var xMapping = "-"

      // Automap
      if (xXmltvFile != "-") {
        if (xEPG["xmltvMap"][xXmltvFile].hasOwnProperty(channelID) == true) {
          thisChannel["x-mapping"] = channelID;
          xMapping = channelID
        } else {
          thisChannel["x-mapping"] = xMapping
        }
      } else {
        thisChannel["x-mapping"] = xMapping

      }
      
      var tr = document.getElementById(xepg);

      if (xMapping == "-") {
        thisChannel["x-active"] = false;
        tr.className = "notActiveEPG"
      } else {
        thisChannel["x-active"]  = true;
        tr.className = "activeEPG"
      }

      // Show data in table
      var td = tr.getElementsByTagName("TD");
      var dataFile = td[td.length - 2].lastChild;
      switch(xXmltvFile) {
        case "-":           dataFile.innerHTML = xXmltvFile; break;
        case "xTeVe Dummy": dataFile.innerHTML = xXmltvFile; break;
        default:            dataFile.innerHTML = getValueFromProviderFile(xXmltvFile, "xmltv", "name"); break;
      }

      //xXmltvFile.replace(/^.*[\\\/]/, '');

      var dataXmltvID = td[td.length - 1].lastChild;
      dataXmltvID.innerHTML = xMapping;
      
      mappingDetail(xepg);

    }, 10);
  }

  if (bulk == true) {
    var key = "x-xmltv-file"
    var value = xXmltvFile
    saveBulk(key, value);

    var key = "x-mapping"
    var value = "-"
    saveBulk(key, value);
    mappingDetail(xepg);
    return
  }

  return
}

function mappingChannel(xepg, elm) {
  var thisChannel = xEPG["epgMapping"][xepg];
  //var xMapping      = elm.value;
  var xMapping      = elm.options[elm.selectedIndex].value
  
  if (bulk == false) {
    
    thisChannel["x-mapping"] = xMapping;

    var tr = document.getElementById(xepg);

    if (xMapping == "-") {
      thisChannel["x-active"] = false;
      tr.className = "notActiveEPG"
    } else {
      thisChannel["x-active"]  = true;
      tr.className = "activeEPG"
    }

    // Show data in table
    var td = tr.getElementsByTagName("TD");
    var dataXmltvID = td[td.length - 1].lastChild;
    dataXmltvID.innerHTML = xMapping;
    //console.log(td[td.length - 1]);
    //console.log(xMapping, elm);

    createSearchObj();
    searchInMapping();
    updateChannelLogo(xepg)
    mappingDetail(xepg);
    return
  }

  if (bulk == true) {

    var key = "x-mapping"
    var value = xMapping
    saveBulk(key, value);

    mappingDetail(xepg);
    return
  }

  return
}


function createNewTD(newItem, elm) {
  var newTD = new Object();
  newTD["_element"] = "TD";
  
  elm.appendChild(createElement(newTD));
  var td = elm.lastChild; 
  
  switch(newItem["_element"]) {
    case "SELECT":
      td.appendChild(createElement(newItem));
      var td = elm.lastChild.lastChild; 
      var values = newItem["_optionValues"];
      for (var i = 0; i < values.length; i++) {
        //console.log(item);
        var newEntry = new Object;
        newEntry["_element"]  = "OPTION";
        newEntry["_text"]     = values[i];
        newEntry["value"]     = values[i];
        td.appendChild(createElement(newEntry));
      }
      td.value = newItem["value"];

      break;
    
    default: 
      
      td.appendChild(createElement(newItem));
      break;
  }
  
}

function saveXEPG() {
  if (mappingError == true) {
    alert("Data could not be saved, errors in the XEPG data.");
    return;
  }
  showLoadingScreen(true);

  var data = new Object();
  data["epgMapping"] = xEPG["epgMapping"];
  data["cmd"] = "saveEpgMapping";
  //console.log(data);
  xTeVe(data);
}

function bulkEdit() {
  bulk = !bulk;
  var className;

  var elm = document.getElementsByClassName("bulk");

  switch(bulk) {
    case true: 
      className = "bulk showBulk";
      break;

    case false: 
      className = "bulk hideBulk";
      bulkEditAll = false;
      break;
  }

  for (var i = 0; i < elm.length; i++) {
    elm[i].className = className;
    elm[i].checked = false;
  }

}

function bulkEditAllChannels() {

  var allTR = document.getElementById("id_mapping").getElementsByTagName("TR");

  for (var i = 1; i < allTR.length; i++) {
    if (allTR[i].style.display != "none") {
      switch(bulkEditAll) {
        case false: allTR[i].firstChild.firstChild.checked = true; break;
        case true: allTR[i].firstChild.firstChild.checked = false; break; 
      }

    }
    
  }

  bulkEditAll = !bulkEditAll;
}

function sortTable(columm) {
  //console.log(columm);
  if (columm == columnToSort) {
    //return;
  }

  var table       = document.getElementById("id_mapping");
  var tableHead   = table.getElementsByTagName("TR")[0];
  var tableItems  = tableHead.getElementsByTagName("TD");
  
  var sortObj = new Object();
  var x, xValue;
  var tableHeader
  var sortByString = false

  if (columm > 0 && columnToSort > 0)  {
    tableItems[columnToSort].className = "pointer";
    tableItems[columm].className = "sortThis";
  }

  columnToSort = columm;

  var rows = table.rows;

  if (rows[1] != undefined) {
    tableHeader = rows[0]

    x = rows[1].getElementsByTagName("TD")[columm];
    
    for (i = 1; i < rows.length; i++) {

      x = rows[i].getElementsByTagName("TD")[columm];

      switch(x.childNodes[0].tagName.toLowerCase()) {
        case "input":
          xValue = x.getElementsByTagName("INPUT")[0].value.toLowerCase();
          break;

        case "p":
          xValue = x.getElementsByTagName("P")[0].innerText.toLowerCase();
          break;
        
        default: console.log(x.childNodes[0].tagName);
      }

      if (xValue == "" || xValue == NaN) {
        xValue = i
        sortObj[i] = rows[i];
      
      } else {

        switch(isNaN(xValue)) {
          case false: 

            xValue = parseFloat(xValue);
            sortObj[xValue] = rows[i]
            break;

          case true:

            sortByString = true
            sortObj[xValue.toLowerCase() + i] = rows[i]
            break;

        }

      }
    
    }

    while (table.firstChild) {
      table.removeChild(table.firstChild);
    }
    
    var sortValues = getObjKeys(sortObj)
    if (sortByString == true) {
      sortValues.sort()
    } else {
      function sortFloat(a, b) { 
        return a - b; 
      }
      sortValues.sort(sortFloat);
    }

    table.appendChild(tableHeader)
    
    for (var i = 0; i < sortValues.length; i++) {
     
      table.appendChild(sortObj[sortValues[i]])

    }
    
  }

}


function sortTable_old(columm) {
  showLoadingScreen(true);
  
  setTimeout(function(){ 

    var table, rows, switching, i, x, y, shouldSwitch;
    table = document.getElementById("id_mapping");

    var tableHead = table.getElementsByTagName("TR")[0];
    var tableItems = tableHead.getElementsByTagName("TD");

    if (columm > 0)  {
      tableItems[columnToSort].className = "pointer";
      tableItems[columm].className = "sortThis";
    }
    
    columnToSort = columm;

    /*
    for (var i = 0; i < tableItems.length; i++) {
      if (tableItems[i].className != undefined) {
        tableItems[i].className = "pointer"
      }

    }
    */

    

    console.log(tableItems); 

    switching = true;
    while (switching) {
      switching = false;
      rows = table.rows;
      for (i = 1; i < (rows.length - 1); i++) {
        shouldSwitch = false;

        x = rows[i].getElementsByTagName("TD")[columm];
        y = rows[i + 1].getElementsByTagName("TD")[columm];

        switch(x.childNodes[0].tagName.toLowerCase()) {
          case "input":
            xValue = x.getElementsByTagName("INPUT")[0].value.toLowerCase();
            yValue = y.getElementsByTagName("INPUT")[0].value.toLowerCase();
            break;

          case "p":
            xValue = x.getElementsByTagName("P")[0].innerText.toLowerCase();
            yValue = y.getElementsByTagName("P")[0].innerText.toLowerCase();
            break;
          
          default: console.log(x.childNodes[0].tagName);
        }

        
        switch(isNaN(xValue)) {
          case false: xValue = parseFloat(xValue) ; break;
        }

        switch(isNaN(yValue)) {
          case false: yValue = parseFloat(yValue) ; break;
        }
        

        if (xValue > yValue) {
          shouldSwitch = true;
          break;
        }

      }
      if (shouldSwitch) {
        rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
        switching = true;
      }
    }
    createSearchObj()
    
    showLoadingScreen(false);
  }, 20);

}

function showXEPG() {
  var url = location.protocol + "//" + location.hostname + ":" + location.port + "/xmltv/xteve.xml"
  var win = window.open(url, '_blank');
  win.focus();
}