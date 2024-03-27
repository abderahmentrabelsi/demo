import React, { useEffect, useState } from 'react';
import { TodoService } from '../services/TodoService.jsx';
import TodoCard from '../component/TodoCard.jsx';
import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import Form from 'react-bootstrap/Form';

const TodosPage = () => {
    const [todos, setTodos] = useState([]);
    const [error, setError] = useState(null);
    const [showModal, setShowModal] = useState(false);
    const [newTodo, setNewTodo] = useState({ title: '', detail: '', done: false });

    // Refactor fetchTodos to be callable
    const fetchTodos = async () => {
        try {
            const fetchedTodos = await TodoService.getTodos();
            setTodos(fetchedTodos);
            setError(null); // Reset error state on successful fetch
        } catch (error) {
            setError("Failed to fetch Todos. Please try again later.");
        }
    };

    useEffect(() => {
        fetchTodos();
    }, []);

    const handleAddTodo = async (e) => {
        e.preventDefault();
        try {
            await TodoService.addTodo(newTodo);
            await fetchTodos(); // Refetch todos to update the list
            setNewTodo({ title: '', detail: '', done: false }); // Reset the form
            setShowModal(false); // Close the modal
        } catch (error) {
            console.error("Failed to add Todo:", error);
        }
    };


    const handleModalClose = () => setShowModal(false);
    const handleModalShow = () => setShowModal(true);




    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setNewTodo(prevTodo => ({
            ...prevTodo,
            [name]: type === 'checkbox' ? checked : value
        }));
    };

    return (
        <Container>
            <h1 style={{ display: 'inline-block', marginRight: '1000px' }}>My Todos</h1>
            <Button variant="primary" onClick={handleModalShow} style={{ marginBottom: '10px' }}>Add</Button>

            <Modal show={showModal} onHide={handleModalClose}>
                <Modal.Header closeButton>
                    <Modal.Title>Add Todo</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <Form onSubmit={handleAddTodo}>
                        <Form.Group className="mb-3">
                            <Form.Label>Title</Form.Label>
                            <Form.Control
                                type="text"
                                placeholder="Enter title"
                                name="title"
                                value={newTodo.title}
                                onChange={handleChange}
                            />
                        </Form.Group>
                        <Form.Group className="mb-3">
                            <Form.Label>Detail</Form.Label>
                            <Form.Control
                                type="text"
                                placeholder="Enter detail"
                                name="detail"
                                value={newTodo.detail}
                                onChange={handleChange}
                            />
                        </Form.Group>
                        <Form.Group className="mb-3" controlId="formBasicCheckbox">
                            <Form.Check
                                type="checkbox"
                                label="Done"
                                name="done"
                                checked={newTodo.done}
                                onChange={handleChange}
                            />
                        </Form.Group>
                        <Button variant="primary" type="submit">
                            Submit
                        </Button>
                    </Form>
                </Modal.Body>
            </Modal>

            <div className="todos-container">
                {error && <p>{error}</p>}
                {todos.map(todo => (
                    <TodoCard key={todo.ID} todo={todo} />
                ))}
            </div>
        </Container>
    );
};

export default TodosPage;
