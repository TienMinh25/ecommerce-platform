import React from 'react';
import {
    Box,
    Text,
    Flex,
    VStack,
    Divider,
    Grid,
    GridItem,
    useColorModeValue
} from '@chakra-ui/react';
import OrderPriceBreakdown from './OrderPriceBreakdown';

const OrderDetailsSection = ({ order }) => {
    const borderColor = useColorModeValue('gray.200', 'gray.700');

    // Format date
    const formatDate = (dateString) => {
        if (!dateString) return '';
        const date = new Date(dateString);
        return date.toLocaleDateString('vi-VN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit'
        });
    };

    // Format payment method
    const formatPaymentMethod = (method) => {
        const methodConfig = {
            'cod': 'Thanh toán khi nhận hàng',
            'momo': 'Thanh toán bằng MoMo',
            'vnpay': 'Thanh toán bằng VNPay',
            'zalopay': 'Thanh toán bằng ZaloPay'
        };
        return methodConfig[method] || method;
    };

    return (
        <Box
            p={4}
            bg="gray.50"
            borderTopWidth="1px"
            borderColor={borderColor}
        >
            <Grid templateColumns="repeat(2, 1fr)" gap={6}>
                <GridItem>
                    <VStack align="start" spacing={3}>
                        <Box>
                            <Text fontWeight="bold" mb={1}>Mã vận đơn:</Text>
                            <Text>{order.tracking_number}</Text>
                        </Box>

                        <Box>
                            <Text fontWeight="bold" mb={1}>Địa chỉ giao hàng:</Text>
                            <Text>{order.shipping_address}</Text>
                        </Box>

                        <Box>
                            <Text fontWeight="bold" mb={1}>Người nhận:</Text>
                            <Text>{order.recipient_name}</Text>
                            <Text color="gray.600">{order.recipient_phone}</Text>
                        </Box>

                        <Box>
                            <Text fontWeight="bold" mb={1}>Phương thức vận chuyển:</Text>
                            <Text>Vận chuyển nhanh</Text>
                        </Box>

                        <Box>
                            <Text fontWeight="bold" mb={1}>Phương thức thanh toán:</Text>
                            <Text>{formatPaymentMethod(order.shipping_method)}</Text>
                        </Box>
                    </VStack>
                </GridItem>

                <GridItem>
                    <VStack align="start" spacing={3}>
                        <Box>
                            <Text fontWeight="bold" mb={1}>Ngày giao dự kiến:</Text>
                            <Text>{formatDate(order.estimated_delivery_date)}</Text>
                        </Box>

                        {order.actual_delivery_date && (
                            <Box>
                                <Text fontWeight="bold" mb={1}>Ngày giao thực tế:</Text>
                                <Text>{formatDate(order.actual_delivery_date)}</Text>
                            </Box>
                        )}

                        {order.notes && (
                            <Box>
                                <Text fontWeight="bold" mb={1}>Ghi chú:</Text>
                                <Text>{order.notes}</Text>
                            </Box>
                        )}

                        {order.cancelled_reason && (
                            <Box>
                                <Text fontWeight="bold" mb={1}>Lý do hủy:</Text>
                                <Text color="red.500">{order.cancelled_reason}</Text>
                            </Box>
                        )}
                    </VStack>
                </GridItem>
            </Grid>

            <Divider my={4} />

            {/* Price breakdown component */}
            <OrderPriceBreakdown order={order} />
        </Box>
    );
};

export default OrderDetailsSection;