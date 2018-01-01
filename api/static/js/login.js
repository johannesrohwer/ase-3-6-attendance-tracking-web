var login = new Vue({
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

            // Send post request.
            body = {
                "id": this.id,
                "password": this.pwd
            };

            url = "/api/login";
            params = {
                method: 'POST',
                body: JSON.stringify(body),cu
                headers: createAuthorizationHeader()
            };

            fetch(url, params)
                .then(response => {
                    return response.ok ? response : Promise.reject(response.statusText);
                })
                .then(response => response.json())
                .then(data => {
                    sessionStorage.userID = self.id;
                    sessionStorage.token = data.token;
                    window.location.replace("/dashboard");
                })
                .catch(error => {
                    console.log(error);
                    self.alertMessage = error
                })
        }
    }
});