import React, { useState } from 'react';
import {
    Box,
    Text,
    Flex,
    Checkbox,
    Stack,
    Spinner,
    Center,
    Heading,
    Divider,
    Button,
} from '@chakra-ui/react';
import { ChevronDownIcon, ChevronUpIcon } from '@chakra-ui/icons';
import { FaFilter } from 'react-icons/fa';

const ProductFilterSidebar = ({
                                  categories,
                                  selectedCategories,
                                  minRating,
                                  isLoadingCategories,
                                  handleCategoryToggle,
                                  handleRatingChange,
                              }) => {
    const [showAllCategories, setShowAllCategories] = useState(false);

    // Default: only show first 6 categories
    const visibleCategories = showAllCategories
        ? categories
        : categories.slice(0, 6);

    const hasMoreCategories = categories.length > 6;

    return (
        <Box bg="white" p={4} borderRadius="md" boxShadow="sm">
            {/* Title - "Search filters" with filter icon */}
            <Heading as="h3" size="md" mb={4} display="flex" alignItems="center">
                <FaFilter style={{ marginRight: '8px' }} />
                Bộ lọc tìm kiếm
            </Heading>

            {/* Categories Section */}
            <Box mb={4}>
                <Text fontWeight="bold" mb={2}>
                    Danh mục
                </Text>

                {/* Categories Content */}
                <Stack spacing={2} mb={2}>
                    {isLoadingCategories ? (
                        <Center py={4}>
                            <Spinner size="sm" />
                        </Center>
                    ) : categories.length > 0 ? (
                        visibleCategories.map((category) => (
                            <Checkbox
                                key={category.id}
                                isChecked={selectedCategories.includes(category.id)}
                                onChange={() => handleCategoryToggle(category.id)}
                                colorScheme="brand"
                            >
                                <Flex justify="space-between" width="100%">
                                    <Text>{category.name}</Text>
                                    {category.product_count !== undefined && (
                                        <Text fontSize="sm" color="gray.500" marginLeft={2}>
                                            ({category.product_count})
                                        </Text>
                                    )}
                                </Flex>
                            </Checkbox>
                        ))
                    ) : (
                        <Text color="gray.500">No categories found</Text>
                    )}
                </Stack>

                {/* "See more" button - only shown if there are more than 6 categories */}
                {hasMoreCategories && (
                    <Button
                        onClick={() => setShowAllCategories(!showAllCategories)}
                        variant="ghost"
                        size="sm"
                        leftIcon={showAllCategories ? <ChevronUpIcon /> : <ChevronDownIcon />}
                        color="blue.500"
                        p={0}
                        mt={1}
                        _hover={{ bg: 'transparent', textDecoration: 'underline' }}
                    >
                        {showAllCategories ? 'Thu gọn' : 'Xem thêm'}
                    </Button>
                )}
            </Box>

            {/* Divider between sections */}
            <Divider my={4} borderColor="gray.300" />

            {/* Filter by Rating - Always shown */}
            <Box>
                <Text fontWeight="bold" mb={3}>
                    Đánh giá
                </Text>
                <Stack spacing={2}>
                    {[5, 4, 3, 2, 1].map((rating) => (
                        <Checkbox
                            key={rating}
                            isChecked={minRating === rating}
                            onChange={() => handleRatingChange(rating)}
                            colorScheme="brand"
                        >
                            {rating} sao trở lên
                        </Checkbox>
                    ))}
                </Stack>
            </Box>
        </Box>
    );
};

export default ProductFilterSidebar;