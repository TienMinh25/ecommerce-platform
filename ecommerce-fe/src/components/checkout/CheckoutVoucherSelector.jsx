import React, { useState, useEffect, useRef, useCallback } from 'react';
import {
    Box,
    Button,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalBody,
    ModalCloseButton,
    VStack,
    HStack,
    Text,
    Badge,
    useDisclosure,
    Icon,
    Flex,
    useToast,
    Skeleton,
    SkeletonText,
    Radio,
    RadioGroup,
    Spinner,
    Center,
} from '@chakra-ui/react';
import { FiTag, FiPercent, FiDollarSign, FiClock, FiGift, FiChevronRight } from 'react-icons/fi';
import couponService from "../../services/couponService.js";

const CheckoutVoucherSelector = ({ selectedVoucher, onVoucherSelect, cartTotal }) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [vouchers, setVouchers] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isLoadingMore, setIsLoadingMore] = useState(false);
    const [selectedVoucherId, setSelectedVoucherId] = useState(selectedVoucher?.id || '');
    const [currentPage, setCurrentPage] = useState(1);
    const [hasMore, setHasMore] = useState(true);
    const toast = useToast();
    const scrollContainerRef = useRef(null);
    const loadingRef = useRef(null);

    const ITEMS_PER_PAGE = 10;

    // Fetch vouchers when modal opens
    useEffect(() => {
        if (isOpen) {
            resetAndFetchVouchers();
        }
    }, [isOpen]);

    // Update selected voucher when prop changes
    useEffect(() => {
        setSelectedVoucherId(selectedVoucher?.id || '');
    }, [selectedVoucher]);

    // Reset pagination and fetch initial vouchers
    const resetAndFetchVouchers = () => {
        setVouchers([]);
        setCurrentPage(1);
        setHasMore(true);
        fetchVouchers(1, true);
    };

    // Fetch vouchers with pagination
    const fetchVouchers = async (page = 1, isInitial = false) => {
        if (isInitial) {
            setIsLoading(true);
        } else {
            setIsLoadingMore(true);
        }

        try {
            const response = await couponService.getCouponsForClient({
                limit: ITEMS_PER_PAGE,
                page: page
            });

            if (response && response.data) {
                if (isInitial) {
                    setVouchers(response.data);
                } else {
                    setVouchers(prev => [...prev, ...response.data]);
                }

                // Check if there are more pages
                const hasMoreData = response.metadata?.pagination?.has_next ||
                    (response.data.length === ITEMS_PER_PAGE);
                setHasMore(hasMoreData);

                if (hasMoreData) {
                    setCurrentPage(page + 1);
                }
            }
        } catch (error) {
            console.error('Error fetching vouchers:', error);
            toast({
                title: 'Lỗi tải voucher',
                description: 'Không thể tải danh sách voucher',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            if (isInitial) {
                setIsLoading(false);
            } else {
                setIsLoadingMore(false);
            }
        }
    };

    // Load more vouchers
    const loadMoreVouchers = useCallback(() => {
        if (!isLoadingMore && hasMore) {
            fetchVouchers(currentPage, false);
        }
    }, [currentPage, hasMore, isLoadingMore]);

    // Intersection Observer for infinite scroll
    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                const target = entries[0];
                if (target.isIntersecting && hasMore && !isLoadingMore) {
                    loadMoreVouchers();
                }
            },
            {
                root: scrollContainerRef.current,
                rootMargin: '100px',
                threshold: 0.1,
            }
        );

        if (loadingRef.current) {
            observer.observe(loadingRef.current);
        }

        return () => {
            if (loadingRef.current) {
                observer.unobserve(loadingRef.current);
            }
        };
    }, [loadMoreVouchers, hasMore, isLoadingMore]);

    // Format currency
    const formatCurrency = (value) => {
        return new Intl.NumberFormat('vi-VN', {
            style: 'currency',
            currency: 'VND'
        }).format(value);
    };

    // Format date
    const formatDate = (dateString) => {
        try {
            const date = new Date(dateString);
            return date.toLocaleDateString('vi-VN');
        } catch (e) {
            return dateString;
        }
    };

    // Calculate discount amount based on cart total and voucher limits
    const calculateDiscount = (voucher) => {
        if (!voucher || cartTotal === 0) return 0;

        let discount = 0;
        if (voucher.discount_type === 'percentage') {
            discount = (cartTotal * voucher.discount_value) / 100;
        } else {
            discount = voucher.discount_value;
        }

        // Apply maximum discount limit
        if (voucher.maximum_discount_amount) {
            discount = Math.min(discount, voucher.maximum_discount_amount);
        }

        return discount;
    };

    // Get discount display text
    const getDiscountText = (voucher) => {
        if (voucher.discount_type === 'percentage') {
            return `${voucher.discount_value}%`;
        } else {
            return formatCurrency(voucher.discount_value);
        }
    };

    // Check if voucher is applicable based on minimum order amount
    const isVoucherApplicable = (voucher) => {
        return cartTotal >= (voucher.minimum_order_amount || 0);
    };

    // Handle voucher selection
    const handleVoucherSelect = (voucherId) => {
        setSelectedVoucherId(voucherId);
    };

    // Apply selected voucher
    const handleApplyVoucher = () => {
        const selectedVoucherData = vouchers.find(v => v.id.toString() === selectedVoucherId);
        if (selectedVoucherData) {
            onVoucherSelect(selectedVoucherData);
            toast({
                title: 'Áp dụng voucher thành công',
                description: `Bạn đã tiết kiệm ${formatCurrency(calculateDiscount(selectedVoucherData))}`,
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
        } else {
            onVoucherSelect(null);
        }
        onClose();
    };

    // Remove selected voucher
    const handleRemoveVoucher = () => {
        setSelectedVoucherId('');
        onVoucherSelect(null);
        onClose();
    };

    return (
        <>
            {/* Voucher Selector Button */}
            <Flex
                align="center"
                justify="space-between"
                p={3}
                bg="white"
                borderWidth="1px"
                borderColor="gray.200"
                borderRadius="md"
                cursor="pointer"
                onClick={onOpen}
                _hover={{ bg: "gray.50" }}
                transition="background-color 0.2s"
            >
                <HStack spacing={3}>
                    <Icon as={FiTag} color="orange.500" boxSize={5} />
                    <VStack align="start" spacing={0}>
                        <Text fontSize="sm" fontWeight="medium" color="orange.500">
                            Minh Plaza Voucher
                        </Text>
                        {selectedVoucher ? (
                            <Text fontSize="xs" color="gray.600">
                                Tiết kiệm {formatCurrency(calculateDiscount(selectedVoucher))}
                            </Text>
                        ) : (
                            <Text fontSize="xs" color="gray.500">
                                Chọn voucher
                            </Text>
                        )}
                    </VStack>
                </HStack>
                <Icon as={FiChevronRight} color="gray.400" />
            </Flex>

            {/* Voucher Selection Modal */}
            <Modal isOpen={isOpen} onClose={onClose} size="lg" scrollBehavior="inside">
                <ModalOverlay />
                <ModalContent maxH="90vh">
                    <ModalHeader borderBottom="1px" borderColor="gray.200">
                        <HStack>
                            <Icon as={FiGift} color="orange.500" />
                            <Text>Chọn Minh Plaza Voucher</Text>
                        </HStack>
                    </ModalHeader>
                    <ModalCloseButton />

                    <ModalBody p={0} ref={scrollContainerRef} maxH="60vh" overflowY="auto">
                        <VStack spacing={0} align="stretch">
                            {/* Available Vouchers */}
                            <Box p={4}>
                                <Text fontSize="sm" fontWeight="medium" mb={3}>
                                    Voucher khả dụng
                                </Text>

                                {isLoading ? (
                                    <VStack spacing={3}>
                                        {[1, 2, 3, 4, 5].map(i => (
                                            <Box key={i} p={3} borderWidth="1px" borderRadius="md" w="100%">
                                                <HStack spacing={3}>
                                                    <Skeleton height="40px" width="40px" />
                                                    <VStack align="start" spacing={1} flex="1">
                                                        <Skeleton height="16px" width="60%" />
                                                        <SkeletonText noOfLines={2} spacing="2" width="80%" />
                                                    </VStack>
                                                </HStack>
                                            </Box>
                                        ))}
                                    </VStack>
                                ) : vouchers.length === 0 && !isLoading ? (
                                    <Box textAlign="center" py={8}>
                                        <Icon as={FiTag} boxSize={12} color="gray.300" mb={2} />
                                        <Text color="gray.500">Không có voucher khả dụng</Text>
                                        <Text fontSize="sm" color="gray.400">
                                            Hãy mua sắm thêm để mở khóa các voucher hấp dẫn
                                        </Text>
                                    </Box>
                                ) : (
                                    <RadioGroup value={selectedVoucherId} onChange={handleVoucherSelect}>
                                        <VStack spacing={3} align="stretch">
                                            {vouchers.map((voucher) => {
                                                const discount = calculateDiscount(voucher);
                                                const applicable = isVoucherApplicable(voucher);

                                                return (
                                                    <Box
                                                        key={voucher.id}
                                                        p={3}
                                                        borderWidth="1px"
                                                        borderColor={selectedVoucherId == voucher.id ? "orange.300" : "gray.200"}
                                                        borderRadius="md"
                                                        bg={selectedVoucherId == voucher.id ? "orange.50" : "white"}
                                                        opacity={applicable ? 1 : 0.6}
                                                        cursor={applicable ? "pointer" : "not-allowed"}
                                                        onClick={() => applicable && handleVoucherSelect(voucher.id.toString())}
                                                        _hover={applicable ? { borderColor: "orange.200", bg: "orange.25" } : {}}
                                                    >
                                                        <HStack spacing={3} align="start">
                                                            <Radio
                                                                value={voucher.id.toString()}
                                                                colorScheme="orange"
                                                                isDisabled={!applicable}
                                                            />

                                                            <Box
                                                                bg="orange.500"
                                                                color="white"
                                                                p={2}
                                                                borderRadius="md"
                                                                minW="40px"
                                                                textAlign="center"
                                                            >
                                                                <Icon
                                                                    as={voucher.discount_type === 'percentage' ? FiPercent : FiDollarSign}
                                                                    boxSize={4}
                                                                />
                                                            </Box>

                                                            <VStack align="start" spacing={1} flex="1">
                                                                <HStack spacing={2}>
                                                                    <Text fontSize="sm" fontWeight="bold">
                                                                        {voucher.name}
                                                                    </Text>
                                                                    <Badge colorScheme="orange" size="sm">
                                                                        {getDiscountText(voucher)}
                                                                    </Badge>
                                                                </HStack>

                                                                {voucher.maximum_discount_amount && (
                                                                    <Text fontSize="xs" color="gray.600">
                                                                        Giảm tối đa {formatCurrency(voucher.maximum_discount_amount)}
                                                                    </Text>
                                                                )}

                                                                {voucher.minimum_order_amount > 0 && (
                                                                    <Text fontSize="xs" color="gray.500">
                                                                        Đơn tối thiểu {formatCurrency(voucher.minimum_order_amount)}
                                                                    </Text>
                                                                )}

                                                                {voucher.end_date && (
                                                                    <HStack spacing={1} fontSize="xs" color="gray.500">
                                                                        <Icon as={FiClock} />
                                                                        <Text>HSD: {formatDate(voucher.end_date)}</Text>
                                                                    </HStack>
                                                                )}

                                                                {applicable && discount > 0 && (
                                                                    <Text fontSize="xs" color="green.600" fontWeight="medium">
                                                                        Tiết kiệm: {formatCurrency(discount)}
                                                                    </Text>
                                                                )}

                                                                {!applicable && voucher.minimum_order_amount > cartTotal && (
                                                                    <Text fontSize="xs" color="red.500">
                                                                        Cần mua thêm {formatCurrency(voucher.minimum_order_amount - cartTotal)}
                                                                    </Text>
                                                                )}
                                                            </VStack>
                                                        </HStack>
                                                    </Box>
                                                );
                                            })}

                                            {/* Loading more indicator */}
                                            {hasMore && (
                                                <Box ref={loadingRef} py={4}>
                                                    {isLoadingMore ? (
                                                        <Center>
                                                            <VStack spacing={2}>
                                                                <Spinner size="sm" color="orange.500" />
                                                                <Text fontSize="xs" color="gray.500">
                                                                    Đang tải thêm voucher...
                                                                </Text>
                                                            </VStack>
                                                        </Center>
                                                    ) : (
                                                        <Center>
                                                            <Button
                                                                variant="ghost"
                                                                size="sm"
                                                                onClick={loadMoreVouchers}
                                                                color="orange.500"
                                                            >
                                                                Tải thêm voucher
                                                            </Button>
                                                        </Center>
                                                    )}
                                                </Box>
                                            )}

                                            {/* End of list indicator */}
                                            {!hasMore && vouchers.length > 0 && (
                                                <Center py={4}>
                                                    <Text fontSize="xs" color="gray.500">
                                                        Đã hiển thị tất cả voucher
                                                    </Text>
                                                </Center>
                                            )}
                                        </VStack>
                                    </RadioGroup>
                                )}
                            </Box>
                        </VStack>
                    </ModalBody>

                    {!isLoading && vouchers.length > 0 && (
                        <Box p={4} borderTop="1px" borderColor="gray.200" bg="white">
                            <HStack spacing={2}>
                                <Button
                                    variant="outline"
                                    onClick={handleRemoveVoucher}
                                    flex="1"
                                >
                                    Trở lại
                                </Button>
                                <Button
                                    colorScheme="orange"
                                    onClick={handleApplyVoucher}
                                    flex="1"
                                    isDisabled={!selectedVoucherId}
                                >
                                    OK
                                </Button>
                            </HStack>
                        </Box>
                    )}
                </ModalContent>
            </Modal>
        </>
    );
};

export default CheckoutVoucherSelector;