// Register page JS file

document.addEventListener("DOMContentLoaded", () => {
  handlePage();
});


function handlePage() {
  switch(window.location.pathname) {
    case "/register":
      handleRegister();
      break;

  }
}

function handleRegister() {
  // if register form was loaded it means everything is fine
  // if not, it means register is only available by invitation
  if(form = document.getElementById("regForm")) {
    form.querySelectorAll("button[type=submit]")[0].addEventListener("click", function(e) {
      e.preventDefault();
      if((form.querySelectorAll("input[type=password]"[0].value))== (form.querySelectorAll("input[type=password]")[1].value)) {
        // passwords match
        firstname = form.querySelectorAll('input[name=first_name]')[0].value;
        lastname = form.querySelectorAll('input[name=last_name]')[0].value;
        email = form.querySelectorAll("input[name=email]")[0].value;
        password = form.querySelectorAll("input[name=password]")[0].value;
        var pwdHash = new jsSHA("SHA-256", "TEXT");
        pwdHash.update(password);
        pwdHash = pwdHash.getHash("HEX");

        var request = new XMLHttpRequest();
        request.open("POST", window.location, true);
        request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        request.send("first_name=" + firstname + "&last_name=" + lastname + "&email=" + email + "&password=" + pwdHash);
        request.onreadystatechange = function() {
          if(request.readyState == 4) {
            if(request.status == 200) {// success
              console.log(request.responseText);
            } else {
              alert(request.status + ":" + request.responseText);
            }
          }
        }
      } else {
        // passwords don't match
      }
    });
  }
}
