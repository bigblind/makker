import React from "react";

import {games, users} from "../../api/index";
import channels from "../../channels";
import WaitingArea from "./WaitingArea"

export default class GameView extends React.Component {
    constructor(props){
        super(props);

        this.state = {
            instance: games.getInstance(props.match.params.instanceId)
        };

        games.refreshInstance(props.match.params.instanceId)
    }

    componentDidMount(){
        games.on("instances", () => {
            console.log("instances updated!");
            this.setState({instance: games.getInstance(this.props.match.params.instanceId)})
        })
    }

    componentDidUpdate(prevProps, prevState) {
        if(this.state.instance && (!prevState.instance || prevState.instance.id !== this.state.instance.id)){
            this.connectPublicChannel();
            let userInGame = this.state.instance.players.filter((p) => p.user_id === users.getUserData().id).length > 0;
            console.log("userInGame", userInGame, "userData", users.getUserData());
            if(userInGame){
                this.ensurePrivateConnection();
            }
            this.setState({userInGame});
        }
    }

    connectPublicChannel() {
        channels.then((c) => {
            if(this.publicChannelConnection){
                c.unsubscribe(this.publicChannelConnection);
            }

            this.publicChannelConnection = this.state.instance.public_channel;
            c.subscribe(this.publicChannelConnection);
        });
    }

    ensurePrivateConnection(){
        channels.then((c) => {
            if(this.privateChannelConnection !== this.state.instance.private_channel) {
                if(this.privateChannelConnection){
                    c.unsubscribe(this.privateChannelConnection);
                }

                this.privateChannelConnection = this.state.instance.private_channel;
                c.subscribe(this.privateChannelConnection);
            };
        });
    }

    render(){
        if(!this.state.instance){
            return "..."
        }

        if(this.state.instance.state === 0) {
            return <WaitingArea instance={this.state.instance} />
        }
    }
}
