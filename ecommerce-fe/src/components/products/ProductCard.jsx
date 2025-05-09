import React from 'react';
import {
  Box,
  Image,
  Text,
  Badge,
  Flex,
  Icon,
  AspectRatio,
  HStack
} from '@chakra-ui/react';
import { FaStar, FaRegStar } from 'react-icons/fa';
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
  // Extract values from API data
  const id = product.product_id;
  const name = product.product_name;
  const image = product.product_thumbnail;
  const rating = product.product_average_rating;
  const reviewCount = product.product_total_reviews;
  const regularPrice = product.product_price;
  const discountPrice = product.product_discount_price > 0 ? product.product_discount_price : null;

  // Calculate discount percentage
  const discountPercentage = calculateDiscount(regularPrice, discountPrice);

  // Generate star rating
  const renderStars = () => {
    const stars = [];
    const ratingValue = rating || 0;

    for (let i = 1; i <= 5; i++) {
      if (i <= Math.round(ratingValue)) {
        stars.push(<Icon key={i} as={FaStar} color="yellow.400" boxSize={3} />);
      } else {
        stars.push(<Icon key={i} as={FaRegStar} color="yellow.400" boxSize={3} />);
      }
    }

    return stars;
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

        {/* Product image */}
        <AspectRatio ratio={1} w="100%">
          <Image
              src={image}
              alt={name}
              objectFit="contain"
              bg="gray.50"
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
              {rating ? rating.toFixed(1) : '0'} ({reviewCount || 0})
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