// Vue app holding the dashboard bindings.
var app = new Vue({
    el: "#dashboard",
    data: {
        userID: sessionStorage.userID,
        name: '',
        group: sessionStorage.groupID,
        time: '',
        place: '',
        attendances: [],
        currentWeek: '',
        baseWeek: ''
    },
    // Fetch all needed information for the dashboard before the page mounts.
    beforeMount() {
        let self = this;

        // Set URLs.
        studentURL = "/api/students/" + self.userID;
        groupURL = "";
        attendanceURL = "/api/attendances/for/" + self.userID;

        // Fetch current week and prepare "empty" attendances
        let weekURL = "/api/week/current";
        let weekPromise = fetch(weekURL)
            .then(response => {
                return response.ok ? response : Promise.reject(response.statusText);
            })
            .then(response => response.json())
            .then(data => {
                self.currentWeek = data.current_week;
                self.baseWeek = data.base_week;
            });


        // Fetch data about the current student and the attendance data
        let studentURL = "/api/students/" + self.userID;
        let attendanceURL = "/api/attendances/for/" + self.userID + "?missing_attendances=true";

        let params = {
            method: 'GET',
            headers: createAuthorizationHeader()
        };

        // Fetch data about the current student, then get information about the group and attendances.
        fetch(studentURL, params)
            .then(response => {
                return response.ok ? response : Promise.reject(response.statusText);
            })
            .then(response => response.json())
            .then(data => {
                self.name = data.name;
                self.group = data.group_id;
                return "/api/groups/" + self.group;
            })
            .then(groupURL => fetch(groupURL, params))
            .then(response => response.json())
            .then(data => {
                self.place = data.place;
                self.time = data.time;
            })
            .then(() => fetch(attendanceURL, params))
            .then(response => response.json())
            .then(data => {
                self.attendances = data
            }).catch(error => {
                console.log(error)
            });

    }
});


function isStudentDashboard() {
    // TODO: Currently, only the student dashboard is shown, we don't have an instructor dashboard yet.
    return true
}
