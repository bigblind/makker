import ApiBase from "./base"

export default class ConfigAPI extends ApiBase {
    constructor() {
        super();

        this.config = ConfigAPI.makeRequest("GET", "/config")
    }

    getConfig() {
        return this.config;
    }
}