import React, { useEffect, useRef, useState } from 'react';
import {
    Box,
    Button,
    Collapse,
    Flex,
    FormControl,
    FormLabel,
    IconButton,
    Input,
    InputGroup,
    InputRightElement,
    Menu,
    MenuButton,
    MenuItem,
    MenuList,
    Text,
    useColorModeValue,
    useDisclosure,
    VStack,
    HStack,
    Switch,
} from '@chakra-ui/react';
import { FiCalendar, FiChevronDown, FiFilter, FiX } from 'react-icons/fi';
import {ModernToggleGroup} from "../../../toggle/main.jsx";

// Improved Dropdown Component
const CustomDropdown = ({ options, value, onChange, name, placeholder }) => {
    const selectedOption = options.find(option => option.value === value);
    const menuButtonRef = useRef(null);
    const [menuWidth, setMenuWidth] = useState(0);

    useEffect(() => {
        const updateWidth = () => {
            if (menuButtonRef.current) {
                setMenuWidth(menuButtonRef.current.offsetWidth);
            }
        };

        updateWidth();
        window.addEventListener('resize', updateWidth);

        return () => {
            window.removeEventListener('resize', updateWidth);
        };
    }, []);

    return (
        <Box position="relative" width="100%" zIndex="dropdown">
            <Menu
                isLazy
                gutter={0}
                strategy="fixed"
                autoSelect={false}
                closeOnBlur={true}
                closeOnSelect={true}
                matchWidth={true}
            >
                <MenuButton
                    as={Button}
                    ref={menuButtonRef}
                    rightIcon={<FiChevronDown />}
                    width="100%"
                    justifyContent="space-between"
                    textAlign="left"
                    variant="outline"
                    color={value ? "black" : "gray.500"}
                    fontWeight="normal"
                    height="36px"
                    borderRadius="md"
                    _focus={{ boxShadow: "outline" }}
                >
                    {value ? selectedOption?.label : placeholder}
                </MenuButton>
                <MenuList
                    minW={`${menuWidth}px`}
                    width="100%"
                    maxHeight="200px"
                    overflowY="auto"
                    overflowX="hidden"
                    zIndex={2000}
                    borderRadius="md"
                    boxShadow="lg"
                    border="1px solid"
                    borderColor="gray.200"
                    py={1}
                >
                    {options.map((option) => (
                        <MenuItem
                            key={option.value}
                            onClick={() => onChange({ target: { name, value: option.value } })}
                            bg={value === option.value ? "blue.50" : "transparent"}
                            color={value === option.value ? "blue.600" : "inherit"}
                            _hover={{ bg: "gray.100" }}
                            px={3}
                            py={2}
                            width="100%"
                        >
                            {option.label}
                        </MenuItem>
                    ))}
                </MenuList>
            </Menu>
        </Box>
    );
};

