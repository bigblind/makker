import React from "react"
import {Table, Badge, Button} from "reactstrap";
import {Link, withRouter} from "react-router-dom";

import games from "../../api/games";

export default withRouter(class InstancesList extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            starting: false
        };
        if(props.game){
            this.state = {
                instances: games.getInstances(props.game)
            }
        }

        this.handleInstancesUpdate = this.handleInstancesUpdate.bind(this)
    }

    componentDidMount() {
        games.on("instances", this.handleInstancesUpdate);
    }

    componentWillUnmount(){
        games.off("instances", this.handleInstancesUpdate)
    }

    static getDerivedStateFromProps(props, state) {
        return {
            ...state,
            instances: games.getInstances(props.game)
        }
    }

    handleInstancesUpdate({game, instances}){
        if(game === this.props.game){
            this.setState({
                instances
            });
        }
    }

    startGame(){
        this.setState({starting: true});
        games.createInstance(this.props.game).then(
            instance => {
                this.setState({starting: false});
                this.props.history.push(`/instances/${instance.id}`)
            }
        )
    }

    render() {
        let instances = this.state.instances || [];

        return <Table>
            <thead>
                <tr>
                    <th>View</th>
                    <th>Players</th>
                    <th>Created</th>
                </tr>
            </thead>
            <tbody>
            {instances.map((i) => {
                return <tr key={i.id}>
                    <td><Link to={`/instances/${i.id}`}>View</Link></td>
                    <td>{i.players.map((p) => <Badge color="primary" key={p.id}>{p.name}</Badge>)}</td>
                    <td>{i.created_at}</td>
                </tr>
            })}
            <tr>
                <td colSpan="2">
                    <Button color="success" active={this.state.starting} onClick={this.startGame.bind(this)}>Start new game</Button>
                </td>
            </tr>
            </tbody>
        </Table>
    }
});

