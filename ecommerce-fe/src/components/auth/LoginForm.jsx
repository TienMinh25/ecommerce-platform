import {useState} from 'react';
import {Link as RouterLink, useLocation, useNavigate} from 'react-router-dom';
import {useForm} from 'react-hook-form';
import {yupResolver} from '@hookform/resolvers/yup';
import {
  Box,
  Button,
  Checkbox,
  Divider,
  Flex,
  FormControl,
  FormErrorMessage,
  FormLabel,
  HStack,
  Icon,
  IconButton,
  Input,
  InputGroup,
  InputRightElement,
  Link,
  Stack,
  Text,
  useToast,
} from '@chakra-ui/react';
import {ViewIcon, ViewOffIcon} from '@chakra-ui/icons';
import {MdEmail, MdVerified} from 'react-icons/md';
import {loginSchema} from '../../utils/validation';
import {useAuth} from '../../hooks/useAuth';
import SocialLogin from './SocialLogin';

const LoginForm = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [needVerification, setNeedVerification] = useState(false);
  const [loginEmail, setLoginEmail] = useState('');
  const [rememberMe, setRememberMe] = useState(true); // Default to checked
  const navigate = useNavigate();
  const location = useLocation();
  const toast = useToast();
  const { login, resendVerifyEmailOTP } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm({
    resolver: yupResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  const handleVerifyEmail = async () => {
    await resendVerifyEmailOTP({
      "email": loginEmail,
    })
    navigate('/verify-email', {
      replace: true,
      state: {
        email: loginEmail,
        isRegister: false
      }
    });
  };

  const onSubmit = async (data) => {
    setIsLoading(true);
    setNeedVerification(false); // Reset trạng thái xác thực

    try {
      // Lưu email để dùng cho nút xác thực
      setLoginEmail(data.email);

      const result = await login({
        ...data,
        remember: rememberMe // Thêm trạng thái remember me vào dữ liệu gửi đi
      });

      if (result.success) {
        toast({
          title: 'Đăng nhập thành công',
          status: 'success',
          duration: 3000,
          isClosable: true,
        });

        // Redirect to the page they tried to visit or to home
        const from = location.state?.from || '/';
        navigate(from, { replace: true });
      } else if (result.needVerification) {
        // Đặt trạng thái cần xác thực, không chuyển hướng ngay
        setNeedVerification(true);
      } else {
        throw new Error(result.error || 'Đăng nhập thất bại');
      }
    } catch (error) {
      toast({
        title: 'Đăng nhập thất bại',
        description: error.message,
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
      <Box as='form' onSubmit={handleSubmit(onSubmit)} w='full'>
        <Stack spacing={4}>
          <Box as='h1' fontSize='2xl' fontWeight='bold' mb={2} color="gray.800">
            Đăng nhập
          </Box>
          <Text color='gray.700' mb={6} fontWeight="medium">
            Chào mừng bạn quay trở lại!
          </Text>

          {needVerification && (
              <Box
                  p={5}
                  borderRadius="lg"
                  bg="orange.50"
                  borderWidth="1px"
                  borderColor="orange.200"
                  boxShadow="md"
                  mb={6}
                  position="relative"
                  overflow="hidden"
              >
                <Box
                    position="absolute"
                    top={0}
                    left={0}
                    height="full"
                    width="4px"
                    bg="orange.400"
                />

                <Flex direction="column" gap={3}>
                  <Flex align="center" gap={2}>
                    <Icon as={MdVerified} color="orange.500" boxSize={5} />
                    <Text fontSize="lg" fontWeight="bold" color="orange.800">
                      Tài khoản chưa được xác thực
                    </Text>
                  </Flex>

                  <Text color="gray.700" pl={7}>
                    Để đảm bảo an toàn cho tài khoản của bạn, vui lòng xác thực email <strong>{loginEmail}</strong> trước khi tiếp tục.
                  </Text>

                  <HStack spacing={4} mt={2} justify="flex-end">
                    <Button
                        variant="outline"
                        colorScheme="orange"
                        size="md"
                        onClick={() => setNeedVerification(false)}
                    >
                      Để sau
                    </Button>
                    <Button
                        leftIcon={<MdEmail />}
                        colorScheme="brand"
                        size="md"
                        fontWeight="bold"
                        boxShadow="sm"
                        _hover={{
                          transform: 'translateY(-1px)',
                          boxShadow: 'md',
                        }}
                        onClick={handleVerifyEmail}
                    >
                      Xác thực ngay
                    </Button>
                  </HStack>
                </Flex>
              </Box>
          )}

          <FormControl isInvalid={errors.email} id='email'>
            <FormLabel fontWeight="medium" color="gray.700">Email</FormLabel>
            <Input
                {...register('email')}
                type='email'
                placeholder='your@email.com'
                size='lg'
                bg="white"
                borderColor="gray.300"
                _hover={{ borderColor: 'brand.400' }}
                _focus={{
                  borderColor: 'brand.500',
                  boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)'
                }}
            />
            <FormErrorMessage fontWeight="medium">{errors.email?.message}</FormErrorMessage>
          </FormControl>

          <FormControl isInvalid={errors.password} id='password'>
            <FormLabel fontWeight="medium" color="gray.700">Mật khẩu</FormLabel>
            <InputGroup size='lg'>
              <Input
                  {...register('password')}
                  type={showPassword ? 'text' : 'password'}
                  placeholder='********'
                  bg="white"
                  borderColor="gray.300"
                  _hover={{ borderColor: 'brand.400' }}
                  _focus={{
                    borderColor: 'brand.500',
                    boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)'
                  }}
              />
              <InputRightElement>
                <IconButton
                    aria-label={showPassword ? 'Ẩn mật khẩu' : 'Hiện mật khẩu'}
                    icon={showPassword ? <ViewOffIcon /> : <ViewIcon />}
                    variant='ghost'
                    color="gray.500"
                    _hover={{ color: 'brand.500', bg: 'brand.50' }}
                    onClick={() => setShowPassword(!showPassword)}
                />
              </InputRightElement>
            </InputGroup>
            <FormErrorMessage fontWeight="medium">{errors.password?.message}</FormErrorMessage>
          </FormControl>

          <Stack
              direction={{ base: 'column', sm: 'row' }}
              justify='space-between'
              align='center'
              pt={2}
          >
            <Checkbox
                colorScheme='brand'
                isChecked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
            >
              <Text color="gray.700" fontWeight="medium">Ghi nhớ đăng nhập</Text>
            </Checkbox>
            <Link
                as={RouterLink}
                to='/forgot-password'
                color='brand.600'
                fontWeight='semibold'
                _hover={{ color: 'brand.700', textDecoration: 'underline' }}
            >
              Quên mật khẩu?
            </Link>
          </Stack>

          <Button
              type='submit'
              size='lg'
              colorScheme='brand'
              isLoading={isLoading}
              loadingText='Đang đăng nhập...'
              w='full'
              mt={6}
              fontWeight="bold"
              fontSize="md"
              py={6}
              boxShadow="md"
              _hover={{
                transform: 'translateY(-1px)',
                boxShadow: 'lg',
              }}
          >
            Đăng nhập
          </Button>

          <Divider my={6} borderColor="gray.300" />

          <SocialLogin />

          <Text mt={4} textAlign='center' color="gray.700">
            Chưa có tài khoản?{' '}
            <Link
                as={RouterLink}
                to='/register'
                color='brand.600'
                fontWeight='bold'
                _hover={{ color: 'brand.700', textDecoration: 'underline' }}
            >
              Đăng ký ngay
            </Link>
          </Text>
        </Stack>
      </Box>
  );
};

export default LoginForm;