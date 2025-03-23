import { CheckIcon, ChevronRightIcon, StarIcon } from '@chakra-ui/icons';
import {
  Avatar,
  Badge,
  Box,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  Button,
  Container,
  Divider,
  Flex,
  Grid,
  GridItem,
  Heading,
  HStack,
  Icon,
  Image,
  List,
  ListIcon,
  ListItem,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  SimpleGrid,
  Tab,
  Table,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Tbody,
  Td,
  Text,
  Tr,
  useMediaQuery,
  useToast,
  VStack,
} from '@chakra-ui/react';
import { useEffect, useState } from 'react';
import {
  FaExchangeAlt,
  FaHeart,
  FaInfoCircle,
  FaRegHeart,
  FaShieldAlt,
  FaShoppingCart,
  FaTruck,
} from 'react-icons/fa';
import { Link as RouterLink, useNavigate, useParams } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import ProductCard from '../components/products/ProductCard';

const ProductDetail = () => {
  const { id } = useParams();
  const [product, setProduct] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [selectedImage, setSelectedImage] = useState(0);
  const [isLiked, setIsLiked] = useState(false);
  const [relatedProducts, setRelatedProducts] = useState([]);
  const [selectedColor, setSelectedColor] = useState('');
  const [selectedSize, setSelectedSize] = useState('');
  const navigate = useNavigate();
  const toast = useToast();
  const { isAuthenticated } = useAuth();
  const [isLargerThan768] = useMediaQuery('(min-width: 768px)');

  useEffect(() => {
    const fetchProduct = async () => {
      setIsLoading(true);
      try {
        // TODO: Implement actual API call
        // For now, use dummy data

        // Simulate network delay
        setTimeout(() => {
          // Mock product data
          const mockProduct = {
            id: parseInt(id),
            name: 'Áo khoác denim cao cấp',
            description:
              'Áo khoác denim chất liệu cao cấp, thiết kế hiện đại, phù hợp với nhiều phong cách thời trang khác nhau. Sản phẩm được sản xuất từ vải denim cao cấp, bền đẹp và thoáng mát.',
            price: 850000,
            originalPrice: 1050000,
            rating: 4.7,
            reviewCount: 182,
            inStock: true,
            sku: 'DEN-JAC-001',
            category: 'Áo khoác',
            brand: 'Brand Name',
            colors: ['Xanh đậm', 'Xanh nhạt', 'Đen'],
            sizes: ['S', 'M', 'L', 'XL'],
            images: [
              'https://via.placeholder.com/600x600?text=Denim+Jacket+1',
              'https://via.placeholder.com/600x600?text=Denim+Jacket+2',
              'https://via.placeholder.com/600x600?text=Denim+Jacket+3',
              'https://via.placeholder.com/600x600?text=Denim+Jacket+4',
            ],
            details: [
              'Chất liệu: 100% cotton denim',
              'Kiểu dáng: Regular fit',
              'Xuất xứ: Việt Nam',
              'Mùa: Xuân/Thu',
              'Bảo quản: Giặt máy nhiệt độ thấp, không tẩy',
            ],
            specifications: {
              'Chất liệu': '100% cotton denim',
              'Kiểu dáng': 'Regular fit',
              'Xuất xứ': 'Việt Nam',
              Mùa: 'Xuân/Thu',
              'Họa tiết': 'Trơn',
              'Bảo quản': 'Giặt máy nhiệt độ thấp, không tẩy',
            },
            reviews: [
              {
                id: 1,
                user: 'Nguyễn Văn A',
                avatar: 'https://via.placeholder.com/40x40',
                rating: 5,
                date: '15/02/2025',
                comment:
                  'Sản phẩm chất lượng tốt, đúng kích cỡ, giao hàng nhanh. Rất hài lòng!',
              },
              {
                id: 2,
                user: 'Trần Thị B',
                avatar: 'https://via.placeholder.com/40x40',
                rating: 4,
                date: '10/02/2025',
                comment:
                  'Áo đẹp, chất vải tốt. Tuy nhiên size hơi rộng một chút so với mô tả.',
              },
              {
                id: 3,
                user: 'Lê Văn C',
                avatar: 'https://via.placeholder.com/40x40',
                rating: 5,
                date: '05/02/2025',
                comment:
                  'Mẫu mã đẹp, thiết kế trẻ trung, giá cả hợp lý. Sẽ ủng hộ shop lần sau.',
              },
            ],
          };

          // Mock related products
          const mockRelatedProducts = [
            {
              id: 101,
              name: 'Áo khoác denim vintage',
              image: 'https://via.placeholder.com/300x300?text=Related+1',
              price: 750000,
              originalPrice: 900000,
              rating: 4.5,
              reviewCount: 98,
            },
            {
              id: 102,
              name: 'Áo khoác denim oversize',
              image: 'https://via.placeholder.com/300x300?text=Related+2',
              price: 820000,
              originalPrice: 950000,
              rating: 4.6,
              reviewCount: 124,
            },
            {
              id: 103,
              name: 'Áo khoác jean nữ',
              image: 'https://via.placeholder.com/300x300?text=Related+3',
              price: 680000,
              originalPrice: 800000,
              rating: 4.4,
              reviewCount: 76,
            },
            {
              id: 104,
              name: 'Áo khoác denim phối nón',
              image: 'https://via.placeholder.com/300x300?text=Related+4',
              price: 890000,
              originalPrice: 1100000,
              rating: 4.8,
              reviewCount: 65,
            },
          ];

          setProduct(mockProduct);
          setSelectedColor(mockProduct.colors[0]);
          setSelectedSize(mockProduct.sizes[1]); // Default to M
          setRelatedProducts(mockRelatedProducts);
          setIsLoading(false);
        }, 1000);
      } catch (error) {
        setError(error);
        setIsLoading(false);
      }
    };

    if (id) {
      fetchProduct();
    }
  }, [id]);

  const handleAddToCart = () => {
    if (!selectedColor || !selectedSize) {
      toast({
        title: 'Vui lòng chọn',
        description: !selectedColor
          ? 'Vui lòng chọn màu sắc'
          : 'Vui lòng chọn kích thước',
        status: 'warning',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    // TODO: Implement add to cart functionality
    toast({
      title: 'Thêm vào giỏ hàng',
      description: `Đã thêm ${quantity} ${product.name} (${selectedColor}, ${selectedSize}) vào giỏ hàng`,
      status: 'success',
      duration: 3000,
      isClosable: true,
    });
  };

  const handleBuyNow = () => {
    if (!selectedColor || !selectedSize) {
      toast({
        title: 'Vui lòng chọn',
        description: !selectedColor
          ? 'Vui lòng chọn màu sắc'
          : 'Vui lòng chọn kích thước',
        status: 'warning',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    // TODO: Implement buy now functionality
    handleAddToCart();
    navigate('/checkout');
  };

  const handleToggleFavorite = () => {
    if (!isAuthenticated) {
      toast({
        title: 'Yêu cầu đăng nhập',
        description: 'Vui lòng đăng nhập để lưu sản phẩm yêu thích',
        status: 'info',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    setIsLiked(!isLiked);

    toast({
      title: isLiked ? 'Đã xóa khỏi yêu thích' : 'Đã thêm vào yêu thích',
      status: 'success',
      duration: 3000,
      isClosable: true,
    });
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND',
    }).format(price);
  };

  if (isLoading) {
    return (
      <Container maxW='container.xl' py={10}>
        <Box textAlign='center'>Đang tải thông tin sản phẩm...</Box>
      </Container>
    );
  }

  if (error || !product) {
    return (
      <Container maxW='container.xl' py={10}>
        <Box textAlign='center'>
          Không thể tải thông tin sản phẩm. Vui lòng thử lại sau.
        </Box>
      </Container>
    );
  }

  const discountPercent = product.originalPrice
    ? Math.round(
        ((product.originalPrice - product.price) / product.originalPrice) * 100,
      )
    : 0;

  return (
    <Container maxW='container.xl' py={8}>
      {/* Breadcrumb */}
      <Breadcrumb
        separator={<ChevronRightIcon color='gray.500' />}
        mb={8}
        fontSize='sm'
      >
        <BreadcrumbItem>
          <BreadcrumbLink as={RouterLink} to='/'>
            Trang chủ
          </BreadcrumbLink>
        </BreadcrumbItem>
        <BreadcrumbItem>
          <BreadcrumbLink as={RouterLink} to='/products'>
            Sản phẩm
          </BreadcrumbLink>
        </BreadcrumbItem>
        <BreadcrumbItem isCurrentPage>
          <BreadcrumbLink href='#'>{product.name}</BreadcrumbLink>
        </BreadcrumbItem>
      </Breadcrumb>

      {/* Product Info */}
      <Grid templateColumns={{ base: '1fr', md: '1fr 1fr' }} gap={10} mb={16}>
        {/* Product Images */}
        <GridItem>
          <Box position='relative' mb={4}>
            {discountPercent > 0 && (
              <Badge
                position='absolute'
                top='10px'
                left='10px'
                colorScheme='red'
                variant='solid'
                borderRadius='md'
                px={2}
                py={1}
                zIndex={1}
              >
                -{discountPercent}%
              </Badge>
            )}
            <Image
              src={product.images[selectedImage]}
              alt={product.name}
              borderRadius='md'
              width='100%'
              height='auto'
              objectFit='cover'
            />
          </Box>
          <SimpleGrid columns={4} spacing={4}>
            {product.images.map((image, index) => (
              <Box
                key={index}
                borderWidth={selectedImage === index ? '2px' : '1px'}
                borderColor={selectedImage === index ? 'brand.500' : 'gray.200'}
                borderRadius='md'
                overflow='hidden'
                cursor='pointer'
                onClick={() => setSelectedImage(index)}
              >
                <Image src={image} alt={`${product.name} ${index + 1}`} />
              </Box>
            ))}
          </SimpleGrid>
        </GridItem>

        {/* Product Details */}
        <GridItem>
          <VStack align='stretch' spacing={6}>
            <Box>
              <Heading as='h1' size='xl' mb={2}>
                {product.name}
              </Heading>

              <HStack mb={4}>
                <Box display='flex' alignItems='center'>
                  {Array(5)
                    .fill('')
                    .map((_, i) => (
                      <StarIcon
                        key={i}
                        color={
                          i < Math.floor(product.rating)
                            ? 'yellow.400'
                            : 'gray.300'
                        }
                      />
                    ))}
                  <Text ml={2} color='gray.600'>
                    ({product.rating}) - {product.reviewCount} đánh giá
                  </Text>
                </Box>
              </HStack>

              <Box mb={6}>
                <Flex align='baseline'>
                  <Heading as='h2' size='xl' color='brand.500' mr={3}>
                    {formatPrice(product.price)}
                  </Heading>
                  {product.originalPrice &&
                    product.originalPrice > product.price && (
                      <Text as='s' color='gray.500' fontSize='lg'>
                        {formatPrice(product.originalPrice)}
                      </Text>
                    )}
                  {discountPercent > 0 && (
                    <Badge colorScheme='red' ml={2} fontSize='sm'>
                      Giảm {discountPercent}%
                    </Badge>
                  )}
                </Flex>
              </Box>
            </Box>

            <Divider />

            {/* Description */}
            <Box>
              <Text color='gray.700'>{product.description}</Text>
            </Box>

            <Divider />

            {/* Colors */}
            <Box>
              <Text fontWeight='semibold' mb={2}>
                Màu sắc: {selectedColor}
              </Text>
              <HStack spacing={2}>
                {product.colors.map((color) => (
                  <Button
                    key={color}
                    size='sm'
                    variant={selectedColor === color ? 'solid' : 'outline'}
                    colorScheme={selectedColor === color ? 'brand' : 'gray'}
                    onClick={() => setSelectedColor(color)}
                  >
                    {color}
                  </Button>
                ))}
              </HStack>
            </Box>

            {/* Sizes */}
            <Box>
              <Flex justify='space-between' align='center' mb={2}>
                <Text fontWeight='semibold'>Kích thước: {selectedSize}</Text>
                <Text
                  as={RouterLink}
                  to='#'
                  color='brand.500'
                  fontSize='sm'
                  textDecoration='underline'
                >
                  Hướng dẫn chọn size
                </Text>
              </Flex>
              <HStack spacing={2}>
                {product.sizes.map((size) => (
                  <Button
                    key={size}
                    size='sm'
                    variant={selectedSize === size ? 'solid' : 'outline'}
                    colorScheme={selectedSize === size ? 'brand' : 'gray'}
                    onClick={() => setSelectedSize(size)}
                  >
                    {size}
                  </Button>
                ))}
              </HStack>
            </Box>

            {/* Quantity */}
            <Box>
              <Text fontWeight='semibold' mb={2}>
                Số lượng:
              </Text>
              <NumberInput
                max={10}
                min={1}
                value={quantity}
                onChange={(value) => setQuantity(parseInt(value))}
                size='md'
                maxW={32}
              >
                <NumberInputField />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
              <Text mt={2} fontSize='sm' color='gray.500'>
                {product.inStock ? 'Còn hàng' : 'Hết hàng'}
              </Text>
            </Box>

            {/* Actions */}
            <HStack spacing={4} pt={4}>
              <Button
                leftIcon={<FaShoppingCart />}
                colorScheme='brand'
                variant='outline'
                size='lg'
                flex='1'
                onClick={handleAddToCart}
                isDisabled={!product.inStock}
              >
                Thêm vào giỏ
              </Button>
              <Button
                colorScheme='brand'
                size='lg'
                flex='1'
                onClick={handleBuyNow}
                isDisabled={!product.inStock}
              >
                Mua ngay
              </Button>
              <IconButton
                aria-label='Add to favorites'
                icon={isLiked ? <FaHeart /> : <FaRegHeart />}
                colorScheme={isLiked ? 'red' : 'gray'}
                variant='outline'
                size='lg'
                onClick={handleToggleFavorite}
              />
            </HStack>

            {/* Product Info */}
            <Box borderWidth='1px' borderRadius='md' p={4} mt={4} bg='gray.50'>
              <SimpleGrid columns={{ base: 1, md: 2 }} spacing={4}>
                <Flex align='center'>
                  <Icon as={FaTruck} mr={2} color='brand.500' />
                  <Box>
                    <Text fontWeight='semibold'>Miễn phí vận chuyển</Text>
                    <Text fontSize='sm' color='gray.600'>
                      Cho đơn hàng từ 500k
                    </Text>
                  </Box>
                </Flex>
                <Flex align='center'>
                  <Icon as={FaExchangeAlt} mr={2} color='brand.500' />
                  <Box>
                    <Text fontWeight='semibold'>Đổi trả miễn phí</Text>
                    <Text fontSize='sm' color='gray.600'>
                      Trong 30 ngày
                    </Text>
                  </Box>
                </Flex>
                <Flex align='center'>
                  <Icon as={FaShieldAlt} mr={2} color='brand.500' />
                  <Box>
                    <Text fontWeight='semibold'>Bảo hành chính hãng</Text>
                    <Text fontSize='sm' color='gray.600'>
                      12 tháng
                    </Text>
                  </Box>
                </Flex>
                <Flex align='center'>
                  <Icon as={FaInfoCircle} mr={2} color='brand.500' />
                  <Box>
                    <Text fontWeight='semibold'>Mã sản phẩm</Text>
                    <Text fontSize='sm' color='gray.600'>
                      {product.sku}
                    </Text>
                  </Box>
                </Flex>
              </SimpleGrid>
            </Box>
          </VStack>
        </GridItem>
      </Grid>

      {/* Product Tabs */}
      <Box
        mb={16}
        borderWidth='1px'
        borderRadius='lg'
        overflow='hidden'
        bg='white'
      >
        <Tabs colorScheme='brand' defaultIndex={0}>
          <TabList px={{ base: 2, md: 6 }}>
            <Tab fontWeight='semibold' py={4}>
              Chi tiết sản phẩm
            </Tab>
            <Tab fontWeight='semibold' py={4}>
              Thông số kỹ thuật
            </Tab>
            <Tab fontWeight='semibold' py={4}>
              Đánh giá ({product.reviews.length})
            </Tab>
          </TabList>

          <TabPanels>
            {/* Product Details */}
            <TabPanel p={{ base: 4, md: 8 }}>
              <Box>
                <Heading as='h3' size='md' mb={4}>
                  Thông tin chi tiết
                </Heading>
                <Text mb={6}>{product.description}</Text>

                <List spacing={2}>
                  {product.details.map((detail, index) => (
                    <ListItem key={index} display='flex' alignItems='center'>
                      <ListIcon as={CheckIcon} color='brand.500' />
                      <Text>{detail}</Text>
                    </ListItem>
                  ))}
                </List>
              </Box>
            </TabPanel>

            {/* Specifications */}
            <TabPanel p={{ base: 4, md: 8 }}>
              <Box>
                <Heading as='h3' size='md' mb={4}>
                  Thông số kỹ thuật
                </Heading>
                <Table variant='simple'>
                  <Tbody>
                    {Object.entries(product.specifications).map(
                      ([key, value]) => (
                        <Tr key={key}>
                          <Td fontWeight='semibold' width='30%'>
                            {key}
                          </Td>
                          <Td>{value}</Td>
                        </Tr>
                      ),
                    )}
                  </Tbody>
                </Table>
              </Box>
            </TabPanel>

            {/* Reviews */}
            <TabPanel p={{ base: 4, md: 8 }}>
              <Box>
                <Flex justify='space-between' align='center' mb={6}>
                  <Heading as='h3' size='md'>
                    Đánh giá từ khách hàng
                  </Heading>
                  <Button colorScheme='brand' variant='outline'>
                    Viết đánh giá
                  </Button>
                </Flex>

                <Box mb={8}>
                  <Flex
                    align='center'
                    justify='space-between'
                    p={4}
                    borderWidth='1px'
                    borderRadius='md'
                    bg='gray.50'
                  >
                    <Box>
                      <Text fontSize='xl' fontWeight='bold'>
                        {product.rating}/5
                      </Text>
                      <HStack spacing={1}>
                        {Array(5)
                          .fill('')
                          .map((_, i) => (
                            <StarIcon
                              key={i}
                              color={
                                i < Math.floor(product.rating)
                                  ? 'yellow.400'
                                  : 'gray.300'
                              }
                            />
                          ))}
                      </HStack>
                      <Text color='gray.600' fontSize='sm'>
                        {product.reviewCount} đánh giá
                      </Text>
                    </Box>

                    <Box display={{ base: 'none', md: 'block' }}>
                      <SimpleGrid columns={5} spacing={4}>
                        {[5, 4, 3, 2, 1].map((num) => (
                          <Flex key={num} align='center'>
                            <Text mr={2}>{num}</Text>
                            <StarIcon color='yellow.400' />
                            <Box
                              w='100px'
                              h='8px'
                              bg='gray.200'
                              borderRadius='full'
                              ml={2}
                              position='relative'
                              overflow='hidden'
                            >
                              <Box
                                position='absolute'
                                h='100%'
                                w={
                                  num === 5 ? '70%' : num === 4 ? '20%' : '10%'
                                }
                                bg='brand.500'
                                borderRadius='full'
                              />
                            </Box>
                          </Flex>
                        ))}
                      </SimpleGrid>
                    </Box>
                  </Flex>
                </Box>

                <VStack spacing={4} align='stretch'>
                  {product.reviews.map((review) => (
                    <Box
                      key={review.id}
                      p={4}
                      borderWidth='1px'
                      borderRadius='md'
                    >
                      <Flex mb={4}>
                        <Avatar
                          size='sm'
                          name={review.user}
                          src={review.avatar}
                          mr={4}
                        />
                        <Box flex='1'>
                          <Flex justify='space-between' align='center'>
                            <Text fontWeight='bold'>{review.user}</Text>
                            <Text fontSize='sm' color='gray.500'>
                              {review.date}
                            </Text>
                          </Flex>
                          <HStack spacing={1}>
                            {Array(5)
                              .fill('')
                              .map((_, i) => (
                                <StarIcon
                                  key={i}
                                  size='sm'
                                  color={
                                    i < review.rating
                                      ? 'yellow.400'
                                      : 'gray.300'
                                  }
                                />
                              ))}
                          </HStack>
                        </Box>
                      </Flex>
                      <Text>{review.comment}</Text>
                    </Box>
                  ))}
                </VStack>
              </Box>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>

      {/* Related Products */}
      <Box mb={16}>
        <Heading as='h2' size='xl' mb={6}>
          Sản phẩm liên quan
        </Heading>
        <SimpleGrid columns={{ base: 2, md: 4 }} spacing={4}>
          {relatedProducts.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </SimpleGrid>
      </Box>
    </Container>
  );
};

export default ProductDetail;
