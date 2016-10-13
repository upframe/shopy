'use strict';

// Find a better way to do this
if (document.querySelector("main.covered")) {
    document.querySelector("body").classList.add("covered");
}

document.addEventListener("DOMContentLoaded", () => {
    let thing;

    // if register form was loaded it means everything is fine
    // if not, it means register is only available by invitation
    if (thing = document.getElementById("register")) {
        thing.addEventListener("submit", registerHandler);
    }

    if (thing = document.getElementById("login")) {
        thing.addEventListener("submit", loginHandler);
    }

    if (thing = document.getElementById("settings")) {
        thing.addEventListener('submit', submitSettings);
    }

    if (thing = document.querySelector('body>nav input[type="checkbox"]')) {
        document.addEventListener('click', function(event) {
            if (event.target.id != "dropdown" && event.target.parentElement.getAttribute("for") != "dropdown") {
                document.querySelector('body>nav input[type="checkbox"]').checked = false;
            }
        });
    }

    if (thing = document.getElementById("store")) {
        initializeStore();
    }

    if (thing = document.getElementById("cart")) {
        initializeCart();
    }

    if ((thing = document.getElementById("deactivate")) && window.location.pathname == "/settings") {
        thing.addEventListener("click", deactivateAccount);
    }

    if (thing = document.getElementById("copy-ref")) {
        thing.addEventListener("click", copyReferral);
    }

    if (thing = document.getElementById("reset")) {
        thing.addEventListener("submit", resetEmailForm);
    }

    if (thing = document.getElementById("reset-form")) {
        thing.addEventListener("submit", resetForm);
    }

    if (window.location.pathname === "/checkout/discounts") {
        document.getElementById("promocode").addEventListener("keyup", validateCoupon)
    }
});

function copyReferral(event) {
    event.preventDefault();

    var input = document.createElement('textarea');
    document.body.appendChild(input);
    input.value = this.dataset.copy;
    input.focus();
    input.select();
    document.execCommand('Copy');
    input.remove();
}

function validateCoupon(e) {
    let el = document.getElementById("promocode");

    if (this.value == "") {
        el.classList.remove("error");
        el.classList.remove("works");
        return;
    }

    let request = new XMLHttpRequest();
    request.open("POST", "/coupon/validate", true);
    request.send(this.value);
    request.onreadystatechange = function() {
        if (request.readyState == 4) {
            switch (request.status) {
                case 200:
                    //useCoupon(request.responseText);
                    el.classList.add("works");
                    el.classList.remove("error");
                    break;
                case 404:
                    el.classList.remove("works");
                    el.classList.add("error");
                    break;
            }
        }
    }
}
/*
function useCoupon(discount) {
    document.getElementById("price").innerHTML = "Total: € " + (document.getElementById("price").innerHTML.substr(7) * (1 - (discount / 100)));
} */

function initializeStore() {
    let request = function() {
        cartRequest(
            "POST",
            window.location.origin + "/cart/" + this.parentElement.dataset.id,
            ""
        );
    }

    Array.from(document.querySelectorAll(".btnBuy")).forEach((btn) => {
        btn.addEventListener("click", request);
    });
}

function initializeCart() {
    let request = function() {
        cartRequest(
            "DELETE",
            window.location.origin + "/cart/" + this.parentElement.parentElement.dataset.id,
            "",
            this.parentElement.parentElement.dataset.id
        );
    }

    Array.from(document.querySelectorAll(".btnRemove")).forEach((btn) => {
        btn.addEventListener("click", request);
    });
}

function cartRequest(method, link, data, itemID) {
    let request = new XMLHttpRequest();
    request.open(method, link, true);
    request.send(data);
    request.onreadystatechange = function() {
        if (request.readyState == 4) {
            switch (request.status) {
                case 200:
                    if (method == "DELETE") {
                        let item = document.querySelector('tr[data-id="' + itemID + '"]');
                        item.children[0].innerHTML--;

                        let price = item.querySelector('td[data-name="Price"]');
                        price.innerHTML -= item.dataset.price;

                        let total = document.querySelector("#total");
                        total.innerHTML -= item.dataset.price;

                        if (item.children[0].innerHTML == 0) {
                            item.parentElement.removeChild(item);
                        }
                    }
                    break;
                default:
                    console.log(request.status + ": Bad request");
            }
        }
    }
}

function deactivateAccount(event) {
    event.preventDefault();
    if (confirm("Are you sure you want to deactivate your account?")) {
        let request = new XMLHttpRequest();
        request.open("POST", "/settings/deactivate", true);
        request.send();
        request.onreadystatechange = function() {
            if (request.readyState == 4) {
                switch (request.status) {
                    case 200:
                        alert("Check your email please.");
                        break;
                    default:
                        alert("Something wrong happened.");
                }
            }
        }
    }
}

