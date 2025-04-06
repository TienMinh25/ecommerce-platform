import api from './api.js'; // Import the axios instance

const addressTypeService = {
    // Get list of address types with pagination
    getAddressTypes: async (page = 1, limit = 10) => {
        try {
            const response = await api.get('/address-types', {
                params: { page, limit }
            });
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Get a single address type by ID
    getAddressTypeById: async (id) => {
        try {
            const response = await api.get(`/address-types/${id}`);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Create a new address type
    createAddressType: async (addressTypeData) => {
        try {
            const response = await api.post('/address-types', addressTypeData);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Update an address type
    updateAddressType: async (id, addressTypeData) => {
        try {
            const response = await api.patch(`/address-types/${id}`, addressTypeData);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    // Delete an address type
    deleteAddressType: async (id) => {
        try {
            const response = await api.delete(`/address-types/${id}`);
            return response.data;
        } catch (error) {
            throw error;
        }
    }
};

export default addressTypeService;