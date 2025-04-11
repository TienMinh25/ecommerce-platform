import React, {useState, useEffect} from 'react';
import {
    Button,
    FormControl,
    FormLabel,
    HStack,
    Menu,
    MenuButton,
    MenuList,
    Select,
    useColorModeValue,
    VStack,
    Flex,
    Divider,
} from '@chakra-ui/react';
import {FiFilter, FiCheck} from 'react-icons/fi';

/**
 * Role filter dropdown component for managing sorting options
 *
 * @param {Object} filters - Current filter state
 * @param {Function} onFiltersChange - Callback for filter changes
 * @param {Function} onApplyFilters - Callback when filters are applied
 */
const RoleFilterDropdown = ({ filters, onFiltersChange, onApplyFilters }) => {
    // Local state for the component
    const [localFilters, setLocalFilters] = useState({
        sortBy: filters.sortBy || 'name',
        sortOrder: filters.sortOrder || 'asc'
    });

    // Update local state when props change
    useEffect(() => {
        setLocalFilters({
            sortBy: filters.sortBy || 'name',
            sortOrder: filters.sortOrder || 'asc'
        });
    }, [filters]);

    // UI colors
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');

    // Handle filter changes
    const handleFilterChange = (e) => {
        const { name, value } = e.target;
        setLocalFilters((prev) => ({ ...prev, [name]: value }));
    };

    // Apply filters
    const handleApply = () => {
        onFiltersChange(localFilters);
        onApplyFilters(localFilters);
    };

    // Reset filters
    const handleReset = () => {
        const resetFilters = {
            sortBy: 'name',
            sortOrder: 'asc'
        };
        setLocalFilters(resetFilters);
        onFiltersChange(resetFilters);
        onApplyFilters(resetFilters);
    };

    return (
        <Menu closeOnSelect={false}>
            <MenuButton
                as={Button}
                leftIcon={<FiFilter />}
                size="sm"
                variant="outline"
                colorScheme="gray"
            >
                Filter
            </MenuButton>
            <MenuList
                bg={bgColor}
                borderColor={borderColor}
                p={4}
                minW="300px"
                shadow="lg"
                borderRadius="md"
            >
                <VStack spacing={4} align="stretch">
                    <FormControl>
                        <FormLabel fontSize="sm" fontWeight="medium">Sort By</FormLabel>
                        <Select
                            name="sortBy"
                            value={localFilters.sortBy}
                            onChange={handleFilterChange}
                            size="sm"
                            borderRadius="md"
                        >
                            <option value="name">Name</option>
                        </Select>
                    </FormControl>

                    <FormControl>
                        <FormLabel fontSize="sm" fontWeight="medium">Sort Order</FormLabel>
                        <Select
                            name="sortOrder"
                            value={localFilters.sortOrder}
                            onChange={handleFilterChange}
                            size="sm"
                            borderRadius="md"
                        >
                            <option value="asc">Ascending</option>
                            <option value="desc">Descending</option>
                        </Select>
                    </FormControl>

                    <Divider my={1} />

                    <Flex justify="space-between">
                        <Button
                            size="sm"
                            variant="ghost"
                            onClick={handleReset}
                        >
                            Reset
                        </Button>
                        <Button
                            size="sm"
                            colorScheme="blue"
                            onClick={handleApply}
                            leftIcon={<FiCheck size={14} />}
                        >
                            Apply
                        </Button>
                    </Flex>
                </VStack>
            </MenuList>
        </Menu>
    );
};

export default RoleFilterDropdown;