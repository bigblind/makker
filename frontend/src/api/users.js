import APIBase from "./base"

class UsersAPI extends APIBase {
    constructor() {
        super();
        this.userData = null;
        let serialized = localStorage.userData;
        if(serialized){
            this.userData = JSON.parse(serialized);
        }
    }

    getUserData(retried) {
        if(!this.userData) {
            UsersAPI.makeRequest("GET", "/users/me").then((userData) => {
                if(userData.id === "" && !retried) {
                    this.getUserData(true) // retry once
                }
                this.userData = userData;
                localStorage.setItem("userData", JSON.stringify(userData));
                this.emit("userData", userData);
            })
        }

        return this.userData
    }
}

export default new UsersAPI();