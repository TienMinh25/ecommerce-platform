import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { useEffect, useState } from 'react';
import { Box, Center, Spinner, Text } from '@chakra-ui/react';

const PrivateRoute = () => {
  const { user, isLoading, authCheckComplete } = useAuth();
  const location = useLocation();
  const [showTimeout, setShowTimeout] = useState(false);

  // Hiển thị thông báo nếu việc kiểm tra xác thực mất quá nhiều thời gian
  useEffect(() => {
    let timeoutId;

    if (isLoading) {
      timeoutId = setTimeout(() => {
        setShowTimeout(true);
      }, 5000); // 5 giây
    }

    return () => {
      if (timeoutId) clearTimeout(timeoutId);
    };
  }, [isLoading]);

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

          {showTimeout && (
              <Box mt={4} p={4} bg="yellow.50" borderRadius="md" maxW="md" textAlign="center">
                <Text color="yellow.700">
                  Quá trình xác thực đang mất nhiều thời gian hơn dự kiến.
                  Vui lòng đợi hoặc thử lại sau.
                </Text>
              </Box>
          )}
        </Center>
    );
  }

  // Sau khi kiểm tra xong, nếu không có user thì chuyển hướng
  if (user === null) {
    return <Navigate to='/login' state={{ from: location.pathname }} replace />;
  }

  return <Outlet />;
};

export default PrivateRoute;