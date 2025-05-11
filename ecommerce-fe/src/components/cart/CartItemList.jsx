import React from 'react';
import {
    VStack,
    HStack,
    Box,
    Checkbox,
    Image,
    Text,
    Button,
    Divider,
    NumberInput,
    NumberInputField,
    NumberInputStepper,
    NumberIncrementStepper,
    NumberDecrementStepper,
    Skeleton,
    SkeletonText,
    Icon,
    Flex
} from '@chakra-ui/react';
import { FaTrash } from 'react-icons/fa';

const CartItemList = ({
                          items,
                          isLoading,
                          selectedItems,
                          onSelectItem,
                          onSelectAll,
                          onUpdateQuantity,
                          onDeleteItems
                      }) => {
    const isAllSelected = items.length > 0 && selectedItems.length === items.length;

    if (isLoading) {
        return <CartItemSkeleton />;
    }

    if (items.length === 0) {
        return (
            <Box textAlign="center" py={10}>
                <Text fontSize="lg">Giỏ hàng của bạn đang trống</Text>
                <Button colorScheme="brand" mt={4} as="a" href="/">
                    Tiếp tục mua sắm
                </Button>
            </Box>
        );
    }

    // Format currency function
    const formatPrice = (price) => {
        return new Intl.NumberFormat('vi-VN', {
            style: 'currency',
            currency: 'VND',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        }).format(price);
    };

    return (
        <VStack align="stretch" spacing={4} bg="white" p={4} borderRadius="md" borderWidth="1px">
            {/* Header */}
            <HStack justify="space-between">
                <Checkbox
                    isChecked={isAllSelected}
                    onChange={(e) => onSelectAll(e.target.checked)}
                >
                    <Text fontWeight="bold">Chọn tất cả ({items.length} sản phẩm)</Text>
                </Checkbox>
                {selectedItems.length > 0 && (
                    <Button
                        leftIcon={<Icon as={FaTrash} />}
                        colorScheme="red"
                        variant="ghost"
                        size="sm"
                        onClick={() => onDeleteItems(selectedItems)}
                    >
                        Xóa
                    </Button>
                )}
            </HStack>

            <Divider />

            {/* Item list */}
            {items.map((item) => (
                <Box key={item.cart_item_id}>
                    <HStack spacing={4} align="center" py={2}>
                        <Checkbox
                            isChecked={selectedItems.includes(item.cart_item_id)}
                            onChange={(e) => onSelectItem(item.cart_item_id, e.target.checked)}
                        />

                        <Image
                            src={item.product_variant_thumbnail}
                            alt={item.product_name}
                            boxSize="80px"
                            objectFit="cover"
                            borderRadius="md"
                            fallbackSrc="https://via.placeholder.com/80"
                        />

                        <VStack flex="1" align="start" spacing={1}>
                            <Text fontWeight="medium" noOfLines={2}>{item.product_name}</Text>
                            <Text fontSize="sm" color="gray.500">Variant: {item.product_variant_id}</Text>
                        </VStack>

                        <HStack spacing={6}>
                            <Text fontWeight="bold" color={item.discount_price > 0 ? "red.500" : "gray.700"}>
                                {formatPrice(item.discount_price > 0 ? item.discount_price : item.price)}
                            </Text>

                            <NumberInput
                                min={1}
                                max={99}
                                value={item.quantity}
                                onChange={(valueString) => {
                                    const value = parseInt(valueString);
                                    if (!isNaN(value)) {
                                        onUpdateQuantity(item.cart_item_id, value, item.product_variant_id);
                                    }
                                }}
                                size="sm"
                                maxW="100px"
                            >
                                <NumberInputField />
                                <NumberInputStepper>
                                    <NumberIncrementStepper />
                                    <NumberDecrementStepper />
                                </NumberInputStepper>
                            </NumberInput>

                            <Button
                                size="sm"
                                variant="ghost"
                                colorScheme="red"
                                onClick={() => onDeleteItems([item.cart_item_id])}
                            >
                                <Icon as={FaTrash} />
                            </Button>
                        </HStack>
                    </HStack>
                    <Divider />
                </Box>
            ))}
        </VStack>
    );
};

// Skeleton cho trạng thái loading
const CartItemSkeleton = () => (
    <VStack align="stretch" spacing={4} bg="white" p={4} borderRadius="md" borderWidth="1px">
        <Flex justify="space-between" align="center">
            <Skeleton height="20px" width="200px" />
            <Skeleton height="20px" width="50px" />
        </Flex>
        <Divider />
        {[1, 2, 3].map((item) => (
            <Box key={item}>
                <HStack spacing={4} py={2}>
                    <Skeleton height="20px" width="20px" />
                    <Skeleton height="80px" width="80px" />
                    <VStack flex="1" align="start" spacing={1}>
                        <Skeleton height="20px" width="80%" />
                        <Skeleton height="16px" width="40%" />
                    </VStack>
                    <HStack spacing={6}>
                        <Skeleton height="20px" width="80px" />
                        <Skeleton height="30px" width="100px" />
                        <Skeleton height="20px" width="20px" />
                    </HStack>
                </HStack>
                <Divider />
            </Box>
        ))}
    </VStack>
);

export default CartItemList;