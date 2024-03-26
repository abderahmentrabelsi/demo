// eslint-disable-next-line no-unused-vars
import React from 'react';
import PropTypes from 'prop-types';

const TodoCard = ({ todo }) => {
    // Determine the card's background color based on the todo's done status
    console.log(todo)
    const cardBgColor = todo.done ? 'bg-success' : 'bg-danger';

    return (
        <div className={`card text-white ${cardBgColor} mb-3`} style={{ maxWidth: '18rem' }}>
            <div className="card-header">{todo.title}</div>
            <div className="card-body">
                <h5 className="card-title">Details</h5>
                <p className="card-text">{todo.detail}</p>
                <button className="btn btn-primary" disabled={todo.done}>Mark as Done</button>
            </div>
        </div>
    );
};

TodoCard.propTypes = {
    todo: PropTypes.shape({
        title: PropTypes.string.isRequired,
        detail: PropTypes.string.isRequired,
        done: PropTypes.bool.isRequired
    }).isRequired
};

export default TodoCard;
