import ApiBase from "./base"

class ConfigAPI extends ApiBase {
    constructor() {
        super();

        this.config = ConfigAPI.makeRequest("GET", "/config")
    }

    getConfig() {
        return this.config;
    }
}

export default new ConfigAPI();