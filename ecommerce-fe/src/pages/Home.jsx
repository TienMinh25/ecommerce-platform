import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import {
  Box,
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
import { useEffect, useRef, useState } from 'react';
import ProductGrid from '../components/products/ProductGrid';
import CategorySlider from "../components/category/CategorySlider.jsx";
import productService from '../services/productService';

import slideImage from '../assets/images/sale.png';
import freeshipImage from '../assets/images/freeship.png';
import newCollectionImage from '../assets/images/new_collection.png';

const Home = () => {
  const [currentSlide, setCurrentSlide] = useState(0);
  const sliderRef = useRef(null);
  const [featuredProducts, setFeaturedProducts] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

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

  useEffect(() => {
    const fetchFeaturedProducts = async () => {
      try {
        setIsLoading(true);
        const response = await productService.getFeaturedProducts(24);
        setFeaturedProducts(response.data.data || []);
        setIsLoading(false);
      } catch (err) {
        console.error('Error fetching featured products:', err);
        setError('Không thể tải sản phẩm. Vui lòng thử lại sau.');
        setIsLoading(false);
      }
    };

    fetchFeaturedProducts();
  }, []);

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
          {/* Categories Slider - Using the CategorySlider component */}
          <CategorySlider />

          {/* Popular Products Section */}
          <Box mb={16}>
            <Flex justify='space-between' align='center' mb={8}>
              <Heading as='h2' size='lg'>
                Sản phẩm nổi bật
              </Heading>
            </Flex>

            {/* Product Grid */}
            <ProductGrid
                products={featuredProducts}
                isLoading={isLoading}
                error={error}
            />
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