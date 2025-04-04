import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { Center, Spinner, Text } from '@chakra-ui/react';

const PublicRoute = () => {
  const { user, isLoading, authCheckComplete } = useAuth();
  const location = useLocation();

  // Lấy đường dẫn chuyển hướng từ state (nếu có)
  const from = location.state?.from || '/';

  // Nếu đang kiểm tra trạng thái xác thực, hiển thị loading
  if (isLoading || !authCheckComplete) {
    return (
        <Center flexDirection="column" h="100vh" bg="gray.50">
          <Spinner
              thickness="4px"
              speed="0.65s"
              emptyColor="gray.200"
              color="brand.500"
              size="xl"
              mb={4}
          />
          <Text fontWeight="medium" color="gray.600">
            Đang xác thực...
          </Text>
        </Center>
    );
  }

  // Nếu người dùng đã đăng nhập và đang cố gắng truy cập trang công khai như đăng nhập/đăng ký
  // Chuyển hướng họ đến trang chính hoặc trang họ đến từ đó
  if (user !== null) {
    // Không chuyển hướng nếu người dùng đang truy cập OAuthCallbackPage
    if (location.pathname === '/oauth') {
      return <Outlet />;
    }

    return <Navigate to={from} replace />;
  }

  return <Outlet />;
};

export default PublicRoute;