import React, {useState} from "react";

import {Button, Image, Form} from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css"


const  AddImageForm = (props) => {
    const [currentFile, setFile] = useState("");
    const [previewImg, setPreview] = useState()
    const [currentName, setImageName] = useState("");
    const [currentTags, setTags] = useState();

    const refreshImage = (e) => {
        const image = e.target.files[0];
        showImagePreview(image)
    }

    const handleName = (e) => {
        const name = e.target.value;
        setImageName(name)
    }

    const handleTags = (e) => {
        const tagString = e.target.value;
        const tags = tagString.split(/(\s+)/).filter((e) => { return e.trim().length > 0; });
        setTags(tags);
    }

    const showImagePreview = (file) => {
        const reader = new FileReader();
        reader.readAsDataURL(file)
        reader.onloadend = () => {
            setPreview(reader.result)
        }
    }
    
    const submitImage = (e) => {
        e.preventDefault();
        if(!previewImg) return;
        if(!currentName) return;

        var imgTags = [];
        for (var i = 0; i < currentTags.length; i++) {
            imgTags[i] = {
                tag: currentTags[i]
            }
        }

        // Image request to the server
        postImage(previewImg, imgTags);
        setPreview()
        props.refreshImages();
    } 

    // Encode to base64 and send it to the server
    const postImage = async (base64EncodedImage, imgTags) => { 
        try {
            await fetch("/image", {
                method: "POST",
                body: JSON.stringify({
                    name: currentName,
                    image: base64EncodedImage,
                    ownerID: 1,
                    tags: imgTags
                }),
                headers: {"Content-type": "application/json"}
            })
        } catch (error) {
            console.error(error)
        }
    }

    return(
        <div>
        <Form onSubmit={submitImage}>
            <Form.Label>Post Image</Form.Label>
            
            <Form.Control 
                type="file" 
                name="image" 
                onChange={refreshImage} 
                value={currentFile}
            />
   
            <Form.Control  
                className="mt-2"
                type="text" 
                name="image" 
                placeholder="Image Name"
                onChange={handleName} 
                value={currentName}
            />

            <Form.Control  
                className="mt-2"
                type="text" 
                name="tags" 
                placeholder="Tags (Separated by spaces)"
                onChange={handleTags}
            />

            {previewImg && (
            <Image className="mt-2"
                src={previewImg} 
                alt="" 
                style={{height: '300px'}}/>
            )}
            <div className="mt-2">
                <Button className="mt-2" type="submit">Upload</Button>
            </div>
        </Form>
        </div>
    )
}

export default AddImageForm;