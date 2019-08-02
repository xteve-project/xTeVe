function openUsers(elm) {
  colomnSort = 0;

  var newDiv = document.getElementById("settings");
  
  var newEntry = new Object();
  newEntry["_element"]  = "HR";
  newDiv.appendChild(createElement(newEntry));

  var newEntry = new Object();
  newEntry["_element"]  = "INPUT";
  newEntry["type"] = "button";
  newEntry["class"] = "button";
  newEntry["value"] = "New";
  newEntry["onclick"] = "userDetail(0)";
  newDiv.appendChild(createElement(newEntry));

  var div = document.getElementById("settings");

   // Build table
  var newTable = new Object();
  newTable["_element"]  = "TABLE";
  newTable["id"]        = "id_mapping";
  newTable["class"]     = "table-mapping";
  div.appendChild(createElement(newTable));

  setTimeout(function(){ 
    createUsersTable(); 
  }, 10);
}

function createUsersTable() {
  var table = document.getElementById("id_mapping");
  table.innerHTML = "";
  var newTR = new Object();
  newTR["_element"] = "TR";
  newTR["class"]    = "table-mapping-header";
  table.appendChild(createElement(newTR));

  var tr = table.lastChild;
  var trHeadlines = new Array("Username", "Password", "WEB", "PMS", "M3U", "XML", "API")

  for (var i = 0; i < trHeadlines.length; i++) {
    var newTD = new Object();
    newTD["_element"] = "TD";
    newTD["_text"]    = trHeadlines[i];
    tr.appendChild(createElement(newTD));
  }


  // Sort users
  var userIds = getObjKeys(users);

  var userObj = new Object();

  for (var i = 0; i < userIds.length; i++) {
    var username = users[userIds[i]]["data"]["username"];
    userObj[username] = userIds[i];
  }

  var allUsers = getObjKeys(userObj);
  allUsers.sort();
  // --

  for (var i = 0; i < allUsers.length; i++) {
    var table     = document.getElementById("id_mapping");
    var userID    = userObj[allUsers[i]];
    var username  = allUsers[i];
    var item      = users[userID]["data"];

    // Create TR
    var newTR = new Object();
    newTR["_element"]       = "TR";
    newTR["class"]          = "";
    newTR["id"]             = userID;
    newTR["onclick"]        = 'javascript: userDetail("' + userID + '");';
    table.appendChild(createElement(newTR));

    var tr = table.lastChild;

    // Create username TD
    var newTD = new Object();
    newTD["_element"] = "P";
    newTD["_text"]    = username;
    createNewTD(newTD, tr);

    // Create password TD
    var newTD = new Object();
    newTD["_element"] = "P";
    newTD["_text"]    = ".....";
    createNewTD(newTD, tr);

    // Create web access
    var newTD = new Object();
    newTD["_element"] = "P";
    switch(item["authentication.web"]){
      case true: newTD["_text"]    = "✓"; break;
      default:   newTD["_text"]    = "-"; break;
    }
    createNewTD(newTD, tr);

    // Create PMS access
    var newTD = new Object();
    newTD["_element"] = "P";
    switch(item["authentication.pms"]){
      case true: newTD["_text"]    = "✓"; break;
      default:   newTD["_text"]    = "-"; break;
    }
    createNewTD(newTD, tr);

    // Create M3U access
    var newTD = new Object();
    newTD["_element"] = "P";
    switch(item["authentication.m3u"]){
      case true: newTD["_text"]    = "✓"; break;
      default:   newTD["_text"]    = "-"; break;
    }
    createNewTD(newTD, tr);

    // Create XMLTV access
    var newTD = new Object();
    newTD["_element"] = "P";
    switch(item["authentication.xml"]){
      case true: newTD["_text"]    = "✓"; break;
      default:   newTD["_text"]    = "-"; break;
    }
    createNewTD(newTD, tr);

    // Create API access
    var newTD = new Object();
    newTD["_element"] = "P";

    switch(item["authentication.api"]){
      case true: newTD["_text"]    = "✓"; break;
      default:   newTD["_text"]    = "-"; break;
    }
    createNewTD(newTD, tr);

  }

  // usage Info  
  var div = document.getElementById("settings");
  switch(menu[activeMenu.id].hasOwnProperty("_usage")) {
    case true: 
      var usageItem = new Object();
      usageItem["_element"] = "PRE"
      usageItem["_text"]    = menu[activeMenu.id]["_usage"];

      var newHR = new Object();
      newHR["_element"] = "HR"
      div.appendChild(createElement(newHR));
      div.appendChild(createElement(usageItem));
  }

  sortTable(0);
}

