import api from './api';

// User API Services
const userService = {
    // Get users with pagination and filters
    getUsers: async (params) => {
        try {
            const response = await api.get('/users', { params });
            return response.data;
        } catch (error) {
            console.error('Error fetching users:', error);
            throw error;
        }
    },

    // Create new user
    createUser: async (userData) => {
        try {
            const response = await api.post('/users', userData);
            return response.data;
        } catch (error) {
            console.error('Error creating user:', error);
            throw error;
        }
    },

    // Update user by admin
    updateUser: async (userId, updateData) => {
        try {
            await api.patch(`/users/${userId}`, updateData);
        } catch (error) {
            console.error(`Error updating user with id ${userId}:`, error);
            throw error;
        }
    },

    // Delete user by admin
    deleteUser: async (userId) => {
        try {
            await api.delete(`/users/${userId}`);
        } catch (error) {
            console.error(`Error deleting user with id ${userId}:`, error);
            throw error;
        }
    },
};

export default userService;