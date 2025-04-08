import api from './api.js'; // Import the axios instance

const permissionService = {
    // Get list of permissions with pagination
    getPermissions: async (page = 1, limit = 10, getAll = false) => {
        try {
            const response = await api.get('/permissions', {
                params: { page, limit, getAll }
            });
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Get a single permission by ID
    getPermissionById: async (id) => {
        try {
            const response = await api.get(`/permissions/${id}`);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Create a new permission
    createPermission: async (permissionData) => {
        try {
            const response = await api.post('/permissions', permissionData);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Update a permission
    updatePermission: async (id, permissionData) => {
        try {
            const response = await api.patch(`/permissions/${id}`, permissionData);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Delete a permission
    deletePermission: async (id) => {
        try {
            const response = await api.delete(`/permissions/${id}`);
            return response.data;
        } catch (error) {
            throw error;
        }
    }
};

export default permissionService;