

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
                    console.log(data)
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
                if (!this.validateName()) {
                    err = true
                }

                if (!this.validateMatriculationNumber()) {
                    err = true
                }

                if (!this.validatePassword()) {
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

                console.log(params.body)

                fetch(url, params)
                    .then((resp) => resp.json())
                    .then(function (data) {
                        sessionStorage.login = JSON.stringify(data)
                        window.location.replace("/dashboard");
                    })
                    .catch(function (error) {
                        console.log(error)
                    })

            },
            validateMatriculationNumber: function() {
                if (this.id.length == 8 && (this.id.match(/^[0-9]+$/) != null)) {
                    $('#matriculationNumberForm').removeClass("is-invalid")
                    return true
                }
                $('#matriculationNumberForm').addClass("is-invalid")
                return false
            },
            validateName: function() {
                if (this.name) {
                    $('#fullNameForm').removeClass("is-invalid")
                    return true
                }
                $('#fullNameForm').addClass("is-invalid")
                return false
            },

            // Checks if the two input passwords really are the same.
            validatePassword: function() {
                if (this.pwd == this.conf_pwd) {
                    $("#confPwdForm").removeClass("is-invalid")
                    $("#pwdForm").removeClass("is-invalid")
                    return true
                }
                $("#confPwdForm").addClass("is-invalid")
                $("#pwdForm").addClass("is-invalid")
                return false
            }

        }
    })


