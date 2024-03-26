// eslint-disable-next-line no-unused-vars
import React from 'react';
import PropTypes from 'prop-types';
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/Button';

const TodoCard = ({ todo }) => (
    <Card style={{ width: 'auto', marginBottom: 'auto' }}>
        <Card.Body>
            <Card.Title>{todo.title}</Card.Title>
            <Card.Text>
                {todo.detail}
            </Card.Text>
            <Button variant="primary" disabled={todo.done}>Mark as Done</Button>
        </Card.Body>
    </Card>
);

TodoCard.propTypes = {
    todo: PropTypes.shape({
        title: PropTypes.string.isRequired,
        detail: PropTypes.string.isRequired,
        done: PropTypes.bool.isRequired
    }).isRequired
};

export default TodoCard;
