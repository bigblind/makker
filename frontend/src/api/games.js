import APIBase from "./base"
import GamesList from "../games/"

export default class GamesAPI extends APIBase {
    constructor() {
        super();
        this.instances = {}; // maps from a list of games to a list of instances
        this.init();
    }

    init(){
        GamesList.listGames().forEach((g) => {
            GamesAPI.makeRequest("GET", `/games/${g}/instances`).then((instances) => {
                console.log("loaded instances for", g, instances);
                this.instances[g] = instances;
                this.emit("instances", {
                    game: g,
                    instances
                })
            })
        });
    }

    getInstances(game){
        return this.instances[game] || [];
    }
}