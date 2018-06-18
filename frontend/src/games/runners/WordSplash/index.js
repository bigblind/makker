import React from "react";

export class RunnerView extends React.Component {
    render(){
        return <p>
            {JSON.stringify(this.props, 2)}
        </p>
    }
}