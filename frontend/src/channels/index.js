import Pusher from "pusher-js";

import {config} from "../api";

export default config.getConfig().then((cfg) => {
    return new Pusher(cfg.pusher_key, {
        "cluster": cfg.pusher_cluster,
        'authEndpoint': "/api/channels/auth",
    })
})