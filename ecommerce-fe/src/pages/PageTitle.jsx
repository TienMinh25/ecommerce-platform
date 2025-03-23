import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';

const PageTitle = ({ title, suffix = 'Minh Plaza' }) => {
  const location = useLocation();

  useEffect(() => {
    let pageTitle = title;

    if (!pageTitle) {
      const path = location.pathname;

      // Map các path sang tiêu đề
      const pageTitles = {
        '/': 'Trang chủ',
        '/products': 'Sản phẩm',
        '/cart': 'Giỏ hàng',
        '/checkout': 'Thanh toán',
        '/account': 'Tài khoản',
        '/login': 'Đăng nhập',
        '/register': 'Đăng ký',
        '/about': 'Về chúng tôi',
        '/contact': 'Liên hệ',
      };

      pageTitle = pageTitles[path] || 'Mua sắm trực tuyến';

      if (path.startsWith('/product/')) {
        const productName = path.split('/').pop().replace(/-/g, ' ');
        pageTitle = productName.charAt(0).toUpperCase() + productName.slice(1);
      }
    }

    document.title = pageTitle ? `${pageTitle} | ${suffix}` : suffix;
  }, [title, suffix, location]);

  return null;
};

export default PageTitle;
