function login() {
  var err:Boolean = false
  var data = new Object()
  var div:any = document.getElementById("content")
  var form:any = document.getElementById("authentication")

  var inputs:any = div.getElementsByTagName("INPUT")

  for (var i = inputs.length - 1; i >= 0; i--) {
    
    var key:string = (inputs[i] as HTMLInputElement).name
    var value:string = (inputs[i] as HTMLInputElement).value

    if (value.length == 0) {
      inputs[i].style.borderColor = "red"
      err = true
    }

    data[key] = value

  }

  if (err == true) {
    data = new Object()
    return
  }

  if (data.hasOwnProperty("confirm")) {

    if (data["confirm"] != data["password"]) {
      alert("sdafsd")
      document.getElementById('password').style.borderColor = "red"
      document.getElementById('confirm').style.borderColor = "red"

      document.getElementById("err").innerHTML = "{{.account.failed}}"
      return
    }

  }

  form.submit();

}