import React from "react";

import { Navbar } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css"

const Header = (props) => (
    <Navbar expand="lg" bg="dark" variant="dark">
        <Navbar.Brand>Fime</Navbar.Brand>
        <Navbar.Text className="tagline">
            {props.tagline}
        </Navbar.Text >
    </Navbar>
);

export default Header;