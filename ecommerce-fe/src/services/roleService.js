import api from './api'; // Assuming api.js is in the same directory

/**
 * Service to handle role-related API calls
 */
const roleService = {
    /**
     * Fetch all roles from the API
     * @returns {Promise} Promise that resolves to an array of role objects
     */
    getRoles: async () => {
        try {
            const response = await api.get('/roles');

            // Return the data property from the response
            return response.data.data;
        } catch (error) {
            console.error('Error fetching roles:', error);
            throw error;
        }
    }
};

export default roleService;