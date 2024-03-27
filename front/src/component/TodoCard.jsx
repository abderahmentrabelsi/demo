// eslint-disable-next-line no-unused-vars
import React from 'react';
import PropTypes from 'prop-types';
import { useNavigate } from 'react-router-dom';


const TodoCard = ({ todo }) => {
    //console.log(todo)
    const navigate = useNavigate();

    const cardBgColor = todo.done ? 'bg-success' : 'bg-danger';
    const navigateToDetails = () => navigate(`/todos/${todo.ID}`);

    return (
        <div onClick={navigateToDetails} className={`card text-white ${cardBgColor} mb-3`} style={{ maxWidth: '18rem', cursor: 'pointer' }}>
            <div className="card-header">{todo.title}</div>
            <div className="card-body">
                <h5 className="card-title">Details</h5>
                <button className="bi bi-check2-circle" hidden={todo.done}>
                </button>

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
