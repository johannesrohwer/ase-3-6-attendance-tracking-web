// Vue component that holds the creating of groups bindings.
var createGroup = new Vue({
    el: '#createGroup',
    data: {
        group_id: '',
        name: '',
        time: '',
        place: '',
        instructor_id: '',
        instructor_name: '',
        instructor_list: []
    },
    beforeMount() {
        let self = this;
        let url = "/api/instructors";
        let params = {
            method: 'GET',
            headers: createAuthorizationHeader()
        };

        console.log(new Headers());
        fetch(url, params)
            .then(response => {
                console.log(response);
                return response.ok ? response : Promise.reject(response.statusText);
            })
            .then(response => response.json())
            .then(data => {
                self.instructor_list = data;
            }).catch(err => {
                console.log(err)
        })

    },
    methods: {
        // Submits the form to the REST API. Previously performs minor input validation.
        // TODO: improve input validation, e.g. the Room and Time slot format.
        submitData: function (event) {

            // Perform (light) input validation.
            let err = false;

            if (!isValidGroupNumber(this.group_id)) {
                err = true
            }

            if (err) {
                alert("Please make sure you have filled out all fields correctly.");
                return
            }

            // Lookup correct instructor ID
            let self = this;
            let ins_id = this.instructor_list.filter(e => {
                return e.name === self.instructor_name
            })[0].id;

            // Fix body and params for POST request, then send it.
            body = {
                "id": this.group_id,
                "time": this.time,
                "place": this.place,
                "instructor_id": ins_id
            };

            url = "/api/groups";
            params = {
                method: 'POST',
                body: JSON.stringify(body),
                headers: new Headers()
            };

            fetch(url, params)
                .then(response => {
                    return response.ok ? response : Promise.reject(response.statusText);
                })
                .then(response => response.json())
                .then(function () {
                    window.location.replace("/dashboard");
                })
                .catch(function (error) {
                    console.log(error);
                    alert("There has been an issue creating your tutorial group. Please try again.")
                })
        }
    }
});


let hideMenu = () => {
    let hiddenElements = [];
    hiddenElements.push($("#menu_login"));
    hiddenElements.push($("#menu_signup"));
    hiddenElements.forEach((e) => {
        e.hide();
    })
};

hideMenu();