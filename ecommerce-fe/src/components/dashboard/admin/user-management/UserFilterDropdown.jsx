import React, {useEffect, useRef, useState} from 'react';
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
import {FiCalendar, FiChevronDown, FiFilter, FiX} from 'react-icons/fi';
import roleService from "../../../../services/roleService.js";
import {ModernToggleGroup} from "../../../toggle/main.jsx";

// Improved Dropdown Component with correct positioning and better styling
const CustomDropdown = ({options, value, onChange, name, placeholder}) => {
    const selectedOption = options.find(option => option.value === value);
    const menuButtonRef = useRef(null);
    const [menuWidth, setMenuWidth] = useState(0);

    // Update the width when the component mounts or the window resizes
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
                matchWidth={true} // Force menu to match button width
            >
                <MenuButton
                    as={Button}
                    ref={menuButtonRef}
                    rightIcon={<FiChevronDown/>}
                    width="100%"
                    justifyContent="space-between"
                    textAlign="left"
                    variant="outline"
                    color={value ? "black" : "gray.500"}
                    fontWeight="normal"
                    height="36px"
                    borderRadius="md"
                    _focus={{boxShadow: "outline"}}
                >
                    {value ? selectedOption?.label : placeholder}
                </MenuButton>
                <MenuList
                    minW={`${menuWidth}px`} // Set minimum width to match button
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
                    sx={{
                        // Ensure menu list takes full width
                        '& > button, & > a': {
                            width: '100%'
                        }
                    }}
                >
                    {options.map((option) => (
                        <MenuItem
                            key={option.value}
                            onClick={() => onChange({target: {name, value: option.value}})}
                            bg={value === option.value ? "blue.50" : "transparent"}
                            color={value === option.value ? "blue.600" : "inherit"}
                            _hover={{bg: "gray.100"}}
                            px={3}
                            py={2}
                            width="100%" // Make each item full width
                        >
                            {option.label}
                        </MenuItem>
                    ))}
                </MenuList>
            </Menu>
        </Box>
    );
};

