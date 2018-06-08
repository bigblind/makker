import React from "react"
import {Card, CardDeck, CardTitle, CardText, CardBody} from "reactstrap";

import GamesList from "../index"

export default class GamesListView extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            games: GamesList.listGames(),
            selected: null
        }
    }

    render() {
        return <CardDeck>
            {this.state.games.map((g) => {
                let props = {}
                if(this.state.selected === g) {
                    props = {outline: true, color: "primary"}
                }
                return <Card {...props} onClick={() => alert(g)} key={g}>
                    <CardBody>
                        <CardTitle>{g}</CardTitle>
                    </CardBody>
                </Card>
            })}
            <Card>
                <CardBody>
                    <CardTitle>More Games Coming Soon!</CardTitle>
                    <CardText>They're in the making :)</CardText>
                </CardBody>
            </Card>
        </CardDeck>
    }
}