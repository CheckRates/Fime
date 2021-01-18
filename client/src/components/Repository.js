import React from "react";
import ImagePost from "./ImagePost";

import {Container, Row } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css"

const Repository = (props) => {


    return (
        <div>
            <h3>Repo</h3>
            <Container className="images">                
                <Row>
                    {Object.keys(props.images).map(key => 
                        <ImagePost
                            key={key}
                            deleteImage={props.deleteImage}
                            info={props.images[key]}
                        />
                )}
                </Row>
            </Container>
        </div>
    )
}

export default Repository;