const UserFilterDropdown = ({ filters, onFiltersChange, onApplyFilters }) => {
    // State for role options
    const [roleOptions, setRoleOptions] = useState([
        {label: 'Admin', value: '1'},
        {label: 'Supplier', value: '3'},
        {label: 'Deliverer', value: '4'},
        {label: 'User', value: '2'}
    ]);
    const [isLoadingRoles, setIsLoadingRoles] = useState(false);
    const [roleError, setRoleError] = useState(null);

    // Dropdown visibility and outside click handling
    const {isOpen, onToggle, onClose} = useDisclosure();
    const filterRef = useRef(null);
    const dropdownRef = useRef(null);

    // Styling variables with more refined palette
    const filterBgColor = useColorModeValue('white', 'gray.800');
    const filterBorderColor = useColorModeValue('gray.200', 'gray.600');
    const headerBgColor = useColorModeValue('gray.50', 'gray.700');
    const buttonHoverColor = useColorModeValue('blue.50', 'blue.700');

    // Input change handlers
    const handleChange = (e) => {
        const {name, value} = e.target;
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
        const formattedFilters = {...filters};

        if (formattedFilters.updatedAtStartFrom) {
            // Convert to UTC format: YYYY-MM-DDT00:00:00.000Z
            const startDate = new Date(formattedFilters.updatedAtStartFrom);
            formattedFilters.updatedAtStartFrom = startDate.toISOString();
        }

        if (formattedFilters.updatedAtEndFrom) {
            // Convert to UTC format: YYYY-MM-DDT23:59:59.999Z
            const endDate = new Date(formattedFilters.updatedAtEndFrom);
            endDate.setHours(23, 59, 59, 999);
            formattedFilters.updatedAtEndFrom = endDate.toISOString();
        }

        // Validate dates - don't allow future dates
        const now = new Date();
        if (formattedFilters.updatedAtStartFrom && new Date(formattedFilters.updatedAtStartFrom) > now) {
            console.warn("Start date is in the future, using current date instead");
            formattedFilters.updatedAtStartFrom = now.toISOString();
        }

        if (formattedFilters.updatedAtEndFrom && new Date(formattedFilters.updatedAtEndFrom) > now) {
            console.warn("End date is in the future, using current date instead");
            formattedFilters.updatedAtEndFrom = now.toISOString();
        }

        const activeFilters = Object.fromEntries(
            Object.entries(formattedFilters).filter(([_, v]) => v !== '')
        );

        onApplyFilters(activeFilters);
        onClose();
    };

    const handleResetFilters = () => {
        onFiltersChange({
            sortBy: '',
            sortOrder: 'asc',
            emailVerify: '',
            phoneVerify: '',
            status: '',
            updatedAtStartFrom: '',
            updatedAtEndFrom: '',
            roleID: '',
        });
    };

    // Fetch roles from API using roleService
    const fetchRoles = async () => {
        setIsLoadingRoles(true);
        setRoleError(null);

        try {
            // Use the roleService to fetch roles
            const roles = await roleService.getRoles();

            // Transform API response to match our dropdown format
            if (roles && Array.isArray(roles)) {
                const formattedRoles = roles.map(role => ({
                    label: role.name,
                    value: role.id.toString()
                }));
                setRoleOptions(formattedRoles);
            } else {
                // Fallback to default roles if API response is unexpected
                setRoleOptions([
                    {label: 'Admin', value: '1'},
                    {label: 'Supplier', value: '3'},
                    {label: 'Deliverer', value: '4'},
                    {label: 'User', value: '2'}
                ]);
            }
        } catch (error) {
            console.error("Error fetching roles:", error);
            setRoleError(error.message);

            // Fallback to default roles on error
            setRoleOptions([
                {label: 'Admin', value: '1'},
                {label: 'Supplier', value: '3'},
                {label: 'Deliverer', value: '4'},
                {label: 'User', value: '2'}
            ]);
        } finally {
            setIsLoadingRoles(false);
        }
    };

    const sortByOptions = [
        {label: 'Họ và tên', value: 'fullname'},
        {label: 'Email', value: 'email'},
        {label: 'Ngày sinh nhật', value: 'birthdate'},
        {label: 'Thời gian cập nhật', value: 'updated_at'},
        {label: 'Số điện thoại', value: 'phone'}
    ];

    const statusOptions = [
        {label: 'Active', value: 'active'},
        {label: 'Inactive', value: 'inactive'}
    ];

    const sortOrderOptions = [
        {label: 'Tăng dần', value: 'asc'},
        {label: 'Giảm dần', value: 'desc'}
    ];

    // Fetch roles when the filter dropdown is opened
    useEffect(() => {
        if (isOpen) {
            fetchRoles();
        }
    }, [isOpen]);

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
                leftIcon={<FiFilter/>}
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
                            icon={<FiX/>}
                            variant="ghost"
                            size="sm"
                            aria-label="Close filters"
                            onClick={onClose}
                            color="gray.600"
                            _hover={{bg: "gray.200"}}
                        />
                    </Flex>

                    {/* Filter Form - with improved scroll */}
                    <Box
                        p={4}
                        maxH="60vh"
                        overflowY="auto"
                        overflowX="hidden"
                        sx={{
                            '&::-webkit-scrollbar': {
                                width: '6px',
                                height: '0px', // Hide horizontal scrollbar
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
                            {/* Role Selection */}
                            <FormControl width={'100%'}>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Vai trò người dùng
                                </FormLabel>
                                {isLoadingRoles ? (
                                    <Flex align="center" justify="center" height="36px" borderWidth="1px" borderRadius="md" borderColor="gray.200" px={3}>
                                        <Text fontSize="sm" color="gray.500">Đang tải vai trò...</Text>
                                    </Flex>
                                ) : roleError ? (
                                    <Flex align="center" justify="center" height="36px" borderWidth="1px" borderRadius="md" borderColor="red.200" bg="red.50" px={3}>
                                        <Text fontSize="sm" color="red.500">Thất bại khi lấy các vai trò</Text>
                                    </Flex>
                                ) : (
                                    <CustomDropdown
                                        options={roleOptions}
                                        value={filters.roleID}
                                        onChange={handleChange}
                                        name="roleID"
                                        placeholder="Chọn vai trò người dùng để lọc"
                                    />
                                )}
                            </FormControl>

                            {/* Email Verification */}
                            <FormControl>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Xác thực email
                                </FormLabel>
                                <ModernToggleGroup
                                    options={[
                                        {label: 'Đã xác thực', value: 'true'},
                                        {label: 'Chưa xác thực', value: 'false'}
                                    ]}
                                    value={filters.emailVerify}
                                    onChange={handleToggleChange}
                                    name="emailVerify"
                                />
                            </FormControl>

                            {/* Phone Verification */}
                            <FormControl>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Xác thực số điện thoại
                                </FormLabel>
                                <ModernToggleGroup
                                    options={[
                                        {label: 'Đã xác thực', value: 'true'},
                                        {label: 'Chưa xác thực', value: 'false'}
                                    ]}
                                    value={filters.phoneVerify}
                                    onChange={handleToggleChange}
                                    name="phoneVerify"
                                />
                            </FormControl>

                            {/* Status Selection */}
                            <FormControl>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Trạng thái người dùng
                                </FormLabel>
                                <ModernToggleGroup
                                    options={statusOptions}
                                    value={filters.status}
                                    onChange={handleToggleChange}
                                    name="status"
                                />
                            </FormControl>

                            {/* Sort By */}
                            <FormControl>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Sắp xếp bởi
                                </FormLabel>
                                <CustomDropdown
                                    options={sortByOptions}
                                    value={filters.sortBy}
                                    onChange={handleChange}
                                    name="sortBy"
                                    placeholder="Chọn trường để sắp xếp"
                                />
                            </FormControl>

                            {/* Sort Order */}
                            <FormControl>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Sắp xếp theo
                                </FormLabel>
                                <ModernToggleGroup
                                    options={sortOrderOptions}
                                    value={filters.sortOrder}
                                    onChange={handleToggleChange}
                                    name="sortOrder"
                                />
                            </FormControl>

                            {/* Date Range */}
                            <FormControl>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Lọc theo thời gian cập nhật
                                </FormLabel>
                                <VStack spacing={3} align="stretch">
                                    <Box>
                                        <Text fontSize="xs" mb={1} color="gray.600">
                                            Từ
                                        </Text>
                                        <InputGroup size="sm">
                                            <Input
                                                type="date"
                                                name="updatedAtStartFrom"
                                                value={filters.updatedAtStartFrom}
                                                onChange={handleChange}
                                                height="36px"
                                                borderRadius="md"
                                                max={new Date().toISOString().split('T')[0]} // Limit to current date
                                            />
                                            <InputRightElement
                                                pointerEvents="none"
                                                height="36px"
                                                children={<FiCalendar color="gray.400"/>}
                                            />
                                        </InputGroup>
                                    </Box>
                                    <Box>
                                        <Text fontSize="xs" mb={1} color="gray.600">
                                            Đến
                                        </Text>
                                        <InputGroup size="sm">
                                            <Input
                                                type="date"
                                                name="updatedAtEndFrom"
                                                value={filters.updatedAtEndFrom}
                                                onChange={handleChange}
                                                height="36px"
                                                borderRadius="md"
                                                max={new Date().toISOString().split('T')[0]} // Limit to current date
                                                min={filters.updatedAtStartFrom} // Can't be before start date
                                                disabled={!filters.updatedAtStartFrom} // Disable if no start date
                                            />
                                            <InputRightElement
                                                pointerEvents="none"
                                                height="36px"
                                                children={<FiCalendar color="gray.400"/>}
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
                            _hover={{bg: buttonHoverColor}}
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
                            _hover={{bg: "blue.600"}}
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

export default UserFilterDropdown;