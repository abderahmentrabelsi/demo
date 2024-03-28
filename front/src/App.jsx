import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import TodosPage from './pages/TodosPage.jsx';
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import TodoDetailsPage from "./pages/TodoDetailsPage.jsx";
import 'bootstrap-icons/font/bootstrap-icons.css';

function App() {
    const [isDarkMode, setIsDarkMode] = useState(false);

    const toggleDarkMode = () => {
        setIsDarkMode(!isDarkMode);
    };

    useEffect(() => {
        if (isDarkMode) {
            document.documentElement.style.setProperty('--background-color', '#242424');
            document.documentElement.style.setProperty('--text-color', '#ffffff');
        } else {
            document.documentElement.style.setProperty('--background-color', '#ffffff');
            document.documentElement.style.setProperty('--text-color', '#000000');
        }
    }, [isDarkMode]);

    return (
        <Router>
            <nav>
                <div className={`toggle-button ${isDarkMode ? 'dark-mode' : ''}`} onClick={toggleDarkMode}>
                    <div className="toggle-switch"></div>
                </div>
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