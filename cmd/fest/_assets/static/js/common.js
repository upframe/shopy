'use strict';

Object.prototype.serialize = function() {
    var str = [];
    for (var p in this) {
        if (this.hasOwnProperty(p)) {
            str.push(encodeURIComponent(p) + "=" + encodeURIComponent(this[p]));
        }
    }

    return str.join("&");
}

function formError(message, type) {
    let error = document.getElementById("form-error");

    error.classList.remove("warning");
    error.classList.remove("success");
    error.classList.remove("error");

    error.classList.add(type);
    error.innerHTML = message;
    error.classList.add("shake");

    setTimeout(() => {
        error.classList.remove("shake");
    }, 830);
}

function copyFormToObject(form) {
    let object = new Object();
    object.ID = 0;

    let inputs = form.querySelectorAll('input, textarea');

    Array.from(inputs).forEach((input) => {
        let name = input.name;

        if (typeof name == 'undefined' || name == null || name == "") {
            return;
        }

        if (name.indexOf(".") !== -1) {
            let parts = name.split(".");

            if (parts.length > 2) {
                console.log("Invalid.");
                return;
            }

            let val = inputToValue(input);

            if (val == null) {
                object[parts[0]] = null;
            } else {
                object[parts[0]] = new Object();
                object[parts[0]][parts[1]] = inputToValue(input);
            }
        } else {
            object[name] = inputToValue(input);
        }
    })

    if (isNaN(object.ID)) object.ID = 0;
    return object;
}

function inputToValue(input) {
    switch (input.type) {
        case "number":
            if (input.value == "") {
                return null;
            }

            return parseInt(input.value);
        case "datetime-local":
            return (new Date(input.value)).toISOString();
            break;
        case "checkbox":
            return input.checked;
            break;
        default:
            return input.value;
    }
}