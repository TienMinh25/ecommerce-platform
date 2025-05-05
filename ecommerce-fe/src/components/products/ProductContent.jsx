import React from 'react';
import {
    Box,
    HStack,
    Badge,
} from '@chakra-ui/react';
import ProductGrid from './ProductGrid';
import Pagination from '../common/Pagination';

const ProductContent = ({
                            products,
                            isLoading,
                            error,
                            metadata,
                            categories,
                            selectedCategories,
                            minRating,
                            handleCategoryToggle,
                            handleRatingChange,
                            handlePageChange,
                        }) => {
    return (
        <Box>
            {/* Hiển thị filter đã chọn */}
            {(selectedCategories.length > 0 || minRating) && (
                <HStack spacing={2} mb={4} flexWrap="wrap">
                    {selectedCategories.map((categoryId) => {
                        const category = categories.find((c) => c.id === categoryId);
                        if (!category) return null;

                        return (
                            <Badge
                                key={categoryId}
                                py={1}
                                px={2}
                                borderRadius="full"
                                colorScheme="brand"
                                cursor="pointer"
                                onClick={() => handleCategoryToggle(categoryId)}
                            >
                                {category.name} ✕
                            </Badge>
                        );
                    })}

                    {minRating && (
                        <Badge
                            py={1}
                            px={2}
                            borderRadius="full"
                            colorScheme="yellow"
                            cursor="pointer"
                            onClick={() => handleRatingChange(minRating)}
                        >
                            {minRating} sao trở lên ✕
                        </Badge>
                    )}
                </HStack>
            )}

            {/* Lưới sản phẩm */}
            <ProductGrid
                products={products}
                isLoading={isLoading}
                error={error}
            />

            {/* Phân trang */}
            {!isLoading && metadata.total_pages > 1 && (
                <Box mt={8} display="flex" justifyContent="center">
                    <Pagination
                        currentPage={metadata.page}
                        totalPages={metadata.total_pages}
                        onPageChange={handlePageChange}
                    />
                </Box>
            )}
        </Box>
    );
};

export default ProductContent;