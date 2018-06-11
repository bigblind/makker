import APIBase from "./base"
import GamesList from "../games/"

export default class GamesAPI extends APIBase {
    constructor() {
        super();
        this.instancesByGame = {}; // maps from a list of games to a list of instance IDs
        this.instances = {}; // maps from instance ids to instances
        this.init();
    }

    init(){
        GamesList.listGames().forEach((g) => {
            GamesAPI.makeRequest("GET", `/games/${g}/instances`).then((instances) => {
                console.log("loaded instances for", g, instances);

                let ids = instances.map((i) => {
                    this.instances[i.id] = i;
                    return i.id;
                });

                this.instancesByGame[g] = ids;

                this.emit("instances", {
                    game: g,
                    instances
                })
            })
        });
    }

    getInstances(game){
        let ids = this.instances[game] || [];
        return ids.map((id) => {
            return this.instances[id];
        })
    }

    getInstance(id) {
        return this.instances[id]
    }

    refreshInstance(id) {
        return GamesAPI.makeRequest("GET", `/games/instances/${id}`).then((instance) => {
            let game = instance.game_info.name;
            this.instances[id] = instance;
            this.emit("instances", {
                game,
                instances: this.getInstances(game)
            })
        })
    }
}