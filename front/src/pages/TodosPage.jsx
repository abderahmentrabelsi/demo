// src/pages/TodosPage.js

import  { useEffect, useState } from 'react';
import { TodoService } from '../services/TodoService.jsx';
import TodoCard from '../component/TodoCard.jsx';
import Container from 'react-bootstrap/Container';

const TodosPage = () => {
    const [todos, setTodos] = useState([]);

    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchTodos = async () => {
            try {
                const fetchedTodos = await TodoService.getTodos();
                setTodos(fetchedTodos);
                setError(null); // Reset error state on successful fetch
            } catch (error) {
                setError("Failed to fetch Todos. Please try again later.");
            }
        };

        fetchTodos();
    }, []);


    return (
        <Container>
            <h1>My Todos</h1>

            <div className="todos-container">
            {error && <p>{error}</p>}
            {todos.map(todo => (
                <TodoCard key={todo.ID} todo={todo} />
            ))}
            </div>
        </Container>
    );

}

export default TodosPage;
