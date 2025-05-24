import React, { useState, useEffect, useRef, useCallback } from 'react';
import {
    Box,
    VStack,
    HStack,
    Text,
    Icon,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalBody,
    ModalCloseButton,
    useDisclosure,
    Radio,
    RadioGroup,
    Badge,
    Spinner,
    Center,
    useToast,
    Skeleton,
    SkeletonText,
    Button,
    Flex
} from '@chakra-ui/react';
import { FiMapPin, FiChevronRight } from 'react-icons/fi';
import userMeService from '../../services/userMeService';

const AddressSelector = ({ selectedAddress, onAddressSelect, orderTotal }) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [addresses, setAddresses] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isLoadingMore, setIsLoadingMore] = useState(false);
    const [selectedAddressId, setSelectedAddressId] = useState(selectedAddress?.id || '');
    const [currentPage, setCurrentPage] = useState(1);
    const [hasMore, setHasMore] = useState(true);
    const scrollContainerRef = useRef(null);
    const loadingRef = useRef(null);
    const toast = useToast();

    const ITEMS_PER_PAGE = 10;

    useEffect(() => {
        if (isOpen) {
            resetAndFetchAddresses();
        }
    }, [isOpen]);

    useEffect(() => {
        setSelectedAddressId(selectedAddress?.id || '');
    }, [selectedAddress]);

    const resetAndFetchAddresses = () => {
        setAddresses([]);
        setCurrentPage(1);
        setHasMore(true);
        fetchAddresses(1, true);
    };

    const fetchAddresses = async (page = 1, isInitial = false) => {
        if (isInitial) {
            setIsLoading(true);
        } else {
            setIsLoadingMore(true);
        }

        try {
            const response = await userMeService.getAddresses({
                limit: ITEMS_PER_PAGE,
                page: page
            });

            if (response && response.data) {
                if (isInitial) {
                    setAddresses(response.data);
                } else {
                    setAddresses(prev => [...prev, ...response.data]);
                }

                const hasMoreData = response.metadata?.pagination?.has_next ||
                    (response.data.length === ITEMS_PER_PAGE);
                setHasMore(hasMoreData);

                if (hasMoreData) {
                    setCurrentPage(page + 1);
                }
            }
        } catch (error) {
            console.error('Error fetching addresses:', error);
            toast({
                title: 'Lỗi tải địa chỉ',
                description: 'Không thể tải danh sách địa chỉ',
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

    const loadMoreAddresses = useCallback(() => {
        if (!isLoadingMore && hasMore) {
            fetchAddresses(currentPage, false);
        }
    }, [currentPage, hasMore, isLoadingMore]);

    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                const target = entries[0];
                if (target.isIntersecting && hasMore && !isLoadingMore) {
                    loadMoreAddresses();
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
    }, [loadMoreAddresses, hasMore, isLoadingMore]);

    const handleAddressSelect = (addressId) => {
        setSelectedAddressId(addressId);
    };

    const handleApplyAddress = () => {
        const selectedAddressData = addresses.find(a => a.id.toString() === selectedAddressId);
        if (selectedAddressData) {
            onAddressSelect(selectedAddressData);
        }
        onClose();
    };

    return (
        <>
            <Flex
                align="center"
                justify="space-between"
                p={4}
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
                    <Icon as={FiMapPin} color="blue.500" boxSize={5} />
                    <VStack align="start" spacing={0}>
                        <Text fontSize="sm" fontWeight="medium" color="blue.500">
                            Địa Chỉ Nhận Hàng
                        </Text>
                        {selectedAddress ? (
                            <VStack align="start" spacing={1}>
                                <Text fontSize="sm" fontWeight="semibold">
                                    {selectedAddress.recipient_name} | {selectedAddress.phone}
                                </Text>
                                <Text fontSize="xs" color="gray.600">
                                    {selectedAddress.address_line}, {selectedAddress.ward}, {selectedAddress.district}, {selectedAddress.province}
                                </Text>
                                {selectedAddress.is_default && (
                                    <Badge colorScheme="red" size="sm">Mặc định</Badge>
                                )}
                            </VStack>
                        ) : (
                            <Text fontSize="xs" color="gray.500">
                                Chọn địa chỉ nhận hàng
                            </Text>
                        )}
                    </VStack>
                </HStack>
                <Icon as={FiChevronRight} color="gray.400" />
            </Flex>

            <Modal isOpen={isOpen} onClose={onClose} size="lg" scrollBehavior="inside">
                <ModalOverlay />
                <ModalContent maxH="90vh">
                    <ModalHeader borderBottom="1px" borderColor="gray.200">
                        <HStack>
                            <Icon as={FiMapPin} color="blue.500" />
                            <Text>Chọn Địa Chỉ Nhận Hàng</Text>
                        </HStack>
                    </ModalHeader>
                    <ModalCloseButton />

                    <ModalBody p={0} ref={scrollContainerRef} maxH="60vh" overflowY="auto">
                        <VStack spacing={0} align="stretch">
                            <Box p={4}>
                                <Text fontSize="sm" fontWeight="medium" mb={3}>
                                    Địa chỉ của tôi
                                </Text>

                                {isLoading ? (
                                    <VStack spacing={3}>
                                        {[1, 2, 3].map(i => (
                                            <Box key={i} p={3} borderWidth="1px" borderRadius="md" w="100%">
                                                <HStack spacing={3}>
                                                    <Skeleton height="16px" width="16px" />
                                                    <VStack align="start" spacing={1} flex="1">
                                                        <Skeleton height="16px" width="60%" />
                                                        <SkeletonText noOfLines={2} spacing="2" width="80%" />
                                                    </VStack>
                                                </HStack>
                                            </Box>
                                        ))}
                                    </VStack>
                                ) : addresses.length === 0 && !isLoading ? (
                                    <Box textAlign="center" py={8}>
                                        <Icon as={FiMapPin} boxSize={12} color="gray.300" mb={2} />
                                        <Text color="gray.500">Chưa có địa chỉ nào</Text>
                                        <Text fontSize="sm" color="gray.400">
                                            Thêm địa chỉ để tiếp tục mua hàng
                                        </Text>
                                    </Box>
                                ) : (
                                    <RadioGroup value={selectedAddressId} onChange={handleAddressSelect}>
                                        <VStack spacing={3} align="stretch">
                                            {addresses.map((address) => (
                                                <Box
                                                    key={address.id}
                                                    p={3}
                                                    borderWidth="1px"
                                                    borderColor={selectedAddressId == address.id ? "blue.300" : "gray.200"}
                                                    borderRadius="md"
                                                    bg={selectedAddressId == address.id ? "blue.50" : "white"}
                                                    cursor="pointer"
                                                    onClick={() => handleAddressSelect(address.id.toString())}
                                                    _hover={{ borderColor: "blue.200", bg: "blue.25" }}
                                                >
                                                    <HStack spacing={3} align="start">
                                                        <Radio
                                                            value={address.id.toString()}
                                                            colorScheme="blue"
                                                        />

                                                        <VStack align="start" spacing={1} flex="1">
                                                            <HStack spacing={2} flexWrap="wrap">
                                                                <Text fontSize="sm" fontWeight="bold">
                                                                    {address.recipient_name}
                                                                </Text>
                                                                <Text fontSize="sm" color="gray.600">
                                                                    {address.phone}
                                                                </Text>
                                                                {address.is_default && (
                                                                    <Badge colorScheme="red" size="sm">
                                                                        Mặc định
                                                                    </Badge>
                                                                )}
                                                                <Badge colorScheme="blue" size="sm">
                                                                    {address.address_type?.name || 'Nhà'}
                                                                </Badge>
                                                            </HStack>

                                                            <Text fontSize="xs" color="gray.600">
                                                                {address.address_line}
                                                            </Text>
                                                            <Text fontSize="xs" color="gray.500">
                                                                {address.ward}, {address.district}, {address.province}
                                                            </Text>
                                                        </VStack>
                                                    </HStack>
                                                </Box>
                                            ))}

                                            {hasMore && (
                                                <Box ref={loadingRef} py={4}>
                                                    {isLoadingMore ? (
                                                        <Center>
                                                            <VStack spacing={2}>
                                                                <Spinner size="sm" color="blue.500" />
                                                                <Text fontSize="xs" color="gray.500">
                                                                    Đang tải thêm địa chỉ...
                                                                </Text>
                                                            </VStack>
                                                        </Center>
                                                    ) : (
                                                        <Center>
                                                            <Button
                                                                variant="ghost"
                                                                size="sm"
                                                                onClick={loadMoreAddresses}
                                                                color="blue.500"
                                                            >
                                                                Tải thêm địa chỉ
                                                            </Button>
                                                        </Center>
                                                    )}
                                                </Box>
                                            )}

                                            {!hasMore && addresses.length > 0 && (
                                                <Center py={4}>
                                                    <Text fontSize="xs" color="gray.500">
                                                        Đã hiển thị tất cả địa chỉ
                                                    </Text>
                                                </Center>
                                            )}
                                        </VStack>
                                    </RadioGroup>
                                )}
                            </Box>
                        </VStack>
                    </ModalBody>

                    {!isLoading && addresses.length > 0 && (
                        <Box p={4} borderTop="1px" borderColor="gray.200" bg="white">
                            <HStack spacing={2}>
                                <Button
                                    variant="outline"
                                    onClick={onClose}
                                    flex="1"
                                >
                                    Hủy
                                </Button>
                                <Button
                                    colorScheme="blue"
                                    onClick={handleApplyAddress}
                                    flex="1"
                                    isDisabled={!selectedAddressId}
                                >
                                    Xác nhận
                                </Button>
                            </HStack>
                        </Box>
                    )}
                </ModalContent>
            </Modal>
        </>
    );
};

export default AddressSelector;