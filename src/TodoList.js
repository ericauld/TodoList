import './TodoList.css';
import React from 'react';

class TodoList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {items: []};
  }
  
  render() {
    return (
      <div>
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
