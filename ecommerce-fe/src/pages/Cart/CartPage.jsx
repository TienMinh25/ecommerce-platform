import React, { useState } from 'react';
import {
    Box,
    Container,
    Checkbox,
    Text,
    Image,
    HStack,
    VStack,
    Flex,
    Button,
    Icon,
    Input,
    IconButton,
    Grid,
    GridItem,
    Heading,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalBody,
    ModalFooter,
    useDisclosure,
    SkeletonText,
    Skeleton,
} from '@chakra-ui/react';
import { FaTrash, FaMinus, FaPlus } from 'react-icons/fa';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import PageTitle from "../PageTitle.jsx";
import { useCart } from "../../context/CartContext.jsx";
import VoucherSelector from "../../components/cart/VoucherSelector.jsx";

const CartPage = () => {
    const {
        cartItems,
        selectedItems,
        isLoading,
        updateCartItem,
        deleteCartItems,
        toggleSelectItem,
        toggleSelectAll,
        calculateTotal,
        formatPrice
    } = useCart();

    const [itemToDelete, setItemToDelete] = useState(null);
    const [selectedVoucher, setSelectedVoucher] = useState(null);
    const { isOpen, onOpen, onClose } = useDisclosure();
    const navigate = useNavigate();

    // Calculate cart total for selected items
    const cartTotal = calculateTotal();

    // Calculate discount from voucher
    const calculateVoucherDiscount = () => {
        if (!selectedVoucher || cartTotal === 0) return 0;

        let discount = 0;
        if (selectedVoucher.discount_type === 'percentage') {
            discount = (cartTotal * selectedVoucher.discount_value) / 100;
        } else {
            discount = selectedVoucher.discount_value;
        }

        // Apply maximum discount limit
        if (selectedVoucher.maximum_discount_amount) {
            discount = Math.min(discount, selectedVoucher.maximum_discount_amount);
        }

        return discount;
    };

    // Calculate final total after voucher discount
    const finalTotal = cartTotal - calculateVoucherDiscount();

    const openDeleteConfirmation = (cartItemId, productName) => {
        setItemToDelete({ id: cartItemId, name: productName });
        onOpen();
    };

    const handleDeleteItem = async () => {
        if (!itemToDelete) return;
        await deleteCartItems([itemToDelete.id]);
        onClose();
    };

    const handleDeleteSelected = async () => {
        if (selectedItems.length === 0) return;
        await deleteCartItems(selectedItems);
    };

    const handleVoucherSelect = (voucher) => {
        setSelectedVoucher(voucher);
    };

    const handleCheckout = () => {
        if (selectedItems.length === 0) {
            return;
        }
        // Pass voucher data to checkout if needed
        navigate('/checkout', {
            state: {
                selectedVoucher,
                voucherDiscount: calculateVoucherDiscount(),
                finalTotal
            }
        });
    };

    // Loading skeleton
    if (isLoading) {
        return (
            <Container maxW="container.xl" py={6}>
                <PageTitle title="Giỏ hàng" />
                <Box bg="white" borderRadius="sm" boxShadow="sm" overflow="hidden">
                    <HStack justify="space-between" p={4} borderBottomWidth="1px" bgColor="gray.50">
                        <Skeleton height="24px" width="200px" />
                        <Skeleton height="24px" width="100px" />
                    </HStack>

                    {[1, 2, 3].map(i => (
                        <Box key={i} p={4} borderBottomWidth="1px">
                            <HStack spacing={4} align="flex-start">
                                <Skeleton height="24px" width="24px" />
                                <Skeleton height="80px" width="80px" />
                                <VStack align="start" spacing={2} flex="1">
                                    <Skeleton height="20px" width="80%" />
                                    <Skeleton height="16px" width="50%" />
                                </VStack>
                                <SkeletonText mt={2} noOfLines={2} width="100px" />
                                <Skeleton height="35px" width="120px" />
                                <Skeleton height="24px" width="80px" />
                            </HStack>
                        </Box>
                    ))}
                </Box>
            </Container>
        );
    }

    // Empty cart
    if (cartItems.length === 0) {
        return (
            <Container maxW="container.xl" py={6}>
                <PageTitle title="Giỏ hàng" />
                <Box
                    bg="white"
                    p={8}
                    borderRadius="sm"
                    textAlign="center"
                    boxShadow="sm"
                >
                    <Heading size="md" mb={6} color="gray.600">Giỏ hàng của bạn đang trống</Heading>
                    <Button
                        as={RouterLink}
                        to="/"
                        colorScheme="red"
                        size="lg"
                    >
                        Tiếp tục mua sắm
                    </Button>
                </Box>
            </Container>
        );
    }

    return (
        <Container maxW="container.xl" py={6}>
            <PageTitle title="Giỏ hàng" />

            {/* Header - Desktop */}
            <Grid
                templateColumns="50px 3fr 1fr 1fr 1fr 1fr"
                gap={2}
                alignItems="center"
                p={3}
                bg="white"
                borderWidth="1px"
                borderColor="gray.200"
                mb="3px"
                display={{ base: 'none', md: 'grid' }}
            >
                <GridItem>
                    <Checkbox
                        isChecked={selectedItems.length === cartItems.length && cartItems.length > 0}
                        onChange={toggleSelectAll}
                        colorScheme="red"
                    />
                </GridItem>
                <GridItem>
                    <Text fontWeight="medium">Sản Phẩm</Text>
                </GridItem>
                <GridItem>
                    <Text fontWeight="medium">Đơn Giá</Text>
                </GridItem>
                <GridItem>
                    <Text fontWeight="medium">Số Lượng</Text>
                </GridItem>
                <GridItem>
                    <Text fontWeight="medium">Số Tiền</Text>
                </GridItem>
                <GridItem>
                    <Text fontWeight="medium">Thao Tác</Text>
                </GridItem>
            </Grid>

            {/* Header - Mobile */}
            <Box
                display={{ base: 'flex', md: 'none' }}
                justifyContent="space-between"
                alignItems="center"
                p={3}
                bg="white"
                borderWidth="1px"
                borderColor="gray.200"
                mb="3px"
            >
                <HStack>
                    <Checkbox
                        isChecked={selectedItems.length === cartItems.length && cartItems.length > 0}
                        onChange={toggleSelectAll}
                        colorScheme="red"
                    />
                    <Text fontWeight="medium">Chọn tất cả ({cartItems.length})</Text>
                </HStack>
                {selectedItems.length > 0 && (
                    <Button
                        leftIcon={<Icon as={FaTrash} />}
                        colorScheme="red"
                        variant="ghost"
                        size="sm"
                        onClick={handleDeleteSelected}
                    >
                        Xóa
                    </Button>
                )}
            </Box>

            {/* Cart items */}
            {cartItems.map((item) => (
                <Box
                    key={item.cart_item_id}
                    bg="white"
                    mb="3px"
                    borderWidth="1px"
                    borderColor="gray.200"
                >
                    {/* Desktop view */}
                    <Grid
                        templateColumns="50px 3fr 1fr 1fr 1fr 1fr"
                        gap={2}
                        alignItems="center"
                        p={3}
                        display={{ base: 'none', md: 'grid' }}
                    >
                        <GridItem>
                            <Checkbox
                                isChecked={selectedItems.includes(item.cart_item_id)}
                                onChange={() => toggleSelectItem(item.cart_item_id)}
                                colorScheme="red"
                            />
                        </GridItem>

                        <GridItem>
                            <HStack>
                                <Image
                                    src={item.product_variant_thumbnail || 'https://via.placeholder.com/80'}
                                    alt={item.product_name}
                                    boxSize="80px"
                                    objectFit="cover"
                                    borderRadius="sm"
                                />
                                <VStack align="flex-start" spacing={1}>
                                    <Text fontWeight="medium" noOfLines={2}>{item.product_name}</Text>
                                    <HStack>
                                        <Text fontSize="sm" color="gray.500">Phân Loại Hàng:</Text>
                                        <Text fontSize="sm">
                                            {item.variant_name ?
                                                `${item.variant_name}` :
                                                item.attribute_values && item.attribute_values.length > 0 ?
                                                    item.attribute_values.map(attr => `${attr.attribute_value}`).join(', ') :
                                                    'Mặc định'}
                                        </Text>
                                    </HStack>
                                </VStack>
                            </HStack>
                        </GridItem>

                        <GridItem>
                            {item.discount_price > 0 ? (
                                <VStack align="flex-start" spacing={1}>
                                    <Text as="s" color="gray.500">
                                        {formatPrice(item.price)}
                                    </Text>
                                    <Text fontWeight="medium" color="red.500">
                                        {formatPrice(item.discount_price)}
                                    </Text>
                                </VStack>
                            ) : (
                                <Text fontWeight="medium">
                                    {formatPrice(item.price)}
                                </Text>
                            )}
                        </GridItem>

                        <GridItem>
                            <HStack w="120px">
                                <IconButton
                                    icon={<FaMinus />}
                                    aria-label="Decrease quantity"
                                    size="sm"
                                    borderRadius="md"
                                    colorScheme="blue"
                                    variant="outline"
                                    onClick={() => updateCartItem(item.cart_item_id, {
                                        product_variant_id: item.product_variant_id,
                                        quantity: item.quantity - 1
                                    })}
                                    _hover={{ bg: 'blue.50' }}
                                />
                                <Input
                                    value={item.quantity}
                                    readOnly
                                    textAlign="center"
                                    minW="40px"
                                    borderRadius="md"
                                    borderColor="gray.300"
                                    p={1}
                                />
                                <IconButton
                                    icon={<FaPlus />}
                                    aria-label="Increase quantity"
                                    size="sm"
                                    borderRadius="md"
                                    colorScheme="blue"
                                    variant="outline"
                                    onClick={() => updateCartItem(item.cart_item_id, {
                                        product_variant_id: item.product_variant_id,
                                        quantity: item.quantity + 1
                                    })}
                                    _hover={{ bg: 'blue.50' }}
                                />
                            </HStack>
                        </GridItem>

                        <GridItem>
                            <Text fontWeight="bold" color="red.500">
                                {formatPrice((item.discount_price > 0 ? item.discount_price : item.price) * item.quantity)}
                            </Text>
                        </GridItem>

                        <GridItem>
                            <Button
                                colorScheme="blue"
                                variant="ghost"
                                size="sm"
                                onClick={() => openDeleteConfirmation(item.cart_item_id, item.product_name)}
                                color="blue.500"
                            >
                                Xóa
                            </Button>
                        </GridItem>
                    </Grid>

                    {/* Mobile view */}
                    <Box display={{ base: 'block', md: 'none' }} p={3}>
                        <HStack align="flex-start" spacing={3}>
                            <Checkbox
                                isChecked={selectedItems.includes(item.cart_item_id)}
                                onChange={() => toggleSelectItem(item.cart_item_id)}
                                colorScheme="red"
                                mt={1}
                            />

                            <Image
                                src={item.product_variant_thumbnail || 'https://via.placeholder.com/80'}
                                alt={item.product_name}
                                boxSize="80px"
                                objectFit="cover"
                                borderRadius="sm"
                            />

                            <VStack align="flex-start" spacing={2} flex="1">
                                <Text fontWeight="medium" noOfLines={2}>{item.product_name}</Text>
                                <HStack>
                                    <Text fontSize="sm" color="gray.500">Phân Loại Hàng:</Text>
                                    <Text fontSize="sm">
                                        {item.variant_name ?
                                            `${item.variant_name}` :
                                            item.attribute_values && item.attribute_values.length > 0 ?
                                                item.attribute_values.map(attr => `${attr.attribute_value}`).join(', ') :
                                                'Mặc định'}
                                    </Text>
                                </HStack>

                                <Flex justify="space-between" w="100%" mt={2}>
                                    <HStack spacing={1}>
                                        <IconButton
                                            icon={<FaMinus />}
                                            aria-label="Decrease quantity"
                                            size="sm"
                                            borderRadius="md"
                                            colorScheme="blue"
                                            variant="outline"
                                            onClick={() => updateCartItem(item.cart_item_id, {
                                                product_variant_id: item.product_variant_id,
                                                quantity: item.quantity - 1
                                            })}
                                            _hover={{ bg: 'blue.50' }}
                                        />
                                        <Input
                                            value={item.quantity}
                                            readOnly
                                            textAlign="center"
                                            w="40px"
                                            borderRadius="md"
                                            p={1}
                                        />
                                        <IconButton
                                            icon={<FaPlus />}
                                            aria-label="Increase quantity"
                                            size="sm"
                                            borderRadius="md"
                                            colorScheme="blue"
                                            variant="outline"
                                            onClick={() => updateCartItem(item.cart_item_id, {
                                                product_variant_id: item.product_variant_id,
                                                quantity: item.quantity + 1
                                            })}
                                            _hover={{ bg: 'blue.50' }}
                                        />
                                    </HStack>

                                    <Button
                                        colorScheme="blue"
                                        variant="ghost"
                                        size="sm"
                                        onClick={() => openDeleteConfirmation(item.cart_item_id, item.product_name)}
                                        color="blue.500"
                                    >
                                        Xóa
                                    </Button>
                                </Flex>

                                <Flex justify="space-between" w="100%" mt={1}>
                                    {item.discount_price > 0 ? (
                                        <VStack align="flex-start" spacing={0}>
                                            <Text as="s" color="gray.500" fontSize="sm">
                                                {formatPrice(item.price)}
                                            </Text>
                                            <Text fontWeight="medium" color="red.500">
                                                {formatPrice(item.discount_price)}
                                            </Text>
                                        </VStack>
                                    ) : (
                                        <Text fontWeight="medium">
                                            {formatPrice(item.price)}
                                        </Text>
                                    )}

                                    <Text fontWeight="bold" color="red.500">
                                        {formatPrice((item.discount_price > 0 ? item.discount_price : item.price) * item.quantity)}
                                    </Text>
                                </Flex>
                            </VStack>
                        </HStack>
                    </Box>
                </Box>
            ))}

            {/* Cart summary and checkout */}
            <Flex
                direction={{ base: 'column', md: 'row' }}
                justify="space-between"
                align="center"
                bg="white"
                p={4}
                borderWidth="1px"
                borderColor="gray.200"
                position="sticky"
                bottom={0}
                mt={6}
            >
                <HStack spacing={4}>
                    <Checkbox
                        isChecked={selectedItems.length === cartItems.length && cartItems.length > 0}
                        onChange={toggleSelectAll}
                        colorScheme="red"
                    />
                    <Text>Chọn Tất Cả ({cartItems.length})</Text>
                    <Button
                        colorScheme="blue"
                        variant="ghost"
                        size="sm"
                        onClick={handleDeleteSelected}
                        isDisabled={selectedItems.length === 0}
                        display={{ base: 'none', md: 'inline-flex' }}
                        color="blue.500"
                    >
                        Xóa
                    </Button>
                </HStack>

                <VStack spacing={3} align={{ base: 'stretch', md: 'flex-end' }} mt={{ base: 4, md: 0 }}>
                    {/* Voucher Section - Always visible */}
                    <Box w={{ base: '100%', md: '400px' }}>
                        <VoucherSelector
                            selectedVoucher={selectedVoucher}
                            onVoucherSelect={handleVoucherSelect}
                            cartTotal={cartTotal}
                        />
                    </Box>

                    {/* Show voucher discount if applied */}
                    {selectedVoucher && calculateVoucherDiscount() > 0 && (
                        <HStack spacing={2} justify={{ base: 'space-between', md: 'flex-end' }} w="100%">
                            <Text fontSize="sm" color="gray.600">
                                Voucher giảm:
                            </Text>
                            <Text fontSize="sm" color="green.600" fontWeight="medium">
                                -{formatPrice(calculateVoucherDiscount())}
                            </Text>
                        </HStack>
                    )}

                    <HStack spacing={6} w="100%" justify={{ base: 'space-between', md: 'flex-end' }}>
                        <Box textAlign={{ base: 'left', md: 'right' }}>
                            <Text>
                                Tổng thanh toán ({selectedItems.length} Sản phẩm):
                                <Text as="span" fontWeight="bold" color="red.500" fontSize="xl" ml={2}>
                                    {formatPrice(finalTotal)}
                                </Text>
                            </Text>
                            {selectedVoucher && calculateVoucherDiscount() > 0 && (
                                <Text fontSize="sm" color="gray.500" as="s">
                                    {formatPrice(cartTotal)}
                                </Text>
                            )}
                        </Box>

                        <Button
                            colorScheme="red"
                            size="lg"
                            isDisabled={selectedItems.length === 0}
                            onClick={handleCheckout}
                            px={8}
                        >
                            Mua Hàng
                        </Button>
                    </HStack>
                </VStack>
            </Flex>

            {/* Delete Confirmation Modal */}
            <Modal isOpen={isOpen} onClose={onClose} isCentered>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader color="red.500" fontSize="xl" textAlign="center">
                        Bạn chắc chắn muốn bỏ sản phẩm này?
                    </ModalHeader>

                    <ModalBody py={6}>
                        <Text textAlign="center">
                            {itemToDelete?.name}
                        </Text>
                    </ModalBody>

                    <ModalFooter justifyContent="center" pb={6}>
                        <Button
                            colorScheme="red"
                            mr={3}
                            onClick={handleDeleteItem}
                            w="150px"
                        >
                            Có
                        </Button>
                        <Button
                            variant="outline"
                            onClick={onClose}
                            w="150px"
                        >
                            Không
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </Container>
    );
};

export default CartPage;