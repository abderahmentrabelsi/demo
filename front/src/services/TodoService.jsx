// src/services/TodoService.js

import axios from 'axios';

const baseURL = 'http://localhost:8086';

export const TodoService = {
    getTodos: async () => {
        try {
            const response = await axios.get(`${baseURL}/todos`);
            return response.data.Todos;
        } catch (error) {
            console.error("Error fetching Todos:", error);
            throw error;
        }
    },

    getTodoById: async (id) => {
        try {
            const response = await axios.get(`${baseURL}/todos/${id}`);
            return response.data;
        } catch (error) {
            console.error("Error fetching Todo by ID:", error);
            throw error;
        }
    },
    addTodo: async (todo) => {
        try {
            const response = await axios.post(`${baseURL}/todos`, todo);
            // Ensure this returns the full todo item, including any server-generated fields like ID
            return response.data; // This should be the new todo item as returned by your API
        } catch (error) {
            console.error("Error adding new Todo:", error);
            throw error;
        }
    },

};
