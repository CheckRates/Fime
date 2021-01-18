import React from "react";
import Tag from "./Tag";

import {Row, Container } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css"

const TagDisplay = (props) => {
    return (
        <div>
            <h3>Tags</h3>
            <Container className="tags">                
                <Row>
                    {Object.keys(props.tags).map(key => 
                        <Tag
                            key={key} 
                            index={key}
                            tag={props.tags[key]}
                        />
                    )}
                </Row>
            </Container>
        </div>
    )
}

export default TagDisplay;