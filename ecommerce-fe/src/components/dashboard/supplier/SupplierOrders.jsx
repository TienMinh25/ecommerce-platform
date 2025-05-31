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
    useDisclosure,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    Textarea,
    useToast,
} from '@chakra-ui/react';
import { ChevronDownIcon, ChevronUpIcon } from '@chakra-ui/icons';
import supplierService from '../../../services/supplierService';

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
    const [updating, setUpdating] = useState(false);

    const { isOpen: isStatusModalOpen, onOpen: onStatusModalOpen, onClose: onStatusModalClose } = useDisclosure();
    const { isOpen: isCancelModalOpen, onOpen: onCancelModalOpen, onClose: onCancelModalClose } = useDisclosure();
    const toast = useToast();

    // Status tabs for supplier - based on acceptStatus from backend
    const statusTabs = [
        { key: null, label: 'Tất cả' },
        { key: 'pending', label: 'Chờ xác nhận' },
        { key: 'confirmed', label: 'Đã xác nhận' },
        { key: 'processing', label: 'Đang chuẩn bị' },
        { key: 'ready_to_ship', label: 'Sẵn sàng giao' },
        { key: 'in_transit', label: 'Đang vận chuyển' },
        { key: 'out_for_delivery', label: 'Sắp giao' },
        { key: 'delivered', label: 'Đã giao' },
        { key: 'cancelled', label: 'Đã hủy' },
        { key: 'refunded', label: 'Đã hoàn tiền' },
    ];

    // Colors
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const redColor = useColorModeValue('red.500', 'red.300');
    const greenColor = useColorModeValue('green.500', 'green.300');

    // Fetch orders from API
    const fetchOrders = async (page = 1, status = null) => {
        setLoading(true);
        setError(null);

        try {
            // Build query parameters
            const params = {
                page: page,
                limit: pagination.limit
            };

            // Add status filter if provided
            if (status) {
                params.status = status;
            }

            // Call API through supplierService
            const response = await supplierService.getSupplierOrders(params);

            if (response && response.data) {
                setOrders(response.data);

                // Update pagination from response metadata
                if (response.metadata && response.metadata.pagination) {
                    setPagination(response.metadata.pagination);
                } else {
                    // Fallback pagination if not provided
                    setPagination({
                        page: page,
                        limit: params.limit,
                        total_items: response.data.length,
                        total_pages: Math.ceil(response.data.length / params.limit),
                        has_next: false,
                        has_previous: false
                    });
                }
            } else {
                setOrders([]);
                setPagination({
                    page: 1,
                    limit: 10,
                    total_items: 0,
                    total_pages: 0,
                    has_next: false,
                    has_previous: false
                });
            }
        } catch (err) {
            console.error('Error fetching supplier orders:', err);
            setError(err.response?.data?.error?.message || 'Có lỗi xảy ra khi tải đơn hàng');
            setOrders([]);
            setPagination({
                page: 1,
                limit: 10,
                total_items: 0,
                total_pages: 0,
                has_next: false,
                has_previous: false
            });
        } finally {
            setLoading(false);
        }
    };

    // Fetch orders when tab changes
    useEffect(() => {
        const currentStatus = statusTabs[activeTab]?.key;
        fetchOrders(1, currentStatus);
    }, [activeTab]);

    // Status badges with statuses matching supplier acceptStatus
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

        setUpdating(true);
        try {
            // Call API to update order status using the correct endpoint
            await supplierService.updateOrderStatus(
                selectedOrder.order_item_id,
                statusToUpdate
            );

            // Show success toast
            toast({
                title: 'Cập nhật thành công',
                description: `Đã cập nhật trạng thái đơn hàng thành "${getStatusLabel(statusToUpdate)}"`,
                status: 'success',
                duration: 3000,
                isClosable: true,
            });

            // Reload data with current tab status and pagination
            const currentStatus = statusTabs[activeTab]?.key;
            await fetchOrders(pagination.page, currentStatus);

            onStatusModalClose();
            setSelectedOrder(null);
            setStatusToUpdate('');
        } catch (error) {
            console.error('Error updating status:', error);
            const errorMessage = error.response?.data?.error?.message || 'Có lỗi xảy ra khi cập nhật trạng thái';
            setError(errorMessage);

            // Show error toast
            toast({
                title: 'Cập nhật thất bại',
                description: errorMessage,
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setUpdating(false);
        }
    };

    const confirmCancelOrder = async () => {
        if (!selectedOrder) return;

        setUpdating(true);
        try {
            // Call API to cancel order (no cancel reason in API)
            await supplierService.updateOrderStatus(
                selectedOrder.order_item_id,
                'cancelled'
            );

            // Show success toast
            toast({
                title: 'Hủy đơn hàng thành công',
                description: 'Đã hủy đơn hàng thành công',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });

            // Reload data with current tab status and pagination
            const currentStatus = statusTabs[activeTab]?.key;
            await fetchOrders(pagination.page, currentStatus);

            onCancelModalClose();
            setSelectedOrder(null);
            setStatusToUpdate('');
            setCancelReason('');
        } catch (error) {
            console.error('Error cancelling order:', error);
            const errorMessage = error.response?.data?.error?.message || 'Có lỗi xảy ra khi hủy đơn hàng';
            setError(errorMessage);

            // Show error toast
            toast({
                title: 'Hủy đơn hàng thất bại',
                description: errorMessage,
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setUpdating(false);
        }
    };

    // Get available status transitions based on current status
    const getAvailableStatuses = (currentStatus) => {
        // Business logic for supplier status transitions:
        if (currentStatus === 'pending') {
            // Pending orders can be confirmed or cancelled
            return ['confirmed', 'cancelled'];
        } else if (currentStatus === 'confirmed') {
            // Confirmed orders can start processing
            return ['processing'];
        }
        // All other statuses are view-only for suppliers
        return [];
    };

    // Get status label for buttons
    const getStatusLabel = (status) => {
        const labels = {
            'confirmed': 'Xác nhận đơn hàng',
            'cancelled': 'Hủy đơn hàng',
            'processing': 'Bắt đầu chuẩn bị'
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
                            name={order.recipient_name}
                        />
                        <VStack align="start" spacing={1}>
                            <Text fontWeight="bold" noOfLines={1}>
                                Khách hàng: {order.recipient_name}
                            </Text>
                            <Text fontSize="sm" color="gray.600">
                                Mã đơn: {order.order_item_id} • {formatDate(order.created_at)}
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
                                {formatCurrency(order.unit_price * order.quantity)}
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
                            {formatCurrency(order.total_price + order.shipping_fee + order.tax_amount - order.discount_amount)}
                        </Text>

                        {/* Order actions - Show for pending and confirmed orders */}
                        {availableStatuses.length > 0 && (
                            <HStack spacing={2} wrap="wrap">
                                {availableStatuses.map(status => (
                                    <Button
                                        key={status}
                                        size="sm"
                                        colorScheme={status === 'cancelled' ? 'red' : status === 'processing' ? 'green' : 'blue'}
                                        variant={status === 'cancelled' ? 'outline' : 'solid'}
                                        onClick={() => handleStatusUpdate(order, status)}
                                        isLoading={updating && selectedOrder?.order_item_id === order.order_item_id}
                                        isDisabled={updating}
                                    >
                                        {getStatusLabel(status)}
                                    </Button>
                                ))}
                            </HStack>
                        )}
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

                            {order.tracking_number && (
                                <Box>
                                    <Text fontWeight="bold" mb={1}>Mã vận đơn:</Text>
                                    <Text>{order.tracking_number}</Text>
                                </Box>
                            )}

                            {order.estimated_delivery_date && (
                                <Box>
                                    <Text fontWeight="bold" mb={1}>Ngày giao dự kiến:</Text>
                                    <Text>{formatDate(order.estimated_delivery_date)}</Text>
                                </Box>
                            )}

                            <Box>
                                <Text fontWeight="bold" mb={1}>Phương thức thanh toán:</Text>
                                <Text>
                                    {order.shipping_method === 'cod' ? 'Thanh toán khi nhận hàng' :
                                        order.shipping_method === 'momo' ? 'Thanh toán bằng MoMo' :
                                            'Thanh toán trực tuyến'}
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
                                    {order.tax_amount > 0 && (
                                        <Flex justify="space-between">
                                            <Text>Thuế:</Text>
                                            <Text>{formatCurrency(order.tax_amount)}</Text>
                                        </Flex>
                                    )}
                                    <Divider />
                                    <Flex justify="space-between" fontWeight="bold">
                                        <Text>Tổng thanh toán:</Text>
                                        <Text color={redColor}>
                                            {formatCurrency(order.total_price + order.shipping_fee + order.tax_amount - order.discount_amount)}
                                        </Text>
                                    </Flex>
                                </VStack>
                            </Box>

                            {/* Show cancel reason if order is cancelled */}
                            {order.status === 'cancelled' && order.cancelled_reason && (
                                <Box>
                                    <Text fontWeight="bold" mb={1}>Lý do hủy:</Text>
                                    <Text color="red.600">{order.cancelled_reason}</Text>
                                </Box>
                            )}
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

                                                {/* Pagination */}
                                                {pagination.total_pages > 1 && (
                                                    <Flex justify="center" mt={6} gap={2}>
                                                        <Button
                                                            size="sm"
                                                            onClick={() => {
                                                                const currentStatus = statusTabs[activeTab]?.key;
                                                                fetchOrders(pagination.page - 1, currentStatus);
                                                            }}
                                                            isDisabled={!pagination.has_previous || loading}
                                                        >
                                                            Trước
                                                        </Button>

                                                        <Text mx={4} alignSelf="center">
                                                            Trang {pagination.page} / {pagination.total_pages}
                                                        </Text>

                                                        <Button
                                                            size="sm"
                                                            onClick={() => {
                                                                const currentStatus = statusTabs[activeTab]?.key;
                                                                fetchOrders(pagination.page + 1, currentStatus);
                                                            }}
                                                            isDisabled={!pagination.has_next || loading}
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

            {/* Status Update Modal */}
            <Modal isOpen={isStatusModalOpen} onClose={onStatusModalClose} isCentered>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Xác nhận cập nhật trạng thái</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <Text>
                            Bạn có chắc chắn muốn{' '}
                            <Text as="span" fontWeight="bold" color="blue.500">
                                {getStatusLabel(statusToUpdate)}
                            </Text>{' '}
                            cho đơn hàng{' '}
                            <Text as="span" fontWeight="bold">
                                {selectedOrder?.order_item_id}
                            </Text>?
                        </Text>
                    </ModalBody>
                    <ModalFooter>
                        <Button variant="ghost" mr={3} onClick={onStatusModalClose} isDisabled={updating}>
                            Hủy
                        </Button>
                        <Button
                            colorScheme="blue"
                            onClick={confirmStatusUpdate}
                            isLoading={updating}
                        >
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
                        <Text>
                            Bạn có chắc chắn muốn hủy đơn hàng{' '}
                            <Text as="span" fontWeight="bold">
                                {selectedOrder?.order_item_id}
                            </Text>?
                        </Text>
                    </ModalBody>
                    <ModalFooter>
                        <Button variant="ghost" mr={3} onClick={onCancelModalClose} isDisabled={updating}>
                            Không
                        </Button>
                        <Button
                            colorScheme="red"
                            onClick={confirmCancelOrder}
                            isLoading={updating}
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