import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import TodosPage from './pages/TodosPage.jsx'; // Ensure correct path
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';


function App() {
    return (
        <Router>
            <div>
                <nav>
                    <Link to="/">Home</Link> | <Link to="/todos">Todos</Link>
                </nav>
                <Routes>
                    <Route path="/todos" element={<TodosPage />} />
                    <Route path="/" element={<Home />} />
                </Routes>
            </div>
        </Router>
    );
}

function Home() {
    return (
        <div>
            <h2>Home Page</h2>
            {/* Home page content */}
        </div>
    );
}

export default App;
