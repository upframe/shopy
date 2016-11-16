'use strict';

var editor,
    monthNames = [
        "Jan", "Feb", "Mar",
        "Apr", "May", "Jun", "Jul",
        "Aug", "Sep", "Oct",
        "Nov", "Dec"
    ];

document.addEventListener("DOMContentLoaded", () => {
    // Initializes the Single Form variable
    editor = document.getElementById("editor");

    let thing;

    if (editor) {
        // Add an event listener to the Single Form
        document.addEventListener("click", pageClick)
        editor.addEventListener('submit', submitHandler);
        editor.addEventListener('keyup', escapeHandler)

        // Get all the edit buttons and initialize them
        var btns = document.getElementsByClassName("edit");
        Array.from(btns).forEach(editHandler);

        // Initialize add button
        if (thing = document.getElementById("add")) {
            thing.addEventListener("click", newHandler);
        }
        if (thing = document.getElementById("deactivate")) {
            thing.addEventListener("click", deactivateHandler);
        }
        if (thing = document.getElementById("edit")) {
            thing.addEventListener("click", editMultipleHandler);
        }
        if (thing = document.getElementById("activate")) {
            thing.addEventListener("click", activateHandler);
        }

        document.getElementById('expand').addEventListener('click', function(event) {
            event.preventDefault();
            editor.classList.toggle("show");
        });

        let rows = document.querySelectorAll('tbody tr');
        Array.from(rows).forEach((row) => {
            row.addEventListener('click', function(event) {
                if (event.target.innerHTML == 'mode_edit') return;
                this.classList.toggle('highlight');
                refreshButtons();
            });
        });
    }

    highlight();
    refreshButtons();
});

function highlight() {
    let hash = window.location.hash.replace('#', '');

    if (hash === '') return;

    let items = hash.split(',')

    document.querySelector('tr[data-id="' + items[0] + '"]').scrollIntoView();
    for (var i = 0; i < items.length; i++) {
        let row = document.querySelector('tr[data-id="' + items[i] + '"]');
        if (typeof row == 'undefined' || row == null) continue;
        row.classList.add('highlight');
    }
}

function refreshButtons() {
    let selected = document.querySelectorAll('tr.highlight'),
        deactivate = document.getElementById("deactivate"),
        edit = document.getElementById("edit"),
        activate = document.getElementById("activate");

    if (selected.length == 0) {
        if (deactivate) deactivate.setAttribute("disabled", "true");
        if (edit) edit.setAttribute("disabled", "true");
        if (activate) activate.setAttribute("disabled", "true");
        return;
    }

    if (deactivate) deactivate.removeAttribute("disabled");
    if (edit) edit.removeAttribute("disabled");
    if (activate) activate.removeAttribute("disabled");
}

function pageClick(event) {
    for (let i = 0; i < event.path.length; i++) {
        if (event.path[i].id == "add" ||
            event.path[i].id == "editor" ||
            event.path[i].id == "delete" ||
            event.path[i].id == "edit") {
            return true;
        }

        if (event.path[i].className == "edit") {
            return true;
        }
    }

    if (editor.classList.contains("show")) {
        editor.classList.remove("show");
    }
}

function escapeHandler(event) {
    if (event.key == "Escape") {
        document.getElementById("editor").classList.remove("show");
    }
}

function newHandler(event) {
    if (!editor.classList.contains("show")) {
        editor.classList.add('show');
    }

    clearForm(editor);
    editor.children[1].children[3].focus();
}

function editHandler(btn) {
    btn.addEventListener("click", function(e) {
        e.preventDefault();
        editor.querySelector('#edit-text').style.display = "block";
        editor.querySelector('#new-text').style.display = "none";
        editor.classList.add("show");
        copyRowToForm(btn.parentElement.parentElement);
        editor.children[1].children[3].focus();
    });
}

function editMultipleHandler(event) {
    event.preventDefault();

    var div = editor.children[1];

    editor.querySelector('#edit-text').style.display = "block";
    editor.querySelector('#new-text').style.display = "none";
    document.getElementById('barID').innerHTML = "multiple";

    for (var x = 0; x < div.childElementCount; x++) {
        let type = div.children[x].type;

        if (typeof type == "undefined") {
            continue;
        }

        div.children[x].value = "";

        switch (type) {
            case "checkbox":
                div.children[x].dataset.initial = false;
                div.children[x].checked = false;
                break;
            default:
                div.children[x].placeholder = "(multiple values)";
                break;
        }
    }

    editor.classList.add("show");
}

