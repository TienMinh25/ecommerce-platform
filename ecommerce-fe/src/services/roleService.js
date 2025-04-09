import api from './api'; // Assuming api.js is in the same directory

/**
 * Service to handle role-related API calls
 */
const roleService = {
    getRoles: async (params = {}) => {
        try {
            const response = await api.get('/roles', { params });
            return response.data;
        } catch (error) {
            console.error('Error fetching roles:', error);
            throw error;
        }
    },
    updateRolePermissions: async (roleId, formattedPermissions) => {

    },
};

export default roleService;