import React from "react";
import {users} from "../api"

import {Navbar, NavbarBrand, Nav, NavLink, NavItem} from "reactstrap";

export default class MenuBar extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            userData: users.getUserData()
        };

        this.onUserData = (d) => this.setState({userData: d});
        console.log("render", this.render)
    }

    componentDidMount(){
        users.on("userData", this.onUserData);
    }

    componentWillUnmount(){
        users.off("userData", this.onUserData);
    }

    render() {
        let userMenu = "";
        if(this.state.userData) {
            userMenu = <Nav className="ml-auto">
                <NavItem>
                    <NavLink href="/profile">{this.state.userData.name}</NavLink>
                </NavItem>
            </Nav>
        }

        let resp = (<Navbar color="dark" dark>
            <NavbarBrand href="/">MAKKER</NavbarBrand>
            {userMenu}
        </Navbar>);

        return resp
    }
}

