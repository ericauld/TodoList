import './TodoList.css';
import React from 'react';

class TodoList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      items: [],
      inputBox: ''
    };
  }

  createNewItem = (e) => {
    if (e.key === 'Enter') {
      var formData = new FormData()
      formData.append('Title', this.state.inputBox)
      fetch("/api/newItem", { method: "POST", body: formData })
      if (this.state.items) {
        this.setState((state, props) => ({
          items: [...state.items, { Title: this.state.inputBox }],
          inputBox: ''
        }))
      }
      else {
        this.setState((state, props) => ({
          items: [{ Title: this.state.inputBox }],
          inputBox: ''
        }))
      }
    }
  }

  deleteItem = (title) => {
    fetch("/api/deleteItem", {
      method: "DELETE",
      headers: { 'Content-Type': 'application/json; charset=UTF-8' },
      body: JSON.stringify({ Title: title })
    })
    var arr = [...this.state.items]
    var index = this.state.items.findIndex(x => x.Title === title)
    if (index !== -1) {
      arr.splice(index, 1)
      this.setState({ items: arr })
    }
  }

  handleChange = (e) => { this.setState({ inputBox: e.target.value }); }

  render() {
    return (
      <div className="main">
        {this.renderTitleAndTextBox()}
        <div className="listWrapper">
          <ul className="taskList">
            {(() => {
              if (this.state.error)
                return this.showErrorMessage()
              else
                return this.renderItems()
            }
            )()}
          </ul>
        </div>
      </div>
    );
  }

  showErrorMessage() {
    return (
      <div>
        There's been an error loading the page. The error message
        says, "{this.state.error.message}"
      </div>
    );
  }

  renderTitleAndTextBox() {
    return (
      <div>
        Todo list
        <div>
          <input
            value={this.state.inputBox}
            onChange={this.handleChange}
            onKeyPress={this.createNewItem}
            placeholder="Enter a new task" />
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
      .then(res => {
        if (res.ok)
          return res.json()
        else
          throw new Error("The API call to get the todos list returned" 
          + " with a response that was not OK.")
      })
      .then((result) => { this.setState({ items: result }); })
      .catch((err) => { this.setState({ error: err }); })
  }
}

export default TodoList;
