import React from 'react';

import users from '../api/users';

// Higer Order Component (HOC) that passes user data to its wrapped component. It ensures its wrapped component is only rendered after user data hsb been loaded.
export default function withUserData(Component) {
    class UserDataWrapper extends React.Component {
        constructor(props) {
            super(props);

            this.state = {
                user: users.getUserData()
            };

            this.handleUserData = data => this.setState({ user: data });
        }

        componentDidMount() {
            users.on('userData', this.handleUserData);
        }

        componentWillUnmount() {
            users.off('userData', this.handleUserData);
        }

        render() {
            if (!this.state.user) {
                return null;
            }

            return <Component {...this.props} user={this.state.user} />;
        }
    }

    return UserDataWrapper;
}
