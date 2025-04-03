import { VStack, HStack, Button, Text, useToast } from '@chakra-ui/react';
import { FaGoogle, FaFacebook } from 'react-icons/fa';
import { useAuth } from '../../hooks/useAuth';
import { useState } from 'react';

const SocialLogin = ({ buttonText = 'Đăng nhập với' }) => {
  const [isLoading, setIsLoading] = useState({
    google: false,
    facebook: false,
  });
  const toast = useToast();
  const { socialLogin } = useAuth();

  const handleSocialLogin = async (provider) => {
    setIsLoading((prev) => ({ ...prev, [provider]: true }));

    try {
      // Lưu provider vào localStorage để sử dụng khi callback
      localStorage.setItem('oauth_provider', provider);

      // Lấy URL xác thực từ backend
      const result = await socialLogin(null, null, provider, true);

      if (result.url) {
        // Chuyển hướng đến URL xác thực
        window.location.href = result.url;
      } else {
        throw new Error(`Không thể lấy URL đăng nhập cho ${provider}`);
      }
    } catch (error) {
      toast({
        title: 'Đăng nhập thất bại',
        description: error.message,
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
      setIsLoading((prev) => ({ ...prev, [provider]: false }));
      localStorage.removeItem('oauth_provider');
    }
  };

  return (
      <VStack spacing={4} width='full'>
        <Text textAlign='center' color='gray.700' fontSize='sm' fontWeight="medium">
          {buttonText}
        </Text>

        <HStack spacing={4} width='full'>
          <Button
              onClick={() => handleSocialLogin('google')}
              leftIcon={<FaGoogle />}
              colorScheme='red'
              variant='outline'
              size='lg'
              flex='1'
              isLoading={isLoading.google}
              _hover={{
                bg: 'red.100',
              }}
          >
            Google
          </Button>

          <Button
              onClick={() => handleSocialLogin('facebook')}
              leftIcon={<FaFacebook />}
              colorScheme='facebook'
              variant='outline'
              size='lg'
              flex='1'
              isLoading={isLoading.facebook}
              _hover={{
                bg: 'blue.100',
              }}
          >
            Facebook
          </Button>
        </HStack>
      </VStack>
  );
};

export default SocialLogin;