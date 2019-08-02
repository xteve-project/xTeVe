function login() {
    var err = false;
    var data = new Object();
    var div = document.getElementById("content");
    var form = document.getElementById("authentication");
    var inputs = div.getElementsByTagName("INPUT");
    console.log(inputs);
    for (var i = inputs.length - 1; i >= 0; i--) {
        var key = inputs[i].name;
        var value = inputs[i].value;
        if (value.length == 0) {
            inputs[i].style.borderColor = "red";
            err = true;
        }
        data[key] = value;
    }
    if (err == true) {
        data = new Object();
        return;
    }
    if (data.hasOwnProperty("confirm")) {
        if (data["confirm"] != data["password"]) {
            alert("sdafsd");
            document.getElementById('password').style.borderColor = "red";
            document.getElementById('confirm').style.borderColor = "red";
            document.getElementById("err").innerHTML = "{{.account.failed}}";
            return;
        }
    }
    console.log(data);
    form.submit();
}
