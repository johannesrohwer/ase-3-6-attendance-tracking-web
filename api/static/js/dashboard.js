var app = new Vue({
    el: "#dashboard",
    data: {
        userID: sessionStorage.userID,
        name: '',
        group: sessionStorage.groupID,
        time: '',
        place: '',
        items: []
    },
    beforeMount() {
        var self = this
        studentURL = "/api/students/" + self.userID;
        groupURL = "";
        attendanceURL = "/api/attendances/for/" + self.userID;

        params = {
            method: 'GET',
            headers: createAuthorizationHeader()
        };

        // Fetch data about the current student
        fetch(studentURL, params)
            .then(response => { return response.ok ? response : Promise.reject(response.statusText);})
            .then(response => response.json())
            .then(data => {
                self.name = data.name;
                self.group = data.group_id;
                groupURL = "/api/groups/" + self.group;
            })
            .then(() => fetch(groupURL, params))
            .then(response => response.json())
            .then(function (data) {
                self.place = data.place;
                self.time = data.time;
            })
        . then(() => fetch(attendanceURL, params))
            .then(response => response.json())
            .then(function (data) {
                self.items = data
            }).catch(function(error) {
                console.log(error)
            }
        )
    }

});


function isStudentDashboard() {
    // TODO: FIXME
    return true
}