function userDetail(userID) {
  showPopUpElement('user-detail');
  setTimeout(function(){ 
    showElement("popup", true);
  }, 10);
  var defaultUser;

  document.getElementById("saveUserDetail").setAttribute("onclick", 'javascript: saveUserDetail("' + userID + '", false)');
  document.getElementById("deleteUserDetail").setAttribute("onclick", 'javascript: saveUserDetail("' + userID + '", true)');

  var data = new Object();
  
  switch(userID) {
    case 0:   // New User
      data["username"] = "";
      data["authentication.web"] = false;
      data["authentication.pms"] = true;
      data["authentication.xml"] = true;
      data["authentication.m3u"] = false;
      data["authentication.api"] = false;
      data["defaultUser"]        = false;
      setTimeout(function(){ 
        showElement("deleteUserDetail", false)
      }, 1);

      break; 
    
    default: 
      data = users[userID]["data"]; 
      showElement("deleteUserDetail", true)   
      document.getElementById("deleteUserDetail").className = "delete";
   
      break
  }
  

  var username = data["username"];
  data["password"] = "";
  data["confirm"] = "";

  var keys = getObjKeys(data);
  defaultUser = data["defaultUser"];
  if (data.hasOwnProperty("defaultUser")) {
    defaultUser = JSON.parse(data["defaultUser"]);
  }

  for (var i = 0; i < keys.length; i++) {

    if(document.getElementById(keys[i])){
      var td = document.getElementById(keys[i])
    } else {
      var td = undefined;
    }

    var newItem = new Object();

    newItem["_element"] = "INPUT";
    
    newItem["value"]    = data[keys[i]];
    newItem["name"]     = keys[i];





    switch(keys[i].indexOf("authentication")) {
      case -1: 
        if (keys[i] == "password" || keys[i] == "confirm") {
          newItem["type"]     = "password";
        } else {
          newItem["type"]     = "text";
        }
        break;

      default: 
        newItem["type"]     = "checkbox";
        
        if (keys[i] == "authentication.web" && defaultUser == true) {
          newItem["onclick"] = "return false"; 
        }
        
        if (data[keys[i]] == true) {          
          newItem["checked"]  = data[keys[i]];  
        }
        
        break;
    }

    switch(keys[i]) {
      case "defaultUser": 
        //if (data[keys[i]] == true) {
        newItem["type"]     = "hidden";
        //}
    }


    if (td != undefined) {
      td.innerHTML = "";
      var element = createNewElement(newItem)
      //console.log(element);
      td.appendChild(element);
    }


  }


  if (defaultUser == true) {
    showElement("deleteUserDetail", false)
  } else {
    showElement("deleteUserDetail", true)
    document.getElementById("deleteUserDetail").className = "delete";
  }

}

function saveUserDetail(userID, deleteUser) {

  var inputs = document.getElementById("user-detail-table").getElementsByTagName("INPUT");

  var newUserData = new Object();
  for (var i = 0; i < inputs.length; i++) {
    switch(inputs[i].type) {
      case "checkbox":  newUserData[inputs[i].name] = inputs[i].checked; break;
      default:          newUserData[inputs[i].name] = inputs[i].value; break;
    }
    
    if (inputs["username"].value.length == 0) {
      inputs["username"].style.border = "solid 1px red";
      return;
    }

    switch(userID) {
      case "0": 
        if (inputs["password"].value.length == 0) {
          console.log(inputs["password"].value.length);
          inputs["password"].style.border = "solid 1px red";
          return
        }
        break;
    }

    if (inputs["password"].value.length > 0) {
      if (inputs["password"].value != inputs["confirm"].value) {
        inputs["password"].style.border = "solid 1px red";
        inputs["confirm"].style.border = "solid 1px red";
        return;
      }
    }
    
  }

  var data = new Object();
  
  switch(userID) {
    case "0":
      //data = newUserData
      data["userData"]  = newUserData
      data["cmd"]       = "saveNewUser"; break;
    
    default:  
      var thisUser      = new Object();
      
      if (deleteUser == true) {
        if (confirm('Delete the selected user?')) {
          data["deleteUser"] = true;
        } else {
          showElement("popup", false);
          return
        }
      }
      
      thisUser[userID]  = newUserData;

      data["userData"] = thisUser;
      data["cmd"] = "saveUserData"; break;
  }
  
  xTeVe(data);
  //createUsersTable()
  showElement("popup", false);
}


