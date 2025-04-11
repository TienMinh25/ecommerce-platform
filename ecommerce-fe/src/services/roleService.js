import api from './api'; // Assuming api.js is in the same directory

/**
 * Service to handle role-related API calls
 */
const roleService = {
    // Get roles with pagination and filters
    getRoles: async (params = {}) => {
        try {
            const response = await api.get('/roles', { params });
            return response.data;
        } catch (error) {
            console.error('Error fetching roles:', error);
            throw error;
        }
    },

    // Get a single role by ID
    getRoleById: async (roleId) => {
        try {
            const response = await api.get(`/roles/${roleId}`);
            return response.data;
        } catch (error) {
            console.error(`Error fetching role with id ${roleId}:`, error);
            throw error;
        }
    },

    // Create a new role
    createRole: async (roleData) => {
        try {
            const response = await api.post('/roles', roleData);
            return response.data;
        } catch (error) {
            console.error('Error creating role:', error);
            throw error;
        }
    },

    // Update a role's basic info
    updateRole: async (roleId, roleData) => {
        try {
            const response = await api.patch(`/roles/${roleId}`, roleData);
            return response.data;
        } catch (error) {
            console.error(`Error updating role with id ${roleId}:`, error);
            throw error;
        }
    },

    // Update role permissions
    updateRolePermissions: async (roleId, permissionsData) => {
        try {
            const response = await api.patch(`/roles/${roleId}`, permissionsData);
            return response.data;
        } catch (error) {
            console.error(`Error updating permissions for role ${roleId}:`, error);
            throw error;
        }
    },

    // Delete a role
    deleteRole: async (roleId) => {
        try {
            const response = await api.delete(`/roles/${roleId}`);
            return response.data;
        } catch (error) {
            console.error(`Error deleting role with id ${roleId}:`, error);
            throw error;
        }
    }
};

export default roleService;