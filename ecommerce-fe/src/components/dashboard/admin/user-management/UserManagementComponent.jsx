import React, { useCallback, useEffect, useState } from 'react';
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
    Menu,
    MenuButton,
    MenuItem,
    MenuList,
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
    FiChevronDown,
    FiChevronLeft,
    FiChevronRight,
    FiEdit2,
    FiMail,
    FiPlus,
    FiRefreshCw,
    FiSearch,
    FiTrash2,
    FiUser,
    FiX,
} from 'react-icons/fi';
import { RiVerifiedBadgeFill } from "react-icons/ri";
import userService from "../../../../services/userService.js";
import UserFilterDropdown from './UserFilterDropdown.jsx';
import CreateUserModal from "./CreateUserModal.jsx";
import UserSearchComponent from './UserSearchComponent.jsx';
import EditUserModal from "./EditUserModal.jsx";

const UserManagementComponent = () => {
    // Theme colors - Define ALL color mode values at the TOP of the component
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const scrollTrackBg = useColorModeValue('#f1f1f1', '#2d3748');
    const scrollThumbBg = useColorModeValue('#c1c1c1', '#4a5568');
    const scrollThumbHoverBg = useColorModeValue('#a1a1a1', '#718096');
    const tableBgColor = useColorModeValue('gray.50', 'gray.900');
    const tableTextColor = useColorModeValue('gray.600', 'gray.300');
    const tableRowHoverBg = useColorModeValue('blue.50', 'gray.700');
    const tableRowActiveBg = useColorModeValue('blue.100', 'gray.600');
    const footerBgGradient = useColorModeValue(
        "linear(to-r, white, gray.50, white)",
        "linear(to-r, gray.800, gray.700, gray.800)"
    );
    const textColorPrimary = useColorModeValue('gray.800', 'white');
    const textColorSecondary = useColorModeValue('gray.500', 'gray.400');
    const textColorTertiary = useColorModeValue('gray.600', 'gray.300');

    // Toast for notifications
    const toast = useToast();

    // State for users and loading
    const [users, setUsers] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    // State for pagination
    const [totalCount, setTotalCount] = useState(0);
    const [currentPage, setCurrentPage] = useState(1);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [totalPages, setTotalPages] = useState(1);

    const { isOpen: isCreateModalOpen, onOpen: onOpenCreateModal, onClose: onCloseCreateModal } = useDisclosure();
    const { isOpen: isEditModalOpen, onOpen: onOpenEditModal, onClose: onCloseEditModal } = useDisclosure();
    const [selectedUser, setSelectedUser] = useState(null);

    // Unified filter state
    const [filters, setFilters] = useState({
        sortBy: '',
        sortOrder: 'asc',
        emailVerify: '',
        phoneVerify: '',
        status: '',
        updatedAtStartFrom: '',
        updatedAtEndFrom: '',
        roleID: '',
    });

    // State for active filter parameters that will be sent to the API
    const [activeFilterParams, setActiveFilterParams] = useState({});

    // Search state (will be managed by UserSearchComponent)
    const [searchParams, setSearchParams] = useState({
        searchBy: '',
        searchValue: ''
    });

    // Create params object with only needed parameters
    const createRequestParams = useCallback(() => {
        const params = {
            page: currentPage,
            limit: rowsPerPage,
        };

        Object.entries(activeFilterParams).forEach(([key, value]) => {
            if (value !== '') {
                params[key] = value;
            }
        });

        if (searchParams.searchBy && searchParams.searchValue && searchParams.searchValue.trim() !== '') {
            params.searchBy = searchParams.searchBy;
            params.searchValue = searchParams.searchValue;
        }

        return params;
    }, [currentPage, rowsPerPage, activeFilterParams, searchParams]);

    // Fetch users from API
    const fetchUsers = useCallback(async () => {
        setIsLoading(true);
        setIsError(false);

        try {
            const params = createRequestParams();
            const response = await userService.getUsers(params);

            if (response && response.data) {
                setUsers(response.data || []);

                if (response.metadata) {
                    if (response.metadata.pagination && response.metadata.pagination.total_items) {
                        setTotalCount(response.metadata.pagination.total_items);
                    } else if (response.metadata.total_count) {
                        setTotalCount(response.metadata.total_count);
                    } else {
                        setTotalCount(response.data.length);
                    }

                    if (response.metadata.pagination) {
                        const pagination = response.metadata.pagination;
                        setTotalPages(pagination.total_pages || 1);
                    } else {
                        setTotalPages(Math.max(1, Math.ceil((response.data.length || 0) / rowsPerPage)));
                    }
                } else {
                    setTotalCount(response.data.length);
                    setTotalPages(Math.max(1, Math.ceil((response.data.length || 0) / rowsPerPage)));
                }
            } else {
                setUsers([]);
                setTotalCount(0);
                setTotalPages(1);
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
            setTotalPages(1);
        } finally {
            setIsLoading(false);
        }
    }, [createRequestParams, rowsPerPage, toast]);

    // Handle search from UserSearchComponent
    const handleSearch = (searchField, searchValue) => {
        setSearchParams({
            searchBy: searchField,
            searchValue: searchValue
        });
    };

    // Handle filter changes (used by UserFilterDropdown)
    const handleFiltersChange = (updatedFilters) => {
        setFilters(updatedFilters);
    };

    // Apply filters from modal
    const handleApplyFilters = (filteredParams) => {
        setActiveFilterParams(filteredParams);
        setCurrentPage(1);
    };

    // Load users when component mounts or when dependencies change
    useEffect(() => {
        if (searchParams.searchBy || searchParams.searchValue) {
            setCurrentPage(1);
        }
        fetchUsers();
    }, [fetchUsers, searchParams]);

    // Handle user created and reload current page
    const handleUserCreated = async () => {
        await fetchUsers(); // Reload dữ liệu đúng trang hiện tại
    };

    // Handle user updated and reload current page
    const handleUserUpdated = async () => {
        await fetchUsers(); // Reload dữ liệu đúng trang hiện tại
    };

    // Handle user deleted and reload current page
    const handleDeleteUser = async (userId) => {
        if (window.confirm('Are you sure you want to delete this user?')) {
            try {
                await userService.deleteUser(userId);
                toast({
                    title: "User deleted successfully",
                    status: "success",
                    duration: 3000,
                    isClosable: true,
                });
                await fetchUsers(); // Reload dữ liệu đúng trang hiện tại
            } catch (error) {
                toast({
                    title: "Failed to delete user",
                    description: error.response?.data?.error?.message || 'An unexpected error occurred',
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            }
        }
    };

    // Handle opening Edit Modal
    const handleOpenEditModal = (user) => {
        setSelectedUser(user);
        onOpenEditModal();
    };

    // Generate pagination range
    const generatePaginationRange = (current, total) => {
        current = Math.max(1, Math.min(current, total));

        if (total <= 5) {
            return Array.from({ length: total }, (_, i) => i + 1);
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
        if (!roleName) return 'gray';

        switch (roleName.toLowerCase()) {
            case 'admin':
                return 'purple';
            case 'supplier':
                return 'blue';
            case 'deliverer':
                return 'green';
            case 'user':
            case 'customer':
                return 'gray';
            default:
                return 'gray';
        }
    };

    // Check if there are any active filters
    const hasActiveFilters = Object.values(filters).some(value => value !== '');

    // Check if search is active
    const hasActiveSearch = searchParams.searchBy && searchParams.searchValue;

    return (
        <Container maxW="container.xl" py={6}>
            <Flex
                justifyContent="space-between"
                alignItems="center"
                p={4}
                mb={4}
                flexDir={{ base: 'column', md: 'row' }}
                gap={{ base: 4, md: 0 }}
            >
                <Flex
                    flex={{ md: 1 }}
                    direction={{ base: "column", sm: "row" }}
                    gap={3}
                    align={{ base: "stretch", sm: "center" }}
                >
                    <UserSearchComponent
                        onSearch={handleSearch}
                        isLoading={isLoading}
                    />
                    <Tooltip label="Refresh data" hasArrow>
                        <IconButton
                            icon={<FiRefreshCw />}
                            onClick={fetchUsers}
                            aria-label="Refresh data"
                            variant="ghost"
                            colorScheme="blue"
                            size="sm"
                            isLoading={isLoading}
                        />
                    </Tooltip>
                </Flex>

                <HStack spacing={2}>
                    <UserFilterDropdown
                        filters={filters}
                        onFiltersChange={handleFiltersChange}
                        onApplyFilters={handleApplyFilters}
                    />
                    <Button
                        leftIcon={<FiPlus />}
                        colorScheme="blue"
                        size="sm"
                        borderRadius="md"
                        fontWeight="normal"
                        px={4}
                        shadow="md"
                        onClick={onOpenCreateModal}
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

            {isError && (
                <Alert status="error" variant="left-accent" mb={4} borderRadius="md">
                    <AlertIcon />
                    <Text>{errorMessage || 'An error occurred while fetching users'}</Text>
                </Alert>
            )}

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
                    maxH={{ base: "60vh", lg: "calc(100vh - 250px)" }}
                    borderBottomWidth="1px"
                    borderColor={borderColor}
                >
                    <Table variant="simple" size="md" colorScheme="gray" style={{ borderCollapse: 'separate', borderSpacing: '0' }}>
                        <Thead bg={tableBgColor} position="sticky" top={0} zIndex={1}>
                            <Tr>
                                <Th
                                    py={4}
                                    borderTopLeftRadius="md"
                                    fontSize="xs"
                                    color={tableTextColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    User
                                </Th>
                                <Th
                                    py={4}
                                    display={{ base: "none", md: "table-cell" }}
                                    fontSize="xs"
                                    color={tableTextColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Contact Info
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={tableTextColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Role
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={tableTextColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Status
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={tableTextColor}
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
                                    color={tableTextColor}
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
                                        _hover={{ bg: tableRowHoverBg }}
                                        transition="background-color 0.2s"
                                        borderBottomWidth="1px"
                                        borderColor={borderColor}
                                        _active={{ bg: tableRowActiveBg }}
                                        h="60px"
                                        bg={bgColor}
                                    >
                                        <Td>
                                            <HStack spacing={3}>
                                                <Avatar size="sm" name={user.fullname} src={user.avatar_url || "/api/placeholder/40/40"} />
                                                <Box>
                                                    <Text
                                                        fontWeight="medium"
                                                        fontSize="sm"
                                                        color={textColorPrimary}
                                                    >
                                                        {user.fullname}
                                                    </Text>
                                                    <Text
                                                        fontSize="xs"
                                                        color={textColorSecondary}
                                                        display={{ base: "none", lg: "block" }}
                                                    >
                                                        {user.email}
                                                    </Text>
                                                </Box>
                                            </HStack>
                                        </Td>
                                        <Td display={{ base: "none", md: "table-cell" }}>
                                            <VStack align="start" spacing={0.5}>
                                                <HStack spacing={1} align="center">
                                                    <FiMail size={12} />
                                                    <Text
                                                        fontSize="sm"
                                                        color={textColorTertiary}
                                                    >
                                                        {user.email}
                                                    </Text>
                                                    {user.email_verify ? (
                                                        <Tooltip label="Email verified" hasArrow>
                                                            <Box display="inline-block" ml={1} color="blue.400">
                                                                <RiVerifiedBadgeFill size={17} />
                                                            </Box>
                                                        </Tooltip>
                                                    ) : (
                                                        <Tooltip label="Email not verified" hasArrow>
                                                            <Badge colorScheme="red" variant="outline" fontSize="2xs" p={0.5}>
                                                                <FiX size={13} />
                                                            </Badge>
                                                        </Tooltip>
                                                    )}
                                                </HStack>
                                                <HStack spacing={1} align="center">
                                                    <Text
                                                        fontSize="sm"
                                                        color={textColorSecondary}
                                                    >
                                                        {user.phone || 'No phone'}
                                                    </Text>
                                                    {user.phone && user.phone_verify && (
                                                        <Tooltip label="Phone verified" hasArrow>
                                                            <Box display="inline-block" ml={1} color="blue.400">
                                                                <RiVerifiedBadgeFill size={15} />
                                                            </Box>
                                                        </Tooltip>
                                                    )}
                                                </HStack>
                                            </VStack>
                                        </Td>
                                        <Td>
                                            <HStack spacing={1} flexWrap="wrap">
                                                {user.roles && user.roles.length > 0 ? (
                                                    user.roles.map((role, index) => (
                                                        <Tag
                                                            key={`${user.id}-role-${index}`}
                                                            size="md"
                                                            variant="subtle"
                                                            colorScheme={getRoleColor(role.name)}
                                                            borderRadius="md"
                                                            mb={1}
                                                        >
                                                            <TagLeftIcon as={FiUser} boxSize="12px" />
                                                            <TagLabel fontSize="xs" fontWeight="medium">{role.name}</TagLabel>
                                                        </Tag>
                                                    ))
                                                ) : (
                                                    <Tag
                                                        size="md"
                                                        variant="subtle"
                                                        colorScheme="gray"
                                                        borderRadius="md"
                                                    >
                                                        <TagLeftIcon as={FiUser} boxSize="12px" />
                                                        <TagLabel fontSize="xs" fontWeight="medium">No Role</TagLabel>
                                                    </Tag>
                                                )}
                                            </HStack>
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
                                                color={textColorSecondary}
                                            >
                                                {formatDate(user.updated_at)}
                                            </Text>
                                        </Td>
                                        <Td textAlign="right">
                                            <HStack spacing={1} justifyContent="flex-end">
                                                <Tooltip label="Edit user" hasArrow>
                                                    <IconButton
                                                        icon={<FiEdit2 size={15} />}
                                                        size="sm"
                                                        variant="ghost"
                                                        colorScheme="blue"
                                                        aria-label="Edit user"
                                                        borderRadius="md"
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            handleOpenEditModal(user);
                                                        }}
                                                    />
                                                </Tooltip>
                                                <Tooltip label="Delete user" hasArrow>
                                                    <IconButton
                                                        icon={<FiTrash2 size={15} />}
                                                        size="sm"
                                                        variant="ghost"
                                                        colorScheme="red"
                                                        aria-label="Delete user"
                                                        borderRadius="md"
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            handleDeleteUser(user.id);
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
                                                <FiSearch size={36} />
                                            </Box>
                                            <Text fontWeight="normal" color="gray.500" fontSize="md">No users found</Text>
                                            <Text color="gray.400" fontSize="sm" mt={1}>Try a different search term or filter</Text>
                                        </Flex>
                                    </Td>
                                </Tr>
                            )}
                        </Tbody>
                    </Table>
                </Box>

                <Box
                    borderTop="1px"
                    borderColor={borderColor}
                    bg={useColorModeValue('gray.50', 'gray.800')}
                    bgGradient={footerBgGradient}
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
                        flexWrap={{ base: "wrap", md: "nowrap" }}
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
                                    rightIcon={<FiChevronDown />}
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
                            <HStack spacing={1} justify="center" width={{ base: "100%", md: "auto" }}>
                                <IconButton
                                    icon={<FiChevronLeft />}
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
                                    icon={<FiChevronRight />}
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

            <CreateUserModal
                isOpen={isCreateModalOpen}
                onClose={onCloseCreateModal}
                onUserCreated={handleUserCreated}
            />

            {selectedUser && (
                <EditUserModal
                    isOpen={isEditModalOpen}
                    onClose={onCloseEditModal}
                    user={selectedUser}
                    onUserUpdated={handleUserUpdated}
                />
            )}
        </Container>
    );
};

export default UserManagementComponent;