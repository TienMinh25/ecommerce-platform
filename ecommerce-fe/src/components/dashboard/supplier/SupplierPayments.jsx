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
    VStack,
    HStack,
    Divider,
    useColorModeValue,
    Spinner,
    Alert,
    AlertIcon,
    Card,
    CardBody,
    Stat,
    StatLabel,
    StatNumber,
    StatHelpText,
    StatArrow,
    Table,
    Thead,
    Tbody,
    Tr,
    Th,
    Td,
    TableContainer,
    Select,
    SimpleGrid,
    Icon,
} from '@chakra-ui/react';
import { SearchIcon, CalendarIcon } from '@chakra-ui/icons';
import { FiDollarSign, FiTrendingUp, FiClock, FiCheckCircle } from 'react-icons/fi';

const SupplierPayments = () => {
    const [searchQuery, setSearchQuery] = useState('');
    const [activeTab, setActiveTab] = useState(0);
    const [payments, setPayments] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [selectedMonth, setSelectedMonth] = useState(new Date().getMonth() + 1);
    const [selectedYear, setSelectedYear] = useState(new Date().getFullYear());

    // Payment status tabs
    const statusTabs = [
        { key: null, label: 'Tất cả' },
        { key: 'pending', label: 'Chờ thanh toán' },
        { key: 'processing', label: 'Đang xử lý' },
        { key: 'completed', label: 'Đã thanh toán' },
        { key: 'failed', label: 'Thất bại' },
    ];

    // Colors
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const greenColor = useColorModeValue('green.500', 'green.300');
    const redColor = useColorModeValue('red.500', 'red.300');
    const blueColor = useColorModeValue('blue.500', 'blue.300');

    // Mock payment data
    const mockPayments = [
        {
            payment_id: 'PAY001',
            order_id: 'ORD001',
            amount: 32990000,
            commission_rate: 5,
            commission_amount: 1649500,
            net_amount: 31340500,
            status: 'completed',
            payment_method: 'bank_transfer',
            created_at: '2025-01-30T10:00:00Z',
            processed_at: '2025-01-30T14:00:00Z',
            customer_name: 'Nguyễn Văn A',
            product_name: 'iPhone 14 Pro Max 256GB'
        },
        {
            payment_id: 'PAY002',
            order_id: 'ORD002',
            amount: 59980000,
            commission_rate: 5,
            commission_amount: 2999000,
            net_amount: 56981000,
            status: 'processing',
            payment_method: 'momo',
            created_at: '2025-01-29T14:30:00Z',
            processed_at: null,
            customer_name: 'Trần Thị B',
            product_name: 'Samsung Galaxy S24 Ultra'
        },
        {
            payment_id: 'PAY003',
            order_id: 'ORD003',
            amount: 52990000,
            commission_rate: 5,
            commission_amount: 2649500,
            net_amount: 50340500,
            status: 'pending',
            payment_method: 'vnpay',
            created_at: '2025-01-28T09:15:00Z',
            processed_at: null,
            customer_name: 'Lê Văn C',
            product_name: 'MacBook Pro 14 inch M3'
        }
    ];

    // Mock summary data
    const mockSummary = {
        total_revenue: 145960000,
        total_commission: 7298000,
        total_net_amount: 138662000,
        pending_amount: 50340500,
        completed_amount: 88321500,
        total_transactions: 3,
        growth_rate: 12.5
    };

    // Fetch payments (mock implementation)
    const fetchPayments = async (status = null, keyword = null) => {
        setLoading(true);
        setError(null);

        try {
            // Simulate API delay
            await new Promise(resolve => setTimeout(resolve, 1000));

            let filteredPayments = [...mockPayments];

            // Filter by status
            if (status) {
                filteredPayments = filteredPayments.filter(payment => payment.status === status);
            }

            // Filter by keyword
            if (keyword && keyword.trim()) {
                const searchTerm = keyword.trim().toLowerCase();
                filteredPayments = filteredPayments.filter(payment =>
                    payment.payment_id.toLowerCase().includes(searchTerm) ||
                    payment.order_id.toLowerCase().includes(searchTerm) ||
                    payment.customer_name.toLowerCase().includes(searchTerm) ||
                    payment.product_name.toLowerCase().includes(searchTerm)
                );
            }

            setPayments(filteredPayments);
        } catch (err) {
            setError('Có lỗi xảy ra khi tải dữ liệu thanh toán');
            setPayments([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        const currentStatus = statusTabs[activeTab]?.key;
        fetchPayments(currentStatus, searchQuery);
    }, [activeTab]);

    useEffect(() => {
        const timeoutId = setTimeout(() => {
            const currentStatus = statusTabs[activeTab]?.key;
            fetchPayments(currentStatus, searchQuery);
        }, 500);

        return () => clearTimeout(timeoutId);
    }, [searchQuery]);

    // Format currency
    const formatCurrency = (amount) => {
        return `₫${amount.toLocaleString('vi-VN')}`;
    };

    // Format date
    const formatDate = (dateString) => {
        if (!dateString) return 'Chưa xử lý';
        const date = new Date(dateString);
        return date.toLocaleDateString('vi-VN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    // Status badges
    const renderStatusBadge = (status) => {
        const statusConfig = {
            'pending': {
                text: 'Chờ thanh toán',
                bg: '#FED7AA',
                color: '#C2410C'
            },
            'processing': {
                text: 'Đang xử lý',
                bg: '#E9D5FF',
                color: '#6B21A8'
            },
            'completed': {
                text: 'Đã thanh toán',
                bg: '#D1FAE5',
                color: '#065F46'
            },
            'failed': {
                text: 'Thất bại',
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

    // Get payment method label
    const getPaymentMethodLabel = (method) => {
        const methods = {
            'bank_transfer': 'Chuyển khoản ngân hàng',
            'momo': 'Ví MoMo',
            'vnpay': 'VNPay',
            'zalopay': 'ZaloPay',
            'cod': 'Thanh toán khi nhận hàng'
        };
        return methods[method] || method;
    };

    return (
        <Box>
            <Heading as="h1" size="lg" mb={6}>
                Quản Lý Thanh Toán
            </Heading>

            {/* Summary Cards */}
            <SimpleGrid columns={{ base: 1, md: 2, lg: 4 }} spacing={6} mb={8}>
                <Card>
                    <CardBody>
                        <Stat>
                            <StatLabel>
                                <HStack>
                                    <Icon as={FiDollarSign} color={greenColor} />
                                    <Text>Tổng doanh thu</Text>
                                </HStack>
                            </StatLabel>
                            <StatNumber color={greenColor}>
                                {formatCurrency(mockSummary.total_revenue)}
                            </StatNumber>
                            <StatHelpText>
                                <StatArrow type="increase" />
                                {mockSummary.growth_rate}% so với tháng trước
                            </StatHelpText>
                        </Stat>
                    </CardBody>
                </Card>

                <Card>
                    <CardBody>
                        <Stat>
                            <StatLabel>
                                <HStack>
                                    <Icon as={FiTrendingUp} color={blueColor} />
                                    <Text>Số tiền thực nhận</Text>
                                </HStack>
                            </StatLabel>
                            <StatNumber color={blueColor}>
                                {formatCurrency(mockSummary.total_net_amount)}
                            </StatNumber>
                            <StatHelpText>
                                Sau khi trừ hoa hồng {mockSummary.total_commission.toLocaleString('vi-VN')}đ
                            </StatHelpText>
                        </Stat>
                    </CardBody>
                </Card>

                <Card>
                    <CardBody>
                        <Stat>
                            <StatLabel>
                                <HStack>
                                    <Icon as={FiClock} color="orange.500" />
                                    <Text>Chờ thanh toán</Text>
                                </HStack>
                            </StatLabel>
                            <StatNumber color="orange.500">
                                {formatCurrency(mockSummary.pending_amount)}
                            </StatNumber>
                            <StatHelpText>
                                Đang chờ xử lý
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
                                    <Text>Đã thanh toán</Text>
                                </HStack>
                            </StatLabel>
                            <StatNumber color={greenColor}>
                                {formatCurrency(mockSummary.completed_amount)}
                            </StatNumber>
                            <StatHelpText>
                                {mockSummary.total_transactions} giao dịch thành công
                            </StatHelpText>
                        </Stat>
                    </CardBody>
                </Card>
            </SimpleGrid>

            {/* Filters */}
            <Box
                bg={bgColor}
                borderRadius="md"
                borderWidth="1px"
                borderColor={borderColor}
                p={4}
                mb={6}
            >
                <Flex gap={4} direction={{ base: "column", md: "row" }}>
                    <InputGroup flex="2">
                        <InputLeftElement pointerEvents="none">
                            <SearchIcon color="gray.300" />
                        </InputLeftElement>
                        <Input
                            placeholder="Tìm kiếm theo mã thanh toán, đơn hàng, khách hàng..."
                            value={searchQuery}
                            onChange={e => setSearchQuery(e.target.value)}
                            bg={bgColor}
                            borderColor={borderColor}
                        />
                    </InputGroup>

                    <Select
                        value={selectedMonth}
                        onChange={(e) => setSelectedMonth(parseInt(e.target.value))}
                        w={{ base: "full", md: "150px" }}
                    >
                        {Array.from({ length: 12 }, (_, i) => (
                            <option key={i + 1} value={i + 1}>
                                Tháng {i + 1}
                            </option>
                        ))}
                    </Select>

                    <Select
                        value={selectedYear}
                        onChange={(e) => setSelectedYear(parseInt(e.target.value))}
                        w={{ base: "full", md: "120px" }}
                    >
                        {Array.from({ length: 5 }, (_, i) => {
                            const year = new Date().getFullYear() - i;
                            return (
                                <option key={year} value={year}>
                                    {year}
                                </option>
                            );
                        })}
                    </Select>
                </Flex>
            </Box>

            {/* Payment Table */}
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
                            <TabPanel key={tabIndex} p={0}>
                                {error && (
                                    <Alert status="error" m={4}>
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
                                        {payments.length > 0 ? (
                                            <TableContainer>
                                                <Table variant="simple">
                                                    <Thead bg="gray.50">
                                                        <Tr>
                                                            <Th>Mã thanh toán</Th>
                                                            <Th>Đơn hàng</Th>
                                                            <Th>Khách hàng</Th>
                                                            <Th>Sản phẩm</Th>
                                                            <Th isNumeric>Số tiền</Th>
                                                            <Th isNumeric>Hoa hồng</Th>
                                                            <Th isNumeric>Thực nhận</Th>
                                                            <Th>Phương thức</Th>
                                                            <Th>Trạng thái</Th>
                                                            <Th>Ngày tạo</Th>
                                                            <Th>Ngày xử lý</Th>
                                                        </Tr>
                                                    </Thead>
                                                    <Tbody>
                                                        {payments.map((payment) => (
                                                            <Tr key={payment.payment_id}>
                                                                <Td fontWeight="medium">
                                                                    {payment.payment_id}
                                                                </Td>
                                                                <Td>
                                                                    <Text color={blueColor} fontWeight="medium">
                                                                        {payment.order_id}
                                                                    </Text>
                                                                </Td>
                                                                <Td>{payment.customer_name}</Td>
                                                                <Td>
                                                                    <Text noOfLines={1} maxW="200px">
                                                                        {payment.product_name}
                                                                    </Text>
                                                                </Td>
                                                                <Td isNumeric fontWeight="medium">
                                                                    {formatCurrency(payment.amount)}
                                                                </Td>
                                                                <Td isNumeric color={redColor}>
                                                                    -{formatCurrency(payment.commission_amount)}
                                                                    <Text fontSize="xs" color="gray.500">
                                                                        ({payment.commission_rate}%)
                                                                    </Text>
                                                                </Td>
                                                                <Td isNumeric fontWeight="bold" color={greenColor}>
                                                                    {formatCurrency(payment.net_amount)}
                                                                </Td>
                                                                <Td>
                                                                    <Text fontSize="sm">
                                                                        {getPaymentMethodLabel(payment.payment_method)}
                                                                    </Text>
                                                                </Td>
                                                                <Td>
                                                                    {renderStatusBadge(payment.status)}
                                                                </Td>
                                                                <Td>
                                                                    <Text fontSize="sm">
                                                                        {formatDate(payment.created_at)}
                                                                    </Text>
                                                                </Td>
                                                                <Td>
                                                                    <Text fontSize="sm">
                                                                        {formatDate(payment.processed_at)}
                                                                    </Text>
                                                                </Td>
                                                            </Tr>
                                                        ))}
                                                    </Tbody>
                                                </Table>
                                            </TableContainer>
                                        ) : (
                                            <VStack spacing={4} py={10}>
                                                <Icon as={FiDollarSign} boxSize={12} color="gray.400" />
                                                <Text color="gray.500">Chưa có dữ liệu thanh toán</Text>
                                            </VStack>
                                        )}
                                    </>
                                )}
                            </TabPanel>
                        ))}
                    </TabPanels>
                </Tabs>
            </Box>

            {/* Payment Summary for Selected Period */}
            <Card mt={8}>
                <CardBody>
                    <Heading size="md" mb={4}>
                        Thống kê thanh toán tháng {selectedMonth}/{selectedYear}
                    </Heading>
                    <SimpleGrid columns={{ base: 1, md: 3 }} spacing={6}>
                        <Box>
                            <Text color="gray.600" fontSize="sm">Tổng giao dịch</Text>
                            <Text fontSize="2xl" fontWeight="bold">
                                {mockSummary.total_transactions}
                            </Text>
                        </Box>
                        <Box>
                            <Text color="gray.600" fontSize="sm">Tổng doanh thu</Text>
                            <Text fontSize="2xl" fontWeight="bold" color={greenColor}>
                                {formatCurrency(mockSummary.total_revenue)}
                            </Text>
                        </Box>
                        <Box>
                            <Text color="gray.600" fontSize="sm">Thực nhận sau hoa hồng</Text>
                            <Text fontSize="2xl" fontWeight="bold" color={blueColor}>
                                {formatCurrency(mockSummary.total_net_amount)}
                            </Text>
                        </Box>
                    </SimpleGrid>

                    <Divider my={4} />

                    <Box>
                        <Text color="gray.600" fontSize="sm" mb={2}>Chi tiết hoa hồng</Text>
                        <HStack justify="space-between">
                            <Text>Tỷ lệ hoa hồng trung bình:</Text>
                            <Text fontWeight="medium">5%</Text>
                        </HStack>
                        <HStack justify="space-between">
                            <Text>Tổng hoa hồng đã trừ:</Text>
                            <Text fontWeight="medium" color={redColor}>
                                {formatCurrency(mockSummary.total_commission)}
                            </Text>
                        </HStack>
                    </Box>
                </CardBody>
            </Card>
        </Box>
    );
};

export default SupplierPayments;