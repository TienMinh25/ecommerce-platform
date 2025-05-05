import React, { useState, useEffect, useCallback } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  Box,
  Container,
  Heading,
  Text,
  Grid,
  GridItem,
  useBreakpointValue,
} from '@chakra-ui/react';
import productService from '../services/productService';
import categoryService from '../services/categoryService';

// Import components
import ProductFilterSidebar from '../components/products/ProductFilterSidebar';
import ProductContent from '../components/products/ProductContent';

const ProductListing = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const queryParams = new URLSearchParams(location.search);
  const isMobile = useBreakpointValue({ base: true, md: false });

  // Get state from navigation (if any)
  const categoryFromHome = location.state?.selectedCategory;

  // Parse query params from URL
  const keyword = queryParams.get('keyword') || '';
  const categoryIdsParam = queryParams.get('category_ids') || '';
  const pageParam = parseInt(queryParams.get('page')) || 1;
  const ratingParam = parseInt(queryParams.get('rating')) || null;

  // State
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [isLoadingProducts, setIsLoadingProducts] = useState(false);
  const [isLoadingCategories, setIsLoadingCategories] = useState(false);
  const [error, setError] = useState(null);
  const [metadata, setMetadata] = useState({
    total_items: 0,
    total_pages: 0,
    page: 1,
    limit: 20,
    has_next: false,
    has_previous: false
  });

  // State for filters
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [minRating, setMinRating] = useState(null);
  const [hasInitializedFromURL, setHasInitializedFromURL] = useState(false);
  const [hasLoadedCategories, setHasLoadedCategories] = useState(false);

  // Parse URL params on first load
  useEffect(() => {
    if (!hasInitializedFromURL) {
      // Parse category IDs from URL
      if (categoryIdsParam) {
        const categoryIds = categoryIdsParam
            .split(',')
            .map(id => parseInt(id))
            .filter(id => !isNaN(id));

        if (categoryIds.length > 0) {
          setSelectedCategories(categoryIds);
        }
      }

      // Parse rating filter
      if (ratingParam) {
        setMinRating(ratingParam);
      }

      setHasInitializedFromURL(true);
    }
  }, [categoryIdsParam, ratingParam, hasInitializedFromURL]);

  // Handle category from Home - only run once when receiving state
  useEffect(() => {
    if (categoryFromHome && categoryFromHome.id && !hasLoadedCategories) {
      setSelectedCategories([categoryFromHome.id]);

      // Load subcategories of the selected category
      const fetchSubCategories = async () => {
        try {
          setHasLoadedCategories(true);
          setIsLoadingCategories(true);

          const response = await categoryService.getSubCategories(categoryFromHome.id);
          if (response?.data?.data) {
            setCategories(response.data.data);
          }
        } catch (err) {
          console.error('Error loading subcategories:', err);
        } finally {
          setIsLoadingCategories(false);
        }
      };

      fetchSubCategories();
    }
  }, [categoryFromHome, hasLoadedCategories]);

  // Fetch categories based on search keyword or initial load
  useEffect(() => {
    // Only fetch categories for keyword search
    if (keyword && hasInitializedFromURL) {
      const fetchCategoriesByKeyword = async () => {
        setIsLoadingCategories(true);
        try {
          const response = await categoryService.getCategoriesByKeyword(keyword);
          if (response?.data?.data) {
            setCategories(response.data.data);
            setHasLoadedCategories(true);
          }
        } catch (err) {
          console.error('Error loading categories by keyword:', err);
        } finally {
          setIsLoadingCategories(false);
        }
      };

      fetchCategoriesByKeyword();
    }
    // Skip category fetch if we already have categories from Home
  }, [keyword, hasInitializedFromURL, categoryFromHome]);

  // Fetch products after state is initialized
  useEffect(() => {
    // Only fetch after URL params are parsed
    if (!hasInitializedFromURL) {
      return;
    }

    const fetchProducts = async () => {
      setIsLoadingProducts(true);
      setError(null);

      try {
        // Prepare parameters for API
        const options = {
          page: pageParam,
          limit: 20
        };

        // Add min_rating if available
        if (minRating) {
          options.minRating = minRating;
        }

        // Add category_ids if categories are selected
        if (selectedCategories.length > 0) {
          options.categoryIds = selectedCategories;
        }

        // Add keyword if searching
        if (keyword) {
          options.keyword = keyword;
        }

        // Call products API
        const productsResponse = await productService.getProductsByCriteria(options);

        // Process results
        if (productsResponse?.data?.data) {
          setProducts(productsResponse.data.data);

          // Update pagination metadata
          if (productsResponse.data.metadata?.pagination) {
            setMetadata(productsResponse.data.metadata.pagination);
          }
        }
      } catch (err) {
        console.error('Error loading products:', err);
        setError('Có lỗi xảy ra khi tải dữ liệu sản phẩm. Vui lòng thử lại sau.');
      } finally {
        setIsLoadingProducts(false);
      }
    };

    fetchProducts();
  }, [pageParam, selectedCategories, keyword, minRating, hasInitializedFromURL]);

  // Handle category toggle
  const handleCategoryToggle = useCallback((categoryId) => {
    // Check if category is already selected
    const isSelected = selectedCategories.includes(categoryId);
    let newSelectedCategories;

    if (isSelected) {
      // Deselect category
      newSelectedCategories = selectedCategories.filter(id => id !== categoryId);
    } else {
      // Add category to selected list
      newSelectedCategories = [...selectedCategories, categoryId];
    }

    // Update state
    setSelectedCategories(newSelectedCategories);

    // Update URL for bookmarking
    const currentParams = new URLSearchParams(location.search);
    if (newSelectedCategories.length > 0) {
      currentParams.set('category_ids', newSelectedCategories.join(','));
    } else {
      currentParams.delete('category_ids');
    }
    currentParams.delete('page'); // Reset to page 1
    navigate({ search: currentParams.toString() }, { replace: true });
  }, [selectedCategories, location.search, navigate]);

  // Handle page change
  const handlePageChange = useCallback((page) => {
    // Update URL
    const currentParams = new URLSearchParams(location.search);
    currentParams.set('page', page);
    navigate({ search: currentParams.toString() });
  }, [location.search, navigate]);

  // Handle rating change
  const handleRatingChange = useCallback((rating) => {
    // Toggle rating filter
    const newRating = minRating === rating ? null : rating;

    // Update state
    setMinRating(newRating);

    // Update URL for bookmarking
    const currentParams = new URLSearchParams(location.search);
    if (newRating) {
      currentParams.set('rating', newRating);
    } else {
      currentParams.delete('rating');
    }
    currentParams.delete('page'); // Reset to page 1
    navigate({ search: currentParams.toString() }, { replace: true });
  }, [minRating, location.search, navigate]);

  // Create title based on current filters
  const getPageTitle = () => {
    if (keyword) {
      return `Kết quả tìm kiếm: "${keyword}"`;
    } else if (categoryFromHome) {
      return categoryFromHome.name;
    } else if (selectedCategories.length === 1) {
      // Find the name of the selected category
      const selectedCategory = categories.find(c => c.id === selectedCategories[0]);
      if (selectedCategory) {
        return selectedCategory.name;
      }
      return "Sản phẩm theo danh mục";
    } else if (selectedCategories.length > 1) {
      return "Sản phẩm theo danh mục";
    } else {
      return "Tất cả sản phẩm";
    }
  };

  return (
      <Container maxW="container.xl" py={8}>
        <Box mb={8}>
          <Heading as="h1" size="xl" mb={2}>
            {getPageTitle()}
          </Heading>
          {!isLoadingProducts && (
              <Text color="gray.600">
                {metadata.total_items} sản phẩm
              </Text>
          )}
        </Box>

        {/* Main Content Grid */}
        <Grid templateColumns={{ base: '1fr', md: '250px 1fr' }} gap={8}>
          {/* Filter Sidebar */}
          {!isMobile && (
              <GridItem>
                <ProductFilterSidebar
                    categories={categories}
                    selectedCategories={selectedCategories}
                    minRating={minRating}
                    isLoadingCategories={isLoadingCategories}
                    handleCategoryToggle={handleCategoryToggle}
                    handleRatingChange={handleRatingChange}
                />
              </GridItem>
          )}

          {/* Products */}
          <GridItem>
            <ProductContent
                products={products}
                isLoading={isLoadingProducts}
                error={error}
                metadata={metadata}
                categories={categories}
                selectedCategories={selectedCategories}
                minRating={minRating}
                handleCategoryToggle={handleCategoryToggle}
                handleRatingChange={handleRatingChange}
                handlePageChange={handlePageChange}
            />
          </GridItem>
        </Grid>
      </Container>
  );
};

export default ProductListing;