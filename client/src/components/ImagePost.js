import React from 'react';

class ImagePost extends React.Component {
    render() {
        return (
            <div className="single-image">
                <img src="https://camo.githubusercontent.com/e83482e80ab2ed5bd3d9781a3b4602cdcb3407aa/687474703a2f2f692e696d6775722e636f6d2f485379686177742e6a7067" alt="Cool Gopher"></img>
                <h3 className="imageName">Test</h3>
                <h3 className="imageDate">2020-09-29</h3>
            </div>
        )
    }
}

export default ImagePost;