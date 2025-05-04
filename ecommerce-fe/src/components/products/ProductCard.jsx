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
  AspectRatio
} from '@chakra-ui/react';
import { FaStar, FaHeart, FaShoppingCart } from 'react-icons/fa';
import { Link as RouterLink } from 'react-router-dom';

// Format price with VND currency
const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price);
};

// Calculate discount percentage
const calculateDiscount = (originalPrice, currentPrice) => {
  if (!originalPrice || !currentPrice || originalPrice <= currentPrice) return null;
  const discount = Math.round(((originalPrice - currentPrice) / originalPrice) * 100);
  return discount;
};

const ProductCard = ({ product }) => {
  const toast = useToast();

  // Calculate discount if available
  const discountPercentage = calculateDiscount(product.originalPrice, product.price);

  // Handle quick actions
  const handleAddToCart = (e) => {
    e.preventDefault();
    e.stopPropagation();

    toast({
      title: 'Đã thêm vào giỏ hàng',
      description: `${product.name} đã được thêm vào giỏ hàng của bạn.`,
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
      description: `${product.name} đã được thêm vào danh sách yêu thích của bạn.`,
      status: 'success',
      duration: 3000,
      isClosable: true,
    });
  };

  return (
      <Box
          as={RouterLink}
          to={`/products/${product.id}`}
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
              src={product.image}
              alt={product.name}
              objectFit="contain"
              bg="gray.50"
          />
        </AspectRatio>

        {/* Product info */}
        <Box p={4} flex="1" display="flex" flexDirection="column">
          <Text fontSize="sm" color="gray.500" mb={1}>
            {product.brand || 'Thương hiệu'}
          </Text>

          <Text
              fontWeight="semibold"
              fontSize="md"
              mb={2}
              noOfLines={2}
              flex="1"
          >
            {product.name}
          </Text>

          {/* Rating */}
          <Flex align="center" mb={2}>
            <Flex align="center">
              <Icon as={FaStar} color="yellow.400" mr={1} boxSize={3} />
              <Text fontSize="sm" fontWeight="medium">
                {product.rating}
              </Text>
            </Flex>
            <Text fontSize="xs" color="gray.500" ml={1}>
              ({product.reviewCount} đánh giá)
            </Text>
          </Flex>

          {/* Price */}
          <Flex align="baseline">
            <Text fontWeight="bold" fontSize="md" color="brand.500">
              {formatPrice(product.price)}
            </Text>
            {product.originalPrice && (
                <Text
                    fontSize="sm"
                    color="gray.500"
                    textDecoration="line-through"
                    ml={2}
                >
                  {formatPrice(product.originalPrice)}
                </Text>
            )}
          </Flex>
        </Box>
      </Box>
  );
};

export default ProductCard;