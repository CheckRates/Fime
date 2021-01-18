import React from "react";

import { Button } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css"

const Tag = (props) => (
    <Button className="m-1 rounded" 
        expand="lg" bg="blue">{props.tag.Name}</Button>
);

export default Tag;