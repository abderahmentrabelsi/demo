// src/pages/TodoDetailsPage.jsx

// eslint-disable-next-line no-unused-vars
import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { TodoService } from '../services/TodoService';

const TodoDetailsPage = () => {
    const { todoId } = useParams();
    const [todo, setTodo] = useState(null);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchTodoDetails = async () => {
            try {
                const fetchedTodo = await TodoService.getTodoById(todoId);
                setTodo(fetchedTodo);
            } catch (error) {
                setError("Failed to fetch Todo details. Please try again later.");
            }
        };

        fetchTodoDetails();
    }, [todoId]);

    if (error) return <div>{error}</div>;
    if (!todo) return <div>Loading...</div>;

    return (
        <div>
            <p>{todo.detail}</p>
        </div>
    );
};

export default TodoDetailsPage;
