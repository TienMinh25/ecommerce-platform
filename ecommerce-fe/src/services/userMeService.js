import api from './api';

const userMeService = {
    /**
     * Get current user profile
     * @returns {Promise} Promise object with user data
     */
    getProfile: async () => {
        try {
            const response = await api.get('/users/me');
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Update user profile
     * @param {Object} userData - User data to update
     * @returns {Promise} Promise object with updated user data
     */
    updateProfile: async (userData) => {
        try {
            const response = await api.patch('/users/me', userData);
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Change user password
     * @param {Object} passwordData - Object containing old and new password
     * @returns {Promise} Promise object with result
     */
    changePassword: async (passwordData) => {
        try {
            const response = await api.post('/auth/change-password', passwordData);
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Get user addresses
     * @returns {Promise} Promise object with user addresses
     */
    getAddresses: async () => {
        try {
            const response = await api.get('/users/me/addresses');
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Add new address
     * @param {Object} addressData - Address data to add
     * @returns {Promise} Promise object with result
     */
    addAddress: async (addressData) => {
        try {
            const response = await api.post('/users/me/addresses', addressData);
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Update existing address
     * @param {number} addressId - ID of the address to update
     * @param {Object} addressData - Address data to update
     * @returns {Promise} Promise object with result
     */
    updateAddress: async (addressId, addressData) => {
        try {
            const response = await api.patch(`/users/me/addresses/${addressId}`, addressData);
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Delete address
     * @param {number} addressId - ID of the address to delete
     * @returns {Promise} Promise object with result
     */
    deleteAddress: async (addressId) => {
        try {
            const response = await api.delete(`/users/me/addresses/${addressId}`);
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Set address as default
     * @param {number} addressId - ID of the address to set as default
     * @returns {Promise} Promise object with result
     */
    setDefaultAddress: async (addressId) => {
        try {
            const response = await api.patch(`/users/me/addresses/${addressId}/default`);
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Upload avatar
     * @param {File} file - Avatar file to upload
     * @returns {Promise} Promise object with result
     */
    uploadAvatar: async (file) => {
        try {
            const formData = new FormData();
            formData.append('avatar', file);

            const response = await api.post('/users/me/avatar', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    getPresignedUrl: async (presignedRequest) => {
        try {
            const response = await api.post('/users/me/avatars/get-presigned-url', presignedRequest);
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Lấy cài đặt thông báo
     * @returns {Promise} Promise object với dữ liệu cài đặt thông báo
     */
    getNotificationSettings: async () => {
        try {
            const response = await api.get('/users/me/notification-settings');
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Cập nhật cài đặt thông báo
     * @param {Object} settings - Object chứa cài đặt email_setting và in_app_setting
     * @returns {Promise} Promise object với dữ liệu cài đặt thông báo đã cập nhật
     */
    updateNotificationSettings: async (settings) => {
        try {
            const response = await api.post('/users/me/notification-settings', settings);
            return response.data.data;
        } catch (error) {
            throw error;
        }
    }
};

export default userMeService;