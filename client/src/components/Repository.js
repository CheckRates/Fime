import React from "react";
import AddImageForm from "./AddImageForm";
import ImagePost from "./ImagePost";

class Repository extends React.Component {
    render() {
        return (
            <div className="repository">
                <h3>Image Repository</h3>
                <AddImageForm addImage={this.props.addImage} />   
                <button onClick={this.props.loadSample}>
                    Load Image Samples
                </button>   
                <ul className="images">
                    {Object.keys(this.props.images).map(key => <ImagePost key={key}>{key}</ImagePost>)}
                </ul>
            </div>
        )
    }
}

export default Repository;