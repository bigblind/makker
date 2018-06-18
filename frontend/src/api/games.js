import APIBase from "./base";
import config from "./config";
import GamesList from "../games/";
import channels from "../channels";

class GamesAPI extends APIBase {
    constructor() {
        super();
        this.instancesByGame = {}; // maps from a list of games to a list of instance IDs
        this.instances = {}; // maps from instance ids to instances
        this.initGames();
        this.connectToLobby();
    }

    initGames(){
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

    connectToLobby(){

        Promise.all([config.getConfig(), channels]).then(([config, channels]) => {
            let chan = channels.subscribe(config.lobby_channel);

            chan.bind("created", async (id) => {
                console.log("created fired");
                let inst = await this.refreshInstance(id, true);
                this.instancesByGame[inst.game_info.name].push(id);
                this.emit("instances", {
                    game: inst.game_info.name,
                    instances: this.getInstances(inst.game_info.name)
                })
            });

            chan.bind("update", async (id) => {
                let inst = await this.refreshInstance(id);
                if(inst.state !== 0) {
                    let ids = this.instancesByGame[inst.game_info.name];
                    let index = ids.indexOf(id);
                    if(index > -1){
                        ids.splice(index, 1)
                    }
                }
            })
        })
    }

    getInstances(game){
        let ids = this.instancesByGame[game] || [];
        return ids.map((id) => {
            return this.instances[id];
        })
    }

    getInstance(id) {
        return this.instances[id]
    }

    refreshInstance(id, noEmit) {
        return GamesAPI.makeRequest("GET", `/games/instances/${id}`).then((instance) => {
            let game = instance.game_info.name;
            this.instances[id] = instance;
            !noEmit && this.emit("instances", {
                game,
                instances: this.getInstances(game)
            });
            return instance;
        })
    }

    startGame(id){
        return GamesAPI.makeRequest("POST", `/games/instances/${id}/start`);
    }
}

export default new GamesAPI();