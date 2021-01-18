import React from 'react';
import Header from "./Header";
import Repository from "./Repository";
import TagDisplay from "./TagDisplay";
import SearchBar from "./SearchBar";
import AddImageForm from './AddImageForm';

import {Container} from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css"


class App extends React.Component {
  state = {
    images: {},
    tags:{},
    searchResults: {}
  };

  deleteImage = (key) => {
    const images = {...this.state.images};
    
    this.eraseImage(images[key].image.ID)
    delete images[key]
    this.setState({
      images: images
    });    
  };

  searchImage = (param) => {
    const images = [...this.state.images];

    const results = images.filter(img =>
      img.image.Name.toLowerCase().includes(param)
    );

    if(param) {
      this.setState({images: results});
    } else {
      this.setState({...this.state.images});
    }
  }

  getImage = async(id) => {
    try {
        const res = await fetch("/images/1?page=1&size=10");
        const imgs = await res.json();
        if(!imgs) return;
      
        const imgDict = Object.assign({}, ...imgs.map((i) => ({[i.image.CreatedAt]: i})));
        //console.log(imgDict)
        this.setState({
          images: imgDict
        });
    } 
    catch (error) {
      console.error(error)
    }
  }

  eraseImage = async(id) => {
    try {
        const res = await fetch("/image/" + id, {
          method: "DELETE",
        });
    } 
    catch (error) {
      console.error(error)
    }
  }

  getTags = async(id) => {
    try {
        const res = await fetch("/tags/1?page=1&size=10");
        const tgs = await res.json();
        if(!tgs) return;
        this.setState({
          tags: tgs
        });
    } 
    catch (error) {
      console.error(error)
    }
  }

  componentDidMount() {
    this.getImage();
    this.getTags();
  }

  render() {
    return (
      <div className="fime-app">
        <Header tagline="Easy to find images"/>  
      <Container>
        <div className="mt-5">
          <AddImageForm refreshImages={this.getImage}/>
        </div>
        <div className="mt-5">
          <TagDisplay tags={this.state.tags}/>
        </div>
        <div className="mt-5">
          <Repository 
            images={this.state.images}
            addImage={this.addImage} 
            deleteImage={this.deleteImage}
          />
        </div>
        </Container>  
      </div>
    )
  }
}

export default App;