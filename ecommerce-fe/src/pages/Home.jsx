import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import {
  Box,
  Button,
  Container,
  Flex,
  Grid,
  GridItem,
  Heading,
  IconButton,
  Stack,
  Text,
  VStack
} from '@chakra-ui/react';
import { useRef, useState } from 'react';
import { FaArrowRight } from 'react-icons/fa';
import { Link as RouterLink } from 'react-router-dom';
import ProductCard from '../components/products/ProductCard';

import slideImage from '../assets/images/sale.png';
import freeshipImage from '../assets/images/freeship.png';
import newCollectionImage from '../assets/images/new_collection.png';
import CategorySlider from "../components/category/CategorySlider.jsx";

const Home = () => {
  const [currentSlide, setCurrentSlide] = useState(0);
  const sliderRef = useRef(null);

  const slides = [
    {
      id: 1,
      title: 'Siêu sale',
      description: 'Giảm đến 50% toàn bộ sản phẩm',
      image: slideImage,
    },
    {
      id: 2,
      title: 'Bộ sưu tập mới',
      description: 'Khám phá xu hướng thời trang 2025',
      image: newCollectionImage,
    },
    {
      id: 3,
      title: 'Miễn phí vận chuyển',
      description: 'Cho đơn hàng từ 299.000đ',
      image: freeshipImage,
    },
  ];

  // Extended popular products to 16 items
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
    {
      id: 9,
      name: 'Áo khoác denim nữ',
      image: 'https://via.placeholder.com/300x300?text=Jacket',
      price: 799000,
      originalPrice: 999000,
      rating: 4.7,
      reviewCount: 143,
    },
    {
      id: 10,
      name: 'Quần short nam thể thao',
      image: 'https://via.placeholder.com/300x300?text=Shorts',
      price: 259000,
      originalPrice: 350000,
      rating: 4.4,
      reviewCount: 92,
    },
    {
      id: 11,
      name: 'Áo sơ mi nam dài tay',
      image: 'https://via.placeholder.com/300x300?text=Shirt',
      price: 429000,
      originalPrice: 529000,
      rating: 4.6,
      reviewCount: 178,
    },
    {
      id: 12,
      name: 'Dép nữ đi trong nhà',
      image: 'https://via.placeholder.com/300x300?text=Slippers',
      price: 129000,
      originalPrice: 179000,
      rating: 4.3,
      reviewCount: 146,
    },
    {
      id: 13,
      name: 'Váy đầm dự tiệc',
      image: 'https://via.placeholder.com/300x300?text=PartyDress',
      price: 850000,
      originalPrice: 1100000,
      rating: 4.8,
      reviewCount: 87,
    },
    {
      id: 14,
      name: 'Balo laptop thời trang',
      image: 'https://via.placeholder.com/300x300?text=Backpack',
      price: 590000,
      originalPrice: 790000,
      rating: 4.7,
      reviewCount: 124,
    },
    {
      id: 15,
      name: 'Khăn quàng cổ len',
      image: 'https://via.placeholder.com/300x300?text=Scarf',
      price: 199000,
      originalPrice: 280000,
      rating: 4.5,
      reviewCount: 58,
    },
    {
      id: 16,
      name: 'Găng tay da nam',
      image: 'https://via.placeholder.com/300x300?text=Gloves',
      price: 399000,
      originalPrice: 499000,
      rating: 4.6,
      reviewCount: 42,
    },
  ];

  const prevSlide = () => {
    setCurrentSlide((s) => (s === 0 ? slides.length - 1 : s - 1));
  };

  const nextSlide = () => {
    setCurrentSlide((s) => (s === slides.length - 1 ? 0 : s + 1));
  };

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
                  <Box
                      position='absolute'
                      top='0'
                      left='0'
                      right='0'
                      bottom='0'
                      backgroundImage={`url(${slide.image})`}
                      backgroundSize='cover'
                      backgroundPosition='center'
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
              color='black'
              boxShadow='md'
              opacity='0.8'
              _hover={{ opacity: 1, bg: 'white' }}
              zIndex='1'
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
              color='black'
              boxShadow='md'
              opacity='0.8'
              _hover={{ opacity: 1, bg: 'white' }}
              zIndex='1'
          />
        </Box>

        <Container maxW='container.xl' py={8}>
          {/* Categories Slider - Using the updated component with API service */}
          <CategorySlider />

          {/* Popular Products - 16 products */}
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