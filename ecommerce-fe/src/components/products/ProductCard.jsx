import { StarIcon } from '@chakra-ui/icons';
import {
  Badge,
  Box,
  Flex,
  HStack,
  IconButton,
  Image,
  Text,
  Tooltip,
  useToast,
} from '@chakra-ui/react';
import { useState } from 'react';
import { FaEye, FaHeart, FaRegHeart, FaShoppingCart } from 'react-icons/fa';
import { Link as RouterLink } from 'react-router-dom';
import useAuth from '../../hooks/useAuth';

const ProductCard = ({ product }) => {
  const [isLiked, setIsLiked] = useState(false);
  const toast = useToast();
  const { isAuthenticated } = useAuth();

  const formatPrice = (price) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND',
    }).format(price);
  };

  const calculateDiscount = (originalPrice, currentPrice) => {
    if (!originalPrice || originalPrice <= currentPrice) return null;
    const discount = Math.round(
      ((originalPrice - currentPrice) / originalPrice) * 100,
    );
    return discount > 0 ? `-${discount}%` : null;
  };

  const handleAddToCart = (e) => {
    e.preventDefault();
    e.stopPropagation();

    // TODO: Add to cart logic
    toast({
      title: 'Thêm vào giỏ hàng',
      description: `Đã thêm ${product.name} vào giỏ hàng`,
      status: 'success',
      duration: 3000,
      isClosable: true,
    });
  };

  const handleToggleFavorite = (e) => {
    e.preventDefault();
    e.stopPropagation();

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

    // TODO: Add/remove from favorites logic
    toast({
      title: isLiked ? 'Đã xóa khỏi yêu thích' : 'Đã thêm vào yêu thích',
      description: isLiked
        ? `Đã xóa ${product.name} khỏi danh sách yêu thích`
        : `Đã thêm ${product.name} vào danh sách yêu thích`,
      status: 'success',
      duration: 3000,
      isClosable: true,
    });
  };

  const discountBadge = calculateDiscount(product.originalPrice, product.price);

  return (
    <Box
      as={RouterLink}
      to={`/products/${product.id}`}
      borderWidth='1px'
      borderRadius='lg'
      overflow='hidden'
      bg='white'
      transition='transform 0.3s, box-shadow 0.3s'
      _hover={{
        transform: 'translateY(-5px)',
        boxShadow: 'lg',
        textDecoration: 'none',
      }}
      position='relative'
      h='100%'
      display='flex'
      flexDirection='column'
    >
      {/* Discount Badge */}
      {discountBadge && (
        <Badge
          position='absolute'
          top='10px'
          left='10px'
          colorScheme='red'
          variant='solid'
          borderRadius='md'
          px={2}
          py={1}
          fontSize='xs'
          zIndex='1'
        >
          {discountBadge}
        </Badge>
      )}

      {/* Product Image */}
      <Box position='relative' overflow='hidden'>
        <Image
          src={product.image}
          alt={product.name}
          h='200px'
          w='100%'
          objectFit='cover'
          transition='transform 0.5s'
          _groupHover={{ transform: 'scale(1.05)' }}
        />

        {/* Hover Actions */}
        <Flex
          position='absolute'
          bottom='10px'
          left='0'
          right='0'
          justifyContent='center'
          opacity='0'
          transition='opacity 0.3s'
          _groupHover={{ opacity: 1 }}
          zIndex='1'
        >
          <HStack spacing={2}>
            <Tooltip label='Thêm vào giỏ hàng' placement='top'>
              <IconButton
                aria-label='Add to cart'
                icon={<FaShoppingCart />}
                colorScheme='brand'
                size='sm'
                borderRadius='full'
                onClick={handleAddToCart}
              />
            </Tooltip>

            <Tooltip
              label={isLiked ? 'Xóa khỏi yêu thích' : 'Thêm vào yêu thích'}
              placement='top'
            >
              <IconButton
                aria-label='Toggle favorite'
                icon={isLiked ? <FaHeart /> : <FaRegHeart />}
                colorScheme={isLiked ? 'red' : 'gray'}
                size='sm'
                borderRadius='full'
                onClick={handleToggleFavorite}
              />
            </Tooltip>

            <Tooltip label='Xem nhanh' placement='top'>
              <IconButton
                as={RouterLink}
                to={`/products/${product.id}`}
                aria-label='Quick view'
                icon={<FaEye />}
                colorScheme='gray'
                size='sm'
                borderRadius='full'
                onClick={(e) => e.stopPropagation()}
              />
            </Tooltip>
          </HStack>
        </Flex>
      </Box>

      {/* Product Info */}
      <Box p={4} flex='1' display='flex' flexDirection='column'>
        <Text
          fontWeight='semibold'
          as='h3'
          lineHeight='tight'
          noOfLines={2}
          mb={2}
          flex='1'
        >
          {product.name}
        </Text>

        {/* Rating */}
        <Flex alignItems='center' mb={2}>
          <HStack spacing={1}>
            {Array(5)
              .fill('')
              .map((_, i) => (
                <StarIcon
                  key={i}
                  color={
                    i < Math.floor(product.rating) ? 'yellow.400' : 'gray.300'
                  }
                  size='sm'
                />
              ))}
            {product.rating % 1 >= 0.5 && (
              <StarIcon color='yellow.400' size='sm' />
            )}
          </HStack>
          <Text ml={2} fontSize='sm' color='gray.600'>
            ({product.reviewCount})
          </Text>
        </Flex>

        {/* Price */}
        <Box>
          <Flex align='baseline'>
            <Text fontWeight='bold' fontSize='lg' color='brand.500'>
              {formatPrice(product.price)}
            </Text>
            {product.originalPrice && product.originalPrice > product.price && (
              <Text as='s' color='gray.500' fontSize='sm' ml={2}>
                {formatPrice(product.originalPrice)}
              </Text>
            )}
          </Flex>
        </Box>
      </Box>
    </Box>
  );
};

export default ProductCard;
