import React from "react"
import {Table, Badge} from "reactstrap";

import {games} from "../../api";

export default class InstancesList extends React.Component {
    constructor(props){
        super(props);
        this.state = {};
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
        games.off("instances", this.handleInstancesUpdate())
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

    render() {
        let instances = this.state.instances || [];

        return <Table>
            <thead>
                <tr>
                    <th>Players</th>
                    <th>Created</th>
                </tr>
            </thead>
            <tbody>
            {instances.map((i) => {
                return <tr key={i.id}>
                    <td>{i.players.map((p) => <Badge color="primary" key={p.id}>{p.name}</Badge>)}</td>
                    <td>{i.created_at}</td>
                </tr>
            })}
            </tbody>
        </Table>
    }
}

