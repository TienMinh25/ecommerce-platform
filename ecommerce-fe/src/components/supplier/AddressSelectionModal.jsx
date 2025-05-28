import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    VStack,
    HStack,
    Text,
    Button,
    Card,
    CardBody,
    Badge,
    Icon,
    Spinner,
    Center,
    Flex,
    Box,
    Radio,
    RadioGroup,
    Stack,
    useColorModeValue,
    useToast
} from '@chakra-ui/react';
import { useState, useCallback, useEffect } from 'react';
import { FaMapMarkerAlt, FaHome, FaBuilding } from 'react-icons/fa';
import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import userMeService from "../../services/userMeService.js";

const AddressSelectionModal = ({ isOpen, onClose, onSelectAddress, selectedAddressId }) => {
    const [addresses, setAddresses] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 5,
        totalItems: 0,
        totalPages: 0,
        hasNext: false,
        hasPrevious: false
    });
    const [tempSelectedId, setTempSelectedId] = useState(selectedAddressId);
    const toast = useToast();

    const cardBg = useColorModeValue('white', 'gray.800');
    const cardBorderColor = useColorModeValue('gray.200', 'gray.700');
    const selectedBg = useColorModeValue('blue.50', 'blue.900');
    const selectedBorderColor = useColorModeValue('blue.300', 'blue.600');

    const fetchAddresses = useCallback(async () => {
        setIsLoading(true);
        try {
            const response = await userMeService.getAddresses({
                page: pagination.page,
                limit: pagination.limit
            });

            setAddresses(response.data || []);

            if (response?.metadata?.pagination) {
                setPagination(prev => ({
                    ...prev,
                    totalItems: response.metadata.pagination.total_items,
                    totalPages: response.metadata.pagination.total_pages,
                    hasNext: response.metadata.pagination.has_next,
                    hasPrevious: response.metadata.pagination.has_previous
                }));
            }
        } catch (error) {
            console.error('Error fetching addresses:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể tải danh sách địa chỉ',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    }, [pagination.page, pagination.limit, toast]);

    useEffect(() => {
        if (isOpen) {
            fetchAddresses();
            setTempSelectedId(selectedAddressId);
        }
    }, [isOpen, fetchAddresses, selectedAddressId]);

    const handlePageChange = (newPage) => {
        setPagination(prev => ({
            ...prev,
            page: newPage
        }));
    };

    const handleSelectAddress = () => {
        const selectedAddress = addresses.find(addr => addr.id === tempSelectedId);
        if (selectedAddress) {
            onSelectAddress(selectedAddress);
            onClose();
        }
    };

    const getAddressTypeIcon = (addressType) => {
        switch (addressType?.toLowerCase()) {
            case 'home':
                return FaHome;
            case 'office':
            case 'work':
                return FaBuilding;
            default:
                return FaMapMarkerAlt;
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} size="xl" scrollBehavior="inside">
            <ModalOverlay />
            <ModalContent maxH="80vh">
                <ModalHeader>
                    <HStack spacing={3}>
                        <Icon as={FaMapMarkerAlt} color="blue.500" />
                        <Text>Chọn địa chỉ kinh doanh</Text>
                    </HStack>
                </ModalHeader>
                <ModalCloseButton />

                <ModalBody>
                    {isLoading ? (
                        <Center py={10}>
                            <VStack spacing={4}>
                                <Spinner size="xl" color="blue.500" />
                                <Text>Đang tải danh sách địa chỉ...</Text>
                            </VStack>
                        </Center>
                    ) : (
                        <VStack spacing={4} align="stretch">
                            {addresses.length === 0 ? (
                                <Center py={10}>
                                    <VStack spacing={4}>
                                        <Icon as={FaMapMarkerAlt} w={12} h={12} color="gray.400" />
                                        <VStack spacing={2}>
                                            <Text fontSize="lg" fontWeight="medium" color="gray.600">
                                                Chưa có địa chỉ nào
                                            </Text>
                                            <Text fontSize="sm" color="gray.500" textAlign="center">
                                                Bạn cần tạo địa chỉ trước khi đăng ký làm nhà cung cấp
                                            </Text>
                                            <Button
                                                colorScheme="blue"
                                                size="md"
                                                mt={4}
                                                onClick={() => {
                                                    window.open('/user/account/addresses', '_blank');
                                                }}
                                                leftIcon={<Icon as={FaMapMarkerAlt} />}
                                            >
                                                Tạo địa chỉ mới
                                            </Button>
                                        </VStack>
                                    </VStack>
                                </Center>
                            ) : (
                                <RadioGroup value={tempSelectedId?.toString()} onChange={(value) => setTempSelectedId(parseInt(value))}>
                                    <Stack spacing={3}>
                                        {addresses.map(address => (
                                            <Card
                                                key={address.id}
                                                variant="outline"
                                                cursor="pointer"
                                                onClick={() => setTempSelectedId(address.id)}
                                                bg={tempSelectedId === address.id ? selectedBg : cardBg}
                                                borderColor={tempSelectedId === address.id ? selectedBorderColor : cardBorderColor}
                                                borderWidth={tempSelectedId === address.id ? "2px" : "1px"}
                                                _hover={{
                                                    shadow: "md",
                                                    transform: "translateY(-1px)",
                                                    transition: "all 0.2s"
                                                }}
                                                transition="all 0.2s"
                                            >
                                                <CardBody p={4}>
                                                    <HStack spacing={4} align="flex-start">
                                                        <Radio
                                                            value={address.id.toString()}
                                                            colorScheme="blue"
                                                            size="lg"
                                                            mt={1}
                                                        />

                                                        <Box flex={1}>
                                                            <HStack spacing={3} mb={2} flexWrap="wrap">
                                                                <HStack spacing={2}>
                                                                    <Icon
                                                                        as={getAddressTypeIcon(address.address_type)}
                                                                        color="blue.500"
                                                                        boxSize={4}
                                                                    />
                                                                    <Text fontWeight="bold" fontSize="md">
                                                                        {address.recipient_name}
                                                                    </Text>
                                                                </HStack>

                                                                <Text color="gray.600" fontSize="sm">
                                                                    {address.phone}
                                                                </Text>

                                                                {address.is_default && (
                                                                    <Badge colorScheme="red" size="sm">
                                                                        Mặc định
                                                                    </Badge>
                                                                )}

                                                                <Badge colorScheme="blue" variant="subtle" size="sm">
                                                                    {address.address_type || 'Home'}
                                                                </Badge>
                                                            </HStack>

                                                            <VStack spacing={1} align="flex-start">
                                                                <Text fontSize="sm" color="gray.700">
                                                                    <Text as="span" fontWeight="medium">Địa chỉ:</Text> {address.street}
                                                                </Text>
                                                                <Text fontSize="sm" color="gray.600">
                                                                    {address.ward}, {address.district}, {address.province}
                                                                </Text>
                                                                <Text fontSize="sm" color="gray.600">
                                                                    {address.country}
                                                                </Text>
                                                            </VStack>
                                                        </Box>
                                                    </HStack>
                                                </CardBody>
                                            </Card>
                                        ))}
                                    </Stack>
                                </RadioGroup>
                            )}

                            {/* Pagination */}
                            {addresses.length > 0 && pagination.totalPages > 1 && (
                                <Flex justify="center" pt={4}>
                                    <HStack spacing={2}>
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() => handlePageChange(pagination.page - 1)}
                                            isDisabled={!pagination.hasPrevious || isLoading}
                                        >
                                            <ChevronLeftIcon />
                                        </Button>

                                        <Text fontSize="sm" px={3}>
                                            Trang {pagination.page} / {pagination.totalPages}
                                        </Text>

                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() => handlePageChange(pagination.page + 1)}
                                            isDisabled={!pagination.hasNext || isLoading}
                                        >
                                            <ChevronRightIcon />
                                        </Button>
                                    </HStack>
                                </Flex>
                            )}
                        </VStack>
                    )}
                </ModalBody>

                <ModalFooter>
                    <Button mr={3} onClick={onClose}>
                        Hủy
                    </Button>
                    <Button
                        colorScheme="blue"
                        onClick={handleSelectAddress}
                        isDisabled={!tempSelectedId || addresses.length === 0}
                    >
                        Chọn địa chỉ này
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default AddressSelectionModal;