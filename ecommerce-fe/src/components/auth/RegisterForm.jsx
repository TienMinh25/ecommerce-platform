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
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
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
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
      phone: '',
      agreeTerms: false,
    },
  });

  const onSubmit = async (data) => {
    setIsLoading(true);
    try {
      // TODO: Call register API
      const result = await signUp(data);

      if (result.success) {
        toast({
          title: 'Đăng ký thành công',
          description: 'Chào mừng bạn đến với Shop!',
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
        navigate('/', { replace: true });
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
        <Box as='h1' fontSize='2xl' fontWeight='bold' mb={2}>
          Tạo tài khoản mới
        </Box>
        <Text color='gray.600' mb={6}>
          Đăng ký để mua sắm và nhận nhiều ưu đãi hơn
        </Text>

        <FormControl isInvalid={errors.name} id='name'>
          <FormLabel>Họ và tên</FormLabel>
          <Input
            {...register('name')}
            type='text'
            placeholder='Nguyễn Văn A'
            size='lg'
          />
          <FormErrorMessage>{errors.name?.message}</FormErrorMessage>
        </FormControl>

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

        <FormControl isInvalid={errors.phone} id='phone'>
          <FormLabel>Số điện thoại</FormLabel>
          <Input
            {...register('phone')}
            type='tel'
            placeholder='0911234567'
            size='lg'
          />
          <FormErrorMessage>{errors.phone?.message}</FormErrorMessage>
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
                aria-label={showPassword ? 'Ẩn mật khẩu' : 'Hiện mật khẩu'}
                icon={showPassword ? <ViewOffIcon /> : <ViewIcon />}
                variant='ghost'
                onClick={() => setShowPassword(!showPassword)}
              />
            </InputRightElement>
          </InputGroup>
          <FormErrorMessage>{errors.password?.message}</FormErrorMessage>
        </FormControl>

        <FormControl isInvalid={errors.confirmPassword} id='confirmPassword'>
          <FormLabel>Xác nhận mật khẩu</FormLabel>
          <InputGroup size='lg'>
            <Input
              {...register('confirmPassword')}
              type={showConfirmPassword ? 'text' : 'password'}
              placeholder='********'
            />
            <InputRightElement>
              <IconButton
                aria-label={
                  showConfirmPassword ? 'Ẩn mật khẩu' : 'Hiện mật khẩu'
                }
                icon={showConfirmPassword ? <ViewOffIcon /> : <ViewIcon />}
                variant='ghost'
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              />
            </InputRightElement>
          </InputGroup>
          <FormErrorMessage>{errors.confirmPassword?.message}</FormErrorMessage>
        </FormControl>

        <FormControl isInvalid={errors.agreeTerms} id='agreeTerms'>
          <Stack direction='row' align='start' spacing={1} pt={4}>
            <Checkbox
              {...register('agreeTerms')}
              colorScheme='brand'
              size='lg'
            />
            <Text fontSize='sm' color='gray.600'>
              Tôi đồng ý với{' '}
              <Link color='brand.500' fontWeight='semibold'>
                Điều khoản sử dụng
              </Link>{' '}
              và{' '}
              <Link color='brand.500' fontWeight='semibold'>
                Chính sách bảo mật
              </Link>{' '}
              của Shop
            </Text>
          </Stack>
          <FormErrorMessage>{errors.agreeTerms?.message}</FormErrorMessage>
        </FormControl>

        <Button
          type='submit'
          size='lg'
          colorScheme='brand'
          isLoading={isLoading}
          loadingText='Đang đăng ký...'
          w='full'
          mt={6}
        >
          Đăng ký
        </Button>

        <Divider my={6} />

        <SocialLogin buttonText='Đăng ký với' />

        <Text mt={4} textAlign='center'>
          Đã có tài khoản?{' '}
          <Link
            as={RouterLink}
            to='/login'
            color='brand.500'
            fontWeight='semibold'
          >
            Đăng nhập
          </Link>
        </Text>
      </Stack>
    </Box>
  );
};

export default RegisterForm;
