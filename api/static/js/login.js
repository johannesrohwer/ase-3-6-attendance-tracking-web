let login = new Vue({
    el: '#login',
    data: {
        id: "",
        pwd: "",
        alertMessage: ""

    },
    methods: {
        start_login: function (event) {
            let self = this;

            // Perform (light) input validation.
            if (!isValidMatriculationNumber(this.id)) {
                return
            }

            let body = {
                "id": this.id,
                "password": this.pwd
            };

            let loginURL = "/api/login";
            let params = {
                method: 'POST',
                body: JSON.stringify(body),
                headers: createAuthorizationHeader()
            };

            // Login and acquire a valid token
            fetch(loginURL, params)
                .then(response => {
                    return response.ok ? response : Promise.reject(response.statusText);
                })
                .then(response => response.json())
                .then(data => {
                    sessionStorage.userID = self.id;
                    sessionStorage.token = data.token;
                    sessionStorage.tokenPayload = extractJWTPayload(data.token);
                    window.location.replace("/dashboard");
                })
                .catch(error => {
                    console.log(error);
                    self.alertMessage = error
                })
        }
    }
});


let hideMenu = () => {
    if(!isLoggedIn()) {
        let hiddenElements = [];
        hiddenElements.push($("#menu_dashboard"));
        hiddenElements.push($("#menu_creategroup"));
        hiddenElements.push($("#menu_logout"));
        hiddenElements.forEach((e) => {
            e.hide();
        })

    }
};

hideMenu();