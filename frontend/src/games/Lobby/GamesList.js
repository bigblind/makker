import React from "react"
import {ListGroup, ListGroupItem, ListGroupItemHeading, ListGroupItemText} from "reactstrap";

import GamesList from "../index"

export default class GamesListView extends React.Component {
    constructor(props) {
        super(props)
        let games = GamesList.listGames();
        this.state = {
            selected: games[0],
            games
        };

        this.selectItem = (item, dontSetState) => {
            !dontSetState && this.setState({selected: item});

            if(this.props.onSelected){
                this.props.onSelected(item)
            }
        }

        this.selectItem(games[0], true)
    }

    render() {
        return <ListGroup>
            {this.state.games.map((g) => {
                return <ListGroupItem active={this.state.selected === g} onClick={() => this.selectItem(g)} key={g}>
                    <ListGroupItemHeading>{g}</ListGroupItemHeading>
                </ListGroupItem>
            })}
            <ListGroupItem disabled>
                <ListGroupItemHeading>More Games Coming Soon!</ListGroupItemHeading>
                <ListGroupItemText>They're in the making :)</ListGroupItemText>
            </ListGroupItem>
        </ListGroup>
    }
}