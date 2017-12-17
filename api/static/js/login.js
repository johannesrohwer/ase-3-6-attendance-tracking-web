var login = new Vue({
    el: '#login',
    data: {
        id: "",
        pwd: ""

    },
    methods: {
        start_login: function (event) {

            // Perform (light) input validation.
            if (!validateMatriculationNumber(this.id)) {
                return
            }

            // Send post request.
            body = {
                "id": this.id,
                "password": this.pwd
            };

            url = "http://localhost:8080/api/login";
            params = {
                method: 'POST',
                body: JSON.stringify(body),
                headers: createAuthorizationHeader()
            };

            fetch(url, params)
                .then((resp) => resp.json())
                .then(function (data) {
                    console.log(data)

                    sessionStorage.login = JSON.stringify(data)
                    window.location.replace("/dashboard");
                })
                .catch(function (error) {
                    console.log(error)
                })
        }
    }
});