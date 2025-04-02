import { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
    Box,
    Button,
    FormControl,
    FormLabel,
    FormErrorMessage,
    Input,
    VStack,
    Text,
    Heading,
    useToast,
    Icon,
    Flex,
    HStack,
    PinInput,
    PinInputField,
    InputGroup,
    InputRightElement,
    IconButton,
} from '@chakra-ui/react';
import { ArrowBackIcon, ViewIcon, ViewOffIcon } from '@chakra-ui/icons';
import { MdLockReset } from 'react-icons/md';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useAuth } from '../../hooks/useAuth';
import {resetPasswordSchema} from "../../utils/validation.js";

const ResetPasswordForm = () => {
    const [otp, setOtp] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');
    const location = useLocation();
    const navigate = useNavigate();
    const toast = useToast();
    const { resetPassword, sendPasswordResetOTP } = useAuth();

    // Lấy email từ state của location
    const email = location.state?.email || '';

    // Xóa lỗi khi người dùng thay đổi OTP
    useEffect(() => {
        if (error) {
            setError('');
        }
    }, [otp]);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(resetPasswordSchema),
        defaultValues: {
            password: '',
            confirmPassword: '',
        },
    });

    const handleOtpChange = (value) => {
        setOtp(value);
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
            const result = await sendPasswordResetOTP({
                email: email,
            });

            if (result.success) {
                toast({
                    title: 'Đã gửi lại mã',
                    description: `Mã xác thực đã được gửi lại tới ${email}`,
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });
                setOtp('');
            } else {
                throw new Error(result.error || 'Không thể gửi lại mã xác thực');
            }
        } catch (error) {
            toast({
                title: 'Lỗi',
                description: error.message || 'Không thể gửi lại mã xác thực. Vui lòng thử lại sau.',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    const onSubmit = async (data) => {
        if (otp.length !== 6) {
            setError('Vui lòng nhập đủ 6 chữ số của mã xác thực');
            return;
        }

        setIsLoading(true);
        try {
            const result = await resetPassword({
                email: email,
                otp: otp,
                password: data.password,
            });

            if (result.success) {
                toast({
                    title: 'Đổi mật khẩu thành công',
                    description: 'Mật khẩu của bạn đã được cập nhật. Vui lòng đăng nhập lại.',
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });
                navigate('/login', { replace: true });
            } else {
                throw new Error(result.error || 'Không thể đặt lại mật khẩu');
            }
        } catch (error) {
            setError(error.message);
            toast({
                title: 'Lỗi',
                description: error.message,
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    };

    const handleBackToForgotPassword = () => {
        navigate('/forgot-password');
    };

    // Kiểm tra nếu không có email, quay về trang forgot password
    if (!email) {
        useEffect(() => {
            toast({
                title: 'Lỗi',
                description: 'Không tìm thấy thông tin email. Vui lòng thử lại.',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
            navigate('/forgot-password', { replace: true });
        }, []);
        return null;
    }

    return (
        <Box w="full" maxW="md" mx="auto" p={8} bg="white" borderRadius="xl" boxShadow="lg">
            <VStack spacing={6} align="stretch">
                <Button
                    leftIcon={<ArrowBackIcon />}
                    variant="link"
                    color="gray.600"
                    alignSelf="flex-start"
                    mb={2}
                    onClick={handleBackToForgotPassword}
                >
                    Quay lại
                </Button>

                <VStack spacing={2} align="center" mb={4}>
                    <Flex
                        w="80px"
                        h="80px"
                        borderRadius="full"
                        bg="brand.50"
                        justify="center"
                        align="center"
                        mb={2}
                    >
                        <Icon as={MdLockReset} boxSize={10} color="brand.500" />
                    </Flex>
                    <Heading size="lg" color="gray.800" textAlign="center">
                        Đặt lại mật khẩu
                    </Heading>
                    <Text color="gray.600" textAlign="center">
                        Nhập mã xác thực đã được gửi đến <strong>{email}</strong> và mật khẩu mới của bạn
                    </Text>
                </VStack>

                <Box as="form" onSubmit={handleSubmit(onSubmit)}>
                    <FormControl isInvalid={!!error} mb={6}>
                        <FormLabel fontWeight="medium" color="gray.700" textAlign="center">
                            Mã xác thực gồm 6 chữ số
                        </FormLabel>
                        <HStack justify="center" spacing={4} mb={2}>
                            <PinInput size="lg" value={otp} onChange={handleOtpChange} type="number">
                                <PinInputField
                                    bg="white"
                                    borderColor="gray.300"
                                    _hover={{ borderColor: 'brand.400' }}
                                    _focus={{
                                        borderColor: 'brand.500',
                                        boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)',
                                    }}
                                />
                                <PinInputField
                                    bg="white"
                                    borderColor="gray.300"
                                    _hover={{ borderColor: 'brand.400' }}
                                    _focus={{
                                        borderColor: 'brand.500',
                                        boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)',
                                    }}
                                />
                                <PinInputField
                                    bg="white"
                                    borderColor="gray.300"
                                    _hover={{ borderColor: 'brand.400' }}
                                    _focus={{
                                        borderColor: 'brand.500',
                                        boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)',
                                    }}
                                />
                                <PinInputField
                                    bg="white"
                                    borderColor="gray.300"
                                    _hover={{ borderColor: 'brand.400' }}
                                    _focus={{
                                        borderColor: 'brand.500',
                                        boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)',
                                    }}
                                />
                                <PinInputField
                                    bg="white"
                                    borderColor="gray.300"
                                    _hover={{ borderColor: 'brand.400' }}
                                    _focus={{
                                        borderColor: 'brand.500',
                                        boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)',
                                    }}
                                />
                                <PinInputField
                                    bg="white"
                                    borderColor="gray.300"
                                    _hover={{ borderColor: 'brand.400' }}
                                    _focus={{
                                        borderColor: 'brand.500',
                                        boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)',
                                    }}
                                />
                            </PinInput>
                        </HStack>
                        {error && <FormErrorMessage textAlign="center">{error}</FormErrorMessage>}
                        <Text fontSize="sm" color="gray.600" textAlign="center" mt={2}>
                            Chưa nhận được mã?{' '}
                            <Button
                                variant="link"
                                colorScheme="brand"
                                size="sm"
                                onClick={handleResendOtp}
                                fontWeight="bold"
                            >
                                Gửi lại
                            </Button>
                        </Text>
                    </FormControl>

                    <FormControl isInvalid={errors.password} mb={4}>
                        <FormLabel fontWeight="medium" color="gray.700">
                            Mật khẩu mới
                        </FormLabel>
                        <InputGroup>
                            <Input
                                {...register('password')}
                                type={showPassword ? 'text' : 'password'}
                                placeholder="Nhập mật khẩu mới"
                                size="lg"
                                bg="white"
                                borderColor="gray.300"
                                _hover={{ borderColor: 'brand.400' }}
                                _focus={{
                                    borderColor: 'brand.500',
                                    boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)',
                                }}
                            />
                            <InputRightElement height="100%">
                                <IconButton
                                    aria-label={showPassword ? 'Ẩn mật khẩu' : 'Hiện mật khẩu'}
                                    icon={showPassword ? <ViewOffIcon /> : <ViewIcon />}
                                    variant="ghost"
                                    color="gray.500"
                                    _hover={{ color: 'brand.500', bg: 'brand.50' }}
                                    onClick={() => setShowPassword(!showPassword)}
                                />
                            </InputRightElement>
                        </InputGroup>
                        <FormErrorMessage>{errors.password?.message}</FormErrorMessage>
                    </FormControl>

                    <FormControl isInvalid={errors.confirmPassword} mb={6}>
                        <FormLabel fontWeight="medium" color="gray.700">
                            Xác nhận mật khẩu
                        </FormLabel>
                        <InputGroup>
                            <Input
                                {...register('confirmPassword')}
                                type={showConfirmPassword ? 'text' : 'password'}
                                placeholder="Nhập lại mật khẩu mới"
                                size="lg"
                                bg="white"
                                borderColor="gray.300"
                                _hover={{ borderColor: 'brand.400' }}
                                _focus={{
                                    borderColor: 'brand.500',
                                    boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)',
                                }}
                            />
                            <InputRightElement height="100%">
                                <IconButton
                                    aria-label={showConfirmPassword ? 'Ẩn mật khẩu' : 'Hiện mật khẩu'}
                                    icon={showConfirmPassword ? <ViewOffIcon /> : <ViewIcon />}
                                    variant="ghost"
                                    color="gray.500"
                                    _hover={{ color: 'brand.500', bg: 'brand.50' }}
                                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                />
                            </InputRightElement>
                        </InputGroup>
                        <FormErrorMessage>{errors.confirmPassword?.message}</FormErrorMessage>
                    </FormControl>

                    <Button
                        type="submit"
                        size="lg"
                        w="full"
                        colorScheme="brand"
                        isLoading={isLoading}
                        loadingText="Đang cập nhật..."
                        fontWeight="bold"
                        boxShadow="md"
                        _hover={{
                            transform: 'translateY(-1px)',
                            boxShadow: 'lg',
                        }}
                    >
                        Cập nhật mật khẩu
                    </Button>
                </Box>
            </VStack>
        </Box>
    );
};

export default ResetPasswordForm;