import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import {
  AspectRatio,
  Box,
  Button,
  Container,
  Flex,
  Grid,
  GridItem,
  Heading,
  IconButton,
  Image,
  Stack,
  Text,
  useBreakpointValue,
  VStack
} from '@chakra-ui/react';
import { useRef, useState } from 'react';
import { FaArrowRight } from 'react-icons/fa';
import { Link as RouterLink } from 'react-router-dom';
import ProductCard from '../components/products/ProductCard';

const Home = () => {
  const [currentSlide, setCurrentSlide] = useState(0);
  const sliderRef = useRef(null);

  const slides = [
    {
      id: 1,
      title: 'Siêu sale tháng 6',
      description: 'Giảm đến 70% toàn bộ sản phẩm',
      image: 'https://via.placeholder.com/1200x400?text=Banner+1',
      buttonText: 'Mua ngay',
      buttonLink: '/promotions/summer-sale',
    },
    {
      id: 2,
      title: 'Bộ sưu tập mới',
      description: 'Khám phá xu hướng thời trang 2025',
      image: 'https://via.placeholder.com/1200x400?text=Banner+2',
      buttonText: 'Xem ngay',
      buttonLink: '/collections/new-arrivals',
    },
    {
      id: 3,
      title: 'Miễn phí vận chuyển',
      description: 'Cho đơn hàng từ 299.000đ',
      image: 'https://via.placeholder.com/1200x400?text=Banner+3',
      buttonText: 'Tìm hiểu thêm',
      buttonLink: '/shipping-policy',
    },
  ];

  const featuredCategories = [
    {
      id: 1,
      name: 'Thời trang nam',
      image: 'https://via.placeholder.com/300x300?text=Men',
      link: '/category/men',
    },
    {
      id: 2,
      name: 'Thời trang nữ',
      image: 'https://via.placeholder.com/300x300?text=Women',
      link: '/category/women',
    },
    {
      id: 3,
      name: 'Đồng hồ',
      image: 'https://via.placeholder.com/300x300?text=Watches',
      link: '/category/watches',
    },
    {
      id: 4,
      name: 'Giày dép',
      image: 'https://via.placeholder.com/300x300?text=Shoes',
      link: '/category/shoes',
    },
    {
      id: 5,
      name: 'Túi xách',
      image: 'https://via.placeholder.com/300x300?text=Bags',
      link: '/category/bags',
    },
    {
      id: 6,
      name: 'Phụ kiện',
      image: 'https://via.placeholder.com/300x300?text=Accessories',
      link: '/category/accessories',
    },
  ];

  const popularProducts = [
    {
      id: 1,
      name: 'Áo thun nam basic',
      image: 'https://via.placeholder.com/300x300?text=T-Shirt',
      price: 199000,
      originalPrice: 249000,
      rating: 4.5,
      reviewCount: 120,
    },
    {
      id: 2,
      name: 'Áo sơ mi nữ công sở',
      image: 'https://via.placeholder.com/300x300?text=Blouse',
      price: 349000,
      originalPrice: 449000,
      rating: 4.3,
      reviewCount: 86,
    },
    {
      id: 3,
      name: 'Quần jean nam slim fit',
      image: 'https://via.placeholder.com/300x300?text=Jeans',
      price: 499000,
      originalPrice: 599000,
      rating: 4.7,
      reviewCount: 203,
    },
    {
      id: 4,
      name: 'Đầm nữ dáng suông',
      image: 'https://via.placeholder.com/300x300?text=Dress',
      price: 545000,
      originalPrice: 650000,
      rating: 4.6,
      reviewCount: 154,
    },
    {
      id: 5,
      name: 'Giày thể thao nam',
      image: 'https://via.placeholder.com/300x300?text=Sneakers',
      price: 899000,
      originalPrice: 1200000,
      rating: 4.8,
      reviewCount: 312,
    },
    {
      id: 6,
      name: 'Túi xách nữ thời trang',
      image: 'https://via.placeholder.com/300x300?text=Handbag',
      price: 750000,
      originalPrice: 950000,
      rating: 4.4,
      reviewCount: 98,
    },
    {
      id: 7,
      name: 'Đồng hồ nam cao cấp',
      image: 'https://via.placeholder.com/300x300?text=Watch',
      price: 2490000,
      originalPrice: 2990000,
      rating: 4.9,
      reviewCount: 76,
    },
    {
      id: 8,
      name: 'Kính mát thời trang',
      image: 'https://via.placeholder.com/300x300?text=Sunglasses',
      price: 450000,
      originalPrice: 550000,
      rating: 4.2,
      reviewCount: 65,
    },
  ];

  const prevSlide = () => {
    setCurrentSlide((s) => (s === 0 ? slides.length - 1 : s - 1));
  };

  const nextSlide = () => {
    setCurrentSlide((s) => (s === slides.length - 1 ? 0 : s + 1));
  };

  const columns = useBreakpointValue({ base: 2, md: 3, lg: 4 });

  return (
    <Box>
      {/* Hero Slider */}
      <Box position='relative' overflow='hidden' mb={10}>
        <Flex
          ref={sliderRef}
          transition='transform 0.5s ease'
          transform={`translateX(-${currentSlide * 100}%)`}
        >
          {slides.map((slide) => (
            <Box
              key={slide.id}
              w='100%'
              position='relative'
              minW='100%'
              h={{ base: '200px', md: '300px', lg: '400px' }}
            >
              <Image
                src={slide.image}
                alt={slide.title}
                objectFit='cover'
                w='100%'
                h='100%'
              />
              <Box
                position='absolute'
                top='0'
                left='0'
                right='0'
                bottom='0'
                bg='rgba(0,0,0,0.4)'
                display='flex'
                alignItems='center'
                justifyContent='center'
              >
                <Container maxW='container.lg'>
                  <Stack
                    spacing={4}
                    color='white'
                    textAlign={{ base: 'center', md: 'left' }}
                  >
                    <Heading as='h1' size='xl' fontWeight='bold'>
                      {slide.title}
                    </Heading>
                    <Text fontSize={{ base: 'md', md: 'lg' }}>
                      {slide.description}
                    </Text>
                    <Box>
                      <Button
                        as={RouterLink}
                        to={slide.buttonLink}
                        colorScheme='brand'
                        size='lg'
                      >
                        {slide.buttonText}
                      </Button>
                    </Box>
                  </Stack>
                </Container>
              </Box>
            </Box>
          ))}
        </Flex>

        <IconButton
          aria-label='Previous slide'
          icon={<ChevronLeftIcon boxSize={8} />}
          position='absolute'
          left={{ base: 2, md: 8 }}
          top='50%'
          transform='translateY(-50%)'
          borderRadius='full'
          onClick={prevSlide}
          bg='white'
          opacity='0.8'
          _hover={{ opacity: 1 }}
        />

        <IconButton
          aria-label='Next slide'
          icon={<ChevronRightIcon boxSize={8} />}
          position='absolute'
          right={{ base: 2, md: 8 }}
          top='50%'
          transform='translateY(-50%)'
          borderRadius='full'
          onClick={nextSlide}
          bg='white'
          opacity='0.8'
          _hover={{ opacity: 1 }}
        />
      </Box>

      <Container maxW='container.xl' py={8}>
        {/* Featured Categories */}
        <Box mb={16}>
          <Flex justify='space-between' align='center' mb={8}>
            <Heading as='h2' size='lg'>
              Danh mục nổi bật
            </Heading>
            <Button
              as={RouterLink}
              to='/categories'
              variant='link'
              colorScheme='brand'
              rightIcon={<FaArrowRight />}
            >
              Xem tất cả
            </Button>
          </Flex>

          <Grid
            templateColumns={{
              base: 'repeat(2, 1fr)',
              md: 'repeat(3, 1fr)',
              lg: 'repeat(6, 1fr)',
            }}
            gap={6}
          >
            {featuredCategories.map((category) => (
              <GridItem key={category.id}>
                <Box
                  as={RouterLink}
                  to={category.link}
                  borderRadius='lg'
                  overflow='hidden'
                  transition='transform 0.3s'
                  _hover={{ transform: 'translateY(-5px)' }}
                >
                  <AspectRatio ratio={1}>
                    <Image
                      src={category.image}
                      alt={category.name}
                      objectFit='cover'
                    />
                  </AspectRatio>
                  <Box
                    position='absolute'
                    bottom='0'
                    left='0'
                    right='0'
                    bg='rgba(0,0,0,0.7)'
                    p={3}
                    textAlign='center'
                  >
                    <Text color='white' fontWeight='semibold'>
                      {category.name}
                    </Text>
                  </Box>
                </Box>
              </GridItem>
            ))}
          </Grid>
        </Box>

        {/* Popular Products */}
        <Box mb={16}>
          <Flex justify='space-between' align='center' mb={8}>
            <Heading as='h2' size='lg'>
              Sản phẩm nổi bật
            </Heading>
            <Button
              as={RouterLink}
              to='/products'
              variant='link'
              colorScheme='brand'
              rightIcon={<FaArrowRight />}
            >
              Xem tất cả
            </Button>
          </Flex>

          <Grid
            templateColumns={{
              base: 'repeat(2, 1fr)',
              md: 'repeat(3, 1fr)',
              lg: 'repeat(4, 1fr)',
            }}
            gap={6}
          >
            {popularProducts.map((product) => (
              <GridItem key={product.id}>
                <ProductCard product={product} />
              </GridItem>
            ))}
          </Grid>
        </Box>

        {/* Features */}
        <Box mb={16}>
          <Grid
            templateColumns={{
              base: 'repeat(1, 1fr)',
              md: 'repeat(2, 1fr)',
              lg: 'repeat(4, 1fr)',
            }}
            gap={8}
          >
            <GridItem>
              <VStack align='center' spacing={4}>
                <Box
                  p={4}
                  borderRadius='full'
                  bg='brand.50'
                  color='brand.500'
                  fontSize='2xl'
                >
                  🚚
                </Box>
                <Text fontWeight='bold' fontSize='lg'>
                  Giao hàng miễn phí
                </Text>
                <Text textAlign='center' color='gray.600'>
                  Cho đơn hàng từ 299.000đ
                </Text>
              </VStack>
            </GridItem>

            <GridItem>
              <VStack align='center' spacing={4}>
                <Box
                  p={4}
                  borderRadius='full'
                  bg='brand.50'
                  color='brand.500'
                  fontSize='2xl'
                >
                  🔄
                </Box>
                <Text fontWeight='bold' fontSize='lg'>
                  Đổi trả dễ dàng
                </Text>
                <Text textAlign='center' color='gray.600'>
                  30 ngày đổi trả miễn phí
                </Text>
              </VStack>
            </GridItem>

            <GridItem>
              <VStack align='center' spacing={4}>
                <Box
                  p={4}
                  borderRadius='full'
                  bg='brand.50'
                  color='brand.500'
                  fontSize='2xl'
                >
                  💰
                </Box>
                <Text fontWeight='bold' fontSize='lg'>
                  Thanh toán an toàn
                </Text>
                <Text textAlign='center' color='gray.600'>
                  Nhiều phương thức thanh toán
                </Text>
              </VStack>
            </GridItem>

            <GridItem>
              <VStack align='center' spacing={4}>
                <Box
                  p={4}
                  borderRadius='full'
                  bg='brand.50'
                  color='brand.500'
                  fontSize='2xl'
                >
                  🎁
                </Box>
                <Text fontWeight='bold' fontSize='lg'>
                  Ưu đãi thành viên
                </Text>
                <Text textAlign='center' color='gray.600'>
                  Tích điểm và nhận quà hấp dẫn
                </Text>
              </VStack>
            </GridItem>
          </Grid>
        </Box>
      </Container>
    </Box>
  );
};

export default Home;
