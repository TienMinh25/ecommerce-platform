import React from 'react';
import {
  Box,
  Image,
  Text,
  Badge,
  Flex,
  Icon,
  IconButton,
  useToast,
  AspectRatio,
  HStack
} from '@chakra-ui/react';
import { FaStar, FaRegStar, FaHeart, FaShoppingCart } from 'react-icons/fa';
import { Link as RouterLink } from 'react-router-dom';

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
  if (!originalPrice || !discountPrice || originalPrice <= discountPrice) return null;
  const discount = Math.round(((originalPrice - discountPrice) / originalPrice) * 100);
  return discount;
};

const ProductCard = ({ product }) => {
  const toast = useToast();

  // Check if product is from API or mock data
  const isFromApi = !!product.product_id;

  // Extract values based on data source
  const id = isFromApi ? product.product_id : product.id;
  const name = isFromApi ? product.product_name : product.name;
  const image = isFromApi ? product.product_thumbnail : product.image;
  const rating = isFromApi ? product.product_average_rating : product.rating;
  const reviewCount = isFromApi ? product.product_total_reviews : product.reviewCount;

  // Price handling
  let regularPrice, discountPrice, displayPrice;

  if (isFromApi) {
    // API data
    regularPrice = product.product_price;
    discountPrice = product.product_discount_price > 0 ? product.product_discount_price : null;
    displayPrice = discountPrice || regularPrice;
  } else {
    // Mock data
    regularPrice = product.originalPrice || product.price;
    discountPrice = product.price < product.originalPrice ? product.price : null;
    displayPrice = product.price;
  }

  // Calculate discount percentage
  const discountPercentage = calculateDiscount(regularPrice, discountPrice);

  // Generate star rating
  const renderStars = () => {
    const stars = [];
    const ratingValue = rating || 0;

    for (let i = 1; i <= 5; i++) {
      if (i <= ratingValue) {
        stars.push(<Icon key={i} as={FaStar} color="yellow.400" boxSize={3} />);
      } else {
        stars.push(<Icon key={i} as={FaRegStar} color="yellow.400" boxSize={3} />);
      }
    }

    return stars;
  };

  // Handle quick actions
  const handleAddToCart = (e) => {
    e.preventDefault();
    e.stopPropagation();

    toast({
      title: 'Đã thêm vào giỏ hàng',
      description: `${name} đã được thêm vào giỏ hàng của bạn.`,
      status: 'success',
      duration: 3000,
      isClosable: true,
    });
  };

  const handleAddToWishlist = (e) => {
    e.preventDefault();
    e.stopPropagation();

    toast({
      title: 'Đã thêm vào danh sách yêu thích',
      description: `${name} đã được thêm vào danh sách yêu thích của bạn.`,
      status: 'success',
      duration: 3000,
      isClosable: true,
    });
  };

  return (
      <Box
          as={RouterLink}
          to={`/products/${id}`}
          borderWidth="1px"
          borderRadius="lg"
          overflow="hidden"
          bg="white"
          transition="transform 0.3s, box-shadow 0.3s"
          _hover={{
            transform: 'translateY(-5px)',
            boxShadow: 'lg',
          }}
          position="relative"
          height="100%"
          display="flex"
          flexDirection="column"
      >
        {/* Discount badge if available */}
        {discountPercentage && (
            <Badge
                position="absolute"
                top={2}
                left={2}
                colorScheme="red"
                variant="solid"
                borderRadius="md"
                px={2}
                py={1}
                fontSize="xs"
                fontWeight="bold"
                zIndex={1}
            >
              -{discountPercentage}%
            </Badge>
        )}

        {/* Quick action buttons */}
        <Flex
            position="absolute"
            top={2}
            right={2}
            direction="column"
            gap={2}
            zIndex={1}
            opacity={0}
            transition="opacity 0.3s"
            _groupHover={{ opacity: 1 }}
        >
          <IconButton
              icon={<FaHeart />}
              onClick={handleAddToWishlist}
              aria-label="Add to wishlist"
              size="sm"
              borderRadius="full"
              colorScheme="pink"
              variant="solid"
              boxShadow="md"
          />
          <IconButton
              icon={<FaShoppingCart />}
              onClick={handleAddToCart}
              aria-label="Add to cart"
              size="sm"
              borderRadius="full"
              colorScheme="brand"
              variant="solid"
              boxShadow="md"
          />
        </Flex>

        {/* Product image */}
        <AspectRatio ratio={1} w="100%">
          <Image
              src={image}
              alt={name}
              objectFit="contain"
              bg="gray.50"
              fallbackSrc="https://via.placeholder.com/300x300?text=Product"
          />
        </AspectRatio>

        {/* Product info */}
        <Box p={4} flex="1" display="flex" flexDirection="column">
          <Text
              fontWeight="semibold"
              fontSize="md"
              mb={2}
              noOfLines={2}
              flex="1"
          >
            {name}
          </Text>

          {/* Rating with 5 stars */}
          <Flex align="center" mb={2}>
            <HStack spacing={0} mr={2}>
              {renderStars()}
            </HStack>
            <Text fontSize="xs" color="gray.500">
              ({reviewCount || 0})
            </Text>
          </Flex>

          {/* Price */}
          <Flex align="baseline" flexWrap="wrap">
            {/* Show discount price if available */}
            {discountPrice && (
                <Text fontWeight="bold" fontSize="md" color="red.500" mr={2}>
                  {formatPrice(discountPrice)}
                </Text>
            )}

            {/* Show regular price - with strikethrough if discounted */}
            <Text
                fontWeight={discountPrice ? "normal" : "bold"}
                fontSize={discountPrice ? "sm" : "md"}
                color={discountPrice ? "gray.500" : "brand.500"}
                textDecoration={discountPrice ? "line-through" : "none"}
            >
              {formatPrice(regularPrice)}
            </Text>
          </Flex>
        </Box>
      </Box>
  );
};

export default ProductCard;