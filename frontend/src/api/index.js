import UsersAPI from "./users"
import GamesAPI from "./games"
import ConfigAPI from "./config"

export const users = new UsersAPI();
export const games = new GamesAPI();
export const config = new ConfigAPI();