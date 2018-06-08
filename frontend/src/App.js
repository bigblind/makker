import React, { Component } from 'react';
import {Container} from "reactstrap";
import {BrowserRouter as Router, Route} from "react-router-dom";

import MenuBar from "./layout/MenuBar"

import Lobby from "./games/Lobby"

class App extends Component {
  render() {
      return [
          <MenuBar/>,
          <Container>
              <Router>
                  <Route path="/" exact component={Lobby}/>
              </Router>
          </Container>
      ]
  }
}

export default App;
