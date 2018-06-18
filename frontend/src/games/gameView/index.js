import React from "react";

import games from "../../api/games";
import channels from "../../channels";
import WaitingArea from "./WaitingArea"
import withUserData from "../../users/withUserData";

export default withUserData(class StateManager extends React.Component {
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
            let instance = games.getInstance(this.props.match.params.instanceId);
            let userInGame = instance.players.filter((p) => p.user_id === this.props.user.id).length > 0;
            this.setState({instance, userInGame})
        })
    }

    componentDidUpdate(prevProps, prevState) {
        if(this.state.instance && (!prevState.instance || prevState.instance.id !== this.state.instance.id)){
            if(this.state.userInGame){
                this.joinGame();
            }
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
            return <WaitingArea instance={this.state.instance} userInGame={this.state.userInGame} onJoin={this.joinGame.bind(this)} onLeave={this.leaveGame.bind(this)} />
        }
    }
})
