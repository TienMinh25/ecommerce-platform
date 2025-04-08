import api from './api.js'; // Import the axios instance

const moduleService = {
    // Get list of modules with pagination
    getModules: async (page = 1, limit = 10, getAll = false) => {
        try {
            const response = await api.get('/modules', {
                params: { page, limit, getAll }
            });
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Get a single module by ID
    getModuleById: async (id) => {
        try {
            const response = await api.get(`/modules/${id}`);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Create a new module
    createModule: async (moduleData) => {
        try {
            const response = await api.post('/modules', moduleData);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Update a module
    updateModule: async (id, moduleData) => {
        try {
            const response = await api.patch(`/modules/${id}`, moduleData);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Delete a module
    deleteModule: async (id) => {
        try {
            const response = await api.delete(`/modules/${id}`);
            return response.data;
        } catch (error) {
            throw error;
        }
    }
};

export default moduleService;