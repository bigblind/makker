import React from "react"
import {Row, Col} from "reactstrap"

import GamesList from "./GamesList"
import InstancesList from "./InstancesList"

export default class Lobby extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            selectedGame: null
        }
    }

    render(){
        return [
            <Row>
                <h1>Select a Game</h1>
            </Row>,
            <Row>
                <Col><GamesList onSelected={(g) => this.setState({selectedGame: g})}/></Col>
                <Col><InstancesList game={this.state.selectedGame}/></Col>
            </Row>
        ]
    }
}