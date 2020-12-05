import './TodoList.css';
import React from 'react';

class TodoList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      items: [],
      inputBox: ''
    };

    // this.createNewItem = this.createNewItem.bind(this);
    // this.handleChange = this.createNewItem.bind(this);
  }

  createNewItem = (e) => {
    if (e.key === 'Enter') {
      var formData = new FormData()
      formData.append('Title', this.state.inputBox)
      fetch("/api/newItem", { method: "POST", body: formData })
      this.setState((state, props) => ({
        items: [...state.items, { "Title": this.state.inputBox }],
        inputBox: ''
      }))
    }
  }

  deleteItem = (title) => {
    // var formData = new FormData()
    // formData.append('Title', title)
    fetch("/api/deleteItem", {
      method: "DELETE", 
      headers: {'Content-Type': 'application/json; charset=UTF-8'},
      body: JSON.stringify({Title: title})})
  }

  handleChange = (e) => { this.setState({ inputBox: e.target.value }); }

  render() {
    return (
      <div>
        <input
          value={this.state.inputBox}
          onChange={this.handleChange}
          onKeyPress={this.createNewItem} />
        <ul>
          {this.state.items.map(item => (
            <li>
              {item.Title}
              <button onClick={(e) => this.deleteItem(item.Title)}>
                Delete
              </button>
            </li>
          ))}
        </ul>
      </div>
    );
  }

  componentDidMount() {
    fetch("/api/todos")
      .then(res => res.json())
      .then((result) => { this.setState({ items: result }); })
  }
}

export default TodoList;
