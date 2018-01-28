// Correct dashboard is selected below the definitions
let app = null;

// Register Vue components
Vue.component('group-list', {
    data: () => {
        let data = {};
        data.searchString = "";
        return data
    },
    props: ['groups'],
    computed: {
        filteredGroups: function () {
            let groupFilter = function (searchString) {
                return (group) => {
                    return searchString === "" || group.id.includes(searchString)
                }
            };

            return this.groups.filter(groupFilter(this.searchString));
        }
    },
    template: `<div>
                <div class="input-group">
                    <span class="input-group-addon">&#x1F50D;</span>
                    <input type="text" class="form-control" placeholder="Search group..." v-model="searchString">
                </div>
                <table class="table table-striped">
                    <thead>
                    <tr>
                        <th>Group Number</th>
                        <th>Time</th>
                        <th>Place</th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr v-for="group in filteredGroups">
                        <td>{{group.id}}</td>
                        <td>{{group.time}}</td>
                        <td>{{group.place}}</td>
                    </tr>
                    </tbody>
                </table>
                </div>`
});

// Register Vue components
Vue.component('attendance-list', {
    data: () => {
        let data = {};
        data.searchString = "";
        return data
    },
    props: ['attendances'],
    computed: {
        filteredAttendances: function () {
            let attendanceFilter = function (searchString) {
                return (attendance) => {
                    return searchString === "" || attendance.student_id.includes(searchString)
                }
            };

            return this.attendances.filter(attendanceFilter(this.searchString));
        }
    },
    template: `<div>
                <div class="input-group">
                    <span class="input-group-addon">&#x1F50D;</span>
                    <input type="text" class="form-control" placeholder="Search attendances..." v-model="searchString">
                </div>
                <table class="table table-striped">
                    <thead>
                    <tr>
                        <th>Attendance ID</th>
                        <th>Group ID</th>
                        <th>Student ID</th>
                        <th>Week ID</th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr v-for="attendance in filteredAttendances">
                        <td>{{attendance.id}}</td>
                        <td>{{attendance.group_id}}</td>
                        <td>{{attendance.student_id}}</td>
                        <td>{{attendance.week_id}}</td>
                    </tr>
                    </tbody>
                </table>
                </div>`
});


studentDashboard = () => {
    return new Vue({
        el: "#studentDashboard",
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
            self = this;

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
            let groupURL = "";

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
};

let instructorDashboard = () => {
    return new Vue({
        el: "#instructorDashboard",
        data: {
            userID: sessionStorage.userID,
            groups: [],
            attendances: [],
            name: ''
        },
        beforeMount: function () {

            let self = this;
            let groupURL = "/api/groups";
            let attendanceURL = "/api/attendances";
            let instructorURL = "/api/instructors/" + self.userID;
            let params = {
                method: 'GET',
                headers: createAuthorizationHeader()
            };


            fetch(groupURL, params)
                .then(response => {
                    return response.ok ? response : Promise.reject(response.statusText);
                })
                .then(response => response.json())
                .then(data => {
                    self.groups = data;
                });

            fetch(attendanceURL, params)
                .then(response => {
                    return response.ok ? response : Promise.reject(response.statusText);
                })
                .then(response => response.json())
                .then(data => {
                    self.attendances = data;
                });

            fetch(instructorURL, params)
                .then(response => {
                    return response.ok ? response : Promise.reject(response.statusText);
                })
                .then(response => response.json())
                .then(data => {
                    self.name = data.name;
                })
                .catch(err => {console.log(err)})

        }
    });
};

let hideMenuStudent = () => {
    let hiddenElements = [];
    hiddenElements.push($("#menu_login"));
    hiddenElements.push($("#menu_signup"));
    hiddenElements.push($("#menu_creategroup"));
    hiddenElements.forEach((e) => {
        e.hide();
    })
};


let hideMenuInstructor = () => {
    let hiddenElements = [];
    hiddenElements.push($("#menu_login"));
    hiddenElements.push($("#menu_signup"));
    hiddenElements.forEach((e) => {
        e.hide();
    })
};

let selectDashboard = () => {
    let tokenPayloadObj = JSON.parse(sessionStorage.tokenPayload);
    let permissions = tokenPayloadObj.credentials.permissions;

    let studentDashboardSelector = $("#studentDashboard");
    let instructorDashboardSelector = $("#instructorDashboard");

    // By default: hide everything
    studentDashboardSelector.hide();
    instructorDashboardSelector.hide();

    if (permissions.indexOf("student") !== -1) {
        hideMenuStudent();
        studentDashboardSelector.show();
        return studentDashboard();

    } else if (permissions.indexOf("instructor") !== -1) {
        hideMenuInstructor();
        instructorDashboardSelector.show();
        return instructorDashboard();

    } else {
        // Not authenticated ?
        // TODO: logout and redirect
    }
};


// Load correct dashboard
app = selectDashboard();




