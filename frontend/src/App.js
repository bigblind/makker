import React, { Component } from 'react';
import {Container} from "reactstrap";
import {BrowserRouter as Router, Route} from "react-router-dom";

import MenuBar from "./layout/MenuBar"

class App extends Component {
  render() {
    return (
      <Container >
          <MenuBar/>
          <Router>

          </Router>
      </Container>
    );
  }
}

export default App;
