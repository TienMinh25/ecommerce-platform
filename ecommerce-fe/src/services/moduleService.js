import api from './api.js'; // Import the axios instance

const moduleService = {
    // Get list of modules with pagination or getAll
    getModules: async (params = {}) => {
        try {
            const response = await api.get('/modules', { params });
            return response.data;
        } catch (error) {
            console.error('Error fetching modules:', error);
            throw error;
        }
    },

    // Get a single module by ID
    getModuleById: async (id) => {
        try {
            const response = await api.get(`/modules/${id}`);
            return response.data;
        } catch (error) {
            console.error(`Error fetching module with id ${id}:`, error);
            throw error;
        }
    },

    // Create a new module
    createModule: async (moduleData) => {
        try {
            const response = await api.post('/modules', moduleData);
            return response.data;
        } catch (error) {
            console.error('Error creating module:', error);
            throw error;
        }
    },

    // Update a module
    updateModule: async (id, moduleData) => {
        try {
            const response = await api.patch(`/modules/${id}`, moduleData);
            return response.data;
        } catch (error) {
            console.error(`Error updating module with id ${id}:`, error);
            throw error;
        }
    },

    // Delete a module
    deleteModule: async (id) => {
        try {
            const response = await api.delete(`/modules/${id}`);
            return response.data;
        } catch (error) {
            console.error(`Error deleting module with id ${id}:`, error);
            throw error;
        }
    }
};

export default moduleService;