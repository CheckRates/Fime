import React from "react";

class AddImageForm extends React.Component {
    nameRef = React.createRef();
    urlRef = React.createRef();

    postImage = (event) => {
        event.preventDefault();
        const image = {
            name: this.nameRef.current.value,
            url: this.urlRef.current.value,
        }
        // Finaly add the image and refresh form
        this.props.addImage(image);
        event.currentTarget.reset();
    }

    render() {
        return(
            <form classname="post-image" onSubmit={this.postImage}>
                <input name="name" ref={this.nameRef} type="text" placeholder="Name"/>
                <input name="url" ref={this.urlRef} type="text" placeholder="URL"/>
                <button type="submit">Add Image</button>
            </form>
        )
    }
}

export default AddImageForm;