import APIBase from "./base"

export default class UsersAPI extends APIBase {
    constructor() {
        super()
        this.userData = null;
        this.listeners = new Map();
    }

    on(event, handler) {
        this.listeners.has(event) || this.listeners.set(event, [])
        this.listeners.get(event).push(handler);
    }

    off(event, handler){
        let listeners = this.listeners.get(event);

        if (listeners && listeners.length) {
            for (let i=0; i<listeners.length; i++) {
                if (listeners[i] === handler) {
                    listeners.splice(i, 1);
                    this.listeners.set(event, listeners);
                    return true
                }
            }
        }

        return false;
    }

    emit(event, ...args) {
        let listeners = this.listeners.get(event)
        if (listeners) {
            listeners.forEach((l) => {
                console.log(l)
                l(...args);
            })
        }
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