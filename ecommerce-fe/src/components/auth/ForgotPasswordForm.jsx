import {useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {
    Box,
    Button,
    Flex,
    FormControl,
    FormErrorMessage,
    FormLabel,
    Heading,
    Icon,
    Input,
    Text,
    useToast,
    VStack,
} from '@chakra-ui/react';
import {ArrowBackIcon} from '@chakra-ui/icons';
import {MdEmail} from 'react-icons/md';
import {useForm} from 'react-hook-form';
import {yupResolver} from '@hookform/resolvers/yup';
import {useAuth} from '../../hooks/useAuth';
import {forgotPasswordSchema} from "../../utils/validation.js";

const ForgotPasswordForm = () => {
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();
    const toast = useToast();
    const {sendPasswordResetOTP} = useAuth();

    const {
        register,
        handleSubmit,
        formState: {errors},
    } = useForm({
        resolver: yupResolver(forgotPasswordSchema),
        defaultValues: {
            email: '',
        },
    });

    const onSubmit = async (data) => {
        setIsLoading(true);
        try {
            const result = await sendPasswordResetOTP({
                email: data.email,
            });

            if (result.success) {
                toast({
                    title: 'Đã gửi mã xác thực',
                    description: `Chúng tôi đã gửi mã xác thực đến email ${data.email}`,
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });

                // Chuyển đến trang nhập OTP và mật khẩu mới
                navigate('/reset-password', {
                    state: {
                        email: data.email,
                    },
                });
            } else {
                throw new Error(result.error || 'Không thể gửi mã xác thực, vui lòng thử lại sau');
            }
        } catch (error) {
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

    const handleBackToLogin = () => {
        navigate('/login');
    };

    return (
        <Box w="full" maxW="md" mx="auto" p={8} bg="white" borderRadius="xl" boxShadow="lg">
            <VStack spacing={6} align="stretch">
                <Button
                    leftIcon={<ArrowBackIcon/>}
                    variant="link"
                    color="gray.600"
                    alignSelf="flex-start"
                    mb={2}
                    onClick={handleBackToLogin}
                >
                    Quay lại đăng nhập
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
                        <Icon as={MdEmail} boxSize={10} color="brand.500"/>
                    </Flex>
                    <Heading size="lg" color="gray.800" textAlign="center">
                        Quên mật khẩu?
                    </Heading>
                    <Text color="gray.600" textAlign="center">
                        Nhập email của bạn để nhận mã xác thực đặt lại mật khẩu
                    </Text>
                </VStack>

                <Box as="form" onSubmit={handleSubmit(onSubmit)}>
                    <FormControl isInvalid={errors.email} mb={6}>
                        <FormLabel fontWeight="medium" color="gray.700">
                            Email
                        </FormLabel>
                        <Input
                            {...register('email')}
                            type="email"
                            placeholder="your@email.com"
                            size="lg"
                            bg="white"
                            borderColor="gray.300"
                            _hover={{borderColor: 'brand.400'}}
                            _focus={{
                                borderColor: 'brand.500',
                                boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)',
                            }}
                        />
                        <FormErrorMessage>{errors.email?.message}</FormErrorMessage>
                    </FormControl>

                    <Button
                        type="submit"
                        size="lg"
                        w="full"
                        colorScheme="brand"
                        isLoading={isLoading}
                        loadingText="Đang gửi..."
                        fontWeight="bold"
                        boxShadow="md"
                        _hover={{
                            transform: 'translateY(-1px)',
                            boxShadow: 'lg',
                        }}
                    >
                        Gửi mã xác thực
                    </Button>
                </Box>
            </VStack>
        </Box>
    );
};

export default ForgotPasswordForm;