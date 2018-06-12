import React from "react";
import {Row, Col} from "reactstrap";

const WaitingArea = ({instance}) => {
    if(!instance) {
        return "...";
    }

    return <Row>
        <Col>
            <h1>{instance.game_info.name}</h1>
            <p>Created: {instance.created_at}</p>
            <h2>Players</h2>
            <ul>
                {instance.players.map((p) => {
                    return <li key={p.id}>{p.name}</li>
                })}
            </ul>
        </Col>
    </Row>

};

export default WaitingArea;