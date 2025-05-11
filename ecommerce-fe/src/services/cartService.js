import api from './api';

const cartService = {
    /**
     * Lấy danh sách sản phẩm trong giỏ hàng
     * @returns {Promise} Promise chứa response từ API
     */
    getCartItems: async () => {
        try {
            const response = await api.get('/users/me/carts');
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Thêm sản phẩm vào giỏ hàng
     * @param {Object} cartItem - Dữ liệu sản phẩm cần thêm (product_id, product_variant_id, quantity)
     * @returns {Promise} Promise chứa response từ API
     */
    addToCart: async (cartItem) => {
        try {
            const response = await api.post('/users/me/carts', cartItem);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Cập nhật số lượng sản phẩm trong giỏ hàng
     * @param {String} cartItemId - ID của item trong giỏ hàng
     * @param {Object} data - Dữ liệu cần cập nhật (product_variant_id, quantity)
     * @returns {Promise} Promise chứa response từ API
     */
    updateCartItem: async (cartItemId, data) => {
        try {
            const response = await api.patch(`/users/me/carts/${cartItemId}`, data);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Xóa sản phẩm khỏi giỏ hàng
     * @param {Array} cartItemIds - Mảng các ID của item cần xóa
     * @returns {Promise} Promise chứa response từ API
     */
    deleteCartItems: async (cartItemIds) => {
        try {
            const response = await api.delete('/users/me/carts', {
                data: { cart_item_ids: cartItemIds }
            });
            return response.data;
        } catch (error) {
            throw error;
        }
    }
};

export default cartService;