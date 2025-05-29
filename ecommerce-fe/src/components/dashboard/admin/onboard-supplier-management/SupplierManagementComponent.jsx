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
    Image,
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
    useToast,
    VStack,
} from '@chakra-ui/react';
import {
    FiChevronDown,
    FiChevronLeft,
    FiChevronRight,
    FiEye,
    FiRefreshCw,
    FiSearch,
    FiPause,
    FiPlay,
    FiShield,
} from 'react-icons/fi';
import supplierService from '../../../../services/supplierService.js';
import SupplierFilterDropdown from './SupplierFilterDropdown.jsx';
import SupplierSearchComponent from './SupplierSearchComponent.jsx';
import SupplierDetailModal from './SupplierDetailModal.jsx';

const SupplierManagementComponent = () => {
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

    const toast = useToast();

    // State
    const [suppliers, setSuppliers] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');

    // Detail modal state
    const [selectedSupplierId, setSelectedSupplierId] = useState(null);
    const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);

    // Pagination
    const [totalCount, setTotalCount] = useState(0);
    const [currentPage, setCurrentPage] = useState(1);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [totalPages, setTotalPages] = useState(1);

    // Filters
    const [filters, setFilters] = useState({
        status: '',
    });
    const [activeFilterParams, setActiveFilterParams] = useState({});

    // Search
    const [searchParams, setSearchParams] = useState({
        searchBy: '',
        searchValue: ''
    });

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

    const fetchSuppliers = useCallback(async () => {
        setIsLoading(true);
        setIsError(false);

        try {
            const params = createRequestParams();
            const response = await supplierService.getSuppliers(params);

            if (response && response.data) {
                setSuppliers(response.data || []);

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
                setSuppliers([]);
                setTotalCount(0);
                setTotalPages(1);
            }
        } catch (error) {
            setIsError(true);
            setErrorMessage(error.response?.data?.error?.message || 'Không thể tải danh sách nhà cung cấp');

            toast({
                title: 'Lỗi khi tải danh sách nhà cung cấp',
                description: error.response?.data?.error?.message || 'Đã xảy ra lỗi khi tải danh sách nhà cung cấp',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });

            setSuppliers([]);
            setTotalCount(0);
            setTotalPages(1);
        } finally {
            setIsLoading(false);
        }
    }, [createRequestParams, rowsPerPage, toast]);

    const handleSearch = (searchField, searchValue) => {
        setSearchParams({
            searchBy: searchField,
            searchValue: searchValue
        });
    };

    const handleFiltersChange = (updatedFilters) => {
        setFilters(updatedFilters);
    };

    const handleApplyFilters = (filteredParams) => {
        setActiveFilterParams(filteredParams);
        setCurrentPage(1);
    };

    const handleSupplierClick = (supplier) => {
        setSelectedSupplierId(supplier.id);
        setIsDetailModalOpen(true);
    };

    const handleCloseDetailModal = () => {
        setIsDetailModalOpen(false);
        setSelectedSupplierId(null);
    };

    const handleUpdateSupplierStatus = async (supplierId, currentStatus) => {
        const newStatus = currentStatus === 'active' ? 'suspended' : 'active';

        try {
            await supplierService.updateSupplier(supplierId, { status: newStatus });
            toast({
                title: 'Cập nhật trạng thái thành công',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
            await fetchSuppliers();
        } catch (error) {
            toast({
                title: 'Cập nhật trạng thái thất bại',
                description: error.response?.data?.error?.message || 'Có lỗi xảy ra khi cập nhật trạng thái',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    useEffect(() => {
        if (searchParams.searchBy || searchParams.searchValue) {
            setCurrentPage(1);
        }
        fetchSuppliers();
    }, [fetchSuppliers, searchParams]);

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

    const paginationRange = generatePaginationRange(currentPage, totalPages);

    const formatDate = (dateString) => {
        if (!dateString) return 'N/A';
        try {
            const date = new Date(dateString);
            return date.toLocaleString();
        } catch (e) {
            return dateString;
        }
    };

    const getStatusColor = (status) => {
        switch (status) {
            case 'active':
                return 'green';
            case 'pending':
                return 'yellow';
            case 'suspended':
                return 'red';
            default:
                return 'gray';
        }
    };

    const getStatusText = (status) => {
        switch (status) {
            case 'active':
                return 'Hoạt động';
            case 'pending':
                return 'Chờ duyệt';
            case 'suspended':
                return 'Tạm ngưng';
            default:
                return status;
        }
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
                    <SupplierSearchComponent
                        onSearch={handleSearch}
                        isLoading={isLoading}
                    />
                    <Tooltip label="Làm mới dữ liệu" hasArrow>
                        <IconButton
                            icon={<FiRefreshCw />}
                            onClick={fetchSuppliers}
                            aria-label="Làm mới dữ liệu"
                            variant="ghost"
                            colorScheme="blue"
                            size="sm"
                            isLoading={isLoading}
                        />
                    </Tooltip>
                </Flex>

                <HStack spacing={2}>
                    <SupplierFilterDropdown
                        filters={filters}
                        onFiltersChange={handleFiltersChange}
                        onApplyFilters={handleApplyFilters}
                    />
                </HStack>
            </Flex>

            {isError && (
                <Alert status="error" variant="left-accent" mb={4} borderRadius="md">
                    <AlertIcon />
                    <Text>{errorMessage || 'Đã xảy ra lỗi khi tải danh sách nhà cung cấp'}</Text>
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
                    <Table variant="simple" size="md" colorScheme="gray">
                        <Thead bg={tableBgColor} position="sticky" top={0} zIndex={1}>
                            <Tr>
                                <Th py={4} fontSize="xs" color={tableTextColor} letterSpacing="0.5px" textTransform="uppercase" fontWeight="bold">
                                    Nhà cung cấp
                                </Th>
                                <Th py={4} display={{ base: "none", md: "table-cell" }} fontSize="xs" color={tableTextColor} letterSpacing="0.5px" textTransform="uppercase" fontWeight="bold">
                                    Thông tin liên hệ
                                </Th>
                                <Th py={4} fontSize="xs" color={tableTextColor} letterSpacing="0.5px" textTransform="uppercase" fontWeight="bold">
                                    Mã số thuế
                                </Th>
                                <Th py={4} fontSize="xs" color={tableTextColor} letterSpacing="0.5px" textTransform="uppercase" fontWeight="bold">
                                    Trạng thái
                                </Th>
                                <Th py={4} fontSize="xs" color={tableTextColor} letterSpacing="0.5px" textTransform="uppercase" fontWeight="bold">
                                    Thời gian cập nhật
                                </Th>
                                <Th py={4} textAlign="right" fontSize="xs" color={tableTextColor} letterSpacing="0.5px" textTransform="uppercase" fontWeight="bold">
                                    Hành động
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
                                            <Text color="gray.500" fontSize="sm">Đang tải danh sách nhà cung cấp...</Text>
                                        </Flex>
                                    </Td>
                                </Tr>
                            ) : suppliers.length > 0 ? (
                                suppliers.map((supplier) => (
                                    <Tr
                                        key={supplier.id}
                                        _hover={{ bg: tableRowHoverBg }}
                                        transition="background-color 0.2s"
                                        borderBottomWidth="1px"
                                        borderColor={borderColor}
                                        _active={{ bg: tableRowActiveBg }}
                                        h="60px"
                                        bg={bgColor}
                                        cursor="pointer"
                                        onClick={() => handleSupplierClick(supplier)}
                                    >
                                        <Td>
                                            <HStack spacing={3}>
                                                <Image
                                                    src={supplier.logo_thumbnail_url || "/api/placeholder/40/40"}
                                                    alt={supplier.company_name}
                                                    boxSize="40px"
                                                    borderRadius="md"
                                                    objectFit="cover"
                                                    fallback={<Avatar size="sm" name={supplier.company_name} />}
                                                />
                                                <Box>
                                                    <Text fontWeight="medium" fontSize="sm" color={textColorPrimary}>
                                                        {supplier.company_name}
                                                    </Text>
                                                </Box>
                                            </HStack>
                                        </Td>
                                        <Td display={{ base: "none", md: "table-cell" }}>
                                            <VStack align="start" spacing={0.5}>
                                                <Text fontSize="sm" color={textColorTertiary}>
                                                    {supplier.contact_phone}
                                                </Text>
                                                <Text fontSize="xs" color={textColorSecondary}>
                                                    {supplier.business_address}
                                                </Text>
                                            </VStack>
                                        </Td>
                                        <Td>
                                            <Text fontSize="sm" color={textColorPrimary} fontFamily="mono">
                                                {supplier.tax_id}
                                            </Text>
                                        </Td>
                                        <Td>
                                            <Badge
                                                px={2}
                                                py={1}
                                                borderRadius="full"
                                                colorScheme={getStatusColor(supplier.status)}
                                                textTransform="capitalize"
                                                fontWeight="medium"
                                                fontSize="xs"
                                            >
                                                {getStatusText(supplier.status)}
                                            </Badge>
                                        </Td>
                                        <Td>
                                            <Text fontSize="sm" color={textColorSecondary}>
                                                {formatDate(supplier.updated_at)}
                                            </Text>
                                        </Td>
                                        <Td textAlign="right">
                                            <HStack spacing={1} justifyContent="flex-end">
                                                <Tooltip label="Xem chi tiết" hasArrow>
                                                    <IconButton
                                                        icon={<FiEye size={15} />}
                                                        size="sm"
                                                        variant="ghost"
                                                        colorScheme="blue"
                                                        aria-label="Xem chi tiết"
                                                        borderRadius="md"
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            handleSupplierClick(supplier);
                                                        }}
                                                    />
                                                </Tooltip>
                                                {supplier.status !== 'pending' && (
                                                    <Tooltip
                                                        label={
                                                            supplier.status === 'active'
                                                                ? 'Tạm ngưng hoạt động'
                                                                : 'Kích hoạt hoạt động'
                                                        }
                                                        hasArrow
                                                    >
                                                        <IconButton
                                                            icon={
                                                                supplier.status === 'active'
                                                                    ? <FiPause size={15} />
                                                                    : <FiPlay size={15} />
                                                            }
                                                            size="sm"
                                                            variant="ghost"
                                                            colorScheme={supplier.status === 'active' ? 'orange' : 'green'}
                                                            aria-label={
                                                                supplier.status === 'active'
                                                                    ? 'Tạm ngưng nhà cung cấp'
                                                                    : 'Kích hoạt nhà cung cấp'
                                                            }
                                                            borderRadius="md"
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                handleUpdateSupplierStatus(supplier.id, supplier.status);
                                                            }}
                                                            _hover={{
                                                                bg: supplier.status === 'active'
                                                                    ? 'orange.100'
                                                                    : 'green.100',
                                                                transform: 'scale(1.05)'
                                                            }}
                                                            transition="all 0.2s"
                                                        />
                                                    </Tooltip>
                                                )}
                                                {supplier.status === 'pending' && (
                                                    <Tooltip label="Đang chờ duyệt" hasArrow>
                                                        <IconButton
                                                            icon={<FiShield size={15} />}
                                                            size="sm"
                                                            variant="ghost"
                                                            colorScheme="yellow"
                                                            aria-label="Đang chờ duyệt"
                                                            borderRadius="md"
                                                            isDisabled
                                                            _disabled={{
                                                                opacity: 0.6,
                                                                cursor: 'not-allowed'
                                                            }}
                                                        />
                                                    </Tooltip>
                                                )}
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
                                            <Text fontWeight="normal" color="gray.500" fontSize="md">Không tìm thấy nhà cung cấp nào</Text>
                                            <Text color="gray.400" fontSize="sm" mt={1}>Thử sử dụng từ khóa tìm kiếm hoặc bộ lọc khác</Text>
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
                                Hiển thị {totalCount > 0 ? (currentPage - 1) * rowsPerPage + 1 : 0}-{Math.min(currentPage * rowsPerPage, totalCount)} trên tổng số {totalCount} nhà cung cấp
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
                                    aria-label="Trang trước"
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
                                    aria-label="Trang sau"
                                    borderRadius="md"
                                />
                            </HStack>
                        )}
                    </Flex>
                </Box>
            </Box>

            {/* Detail Modal */}
            <SupplierDetailModal
                isOpen={isDetailModalOpen}
                onClose={handleCloseDetailModal}
                supplierId={selectedSupplierId}
            />
        </Container>
    );
};

export default SupplierManagementComponent;