// src/pages/TodoDetailsPage.jsx
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
                setTodo(fetchedTodo); // Assuming fetchedTodo has the structure { Todo: { ... } }
            } catch (error) {
                setError("Failed to fetch Todo details. Please try again later.");
            }
        };

        fetchTodoDetails();
    }, [todoId]);

    if (error) return <div>{error}</div>;
    if (!todo) return <div>Loading...</div>;

    // Make sure to access the Todo object for detail and CreatedAt
    const { Todo: todoDetails } = todo;

    return (
        <div>
            <h2>Task Details</h2>
            <p><strong>Detail:</strong> {todoDetails.detail}</p>
            <p><strong>Created At:</strong> {new Date(todoDetails.CreatedAt).toLocaleString()}</p>
        </div>
    );
};

export default TodoDetailsPage;
