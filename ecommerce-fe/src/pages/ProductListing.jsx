import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import {
  Box,
  Container,
  Heading,
  Text,
  Grid,
  GridItem,
  Flex,
  Button,
  Select,
  Input,
  InputGroup,
  InputRightElement,
  IconButton,
  Accordion,
  AccordionItem,
  AccordionButton,
  AccordionPanel,
  AccordionIcon,
  Checkbox,
  Stack,
  RangeSlider,
  RangeSliderTrack,
  RangeSliderFilledTrack,
  RangeSliderThumb,
  Drawer,
  DrawerBody,
  DrawerFooter,
  DrawerHeader,
  DrawerOverlay,
  DrawerContent,
  DrawerCloseButton,
  useDisclosure,
  useBreakpointValue,
  Badge,
  HStack,
} from '@chakra-ui/react';
import { SearchIcon, HamburgerIcon } from '@chakra-ui/icons';
import ProductGrid from '../components/products/ProductGrid';

const ProductListing = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [products, setProducts] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [priceRange, setPriceRange] = useState([0, 5000000]);
  const [sortOption, setSortOption] = useState('popular');
  const [searchQuery, setSearchQuery] = useState(
    searchParams.get('search') || '',
  );
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [selectedBrands, setSelectedBrands] = useState([]);

  const { isOpen, onOpen, onClose } = useDisclosure();
  const isMobile = useBreakpointValue({ base: true, md: false });

  // Sample categories for filter
  const categories = [
    { id: 'men', name: 'Thời trang nam' },
    { id: 'women', name: 'Thời trang nữ' },
    { id: 'kids', name: 'Thời trang trẻ em' },
    { id: 'shoes', name: 'Giày dép' },
    { id: 'bags', name: 'Túi xách' },
    { id: 'watches', name: 'Đồng hồ' },
    { id: 'accessories', name: 'Phụ kiện' },
  ];

  // Sample brands for filter
  const brands = [
    { id: 'nike', name: 'Nike' },
    { id: 'adidas', name: 'Adidas' },
    { id: 'puma', name: 'Puma' },
    { id: 'reebok', name: 'Reebok' },
    { id: 'vans', name: 'Vans' },
    { id: 'converse', name: 'Converse' },
    { id: 'gucci', name: 'Gucci' },
    { id: 'zara', name: 'Zara' },
    { id: 'hm', name: 'H&M' },
  ];

  useEffect(() => {
    const fetchProducts = async () => {
      setIsLoading(true);
      try {
        // TODO: Implement actual API call
        // For now, use dummy data
        const dummyProducts = [
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
            name: 'Áo khoác jean nữ',
            image: 'https://via.placeholder.com/300x300?text=Jacket',
            price: 650000,
            originalPrice: 800000,
            rating: 4.6,
            reviewCount: 112,
          },
          {
            id: 10,
            name: 'Quần short nam',
            image: 'https://via.placeholder.com/300x300?text=Shorts',
            price: 350000,
            originalPrice: 400000,
            rating: 4.3,
            reviewCount: 87,
          },
          {
            id: 11,
            name: 'Áo hoodie unisex',
            image: 'https://via.placeholder.com/300x300?text=Hoodie',
            price: 550000,
            originalPrice: 650000,
            rating: 4.7,
            reviewCount: 143,
          },
          {
            id: 12,
            name: 'Váy midi nữ',
            image: 'https://via.placeholder.com/300x300?text=Skirt',
            price: 450000,
            originalPrice: 550000,
            rating: 4.5,
            reviewCount: 95,
          },
        ];

        // Simulate network delay
        setTimeout(() => {
          setProducts(dummyProducts);
          setIsLoading(false);
        }, 1000);
      } catch (error) {
        setError(error);
        setIsLoading(false);
      }
    };

    fetchProducts();
  }, [searchParams]);

  const handleSearch = (e) => {
    e.preventDefault();

    // Update search params
    const params = new URLSearchParams(searchParams);
    if (searchQuery) {
      params.set('search', searchQuery);
    } else {
      params.delete('search');
    }
    setSearchParams(params);
  };

  const handleSortChange = (e) => {
    setSortOption(e.target.value);
    // TODO: Implement sorting logic
  };

  const handlePriceRangeChange = (values) => {
    setPriceRange(values);
  };

  const handleCategoryToggle = (categoryId) => {
    setSelectedCategories((prev) => {
      if (prev.includes(categoryId)) {
        return prev.filter((id) => id !== categoryId);
      } else {
        return [...prev, categoryId];
      }
    });
  };

  const handleBrandToggle = (brandId) => {
    setSelectedBrands((prev) => {
      if (prev.includes(brandId)) {
        return prev.filter((id) => id !== brandId);
      } else {
        return [...prev, brandId];
      }
    });
  };

  const handleResetFilters = () => {
    setPriceRange([0, 5000000]);
    setSelectedCategories([]);
    setSelectedBrands([]);
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND',
      maximumFractionDigits: 0,
    }).format(price);
  };

  // Filter section for desktop view
  const FilterSection = () => (
    <Stack spacing={6}>
      {/* Price Range */}
      <Box>
        <Heading as='h3' size='sm' mb={4}>
          Khoảng giá
        </Heading>
        <RangeSlider
          min={0}
          max={5000000}
          step={100000}
          value={priceRange}
          onChange={handlePriceRangeChange}
          mb={4}
          colorScheme='brand'
        >
          <RangeSliderTrack>
            <RangeSliderFilledTrack />
          </RangeSliderTrack>
          <RangeSliderThumb index={0} />
          <RangeSliderThumb index={1} />
        </RangeSlider>
        <Flex justify='space-between'>
          <Text>{formatPrice(priceRange[0])}</Text>
          <Text>{formatPrice(priceRange[1])}</Text>
        </Flex>
      </Box>

      {/* Categories */}
      <Box>
        <Accordion allowToggle defaultIndex={[0]}>
          <AccordionItem border='none'>
            <AccordionButton px={0} _hover={{ bg: 'transparent' }}>
              <Box as='h3' flex='1' textAlign='left' fontWeight='semibold'>
                Danh mục
              </Box>
              <AccordionIcon />
            </AccordionButton>
            <AccordionPanel px={0}>
              <Stack spacing={2}>
                {categories.map((category) => (
                  <Checkbox
                    key={category.id}
                    isChecked={selectedCategories.includes(category.id)}
                    onChange={() => handleCategoryToggle(category.id)}
                    colorScheme='brand'
                  >
                    {category.name}
                  </Checkbox>
                ))}
              </Stack>
            </AccordionPanel>
          </AccordionItem>
        </Accordion>
      </Box>

      {/* Brands */}
      <Box>
        <Accordion allowToggle defaultIndex={[0]}>
          <AccordionItem border='none'>
            <AccordionButton px={0} _hover={{ bg: 'transparent' }}>
              <Box as='h3' flex='1' textAlign='left' fontWeight='semibold'>
                Thương hiệu
              </Box>
              <AccordionIcon />
            </AccordionButton>
            <AccordionPanel px={0}>
              <Stack spacing={2}>
                {brands.map((brand) => (
                  <Checkbox
                    key={brand.id}
                    isChecked={selectedBrands.includes(brand.id)}
                    onChange={() => handleBrandToggle(brand.id)}
                    colorScheme='brand'
                  >
                    {brand.name}
                  </Checkbox>
                ))}
              </Stack>
            </AccordionPanel>
          </AccordionItem>
        </Accordion>
      </Box>

      {/* Reset Filters Button */}
      <Button
        variant='outline'
        colorScheme='gray'
        size='sm'
        onClick={handleResetFilters}
      >
        Xóa bộ lọc
      </Button>
    </Stack>
  );

  return (
    <Container maxW='container.xl' py={8}>
      <Box mb={8}>
        <Heading as='h1' size='xl' mb={2}>
          Sản phẩm
        </Heading>
        <Text color='gray.600'>
          Khám phá bộ sưu tập sản phẩm đa dạng của chúng tôi
        </Text>
      </Box>

      {/* Search and Filters */}
      <Grid templateColumns={{ base: '1fr', md: '1fr 3fr' }} gap={8}>
        {/* Filter Sidebar - Desktop */}
        {!isMobile && (
          <GridItem>
            <FilterSection />
          </GridItem>
        )}

        {/* Products */}
        <GridItem>
          {/* Search and Sort Bar */}
          <Box mb={6}>
            <Grid
              templateColumns={{ base: '1fr', md: '2fr 1fr' }}
              gap={4}
              alignItems='center'
            >
              <form onSubmit={handleSearch}>
                <InputGroup>
                  <Input
                    placeholder='Tìm kiếm sản phẩm...'
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    bg='white'
                  />
                  <InputRightElement>
                    <IconButton
                      aria-label='Search'
                      icon={<SearchIcon />}
                      type='submit'
                      variant='ghost'
                      colorScheme='brand'
                    />
                  </InputRightElement>
                </InputGroup>
              </form>

              <Flex justify='space-between' align='center'>
                {isMobile && (
                  <Button
                    leftIcon={<HamburgerIcon />}
                    onClick={onOpen}
                    variant='outline'
                    colorScheme='brand'
                    size='md'
                  >
                    Lọc
                  </Button>
                )}

                <Box flex='1'>
                  <Select
                    value={sortOption}
                    onChange={handleSortChange}
                    bg='white'
                  >
                    <option value='popular'>Phổ biến</option>
                    <option value='newest'>Mới nhất</option>
                    <option value='price-asc'>Giá: Thấp đến cao</option>
                    <option value='price-desc'>Giá: Cao đến thấp</option>
                    <option value='rating'>Đánh giá</option>
                  </Select>
                </Box>
              </Flex>
            </Grid>
          </Box>

          {/* Active Filters */}
          {(selectedCategories.length > 0 || selectedBrands.length > 0) && (
            <Box mb={6}>
              <HStack spacing={2} flexWrap='wrap'>
                {selectedCategories.map((categoryId) => {
                  const category = categories.find(
                    (cat) => cat.id === categoryId,
                  );
                  return (
                    <Badge
                      key={categoryId}
                      colorScheme='brand'
                      py={1}
                      px={2}
                      borderRadius='full'
                    >
                      {category?.name} ✕
                    </Badge>
                  );
                })}
                {selectedBrands.map((brandId) => {
                  const brand = brands.find((b) => b.id === brandId);
                  return (
                    <Badge
                      key={brandId}
                      colorScheme='gray'
                      py={1}
                      px={2}
                      borderRadius='full'
                    >
                      {brand?.name} ✕
                    </Badge>
                  );
                })}
              </HStack>
            </Box>
          )}

          {/* Product Grid */}
          <ProductGrid
            products={products}
            isLoading={isLoading}
            error={error}
          />
        </GridItem>
      </Grid>

      {/* Mobile Filter Drawer */}
      <Drawer isOpen={isOpen} placement='left' onClose={onClose}>
        <DrawerOverlay />
        <DrawerContent>
          <DrawerCloseButton />
          <DrawerHeader borderBottomWidth='1px'>Lọc sản phẩm</DrawerHeader>

          <DrawerBody>
            <FilterSection />
          </DrawerBody>

          <DrawerFooter borderTopWidth='1px'>
            <Button variant='outline' mr={3} onClick={onClose}>
              Hủy
            </Button>
            <Button colorScheme='brand' onClick={onClose}>
              Áp dụng
            </Button>
          </DrawerFooter>
        </DrawerContent>
      </Drawer>
    </Container>
  );
};

export default ProductListing;
