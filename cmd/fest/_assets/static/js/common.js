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

        switch (input.type) {
            case "number":
                object[name] = parseInt(input.value);
                break;
            case "datetime-local":
                object[name] = (new Date(input.value)).toISOString();
                break;
            case "checkbox":
                object[name] = input.checked;
                break;
            default:
                object[name] = input.value;
        }
    })

    if (isNaN(object.ID)) object.ID = 0;
    return object;
}
