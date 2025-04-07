import React, {useCallback, useEffect, useState} from 'react';
import {
    Alert,
    AlertIcon,
    Avatar,
    Badge,
    Box,
    Button,
    Container,
    Flex,
    HStack,
    IconButton,
    Input,
    InputGroup,
    InputLeftElement,
    Menu,
    MenuButton,
    MenuItem,
    MenuList,
    Select,
    Table,
    Tag,
    TagLabel,
    TagLeftIcon,
    Tbody,
    Td,
    Text,
    Th,
    Thead,
    Tooltip,
    Tr,
    useColorModeValue,
    useDisclosure,
    useToast,
    VStack,
} from '@chakra-ui/react';
import {
    FiCheck,
    FiChevronDown,
    FiChevronLeft,
    FiChevronRight,
    FiEdit2,
    FiFilter,
    FiMail,
    FiPlus,
    FiRefreshCw,
    FiSearch,
    FiShield,
    FiUser,
    FiX,
} from 'react-icons/fi';
import userService from "../../../services/userService.js";
import UserFilterDropdown from '../../../components/user-management/UserFilterDropdown.jsx';
import UserPermissionsModal from "../../../components/user-management/UserPermissionsModal.jsx";

const UserManagementComponent = () => {
    // State for users and loading
    const [users, setUsers] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    // State for pagination
    const [totalCount, setTotalCount] = useState(0);
    const [currentPage, setCurrentPage] = useState(1);
    const [rowsPerPage, setRowsPerPage] = useState(10);

    // State for selected user for permissions
    const [selectedUser, setSelectedUser] = useState(null);
    const [isLoadingPermissions, setIsLoadingPermissions] = useState(false);

    // State for filters (filtered params that will be sent to the server)
    const [activeFilters, setActiveFilters] = useState({});

    // State for UI filter tracking
    const [uiFilters, setUiFilters] = useState({
        sortBy: '',
        sortOrder: 'asc',
        emailVerify: '',
        phoneVerify: '',
        status: '',
        updatedAtStartFrom: '',
        updatedAtEndFrom: '',
        roleID: '',
    });

    // Search state (quick filter in UI)
    const [searchField, setSearchField] = useState('');
    const [searchQuery, setSearchQuery] = useState('');

    // Modal controls - Must be defined before any conditional logic
    const permissionsModal = useDisclosure();

    // Toast for notifications
    const toast = useToast();

    // Theme colors
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const hoverBgColor = useColorModeValue('blue.50', 'gray.700');
    const scrollTrackBg = useColorModeValue('#f1f1f1', '#2d3748');
    const scrollThumbBg = useColorModeValue('#c1c1c1', '#4a5568');
    const scrollThumbHoverBg = useColorModeValue('#a1a1a1', '#718096');

    // Create params object with only needed parameters - All useCallback functions must be defined in the same order
    const createRequestParams = useCallback(() => {
        const params = {
            page: currentPage,
            limit: rowsPerPage,
        };

        if (Object.keys(activeFilters).length > 0) {
            Object.entries(activeFilters).forEach(([key, value]) => {
                if (value !== '') {
                    params[key] = value;
                }
            });
        }

        return params;
    }, [currentPage, rowsPerPage, activeFilters]);

    // Fetch users from API
    const fetchUsers = useCallback(async () => {
        setIsLoading(true);
        setIsError(false);

        try {
            const params = createRequestParams();
            const response = await userService.getUsers(params);

            if (response && response.data) {
                setUsers(response.data || []);

                // Handle pagination metadata from API
                if (response.metadata) {
                    // Set total count
                    if (response.metadata.pagination && response.metadata.pagination.total_items) {
                        setTotalCount(response.metadata.pagination.total_items);
                    } else if (response.metadata.total_count) {
                        setTotalCount(response.metadata.total_count);
                    } else {
                        setTotalCount(response.data.length);
                    }

                    // Update pagination
                    if (response.metadata.pagination) {
                        const pagination = response.metadata.pagination;
                        setTotalPages(pagination.total_pages || 1);
                    }
                } else {
                    setTotalCount(response.data.length);
                }
            } else {
                setUsers([]);
                setTotalCount(0);
            }
        } catch (error) {
            setIsError(true);
            setErrorMessage(error.response?.data?.error?.message || 'Failed to fetch users');

            toast({
                title: 'Error loading users',
                description: error.response?.data?.error?.message || 'An error occurred while loading users',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });

            setUsers([]);
            setTotalCount(0);
        } finally {
            setIsLoading(false);
        }
    }, [createRequestParams, toast]);

    // Load users when component mounts or when dependencies change
    useEffect(() => {
        fetchUsers();
    }, [fetchUsers]);

    // Open permissions modal and load user
    const handleOpenPermissions = async (user) => {
        setSelectedUser(user);
        setIsLoadingPermissions(true);

        try {
            permissionsModal.onOpen();
        } catch (error) {
            toast({
                title: 'Error loading permissions',
                description: error.message || 'An error occurred while loading user permissions',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setIsLoadingPermissions(false);
        }
    };

    // Save user permissions
    const handleSavePermissions = async (userId, permissions) => {
        try {
            await userService.updateUserPermissions(userId, permissions);
            fetchUsers();
            return true;
        } catch (error) {
            console.error('Error updating permissions:', error);
            throw error;
        }
    };

    // Apply filters from modal
    const handleApplyFilters = (newFilters) => {
        setUiFilters(newFilters);

        const filteredParams = {};
        Object.entries(newFilters).forEach(([key, value]) => {
            if (value !== '') {
                filteredParams[key] = value;
            }
        });

        setActiveFilters(prevFilters => ({
            ...prevFilters,
            ...filteredParams
        }));

        setCurrentPage(1);
    };

    // Handle search
    const handleSearch = () => {
        if (searchQuery && searchField) {
            setActiveFilters(prevFilters => ({
                ...prevFilters,
                searchBy: searchField,
                searchValue: searchQuery
            }));
        } else {
            const { searchBy, searchValue, ...restFilters } = activeFilters;
            setActiveFilters(restFilters);
        }

        setCurrentPage(1);
    };

    // Handle search input change
    const handleSearchChange = (e) => {
        setSearchQuery(e.target.value);
    };

    // Handle search field change
    const handleSearchFieldChange = (e) => {
        setSearchField(e.target.value);
    };

    // Handle key press in search input
    const handleKeyPress = (e) => {
        if (e.key === 'Enter') {
            handleSearch();
        }
    };

    // Clear all filters
    const clearFilters = () => {
        setUiFilters({
            sortBy: '',
            sortOrder: 'asc',
            emailVerify: '',
            phoneVerify: '',
            status: '',
            updatedAtStartFrom: '',
            updatedAtEndFrom: '',
            roleID: '',
        });

        const {searchBy, searchValue, ...rest} = activeFilters;

        const newFilters = {};
        Object.entries(rest).forEach(([key, value]) => {
            if (!['sortBy', 'sortOrder', 'emailVerify', 'phoneVerify', 'status',
                'updatedAtStartFrom', 'updatedAtEndFrom', 'roleID'].includes(key)) {
                newFilters[key] = value;
            }
        });

        setActiveFilters(newFilters);
        setSearchQuery('');
        setSearchField('');
    };

    // Clear search
    const clearSearch = () => {
        setSearchQuery('');
        const {searchBy, searchValue, ...restFilters} = activeFilters;
        setActiveFilters(restFilters);
    };

    // Pagination logic - Use metadata from API when available
    const [totalPages, setTotalPages] = useState(1);

    // Update totalPages when API response includes pagination metadata
    useEffect(() => {
        // Only calculate if we don't have API-provided totalPages
        if (!activeFilters.apiTotalPages) {
            // Fallback to calculated value
            setTotalPages(Math.max(1, Math.ceil(totalCount / rowsPerPage)));
        }
    }, [totalCount, rowsPerPage, activeFilters.apiTotalPages]);

    // Generate pagination range
    const generatePaginationRange = (current, total) => {
        current = Math.max(1, Math.min(current, total));

        if (total <= 5) {
            return Array.from({length: total}, (_, i) => i + 1);
        }

        if (current <= 3) {
            return [1, 2, 3, 4, 5, '...', total];
        }

        if (current >= total - 2) {
            return [1, '...', total - 4, total - 3, total - 2, total - 1, total];
        }

        return [1, '...', current - 1, current, current + 1, '...', total];
    };

    // Get pagination range
    const paginationRange = generatePaginationRange(currentPage, totalPages);

    // Format date string
    const formatDate = (dateString) => {
        if (!dateString) return 'N/A';

        try {
            const date = new Date(dateString);
            return date.toLocaleString();
        } catch (e) {
            return dateString;
        }
    };

    // Get color for role badge
    const getRoleColor = (roleName) => {
        switch (roleName?.toLowerCase()) {
            case 'admin':
                return 'purple';
            case 'supplier':
                return 'blue';
            case 'deliverer':
                return 'green';
            case 'user':
                return 'gray';
            default:
                return 'gray';
        }
    };

    // Check if there are any active UI filters
    const hasActiveUiFilters = Object.values(uiFilters).some(value => value !== '');

    // Check if search is active
    const hasActiveSearch = searchQuery && searchField;

    return (
        <Container maxW="container.xl" py={6}>
            {/* Search and Filter Bar */}
            <Flex
                justifyContent="space-between"
                alignItems="center"
                p={4}
                mb={4}
                flexDir={{base: 'column', md: 'row'}}
                gap={{base: 4, md: 0}}
            >
                <Flex
                    flex={{md: 1}}
                    direction={{base: "column", sm: "row"}}
                    gap={3}
                    align={{base: "stretch", sm: "center"}}
                >
                    {/* Enhanced Search Input */}
                    <Flex
                        borderWidth="1px"
                        borderRadius="lg"
                        overflow="hidden"
                        align="center"
                        bg={bgColor}
                        shadow="sm"
                        flex="1"
                        maxW={{base: "full", lg: "450px"}}
                    >
                        {/* Search Field Dropdown */}
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
                            _focus={{boxShadow: "none"}}
                            fontSize="sm"
                        >
                            <option value="">Select field</option>
                            <option value="fullname">Name</option>
                            <option value="email">Email</option>
                            <option value="phone">Phone</option>
                        </Select>

                        <InputGroup size="md" variant="unstyled">
                            <InputLeftElement pointerEvents="none" h="full" pl={3}>
                                <FiSearch color="gray.400"/>
                            </InputLeftElement>
                            <Input
                                placeholder={searchField ? `Search by ${searchField.toLowerCase()}...` : "Select a field first"}
                                pl={10}
                                pr={2}
                                py={2.5}
                                value={searchQuery}
                                onChange={handleSearchChange}
                                onKeyPress={handleKeyPress}
                                _placeholder={{color: "gray.400"}}
                                isDisabled={!searchField}
                            />
                        </InputGroup>

                        {/* Search actions */}
                        <HStack>
                            {hasActiveSearch && (
                                <Tooltip label="Clear search" hasArrow>
                                    <IconButton
                                        icon={<FiX size={16}/>}
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
                                    icon={<FiSearch size={16}/>}
                                    onClick={handleSearch}
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

                    {/* Refresh Button */}
                    <Tooltip label="Refresh data" hasArrow>
                        <IconButton
                            icon={<FiRefreshCw/>}
                            onClick={fetchUsers}
                            aria-label="Refresh data"
                            variant="ghost"
                            colorScheme="blue"
                            size="sm"
                            isLoading={isLoading}
                        />
                    </Tooltip>
                </Flex>

                {/* Actions */}
                <HStack spacing={2}>
                    {/* Filter Button */}
                    <UserFilterDropdown
                        onApplyFilters={handleApplyFilters}
                        currentFilters={uiFilters}
                    />

                    {/* Create Button */}
                    <Button
                        leftIcon={<FiPlus/>}
                        colorScheme="blue"
                        size="sm"
                        borderRadius="md"
                        fontWeight="normal"
                        px={4}
                        shadow="md"
                        bgGradient="linear(to-r, blue.400, blue.500)"
                        color="white"
                        _hover={{
                            bgGradient: "linear(to-r, blue.500, blue.600)",
                            shadow: 'lg',
                            transform: 'translateY(-1px)'
                        }}
                        _active={{
                            bgGradient: "linear(to-r, blue.600, blue.700)",
                            transform: 'translateY(0)',
                            shadow: 'md'
                        }}
                        transition="all 0.2s"
                    >
                        Create
                    </Button>
                </HStack>
            </Flex>

            {/* Error Alert */}
            {isError && (
                <Alert status="error" variant="left-accent" mb={4} borderRadius="md">
                    <AlertIcon/>
                    <Text>{errorMessage || 'An error occurred while fetching users'}</Text>
                </Alert>
            )}

            {/* Table Container */}
            <Box
                width="100%"
                borderRadius="xl"
                overflow="hidden"
                boxShadow="lg"
                bg={bgColor}
                display="flex"
                flexDirection="column"
                borderWidth="1px"
                borderColor={borderColor}
            >
                {/* Data Table Container with Fixed Height */}
                <Box
                    overflow="auto"
                    sx={{
                        '&::-webkit-scrollbar': {
                            width: '8px',
                            height: '8px',
                        },
                        '&::-webkit-scrollbar-track': {
                            background: scrollTrackBg,
                            borderRadius: '4px',
                        },
                        '&::-webkit-scrollbar-thumb': {
                            background: scrollThumbBg,
                            borderRadius: '4px',
                        },
                        '&::-webkit-scrollbar-thumb:hover': {
                            background: scrollThumbHoverBg,
                        },
                    }}
                    flex="1"
                    minH="300px"
                    maxH={{base: "60vh", lg: "calc(100vh - 250px)"}}
                    borderBottomWidth="1px"
                    borderColor={borderColor}
                >
                    <Table variant="simple" size="md" colorScheme="gray"
                           style={{borderCollapse: 'separate', borderSpacing: '0'}}>
                        <Thead bg={useColorModeValue('gray.50', 'gray.900')} position="sticky" top={0} zIndex={1}>
                            <Tr>
                                <Th
                                    py={4}
                                    borderTopLeftRadius="md"
                                    fontSize="xs"
                                    color={useColorModeValue('gray.600', 'gray.300')}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    User
                                </Th>
                                <Th
                                    py={4}
                                    display={{base: "none", md: "table-cell"}}
                                    fontSize="xs"
                                    color={useColorModeValue('gray.600', 'gray.300')}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Contact Info
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={useColorModeValue('gray.600', 'gray.300')}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Role
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={useColorModeValue('gray.600', 'gray.300')}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Status
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={useColorModeValue('gray.600', 'gray.300')}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Last Updated
                                </Th>
                                <Th
                                    py={4}
                                    textAlign="right"
                                    borderTopRightRadius="md"
                                    fontSize="xs"
                                    color={useColorModeValue('gray.600', 'gray.300')}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Actions
                                </Th>
                            </Tr>
                        </Thead>
                        <Tbody>
                            {isLoading ? (
                                <Tr>
                                    <Td colSpan={6} textAlign="center" py={12}>
                                        <Flex justify="center" align="center" direction="column">
                                            <Box
                                                h="40px"
                                                w="40px"
                                                borderWidth="3px"
                                                borderStyle="solid"
                                                borderColor="blue.500"
                                                borderTopColor="transparent"
                                                borderRadius="50%"
                                                animation="spin 1s linear infinite"
                                                mb={3}
                                                sx={{
                                                    '@keyframes spin': {
                                                        '0%': { transform: 'rotate(0deg)' },
                                                        '100%': { transform: 'rotate(360deg)' }
                                                    }
                                                }}
                                            />
                                            <Text color="gray.500" fontSize="sm">Loading users...</Text>
                                        </Flex>
                                    </Td>
                                </Tr>
                            ) : users.length > 0 ? (
                                users.map((user) => (
                                    <Tr
                                        key={user.id}
                                        _hover={{bg: useColorModeValue('blue.50', 'gray.700')}}
                                        transition="background-color 0.2s"
                                        cursor="pointer"
                                        borderBottomWidth="1px"
                                        borderColor={borderColor}
                                        _active={{bg: useColorModeValue('blue.100', 'gray.600')}}
                                        h="60px"
                                    >
                                        <Td>
                                            <HStack spacing={3}>
                                                <Avatar size="sm" name={user.fullname}
                                                        src={user.avatar_url || "/api/placeholder/40/40"}/>
                                                <Box>
                                                    <Text
                                                        fontWeight="medium"
                                                        fontSize="sm"
                                                        color={useColorModeValue('gray.800', 'white')}
                                                    >
                                                        {user.fullname}
                                                    </Text>
                                                    <Text
                                                        fontSize="xs"
                                                        color={useColorModeValue('gray.500', 'gray.400')}
                                                        display={{base: "none", lg: "block"}}
                                                    >
                                                        {user.email}
                                                    </Text>
                                                </Box>
                                            </HStack>
                                        </Td>
                                        <Td display={{base: "none", md: "table-cell"}}>
                                            <VStack align="start" spacing={0.5}>
                                                <HStack spacing={1} align="center">
                                                    <FiMail size={12}/>
                                                    <Text
                                                        fontSize="sm"
                                                        color={useColorModeValue('gray.600', 'gray.300')}
                                                    >
                                                        {user.email}
                                                    </Text>
                                                    {user.email_verify ? (
                                                        <Tooltip label="Email verified" hasArrow>
                                                            <Badge colorScheme="green" variant="outline" fontSize="2xs"
                                                                   p={0.5}>
                                                                <FiCheck size={10}/>
                                                            </Badge>
                                                        </Tooltip>
                                                    ) : (
                                                        <Tooltip label="Email not verified" hasArrow>
                                                            <Badge colorScheme="red" variant="outline" fontSize="2xs"
                                                                   p={0.5}>
                                                                <FiX size={10}/>
                                                            </Badge>
                                                        </Tooltip>
                                                    )}
                                                </HStack>
                                                <HStack spacing={1} align="center">
                                                    <Text
                                                        fontSize="sm"
                                                        color={useColorModeValue('gray.500', 'gray.400')}
                                                    >
                                                        {user.phone || 'No phone'}
                                                    </Text>
                                                    {user.phone && user.phone_verify && (
                                                        <Tooltip label="Phone verified" hasArrow>
                                                            <Badge colorScheme="green" variant="outline" fontSize="2xs"
                                                                   p={0.5}>
                                                                <FiCheck size={10}/>
                                                            </Badge>
                                                        </Tooltip>
                                                    )}
                                                </HStack>
                                            </VStack>
                                        </Td>
                                        <Td>
                                            <Tag
                                                size="md"
                                                variant="subtle"
                                                colorScheme={getRoleColor(user.role_name)}
                                                borderRadius="md"
                                            >
                                                <TagLeftIcon as={FiUser} boxSize="12px"/>
                                                <TagLabel fontSize="xs" fontWeight="medium">{user.role_name}</TagLabel>
                                            </Tag>
                                        </Td>
                                        <Td>
                                            <Badge
                                                px={2}
                                                py={1}
                                                borderRadius="full"
                                                colorScheme={user.status === 'active' ? 'green' : 'red'}
                                                textTransform="capitalize"
                                                fontWeight="medium"
                                                fontSize="xs"
                                            >
                                                {user.status}
                                            </Badge>
                                        </Td>
                                        <Td>
                                            <Text
                                                fontSize="sm"
                                                color={useColorModeValue('gray.500', 'gray.400')}
                                            >
                                                {formatDate(user.updated_at)}
                                            </Text>
                                        </Td>
                                        <Td textAlign="right">
                                            <HStack spacing={1} justifyContent="flex-end">
                                                <Tooltip label="Edit user" hasArrow>
                                                    <IconButton
                                                        icon={<FiEdit2 size={15}/>}
                                                        size="sm"
                                                        variant="ghost"
                                                        colorScheme="blue"
                                                        aria-label="Edit user"
                                                        borderRadius="md"
                                                    />
                                                </Tooltip>
                                                <Tooltip label="Manage permissions" hasArrow>
                                                    <IconButton
                                                        icon={<FiShield size={15}/>}
                                                        size="sm"
                                                        variant="ghost"
                                                        colorScheme="purple"
                                                        aria-label="Manage permissions"
                                                        borderRadius="md"
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            handleOpenPermissions(user);
                                                        }}
                                                    />
                                                </Tooltip>
                                            </HStack>
                                        </Td>
                                    </Tr>
                                ))
                            ) : (
                                <Tr>
                                    <Td colSpan={6} textAlign="center" py={12}>
                                        <Flex direction="column" align="center" justify="center" py={8}>
                                            <Box color="gray.400" mb={3}>
                                                <FiSearch size={36}/>
                                            </Box>
                                            <Text fontWeight="normal" color="gray.500" fontSize="md">No users
                                                found</Text>
                                            <Text color="gray.400" fontSize="sm" mt={1}>Try a different search term or
                                                filter</Text>
                                        </Flex>
                                    </Td>
                                </Tr>
                            )}
                        </Tbody>
                    </Table>
                </Box>

                {/* Fixed Pagination Section */}
                <Box
                    borderTop="1px"
                    borderColor={borderColor}
                    bg={useColorModeValue('gray.50', 'gray.800')}
                    bgGradient={useColorModeValue(
                        "linear(to-r, white, gray.50, white)",
                        "linear(to-r, gray.800, gray.700, gray.800)"
                    )}
                    position="sticky"
                    bottom="0"
                    width="100%"
                    zIndex="1"
                    boxShadow="0 -2px 6px rgba(0,0,0,0.05)"
                >
                    <Flex
                        justifyContent="space-between"
                        alignItems="center"
                        py={4}
                        px={6}
                        flexWrap={{base: "wrap", md: "nowrap"}}
                        gap={4}
                    >
                        <HStack spacing={1} flexShrink={0}>
                            <Text fontSize="sm" color="gray.600" fontWeight="normal">
                                Showing {totalCount > 0 ? (currentPage - 1) * rowsPerPage + 1 : 0}-{Math.min(currentPage * rowsPerPage, totalCount)} of {totalCount} users
                            </Text>
                            <Menu>
                                <MenuButton
                                    as={Button}
                                    size="xs"
                                    variant="ghost"
                                    rightIcon={<FiChevronDown/>}
                                    ml={2}
                                    fontWeight="normal"
                                    color="gray.600"
                                >
                                    {rowsPerPage} per page
                                </MenuButton>
                                <MenuList minW="120px" shadow="lg" borderRadius="md">
                                    <MenuItem onClick={() => setRowsPerPage(10)}>10 per page</MenuItem>
                                    <MenuItem onClick={() => setRowsPerPage(15)}>15 per page</MenuItem>
                                    <MenuItem onClick={() => setRowsPerPage(20)}>20 per page</MenuItem>
                                    <MenuItem onClick={() => setRowsPerPage(50)}>50 per page</MenuItem>
                                </MenuList>
                            </Menu>
                        </HStack>

                        {totalCount > 0 && (
                            <HStack spacing={1} justify="center" width={{base: "100%", md: "auto"}}>
                                <IconButton
                                    icon={<FiChevronLeft/>}
                                    size="sm"
                                    variant="ghost"
                                    isDisabled={currentPage === 1 || isLoading}
                                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                                    aria-label="Previous page"
                                    borderRadius="md"
                                />

                                {paginationRange.map((page, index) => (
                                    page === '...' ? (
                                        <Text key={`ellipsis-${index}`} mx={1} color="gray.500">...</Text>
                                    ) : (
                                        <Button
                                            key={`page-${page}`}
                                            size="sm"
                                            variant={currentPage === page ? "solid" : "ghost"}
                                            colorScheme={currentPage === page ? "blue" : "gray"}
                                            onClick={() => typeof page === 'number' && setCurrentPage(page)}
                                            borderRadius="md"
                                            minW="32px"
                                            isDisabled={isLoading}
                                        >
                                            {page}
                                        </Button>
                                    )
                                ))}

                                <IconButton
                                    icon={<FiChevronRight/>}
                                    size="sm"
                                    variant="ghost"
                                    isDisabled={currentPage === totalPages || isLoading}
                                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                                    aria-label="Next page"
                                    borderRadius="md"
                                />
                            </HStack>
                        )}
                    </Flex>
                </Box>
            </Box>

            {/* User Permissions Modal */}
            <UserPermissionsModal
                isOpen={permissionsModal.isOpen}
                onClose={permissionsModal.onClose}
                user={selectedUser}
                onSave={handleSavePermissions}
                isLoading={isLoadingPermissions}
            />
        </Container>
    );
};

export default UserManagementComponent;