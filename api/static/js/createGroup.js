// Vue component that holds the creating of groups bindings.
var createGroup = new Vue({
    el: '#createGroup',
    data: {
        group_id: '',
        name: '',
        time: '',
        place: '',
        instructor_id: ''
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

            // Fix body and params for POST request, then send it.
            body = {
                "id": this.group_id,
                "time": this.time,
                "place": this.place,
                "instructor_id": this.instructor_id
            };

            url = "/api/groups";
            params = {
                method: 'POST',
                body: JSON.stringify(body),
                headers: new Headers()
            };

            // TODO: Instead of showing alerts, e.g. redirect the user to a certain page, like the dashboard.
            fetch(url, params)
                .then((resp) => resp.json())
                .then(function (data) {
                    alert("Your group has been set up.");
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