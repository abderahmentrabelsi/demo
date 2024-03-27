import React, { useState } from 'react';

// Simple custom modal component
const Modal = ({ isOpen, onClose, onSubmit }) => {
    const [todo, setTodo] = useState({
        title: '',
        detail: '',
        done: false,
    });

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setTodo(prevTodo => ({
            ...prevTodo,
            [name]: type === 'checkbox' ? checked : value
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        onSubmit(todo);
    };

    if (!isOpen) return null;

    return (
        <div style={{
            position: 'fixed',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            backgroundColor: '#fff',
            padding: '20px',
            zIndex: 1000,
        }}>
            <form onSubmit={handleSubmit}>
                <label>Title:</label>
                <input name="title" value={todo.title} onChange={handleChange} />
                <label>Detail:</label>
                <input name="detail" value={todo.detail} onChange={handleChange} />
                <label>Done:</label>
                <input name="done" type="checkbox" checked={todo.done} onChange={handleChange} />
                <button type="submit">Add Todo</button>
                <button onClick={onClose}>Close</button>
            </form>
        </div>
    );
};
