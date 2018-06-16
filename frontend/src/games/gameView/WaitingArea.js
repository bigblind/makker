import React from "react";
import {Row, Col, Button} from "reactstrap";

const WaitingArea = ({instance, userInGame, onJoin, onLeave}) => {
    if(!instance) {
        return "...";
    }

    let btn = userInGame ?
        <Button color="danger" onClick={onLeave}>Leave Game</Button> :
        <Button color="primary" onClick={onJoin}>Join Game</Button>;

    return <Row>
        <Col>
            <h1>{instance.game_info.name}</h1>
            <p>Created: {instance.created_at}</p>
            <h2>Players</h2>
            <p>{btn}</p>
            <ul>
                {instance.players.map((p) => {
                    return <li key={p.id}>{p.name}</li>
                })}
            </ul>
        </Col>
    </Row>

};

export default WaitingArea;