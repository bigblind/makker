export default class ApiBase {
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
                return r.json();
            } else {
                return r.json().then(function ({error}) {
                    console.error(`${method} /api${path} returned an error: ${error}`);
                    throw error;
                });
            }
        });
    }
}