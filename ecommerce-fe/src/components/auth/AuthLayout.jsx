import { Box, Flex, Image, useBreakpointValue } from '@chakra-ui/react';
import Logo from '../ui/Logo';

const AuthLayout = ({ children }) => {
  const showImage = useBreakpointValue({ base: false, md: true });

  return (
    <Flex minH='100vh' bg='gray.50'>
      {/* Left side - Form */}
      <Flex
        direction='column'
        w={{ base: 'full', md: '50%' }}
        p={{ base: 6, md: 10 }}
        justify='center'
        align='center'
      >
        <Box w='full' maxW='400px'>
          <Box mb={8} textAlign='center'>
            <Logo />
          </Box>
          {children}
        </Box>
      </Flex>

      {/* Right side - Image (visible on md screens and up) */}
      {showImage && (
        <Flex
          w='50%'
          bg='brand.500'
          justify='center'
          align='center'
          position='relative'
          overflow='hidden'
        >
          <Image
            src='/src/assets/images/hero-image.svg'
            alt='Shopping illustration'
            objectFit='cover'
            w='full'
            h='full'
            position='absolute'
            opacity='0.8'
          />
          <Box
            position='absolute'
            maxW='md'
            p={8}
            borderRadius='md'
            bg='rgba(255, 255, 255, 0.3)'
            backdropFilter='blur(5px)'
            textAlign='center'
          >
            <Box
              as='h2'
              fontSize='3xl'
              fontWeight='bold'
              mb={4}
              color='brand.700'
            >
              Chào mừng đến với Minh Plaza
            </Box>
            <Box fontSize='lg' color='gray.700'>
              Nền tảng mua sắm trực tuyến hàng đầu với hàng triệu sản phẩm và ưu
              đãi hấp dẫn.
            </Box>
          </Box>
        </Flex>
      )}
    </Flex>
  );
};

export default AuthLayout;