function deactivateHandler(event) {
    event.preventDefault();

    Array.from(document.querySelectorAll('tr.highlight')).forEach((row) => {
        let link = "/admin/" + window.location.pathname.split("/")[2] + "/" + row.dataset.id
        let request = new XMLHttpRequest();

        request.open("DELETE", link);
        request.onreadystatechange = function() {
            if (request.readyState == 4) {
                if (request.status == 200) {
                    row.querySelector('td[data-name="Deactivated"] input[type="checkbox"]').checked = true;
                    row.classList.remove("highlight");
                    refreshButtons();
                }
            }
        }
        request.send();
    });
}

function activateHandler(event) {
    event.preventDefault();

    Array.from(document.querySelectorAll('tr.highlight')).forEach((row) => {
        let link = "/admin/" + window.location.pathname.split("/")[2] + "/" + row.dataset.id

        copyRowToForm(row);
        let data = copyFormToObject(editor);
        data["Deactivated"] = false;

        let request = new XMLHttpRequest();
        request.open("PUT", link);
        request.send(JSON.stringify(data));
        request.onreadystatechange = function() {
            if (request.readyState == 4) {
                if (request.status == 200) {
                    row.querySelector('td[data-name="Deactivated"] input[type="checkbox"]').checked = false;
                    row.classList.remove("highlight");
                    refreshButtons();
                }
            }
        }
    });
}

function submitHandler(event) {
    event.preventDefault();

    let data = copyFormToObject(this);
    let method = (data.ID == 0) ? 'POST' : 'PUT';
    let link = this.dataset.link;

    if (data.ID != 0) {
        link += "/" + data.ID;
    }

    let request = new XMLHttpRequest();
    request.open(method, link);
    request.send(JSON.stringify(data));
    request.onreadystatechange = function() {
        if (request.readyState == 4) {
            if (request.status == 200) {
                editor.className = "Down";

                if (method == "PUT") {
                    copyFormToRow(document.querySelector('tr[data-id="' + data.ID + '"]'));
                } else {
                    window.location.pathname = "/admin/" + window.location.pathname.split("/")[2] + "/" + request.responseText;
                }
            } else {
                formError("Something went wrong.", "error")
            }
        }
    }
}

function copyFormToRow(row) {
    let form = editor;

    if (typeof form == 'undefined') {
        return;
    }

    let inputs = form.querySelectorAll('input');

    Array.from(inputs).forEach((input) => {
        let name = input.name;

        if (typeof name == 'undefined' || name == null || name == "") {
            return;
        }

        let space = row.querySelector('td[data-name="' + name + '"]');

        if (typeof space == 'undefined' || space == null) {
            return;
        }

        switch (input.type) {
            case "datetime-local":
                space.innerHTML = getPrettyDate(new Date(input.value));
                break;
            case "checkbox":
                space.querySelector('input[type="checkbox"]').checked = input.checked;
                break;
            default:
                space.innerHTML = input.value;
        }
    });
}

// getPrettyDate puts the date in a 22 Oct 99 08:48 UTC pretty format
function getPrettyDate(date) {
    let pretty;

    let normalize = function(number) {
        if (number.length == 1) {
            return "0" + number;
        }

        return number;
    }

    pretty = normalize(date.getUTCDate().toString()) + " "
    pretty += monthNames[date.getUTCMonth()] + " "
    pretty += date.getUTCFullYear().toString().substr(2, 2) + " "
    pretty += normalize(date.getUTCHours().toString()) + ":"
    pretty += normalize(date.getUTCMinutes().toString()) + " UTC"

    return pretty;
}

// Copies the information from a row to the editor form
function copyRowToForm(row) {
    let form = editor;

    if (typeof form == 'undefined') {
        return;
    }

    for (var x = 0; x < row.childElementCount - 1; x++) {
        let data = row.children[x].dataset.name,
            input = form.querySelector("input[name=" + data + "]"),
            barID = form.querySelector("#barID");

        if (data == undefined || input == undefined) {
            continue;
        }

        if (data == "ID") {
            barID.innerHTML = row.children[x].innerHTML;
        }

        switch (input.type) {
            case "datetime-local":
                input.value = new Date(row.children[x].innerHTML).toISOString().substr(0, 16);
                break;
            case "checkbox":
                input.checked = row.children[x].querySelector('input[type="checkbox"]').checked;
                break;
            default:
                input.value = row.children[x].innerHTML;
        }
    }
}

function clearForm(form) {
    var div = form.children[1];

    editor.querySelector('#edit-text').style.display = "none";
    editor.querySelector('#new-text').style.display = "block";

    for (var x = 0; x < div.childElementCount; x++) {
        let type = div.children[x].type;

        if (typeof type == "undefined") {
            continue;
        }

        switch (type) {
            case "checkbox":
                div.children[x].checked = false;
                break;
            case "number":
                div.children[x].value = "0";
                break;
            default:
                div.children[x].value = "";
        }
    }
}
