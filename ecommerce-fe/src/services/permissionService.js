import api from './api.js'; // Import the axios instance

const permissionService = {
    // Get list of permissions with pagination or getAll
    getPermissions: async (params = {}) => {
        try {
            const response = await api.get('/permissions', { params });
            return response.data;
        } catch (error) {
            console.error('Error fetching permissions:', error);
            throw error;
        }
    },

    // Get a single permission by ID
    getPermissionById: async (id) => {
        try {
            const response = await api.get(`/permissions/${id}`);
            return response.data;
        } catch (error) {
            console.error(`Error fetching permission with id ${id}:`, error);
            throw error;
        }
    },

    // Create a new permission
    createPermission: async (permissionData) => {
        try {
            const response = await api.post('/permissions', permissionData);
            return response.data;
        } catch (error) {
            console.error('Error creating permission:', error);
            throw error;
        }
    },

    // Update a permission
    updatePermission: async (id, permissionData) => {
        try {
            const response = await api.patch(`/permissions/${id}`, permissionData);
            return response.data;
        } catch (error) {
            console.error(`Error updating permission with id ${id}:`, error);
            throw error;
        }
    },

    // Delete a permission
    deletePermission: async (id) => {
        try {
            const response = await api.delete(`/permissions/${id}`);
            return response.data;
        } catch (error) {
            console.error(`Error deleting permission with id ${id}:`, error);
            throw error;
        }
    }
};

export default permissionService;