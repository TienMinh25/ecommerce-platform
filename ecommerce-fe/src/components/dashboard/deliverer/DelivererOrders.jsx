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
    useColorModeValue,
    Collapse,
    Spinner,
    Alert,
    AlertIcon,
    useDisclosure,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    Textarea,
    Card,
    CardBody,
    SimpleGrid,
    Stat,
    StatLabel,
    StatNumber,
    StatHelpText,
    Icon,
} from '@chakra-ui/react';
import { SearchIcon, ChevronDownIcon, ChevronUpIcon, PhoneIcon } from '@chakra-ui/icons';
import { FiTruck, FiMapPin, FiClock, FiCheckCircle, FiPackage, FiDollarSign } from 'react-icons/fi';

const DelivererOrders = () => {
    const [searchQuery, setSearchQuery] = useState('');
    const [activeTab, setActiveTab] = useState(0);
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [expandedOrder, setExpandedOrder] = useState(null);
    const [selectedOrder, setSelectedOrder] = useState(null);
    const [statusToUpdate, setStatusToUpdate] = useState('');
    const [deliveryNote, setDeliveryNote] = useState('');

    const { isOpen: isStatusModalOpen, onOpen: onStatusModalOpen, onClose: onStatusModalClose } = useDisclosure();

    // Status tabs for deliverer
    const statusTabs = [
        { key: null, label: 'Tất cả' },
        { key: 'ready_to_ship', label: 'Sẵn sàng lấy hàng' },
        { key: 'in_transit', label: 'Đang vận chuyển' },
        { key: 'out_for_delivery', label: 'Đang giao hàng' },
        { key: 'delivered', label: 'Đã giao thành công' },
        { key: 'delivery_failed', label: 'Giao thất bại' },
    ];

    // Colors
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const redColor = useColorModeValue('red.500', 'red.300');
    const greenColor = useColorModeValue('green.500', 'green.300');
    const blueColor = useColorModeValue('blue.500', 'blue.300');

    // Mock data for deliverer orders
    const mockOrders = [
        {
            order_item_id: 1,
            order_id: "ORD001",
            product_name: "iPhone 14 Pro Max 256GB",
            product_variant_name: "Deep Purple",
            product_variant_thumbnail: "https://via.placeholder.com/80",
            quantity: 1,
            total_price: 32990000,
            status: "ready_to_ship",
            tracking_number: "MP001234567",
            shipping_address: "123 Nguyễn Trãi, Phường Bến Thành, Quận 1, TP.HCM",
            recipient_name: "Nguyễn Văn A",
            recipient_phone: "0901234567",
            estimated_delivery_date: "2025-02-05",
            pickup_address: "456 Lê Lợi, Quận 3, TP.HCM",
            supplier_name: "Tech Store ABC",
            supplier_phone: "0902345678",
            created_at: "2025-01-30T10:00:00Z",
            distance: "5.2 km",
            delivery_fee: 25000,
            notes: "Giao hàng giờ hành chính"
        },
        {
            order_item_id: 2,
            order_id: "ORD002",
            product_name: "Samsung Galaxy S24 Ultra",
            product_variant_name: "Titanium Black 512GB",
            product_variant_thumbnail: "https://via.placeholder.com/80",
            quantity: 2,
            total_price: 59980000,
            status: "in_transit",
            tracking_number: "MP001234568",
            shipping_address: "789 Võ Văn Tần, Phường 6, Quận 3, TP.HCM",
            recipient_name: "Trần Thị B",
            recipient_phone: "0901234568",
            estimated_delivery_date: "2025-02-06",
            pickup_address: "321 Pasteur, Quận 1, TP.HCM",
            supplier_name: "Mobile World",
            supplier_phone: "0903456789",
            created_at: "2025-01-29T14:30:00Z",
            picked_up_at: "2025-01-30T08:00:00Z",
            distance: "3.8 km",
            delivery_fee: 20000,
            notes: "Khách yêu cầu gọi trước 15 phút"
        },
        {
            order_item_id: 3,
            order_id: "ORD003",
            product_name: "MacBook Pro 14 inch M3",
            product_variant_name: "Space Black 1TB",
            product_variant_thumbnail: "https://via.placeholder.com/80",
            quantity: 1,
            total_price: 52990000,
            status: "delivered",
            tracking_number: "MP001234569",
            shipping_address: "456 Điện Biên Phủ, Phường 25, Quận Bình Thạnh, TP.HCM",
            recipient_name: "Lê Văn C",
            recipient_phone: "0901234569",
            estimated_delivery_date: "2025-02-07",
            pickup_address: "654 Nguyễn Thị Minh Khai, Quận 3, TP.HCM",
            supplier_name: "Laptop Center",
            supplier_phone: "0904567890",
            created_at: "2025-01-28T09:15:00Z",
            picked_up_at: "2025-01-29T10:00:00Z",
            delivered_at: "2025-01-29T16:30:00Z",
            distance: "7.1 km",
            delivery_fee: 30000,
            notes: ""
        }
    ];

    // Mock summary statistics
    const mockSummary = {
        total_orders: 15,
        delivered_orders: 12,
        in_progress_orders: 2,
        failed_orders: 1,
        total_earnings: 375000,
        average_rating: 4.8
    };

    // Fetch orders (mock implementation)
    const fetchOrders = async (status = null, keyword = null) => {
        setLoading(true);
        setError(null);

        try {
            await new Promise(resolve => setTimeout(resolve, 1000));

            let filteredOrders = [...mockOrders];

            if (status) {
                filteredOrders = filteredOrders.filter(order => order.status === status);
            }

            if (keyword && keyword.trim()) {
                const searchTerm = keyword.trim().toLowerCase();
                filteredOrders = filteredOrders.filter(order =>
                    order.order_id.toLowerCase().includes(searchTerm) ||
                    order.recipient_name.toLowerCase().includes(searchTerm) ||
                    order.tracking_number.toLowerCase().includes(searchTerm) ||
                    order.supplier_name.toLowerCase().includes(searchTerm)
                );
            }

            setOrders(filteredOrders);
        } catch (err) {
            setError('Có lỗi xảy ra khi tải đơn hàng');
            setOrders([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        const currentStatus = statusTabs[activeTab]?.key;
        fetchOrders(currentStatus, searchQuery);
    }, [activeTab]);

    useEffect(() => {
        const timeoutId = setTimeout(() => {
            const currentStatus = statusTabs[activeTab]?.key;
            fetchOrders(currentStatus, searchQuery);
        }, 500);

        return () => clearTimeout(timeoutId);
    }, [searchQuery]);

    // Status badges
    const renderStatusBadge = (status) => {
        const statusConfig = {
            'ready_to_ship': {
                text: 'Sẵn sàng lấy hàng',
                bg: '#CFFAFE',
                color: '#155E75'
            },
            'in_transit': {
                text: 'Đang vận chuyển',
                bg: '#DBEAFE',
                color: '#1E40AF'
            },
            'out_for_delivery': {
                text: 'Đang giao hàng',
                bg: '#E9D5FF',
                color: '#6B21A8'
            },
            'delivered': {
                text: 'Đã giao thành công',
                bg: '#D1FAE5',
                color: '#065F46'
            },
            'delivery_failed': {
                text: 'Giao thất bại',
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
            <Badge
                px={3}
                py={1}
                borderRadius="md"
                bg={config.bg}
                color={config.color}
                fontWeight="medium"
                fontSize="xs"
            >
                {config.text}
            </Badge>
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
        onStatusModalOpen();
    };

    const confirmStatusUpdate = async () => {
        if (!selectedOrder) return;

        try {
            await new Promise(resolve => setTimeout(resolve, 500));

            const updateData = { status: statusToUpdate };
            if (statusToUpdate === 'in_transit') {
                updateData.picked_up_at = new Date().toISOString();
            } else if (statusToUpdate === 'delivered') {
                updateData.delivered_at = new Date().toISOString();
                updateData.delivery_note = deliveryNote;
            }

            setOrders(prevOrders =>
                prevOrders.map(order =>
                    order.order_item_id === selectedOrder.order_item_id
                        ? { ...order, ...updateData }
                        : order
                )
            );

            onStatusModalClose();
            setSelectedOrder(null);
            setStatusToUpdate('');
            setDeliveryNote('');
        } catch (error) {
            console.error('Error updating status:', error);
        }
    };

    // Get available status transitions
    const getAvailableStatuses = (currentStatus) => {
        const statusFlow = {
            'ready_to_ship': ['in_transit'],
            'in_transit': ['out_for_delivery'],
            'out_for_delivery': ['delivered', 'delivery_failed']
        };
        return statusFlow[currentStatus] || [];
    };

    // Get status label
    const getStatusLabel = (status) => {
        const labels = {
            'in_transit': 'Đã lấy hàng',
            'out_for_delivery': 'Bắt đầu giao hàng',
            'delivered': 'Giao thành công',
            'delivery_failed': 'Giao thất bại'
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
                    <VStack align="start" spacing={1} flex="1">
                        <HStack>
                            <Text fontWeight="bold">Mã đơn: {order.order_id}</Text>
                            <Text fontSize="sm" color="gray.600">
                                • {order.tracking_number}
                            </Text>
                        </HStack>
                        <HStack spacing={4}>
                            <HStack>
                                <Icon as={FiMapPin} color="gray.500" />
                                <Text fontSize="sm" color="gray.600">
                                    {order.distance}
                                </Text>
                            </HStack>
                            <HStack>
                                <Icon as={FiTruck} color="gray.500" />
                                <Text fontSize="sm" color={greenColor} fontWeight="medium">
                                    Phí giao: {formatCurrency(order.delivery_fee)}
                                </Text>
                            </HStack>
                        </HStack>
                    </VStack>

                    <Flex align="center" gap={2} flexShrink={0}>
                        {renderStatusBadge(order.status)}
                    </Flex>
                </Flex>

                {/* Order content */}
                <Box p={4}>
                    {/* Product info */}
                    <Flex mb={4} gap={4}>
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
                            <Text fontSize="sm" color={redColor} fontWeight="bold" mt={1}>
                                Giá trị: {formatCurrency(order.total_price)}
                            </Text>
                        </Flex>
                    </Flex>

                    {/* Delivery info */}
                    <SimpleGrid columns={{ base: 1, md: 2 }} spacing={4} mb={4}>
                        <Box>
                            <Text fontWeight="bold" mb={2} color={blueColor}>
                                <Icon as={FiPackage} mr={2} />
                                Thông tin lấy hàng
                            </Text>
                            <Text fontSize="sm" mb={1}>
                                <strong>Nhà cung cấp:</strong> {order.supplier_name}
                            </Text>
                            <Text fontSize="sm" mb={1}>
                                <strong>Địa chỉ:</strong> {order.pickup_address}
                            </Text>
                            <HStack>
                                <Text fontSize="sm">
                                    <strong>SĐT:</strong> {order.supplier_phone}
                                </Text>
                                <Button
                                    size="xs"
                                    colorScheme="green"
                                    leftIcon={<PhoneIcon />}
                                    as="a"
                                    href={`tel:${order.supplier_phone}`}
                                >
                                    Gọi
                                </Button>
                            </HStack>
                        </Box>

                        <Box>
                            <Text fontWeight="bold" mb={2} color={greenColor}>
                                <Icon as={FiTruck} mr={2} />
                                Thông tin giao hàng
                            </Text>
                            <Text fontSize="sm" mb={1}>
                                <strong>Người nhận:</strong> {order.recipient_name}
                            </Text>
                            <Text fontSize="sm" mb={1}>
                                <strong>Địa chỉ:</strong> {order.shipping_address}
                            </Text>
                            <HStack>
                                <Text fontSize="sm">
                                    <strong>SĐT:</strong> {order.recipient_phone}
                                </Text>
                                <Button
                                    size="xs"
                                    colorScheme="blue"
                                    leftIcon={<PhoneIcon />}
                                    as="a"
                                    href={`tel:${order.recipient_phone}`}
                                >
                                    Gọi
                                </Button>
                            </HStack>
                        </Box>
                    </SimpleGrid>

                    {order.notes && (
                        <Box mb={4} p={3} bg="yellow.50" borderRadius="md" borderLeft="4px solid" borderLeftColor="yellow.400">
                            <Text fontSize="sm">
                                <strong>Ghi chú:</strong> {order.notes}
                            </Text>
                        </Box>
                    )}

                    {/* Timeline */}
                    {(order.picked_up_at || order.delivered_at) && (
                        <Box mb={4}>
                            <Text fontWeight="bold" mb={2}>Lịch sử giao hàng:</Text>
                            <VStack align="start" spacing={2}>
                                <HStack>
                                    <Icon as={FiClock} color="gray.500" />
                                    <Text fontSize="sm">
                                        Tạo đơn: {formatDate(order.created_at)}
                                    </Text>
                                </HStack>
                                {order.picked_up_at && (
                                    <HStack>
                                        <Icon as={FiPackage} color={blueColor} />
                                        <Text fontSize="sm">
                                            Đã lấy hàng: {formatDate(order.picked_up_at)}
                                        </Text>
                                    </HStack>
                                )}
                                {order.delivered_at && (
                                    <HStack>
                                        <Icon as={FiCheckCircle} color={greenColor} />
                                        <Text fontSize="sm">
                                            Giao thành công: {formatDate(order.delivered_at)}
                                        </Text>
                                    </HStack>
                                )}
                            </VStack>
                        </Box>
                    )}

                    {/* Actions */}
                    <Flex
                        justify="space-between"
                        align="center"
                        pt={4}
                        borderTopWidth="1px"
                        borderColor={borderColor}
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

                        <HStack spacing={2}>
                            {availableStatuses.map(status => (
                                <Button
                                    key={status}
                                    size="sm"
                                    colorScheme={
                                        status === 'delivered' ? 'green' :
                                            status === 'delivery_failed' ? 'red' : 'blue'
                                    }
                                    onClick={() => handleStatusUpdate(order, status)}
                                >
                                    {getStatusLabel(status)}
                                </Button>
                            ))}
                        </HStack>
                    </Flex>
                </Box>

                {/* Expandable details */}
                <Collapse in={isExpanded} animateOpacity>
                    <Box
                        p={4}
                        bg="gray.50"
                        borderTopWidth="1px"
                        borderColor={borderColor}
                    >
                        <Text fontWeight="bold" mb={2}>Chi tiết bổ sung:</Text>
                        <Text fontSize="sm" mb={1}>
                            <strong>Ngày giao dự kiến:</strong> {formatDate(order.estimated_delivery_date)}
                        </Text>
                        <Text fontSize="sm">
                            <strong>Phí giao hàng:</strong> {formatCurrency(order.delivery_fee)}
                        </Text>
                    </Box>
                </Collapse>
            </Box>
        );
    };

    return (
        <Box>
            <Heading as="h1" size="lg" mb={6}>
                Quản Lý Giao Hàng
            </Heading>

            {/* Summary Cards */}
            <SimpleGrid columns={{ base: 1, md: 2, lg: 4 }} spacing={6} mb={8}>
                <Card>
                    <CardBody>
                        <Stat>
                            <StatLabel>
                                <HStack>
                                    <Icon as={FiPackage} color={blueColor} />
                                    <Text>Tổng đơn hàng</Text>
                                </HStack>
                            </StatLabel>
                            <StatNumber color={blueColor}>
                                {mockSummary.total_orders}
                            </StatNumber>
                            <StatHelpText>
                                Đơn hàng trong tháng
                            </StatHelpText>
                        </Stat>
                    </CardBody>
                </Card>

                <Card>
                    <CardBody>
                        <Stat>
                            <StatLabel>
                                <HStack>
                                    <Icon as={FiCheckCircle} color={greenColor} />
                                    <Text>Giao thành công</Text>
                                </HStack>
                            </StatLabel>
                            <StatNumber color={greenColor}>
                                {mockSummary.delivered_orders}
                            </StatNumber>
                            <StatHelpText>
                                Tỷ lệ: {((mockSummary.delivered_orders / mockSummary.total_orders) * 100).toFixed(1)}%
                            </StatHelpText>
                        </Stat>
                    </CardBody>
                </Card>

                <Card>
                    <CardBody>
                        <Stat>
                            <StatLabel>
                                <HStack>
                                    <Icon as={FiTruck} color="orange.500" />
                                    <Text>Đang thực hiện</Text>
                                </HStack>
                            </StatLabel>
                            <StatNumber color="orange.500">
                                {mockSummary.in_progress_orders}
                            </StatNumber>
                            <StatHelpText>
                                Đơn hàng đang giao
                            </StatHelpText>
                        </Stat>
                    </CardBody>
                </Card>

                <Card>
                    <CardBody>
                        <Stat>
                            <StatLabel>
                                <HStack>
                                    <Icon as={FiDollarSign} color={greenColor} />
                                    <Text>Thu nhập tháng</Text>
                                </HStack>
                            </StatLabel>
                            <StatNumber color={greenColor}>
                                {formatCurrency(mockSummary.total_earnings)}
                            </StatNumber>
                            <StatHelpText>
                                Đánh giá: ⭐ {mockSummary.average_rating}/5
                            </StatHelpText>
                        </Stat>
                    </CardBody>
                </Card>
            </SimpleGrid>

            {/* Orders List */}
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
                                            placeholder="Tìm kiếm theo mã đơn hàng, người nhận, nhà cung cấp..."
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
                                                <Icon as={FiTruck} boxSize={12} color="gray.400" />
                                                <Text color="gray.500">Chưa có đơn hàng giao</Text>
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
                    <ModalHeader>Cập nhật trạng thái đơn hàng</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <Text mb={4}>
                            Bạn có chắc chắn muốn cập nhật trạng thái đơn hàng{' '}
                            <Text as="span" fontWeight="bold">
                                {selectedOrder?.order_id}
                            </Text>{' '}
                            thành{' '}
                            <Text as="span" fontWeight="bold" color="blue.500">
                                {getStatusLabel(statusToUpdate)}
                            </Text>?
                        </Text>

                        {statusToUpdate === 'delivered' && (
                            <Box>
                                <Text mb={2}>Ghi chú giao hàng (tùy chọn):</Text>
                                <Textarea
                                    placeholder="Nhập ghi chú về việc giao hàng..."
                                    value={deliveryNote}
                                    onChange={(e) => setDeliveryNote(e.target.value)}
                                    rows={3}
                                />
                            </Box>
                        )}
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
        </Box>
    );
};

export default DelivererOrders;