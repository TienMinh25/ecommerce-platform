import { Center, SimpleGrid, Spinner, Text } from '@chakra-ui/react';
import ProductCard from './ProductCard';

const ProductGrid = ({ products, isLoading, error }) => {
  if (isLoading) {
    return (
      <Center py={10}>
        <Spinner
          thickness='4px'
          speed='0.65s'
          emptyColor='gray.200'
          color='brand.500'
          size='xl'
        />
      </Center>
    );
  }

  if (error) {
    return (
      <Center py={10}>
        <Text color='red.500'>
          {error.message || 'Có lỗi xảy ra khi tải danh sách sản phẩm.'}
        </Text>
      </Center>
    );
  }

  if (!products || products.length === 0) {
    return (
      <Center py={10}>
        <Text color='gray.500'>Không tìm thấy sản phẩm nào.</Text>
      </Center>
    );
  }

  return (
    <SimpleGrid
      columns={{ base: 2, md: 3, lg: 4 }}
      spacing={{ base: 4, md: 6 }}
    >
      {products.map((product) => (
        <ProductCard key={product.id} product={product} />
      ))}
    </SimpleGrid>
  );
};

export default ProductGrid;
