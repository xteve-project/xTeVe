function createFirstAccount(elm) {
  var err = false;
  var div = document.getElementById(elm);
  console.log(div);

  var form = document.getElementById('authentication');
  
  const username  = document.getElementById('username');
  const password  = document.getElementById('password');
  const confirm   = document.getElementById('confirm');

  var inputs = div.getElementsByTagName('INPUT')
  console.log(confirm);

  switch(confirm) {
    case null: break;
    
    default: 
      for (var i = 0; i < inputs.length; i++) {
        if (inputs[i].value.length == 0) {
          inputs[i].style.borderColor = 'red';
          err = true
        }
      }

      switch(err) {
        case true: return; break;
        case false: 
          if (password.value != confirm.value) {
            confirm.style.borderColor = 'red';
            return;
          }
          break;
      }
  }


  

  form.submit();
  return;
}