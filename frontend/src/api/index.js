import ConfigAPI from './config';
import UsersAPI from './users';
import GamesAPI from './games';

export const config = new ConfigAPI();
export const users = new UsersAPI();
export const games = new GamesAPI();
