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

            console.log("derp");

            // Perform (light) input validation.
            if (!validateMatriculationNumber(this.id)) {
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
                body: JSON.stringify(body),
                headers: createAuthorizationHeader()
            };

            fetch(url, params)
                .then(response => { return response.ok ? response : Promise.reject(response.statusText);})
                .then(response => response.json())
                .then(function (data) {
                    sessionStorage.userID = self.id;
                    sessionStorage.token = data.token;
                    window.location.replace("/dashboard");
                })
                .catch(function (error) {
                    console.log(error);
                    self.alertMessage = error
                })
        }
    }
});