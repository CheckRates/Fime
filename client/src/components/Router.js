import React from "react";
import {BrowserRouter, Route, Switch} from "react-router-dom";
import Home from "./Home";
import App from "./App";
import NotFound from "./NotFound";

const Router = () => (
    <BrowserRouter>
        <Switch>
            <Route exact path="/" component={Home}/>
            <Route path="/profile/:userID" component={App}/>
            <Route component={NotFound}/>
        </Switch>
    </BrowserRouter>
);

export default Router;