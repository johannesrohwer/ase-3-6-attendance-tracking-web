
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
            submitData: function (event) {
                // TODO: missing input validation before submission

                // send .post request
                body = {
                    "id": this.group_id,
                    "time": this.time,
                    "place": this.place,
                    "instructor_id": this.instructor_id
                }

                url = "/api/groups"
                params = {
                    method: 'POST',
                    body: JSON.stringify(body),
                    headers: new Headers()
                }

                fetch(url, params)
                    .then((resp) => resp.json())
                    .then(function (data) {
                        alert("Your group has been set up.")

                    })
                    .catch(function (error) {
                        console.log(error)
                        alert("There has been an issue creating your tutorial group. Please try again.")
                    })
            },
            validateGroupNumber: function() {
                // TODO: So far only checking if the group number is an integer.
                if (isInt(this.group_id)) {
                    $('#groupForm').removeClass("is-invalid")
                    return
                }
                $('#groupForm').addClass("is-invalid")
                return
            }
        }
    })
