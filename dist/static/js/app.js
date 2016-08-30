// Register page JS file

document.addEventListener("DOMContentLoaded", () => {
    switch (window.location.pathname) {
        case "/register":
            // if register form was loaded it means everything is fine
            // if not, it means register is only available by invitation
            if (form = document.getElementById("register")) {
                form.addEventListener("submit", registerHandler);
            }
            break;
    }
});

Object.prototype.serialize = function() {
    var str = [];
    for (var p in this) {
            if (this.hasOwnProperty(p)) {
                str.push(encodeURIComponent(p) + "=" + encodeURIComponent(this[p]));
            }
    }

    return str.join("&");
}

var registerHandler = function(event) {
    event.preventDefault();

    if (checkRegisterFields(this)) {
        // passwords match. so que let fica condicionado ao if{} e var fica na funcao toda
        var form = new Object();
        // ugly names!!!
        form.first_name = this.querySelectorAll('input[name=first_name]')[0].value,
            form.last_name = this.querySelectorAll('input[name=last_name]')[0].value,
            form.email = this.querySelectorAll("input[name=email]")[0].value,
            form.password = this.querySelectorAll("input[name=password]")[0].value;

        let pwdHash = new jsSHA("SHA-256", "TEXT");
        pwdHash.update(form.password);
        form.password = pwdHash.getHash("HEX");

        var request = new XMLHttpRequest();
        request.open("POST", window.location, true);
        request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        request.send(form.serialize());
        request.onreadystatechange = function() {
            if (request.readyState == 4) {
                switch (request.status) {
                    case 200:
                    case 201:
                        alert("Registered. Check your email, bitch.");
                        break;
                    case 400:
                        alert("Bad request");
                        break;
                    case 403:
                        alert("Forbidden");
                        break;
                    case 409:
                        alert("Conflict");
                        break;
                    case 410:
                        alert("Gone");
                        break;
                    default:
                        alert("Something went wrong.")
                }
            }
        }
    } else {
        alert("passwords don't match or fields are empty")
    }
}

function checkRegisterFields(form) {
    // check if all fields are empty
    for (var x = 0; x < form.children.length - 1; x++) {
        if (form[x].value == "") {
            return false;
        }
    }

    if (form[3].value != form[4].value || form[2].value.search("@") == -1) {
        return false;
    }

    return true;
}
