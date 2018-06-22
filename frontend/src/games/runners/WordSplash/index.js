import React from "react";

export class RunnerView extends React.Component {
    render(){
        return <p>
            <pre>{JSON.stringify(this.props, null, 2)}</pre>
        </p>
    }
}