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
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  SimpleGrid,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  useMediaQuery,
  useToast,
  VStack,
  Skeleton,
  Center,
} from '@chakra-ui/react';
import { useEffect, useState, useMemo } from 'react';
import {
  FaExchangeAlt,
  FaInfoCircle,
  FaShieldAlt,
  FaShoppingCart,
  FaStore,
  FaTruck,
  FaThumbsUp,
} from 'react-icons/fa';
import { Link as RouterLink, useNavigate, useParams } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import productService from '../services/productService';
import PageTitle from "./PageTitle.jsx";
import AddToCartButton from "../components/cart/AddToCartButton.jsx";

const ProductDetail = () => {
  const { id } = useParams();
  const [product, setProduct] = useState(null);
  const [reviews, setReviews] = useState([]);
  const [reviewMetadata, setReviewMetadata] = useState({
    total_items: 0,
    total_pages: 0,
    page: 1,
    limit: 6
  });
  const [isLoadingProduct, setIsLoadingProduct] = useState(true);
  const [isLoadingReviews, setIsLoadingReviews] = useState(true);
  const [error, setError] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [selectedVariantIndex, setSelectedVariantIndex] = useState(0);
  const [selectedAttributes, setSelectedAttributes] = useState({});
  const [activeTab, setActiveTab] = useState(0);
  const navigate = useNavigate();
  const toast = useToast();
  const { isAuthenticated } = useAuth();
  const [isLargerThan768] = useMediaQuery('(min-width: 768px)');

  // Format price with VND currency
  const formatPrice = (price) => {
    if (!price) return '';
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(price);
  };

  // Calculate discount percentage
  const calculateDiscount = (originalPrice, discountPrice) => {
    if (!originalPrice || !discountPrice || originalPrice <= discountPrice) return 0;
    return Math.round(((originalPrice - discountPrice) / originalPrice) * 100);
  };

  // Tính toán danh sách các combinations có sẵn
  const availableCombinations = useMemo(() => {
    if (!product) return {};

    const combinations = {};

    // Lặp qua tất cả variants để xây dựng map các combinations hợp lệ
    product.product_variants.forEach(variant => {
      variant.attribute_values.forEach(attrValue => {
        const attributeName = attrValue.attribute_name;
        const attributeValue = attrValue.attribute_value;

        if (!combinations[attributeName]) {
          combinations[attributeName] = {};
        }

        if (!combinations[attributeName][attributeValue]) {
          combinations[attributeName][attributeValue] = new Set();
        }

        // Thêm tất cả các giá trị thuộc tính khác của variant này
        variant.attribute_values.forEach(otherAttr => {
          if (otherAttr.attribute_name !== attributeName) {
            combinations[attributeName][attributeValue].add(
                `${otherAttr.attribute_name}:${otherAttr.attribute_value}`
            );
          }
        });
      });
    });

    return combinations;
  }, [product]);

  // Cập nhật selectedAttributes khi selectedVariant thay đổi
  useEffect(() => {
    if (product && product.product_variants && product.product_variants.length > 0) {
      const variant = product.product_variants[selectedVariantIndex];
      const attributes = {};

      variant.attribute_values.forEach(attr => {
        attributes[attr.attribute_name] = attr.attribute_value;
      });

      setSelectedAttributes(attributes);
    }
  }, [product, selectedVariantIndex]);

  useEffect(() => {
    const fetchProductData = async () => {
      setIsLoadingProduct(true);
      try {
        const response = await productService.getProductById(id);
        const productData = response.data.data;

        // Find the default variant or use the first one
        const defaultVariantIndex = productData.product_variants.findIndex(
            (variant) => variant.is_default
        );
        setSelectedVariantIndex(defaultVariantIndex > -1 ? defaultVariantIndex : 0);

        setProduct(productData);
        setIsLoadingProduct(false);
      } catch (err) {
        console.error('Error fetching product details:', err);
        setError('Không thể tải thông tin sản phẩm. Vui lòng thử lại sau.');
        setIsLoadingProduct(false);
      }
    };

    const fetchProductReviews = async () => {
      setIsLoadingReviews(true);
      try {
        const response = await productService.getProductReviews(id);
        setReviews(response.data.data || []);
        setReviewMetadata(response.data.metadata.pagination || {});
        setIsLoadingReviews(false);
      } catch (err) {
        console.error('Error fetching product reviews:', err);
        setIsLoadingReviews(false);
      }
    };

    if (id) {
      fetchProductData();
      fetchProductReviews();
    }
  }, [id]);

  const handleReviewPageChange = async (page) => {
    setIsLoadingReviews(true);
    try {
      const response = await productService.getProductReviews(id, page);
      setReviews(response.data.data || []);
      setReviewMetadata(response.data.metadata.pagination || {});
      setIsLoadingReviews(false);
    } catch (err) {
      console.error('Error fetching product reviews:', err);
      setIsLoadingReviews(false);
    }
  };

  // Kiểm tra xem một giá trị thuộc tính có khả dụng không dựa trên các lựa chọn hiện tại
  const isAttributeValueAvailable = (attributeName, attributeValue) => {
    // Nếu chưa có lựa chọn nào, tất cả đều khả dụng
    if (Object.keys(selectedAttributes).length === 0) return true;

    // Nếu đây là thuộc tính đang được xem xét, luôn khả dụng
    if (selectedAttributes[attributeName] === attributeValue) return true;

    // Kiểm tra xem giá trị này có khả dụng với các lựa chọn khác không
    for (const [selectedAttrName, selectedAttrValue] of Object.entries(selectedAttributes)) {
      // Bỏ qua thuộc tính đang xem xét
      if (selectedAttrName === attributeName) continue;

      // Kiểm tra xem có tồn tại variant nào kết hợp giá trị hiện tại với giá trị đang xem xét không
      const key = `${attributeName}:${attributeValue}`;
      if (!availableCombinations[selectedAttrName] ||
          !availableCombinations[selectedAttrName][selectedAttrValue] ||
          !availableCombinations[selectedAttrName][selectedAttrValue].has(key)) {
        return false;
      }
    }

    return true;
  };

  // Tìm variant dựa trên các thuộc tính đã chọn
  const findVariantByAttributes = (attributes) => {
    return product.product_variants.findIndex(variant => {
      // Kiểm tra xem variant này có khớp với tất cả các thuộc tính đã chọn không
      return Object.entries(attributes).every(([name, value]) => {
        return variant.attribute_values.some(
            attr => attr.attribute_name === name && attr.attribute_value === value
        );
      });
    });
  };

  // Xử lý khi chọn một giá trị thuộc tính
  const handleAttributeSelect = (attributeName, attributeValue) => {
    // Cập nhật thuộc tính đã chọn
    const newSelectedAttributes = {
      ...selectedAttributes,
      [attributeName]: attributeValue
    };

    // Tìm variant phù hợp với các thuộc tính đã chọn
    const variantIndex = findVariantByAttributes(newSelectedAttributes);

    if (variantIndex !== -1) {
      setSelectedVariantIndex(variantIndex);
      // Cập nhật selectedAttributes từ variant mới (để đồng bộ tất cả các thuộc tính)
      setSelectedAttributes(newSelectedAttributes);
      // Reset quantity khi đổi variant
      setQuantity(1);
    }
  };

  const handleVariantSelect = (index) => {
    setSelectedVariantIndex(index);
    // Reset quantity khi đổi variant
    setQuantity(1);
  };

  const handleAddToCart = () => {
    toast({
      title: 'Thêm vào giỏ hàng',
      description: `Đã thêm ${quantity} ${product.name} vào giỏ hàng`,
      status: 'success',
      duration: 3000,
      isClosable: true,
    });
  };

  // Updated Buy Now handler to navigate to checkout
  const handleBuyNow = () => {
    if (!product || !selectedVariant) {
      toast({
        title: 'Lỗi',
        description: 'Không thể xác định sản phẩm',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    // Create order item from current selection
    const orderItem = {
      cart_item_id: `temp_${Date.now()}`, // Temporary ID for direct purchase
      product_id: product.product_id,
      product_name: product.name,
      product_variant_id: selectedVariant.product_variant_id,
      product_variant_thumbnail: selectedVariant.thumbnail_url,
      variant_name: selectedVariant.variant_name,
      price: selectedVariant.price,
      discount_price: selectedVariant.discount_price || 0,
      quantity: quantity,
      attribute_values: selectedVariant.attribute_values || []
    };

    // Navigate to checkout with product data
    navigate('/checkout', {
      state: {
        cartItems: [orderItem],
        fromProductDetail: true,
        selectedVoucher: null,
        voucherDiscount: 0,
        finalTotal: (selectedVariant.discount_price > 0 ? selectedVariant.discount_price : selectedVariant.price) * quantity
      }
    });
  };

  // Loading skeleton for product details
  if (isLoadingProduct) {
    return (
        <Container maxW='container.xl' py={8}>
          <Grid templateColumns={{ base: '1fr', md: '1fr 1fr' }} gap={10}>
            <GridItem>
              <Skeleton height="400px" borderRadius="md" mb={4} />
              <SimpleGrid columns={4} spacing={4}>
                {[...Array(4)].map((_, i) => (
                    <Skeleton key={i} height="80px" borderRadius="md" />
                ))}
              </SimpleGrid>
            </GridItem>
            <GridItem>
              <VStack align="stretch" spacing={6}>
                <Skeleton height="40px" width="80%" />
                <Skeleton height="24px" width="40%" />
                <Skeleton height="40px" width="60%" />
                <Divider />
                <Skeleton height="80px" />
                <Skeleton height="40px" />
                <Skeleton height="40px" />
                <Skeleton height="60px" />
              </VStack>
            </GridItem>
          </Grid>
        </Container>
    );
  }

  // Error state
  if (error || !product) {
    return (
        <Container maxW='container.xl' py={10}>
          <Box textAlign='center'>
            <Heading size="md" color="red.500" mb={4}>
              {error || 'Không thể tải thông tin sản phẩm'}
            </Heading>
            <Button as={RouterLink} to="/products" colorScheme="brand">
              Quay lại danh sách sản phẩm
            </Button>
          </Box>
        </Container>
    );
  }

  // Get current selected variant
  const selectedVariant = product.product_variants[selectedVariantIndex];
  const discountPercent = calculateDiscount(selectedVariant.price, selectedVariant.discount_price);

  // Update page title
  return (
      <Container maxW='container.xl' py={8}>
        {/* Page Title (for SEO) */}
        <PageTitle title={product.name} />

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
          {product.category_name && (
              <BreadcrumbItem>
                <BreadcrumbLink as={RouterLink} to={`/products?category_ids=${product.category_id}`}>
                  {product.category_name}
                </BreadcrumbLink>
              </BreadcrumbItem>
          )}
          <BreadcrumbItem isCurrentPage>
            <BreadcrumbLink>{product.name}</BreadcrumbLink>
          </BreadcrumbItem>
        </Breadcrumb>

        {/* Section 1: Product Info */}
        <Grid templateColumns={{ base: '1fr', md: '1fr 1fr' }} gap={10} mb={8}>
          {/* Product Images */}
          <GridItem>
            <Box position='relative' mb={4}>
              {selectedVariant.discount_price > 0 && (
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
                  src={selectedVariant.thumbnail_url}
                  alt={selectedVariant.alt_text_thumbnail || product.name}
                  borderRadius='md'
                  width='100%'
                  height='auto'
                  objectFit='cover'
              />
            </Box>
            <SimpleGrid columns={4} spacing={4}>
              {product.product_variants.map((variant, index) => (
                  <Box
                      key={variant.product_variant_id}
                      borderWidth={selectedVariantIndex === index ? '2px' : '1px'}
                      borderColor={selectedVariantIndex === index ? 'brand.500' : 'gray.200'}
                      borderRadius='md'
                      overflow='hidden'
                      cursor='pointer'
                      onClick={() => handleVariantSelect(index)}
                  >
                    <Image
                        src={variant.thumbnail_url}
                        alt={variant.alt_text_thumbnail || variant.variant_name}
                    />
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
                                key={`rating-star-${i}`}
                                color={
                                  i < Math.round(product.average_rating)
                                      ? 'yellow.400'
                                      : 'gray.300'
                                }
                            />
                        ))}
                    <Text ml={2} color='gray.600'>
                      {product.average_rating.toFixed(1)} - {product.total_reviews} đánh giá
                    </Text>
                  </Box>
                </HStack>

                <Box mb={6}>
                  <Flex align='baseline'>
                    <Heading as='h2' size='xl' color='brand.500' mr={3}>
                      {selectedVariant.discount_price > 0
                          ? formatPrice(selectedVariant.discount_price)
                          : formatPrice(selectedVariant.price)}
                    </Heading>
                    {selectedVariant.discount_price > 0 && (
                        <Text as='s' color='gray.500' fontSize='lg'>
                          {formatPrice(selectedVariant.price)}
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

              {/* Description - Short version */}
              <Box>
                <Text color='gray.700' noOfLines={3}>{product.description}</Text>
              </Box>

              <Divider />

              {/* Product Attributes/Variants */}
              {product.attributes && product.attributes.length > 0 && (
                  <Box>
                    {product.attributes.map((attribute) => (
                        <Box key={attribute.attribute_id} mb={4}>
                          <Text fontWeight='semibold' mb={2}>
                            {attribute.name}: {selectedAttributes[attribute.name]}
                          </Text>
                          <HStack spacing={2} flexWrap="wrap">
                            {attribute.values.map((value) => {
                              // Xác định xem giá trị này có khả dụng không
                              const isAvailable = isAttributeValueAvailable(attribute.name, value.value);

                              // Kiểm tra xem giá trị này có được chọn không
                              const isSelected = selectedAttributes[attribute.name] === value.value;

                              return (
                                  <Button
                                      key={`${attribute.attribute_id}-${value.option_id}`}
                                      size='sm'
                                      variant={isSelected ? 'solid' : 'outline'}
                                      colorScheme={isSelected ? 'brand' : 'gray'}
                                      onClick={() => isAvailable && handleAttributeSelect(attribute.name, value.value)}
                                      opacity={isAvailable ? 1 : 0.4}
                                      cursor={isAvailable ? 'pointer' : 'not-allowed'}
                                      _hover={{
                                        bg: isAvailable
                                            ? (isSelected ? 'brand.600' : 'gray.100')
                                            : 'transparent',
                                      }}
                                  >
                                    {value.value}
                                  </Button>
                              );
                            })}
                          </HStack>
                        </Box>
                    ))}
                  </Box>
              )}

              {/* Quantity */}
              <Box>
                <Text fontWeight='semibold' mb={2}>
                  Số lượng:
                </Text>
                <Flex align="center">
                  <NumberInput
                      max={selectedVariant.quantity}
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
                  <Text ml={4} fontSize='sm' color='gray.500'>
                    {selectedVariant.quantity} sản phẩm có sẵn
                  </Text>
                </Flex>
              </Box>

              {/* Actions */}
              <HStack spacing={4} pt={4}>
                <AddToCartButton
                    productId={product.product_id}
                    variantId={selectedVariant.product_variant_id}
                    quantity={quantity}
                    colorScheme="brand"
                    variant="outline"
                    size="lg"
                    flex="1"
                    isDisabled={selectedVariant.quantity <= 0}
                    onSuccess={() => toast({
                      title: 'Thêm vào giỏ hàng',
                      description: `Đã thêm ${quantity} ${product.name} vào giỏ hàng`,
                      status: 'success',
                      duration: 3000,
                      isClosable: true,
                    })}
                >
                  Thêm vào giỏ
                </AddToCartButton>
                <Button
                    colorScheme="brand"
                    size="lg"
                    flex="1"
                    onClick={handleBuyNow}
                    isDisabled={selectedVariant.quantity <= 0}
                >
                  Mua ngay
                </Button>
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
                        {selectedVariant.sku}
                      </Text>
                    </Box>
                  </Flex>
                </SimpleGrid>
              </Box>
            </VStack>
          </GridItem>
        </Grid>

        {/* Section 2: Shop Info */}
        {product.supplier && (
            <Box borderWidth="1px" borderRadius="md" p={4} mb={8} bg="white">
              <Flex align="center">
                <Avatar
                    size="md"
                    name={product.supplier.company_name}
                    src={product.supplier.thumbnail}
                    mr={4}
                />
                <Box>
                  <Text fontWeight="bold">{product.supplier.company_name}</Text>
                  <Text fontSize="sm" color="gray.600">{product.supplier.contact_phone}</Text>
                </Box>
                <Flex ml="auto">
                  <Button leftIcon={<FaStore />} colorScheme="brand" size="sm" mr={2}>
                    Xem Shop
                  </Button>
                  <Button variant="outline" colorScheme="brand" size="sm">
                    Chat Ngay
                  </Button>
                </Flex>
              </Flex>
            </Box>
        )}

        {/* Section 3 & 4: Tabs for Description and Reviews */}
        <Box mb={16}>
          <Tabs colorScheme='brand' defaultIndex={0} onChange={(index) => setActiveTab(index)}>
            <TabList>
              <Tab fontWeight='semibold'>Chi tiết sản phẩm</Tab>
              <Tab fontWeight='semibold'>Đánh giá ({product.total_reviews})</Tab>
            </TabList>

            <TabPanels>
              {/* Section 3: Description Tab */}
              <TabPanel p={4}>
                <Box>
                  <Heading as='h3' size='md' mb={4}>
                    Mô tả sản phẩm
                  </Heading>
                  <Text mb={6} whiteSpace="pre-wrap">{product.description}</Text>

                  {/* Product tags as hashtags */}
                  {product.product_tags && product.product_tags.length > 0 && (
                      <Box mb={4} mt={8}>
                        <Heading as="h4" size="sm" mb={3}>Tags:</Heading>
                        <HStack spacing={2} flexWrap="wrap">
                          {product.product_tags.map((tag, index) => (
                              <Badge key={`product-tag-${index}-${tag}`} colorScheme="blue" mr={2} px={2} py={1} borderRadius="md">
                                #{tag}
                              </Badge>
                          ))}
                        </HStack>
                      </Box>
                  )}
                </Box>
              </TabPanel>

              {/* Section 4: Reviews Tab */}
              <TabPanel p={4}>
                <Box>
                  <Heading as="h3" size="md" mb={6}>Đánh giá từ khách hàng</Heading>

                  {/* Average rating */}
                  <Flex
                      direction={{ base: 'column', md: 'row' }}
                      bg="gray.50"
                      p={4}
                      borderRadius="md"
                      mb={6}
                      align="center"
                  >
                    <Box textAlign={{ base: 'center', md: 'left' }} mb={{ base: 4, md: 0 }}>
                      <Text fontSize="3xl" fontWeight="bold">
                        {product.average_rating.toFixed(1)}/5
                      </Text>
                      <HStack justify={{ base: 'center', md: 'flex-start' }}>
                        {Array(5)
                            .fill('')
                            .map((_, i) => (
                                <StarIcon
                                    key={`tab-rating-star-${i}`}
                                    color={i < Math.round(product.average_rating) ? 'yellow.400' : 'gray.300'}
                                />
                            ))}
                      </HStack>
                      <Text color="gray.600" mt={1}>
                        {product.total_reviews} đánh giá
                      </Text>
                    </Box>
                  </Flex>

                  {/* Reviews list */}
                  {isLoadingReviews ? (
                      <VStack spacing={4} align="stretch">
                        {[1, 2, 3].map((_, index) => (
                            <Box key={index} p={4} borderWidth="1px" borderRadius="md">
                              <Flex mb={4}>
                                <Skeleton borderRadius="full" size="40px" mr={4} />
                                <Box flex="1">
                                  <Skeleton height="20px" width="40%" mb={2} />
                                  <Skeleton height="14px" width="30%" />
                                </Box>
                              </Flex>
                              <Skeleton height="16px" width="80%" mb={2} />
                              <Skeleton height="16px" width="60%" />
                            </Box>
                        ))}
                      </VStack>
                  ) : reviews.length > 0 ? (
                      <VStack spacing={4} align="stretch">
                        {reviews.map((review) => (
                            <Box key={review.id} p={4} borderWidth="1px" borderRadius="md">
                              <Flex mb={4}>
                                <Avatar
                                    size="sm"
                                    name={review.user_name}
                                    src={review.user_avatar_url}
                                    mr={4}
                                />
                                <Box flex="1">
                                  <Flex justify="space-between" align="center">
                                    <Text fontWeight="bold">{review.user_name}</Text>
                                    <Text fontSize="sm" color="gray.500">
                                      {new Date(review.created_at).toLocaleDateString('vi-VN')}
                                    </Text>
                                  </Flex>
                                  <HStack spacing={1}>
                                    {Array(5)
                                        .fill('')
                                        .map((_, i) => (
                                            <StarIcon
                                                key={`review-${review.id}-star-${i}`}
                                                size="sm"
                                                color={i < Math.round(review.rating) ? 'yellow.400' : 'gray.300'}
                                            />
                                        ))}
                                    <Text fontSize="sm" color="gray.500" ml={1}>
                                      {review.rating.toFixed(1)}
                                    </Text>
                                  </HStack>
                                </Box>
                              </Flex>
                              <Text>{review.comment}</Text>
                              {review.helpful_votes > 0 && (
                                  <Flex align="center" mt={2}>
                                    <Icon as={FaThumbsUp} color="gray.500" mr={1} />
                                    <Text fontSize="sm" color="gray.500">
                                      {review.helpful_votes} người thấy hữu ích
                                    </Text>
                                  </Flex>
                              )}
                            </Box>
                        ))}

                        {/* Pagination */}
                        {reviewMetadata && reviewMetadata.total_pages > 1 && (
                            <Box mt={6} display="flex" justifyContent="center" alignItems="center">
                              <Button
                                  variant="outline"
                                  onClick={() => handleReviewPageChange(reviewMetadata.page - 1)}
                                  isDisabled={reviewMetadata.page <= 1}
                                  mr={2}
                                  size="sm"
                              >
                                Trang trước
                              </Button>
                              <Text mx={4} fontSize="sm" color="gray.600">
                                Trang {reviewMetadata.page} / {reviewMetadata.total_pages}
                              </Text>
                              <Button
                                  variant="outline"
                                  onClick={() => handleReviewPageChange(reviewMetadata.page + 1)}
                                  isDisabled={reviewMetadata.page >= reviewMetadata.total_pages}
                                  ml={2}
                                  size="sm"
                              >
                                Trang sau
                              </Button>
                            </Box>
                        )}
                      </VStack>
                  ) : (
                      <Center py={10}>
                        <VStack>
                          <Text color="gray.500">Chưa có đánh giá nào cho sản phẩm này</Text>
                          <Text color="gray.400" fontSize="sm">Đánh giá chỉ có thể viết sau khi mua hàng</Text>
                        </VStack>
                      </Center>
                  )}
                </Box>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </Box>
      </Container>
  );
};

export default ProductDetail;