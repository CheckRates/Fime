import React, {useState, useEffect} from "react";

import { Form } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css"

const SearchBar = (props) => {
    const [searchParam, setSearchParam] = useState("");

    const handleSearch = (e) => {
        setSearchParam(e.target.value);
        props.searchImage(searchParam);
    };

    return(
        <div className="search-bar">
        <Form>
            <Form.Control 
                type="text"
                placeholder="Search" 
                value={searchParam} 
                onChange={handleSearch}></Form.Control>
        </Form>
      </div>
    );
}

export default SearchBar;