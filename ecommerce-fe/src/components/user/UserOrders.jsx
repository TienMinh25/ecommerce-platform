import React, { useState, useEffect } from 'react';
import {
    Box,
    Heading,
    Text,
    Flex,
    Tabs,
    TabList,
    Tab,
    TabPanels,
    TabPanel,
    Input,
    InputGroup,
    InputLeftElement,
    Badge,
    Button,
    Image,
    VStack,
    HStack,
    Divider,
    useColorModeValue,
    Link,
    Collapse,
    Spinner,
    Alert,
    AlertIcon,
    Grid,
    GridItem,
    useDisclosure,
    Avatar
} from '@chakra-ui/react';
import { SearchIcon, InfoOutlineIcon, ChevronDownIcon, ChevronUpIcon } from '@chakra-ui/icons';
import userMeService from '../../services/userMeService';
import OrderDetailsSection from "./OrderDetailsSection.jsx";

const UserOrders = () => {
    const [searchQuery, setSearchQuery] = useState('');
    const [activeTab, setActiveTab] = useState(0);
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 10,
        total_items: 0,
        total_pages: 0,
        has_next: false,
        has_previous: false
    });
    const [expandedOrder, setExpandedOrder] = useState(null);

    // Status mapping for tabs
    const statusTabs = [
        { key: null, label: 'Tất cả' },
        { key: 'pending_payment', label: 'Chờ thanh toán' },
        { key: 'pending', label: 'Chờ xác nhận' },
        { key: 'confirmed', label: 'Đã xác nhận' },
        { key: 'processing', label: 'Đang chuẩn bị' },
        { key: 'ready_to_ship', label: 'Sẵn sàng giao' },
        { key: 'in_transit', label: 'Đang vận chuyển' },
        { key: 'out_for_delivery', label: 'Sắp giao' },
        { key: 'delivered', label: 'Đã giao' },
        { key: 'cancelled', label: 'Đã hủy' },
        { key: 'payment_failed', label: 'Thanh toán thất bại' },
        { key: 'refunded', label: 'Đã hoàn tiền' }
    ];

    // Colors
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const redColor = useColorModeValue('red.500', 'red.300');
    const greenColor = useColorModeValue('green.500', 'green.300');

    // Fetch orders
    const fetchOrders = async (page = 1, status = null, keyword = null) => {
        setLoading(true);
        setError(null);

        try {
            const params = {
                page,
                limit: pagination.limit
            };

            if (status) {
                params.status = status;
            }

            if (keyword && keyword.trim()) {
                params.keyword = keyword.trim();
            }

            const response = await userMeService.getOrders(params);
            setOrders(response.data || []);
            setPagination(response.metadata.pagination);
        } catch (err) {
            setError(err.message || 'Có lỗi xảy ra khi tải đơn hàng');
            setOrders([]);
        } finally {
            setLoading(false);
        }
    };

    // Effect to fetch orders when tab or search changes
    useEffect(() => {
        const currentStatus = statusTabs[activeTab]?.key;
        fetchOrders(1, currentStatus, searchQuery);
    }, [activeTab]);

    // Handle search
    const handleSearch = () => {
        const currentStatus = statusTabs[activeTab]?.key;
        fetchOrders(1, currentStatus, searchQuery);
    };

    // Handle search input change with debounce
    useEffect(() => {
        const timeoutId = setTimeout(() => {
            if (searchQuery !== '') {
                handleSearch();
            } else {
                const currentStatus = statusTabs[activeTab]?.key;
                fetchOrders(1, currentStatus, null);
            }
        }, 500);

        return () => clearTimeout(timeoutId);
    }, [searchQuery]);

    // Status badges
    const renderStatusBadge = (status) => {
        const statusConfig = {
            'pending_payment': {
                text: 'Chờ thanh toán',
                bg: '#FEF3C7',
                color: '#92400E'
            },
            'pending': {
                text: 'Chờ xác nhận',
                bg: '#FED7AA',
                color: '#C2410C'
            },
            'confirmed': {
                text: 'Đã xác nhận',
                bg: '#DBEAFE',
                color: '#1E40AF'
            },
            'processing': {
                text: 'Đang chuẩn bị',
                bg: '#E9D5FF',
                color: '#6B21A8'
            },
            'ready_to_ship': {
                text: 'Sẵn sàng giao',
                bg: '#CFFAFE',
                color: '#155E75'
            },
            'in_transit': {
                text: 'Đang vận chuyển',
                bg: '#DBEAFE',
                color: '#1E40AF'
            },
            'out_for_delivery': {
                text: 'Sắp giao',
                bg: '#CCFBF1',
                color: '#134E4A'
            },
            'delivered': {
                text: 'Đã giao',
                bg: '#D1FAE5',
                color: '#065F46'
            },
            'cancelled': {
                text: 'Đã hủy',
                bg: '#FEE2E2',
                color: '#991B1B'
            },
            'payment_failed': {
                text: 'Thanh toán thất bại',
                bg: '#FEE2E2',
                color: '#991B1B'
            },
            'refunded': {
                text: 'Đã hoàn tiền',
                bg: '#F3F4F6',
                color: '#374151'
            }
        };

        const config = statusConfig[status] || {
            text: 'Không xác định',
            bg: '#F3F4F6',
            color: '#374151'
        };

        return (
            <Box
                px={3}
                py={1.5}
                borderRadius="md"
                bg={config.bg}
            >
                <Text
                    color={config.color}
                    fontWeight="medium"
                    fontSize="sm"
                >
                    {config.text}
                </Text>
            </Box>
        );
    };

    // Format currency
    const formatCurrency = (amount) => {
        return `₫${amount.toLocaleString('vi-VN')}`;
    };

    // Convert orders to individual order items (no grouping)
    const convertToOrderItems = (orders) => {
        return orders.map(order => {
            const calculatedAmount = order.total_price + order.tax_amount + order.shipping_fee - order.discount_amount;
            const isPaidAlready = order.shipping_method === 'momo' && order.status !== 'payment_failed';

            return {
                ...order,
                supplier_name: order.supplier_name || 'Cửa hàng không xác định',
                finalAmount: isPaidAlready ? 0 : calculatedAmount
            }
        });
    };

    // Toggle order details
    const toggleOrderDetails = (orderItemId) => {
        setExpandedOrder(expandedOrder === orderItemId ? null : orderItemId);
    };

    // Render single order item
    const renderOrderItem = (order) => {
        const isExpanded = expandedOrder === order.order_item_id;

        return (
            <Box
                key={order.order_item_id}
                mb={6}
                bg={bgColor}
                borderRadius="md"
                borderWidth="1px"
                borderColor={borderColor}
                overflow="hidden"
            >
                {/* Order header */}
                <Flex
                    bg="gray.50"
                    p={4}
                    justify="space-between"
                    align="center"
                    borderBottomWidth="1px"
                    borderColor={borderColor}
                    direction={{ base: "column", md: "row" }}
                    gap={{ base: 3, md: 0 }}
                >
                    <Flex align="center" gap={3} flex="1" minW="0">
                        <Avatar
                            size="sm"
                            src={order.supplier_thumbnail}
                            name={order.supplier_name}
                        />
                        <Text fontWeight="bold" noOfLines={1} flex="1">
                            {order.supplier_name}
                        </Text>
                    </Flex>

                    <Flex align="center" gap={2} flexShrink={0} direction={{ base: "row", md: "row" }}>
                        {order.status === 'delivered' && (
                            <Flex align="center" mr={2} color={greenColor} display={{ base: "none", lg: "flex" }}>
                                <InfoOutlineIcon mr={1} />
                                <Text fontSize="sm">Đơn hàng đã được giao thành công</Text>
                            </Flex>
                        )}
                        {renderStatusBadge(order.status)}
                    </Flex>

                    {/* Mobile buttons */}
                    <Flex gap={2} display={{ base: "flex", md: "none" }} w="full" justify="center">
                        <Button size="xs" colorScheme="red" variant="outline">
                            Chat
                        </Button>
                        <Button size="xs" variant="outline">
                            Xem Shop
                        </Button>
                    </Flex>
                </Flex>

                {/* Order item */}
                <Flex
                    p={4}
                    borderBottomWidth="1px"
                    borderColor={borderColor}
                    align="center"
                >
                    <Flex flex="1" gap={4}>
                        <Image
                            src={order.product_variant_thumbnail || 'https://via.placeholder.com/80'}
                            alt={order.product_name}
                            boxSize="80px"
                            objectFit="cover"
                            borderRadius="md"
                            border="1px solid"
                            borderColor={borderColor}
                        />

                        <Flex direction="column" flex="1">
                            <Text noOfLines={2} fontWeight="medium" mb={1}>
                                {order.product_name}
                            </Text>
                            <Text fontSize="sm" color="gray.600" mb={1}>
                                Phân loại hàng: {order.product_variant_name}
                            </Text>
                            <Text fontSize="sm">x{order.quantity}</Text>
                        </Flex>
                    </Flex>

                    <Flex direction="column" align="flex-end" minW="140px">
                        {order.discount_amount > 0 && (
                            <Text textDecoration="line-through" color="gray.500" fontSize="sm">
                                {formatCurrency(order.unit_price)}
                            </Text>
                        )}
                        <Text color={redColor} fontWeight="bold">
                            {formatCurrency(order.total_price)}
                        </Text>
                    </Flex>
                </Flex>

                {/* Order footer */}
                <Flex
                    p={4}
                    justify="space-between"
                    align="center"
                    bg="gray.50"
                >
                    <Button
                        size="sm"
                        variant="ghost"
                        leftIcon={isExpanded ? <ChevronUpIcon /> : <ChevronDownIcon />}
                        onClick={() => toggleOrderDetails(order.order_item_id)}
                    >
                        {isExpanded ? 'Ẩn chi tiết' : 'Xem chi tiết'}
                    </Button>

                    <Flex align="center" gap={3}>
                        <Text fontSize="sm">Thành tiền:</Text>
                        <Text fontSize="lg" fontWeight="bold" color={redColor}>
                            {formatCurrency(order.finalAmount)}
                        </Text>

                        {/* Order actions */}
                        {order.status === 'cancelled' && (
                            <>
                                <Button size="sm" colorScheme="red">
                                    Mua Lại
                                </Button>
                                <Button size="sm" variant="outline">
                                    Liên Hệ Người Bán
                                </Button>
                            </>
                        )}

                        {order.status === 'delivered' && (
                            <>
                                <Button size="sm" colorScheme="red">
                                    Mua Lại
                                </Button>
                                <Button size="sm" variant="outline">
                                    Đánh Giá
                                </Button>
                                <Button size="sm" variant="outline">
                                    Liên Hệ Người Bán
                                </Button>
                            </>
                        )}
                    </Flex>
                </Flex>

                {/* Expandable order details */}
                <Collapse in={isExpanded} animateOpacity>
                    <OrderDetailsSection order={order} />
                </Collapse>
            </Box>
        );
    };

    // Handle pagination
    const handlePageChange = (newPage) => {
        const currentStatus = statusTabs[activeTab]?.key;
        fetchOrders(newPage, currentStatus, searchQuery);
    };

    return (
        <Box>
            <Heading as="h1" size="lg" mb={6}>
                Đơn Hàng Của Tôi
            </Heading>

            <Box
                bg={bgColor}
                borderRadius="md"
                borderWidth="1px"
                borderColor={borderColor}
                overflow="hidden"
            >
                <Tabs
                    variant="line"
                    colorScheme="red"
                    onChange={index => setActiveTab(index)}
                >
                    <TabList
                        overflowX="auto"
                        overflowY="hidden"
                        maxW="100%"
                        borderBottomWidth="2px"
                        borderColor={borderColor}
                        css={{
                            '&::-webkit-scrollbar': {
                                display: 'none',
                            },
                            '-ms-overflow-style': 'none',
                            'scrollbar-width': 'none',
                        }}
                    >
                        {statusTabs.map((tab, index) => (
                            <Tab
                                key={index}
                                _selected={{
                                    color: redColor,
                                    borderColor: redColor,
                                    borderBottomWidth: "3px"
                                }}
                                whiteSpace="nowrap"
                                minW="fit-content"
                                px={4}
                                py={3}
                                fontSize="sm"
                                fontWeight="medium"
                                position="relative"
                                borderBottomWidth="3px"
                                borderColor="transparent"
                                transition="all 0.2s"
                                _hover={{
                                    bg: "gray.50",
                                    color: redColor
                                }}
                                mr={1}
                            >
                                {tab.label}
                            </Tab>
                        ))}
                    </TabList>

                    <TabPanels>
                        {statusTabs.map((tab, tabIndex) => (
                            <TabPanel key={tabIndex} p={4}>
                                <Box mb={6}>
                                    <InputGroup>
                                        <InputLeftElement pointerEvents="none">
                                            <SearchIcon color="gray.300" />
                                        </InputLeftElement>
                                        <Input
                                            placeholder="Bạn có thể tìm kiếm theo tên sản phẩm"
                                            value={searchQuery}
                                            onChange={e => setSearchQuery(e.target.value)}
                                            bg={bgColor}
                                            borderColor={borderColor}
                                        />
                                    </InputGroup>
                                </Box>

                                {error && (
                                    <Alert status="error" mb={4}>
                                        <AlertIcon />
                                        {error}
                                    </Alert>
                                )}

                                {loading ? (
                                    <Flex justify="center" py={10}>
                                        <Spinner size="lg" color={redColor} />
                                    </Flex>
                                ) : (
                                    <>
                                        {orders.length > 0 ? (
                                            <>
                                                {convertToOrderItems(orders).map(order =>
                                                    renderOrderItem(order)
                                                )}

                                                {/* Pagination */}
                                                {pagination.total_pages > 1 && (
                                                    <Flex justify="center" mt={6} gap={2}>
                                                        <Button
                                                            size="sm"
                                                            onClick={() => handlePageChange(pagination.page - 1)}
                                                            isDisabled={!pagination.has_previous}
                                                        >
                                                            Trước
                                                        </Button>

                                                        <Text mx={4} alignSelf="center">
                                                            Trang {pagination.page} / {pagination.total_pages}
                                                        </Text>

                                                        <Button
                                                            size="sm"
                                                            onClick={() => handlePageChange(pagination.page + 1)}
                                                            isDisabled={!pagination.has_next}
                                                        >
                                                            Sau
                                                        </Button>
                                                    </Flex>
                                                )}
                                            </>
                                        ) : (
                                            <VStack spacing={4} py={10}>
                                                <Image
                                                    src="https://deo.shopeemobile.com/shopee/shopee-pcmall-live-sg/assets/5fafbb923393b712b96488590b8f781f.png"
                                                    alt="No orders"
                                                    boxSize="100px"
                                                />
                                                <Text color="gray.500">Chưa có đơn hàng</Text>
                                            </VStack>
                                        )}
                                    </>
                                )}
                            </TabPanel>
                        ))}
                    </TabPanels>
                </Tabs>
            </Box>
        </Box>
    );
};

export default UserOrders;