import React from 'react';
import Header from "./Header";
import Repository from "./Repository";
import TagDisplay from "./TagDisplay";

// DEBUG: 
import sampleImages from "../util/sampleImages";

class App extends React.Component {
  state = {
    images: {}
  };

  addImage = (image) => {
    const images = {...this.state.images};
    images[`image${Date.now()}`] = image;
    this.setState({
      images: images
    });
  };

  // DEBUG: Function
  loadSample = () => {
    this.setState({images: sampleImages})
  };

  render() {
    return (
      <div className="fime-app">
        <div className="profile">
          <Header tagline="Easy to find memes"/>
        </div>
          <Repository 
            images={this.state.images}
            addImage={this.addImage} 
            loadSample={this.loadSample}/>
          <TagDisplay/>
      </div>
    )
  }
}

export default App;