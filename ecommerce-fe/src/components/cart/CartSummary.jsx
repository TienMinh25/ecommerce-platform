import React from 'react';
import {
    VStack,
    Box,
    Text,
    Button,
    Divider,
    HStack,
    useToast
} from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';

const CartSummary = ({ selectedCount, total, selectedItems, cartItems }) => {
    const navigate = useNavigate();
    const toast = useToast();

    // Format currency function
    const formatPrice = (price) => {
        return new Intl.NumberFormat('vi-VN', {
            style: 'currency',
            currency: 'VND',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        }).format(price);
    };

    const handleCheckout = () => {
        if (selectedCount === 0) {
            toast({
                title: 'Chưa chọn sản phẩm',
                description: 'Vui lòng chọn ít nhất một sản phẩm để tiến hành thanh toán',
                status: 'warning',
                duration: 3000,
                isClosable: true,
            });
            return;
        }

        // Implementation of checkout logic will be added here
        // For now, we just navigate to a placeholder checkout page
        navigate('/checkout');
    };

    return (
        <Box bg="white" p={4} borderRadius="md" borderWidth="1px" position="sticky" top="100px">
            <VStack align="stretch" spacing={4}>
                <Text fontWeight="bold" fontSize="lg">Tóm tắt đơn hàng</Text>

                <Divider />

                <HStack justify="space-between">
                    <Text>Tạm tính ({selectedCount} sản phẩm):</Text>
                    <Text fontWeight="bold">{formatPrice(total)}</Text>
                </HStack>

                <HStack justify="space-between">
                    <Text>Phí vận chuyển:</Text>
                    <Text>{selectedCount > 0 ? 'Sẽ tính sau' : '--'}</Text>
                </HStack>

                <Divider />

                <HStack justify="space-between">
                    <Text fontWeight="bold">Tổng tiền:</Text>
                    <Text fontWeight="bold" fontSize="xl" color="red.500">
                        {formatPrice(total)}
                    </Text>
                </HStack>

                <Button
                    colorScheme="brand"
                    size="lg"
                    onClick={handleCheckout}
                    isDisabled={selectedCount === 0}
                >
                    Thanh toán ({selectedCount})
                </Button>
            </VStack>
        </Box>
    );
};

export default CartSummary;