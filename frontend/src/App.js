import React, { Component } from 'react';
import {Container} from "reactstrap";
import {BrowserRouter as Router, Route} from "react-router-dom";

import MenuBar from "./layout/MenuBar"

import Lobby from "./games/Lobby"
import GameView from "./games/gameView";

class App extends Component {
  render() {
      return [
          <MenuBar key="m" />,
          <Router key="r">
              <Container>
                  <Route path="/" exact component={Lobby}/>
                  <Route path="/instances/:instanceId" component={GameView} />
              </Container>
          </Router>
      ]
  }
}

export default App;
