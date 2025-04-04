import { useEffect, useState, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { Box, Spinner, Text, Center, VStack, Alert, AlertIcon } from '@chakra-ui/react';

const OAuthCallbackPage = () => {
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const location = useLocation();
    const navigate = useNavigate();
    const { socialLogin } = useAuth();
    const processedRef = useRef(false);

    useEffect(() => {
        // Chỉ xử lý một lần khi component mount
        if (processedRef.current) return;

        const handleOAuthCallback = async () => {
            if (processedRef.current) return;
            processedRef.current = true;

            try {
                // Lấy code và state từ URL
                const searchParams = new URLSearchParams(location.search);
                const code = searchParams.get('code');
                const state = searchParams.get('state');

                if (!code || !state) {
                    throw new Error('Thiếu thông tin xác thực từ nhà cung cấp OAuth');
                }

                // Lấy oauth_provider từ localStorage
                const provider = localStorage.getItem('oauth_provider');

                if (!provider) {
                    return () => {
                        processedRef.current = true;
                    }
                }

                console.log('Bắt đầu xác thực OAuth với:', { provider, code, state });

                // Thực hiện đăng nhập với OAuth
                const result = await socialLogin(code, state, provider);

                if (result.success) {
                    // Chuyển hướng đến trang chủ sau khi đăng nhập thành công
                    navigate('/');
                } else {
                    setError(result.error || 'Đăng nhập thất bại');
                }
            } catch (err) {
                console.error('OAuth callback error:', err);
                setError(err.message || 'Đã xảy ra lỗi trong quá trình xác thực');
            } finally {
                setIsLoading(false);
                // Xóa provider từ localStorage sau khi xử lý
                localStorage.removeItem('oauth_provider');
            }
        };

        handleOAuthCallback();

        // Hàm cleanup để đảm bảo không có side effect
        return () => {
            processedRef.current = true;
        };
    }, []); // Chỉ chạy một lần khi component mount, không phụ thuộc vào location, navigate hoặc socialLogin

    // Lấy provider từ localStorage để hiển thị (chỉ UI, không ảnh hưởng xử lý)
    const providerName = localStorage.getItem('oauth_provider') || 'OAuth';

    return (
        <Center minH="100vh" bg="gray.50">
            <VStack spacing={6} p={8} bg="white" boxShadow="md" borderRadius="md" width="100%" maxW="md">
                <Text fontSize="2xl" fontWeight="bold">
                    Đăng nhập với {providerName}
                </Text>

                {isLoading ? (
                    <VStack spacing={4}>
                        <Spinner size="xl" thickness="4px" speed="0.65s" color="blue.500" />
                        <Text>Đang xác thực, vui lòng đợi...</Text>
                    </VStack>
                ) : error ? (
                    <Alert status="error" borderRadius="md">
                        <AlertIcon />
                        {error}
                    </Alert>
                ) : (
                    <Box textAlign="center">
                        <Text color="green.500" fontSize="lg">
                            Đăng nhập thành công! Đang chuyển hướng...
                        </Text>
                    </Box>
                )}
            </VStack>
        </Center>
    );
};

export default OAuthCallbackPage;