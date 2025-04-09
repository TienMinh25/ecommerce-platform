import React, {useState} from 'react';
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
} from '@chakra-ui/react';
import {FiFilter} from 'react-icons/fi';

const RoleFilterDropdown = ({ filters, onFiltersChange, onApplyFilters }) => {
    const [localFilters, setLocalFilters] = useState(filters);

    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');

    const handleFilterChange = (e) => {
        const { name, value } = e.target;
        setLocalFilters((prev) => ({ ...prev, [name]: value }));
    };

    const handleApply = () => {
        onFiltersChange(localFilters);
        onApplyFilters(localFilters); // Gửi filter đã áp dụng
    };

    const hasActiveFilters = Object.values(localFilters).some((value) => value !== '');

    return (
        <Menu closeOnSelect={false}>
            <MenuButton
                as={Button}
                leftIcon={<FiFilter />}
                size="sm"
                variant={hasActiveFilters ? 'solid' : 'outline'}
                colorScheme={hasActiveFilters ? 'blue' : 'gray'}
            >
                Filter
            </MenuButton>
            <MenuList bg={bgColor} borderColor={borderColor} p={4} minW="300px">
                <VStack spacing={4} align="stretch">
                    <FormControl>
                        <FormLabel fontSize="sm">Sort By</FormLabel>
                        <Select
                            name="sortBy"
                            value={localFilters.sortBy}
                            onChange={handleFilterChange}
                            size="sm"
                        >
                            <option value="">None</option>
                            <option value="name">Name</option>
                        </Select>
                    </FormControl>

                    <FormControl>
                        <FormLabel fontSize="sm">Sort Order</FormLabel>
                        <Select
                            name="sortOrder"
                            value={localFilters.sortOrder}
                            onChange={handleFilterChange}
                            size="sm"
                        >
                            <option value="asc">Ascending</option>
                            <option value="desc">Descending</option>
                        </Select>
                    </FormControl>

                    <HStack justify="flex-end">
                        <Button size="sm" colorScheme="blue" onClick={handleApply}>
                            Apply
                        </Button>
                    </HStack>
                </VStack>
            </MenuList>
        </Menu>
    );
};

export default RoleFilterDropdown;