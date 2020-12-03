import './TodoList.css';
import React from 'react';

class TodoList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {items: [],
                  inputBox: ''};
  }
  
  createNewItem = (e) => {
    if (e.key === 'Enter') {
      var formData = new FormData()
      formData.append('Title', this.state.inputBox)
      fetch("/api/newItem", {method: "POST", body: formData})
      this.setState((state, props) => ({
        items: [...state.items, {"Title": this.state.inputBox}],
        inputBox: ''
      }))
    }
  }

  handleChange = (e) => {this.setState({inputBox: e.target.value});}

  render() {
    return (
      <div>
        <input onChange={this.handleChange} onKeyPress={this.createNewItem }/>
        <ul>
          {this.state.items.map(item => (
            <li>
              {item.Title}
            </li>
          ))}
        </ul>
      </div>
    );
  }

  componentDidMount() {
    fetch("/api/todos")
      .then(res => res.json())
      .then((result) => {this.setState({items: result});})
  }
}

export default TodoList;
