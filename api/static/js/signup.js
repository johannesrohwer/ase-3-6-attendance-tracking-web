let student_signup = new Vue({
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
        let self = this;
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
            // Input validation
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


            let body = {
                "id": this.id,
                "group_id": this.group_id.id,
                "name": this.name,
                "password": this.pwd
            };

            let signupURL = "/api/students";
            let params = {
                method: 'POST',
                body: JSON.stringify(body),
                headers: new Headers()
            };


            fetch(signupURL, params)
                .then((resp) => resp.json())
                .then(function (data) {
                    sessionStorage.token = data.token;
                    sessionStorage.tokenPayload = extractJWTPayload(data.token);
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


let instructor_signup = new Vue({
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

            let body = {
                "id": this.id,
                "name": this.name,
                "password": this.pwd
            };

            let instructorSignupURL = "/api/instructors";
            let params = {
                method: 'POST',
                body: JSON.stringify(body),
                headers: new Headers()
            };


            fetch(instructorSignupURL, params)
                .then((resp) => resp.json())
                .then(function (data) {
                    sessionStorage.token = data.token;
                    sessionStorage.tokenPayload = extractJWTPayload(data.token);
                    sessionStorage.userID = data.id;
                    window.location.replace("/dashboard");
                })
                .catch(function (error) {
                    console.log(error)
                })

        }
    }
});


let hideMenu = () => {
    let hiddenElements = [];
    hiddenElements.push($("#menu_logout"));
    hiddenElements.push($("#menu_dashboard"));
    hiddenElements.push($("#menu_creategroup"));
    hiddenElements.forEach((e) => {
        e.hide();
    })
};

hideMenu();