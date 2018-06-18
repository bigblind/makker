import React from "react";
import {Row, Col, Button, Badge} from "reactstrap";

const WaitingArea = ({instance, userInGame, userIsAdmin, onJoin, onLeave, onStart}) => {
    if(!instance) {
        return "...";
    }

    let joinBtn = userInGame ?
        <Button color="danger" onClick={onLeave}>Leave Game</Button> :
        <Button color="primary" onClick={onJoin}>Join Game</Button>;

    let startBtn = userIsAdmin ? <span><Button color="success" onClick={onStart}>Start Game</Button> </span> : null;

    return <Row>
        <Col>
            <h1>{instance.game_info.name}</h1>
            <p>Created: {instance.created_at}</p>
            <h2>Players</h2>
            <p>{startBtn}{joinBtn}</p>
            <ul>
                {instance.players.map((p) => {
                    let admin = instance.admin === p.user_id ? <Badge color="secondary">Game Admin</Badge> : null;
                    return <li key={p.user_id}>{p.name} {admin}</li>
                })}
            </ul>
        </Col>
    </Row>

};

export default WaitingArea;