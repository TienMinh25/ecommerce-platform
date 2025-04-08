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

    // Get single user by id
    getUserById: async (userId) => {
        try {
            const response = await api.get(`/users/${userId}`);
            return response.data;
        } catch (error) {
            console.error(`Error fetching user with id ${userId}:`, error);
            throw error;
        }
    },

    // Get user permissions (assuming an endpoint like this exists)
    getUserPermissions: async (userId) => {
        try {
            const response = await api.get(`/users/${userId}/permissions`);
            return response.data;
        } catch (error) {
            console.error(`Error fetching permissions for user ${userId}:`, error);
            throw error;
        }
    },

    // Update user permissions (placeholder for future implementation)
    updateUserPermissions: async (userId, permissions) => {
        try {
            const response = await api.patch(`/users/${userId}/permissions`, { permissions });
            return response.data;
        } catch (error) {
            console.error(`Error updating permissions for user ${userId}:`, error);
            throw error;
        }
    }
};

export default userService;