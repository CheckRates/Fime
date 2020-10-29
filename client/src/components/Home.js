import React from "react";
import Header from "./Header";

class Home extends React.Component {
    userID = React.createRef();

    goToProfile = (event) => {
        event.preventDefault();
        const profileID = this.userID.current.value;
        // TODO: Redirect to the Auth0 Login Page
        // Should retrieve somewhat a data of the user to identify which
        // profile page should be retrieved
        // It will check the database so that if it matches then we are good
        // to go. For now I will retrieve the user ID from the form
        
        // Redirect
        this.props.history.push(`/profile/${profileID}`)
    }

    componentDidMount() {
        console.log("mounted!")
    }
 
    render() {
        return (
            <div className="home">
                <Header/>
                <form className="login-form" onSubmit={this.goToProfile}>
                    <h3>Please Login to continue</h3>
                    <input 
                     type="text" 
                     ref={this.userID} 
                     placeholder="User ID"
                     required />
                    <button type="submit">Login</button>
                </form>
            </div>
        )
    }
}

export default Home;

