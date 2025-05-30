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
    Button,
    Image,
    VStack,
    HStack,
    Divider,
    useColorModeValue,
    Collapse,
    Spinner,
    Alert,
    AlertIcon,
    Avatar,
    Select,
    useDisclosure,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    Textarea,
} from '@chakra-ui/react';
import { SearchIcon, ChevronDownIcon, ChevronUpIcon } from '@chakra-ui/icons';

const SupplierOrders = () => {
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
    const [selectedOrder, setSelectedOrder] = useState(null);
    const [statusToUpdate, setStatusToUpdate] = useState('');
    const [cancelReason, setCancelReason] = useState('');

    const { isOpen: isStatusModalOpen, onOpen: onStatusModalOpen, onClose: onStatusModalClose } = useDisclosure();
    const { isOpen: isCancelModalOpen, onOpen: onCancelModalOpen, onClose: onCancelModalClose } = useDisclosure();

    // Status tabs for supplier
    const statusTabs = [
        { key: null, label: 'Tất cả' },
        { key: 'pending', label: 'Chờ xác nhận' },
        { key: 'confirmed', label: 'Đã xác nhận' },
        { key: 'processing', label: 'Đang chuẩn bị' },
        { key: 'ready_to_ship', label: 'Sẵn sàng giao' },
        { key: 'cancelled', label: 'Đã hủy' },
    ];

    // Colors
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const redColor = useColorModeValue('red.500', 'red.300');
    const greenColor = useColorModeValue('green.500', 'green.300');

    // Mock data - replace with actual API call later
    const mockOrders = [
        {
            order_item_id: 1,
            order_id: "ORD001",
            product_name: "iPhone 14 Pro Max 256GB",
            product_variant_name: "Deep Purple",
            product_variant_thumbnail: "https://via.placeholder.com/80",
            quantity: 1,
            unit_price: 32990000,
            total_price: 32990000,
            discount_amount: 0,
            shipping_fee: 30000,
            tax_amount: 0,
            status: "pending",
            tracking_number: "MP001234567",
            shipping_address: "123 Nguyễn Trãi, Quận 1, TP.HCM",
            recipient_name: "Nguyễn Văn A",
            recipient_phone: "0901234567",
            shipping_method: "cod",
            estimated_delivery_date: "2025-02-05",
            created_at: "2025-01-30T10:00:00Z",
            customer_name: "Nguyễn Văn A",
            customer_avatar: "https://via.placeholder.com/40",
            notes: "Giao hàng giờ hành chính"
        },
        {
            order_item_id: 2,
            order_id: "ORD002",
            product_name: "Samsung Galaxy S24 Ultra",
            product_variant_name: "Titanium Black 512GB",
            product_variant_thumbnail: "https://via.placeholder.com/80",
            quantity: 2,
            unit_price: 29990000,
            total_price: 59980000,
            discount_amount: 1000000,
            shipping_fee: 0,
            tax_amount: 0,
            status: "confirmed",
            tracking_number: "MP001234568",
            shipping_address: "456 Lê Lợi, Quận 3, TP.HCM",
            recipient_name: "Trần Thị B",
            recipient_phone: "0901234568",
            shipping_method: "momo",
            estimated_delivery_date: "2025-02-06",
            created_at: "2025-01-29T14:30:00Z",
            customer_name: "Trần Thị B",
            customer_avatar: "https://via.placeholder.com/40",
            notes: ""
        },
        {
            order_item_id: 3,
            order_id: "ORD003",
            product_name: "MacBook Pro 14 inch M3",
            product_variant_name: "Space Black 1TB",
            product_variant_thumbnail: "https://via.placeholder.com/80",
            quantity: 1,
            unit_price: 52990000,
            total_price: 52990000,
            discount_amount: 2000000,
            shipping_fee: 50000,
            tax_amount: 0,
            status: "processing",
            tracking_number: "MP001234569",
            shipping_address: "789 Võ Văn Tần, Quận 3, TP.HCM",
            recipient_name: "Lê Văn C",
            recipient_phone: "0901234569",
            shipping_method: "vnpay",
            estimated_delivery_date: "2025-02-07",
            created_at: "2025-01-28T09:15:00Z",
            customer_name: "Lê Văn C",
            customer_avatar: "https://via.placeholder.com/40",
            notes: "Liên hệ trước khi giao"
        }
    ];

    // Fetch orders (mock implementation)
    const fetchOrders = async (page = 1, status = null, keyword = null) => {
        setLoading(true);
        setError(null);

        try {
            // Simulate API delay
            await new Promise(resolve => setTimeout(resolve, 1000));

            let filteredOrders = [...mockOrders];

            // Filter by status
            if (status) {
                filteredOrders = filteredOrders.filter(order => order.status === status);
            }

            // Filter by keyword
            if (keyword && keyword.trim()) {
                const searchTerm = keyword.trim().toLowerCase();
                filteredOrders = filteredOrders.filter(order =>
                    order.product_name.toLowerCase().includes(searchTerm) ||
                    order.customer_name.toLowerCase().includes(searchTerm) ||
                    order.order_id.toLowerCase().includes(searchTerm)
                );
            }

            setOrders(filteredOrders);
            setPagination({
                page: 1,
                limit: 10,
                total_items: filteredOrders.length,
                total_pages: Math.ceil(filteredOrders.length / 10),
                has_next: false,
                has_previous: false
            });
        } catch (err) {
            setError('Có lỗi xảy ra khi tải đơn hàng');
            setOrders([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        const currentStatus = statusTabs[activeTab]?.key;
        fetchOrders(1, currentStatus, searchQuery);
    }, [activeTab]);

    useEffect(() => {
        const timeoutId = setTimeout(() => {
            const currentStatus = statusTabs[activeTab]?.key;
            fetchOrders(1, currentStatus, searchQuery);
        }, 500);

        return () => clearTimeout(timeoutId);
    }, [searchQuery]);

    // Status badges
    const renderStatusBadge = (status) => {
        const statusConfig = {
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
            'cancelled': {
                text: 'Đã hủy',
                bg: '#FEE2E2',
                color: '#991B1B'
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

    // Format date
    const formatDate = (dateString) => {
        if (!dateString) return '';
        const date = new Date(dateString);
        return date.toLocaleDateString('vi-VN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    // Toggle order details
    const toggleOrderDetails = (orderItemId) => {
        setExpandedOrder(expandedOrder === orderItemId ? null : orderItemId);
    };

    // Handle status update
    const handleStatusUpdate = (order, newStatus) => {
        setSelectedOrder(order);
        setStatusToUpdate(newStatus);
        if (newStatus === 'cancelled') {
            onCancelModalOpen();
        } else {
            onStatusModalOpen();
        }
    };

    const confirmStatusUpdate = async () => {
        if (!selectedOrder) return;

        try {
            // Simulate API call
            await new Promise(resolve => setTimeout(resolve, 500));

            // Update order status in local state
            setOrders(prevOrders =>
                prevOrders.map(order =>
                    order.order_item_id === selectedOrder.order_item_id
                        ? { ...order, status: statusToUpdate }
                        : order
                )
            );

            onStatusModalClose();
            setSelectedOrder(null);
            setStatusToUpdate('');
        } catch (error) {
            console.error('Error updating status:', error);
        }
    };

    const confirmCancelOrder = async () => {
        if (!selectedOrder || !cancelReason.trim()) return;

        try {
            // Simulate API call
            await new Promise(resolve => setTimeout(resolve, 500));

            // Update order status in local state
            setOrders(prevOrders =>
                prevOrders.map(order =>
                    order.order_item_id === selectedOrder.order_item_id
                        ? { ...order, status: 'cancelled', cancelled_reason: cancelReason }
                        : order
                )
            );

            onCancelModalClose();
            setSelectedOrder(null);
            setStatusToUpdate('');
            setCancelReason('');
        } catch (error) {
            console.error('Error cancelling order:', error);
        }
    };

    // Get available status transitions
    const getAvailableStatuses = (currentStatus) => {
        const statusFlow = {
            'pending': ['confirmed', 'cancelled'],
            'confirmed': ['processing', 'cancelled'],
            'processing': ['ready_to_ship', 'cancelled'],
            'ready_to_ship': ['cancelled']
        };
        return statusFlow[currentStatus] || [];
    };

    // Get status label
    const getStatusLabel = (status) => {
        const labels = {
            'confirmed': 'Xác nhận đơn hàng',
            'processing': 'Bắt đầu chuẩn bị',
            'ready_to_ship': 'Sẵn sàng giao hàng',
            'cancelled': 'Hủy đơn hàng'
        };
        return labels[status] || status;
    };

    // Render single order item
    const renderOrderItem = (order) => {
        const isExpanded = expandedOrder === order.order_item_id;
        const availableStatuses = getAvailableStatuses(order.status);

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
                            src={order.customer_avatar}
                            name={order.customer_name}
                        />
                        <VStack align="start" spacing={1}>
                            <Text fontWeight="bold" noOfLines={1}>
                                Khách hàng: {order.customer_name}
                            </Text>
                            <Text fontSize="sm" color="gray.600">
                                Mã đơn: {order.order_id} • {formatDate(order.created_at)}
                            </Text>
                        </VStack>
                    </Flex>

                    <Flex align="center" gap={2} flexShrink={0}>
                        {renderStatusBadge(order.status)}
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
                                Phân loại: {order.product_variant_name}
                            </Text>
                            <Text fontSize="sm">Số lượng: x{order.quantity}</Text>
                            {order.notes && (
                                <Text fontSize="sm" color="blue.600" mt={1}>
                                    Ghi chú: {order.notes}
                                </Text>
                            )}
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
                    direction={{ base: "column", md: "row" }}
                    gap={{ base: 3, md: 0 }}
                >
                    <Button
                        size="sm"
                        variant="ghost"
                        leftIcon={isExpanded ? <ChevronUpIcon /> : <ChevronDownIcon />}
                        onClick={() => toggleOrderDetails(order.order_item_id)}
                    >
                        {isExpanded ? 'Ẩn chi tiết' : 'Xem chi tiết'}
                    </Button>

                    <Flex align="center" gap={3} direction={{ base: "column", md: "row" }}>
                        <Text fontSize="sm">Tổng tiền:</Text>
                        <Text fontSize="lg" fontWeight="bold" color={redColor}>
                            {formatCurrency(order.total_price + order.shipping_fee - order.discount_amount)}
                        </Text>

                        {/* Order actions */}
                        <HStack spacing={2}>
                            {availableStatuses.map(status => (
                                <Button
                                    key={status}
                                    size="sm"
                                    colorScheme={status === 'cancelled' ? 'red' : 'blue'}
                                    variant={status === 'cancelled' ? 'outline' : 'solid'}
                                    onClick={() => handleStatusUpdate(order, status)}
                                >
                                    {getStatusLabel(status)}
                                </Button>
                            ))}
                        </HStack>
                    </Flex>
                </Flex>

                {/* Expandable order details */}
                <Collapse in={isExpanded} animateOpacity>
                    <Box
                        p={4}
                        bg="gray.50"
                        borderTopWidth="1px"
                        borderColor={borderColor}
                    >
                        <VStack align="start" spacing={3}>
                            <Box>
                                <Text fontWeight="bold" mb={1}>Thông tin giao hàng:</Text>
                                <Text>{order.shipping_address}</Text>
                                <Text>Người nhận: {order.recipient_name}</Text>
                                <Text>SĐT: {order.recipient_phone}</Text>
                            </Box>

                            <Box>
                                <Text fontWeight="bold" mb={1}>Mã vận đơn:</Text>
                                <Text>{order.tracking_number}</Text>
                            </Box>

                            <Box>
                                <Text fontWeight="bold" mb={1}>Ngày giao dự kiến:</Text>
                                <Text>{formatDate(order.estimated_delivery_date)}</Text>
                            </Box>

                            <Box>
                                <Text fontWeight="bold" mb={1}>Phương thức thanh toán:</Text>
                                <Text>
                                    {order.shipping_method === 'cod' ? 'Thanh toán khi nhận hàng' :
                                        order.shipping_method === 'momo' ? 'Thanh toán bằng MoMo' :
                                            order.shipping_method === 'vnpay' ? 'Thanh toán bằng VNPay' :
                                                order.shipping_method}
                                </Text>
                            </Box>

                            {/* Price breakdown */}
                            <Box w="full">
                                <Text fontWeight="bold" mb={2}>Chi tiết giá:</Text>
                                <VStack align="stretch" spacing={1}>
                                    <Flex justify="space-between">
                                        <Text>Tổng tiền hàng:</Text>
                                        <Text>{formatCurrency(order.total_price)}</Text>
                                    </Flex>
                                    {order.discount_amount > 0 && (
                                        <Flex justify="space-between">
                                            <Text>Giảm giá:</Text>
                                            <Text color="red.500">-{formatCurrency(order.discount_amount)}</Text>
                                        </Flex>
                                    )}
                                    <Flex justify="space-between">
                                        <Text>Phí vận chuyển:</Text>
                                        <Text>{formatCurrency(order.shipping_fee)}</Text>
                                    </Flex>
                                    <Divider />
                                    <Flex justify="space-between" fontWeight="bold">
                                        <Text>Tổng thanh toán:</Text>
                                        <Text color={redColor}>
                                            {formatCurrency(order.total_price + order.shipping_fee - order.discount_amount)}
                                        </Text>
                                    </Flex>
                                </VStack>
                            </Box>
                        </VStack>
                    </Box>
                </Collapse>
            </Box>
        );
    };

    return (
        <Box>
            <Heading as="h1" size="lg" mb={6}>
                Quản Lý Đơn Hàng
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
                    colorScheme="blue"
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
                                    color: 'blue.500',
                                    borderColor: 'blue.500',
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
                                    color: 'blue.500'
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
                                            placeholder="Tìm kiếm theo tên sản phẩm, khách hàng hoặc mã đơn hàng"
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
                                        <Spinner size="lg" color="blue.500" />
                                    </Flex>
                                ) : (
                                    <>
                                        {orders.length > 0 ? (
                                            <>
                                                {orders.map(order => renderOrderItem(order))}
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

            {/* Status Update Modal */}
            <Modal isOpen={isStatusModalOpen} onClose={onStatusModalClose} isCentered>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Xác nhận cập nhật trạng thái</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <Text>
                            Bạn có chắc chắn muốn cập nhật trạng thái đơn hàng{' '}
                            <Text as="span" fontWeight="bold">
                                {selectedOrder?.order_id}
                            </Text>{' '}
                            thành{' '}
                            <Text as="span" fontWeight="bold" color="blue.500">
                                {getStatusLabel(statusToUpdate)}
                            </Text>?
                        </Text>
                    </ModalBody>
                    <ModalFooter>
                        <Button variant="ghost" mr={3} onClick={onStatusModalClose}>
                            Hủy
                        </Button>
                        <Button colorScheme="blue" onClick={confirmStatusUpdate}>
                            Xác nhận
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>

            {/* Cancel Order Modal */}
            <Modal isOpen={isCancelModalOpen} onClose={onCancelModalClose} isCentered>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Hủy đơn hàng</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <Text mb={4}>
                            Bạn có chắc chắn muốn hủy đơn hàng{' '}
                            <Text as="span" fontWeight="bold">
                                {selectedOrder?.order_id}
                            </Text>?
                        </Text>
                        <Text mb={2}>Lý do hủy đơn hàng:</Text>
                        <Textarea
                            placeholder="Nhập lý do hủy đơn hàng..."
                            value={cancelReason}
                            onChange={(e) => setCancelReason(e.target.value)}
                            rows={3}
                        />
                    </ModalBody>
                    <ModalFooter>
                        <Button variant="ghost" mr={3} onClick={onCancelModalClose}>
                            Hủy
                        </Button>
                        <Button
                            colorScheme="red"
                            onClick={confirmCancelOrder}
                            isDisabled={!cancelReason.trim()}
                        >
                            Xác nhận hủy
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </Box>
    );
};

export default SupplierOrders;