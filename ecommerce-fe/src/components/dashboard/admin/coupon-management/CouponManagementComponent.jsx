import React, { useCallback, useEffect, useState } from 'react';
import {
    Alert,
    AlertIcon,
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
    Code,
} from '@chakra-ui/react';
import {
    FiChevronDown,
    FiChevronLeft,
    FiChevronRight,
    FiEdit2,
    FiPlus,
    FiRefreshCw,
    FiSearch,
    FiTrash2,
    FiTag,
    FiPercent,
    FiDollarSign,
    FiCalendar,
    FiUsers,
} from 'react-icons/fi';
import couponService from "../../../../services/couponService.js";
import CouponFilterDropdown from './CouponFilterDropdown.jsx';
import CreateCouponModal from "./CreateCouponModal.jsx";
import CouponSearchComponent from './CouponSearchComponent.jsx';
import EditCouponModal from "./EditCouponModal.jsx";

const CouponManagementComponent = () => {
    // Theme colors
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

    // State for coupons and loading
    const [coupons, setCoupons] = useState([]);
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
    const [selectedCoupon, setSelectedCoupon] = useState(null);

    // Unified filter state
    const [filters, setFilters] = useState({
        discount_type: '',
        start_date: '',
        end_date: '',
        is_active: '',
    });

    // State for active filter parameters that will be sent to the API
    const [activeFilterParams, setActiveFilterParams] = useState({});

    // Search state
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
            params[searchParams.searchBy] = searchParams.searchValue;
        }

        return params;
    }, [currentPage, rowsPerPage, activeFilterParams, searchParams]);

    // Fetch coupons from API
    const fetchCoupons = useCallback(async () => {
        setIsLoading(true);
        setIsError(false);

        try {
            const params = createRequestParams();
            const response = await couponService.getCoupons(params);

            if (response && response.data) {
                setCoupons(response.data || []);

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
                setCoupons([]);
                setTotalCount(0);
                setTotalPages(1);
            }
        } catch (error) {
            setIsError(true);
            setErrorMessage(error.response?.data?.error?.message || 'Failed to fetch coupons');

            toast({
                title: 'Error loading coupons',
                description: error.response?.data?.error?.message || 'An error occurred while loading coupons',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });

            setCoupons([]);
            setTotalCount(0);
            setTotalPages(1);
        } finally {
            setIsLoading(false);
        }
    }, [createRequestParams, rowsPerPage, toast]);

    // Handle search from CouponSearchComponent
    const handleSearch = (searchField, searchValue) => {
        setSearchParams({
            searchBy: searchField,
            searchValue: searchValue
        });
    };

    // Handle filter changes
    const handleFiltersChange = (updatedFilters) => {
        setFilters(updatedFilters);
    };

    // Apply filters from modal
    const handleApplyFilters = (filteredParams) => {
        setActiveFilterParams(filteredParams);
        setCurrentPage(1);
    };

    // Load coupons when component mounts or when dependencies change
    useEffect(() => {
        if (searchParams.searchBy || searchParams.searchValue) {
            setCurrentPage(1);
        }
        fetchCoupons();
    }, [fetchCoupons, searchParams]);

    // Handle coupon created and reload current page
    const handleCouponCreated = async () => {
        await fetchCoupons();
    };

    // Handle coupon updated and reload current page
    const handleCouponUpdated = async () => {
        await fetchCoupons();
    };

    // Handle coupon deleted and reload current page
    const handleDeleteCoupon = async (couponId) => {
        if (window.confirm('Bạn có chắc chắn muốn xoá mã khuyến mãi này?')) {
            try {
                await couponService.deleteCoupon(couponId);
                toast({
                    title: "Xoá mã khuyến mãi thành công",
                    status: "success",
                    duration: 3000,
                    isClosable: true,
                });
                await fetchCoupons();
            } catch (error) {
                toast({
                    title: "Xoá mã khuyến mãi thất bại",
                    description: error.response?.data?.error?.message || 'An unexpected error occurred',
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            }
        }
    };

    // Handle opening Edit Modal
    const handleOpenEditModal = async (coupon) => {
        try {
            // Fetch detailed coupon data
            const response = await couponService.getCouponById(coupon.id);
            setSelectedCoupon(response.data);
            onOpenEditModal();
        } catch (error) {
            toast({
                title: "Lỗi khi tải thông tin mã khuyến mãi",
                description: error.response?.data?.error?.message || 'An unexpected error occurred',
                status: "error",
                duration: 5000,
                isClosable: true,
            });
        }
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
            return date.toLocaleDateString('vi-VN');
        } catch (e) {
            return dateString;
        }
    };

    // Format currency
    const formatCurrency = (value) => {
        return new Intl.NumberFormat('vi-VN', {
            style: 'currency',
            currency: 'VND'
        }).format(value);
    };

    // Get color for discount type badge
    const getDiscountTypeColor = (discountType) => {
        switch (discountType) {
            case 'percentage':
                return 'purple';
            case 'fixed_amount':
                return 'blue';
            default:
                return 'gray';
        }
    };

    // Get status color
    const getStatusColor = (isActive, endDate) => {
        if (!isActive) return 'red';

        if (endDate) {
            const now = new Date();
            const end = new Date(endDate);
            if (end < now) return 'orange'; // Expired
        }

        return 'green'; // Active
    };

    // Get status text
    const getStatusText = (isActive, startDate, endDate) => {
        if (!isActive) return 'Không hoạt động';

        const now = new Date();

        if (startDate) {
            const start = new Date(startDate);
            if (start > now) return 'Chưa bắt đầu';
        }

        if (endDate) {
            const end = new Date(endDate);
            if (end < now) return 'Đã hết hạn';
        }

        return 'Đang hoạt động';
    };

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
                    <CouponSearchComponent
                        onSearch={handleSearch}
                        isLoading={isLoading}
                    />
                    <Tooltip label="Refresh data" hasArrow>
                        <IconButton
                            icon={<FiRefreshCw />}
                            onClick={fetchCoupons}
                            aria-label="Refresh data"
                            variant="ghost"
                            colorScheme="blue"
                            size="sm"
                            isLoading={isLoading}
                        />
                    </Tooltip>
                </Flex>

                <HStack spacing={2}>
                    <CouponFilterDropdown
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
                        Tạo mã khuyến mãi
                    </Button>
                </HStack>
            </Flex>

            {isError && (
                <Alert status="error" variant="left-accent" mb={4} borderRadius="md">
                    <AlertIcon />
                    <Text>{errorMessage || 'An error occurred while fetching coupons'}</Text>
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
                                    Mã khuyến mãi
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
                                    Loại & Giá trị
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={tableTextColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Điều kiện
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={tableTextColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Thời gian
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={tableTextColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Sử dụng
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={tableTextColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Trạng thái
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
                                    Hành động
                                </Th>
                            </Tr>
                        </Thead>
                        <Tbody>
                            {isLoading ? (
                                <Tr>
                                    <Td colSpan={7} textAlign="center" py={12}>
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
                                            <Text color="gray.500" fontSize="sm">Loading coupons...</Text>
                                        </Flex>
                                    </Td>
                                </Tr>
                            ) : coupons.length > 0 ? (
                                coupons.map((coupon) => (
                                    <Tr
                                        key={coupon.id}
                                        _hover={{ bg: tableRowHoverBg }}
                                        transition="background-color 0.2s"
                                        borderBottomWidth="1px"
                                        borderColor={borderColor}
                                        _active={{ bg: tableRowActiveBg }}
                                        h="60px"
                                        bg={bgColor}
                                    >
                                        <Td>
                                            <VStack align="start" spacing={1}>
                                                <HStack spacing={2}>
                                                    <Code colorScheme="blue" fontSize="sm" fontWeight="bold">
                                                        {coupon.code}
                                                    </Code>
                                                </HStack>
                                                <Text
                                                    fontSize="sm"
                                                    color={textColorPrimary}
                                                    fontWeight="medium"
                                                    noOfLines={1}
                                                >
                                                    {coupon.name}
                                                </Text>
                                                {coupon.description && (
                                                    <Text
                                                        fontSize="xs"
                                                        color={textColorSecondary}
                                                        noOfLines={1}
                                                        maxW="200px"
                                                    >
                                                        {coupon.description}
                                                    </Text>
                                                )}
                                            </VStack>
                                        </Td>
                                        <Td display={{ base: "none", md: "table-cell" }}>
                                            <VStack align="start" spacing={1}>
                                                <HStack spacing={2}>
                                                    <Badge
                                                        colorScheme={getDiscountTypeColor(coupon.discount_type)}
                                                        variant="subtle"
                                                        borderRadius="md"
                                                        px={2}
                                                        py={1}
                                                        fontSize="xs"
                                                        fontWeight="bold"
                                                        display="flex"
                                                        alignItems="center"
                                                        gap={1}
                                                    >
                                                        {coupon.discount_type === 'percentage' ? (
                                                            <>
                                                                <FiPercent size={10} />
                                                                {coupon.discount_value}%
                                                            </>
                                                        ) : (
                                                            <>
                                                                <FiDollarSign size={10} />
                                                                {formatCurrency(coupon.discount_value)}
                                                            </>
                                                        )}
                                                    </Badge>
                                                </HStack>
                                                <Text fontSize="xs" color={textColorSecondary}>
                                                    Tối đa: {formatCurrency(coupon.maximum_discount_amount)}
                                                </Text>
                                            </VStack>
                                        </Td>
                                        <Td>
                                            <VStack align="start" spacing={1}>
                                                <Text fontSize="sm" color={textColorTertiary}>
                                                    Đơn tối thiểu:
                                                </Text>
                                                <Text fontSize="sm" color={textColorPrimary} fontWeight="medium">
                                                    {formatCurrency(coupon.minimum_order_amount)}
                                                </Text>
                                            </VStack>
                                        </Td>
                                        <Td>
                                            <VStack align="start" spacing={1}>
                                                <HStack spacing={1} align="center">
                                                    <FiCalendar size={12} />
                                                    <Text fontSize="xs" color={textColorSecondary}>
                                                        {formatDate(coupon.start_date)}
                                                    </Text>
                                                </HStack>
                                                <HStack spacing={1} align="center">
                                                    <FiCalendar size={12} />
                                                    <Text fontSize="xs" color={textColorSecondary}>
                                                        {formatDate(coupon.end_date)}
                                                    </Text>
                                                </HStack>
                                            </VStack>
                                        </Td>
                                        <Td>
                                            <VStack align="start" spacing={1}>
                                                <HStack spacing={1} align="center">
                                                    <FiUsers size={12} />
                                                    <Text fontSize="sm" color={textColorPrimary} fontWeight="medium">
                                                        {coupon.usage_count}/{coupon.usage_limit}
                                                    </Text>
                                                </HStack>
                                                <Box
                                                    w="60px"
                                                    h="4px"
                                                    bg="gray.200"
                                                    borderRadius="full"
                                                    overflow="hidden"
                                                >
                                                    <Box
                                                        h="100%"
                                                        bg={coupon.usage_count >= coupon.usage_limit ? "red.400" : "blue.400"}
                                                        w={`${Math.min((coupon.usage_count / coupon.usage_limit) * 100, 100)}%`}
                                                        borderRadius="full"
                                                        transition="width 0.3s ease"
                                                    />
                                                </Box>
                                            </VStack>
                                        </Td>
                                        <Td>
                                            <Badge
                                                px={2}
                                                py={1}
                                                borderRadius="full"
                                                colorScheme={getStatusColor(coupon.is_active, coupon.end_date)}
                                                textTransform="capitalize"
                                                fontWeight="medium"
                                                fontSize="xs"
                                            >
                                                {getStatusText(coupon.is_active, coupon.start_date, coupon.end_date)}
                                            </Badge>
                                        </Td>
                                        <Td textAlign="right">
                                            <HStack spacing={1} justifyContent="flex-end">
                                                <Tooltip label="Edit coupon" hasArrow>
                                                    <IconButton
                                                        icon={<FiEdit2 size={15} />}
                                                        size="sm"
                                                        variant="ghost"
                                                        colorScheme="blue"
                                                        aria-label="Edit coupon"
                                                        borderRadius="md"
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            handleOpenEditModal(coupon);
                                                        }}
                                                    />
                                                </Tooltip>
                                                <Tooltip label="Delete coupon" hasArrow>
                                                    <IconButton
                                                        icon={<FiTrash2 size={15} />}
                                                        size="sm"
                                                        variant="ghost"
                                                        colorScheme="red"
                                                        aria-label="Delete coupon"
                                                        borderRadius="md"
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            handleDeleteCoupon(coupon.id);
                                                        }}
                                                    />
                                                </Tooltip>
                                            </HStack>
                                        </Td>
                                    </Tr>
                                ))
                            ) : (
                                <Tr>
                                    <Td colSpan={7} textAlign="center" py={12}>
                                        <Flex direction="column" align="center" justify="center" py={8}>
                                            <Box color="gray.400" mb={3}>
                                                <FiTag size={36} />
                                            </Box>
                                            <Text fontWeight="normal" color="gray.500" fontSize="md">No coupons found</Text>
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
                                Hiển thị {totalCount > 0 ? (currentPage - 1) * rowsPerPage + 1 : 0}-{Math.min(currentPage * rowsPerPage, totalCount)} trên tổng số {totalCount} mã khuyến mãi
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
                                    {rowsPerPage} dòng mỗi trang
                                </MenuButton>
                                <MenuList minW="120px" shadow="lg" borderRadius="md">
                                    <MenuItem onClick={() => setRowsPerPage(10)}>10 dòng mỗi trang</MenuItem>
                                    <MenuItem onClick={() => setRowsPerPage(15)}>15 dòng mỗi trang</MenuItem>
                                    <MenuItem onClick={() => setRowsPerPage(20)}>20 dòng mỗi trang</MenuItem>
                                    <MenuItem onClick={() => setRowsPerPage(50)}>50 dòng mỗi trang</MenuItem>
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

            <CreateCouponModal
                isOpen={isCreateModalOpen}
                onClose={onCloseCreateModal}
                onCouponCreated={handleCouponCreated}
            />

            {selectedCoupon && (
                <EditCouponModal
                    isOpen={isEditModalOpen}
                    onClose={onCloseEditModal}
                    coupon={selectedCoupon}
                    onCouponUpdated={handleCouponUpdated}
                />
            )}
        </Container>
    );
};

export default CouponManagementComponent;