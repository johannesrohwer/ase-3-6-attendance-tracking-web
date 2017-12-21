// Checks if the two input passwords really are the same.
function validatePassword(pwd, conf_pwd, pwd_id, conf_pwd_id) {
    if (pwd == conf_pwd) {
        $("#" + pwd_id).removeClass("is-invalid");
        $("#" + conf_pwd_id).removeClass("is-invalid");
        return true
    }
    $("#" + pwd_id).addClass("is-invalid");
    $("#" + conf_pwd_id).addClass("is-invalid");
    return false
}

function validateName(name, id) {
    if (name) {
        $("#" + id).removeClass("is-invalid");
        return true
    }
    $("#" + id).addClass("is-invalid");
    return false
}

function validateMatriculationNumber(matriculationNumber, id) {
    if (matriculationNumber.length == 8 && (matriculationNumber.match(/^[0-9]+$/) != null)) {
        $("#" + id).removeClass("is-invalid");
        return true
    }
    $("#" + id).addClass("is-invalid");
    return false
}

function handleResponseError(err) {
    if (err.status == 401) {
        alert("You are not authorised to view this site.")
    }
}

function createAuthorizationHeader() {

    let token = sessionStorage.token;

    // TODO: Check if empty token really is undefined.
    if (token == undefined) {
        return new Headers()
    }
    let header = new Headers();
    header.append("Authorization", token);
    return header
}

function isInt(value) {
    return !isNaN(value) &&
        parseInt(Number(value)) == value &&
        !isNaN(parseInt(value, 10));
}

function logOut() {
    sessionStorage.clear()
}