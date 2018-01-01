// Checks if the two input passwords really are the same.
function isValidPassword(pwd, conf_pwd, pwd_id, conf_pwd_id) {
    if (pwd == conf_pwd) {
        $("#" + pwd_id).removeClass("is-invalid");
        $("#" + conf_pwd_id).removeClass("is-invalid");
        return true
    }
    $("#" + pwd_id).addClass("is-invalid");
    $("#" + conf_pwd_id).addClass("is-invalid");
    return false
}

// Checks that the inputted name is not empty.
function isValidName(name, id) {
    if (name) {
        $("#" + id).removeClass("is-invalid");
        return true
    }
    $("#" + id).addClass("is-invalid");
    return false
}

// Checks that the matriculation number is valid.
function isValidMatriculationNumber(matriculationNumber, id) {
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

// Creates an authorization header for a request.
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

// Returns if a value is an integer.
function isInt(value) {
    return !isNaN(value) &&
        parseInt(Number(value)) == value &&
        !isNaN(parseInt(value, 10));
}

// Performs log out.
function logOut() {
    sessionStorage.clear()
}

// Checks if an inputted value is a valid group number
function isValidGroupNumber() {
    // TODO: So far only checking if the group number is an integer.
    if (isInt(this.group_id)) {
        $('#groupForm').removeClass("is-invalid")
        return true
    }
    $('#groupForm').addClass("is-invalid")
    return false
}