import {
    Box,
    Container,
    Flex,
    HStack,
    Icon,
    Link,
    Text,
    Tooltip,
  } from '@chakra-ui/react';
  import {
    FaEnvelope,
    FaFacebook,
    FaInstagram,
    FaPhone,
    FaTwitter,
  } from 'react-icons/fa';
  import Logo from '../../components/ui/Logo';
  
  const DashboardFooter = () => {
    return (
      <Box bg='gray.50' color='gray.600' py={4} borderTop="1px" borderColor="gray.200">
        <Container maxW={'container.xl'}>
          <Flex
            direction={{ base: 'column', md: 'row' }}
            justify="space-between"
            align="center"
            gap={4}
          >
            <HStack spacing={4}>
              <Logo size='sm' />
              <Text fontSize="xs">© 2025 Minh Plaza</Text>
            </HStack>
  
            <HStack spacing={3} fontSize="xs">
              <Link href="#">Chính sách bảo mật</Link>
              <Text>•</Text>
              <Link href="#">Điều khoản sử dụng</Link>
            </HStack>
  
            <HStack spacing={4}>
              <Tooltip label="Facebook">
                <Link href='#' isExternal>
                  <Icon as={FaFacebook} w={4} h={4} color='blue.500' />
                </Link>
              </Tooltip>
              <Tooltip label="Instagram">
                <Link href='#' isExternal>
                  <Icon as={FaInstagram} w={4} h={4} color='pink.500' />
                </Link>
              </Tooltip>
              <Tooltip label="Email">
                <Link href='mailto:contact@minhplaza.vn' isExternal>
                  <Icon as={FaEnvelope} w={4} h={4} color='gray.500' />
                </Link>
              </Tooltip>
              <Tooltip label="1900 1234">
                <Link href='tel:19001234' isExternal>
                  <Icon as={FaPhone} w={4} h={4} color='gray.500' />
                </Link>
              </Tooltip>
            </HStack>
          </Flex>
        </Container>
      </Box>
    );
  };
  
  export default DashboardFooter;