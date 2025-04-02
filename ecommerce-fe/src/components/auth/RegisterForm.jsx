import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons';
import {
  Box,
  Button,
  Checkbox,
  Divider,
  FormControl,
  FormErrorMessage,
  FormLabel,
  IconButton,
  Input,
  InputGroup,
  InputRightElement,
  Link,
  Stack,
  Text,
  useToast,
} from '@chakra-ui/react';
import { yupResolver } from '@hookform/resolvers/yup';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { registerSchema } from '../../utils/validation';
import SocialLogin from './SocialLogin';

const RegisterForm = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();
  const toast = useToast();
  const { register: signUp } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: yupResolver(registerSchema),
    defaultValues: {
      email: '',
      password: '',
      full_name: '',
      agreeTerms: false,
    },
  });

  const onSubmit = async (data) => {
    setIsLoading(true);
    try {
      const result = await signUp(data);

      if (result.success) {
        toast({
          title: 'Đăng ký thành công',
          description: 'Vui lòng xác thực email của bạn',
          status: 'success',
          duration: 3000,
          isClosable: true,
        });

        // Chuyển hướng đến trang xác thực email và truyền email qua location state
        navigate('/verify-email', {
          replace: true,
          state: { email: data.email, isRegister: true }
        });
      } else {
        throw new Error(result.error || 'Đăng ký thất bại');
      }
    } catch (error) {
      toast({
        title: 'Đăng ký thất bại',
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
            Tạo tài khoản mới
          </Box>
          <Text color='gray.700' mb={6} fontWeight="medium">
            Đăng ký để mua sắm và nhận nhiều ưu đãi hơn
          </Text>

          <FormControl isInvalid={errors.name} id='name'>
            <FormLabel fontWeight="medium" color="gray.700">Họ và tên</FormLabel>
            <Input
                {...register('full_name')}
                type='text'
                placeholder='Nguyễn Văn A'
                size='lg'
                bg="white"
                borderColor="gray.300"
                _hover={{ borderColor: 'brand.400' }}
                _focus={{
                  borderColor: 'brand.500',
                  boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)'
                }}
            />
            <FormErrorMessage fontWeight="medium">{errors.name?.message}</FormErrorMessage>
          </FormControl>

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

          <FormControl isInvalid={errors.agreeTerms} id='agreeTerms'>
            <Stack direction='row' align='start' spacing={1} pt={4}>
              <Checkbox
                  {...register('agreeTerms')}
                  colorScheme='brand'
                  size='lg'
              />
              <Text fontSize='sm' color='gray.700' fontWeight="medium">
                Tôi đồng ý với{' '}
                <Link
                    color='brand.600'
                    fontWeight='bold'
                    _hover={{ color: 'brand.700', textDecoration: 'underline' }}
                >
                  Điều khoản sử dụng
                </Link>{' '}
                và{' '}
                <Link
                    color='brand.600'
                    fontWeight='bold'
                    _hover={{ color: 'brand.700', textDecoration: 'underline' }}
                >
                  Chính sách bảo mật
                </Link>{' '}
                của Shop
              </Text>
            </Stack>
            <FormErrorMessage fontWeight="medium">{errors.agreeTerms?.message}</FormErrorMessage>
          </FormControl>

          <Button
              type='submit'
              size='lg'
              colorScheme='brand'
              isLoading={isLoading}
              loadingText='Đang đăng ký...'
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
            Đăng ký
          </Button>

          <Divider my={6} borderColor="gray.300" />

          <SocialLogin buttonText='Đăng ký với' />

          <Text mt={4} textAlign='center' color="gray.700">
            Đã có tài khoản?{' '}
            <Link
                as={RouterLink}
                to='/login'
                color='brand.600'
                fontWeight='bold'
                _hover={{ color: 'brand.700', textDecoration: 'underline' }}
            >
              Đăng nhập
            </Link>
          </Text>
        </Stack>
      </Box>
  );
};

export default RegisterForm;