function submitSettings(event) {
    event.preventDefault();
    let inputs = this.querySelectorAll('input'),
        form = new Object();


    var request = new XMLHttpRequest();
    request.open("PUT", window.location, true);
    request.setRequestHeader("Content-type", "application/json; charset=utf-8");
    request.send(JSON.stringify(copyFormToObject(this)));
    request.onreadystatechange = function() {
        if (request.readyState == 4) {
            switch (request.status) {
                case 200:
                    formError("Successfully updated.", "success");
                    break;
                default:
                    formError("Something wrong happened :(", "success");
            }
        }
    }
}

function printMessage(status, hash) {
    let type = (status >= 200 && status < 300) ? "success" : "error";

    if (status == 424) {
        type = "warning";
    }

    if (status in hash) {
        formError(hash[status], type);
    } else {
        formError(hash['default'], type);
    }
}

var registerMessages = {
    200: "You're now registered. Check your email to confirm.",
    201: "You're now registered. Check your email to confirm.",
    400: "Some fields are empty or invalid.",
    403: "The reffer link is invalid.",
    409: "Your email is already registered. Please <a href='/login'>login</a>.",
    410: "It seems that in the meanwhile the person that invited you ran out of invites.",
    'default': "Something went wrong and we are unable to explain it right now."
}

function registerHandler(event) {
    event.preventDefault();

    if (!checkRegisterFields(this)) {
        formError("Passwords doesn't match or some fields are empty.", "error");
        return;
    }

    var form = copyFormToObject(this),
        hash = new jsSHA("SHA-256", "TEXT"),
        request = new XMLHttpRequest();

    hash.update(form.password);
    form.password = hash.getHash("HEX");

    request.open("POST", window.location, true);
    request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    request.send(form.serialize());
    request.onreadystatechange = function() {
        if (request.readyState == 4) {
            printMessage(request.status, registerMessages)
        }
    }
}

var loginMessages = {
    400: "You might have left some fields blank!",
    401: "Incorrect password.",
    404: "We can't find you in our database. <a href='/register'>Register</a> first.",
    423: "Your account is deactivated.",
    424: "Check your email to confirm your account first. <a href='#' onclick='resendConfirmation();'>Resend confirmation.</a>",
    'default': "Something went wrong and we are unable to explain it right now."
}

function loginHandler(event) {
    event.preventDefault();

    let form = copyFormToObject(this),
        hash = new jsSHA("SHA-256", "TEXT"),
        request = new XMLHttpRequest();

    hash.update(form.password);
    form.password = hash.getHash("HEX");

    request.open("POST", window.location, true);
    request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    request.send(form.serialize());
    request.onreadystatechange = function() {
        if (request.readyState == 4) {
            if (request.status == 200) {
                window.location = window.location.protocol + "//" + window.location.hostname;
                return;
            }

            printMessage(request.status, loginMessages)
        }
    }
}

var resetEmailStatus = {
    200: 'Please, check your inbox to continue.',
    'default': "Something went wrong and we are unable to explain it right now."
}

function resetEmailForm(event) {
    event.preventDefault();

    let form = copyFormToObject(this),
        request = new XMLHttpRequest();

    request.open("POST", window.location, true);
    request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    request.send(form.serialize());
    request.onreadystatechange = function() {
        if (request.readyState == 4) {
            printMessage(request.status, resetEmailStatus)
        }
    }
}

function resetForm(event) {
    event.preventDefault();

    let form = copyFormToObject(this),
        hash = new jsSHA("SHA-256", "TEXT"),
        request = new XMLHttpRequest();

    if (form.password != form.confirmpassword) {
        return formError("The passwords don't match.", "error")
    }

    hash.update(form.password);
    form.password = hash.getHash("HEX");

    request.open("POST", window.location, true);
    request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    request.send(form.serialize());
    request.onreadystatechange = function() {
        if (request.readyState == 4) {
            if (request.status == 200) {
                window.location = window.location.protocol + "//" + window.location.hostname + "/login";
                return;
            }

            printMessage(request.status, loginMessages)
        }
    }
}

var resendMessages = {
    200: "Check your email!",
    201: "Check your email!",
    404: "We can't find you in our database. <a href='/register'>Register</a> first.",
    'default': "Something went wrong and we are unable to explain it right now."
}

function resendConfirmation() {
    email = document.querySelector('input[name="email"]');

    if (email.value.search("@") == -1) {
        return formError("Your email is invalid.", "error");
    }

    var request = new XMLHttpRequest();
    request.open("POST", window.location, true);
    request.setRequestHeader("Resend", "true");
    request.setRequestHeader("Email", email.value);
    request.send(form.serialize());
    request.onreadystatechange = function() {
        if (request.readyState == 4) {
            printMessage(request.status, resendMessages)
        }
    }
}

function checkRegisterFields(form) {
    let inputs = form.querySelectorAll("input");

    for (let i = 0; i < inputs.length; i++) {
        if (inputs[i].value == "") {
            console.log(input[i])
            return false;
        }
    }

    if (form.querySelector('input[name="password"]').value != form.querySelector('input[name="password_conf"]').value) {
        return false;
    }

    if (form.querySelector('input[name="email"]').value.search("@") == -1) {
        return false;
    }

    return true;
}