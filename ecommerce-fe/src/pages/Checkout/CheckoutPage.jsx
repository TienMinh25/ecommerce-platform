import React, { useState, useEffect } from 'react';
import {
    Container,
    Box,
    VStack,
    HStack,
    Text,
    Heading,
    Divider,
    Button,
    Flex,
    Image,
    Icon,
    useToast,
    Grid,
    GridItem
} from '@chakra-ui/react';
import {
    FiMapPin,
    FiTruck
} from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';
import PageTitle from '../PageTitle';
import AddressSelector from '../../components/checkout/AddressSelector';
import CheckoutVoucherSelector from '../../components/checkout/CheckoutVoucherSelector';
import PaymentMethodSelector from '../../components/checkout/PaymentMethodSelector';
import userMeService from '../../services/userMeService';
import paymentMethodService from '../../services/paymentMethodService';

const CheckoutPage = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const toast = useToast();

    // Get data from navigation state (from cart or product detail)
    const {
        selectedVoucher: initialVoucher,
        voucherDiscount,
        finalTotal,
        cartItems: passedCartItems
    } = location.state || {};

    // Mock cart items if not passed from previous page
    const [orderItems] = useState(passedCartItems || [
        {
            cart_item_id: 1,
            product_name: "Bộ sản phẩm Nông nghiệp & Vườn tược cao cấp",
            product_variant_thumbnail: "https://via.placeholder.com/80",
            variant_name: "Bạc, Chất liệu: Kim loại",
            price: 1189000,
            discount_price: 0,
            quantity: 1,
            attribute_values: [
                { attribute_name: "Màu sắc", attribute_value: "Bạc" },
                { attribute_name: "Chất liệu", attribute_value: "Kim loại" }
            ]
        }
    ]);

    const [selectedAddress, setSelectedAddress] = useState(null);
    const [selectedVoucher, setSelectedVoucher] = useState(initialVoucher || null);
    const [paymentMethod, setPaymentMethod] = useState('cod');
    const [isProcessing, setIsProcessing] = useState(false);

    // Set default address on component mount
    useEffect(() => {
        fetchDefaultAddress();
    }, []);

    const fetchDefaultAddress = async () => {
        try {
            const response = await userMeService.getAddresses({ page: 1, limit: 10 });
            const defaultAddr = response.data.find(addr => addr.is_default);
            if (defaultAddr) {
                setSelectedAddress(defaultAddr);
            }
        } catch (error) {
            console.error('Error fetching default address:', error);
        }
    };

    // Format currency
    const formatPrice = (price) => {
        return new Intl.NumberFormat('vi-VN', {
            style: 'currency',
            currency: 'VND',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        }).format(price);
    };

    // Calculate totals
    const subtotal = orderItems.reduce((total, item) => {
        const price = item.discount_price > 0 ? item.discount_price : item.price;
        return total + (price * item.quantity);
    }, 0);

    const shippingFee = 18300; // Fixed shipping fee as shown in image

    const calculateVoucherDiscount = () => {
        if (!selectedVoucher || subtotal === 0) return 0;

        let discount = 0;
        if (selectedVoucher.discount_type === 'percentage') {
            discount = (subtotal * selectedVoucher.discount_value) / 100;
        } else {
            discount = selectedVoucher.discount_value;
        }

        if (selectedVoucher.maximum_discount_amount) {
            discount = Math.min(discount, selectedVoucher.maximum_discount_amount);
        }

        return discount;
    };

    const voucherDiscountAmount = calculateVoucherDiscount();
    const totalAmount = subtotal + shippingFee - voucherDiscountAmount;

    const handleAddressSelect = (address) => {
        setSelectedAddress(address);
    };

    const handleVoucherSelect = (voucher) => {
        setSelectedVoucher(voucher);
    };

    const handlePaymentMethodSelect = (methodId) => {
        setPaymentMethod(methodId);
    };

    const handlePlaceOrder = async () => {
        if (!selectedAddress) {
            toast({
                title: 'Thiếu thông tin',
                description: 'Vui lòng chọn địa chỉ nhận hàng',
                status: 'warning',
                duration: 3000,
                isClosable: true,
            });
            return;
        }

        setIsProcessing(true);

        try {
            const orderData = {
                items: orderItems,
                delivery_address: selectedAddress,
                payment_method: paymentMethod,
                voucher: selectedVoucher,
                subtotal: subtotal,
                shipping_fee: shippingFee,
                voucher_discount: voucherDiscountAmount,
                total_amount: totalAmount
            };

            // Call API to create order
            const result = await paymentMethodService.createOrder(orderData);

            toast({
                title: 'Đặt hàng thành công!',
                description: 'Đơn hàng của bạn đã được xác nhận và đang được xử lý',
                status: 'success',
                duration: 5000,
                isClosable: true,
            });

            // Navigate to order success page or orders page
            navigate('/user/account/orders', {
                state: { orderId: result.data?.id }
            });
        } catch (error) {
            console.error('Error placing order:', error);
            toast({
                title: 'Đặt hàng thất bại',
                description: error.message || 'Có lỗi xảy ra khi đặt hàng. Vui lòng thử lại.',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsProcessing(false);
        }
    };

    return (
        <Container maxW="container.xl" py={6}>
            <PageTitle title="Thanh Toán - Minh Plaza" />

            <Heading size="lg" mb={6} color="gray.700">
                Thanh Toán
            </Heading>

            <Grid templateColumns={{ base: '1fr', lg: '2fr 1fr' }} gap={6}>
                {/* Left Column - Order Details */}
                <GridItem>
                    <VStack spacing={6} align="stretch">
                        {/* Delivery Address Section */}
                        <Box bg="white" p={4} borderRadius="md" borderWidth="1px">
                            <HStack mb={3}>
                                <Icon as={FiMapPin} color="red.500" />
                                <Text fontWeight="semibold" color="red.500">
                                    Địa Chỉ Nhận Hàng
                                </Text>
                            </HStack>

                            <AddressSelector
                                selectedAddress={selectedAddress}
                                onAddressSelect={handleAddressSelect}
                                orderTotal={subtotal}
                            />
                        </Box>

                        {/* Products Section */}
                        <Box bg="white" p={4} borderRadius="md" borderWidth="1px">
                            <Text fontWeight="semibold" mb={4}>
                                Sản phẩm
                            </Text>

                            <VStack spacing={4} align="stretch">
                                {orderItems.map((item) => (
                                    <Box key={item.cart_item_id}>
                                        <Flex align="center" spacing={4}>
                                            <Image
                                                src={item.product_variant_thumbnail}
                                                alt={item.product_name}
                                                boxSize="80px"
                                                objectFit="cover"
                                                borderRadius="md"
                                                mr={4}
                                            />

                                            <VStack align="start" spacing={1} flex="1">
                                                <Text fontWeight="medium" noOfLines={2}>
                                                    {item.product_name}
                                                </Text>
                                                <Text fontSize="sm" color="gray.500">
                                                    Phân Loại Hàng: {item.variant_name ||
                                                    (item.attribute_values?.map(attr => attr.attribute_value).join(', ') || 'Mặc định')
                                                }
                                                </Text>
                                                <Text fontSize="sm" color="gray.600">
                                                    x{item.quantity}
                                                </Text>
                                            </VStack>

                                            <VStack align="end" spacing={1}>
                                                {item.discount_price > 0 ? (
                                                    <>
                                                        <Text as="s" color="gray.500" fontSize="sm">
                                                            {formatPrice(item.price)}
                                                        </Text>
                                                        <Text fontWeight="medium" color="red.500">
                                                            {formatPrice(item.discount_price)}
                                                        </Text>
                                                    </>
                                                ) : (
                                                    <Text fontWeight="medium">
                                                        {formatPrice(item.price)}
                                                    </Text>
                                                )}
                                            </VStack>
                                        </Flex>

                                        {orderItems.indexOf(item) < orderItems.length - 1 && (
                                            <Divider mt={4} />
                                        )}
                                    </Box>
                                ))}
                            </VStack>

                            <Divider my={4} />

                            {/* Voucher Section */}
                            <VStack spacing={3} align="stretch">
                                <CheckoutVoucherSelector
                                    selectedVoucher={selectedVoucher}
                                    onVoucherSelect={handleVoucherSelect}
                                    cartTotal={subtotal}
                                />
                            </VStack>

                            <Divider my={4} />

                            {/* Shipping Method */}
                            <Box>
                                <HStack justify="space-between" mb={2}>
                                    <HStack>
                                        <Icon as={FiTruck} color="blue.500" />
                                        <Text fontWeight="medium">Phương thức vận chuyển:</Text>
                                    </HStack>
                                </HStack>

                                <Box p={3} bg="blue.50" borderRadius="md" borderWidth="1px" borderColor="blue.200">
                                    <Flex justify="space-between" align="center">
                                        <VStack align="start" spacing={1}>
                                            <Text fontWeight="medium">Nhanh</Text>
                                            <Text fontSize="xs" color="gray.600">
                                                Đảm bảo nhận hàng từ 25 Tháng 5 - 28 Tháng 5
                                            </Text>
                                            <Text fontSize="xs" color="green.600">
                                                Nhận Voucher trị giá ₫15.000 nếu đơn hàng được giao đến bạn sau ngày 28 Tháng 5 2025.
                                            </Text>
                                        </VStack>
                                        <Text fontWeight="medium">
                                            {formatPrice(shippingFee)}
                                        </Text>
                                    </Flex>
                                </Box>
                            </Box>
                        </Box>

                        {/* Payment Method Section */}
                        <PaymentMethodSelector
                            selectedPaymentMethod={paymentMethod}
                            onPaymentMethodSelect={handlePaymentMethodSelect}
                        />
                    </VStack>
                </GridItem>

                {/* Right Column - Order Summary */}
                <GridItem>
                    <Box bg="white" p={4} borderRadius="md" borderWidth="1px" position="sticky" top="20px">
                        <VStack spacing={4} align="stretch">
                            <Text fontWeight="semibold" fontSize="lg">
                                Tóm tắt đơn hàng
                            </Text>

                            <Divider />

                            <VStack spacing={3} align="stretch">
                                <Flex justify="space-between">
                                    <Text>Tổng số tiền ({orderItems.length} sản phẩm):</Text>
                                    <Text>{formatPrice(subtotal)}</Text>
                                </Flex>

                                <Flex justify="space-between">
                                    <Text>Phí vận chuyển:</Text>
                                    <Text>{formatPrice(shippingFee)}</Text>
                                </Flex>

                                {selectedVoucher && voucherDiscountAmount > 0 && (
                                    <Flex justify="space-between">
                                        <Text>Minh Plaza Voucher giảm giá:</Text>
                                        <Text color="green.600">
                                            -{formatPrice(voucherDiscountAmount)}
                                        </Text>
                                    </Flex>
                                )}

                                <Divider />

                                <Flex justify="space-between" align="center">
                                    <Text fontWeight="bold" fontSize="lg">
                                        Tổng thanh toán:
                                    </Text>
                                    <Text fontWeight="bold" fontSize="xl" color="red.500">
                                        {formatPrice(totalAmount)}
                                    </Text>
                                </Flex>
                            </VStack>

                            <Divider />

                            <Button
                                colorScheme="red"
                                size="lg"
                                w="100%"
                                onClick={handlePlaceOrder}
                                isLoading={isProcessing}
                                loadingText="Đang xử lý..."
                                isDisabled={!selectedAddress}
                            >
                                Đặt Hàng
                            </Button>

                            <Text fontSize="xs" color="gray.500" textAlign="center">
                                Nhấn "Đặt hàng" đồng nghĩa với việc bạn đồng ý tuân theo
                                <Text as="span" color="blue.500" cursor="pointer"> Điều khoản Minh Plaza</Text>
                            </Text>
                        </VStack>
                    </Box>
                </GridItem>
            </Grid>
        </Container>
    );
};

export default CheckoutPage;