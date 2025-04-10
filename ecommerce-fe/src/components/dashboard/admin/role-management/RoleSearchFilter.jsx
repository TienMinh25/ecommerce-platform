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
    Button,
    Select,
} from '@chakra-ui/react';
import {FiSearch, FiX, FiFilter} from 'react-icons/fi';
import RoleFilterDropdown from './RoleFilterDropdown.jsx';

const RoleSearchFilter = ({ filters, onFiltersChange, onApplyFilters, showOnlyFilter = false }) => {
    const [searchQuery, setSearchQuery] = useState(filters.searchValue || '');
    const [searchField, setSearchField] = useState(filters.searchBy || 'name');
    const [typingTimeout, setTypingTimeout] = useState(0);

    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');

    // Nếu chỉ hiển thị filter
    if (showOnlyFilter) {
        return (
            <RoleFilterDropdown filters={filters} onFiltersChange={onFiltersChange} onApplyFilters={onApplyFilters} />
        );
    }

    const handleSearchInputChange = (e) => {
        const inputText = e.target.value;
        setSearchQuery(inputText);

        clearTimeout(typingTimeout);

        setTypingTimeout(
            setTimeout(() => {
                const newFilters = { ...filters, searchValue: inputText.trim(), searchBy: searchField };
                if (!inputText.trim()) {
                    onFiltersChange(newFilters);
                } else if (searchField && inputText.trim()) {
                    onFiltersChange(newFilters);
                }
            }, 300)
        );
    };

    const handleSearchFieldChange = (e) => {
        const newField = e.target.value;
        setSearchField(newField);
        const newFilters = { ...filters, searchBy: newField, searchValue: searchQuery };
        onFiltersChange(newFilters);
    };

    const clearSearch = () => {
        setSearchQuery('');
        const newFilters = { ...filters, searchValue: '' };
        onFiltersChange(newFilters);
    };

    const handleSearchButtonClick = () => {
        const newFilters = { ...filters, searchValue: searchQuery.trim(), searchBy: searchField };
        if (searchField && searchQuery.trim()) {
            onFiltersChange(newFilters);
        }
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
            <Select
                value={searchField}
                onChange={handleSearchFieldChange}
                variant="unstyled"
                size="md"
                w="120px"
                pl={3}
                pr={0}
                py={2.5}
                borderRight="1px"
                borderColor={borderColor}
                borderRadius="0"
                fontSize="sm"
            >
                <option value="name">Name</option>
                <option value="description">Description</option> {/* Có thể thêm field khác nếu API hỗ trợ */}
            </Select>

            <InputGroup size="md" variant="unstyled">
                <InputLeftElement pointerEvents="none" h="full" pl={3}>
                    <FiSearch color="gray.400" />
                </InputLeftElement>
                <Input
                    placeholder={`Search by ${searchField.toLowerCase()}...`}
                    pl={10}
                    pr={2}
                    py={2.5}
                    value={searchQuery}
                    onChange={handleSearchInputChange}
                    _placeholder={{ color: 'gray.400' }}
                />
            </InputGroup>

            <HStack>
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
                <Tooltip label="Search" hasArrow>
                    <IconButton
                        icon={<FiSearch size={16} />}
                        onClick={handleSearchButtonClick}
                        aria-label="Search"
                        variant="ghost"
                        colorScheme="blue"
                        size="sm"
                        mr={2}
                        isDisabled={!searchField || !searchQuery}
                    />
                </Tooltip>
            </HStack>
        </Flex>
    );
};

export default RoleSearchFilter;