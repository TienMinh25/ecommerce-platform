import React from 'react';
import { HStack, Button, IconButton } from '@chakra-ui/react';
import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';

const Pagination = ({ currentPage, totalPages, onPageChange }) => {
    // Calculate visible page numbers based on Shopee-style pagination
    const getVisiblePages = () => {
        const delta = 2; // Number of page buttons to show on either side of current page
        const pages = [];

        // Always include page 1
        pages.push(1);

        // Calculate range around current page
        const leftBound = Math.max(2, currentPage - delta);
        const rightBound = Math.min(totalPages - 1, currentPage + delta);

        // Add ellipsis after page 1 if needed
        if (leftBound > 2) {
            pages.push('...');
        }

        // Add pages around current page
        for (let i = leftBound; i <= rightBound; i++) {
            pages.push(i);
        }

        // Add ellipsis before last page if needed
        if (rightBound < totalPages - 1) {
            pages.push('...');
        }

        // Always include the last page if there's more than one page
        if (totalPages > 1) {
            pages.push(totalPages);
        }

        return pages;
    };

    const visiblePages = getVisiblePages();

    // Safe page change handler to prevent page navigation issues
    const handlePageChange = (page) => {
        if (page < 1 || page > totalPages) return;
        if (typeof onPageChange === 'function') {
            onPageChange(page);
        }
    };

    return (
        <HStack spacing={2} justify="center" my={6}>
            {/* Previous page button */}
            <IconButton
                icon={<ChevronLeftIcon />}
                aria-label="Previous page"
                onClick={() => handlePageChange(currentPage - 1)}
                isDisabled={currentPage === 1}
                variant="outline"
                borderRadius="md"
                colorScheme="brand"
            />

            {/* Page number buttons */}
            {visiblePages.map((page, index) => {
                if (page === '...') {
                    return (
                        <Button
                            key={`ellipsis-${index}`}
                            variant="ghost"
                            size="md"
                            pointerEvents="none"
                        >
                            ...
                        </Button>
                    );
                }

                return (
                    <Button
                        key={`page-${page}`}
                        onClick={() => handlePageChange(page)}
                        colorScheme={page === currentPage ? 'brand' : 'gray'}
                        variant={page === currentPage ? 'solid' : 'outline'}
                        size="md"
                        borderRadius="md"
                    >
                        {page}
                    </Button>
                );
            })}

            {/* Next page button */}
            <IconButton
                icon={<ChevronRightIcon />}
                aria-label="Next page"
                onClick={() => handlePageChange(currentPage + 1)}
                isDisabled={currentPage === totalPages || totalPages === 0}
                variant="outline"
                borderRadius="md"
                colorScheme="brand"
            />
        </HStack>
    );
};

export default Pagination;