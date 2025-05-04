import api from './api';

const categoryService = {
    /**
     * Lấy danh sách tất cả categories
     * @param {Object} params - Các tham số query (không bắt buộc)
     * @returns {Promise} - Promise chứa response từ API
     */
    getAllCategories: (params = {}) => {
        return api.get('/categories', { params });
    },

    /**
     * Lấy danh sách categories con của một category
     * @param {Number} parentId - ID của category cha
     * @returns {Promise} - Promise chứa response từ API
     */
    getSubCategories: (parentId) => {
        return api.get('/categories', { params: { parent_id: parentId } });
    },

    /**
     * Lấy thông tin chi tiết của một category
     * @param {Number} categoryId - ID của category
     * @returns {Promise} - Promise chứa response từ API
     */
    getCategoryById: (categoryId) => {
        return api.get(`/categories/${categoryId}`);
    }
};

export default categoryService;