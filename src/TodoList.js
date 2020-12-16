import './TodoList.css';
import React from 'react';

class TodoList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: false,
      items: [],
      inputBox: ''
    };
  }

  createNewItem = (e) => {
    if (e.key === 'Enter') {
      var formData = new FormData()
      formData.append('Title', this.state.inputBox)
      fetch("/api/newItem", { method: "POST", body: formData })
      if (this.state.items){
      this.setState((state, props) => ({
        items: [...state.items, {Title: this.state.inputBox}],
        inputBox: ''
      }))
    }
    else {
      this.setState((state, props) => ({
        items: [{Title: this.state.inputBox}],
        inputBox: ''
      }))
    }
    }
  }

  deleteItem = (title) => {
    fetch("/api/deleteItem", {
      method: "DELETE", 
      headers: {'Content-Type': 'application/json; charset=UTF-8'}, 
      body: JSON.stringify({Title: title})})
    var arr = [...this.state.items]
    var index = this.state.items.findIndex(x => x.Title === title)
    if (index !== -1) {
      arr.splice(index, 1)
      this.setState({items: arr})
    }
  }

  handleChange = (e) => { this.setState({ inputBox: e.target.value }); }

  render() {
    if (this.state.error) {
      return (<div>Hey, an error!</div>)
    } else
    return (
      <div className="main">
        Todo list
        <div>
          <input
            value={this.state.inputBox}
            onChange={this.handleChange}
            onKeyPress={this.createNewItem} 
            placeholder="Enter a new task"/>
        </div> 
        <div className="listWrapper">
          <ul className="taskList">
            {this.renderItems()}
          </ul>
        </div>
      </div>
    );
  }

  renderItems() {
    if (this.state.items) {
    return this.state.items.map(item => (
      <li className="task">
        {item.Title}
        <span
          className="deleteTaskButton"
          onClick={(e) => this.deleteItem(item.Title)}>
                    x
        </span>
      </li>
    ));
    }
  }

  componentDidMount() {
    fetch("/api/todos")
      .then(res => res.json())
      .then((result) => { this.setState({ items: result }); })
      .catch(() => {this.setState({error: true});})
  }
}

export default TodoList;
