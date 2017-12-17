// Checks if the two input passwords really are the same.
function validatePassword(pwd, conf_pwd) {
    if (pwd == conf_pwd) {
        $("#confPwdForm").removeClass("is-invalid")
        $("#pwdForm").removeClass("is-invalid")
        return true
    }
    $("#confPwdForm").addClass("is-invalid")
    $("#pwdForm").addClass("is-invalid")
    return false
}

function validateName(name) {
    if (name) {
        $('#fullNameForm').removeClass("is-invalid")
        return true
    }
    $('#fullNameForm').addClass("is-invalid")
    return false
}

function validateMatriculationNumber(matriculationNumber) {
    if (matriculationNumber.length == 8 && (matriculationNumber.match(/^[0-9]+$/) != null)) {
        $('#matriculationNumberForm').removeClass("is-invalid")
        return true
    }
    $('#matriculationNumberForm').addClass("is-invalid")
    return false
}

function handleResponseError(err) {
    if (err.status == 401) {
        alert("You are not authorised to view this site.")
    }
}

function createAuthorizationHeader() {

    token = sessionStorage.token

    // TODO: Check if empty token really is undefined.
    if (token == undefined) {
        return new Headers()
    }
    header = new Headers()
    header.append("Authorization", token)
    return header
}

function isInt(value) {
    return !isNaN(value) &&
        parseInt(Number(value)) == value &&
        !isNaN(parseInt(value, 10));
}