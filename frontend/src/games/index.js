class GamesList {
    constructor() {
        this.games = {}

        this.addGame("WordSplash", 1, "WordSplash")
    }

    addGame(name, version, runner) {
        if (!this.games[name]) {
            this.games[name] = {}
        }

        this.games[name][version] = runner;
    }

    getGame(name, version) {
        return this.games[name][version];
    }

    listGames() {
        return Object.keys(this.games);
    }

    loadGame(name, version){
        return import(`./runners/${this.games[name][version]}`);
    }
}

export default new GamesList();