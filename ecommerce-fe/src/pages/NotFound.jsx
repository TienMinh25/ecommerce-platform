import {
  Box,
  Heading,
  Text,
  Button,
  Image,
  VStack,
  Container,
} from '@chakra-ui/react';
import { Link as RouterLink } from 'react-router-dom';

const NotFound = () => {
  return (
    <Container maxW='container.xl' py={20}>
      <VStack spacing={8} textAlign='center'>
        <Box width='100%' display='flex' justifyContent='center' mb={1}>
          <Image
            src='/src/assets/images/not-found.svg'
            alt='Không tìm thấy trang'
            width='auto'
            height='auto'
            maxW='600px'
            objectFit='contain'
            fallbackSrc='https://via.placeholder.com/300x300?text=404'
          />
        </Box>

        <Text fontSize='lg' color='gray.600' maxW='md'>
          Trang bạn đang tìm kiếm có thể đã bị xóa, đổi tên hoặc tạm thời không
          khả dụng.
        </Text>

        <Box pt={4}>
          <Button
            as={RouterLink}
            to='/'
            colorScheme='brand'
            size='lg'
            rightIcon={
              <Box as='span' ml={2}>
                →
              </Box>
            }
          >
            Trở về trang chủ
          </Button>
        </Box>
      </VStack>
    </Container>
  );
};

export default NotFound;
