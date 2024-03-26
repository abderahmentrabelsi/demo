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
};
