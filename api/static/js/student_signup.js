var signup = new Vue({
    el: '#signup',
    data: {
        group_id: '',
        id: '',
        name: '',
        group_options: [],
        pwd: '',
        conf_pwd: ''
    },
    beforeMount() {
        var self = this
        url = "/api/groups"
        fetch(url)
            .then((resp) => resp.json())
            .then(function (data) {
                self.group_options = data

            })
            .catch(function (error) {
                console.log(error)
            })
    },
    methods: {
        submitData: function (event) {
            // input validation
            let err = false
            if (!validateName(this.name)) {
                err = true
            }

            if (!validateMatriculationNumber(this.id)) {
                err = true
            }

            if (!validatePassword(this.pwd, this.conf_pwd)) {
                err = true
            }

            if (err) {
                alert("Please make sure you have filled out all fields correctly.")
                return
            }


            // Send .post request
            body = {
                "id": this.id,
                "group_id": this.group_id.id,
                "name": this.name,
                "password": this.pwd
            }

            url = "http://localhost:8080/api/students"
            params = {
                method: 'POST',
                body: JSON.stringify(body),
                headers: new Headers()
            }


            fetch(url, params)
                .then((resp) => resp.json())
                .then(function (data) {
                    sessionStorage.login = JSON.stringify(data)
                    window.location.replace("/dashboard");
                })
                .catch(function (error) {
                    console.log(error)
                })

        }
    }
})


