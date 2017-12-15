var login = new Vue({
    el: '#login',
    data: {
        id: '',
        pwd: ""

    },
    methods: {
        start_login: function (event) {
            // TODO: input validation

            // send .post request

            body = {
                "id": this.id,
                "password": this.pwd
            }

            url = "http://localhost:8080/api/students/" + this.id
            params = {
                method: 'POST',
                body: JSON.stringify(body),
                headers: new Headers()
            }

            if (this.id.length == 0) {
                alert("This is not a valid matriculation number.")
                return
            }


            fetch(url, params)
                .then((resp) => resp.json())
                .then(function (data) {

                    if (Object.keys(data).length === 0 && data.constructor === Object) {
                        alert("This is not a valid matriculation number.")
                    } else {
                        sessionStorage.login = JSON.stringify(data)
                        window.location.replace("/dashboard");
                    }

                })
                .catch(function (error) {
                    console.log(error)
                })

        }
    }
})