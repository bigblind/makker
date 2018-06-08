import React from "react"
import {Row} from "reactstrap"

import GamesList from "./GamesList"

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
                <GamesList/>
            </Row>
        ]
    }
}