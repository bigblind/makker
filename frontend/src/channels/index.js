import Pusher from "pusher-js";
import config from "../api/config";

export default new Promise((resolve, reject) => {
    config.getConfig().then((cfg) => {
        resolve(new Pusher(cfg.pusher_key, {
            "cluster": cfg.pusher_cluster,
            'authEndpoint': "/api/channels/auth",
        }));
    })
});