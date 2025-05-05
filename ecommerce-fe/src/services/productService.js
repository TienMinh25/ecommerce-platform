import api from './api';
import axios from 'axios';

const productService = {
    /**
     * Lấy danh sách sản phẩm với các tùy chọn lọc
     * @param {Object} options - Các tùy chọn lọc
     * @returns {Promise} - Promise chứa response từ API
     */
    getProducts: (options = {}) => {
        const {
            limit = 20,
            page = 1,
            keyword = null,
            categoryIds = null,
            minRating = null
        } = options;

        const params = {};
        params.limit = limit;
        params.page = page;

        if (keyword) params.keyword = keyword;
        if (categoryIds) {
            if (Array.isArray(categoryIds)) {
                params.category_ids = categoryIds.join(',');
            } else {
                params.category_ids = categoryIds;
            }
        }
        if (minRating) params.min_rating = minRating;

        return api.get('/products', { params });
    },

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
     * Lấy sản phẩm theo danh mục
     * @param {Number|Array} categoryIds - ID danh mục hoặc mảng ID danh mục
     * @param {Number} page - Trang hiện tại
     * @param {Number} limit - Số lượng sản phẩm mỗi trang
     * @returns {Promise} - Promise chứa response từ API
     */
    getProductsByCategory: (categoryIds, page = 1, limit = 20) => {
        // Sử dụng hàm getProductsByCriteria để đảm bảo định dạng chính xác
        return productService.getProductsByCriteria({
            page,
            limit,
            categoryIds
        });
    },

    /**
     * Tìm kiếm sản phẩm theo từ khóa
     * @param {String} keyword - Từ khóa tìm kiếm
     * @param {Number} page - Trang hiện tại
     * @param {Number} limit - Số lượng sản phẩm mỗi trang
     * @returns {Promise} - Promise chứa response từ API
     */
    searchProducts: (keyword, page = 1, limit = 20) => {
        return productService.getProductsByCriteria({
            page,
            limit,
            keyword
        });
    },

    /**
     * Lấy chi tiết sản phẩm theo ID
     * @param {String|Number} productId - ID sản phẩm
     * @returns {Promise} - Promise chứa response từ API
     */
    getProductById: (productId) => {
        return api.get(`/products/${productId}`);
    }
};

export default productService;