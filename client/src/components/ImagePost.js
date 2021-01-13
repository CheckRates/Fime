import React from 'react';

class ImagePost extends React.Component {
    render() {
        const { url, name } = this.props.info;
        return (
            <li className="single-image">
                <img src={url}  alt={name}></img>
                <h3 className="imageName">{name}</h3>
                <button onClick={() => this.props.deleteImage(this.props.index)}>Delete Image</button>
                {/*DEBUG: GET THE DATA FROM API <h3 className="imageDate">{this.image.date}</h3>*/}
            </li>
        )
    }
}

export default ImagePost;