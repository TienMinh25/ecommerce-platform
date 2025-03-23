import { useState } from 'react';
import { Link as RouterLink, useNavigate, useLocation } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  FormErrorMessage,
  Input,
  Stack,
  Checkbox,
  Text,
  Link,
  Divider,
  useToast,
  InputGroup,
  InputRightElement,
  IconButton,
} from '@chakra-ui/react';
import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons';
import { loginSchema } from '../../utils/validation';
import { useAuth } from '../../hooks/useAuth';
import SocialLogin from './SocialLogin';

const LoginForm = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const toast = useToast();
  const { login } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: yupResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  const onSubmit = async (data) => {
    setIsLoading(true);
    try {
      // TODO: Call login API
      const result = await login(data);

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
        <Box as='h1' fontSize='2xl' fontWeight='bold' mb={2}>
          Đăng nhập
        </Box>
        <Text color='gray.600' mb={6}>
          Chào mừng bạn quay trở lại!
        </Text>

        <FormControl isInvalid={errors.email} id='email'>
          <FormLabel>Email</FormLabel>
          <Input
            {...register('email')}
            type='email'
            placeholder='your@email.com'
            size='lg'
          />
          <FormErrorMessage>{errors.email?.message}</FormErrorMessage>
        </FormControl>

        <FormControl isInvalid={errors.password} id='password'>
          <FormLabel>Mật khẩu</FormLabel>
          <InputGroup size='lg'>
            <Input
              {...register('password')}
              type={showPassword ? 'text' : 'password'}
              placeholder='********'
            />
            <InputRightElement>
              <IconButton
                aria-label={showPassword ? 'Hide password' : 'Show password'}
                icon={showPassword ? <ViewOffIcon /> : <ViewIcon />}
                variant='ghost'
                onClick={() => setShowPassword(!showPassword)}
              />
            </InputRightElement>
          </InputGroup>
          <FormErrorMessage>{errors.password?.message}</FormErrorMessage>
        </FormControl>

        <Stack
          direction={{ base: 'column', sm: 'row' }}
          justify='space-between'
          align='center'
          pt={2}
        >
          <Checkbox colorScheme='brand'>Ghi nhớ đăng nhập</Checkbox>
          <Link color='brand.500' fontWeight='semibold'>
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
        >
          Đăng nhập
        </Button>

        <Divider my={6} />

        <SocialLogin />

        <Text mt={4} textAlign='center'>
          Chưa có tài khoản?{' '}
          <Link
            as={RouterLink}
            to='/register'
            color='brand.500'
            fontWeight='semibold'
          >
            Đăng ký ngay
          </Link>
        </Text>
      </Stack>
    </Box>
  );
};

export default LoginForm;
