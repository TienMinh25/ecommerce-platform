import React from 'react';
import {
    Box,
    Text,
    Flex,
    VStack,
    Divider,
    useColorModeValue
} from '@chakra-ui/react';

const OrderPriceBreakdown = ({ order }) => {
    const redColor = useColorModeValue('red.500', 'red.300');

    // Format currency
    const formatCurrency = (amount) => {
        return `₫${amount.toLocaleString('vi-VN')}`;
    };

    // Calculate values
    const subtotal = order.total_price + order.discount_amount;
    const finalAmount = order.total_price + order.tax_amount + order.shipping_fee;

    return (
        <Box>
            <Text fontWeight="bold" mb={3}>Chi tiết thanh toán:</Text>
            <VStack spacing={2} align="stretch">
                <Flex justify="space-between">
                    <Text>Tổng tiền hàng:</Text>
                    <Text>{formatCurrency(subtotal)}</Text>
                </Flex>


                <Flex justify="space-between" color="green.500">
                    <Text>Voucher giảm giá:</Text>
                    <Text>-{formatCurrency(order.discount_amount)}</Text>
                </Flex>


                <Flex justify="space-between" color="green.500">
                    <Text>Phí vận chuyển:</Text>
                    <Text>-{formatCurrency(order.shipping_fee)}</Text>
                </Flex>


                <Flex justify="space-between" color="green.500">
                    <Text>Thuế:</Text>
                    <Text>-{formatCurrency(order.tax_amount)}</Text>
                </Flex>


                <Divider />

                <Flex justify="space-between" fontWeight="bold" fontSize="lg">
                    <Text>Thành tiền:</Text>
                    <Text color={redColor}>{formatCurrency(finalAmount)}</Text>
                </Flex>
            </VStack>
        </Box>
    );
};

export default OrderPriceBreakdown;