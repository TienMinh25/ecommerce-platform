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
        cartItems: passedCartItems,
        fromCart = false,
        fromProductDetail = false
    } = location.state || {};

    // Redirect if no cart items
    useEffect(() => {
        if (!passedCartItems || passedCartItems.length === 0) {
            toast({
                title: 'Không có sản phẩm',
                description: 'Không có sản phẩm nào để thanh toán',
                status: 'warning',
                duration: 3000,
                isClosable: true,
            });
            navigate('/cart');
            return;
        }
    }, [passedCartItems, navigate, toast]);

    // Enhanced cart items with shipping calculation
    const [orderItems] = useState(() => {
        if (!passedCartItems) return [];

        return passedCartItems.map(item => ({
            ...item,
            // Convert field names to match API expectations
            product_variant_name: item.variant_name,
            product_variant_image_url: item.product_variant_thumbnail,
            // Calculate shipping fee per item (base fee 18,300 VND)
            shipping_fee: calculateItemShippingFee(item),
            // Calculate estimated delivery date (current date + 5 days)
            estimated_delivery_date: getEstimatedDeliveryDate()
        }));
    });

    const [selectedAddress, setSelectedAddress] = useState(null);
    const [selectedVoucher, setSelectedVoucher] = useState(initialVoucher || null);
    const [paymentMethod, setPaymentMethod] = useState('cod'); // Default to COD
    const [paymentMethods, setPaymentMethods] = useState([]);
    const [isLoadingPaymentMethods, setIsLoadingPaymentMethods] = useState(true);
    const [isProcessing, setIsProcessing] = useState(false);

    // Calculate shipping fee for individual item
    function calculateItemShippingFee(item) {
        const baseShippingFee = 18300; // Base shipping fee in VND
        const itemValue = item.discount_price > 0 ? item.discount_price : item.price;

        // Quantity discount: More items = lower shipping per item
        let quantityMultiplier = 1;
        if (item.quantity >= 5) {
            quantityMultiplier = 0.6; // 40% off for 5+ items
        } else if (item.quantity >= 3) {
            quantityMultiplier = 0.7; // 30% off for 3+ items
        } else if (item.quantity >= 2) {
            quantityMultiplier = 0.8; // 20% off for 2+ items
        }

        // Value-based discount: Higher value items = lower shipping rate
        let valueMultiplier = 1;
        if (itemValue >= 1000000) { // 1M+
            valueMultiplier = 0.5; // 50% off shipping for expensive items
        } else if (itemValue >= 500000) { // 500k+
            valueMultiplier = 0.7; // 30% off shipping
        } else if (itemValue >= 200000) { // 200k+
            valueMultiplier = 0.85; // 15% off shipping
        }

        // Calculate final shipping fee
        return Math.round(baseShippingFee * quantityMultiplier * valueMultiplier);
    }

    // Get estimated delivery date (current date + 5 days in UTC)
    function getEstimatedDeliveryDate() {
        const deliveryDate = new Date();
        deliveryDate.setDate(deliveryDate.getDate() + 5);
        return deliveryDate.toISOString();
    }

    // Set default address and fetch payment methods on component mount
    useEffect(() => {
        fetchDefaultAddress();
        fetchPaymentMethods();
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

    const fetchPaymentMethods = async () => {
        try {
            setIsLoadingPaymentMethods(true);
            const response = await paymentMethodService.getPaymentMethods();
            if (response && response.data) {
                setPaymentMethods(response.data);

                // Set default payment method to 'cod' if available
                const codMethod = response.data.find(method => method.code === 'cod');
                if (codMethod) {
                    setPaymentMethod('cod');
                }
            }
        } catch (error) {
            console.error('Error fetching payment methods:', error);
            // Keep default 'cod' if API fails
            setPaymentMethod('cod');
        } finally {
            setIsLoadingPaymentMethods(false);
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

    // Format date for shipping display
    const formatShippingDate = (date) => {
        return date.toLocaleDateString('vi-VN', {
            day: 'numeric',
            month: 'long'
        });
    };

    // Calculate shipping dates for display (today + 2 days to today + 5 days)
    const getShippingDates = () => {
        const today = new Date();
        const deliveryStart = new Date();
        const deliveryEnd = new Date();

        deliveryStart.setDate(today.getDate() + 2);
        deliveryEnd.setDate(today.getDate() + 5);

        return {
            startDate: formatShippingDate(deliveryStart),
            endDate: formatShippingDate(deliveryEnd),
            guaranteeDate: formatShippingDate(deliveryEnd)
        };
    };

    const shippingDates = getShippingDates();

    // Calculate total shipping fee based on all items
    const calculateTotalShippingFee = () => {
        // Calculate subtotal to check for free shipping
        const subtotal = orderItems.reduce((total, item) => {
            const price = item.discount_price > 0 ? item.discount_price : item.price;
            return total + (price * item.quantity);
        }, 0);

        // Free shipping for orders over 500k
        if (subtotal >= 500000) {
            return 0;
        }

        // Sum individual item shipping fees
        return orderItems.reduce((total, item) => {
            return total + item.shipping_fee;
        }, 0);
    };

    // Calculate totals
    const subtotal = orderItems.reduce((total, item) => {
        const price = item.discount_price > 0 ? item.discount_price : item.price;
        return total + (price * item.quantity);
    }, 0);

    const shippingFee = calculateTotalShippingFee();
    const isFreeShipping = subtotal >= 500000;

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
        // Validate required fields
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

        if (!paymentMethod) {
            toast({
                title: 'Thiếu thông tin',
                description: 'Vui lòng chọn phương thức thanh toán',
                status: 'warning',
                duration: 3000,
                isClosable: true,
            });
            return;
        }

        setIsProcessing(true);

        try {
            // Prepare checkout data according to API requirements
            const checkoutData = {
                items: orderItems.map(item => ({
                    product_id: item.product_id,
                    product_variant_id: item.product_variant_id,
                    product_name: item.product_name,
                    product_variant_name: item.product_variant_name || item.variant_name,
                    product_variant_image_url: item.product_variant_image_url || item.product_variant_thumbnail,
                    quantity: item.quantity,
                    estimated_delivery_date: item.estimated_delivery_date,
                    shipping_fee: isFreeShipping ? 0 : item.shipping_fee
                })),
                coupon_id: selectedVoucher?.id || null,
                method_type: paymentMethod, // 'cod' or 'momo'
                shipping_address: `${selectedAddress.street}, ${selectedAddress.ward}, ${selectedAddress.district}, ${selectedAddress.province}`,
                recipient_name: selectedAddress.recipient_name,
                recipient_phone: selectedAddress.phone
            };

            console.log('Checkout data to be sent:', checkoutData);

            // Call API to create order
            const result = await paymentMethodService.createOrder(checkoutData);

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

    // Early return if no order items
    if (!orderItems || orderItems.length === 0) {
        return (
            <Container maxW="container.xl" py={6}>
                <PageTitle title="Thanh Toán - Minh Plaza" />
                <Box textAlign="center" py={10}>
                    <Text>Không có sản phẩm để thanh toán</Text>
                    <Button onClick={() => navigate('/cart')} mt={4}>
                        Quay lại giỏ hàng
                    </Button>
                </Box>
            </Container>
        );
    }

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

                        {/* Products Section - Separated by Items */}
                        {orderItems.map((item, index) => (
                            <Box key={item.cart_item_id} bg="white" p={4} borderRadius="md" borderWidth="1px">
                                <VStack spacing={4} align="stretch">
                                    <Flex align="center" spacing={4}>
                                        <Image
                                            src={item.product_variant_thumbnail || item.product_variant_image_url}
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
                                                Phân Loại Hàng: {item.variant_name || item.product_variant_name || 'Mặc định'}
                                            </Text>
                                            <Text fontSize="sm" color="gray.600">
                                                x{item.quantity}
                                            </Text>
                                            {/* Show shipping discount info */}
                                            {!isFreeShipping && (
                                                <Text fontSize="xs" color="blue.600">
                                                    Ship: {formatPrice(item.shipping_fee / item.quantity)}/sp
                                                    {item.quantity >= 2 && (
                                                        <Text as="span" color="green.600" ml={1}>
                                                            (Giảm {item.quantity >= 5 ? '40%' : item.quantity >= 3 ? '30%' : '20%'})
                                                        </Text>
                                                    )}
                                                    {((item.discount_price > 0 ? item.discount_price : item.price) >= 200000) && (
                                                        <Text as="span" color="purple.600" ml={1}>
                                                            (VIP giảm)
                                                        </Text>
                                                    )}
                                                </Text>
                                            )}
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

                                    <Divider />

                                    {/* Individual Shipping Method for this item */}
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
                                                        Đảm bảo nhận hàng từ {shippingDates.startDate} - {shippingDates.endDate}
                                                    </Text>
                                                </VStack>
                                                <VStack align="end" spacing={1}>
                                                    {isFreeShipping ? (
                                                        <Text fontWeight="medium" color="green.600">
                                                            Miễn phí
                                                        </Text>
                                                    ) : (
                                                        <VStack align="end" spacing={0}>
                                                            <Text fontWeight="medium">
                                                                {formatPrice(item.shipping_fee)}
                                                            </Text>
                                                            <Text fontSize="xs" color="gray.500">
                                                                ({formatPrice(item.shipping_fee / item.quantity)}/sp)
                                                            </Text>
                                                        </VStack>
                                                    )}
                                                </VStack>
                                            </Flex>
                                        </Box>
                                    </Box>
                                </VStack>
                            </Box>
                        ))}

                        {/* Voucher Section - Separate */}
                        <Box bg="white" p={4} borderRadius="md" borderWidth="1px">
                            <VStack spacing={3} align="stretch">
                                <CheckoutVoucherSelector
                                    selectedVoucher={selectedVoucher}
                                    onVoucherSelect={handleVoucherSelect}
                                    cartTotal={subtotal}
                                />
                            </VStack>
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
                                    {isFreeShipping ? (
                                        <VStack align="end" spacing={0}>
                                            <Text as="s" color="gray.400" fontSize="sm">
                                                {formatPrice(orderItems.reduce((total, item) => total + item.shipping_fee, 0))}
                                            </Text>
                                            <Text color="green.600" fontWeight="medium">
                                                Miễn phí
                                            </Text>
                                        </VStack>
                                    ) : (
                                        <Text>{formatPrice(shippingFee)}</Text>
                                    )}
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
                                isDisabled={
                                    !selectedAddress ||
                                    !paymentMethod ||
                                    isLoadingPaymentMethods ||
                                    (paymentMethods.length > 0 && !paymentMethods.some(method => method.code === paymentMethod))
                                }
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