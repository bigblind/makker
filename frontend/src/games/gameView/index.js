import React from "react";

import games from "../../api/games";
import gamesList from "../index";
import channels from "../../channels";
import WaitingArea from "./WaitingArea";
import withUserData from "../../users/withUserData";

export default withUserData(class StateManager extends React.Component {
    constructor(props){
        super(props);

        this.state = {
            instance: games.getInstance(props.match.params.instanceId),
            runnerLoaded: false
        };

        games.refreshInstance(props.match.params.instanceId)
    }

    componentDidMount(){
        games.on("instances", () => {
            console.log("instances updated!");
            let instance = games.getInstance(this.props.match.params.instanceId);
            let userInGame = instance.players.filter((p) => p.user_id === this.props.user.id).length > 0;
            let userIsAdmin = instance.admin === this.props.user.id;
            this.setState({instance, userInGame, userIsAdmin});
        })
    }

    componentDidUpdate(prevProps, prevState) {
        if(this.state.instance && (!prevState.instance || prevState.instance.id !== this.state.instance.id)){
            if(this.state.userInGame){
                this.joinGame();
            }

            this.setState({runnerLoaded: false});
            let info = this.state.instance.game_info;
            gamesList.loadGame(info.name, info.version).then((runner) => {
                this.runner = runner;
                this.setState({runnerLoaded: true});
            });
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

    startGame(){
        games.startGame(this.state.instance.id);
    }

    render(){
        if(!this.state.instance){
            return "..."
        }

        if(this.state.instance.state === 0) {
            return <WaitingArea instance={this.state.instance}
                                userInGame={this.state.userInGame}
                                userIsAdmin={this.state.userIsAdmin}
                                onJoin={this.joinGame.bind(this)}
                                onLeave={this.leaveGame.bind(this)}
                                onStart={this.startGame.bind(this)} />
        }

        if(this.state.instance.state === 1){
            if(!this.state.runnerLoaded){
                return "Loading..."
            }

            return <this.runner.RunnerView instance={this.state.instance} />;
        }
    }
})
