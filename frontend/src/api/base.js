import EventEmitter from "./EventEmitter"

export default class ApiBase extends EventEmitter {
    static makeRequest(method, path, body) {
        let options = {
            method: method,
            credentials: "same-origin",
        }

        if(method !== "GET" && body){
            options.body = JSON.stringify(body)
        };

        return fetch("/api" + path, options).then((resp) => {
            if(resp.status === 200) {
                return resp.json();
            } else {
                return resp.json().then(function ({error}) {
                    console.error(`${method} /api${path} returned an error: ${error}`);
                    throw error;
                });
            }
        });
    }
}