const CouponFilterDropdown = ({ filters, onFiltersChange, onApplyFilters }) => {
    // Dropdown visibility and outside click handling
    const { isOpen, onToggle, onClose } = useDisclosure();
    const filterRef = useRef(null);
    const dropdownRef = useRef(null);

    // Styling variables
    const filterBgColor = useColorModeValue('white', 'gray.800');
    const filterBorderColor = useColorModeValue('gray.200', 'gray.600');
    const headerBgColor = useColorModeValue('gray.50', 'gray.700');
    const buttonHoverColor = useColorModeValue('blue.50', 'blue.700');

    // Input change handlers
    const handleChange = (e) => {
        const { name, value } = e.target;
        onFiltersChange({
            ...filters,
            [name]: value
        });
    };

    const handleToggleChange = (name, value) => {
        onFiltersChange({
            ...filters,
            [name]: value
        });
    };

    // Apply and reset filters
    const handleApplyFilters = () => {
        // Convert date strings to UTC format if they exist
        const formattedFilters = { ...filters };

        if (formattedFilters.start_date) {
            const startDate = new Date(formattedFilters.start_date);
            formattedFilters.start_date = startDate.toISOString();
        }

        if (formattedFilters.end_date) {
            const endDate = new Date(formattedFilters.end_date);
            endDate.setHours(23, 59, 59, 999);
            formattedFilters.end_date = endDate.toISOString();
        }

        // Validate dates - don't allow future dates for end_date
        const now = new Date();
        if (formattedFilters.end_date && new Date(formattedFilters.end_date) > now) {
            console.warn("End date is in the future, using current date instead");
            formattedFilters.end_date = now.toISOString();
        }

        const activeFilters = Object.fromEntries(
            Object.entries(formattedFilters).filter(([_, v]) => v !== '')
        );

        onApplyFilters(activeFilters);
        onClose();
    };

    const handleResetFilters = () => {
        onFiltersChange({
            discount_type: '',
            start_date: '',
            end_date: '',
            is_active: '',
        });
    };

    const discountTypeOptions = [
        { label: 'Phần trăm', value: 'percentage' },
        { label: 'Số tiền cố định', value: 'fixed_amount' }
    ];

    const isActiveOptions = [
        { label: 'Đang hoạt động', value: 'true' },
        { label: 'Không hoạt động', value: 'false' }
    ];

    // Click outside handler to close the dropdown
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (
                isOpen &&
                filterRef.current &&
                !filterRef.current.contains(event.target) &&
                dropdownRef.current &&
                !dropdownRef.current.contains(event.target)
            ) {
                onClose();
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [isOpen, onClose]);

    return (
        <Box position="relative" zIndex={1000} ref={filterRef}>
            {/* Filter Button */}
            <Button
                leftIcon={<FiFilter />}
                onClick={onToggle}
                variant={isOpen ? "solid" : "outline"}
                colorScheme={isOpen ? "blue" : "gray"}
                size="sm"
                borderRadius="md"
                px={4}
            >
                Bộ lọc
            </Button>

            {/* Filter Dropdown Panel */}
            <Collapse in={isOpen} animateOpacity>
                <Box
                    position="absolute"
                    top="40px"
                    right="0"
                    width="450px"
                    maxW="calc(100vw - 20px)"
                    bg={filterBgColor}
                    boxShadow="0 4px 20px rgba(0,0,0,0.15)"
                    borderRadius="lg"
                    borderWidth="1px"
                    borderColor={filterBorderColor}
                    overflow="hidden"
                    zIndex={1000}
                    ref={dropdownRef}
                >
                    {/* Header */}
                    <Flex
                        justifyContent="space-between"
                        alignItems="center"
                        p={4}
                        borderBottomWidth="1px"
                        borderColor={filterBorderColor}
                        bg={headerBgColor}
                    >
                        <Text fontWeight="bold" fontSize="lg" color="gray.800">
                            Bộ lọc
                        </Text>
                        <IconButton
                            icon={<FiX />}
                            variant="ghost"
                            size="sm"
                            aria-label="Close filters"
                            onClick={onClose}
                            color="gray.600"
                            _hover={{ bg: "gray.200" }}
                        />
                    </Flex>

                    {/* Filter Form */}
                    <Box
                        p={4}
                        maxH="60vh"
                        overflowY="auto"
                        overflowX="hidden"
                        sx={{
                            '&::-webkit-scrollbar': {
                                width: '6px',
                                height: '0px',
                            },
                            '&::-webkit-scrollbar-track': {
                                width: '6px',
                                background: 'transparent',
                            },
                            '&::-webkit-scrollbar-thumb': {
                                background: '#cbd5e0',
                                borderRadius: '24px',
                            },
                        }}
                    >
                        <VStack spacing={5} align="stretch" width={'100%'}>
                            {/* Discount Type Selection */}
                            <FormControl width={'100%'}>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Loại khuyến mãi
                                </FormLabel>
                                <CustomDropdown
                                    options={discountTypeOptions}
                                    value={filters.discount_type}
                                    onChange={handleChange}
                                    name="discount_type"
                                    placeholder="Chọn loại khuyến mãi"
                                />
                            </FormControl>

                            {/* Active Status */}
                            <FormControl>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Trạng thái hoạt động
                                </FormLabel>
                                <ModernToggleGroup
                                    options={isActiveOptions}
                                    value={filters.is_active}
                                    onChange={handleToggleChange}
                                    name="is_active"
                                />
                            </FormControl>

                            {/* Date Range */}
                            <FormControl>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Lọc theo thời gian khuyến mãi
                                </FormLabel>
                                <VStack spacing={3} align="stretch">
                                    <Box>
                                        <Text fontSize="xs" mb={1} color="gray.600">
                                            Ngày bắt đầu
                                        </Text>
                                        <InputGroup size="sm">
                                            <Input
                                                type="date"
                                                name="start_date"
                                                value={filters.start_date}
                                                onChange={handleChange}
                                                height="36px"
                                                borderRadius="md"
                                            />
                                            <InputRightElement
                                                pointerEvents="none"
                                                height="36px"
                                                children={<FiCalendar color="gray.400" />}
                                            />
                                        </InputGroup>
                                    </Box>
                                    <Box>
                                        <Text fontSize="xs" mb={1} color="gray.600">
                                            Ngày kết thúc
                                        </Text>
                                        <InputGroup size="sm">
                                            <Input
                                                type="date"
                                                name="end_date"
                                                value={filters.end_date}
                                                onChange={handleChange}
                                                height="36px"
                                                borderRadius="md"
                                                min={filters.start_date}
                                                disabled={!filters.start_date}
                                            />
                                            <InputRightElement
                                                pointerEvents="none"
                                                height="36px"
                                                children={<FiCalendar color="gray.400" />}
                                            />
                                        </InputGroup>
                                    </Box>
                                </VStack>
                            </FormControl>
                        </VStack>
                    </Box>

                    {/* Footer Buttons */}
                    <Flex
                        justify="space-between"
                        p={4}
                        borderTopWidth="1px"
                        borderColor={filterBorderColor}
                        bg={headerBgColor}
                    >
                        <Button
                            variant="outline"
                            colorScheme="blue"
                            onClick={handleResetFilters}
                            size="sm"
                            width="45%"
                            borderRadius="md"
                            _hover={{ bg: buttonHoverColor }}
                            height="40px"
                        >
                            Reset bộ lọc
                        </Button>
                        <Button
                            colorScheme="blue"
                            onClick={handleApplyFilters}
                            size="sm"
                            width="45%"
                            borderRadius="md"
                            _hover={{ bg: "blue.600" }}
                            height="40px"
                        >
                            Áp dụng bộ lọc
                        </Button>
                    </Flex>
                </Box>
            </Collapse>
        </Box>
    );
};

export default CouponFilterDropdown;