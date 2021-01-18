import React from 'react';
import Tag from "./Tag";

import {Container, Row, Col, Card, Button} from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css"


const ImagePost = (props) =>  {
    const { image, tags } = props.info;

    const handleDelete = (e) => {
        props.deleteImage(image.CreatedAt)
    }


    return (
        <Col>
            <Card className="h-100 shadow-sm bg-white rounded">
                <Card.Img variant="top" className="mb-3" src={image.URL}  alt={image.Name}/>
                <Card.Body className="d-flex flex-column">
                    <div className="d-flex mb-2 justify-content-between">
                    <Card.Title>{image.Name}</Card.Title>
                </div>
                <Card.Text className="text-secondary">{image.CreatedAt}</Card.Text>
                <Card.Text><b>Tags:</b></Card.Text>
                <Container bg="dark" variant="dark">                
                    <Row>
                        {Object.keys(tags).map(key => 
                            <Tag tag={tags[key]} />
                        )}
                    </Row>
                </Container>
                <Button className="p-2" variant="danger" onClick={() => handleDelete()}>Delete</Button>
            </Card.Body>    
        </Card>
        </Col>
    )
}

export default ImagePost;