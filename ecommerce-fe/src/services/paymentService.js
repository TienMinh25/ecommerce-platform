import api from './api';

const paymentService = {
    /**
     * Get all available payment methods
     * @returns {Promise} Promise object with payment methods data
     */
    getPaymentMethods: async () => {
        try {
            const response = await api.get('/payments/payment-methods');
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Process payment
     * @param {Object} paymentData - Payment information
     * @returns {Promise} Promise object with payment result
     */
    processPayment: async (paymentData) => {
        try {
            const response = await api.post('/payments/process', paymentData);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Create order with payment
     * @param {Object} orderData - Order information including payment method
     * @returns {Promise} Promise object with order result
     */
    createOrder: async (orderData) => {
        try {
            const response = await api.post('/payments/checkout', orderData);
            return response.data;
        } catch (error) {
            throw error;
        }
    }
};

export default paymentService;