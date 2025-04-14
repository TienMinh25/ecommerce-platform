import React, { useState } from 'react';
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
    Link
} from '@chakra-ui/react';
import { SearchIcon, InfoOutlineIcon } from '@chakra-ui/icons';

// Mock data for orders
const mockOrders = [
    {
        id: '123456789',
        status: 'cancelled',
        shop: {
            name: 'Thinkmax Official Store',
            isOfficial: true,
        },
        items: [
            {
                id: 'item1',
                name: 'Thinkmax 2 Mút Đệm Tai Nghe Không Dây Cho Sennheiser Hd4.50 Btnc',
                variant: 'Đen',
                price: 66000,
                originalPrice: 95000,
                quantity: 1,
                image: 'https://via.placeholder.com/80'
            }
        ],
        total: 66000,
        cancelledReason: 'Đã hủy tự động bởi hệ thống Shopee',
        createdAt: '2023-08-15T10:30:00Z'
    },
    {
        id: '987654321',
        status: 'completed',
        shop: {
            name: 'PiSoo - Dép Nam Nữ',
            isFavorite: true,
        },
        items: [
            {
                id: 'item2',
                name: 'Dép Nam và Nữ Quai Ngang Dion Unisex POPULAR Chọn tăng 1 Size',
                variant: 'Trắng Dino,41',
                price: 52000,
                originalPrice: 110000,
                quantity: 1,
                image: 'https://via.placeholder.com/80'
            }
        ],
        total: 52000,
        deliveredAt: '2023-08-10T15:45:00Z'
    }
];

