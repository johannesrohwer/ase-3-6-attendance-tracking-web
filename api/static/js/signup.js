var student_signup = new Vue({
    el: '#student_signup',
    data: {
        group_id: '',
        id: '',
        name: '',
        group_options: [],
        pwd: '',
        conf_pwd: ''
    },
    beforeMount() {
        var self = this;
        url = "/api/groups";
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
            let err = false;
            if (!isValidName(this.name)) {
                err = true
            }

            if (!isValidMatriculationNumber(this.id)) {
                err = true
            }

            if (!isValidPassword(this.pwd, this.conf_pwd)) {
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
            };

            url = "/api/students";
            params = {
                method: 'POST',
                body: JSON.stringify(body),
                headers: new Headers()
            };


            fetch(url, params)
                .then((resp) => resp.json())
                .then(function (data) {
                    sessionStorage.token = data.token;
                    sessionStorage.userID = data.id;
                    sessionStorage.groupID = data.group_id;

                    window.location.replace("/dashboard");
                })
                .catch(function (error) {
                    console.log(error)
                })

        }
    }
});


var instructor_signup = new Vue({
    el: '#instructor_signup',
    data: {
        id: '',
        name: '',
        pwd: '',
        conf_pwd: ''
    },
    methods: {
        submitData: function (event) {
            // input validation
            let err = false;
            if (!isValidName(this.name)) {
                err = true
            }

            if (!isValidMatriculationNumber(this.id)) {
                err = true
            }

            if (!isValidPassword(this.pwd, this.conf_pwd)) {
                err = true
            }

            if (err) {
                alert("Please make sure you have filled out all fields correctly.")
                return
            }


            // Send .post request
            body = {
                "id": this.id,
                "name": this.name,
                "password": this.pwd
            };

            url = "/api/instructors";
            params = {
                method: 'POST',
                body: JSON.stringify(body),
                headers: new Headers()
            };


            fetch(url, params)
                .then((resp) => resp.json())
                .then(function (data) {
                    sessionStorage.token = data.token;
                    sessionStorage.userID = data.id;

                    window.location.replace("/dashboard");
                })
                .catch(function (error) {
                    console.log(error)
                })

        }
    }
});


