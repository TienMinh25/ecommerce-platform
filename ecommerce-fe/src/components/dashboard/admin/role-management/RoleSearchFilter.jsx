import React, {useState, useEffect} from 'react';
import {
    Flex,
    HStack,
    IconButton,
    Input,
    InputGroup,
    InputLeftElement,
    Tooltip,
    useColorModeValue,
} from '@chakra-ui/react';
import {FiSearch, FiX, FiRefreshCw} from 'react-icons/fi';

/**
 * Search component for roles management
 *
 * @param {Object} filters - Current filter state
 * @param {Function} onFiltersChange - Callback for filter changes
 * @param {Function} onApplyFilters - Callback when filters are applied
 * @param {Function} onRefresh - Callback to refresh data
 * @param {boolean} isLoading - Whether data is currently loading
 */
const RoleSearchFilter = ({
                              filters,
                              onFiltersChange,
                              onApplyFilters,
                              onRefresh,
                              isLoading = false
                          }) => {
    // Local state for search value
    const [searchQuery, setSearchQuery] = useState(filters.searchValue || '');
    const [typingTimeout, setTypingTimeout] = useState(0);

    // UI colors
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');

    // Update local state when props change
    useEffect(() => {
        setSearchQuery(filters.searchValue || '');
    }, [filters.searchValue]);

    // Handle search input changes with debounce
    const handleSearchInputChange = (e) => {
        const inputText = e.target.value;
        setSearchQuery(inputText);

        clearTimeout(typingTimeout);

        setTypingTimeout(
            setTimeout(() => {
                const newFilters = {
                    ...filters,
                    searchValue: inputText.trim(),
                    searchBy: 'name' // Always search by name
                };

                onFiltersChange(newFilters);

                // Only apply filters immediately if clearing search
                if (!inputText.trim()) {
                    onApplyFilters(newFilters);
                }
            }, 300)
        );
    };

    // Clear search
    const clearSearch = () => {
        setSearchQuery('');

        const newFilters = {
            ...filters,
            searchValue: ''
        };

        onFiltersChange(newFilters);
        onApplyFilters(newFilters);
    };

    // Apply search
    const handleSearchButtonClick = () => {
        if (!searchQuery.trim()) return;

        const newFilters = {
            ...filters,
            searchValue: searchQuery.trim(),
            searchBy: 'name' // Always search by name
        };

        onFiltersChange(newFilters);
        onApplyFilters(newFilters);
    };

    const hasActiveSearch = searchQuery.trim() !== '';

    return (
        <Flex
            borderWidth="1px"
            borderRadius="lg"
            overflow="hidden"
            align="center"
            bg={bgColor}
            shadow="sm"
            flex="1"
            maxW={{ base: 'full', lg: '450px' }}
        >
            <InputGroup size="md" variant="unstyled">
                <InputLeftElement pointerEvents="none" h="full" pl={3}>
                    <FiSearch color={useColorModeValue('gray.400', 'gray.500')} />
                </InputLeftElement>
                <Input
                    placeholder="Tìm kiếm theo tên..."
                    pl={10}
                    pr={2}
                    py={2.5}
                    value={searchQuery}
                    onChange={handleSearchInputChange}
                    _placeholder={{ color: 'gray.400' }}
                    onKeyPress={(e) => {
                        if (e.key === 'Enter') {
                            handleSearchButtonClick();
                        }
                    }}
                />
            </InputGroup>

            <HStack spacing={1} pr={2}>
                {hasActiveSearch && (
                    <Tooltip label="Clear search" hasArrow>
                        <IconButton
                            icon={<FiX size={16} />}
                            onClick={clearSearch}
                            aria-label="Clear search"
                            variant="ghost"
                            colorScheme="red"
                            size="sm"
                        />
                    </Tooltip>
                )}
                <Tooltip label="Refresh data" hasArrow>
                    <IconButton
                        icon={<FiRefreshCw size={16} />}
                        onClick={onRefresh}
                        aria-label="Refresh data"
                        variant="ghost"
                        colorScheme="blue"
                        size="sm"
                        isLoading={isLoading}
                    />
                </Tooltip>
            </HStack>
        </Flex>
    );
};

export default RoleSearchFilter;