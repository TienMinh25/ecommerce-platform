import { useState, useEffect } from 'react';
import {
    Button,
    Container,
    FormControl,
    FormErrorMessage,
    FormLabel,
    Heading,
    PinInput,
    PinInputField,
    HStack,
    Text,
    VStack,
    useToast,
    Icon,
    Flex
} from '@chakra-ui/react';
import { useLocation, useNavigate } from 'react-router-dom';
import { MdMarkEmailRead, MdLockOutline } from 'react-icons/md';
import useAuth from "../../hooks/useAuth.js";

const EmailVerification = () => {
    const [otp, setOtp] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');
    const {verifyOTP, resendVerifyEmailOTP} = useAuth()
    const toast = useToast();
    const navigate = useNavigate();
    const location = useLocation();

    // Lấy email và isRegister từ location state
    const email = location.state?.email || '';
    const isRegister = location.state?.isRegister || false;

    // Xóa lỗi khi người dùng thay đổi OTP
    useEffect(() => {
        if (error) {
            setError('');
        }
    }, [otp]);

    const handleVerify = async () => {
        if (otp.length !== 6) {
            setError('Vui lòng nhập đủ 6 chữ số');
            return;
        }

        setIsLoading(true);
        setError('');

        try {
            const result = await verifyOTP({
                "email": email,
                "otp": otp,
            })

            if (result.success) {
                toast({
                    title: 'Xác thực thành công',
                    description: 'Tài khoản của bạn đã được kích hoạt',
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });

                // Chuyển hướng đến trang đăng nhập hoặc trang chủ
                navigate('/login', { replace: true });
                return;
            }

            throw new Error(result.error)
        } catch (error) {
            setError(error.message);
        } finally {
            setIsLoading(false);
        }
    };

    const handleResendOtp = async () => {
        if (!email) {
            toast({
                title: 'Lỗi',
                description: 'Không tìm thấy địa chỉ email',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
            return;
        }

        try {
            const result = await resendVerifyEmailOTP({
                "email": email
            })

            if (result.success) {
                toast({
                    title: 'Đã gửi lại mã',
                    description: `Mã xác thực đã được gửi tới ${email}`,
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });

                // Reset OTP input và lỗi khi gửi lại mã
                setOtp('');
                setError('');
            }
        } catch (error) {
            toast({
                title: 'Lỗi',
                description: 'Không thể gửi lại mã xác thực. Vui lòng thử lại sau.',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    // Xử lý thay đổi OTP và xóa lỗi
    const handleOtpChange = (value) => {
        setOtp(value);
        // useEffect sẽ xử lý việc xóa lỗi
    };

    if (!email) {
        return (
            <Container maxW="lg" py={12}>
                <VStack spacing={8} align="center">
                    <Heading>Có lỗi xảy ra</Heading>
                    <Text>Không tìm thấy thông tin email. Vui lòng quay lại trang đăng ký.</Text>
                    <Button colorScheme="brand" onClick={() => navigate('/register')}>
                        Quay lại đăng ký
                    </Button>
                </VStack>
            </Container>
        );
    }

    // Nội dung tiêu đề và icon phụ thuộc vào isRegister
    const renderContent = () => {
        if (isRegister) {
            return {
                icon: MdMarkEmailRead,
                title: "Chúc mừng bạn đã đăng ký thành công!",
                description: "Vui lòng nhập mã xác thực đã được gửi đến"
            };
        } else {
            return {
                icon: MdLockOutline,
                title: "Xác thực tài khoản",
                description: "Để tiếp tục sử dụng tài khoản, vui lòng nhập mã xác thực đã được gửi đến"
            };
        }
    };

    const content = renderContent();

    return (
        <Container maxW="lg" py={12}>
            <VStack spacing={8} align="center">
                <Flex
                    w="120px"
                    h="120px"
                    borderRadius="full"
                    bg="brand.50"
                    justify="center"
                    align="center"
                >
                    <Icon as={content.icon} w={16} h={16} color="brand.500" />
                </Flex>

                <Heading textAlign="center">{content.title}</Heading>

                <Text textAlign="center" fontSize="lg">
                    {content.description}
                </Text>

                <Text fontWeight="bold" fontSize="xl" color="brand.500">
                    {email}
                </Text>

                <FormControl isInvalid={!!error}>
                    <FormLabel textAlign="center">Nhập mã OTP gồm 6 chữ số</FormLabel>
                    <HStack justify="center" spacing="4">
                        <PinInput
                            size="lg"
                            value={otp}
                            onChange={handleOtpChange}
                            type="number"
                        >
                            <PinInputField />
                            <PinInputField />
                            <PinInputField />
                            <PinInputField />
                            <PinInputField />
                            <PinInputField />
                        </PinInput>
                    </HStack>
                    {error && <FormErrorMessage textAlign="center">{error}</FormErrorMessage>}
                </FormControl>

                <VStack spacing={4} w="full">
                    <Button
                        onClick={handleVerify}
                        isLoading={isLoading}
                        loadingText="Đang xác thực..."
                        colorScheme="brand"
                        size="lg"
                        w="full"
                    >
                        Xác thực
                    </Button>

                    <Text>
                        Chưa nhận được mã?{' '}
                        <Button
                            variant="link"
                            colorScheme="brand"
                            onClick={handleResendOtp}
                            isDisabled={isLoading}
                        >
                            Gửi lại
                        </Button>
                    </Text>
                </VStack>
            </VStack>
        </Container>
    );
};

export default EmailVerification;