import APIBase from "./base"

export default class UsersAPI extends APIBase {
    constructor() {
        super()
        this.userData = null;
    }

    getUserData(retried) {
        if(!this.userData) {
            UsersAPI.makeRequest("GET", "/users/me").then((userData) => {
                if(userData.id === "" && !retried) {
                    this.getUserData(true) // retry once
                }
                this.userData = userData;
                this.emit("userData", userData);
            })
        }

        return this.userData
    }
}