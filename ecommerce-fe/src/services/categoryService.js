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
        if (!parentId) {
            return Promise.reject(new Error('Cần có parent_id để lấy danh mục con'));
        }

        return api.get('/categories', {
            params: {
                parent_id: parentId
            }
        });
    },

    /**
     * Lấy danh sách categories dựa trên từ khóa tìm kiếm sản phẩm
     * @param {String} keyword - Từ khóa tìm kiếm
     * @returns {Promise} - Promise chứa response từ API
     */
    getCategoriesByKeyword: (keyword) => {
        if (!keyword) {
            return Promise.reject(new Error('Cần có từ khóa để tìm kiếm danh mục'));
        }

        return api.get('/categories', {
            params: {
                'product_keyword': keyword
            }
        });
    },

    /**
     * Lấy thông tin chi tiết của một category
     * @param {Number} categoryId - ID của category
     * @returns {Promise} - Promise chứa response từ API
     */
    getCategoryById: (categoryId) => {
        if (!categoryId) {
            return Promise.reject(new Error('Cần có category ID'));
        }

        return api.get(`/categories/${categoryId}`);
    }
};

export default categoryService;