import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: "",
      url: ""
    };
  }

  add = () => {
    fetch("/api/url", {
      method: "POST",
      body: JSON.stringify({
        name: this.state.name,
        url: this.state.url
      })
    })
  }

  handleChange = (event) => {
    const name = event.target.name;
    this.setState({[name]: event.target.value});
  }

  render() {
    const links = this.props.links ? "" : "You don't have any links! Try adding some"
    return (
      <div className="App">
        <header className="App-header">
          <h1 className="App-title">Prelink</h1>
        </header>
        <h4>Your links</h4>
        <p className="App-intro">
          {links}
        </p>
        <div>
          <input type='text' name='name' placeholder="Name" value={this.state.name} onChange={this.handleChange}/>
          <input type='text' name='url' placeholder="URL" value={this.state.url} onChange={this.handleChange}/>
          <button type='button' onClick={this.add}>Save</button>
        </div>
      </div>
    );
  }
}

export default App;
