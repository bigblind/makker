import React from "react";

import {games} from "../../api/index";
import channels from "../../channels";
import WaitingArea from "./WaitingArea"
import withUserData from "../../users/withUserData";

export default withUserData(class GameView extends React.Component {
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
            let userInGame = this.state.instance.players.filter((p) => p.user_id === this.props.user.id).length > 0;
            if(userInGame){
                this.joinGame();
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

    joinGame(){
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

    leaveGame(){
        channels.then((c) => {
            c.unsubscribe(this.privateChannelConnection)
            this.privateChannelConnection = null;
        });
    }

    render(){
        if(!this.state.instance){
            return "..."
        }

        if(this.state.instance.state === 0) {
            return <WaitingArea instance={this.state.instance} playerInGame={} onJoin={this.joinGame.bind(this)} onLeave={this.leaveGame.bind(this)} />
        }
    }
})
