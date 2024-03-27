// eslint-disable-next-line no-unused-vars
import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import TodosPage from './pages/TodosPage.jsx'; // Ensure correct path
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import TodoDetailsPage from "./pages/TodoDetailsPage.jsx";
import 'bootstrap-icons/font/bootstrap-icons.css';



function App() {
    return (
        <Router>
                <nav>
                    <Link to="/">Home</Link> | <Link to="/todos">Todos</Link>
                </nav>
                <Routes>
                    <Route path="/todos" element={<TodosPage />} />
                    <Route path="/" element={<Home />} />
                    <Route path="/todos/:todoId" element={<TodoDetailsPage />} />

                </Routes>
        </Router>
    );
}

function Home() {
    return (
        <div className="container my-5">
            <h2>Welcome to Your To-Do List</h2>
            <p className="mt-3 flex-start">Manage your tasks efficiently and never miss a deadline again. Get started by adding your first task.</p>
            <Link to="/todos" className="btn btn-primary mt-3">
                <i className="bi bi-pencil-square"></i> Tasks
            </Link>
        </div>
    );
}

export default App;
