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
      // TODO: Implement actual social login flow
      const result = await socialLogin(provider);

      if (result.success) {
        toast({
          title: 'Đăng nhập thành công',
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
      } else {
        throw new Error(result.error || `Đăng nhập với ${provider} thất bại`);
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
      setIsLoading((prev) => ({ ...prev, [provider]: false }));
    }
  };

  return (
    <VStack spacing={4} width='full'>
      <Text textAlign='center' color='gray.500' fontSize='sm'>
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
        >
          Facebook
        </Button>
      </HStack>
    </VStack>
  );
};

export default SocialLogin;