const UserOrders = () => {
    const [searchQuery, setSearchQuery] = useState('');
    const [activeTab, setActiveTab] = useState(0);
    const [orders, setOrders] = useState(mockOrders);

    // Filter orders based on tab and search query
    const getFilteredOrders = () => {
        let filtered = [...orders];

        // Filter by tab/status
        if (activeTab === 1) {
            filtered = filtered.filter(order => order.status === 'pending');
        } else if (activeTab === 2) {
            filtered = filtered.filter(order => order.status === 'shipping');
        } else if (activeTab === 3) {
            filtered = filtered.filter(order => order.status === 'delivered');
        } else if (activeTab === 4) {
            filtered = filtered.filter(order => order.status === 'completed');
        } else if (activeTab === 5) {
            filtered = filtered.filter(order => order.status === 'cancelled');
        } else if (activeTab === 6) {
            filtered = filtered.filter(order => order.status === 'returned');
        }

        // Filter by search query
        if (searchQuery.trim()) {
            const query = searchQuery.toLowerCase();
            filtered = filtered.filter(order =>
                order.id.toLowerCase().includes(query) ||
                order.items.some(item =>
                    item.name.toLowerCase().includes(query) ||
                    (item.variant && item.variant.toLowerCase().includes(query))
                ) ||
                order.shop.name.toLowerCase().includes(query)
            );
        }

        return filtered;
    };

    // Colors
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const redColor = useColorModeValue('red.500', 'red.300');
    const greenColor = useColorModeValue('green.500', 'green.300');

    // Status badges
    const renderStatusBadge = (status) => {
        let color = 'gray';
        let text = 'Unknown';

        switch (status) {
            case 'pending':
                color = 'yellow';
                text = 'Chờ thanh toán';
                break;
            case 'shipping':
                color = 'blue';
                text = 'Đang vận chuyển';
                break;
            case 'delivered':
                color = 'teal';
                text = 'Chờ giao hàng';
                break;
            case 'completed':
                color = 'green';
                text = 'HOÀN THÀNH';
                break;
            case 'cancelled':
                color = 'red';
                text = 'ĐÃ HỦY';
                break;
            case 'returned':
                color = 'purple';
                text = 'Trả hàng/Hoàn tiền';
                break;
        }

        return (
            <Badge
                colorScheme={color}
                px={2}
                py={1}
                borderRadius="sm"
                textTransform="uppercase"
                fontWeight="bold"
                fontSize="xs"
            >
                {text}
            </Badge>
        );
    };

    // Format currency
    const formatCurrency = (amount) => {
        return `đ${amount.toLocaleString('vi-VN')}`;
    };

    // Render single order
    const renderOrder = (order) => (
        <Box
            key={order.id}
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
            >
                <Flex align="center" gap={2}>
                    {order.shop.isOfficial && (
                        <Badge colorScheme="red" px={2} fontSize="xs">
                            HOT
                        </Badge>
                    )}
                    {order.shop.isFavorite && (
                        <Badge colorScheme="red" px={2} fontSize="xs">
                            Yêu thích
                        </Badge>
                    )}
                    <Text fontWeight="bold">{order.shop.name}</Text>
                    <Button size="xs" colorScheme="red" variant="outline">
                        Chat
                    </Button>
                    <Button size="xs" variant="outline">
                        Xem Shop
                    </Button>
                </Flex>

                <Flex align="center">
                    {order.status === 'completed' && (
                        <Flex align="center" mr={4} color={greenColor}>
                            <InfoOutlineIcon mr={1} />
                            <Text fontSize="sm">Đơn hàng đã được giao thành công</Text>
                        </Flex>
                    )}
                    {renderStatusBadge(order.status)}
                </Flex>
            </Flex>

            {/* Order items */}
            {order.items.map(item => (
                <Flex
                    key={item.id}
                    p={4}
                    borderBottomWidth="1px"
                    borderColor={borderColor}
                    align="center"
                >
                    <Flex flex="1" gap={4}>
                        <Image
                            src={item.image}
                            alt={item.name}
                            boxSize="80px"
                            objectFit="cover"
                            borderRadius="md"
                            border="1px solid"
                            borderColor={borderColor}
                        />

                        <Flex direction="column" flex="1">
                            <Text noOfLines={2} fontWeight="medium" mb={1}>
                                {item.name}
                            </Text>
                            <Text fontSize="sm" color="gray.600" mb={1}>
                                Phân loại hàng: {item.variant}
                            </Text>
                            <Text fontSize="sm">x{item.quantity}</Text>
                        </Flex>
                    </Flex>

                    <Flex direction="column" align="flex-end" minW="140px">
                        <Text textDecoration="line-through" color="gray.500" fontSize="sm">
                            {formatCurrency(item.originalPrice)}
                        </Text>
                        <Text color={redColor} fontWeight="bold">
                            {formatCurrency(item.price)}
                        </Text>
                    </Flex>
                </Flex>
            ))}

            {/* Order footer */}
            <Flex
                p={4}
                justify="space-between"
                align="center"
                bg="gray.50"
            >
                <Box>
                    {order.status === 'cancelled' && (
                        <Flex align="center" color="gray.500">
                            <InfoOutlineIcon mr={2} />
                            <Text fontSize="sm">{order.cancelledReason}</Text>
                        </Flex>
                    )}
                </Box>

                <Flex align="center" gap={3}>
                    <Text fontSize="sm">Thành tiền:</Text>
                    <Text fontSize="lg" fontWeight="bold" color={redColor}>
                        {formatCurrency(order.total)}
                    </Text>

                    {/* Order actions */}
                    {order.status === 'cancelled' && (
                        <>
                            <Button size="sm" colorScheme="red">
                                Mua Lại
                            </Button>
                            <Button size="sm" variant="outline">
                                Xem Chi Tiết Hủy Đơn
                            </Button>
                            <Button size="sm" variant="outline">
                                Liên Hệ Người Bán
                            </Button>
                        </>
                    )}

                    {order.status === 'completed' && (
                        <>
                            <Button size="sm" colorScheme="red">
                                Mua Lại
                            </Button>
                            <Button size="sm" variant="outline">
                                Liên Hệ Người Bán
                            </Button>
                        </>
                    )}
                </Flex>
            </Flex>
        </Box>
    );

    return (
        <Box>
            <Heading as="h1" size="lg" mb={6}>
                Đơn Hàng Của Tôi
            </Heading>

            <Tabs
                variant="line"
                colorScheme="red"
                onChange={index => setActiveTab(index)}
                mb={6}
            >
                <TabList>
                    <Tab _selected={{ color: redColor, borderColor: redColor }}>Tất cả</Tab>
                    <Tab _selected={{ color: redColor, borderColor: redColor }}>Chờ thanh toán</Tab>
                    <Tab _selected={{ color: redColor, borderColor: redColor }}>Vận chuyển</Tab>
                    <Tab _selected={{ color: redColor, borderColor: redColor }}>Chờ giao hàng</Tab>
                    <Tab _selected={{ color: redColor, borderColor: redColor }}>Hoàn thành</Tab>
                    <Tab _selected={{ color: redColor, borderColor: redColor }}>Đã hủy</Tab>
                    <Tab _selected={{ color: redColor, borderColor: redColor }}>Trả hàng/Hoàn tiền</Tab>
                </TabList>

                <TabPanels>
                    {/* We'll use the same panel for all tabs, but filter the orders */}
                    {[0, 1, 2, 3, 4, 5, 6].map(tabIndex => (
                        <TabPanel key={tabIndex} px={0}>
                            <Box mb={6}>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <SearchIcon color="gray.300" />
                                    </InputLeftElement>
                                    <Input
                                        placeholder="Bạn có thể tìm kiếm theo tên Shop, ID đơn hàng hoặc Tên Sản phẩm"
                                        value={searchQuery}
                                        onChange={e => setSearchQuery(e.target.value)}
                                        bg={bgColor}
                                        borderColor={borderColor}
                                    />
                                </InputGroup>
                            </Box>

                            {getFilteredOrders().length > 0 ? (
                                getFilteredOrders().map(order => renderOrder(order))
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
                        </TabPanel>
                    ))}
                </TabPanels>
            </Tabs>
        </Box>
    );
};

export default UserOrders;