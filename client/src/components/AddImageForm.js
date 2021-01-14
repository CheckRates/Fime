import React, {useState} from "react";

const  AddImageForm = () => {
    const [currentFile, setFile] = useState("");
    const [previewImg, setPreview] = useState()

    const refreshImage = (e) => {
        const image = e.target.files[0];
        showImagePreview(image)
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
        // Image request to the server
        postImage(previewImg);
    } 

    // Encode to base64 and send it to the server
    const postImage = async (base64EncodedImage) => {
        console.log(base64EncodedImage)
        try {
            await fetch("/image", {
                method: "POST",
                body: JSON.stringify({
                    name: "test",
                    image: base64EncodedImage,
                    ownerID: 1,
                    tags: [{tag: "pogs"}, {tag: "gamer"}]
                }),
                headers: {"Content-type": "application/json"}
            })
        } catch (error) {
            console.error(error)
        }
    }

    return(
        <div>
        <form onSubmit={submitImage}>
            <input 
                type="file" 
                name="image" 
                onChange={refreshImage} 
                value={currentFile}
            />
            <button className="btn-upload" type="submit">Upload</button>
        </form>
        {previewImg && (
            <img 
                src={previewImg} 
                alt="" 
                style={{height: '300px'}}/>
        )}
        </div>
    )
}

export default AddImageForm;