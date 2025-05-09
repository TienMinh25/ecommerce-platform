import api from './api';
import axios from 'axios';

const productService = {
    /**
     * Lấy danh sách sản phẩm với các tiêu chí cụ thể, sử dụng định dạng API chính xác
     * Hàm này sử dụng format đúng cho category_ids theo yêu cầu API
     * @param {Object} options - Các tùy chọn lọc
     * @returns {Promise} - Promise chứa response từ API
     */
    getProductsByCriteria: (options = {}) => {
        const {
            limit = 20,
            page = 1,
            keyword = null,
            categoryIds = null,
            minRating = null
        } = options;

        // Sử dụng URLSearchParams để xây dựng query string đúng cách
        const params = new URLSearchParams();

        // Thêm tham số cơ bản
        params.append('limit', limit);
        params.append('page', page);

        // Thêm keyword nếu có
        if (keyword) {
            params.append('keyword', keyword);
        }

        // Xử lý category_ids - API cần định dạng category_ids=1&category_ids=2
        if (categoryIds) {
            if (Array.isArray(categoryIds)) {
                // Nếu là mảng, thêm mỗi ID như một tham số riêng
                categoryIds.forEach(id => {
                    params.append('category_ids', id);
                });
            } else if (typeof categoryIds === 'string' && categoryIds.includes(',')) {
                // Nếu là chuỗi phân cách bằng dấu phẩy, chuyển thành mảng và thêm từng ID
                categoryIds.split(',').forEach(id => {
                    params.append('category_ids', id.trim());
                });
            } else {
                // Nếu là một ID đơn
                params.append('category_ids', categoryIds);
            }
        }

        // Thêm min_rating nếu có
        if (minRating) {
            params.append('min_rating', minRating);
        }

        // Sử dụng Axios trực tiếp để giữ đúng định dạng query string
        return api.get(`${api.defaults.baseURL}/products`, {
            params,
            paramsSerializer: params => params.toString(),
            headers: api.defaults.headers
        });
    },

    /**
     * Lấy danh sách sản phẩm nổi bật
     * @param {Number} limit - Số lượng sản phẩm muốn lấy
     * @returns {Promise} - Promise chứa response từ API
     */
    getFeaturedProducts: (limit = 24) => {
        return api.get('/products', {
            params: {
                limit,
                page: 1
            }
        });
    },

    /**
     * Get product details by ID
     * @param {String|Number} productId - Product ID
     * @returns {Promise} - Promise containing API response
     */
    getProductById: (productId) => {
        return api.get(`/products/${productId}`);
    },

    /**
     * Get product reviews by product ID
     * @param {String|Number} productId - Product ID
     * @param {Number} page - Current page
     * @param {Number} limit - Reviews per page
     * @returns {Promise} - Promise containing API response
     */
    getProductReviews: (productId, page = 1, limit = 6) => {
        return api.get(`/products/${productId}/reviews`, {
            params: {
                page,
                limit
            }
        });
    }
};

export default productService;