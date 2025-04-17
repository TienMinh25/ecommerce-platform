import React, { useState } from 'react';
import {
    Box,
    Button,
    Flex,
    Heading,
    Text,
    Divider,
    useDisclosure,
    Badge,
    HStack,
    VStack,
    IconButton,
    useToast,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    FormControl,
    FormLabel,
    Input,
    Select,
    Checkbox,
    useColorModeValue
} from '@chakra-ui/react';
import { AddIcon, EditIcon, DeleteIcon } from '@chakra-ui/icons';

// Mock data for addresses
const mockAddresses = [
    {
        id: 1,
        fullName: 'Lê Văn Tiến Minh',
        phone: '(+84) 865 363 715',
        address: '25/2/29 Phú Minh, Minh Khai, Bắc Từ Liêm',
        ward: 'Phường Minh Khai',
        district: 'Quận Bắc Từ Liêm',
        city: 'Hà Nội',
        isDefault: true
    },
    {
        id: 2,
        fullName: 'Lê Văn Tiến Minh',
        phone: '(+84) 123 456 789',
        address: '123 Nguyễn Trãi',
        ward: 'Phường Thanh Xuân Nam',
        district: 'Quận Thanh Xuân',
        city: 'Hà Nội',
        isDefault: false
    }
];

const UserAddresses = () => {
    const toast = useToast();
    const [addresses, setAddresses] = useState(mockAddresses);
    const [currentAddress, setCurrentAddress] = useState(null);
    const [isDeleting, setIsDeleting] = useState(false);

    // Modal states
    const {
        isOpen: isAddressModalOpen,
        onOpen: onAddressModalOpen,
        onClose: onAddressModalClose
    } = useDisclosure();

    const {
        isOpen: isDeleteModalOpen,
        onOpen: onDeleteModalOpen,
        onClose: onDeleteModalClose
    } = useDisclosure();

    // Form state
    const [formData, setFormData] = useState({
        fullName: '',
        phone: '',
        address: '',
        ward: '',
        district: '',
        city: '',
        isDefault: false
    });

    // Colors
    const cardBg = useColorModeValue('white', 'gray.800');
    const cardBorderColor = useColorModeValue('gray.200', 'gray.700');
    const defaultBadgeBg = useColorModeValue('red.50', 'red.900');
    const defaultBadgeColor = useColorModeValue('red.600', 'red.200');

    // Handle opening address modal
    const handleAddAddress = () => {
        setCurrentAddress(null);
        setFormData({
            fullName: '',
            phone: '',
            address: '',
            ward: '',
            district: '',
            city: '',
            isDefault: false
        });
        onAddressModalOpen();
    };

    // Handle opening edit modal
    const handleEditAddress = (address) => {
        setCurrentAddress(address);
        setFormData({
            fullName: address.fullName,
            phone: address.phone,
            address: address.address,
            ward: address.ward,
            district: address.district,
            city: address.city,
            isDefault: address.isDefault
        });
        onAddressModalOpen();
    };

    // Handle opening delete modal
    const handleConfirmDelete = (address) => {
        setCurrentAddress(address);
        onDeleteModalOpen();
    };

    // Handle form input changes
    const handleInputChange = (e) => {
        const { name, value, checked } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: name === 'isDefault' ? checked : value
        }));
    };

    // Handle form submission
    const handleSubmitAddress = () => {
        if (currentAddress) {
            // Update existing address
            const updatedAddresses = addresses.map(addr =>
                addr.id === currentAddress.id
                    ? { ...formData, id: addr.id }
                    : (formData.isDefault && addr.isDefault ? { ...addr, isDefault: false } : addr)
            );
            setAddresses(updatedAddresses);

            toast({
                title: 'Address updated',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
        } else {
            // Add new address
            const newId = addresses.length > 0 ? Math.max(...addresses.map(a => a.id)) + 1 : 1;
            const newAddress = { ...formData, id: newId };

            // If the new address is set as default, unset default for other addresses
            const updatedAddresses = formData.isDefault
                ? addresses.map(addr => ({ ...addr, isDefault: false }))
                : [...addresses];

            setAddresses([...updatedAddresses, newAddress]);

            toast({
                title: 'Address added',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
        }

        onAddressModalClose();
    };

    // Handle address deletion
    const handleDeleteAddress = () => {
        setIsDeleting(true);

        try {
            // Filter out the address being deleted
            const updatedAddresses = addresses.filter(addr => addr.id !== currentAddress.id);
            setAddresses(updatedAddresses);

            toast({
                title: 'Address deleted',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });

            onDeleteModalClose();
        } catch (error) {
            toast({
                title: 'Error',
                description: 'Failed to delete address',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsDeleting(false);
        }
    };

    // Handle setting address as default
    const handleSetDefault = (addressId) => {
        const updatedAddresses = addresses.map(addr => ({
            ...addr,
            isDefault: addr.id === addressId
        }));

        setAddresses(updatedAddresses);

        toast({
            title: 'Default address updated',
            status: 'success',
            duration: 3000,
            isClosable: true,
        });
    };

    return (
        <Box>
            <Flex justify="space-between" align="center" mb={6}>
                <Box>
                    <Heading as="h1" size="lg">Địa Chỉ Của Tôi</Heading>
                    <Text color="gray.500" mt={1}>Quản lý địa chỉ nhận hàng</Text>
                </Box>
                <Button
                    leftIcon={<AddIcon />}
                    colorScheme="red"
                    onClick={handleAddAddress}
                >
                    Thêm địa chỉ mới
                </Button>
            </Flex>

            <Divider mb={6} />

            {/* Address List */}
            <VStack spacing={4} align="stretch">
                {addresses.length === 0 ? (
                    <Box textAlign="center" py={10}>
                        <Text fontSize="lg" color="gray.500">Bạn chưa có địa chỉ nào</Text>
                    </Box>
                ) : (
                    addresses.map(address => (
                        <Box
                            key={address.id}
                            p={4}
                            borderWidth="1px"
                            borderRadius="md"
                            borderColor={cardBorderColor}
                            bg={cardBg}
                            position="relative"
                        >
                            <HStack spacing={4} mb={3}>
                                <Text fontWeight="bold">{address.fullName}</Text>
                                <Text color="gray.600">{address.phone}</Text>
                                {address.isDefault && (
                                    <Badge
                                        bg={defaultBadgeBg}
                                        color={defaultBadgeColor}
                                        px={2}
                                        py={1}
                                        borderRadius="md"
                                        fontWeight="semibold"
                                    >
                                        Mặc định
                                    </Badge>
                                )}
                            </HStack>

                            <Text fontSize="sm" color="gray.700" mb={3}>
                                {address.address}, {address.ward}, {address.district}, {address.city}
                            </Text>

                            <Flex mt={2}>
                                <HStack spacing={3}>
                                    <Button
                                        size="sm"
                                        variant="outline"
                                        leftIcon={<EditIcon />}
                                        onClick={() => handleEditAddress(address)}
                                    >
                                        Sửa
                                    </Button>
                                    <Button
                                        size="sm"
                                        variant="outline"
                                        leftIcon={<DeleteIcon />}
                                        colorScheme="red"
                                        onClick={() => handleConfirmDelete(address)}
                                    >
                                        Xóa
                                    </Button>
                                    {!address.isDefault && (
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() => handleSetDefault(address.id)}
                                        >
                                            Thiết lập mặc định
                                        </Button>
                                    )}
                                </HStack>
                            </Flex>
                        </Box>
                    ))
                )}
            </VStack>

            {/* Add/Edit Address Modal */}
            <Modal isOpen={isAddressModalOpen} onClose={onAddressModalClose} size="lg">
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>
                        {currentAddress ? 'Cập nhật địa chỉ' : 'Thêm địa chỉ mới'}
                    </ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <VStack spacing={4} align="stretch">
                            <FormControl id="fullName" isRequired>
                                <FormLabel>Họ và tên</FormLabel>
                                <Input
                                    name="fullName"
                                    value={formData.fullName}
                                    onChange={handleInputChange}
                                    placeholder="Họ và tên người nhận"
                                />
                            </FormControl>

                            <FormControl id="phone" isRequired>
                                <FormLabel>Số điện thoại</FormLabel>
                                <Input
                                    name="phone"
                                    value={formData.phone}
                                    onChange={handleInputChange}
                                    placeholder="Số điện thoại"
                                />
                            </FormControl>

                            <HStack spacing={4}>
                                <FormControl id="city" isRequired>
                                    <FormLabel>Tỉnh/Thành phố</FormLabel>
                                    <Select
                                        name="city"
                                        value={formData.city}
                                        onChange={handleInputChange}
                                        placeholder="Chọn Tỉnh/Thành phố"
                                    >
                                        <option value="Hà Nội">Hà Nội</option>
                                        <option value="Hồ Chí Minh">Hồ Chí Minh</option>
                                        <option value="Đà Nẵng">Đà Nẵng</option>
                                        {/* Add more cities as needed */}
                                    </Select>
                                </FormControl>

                                <FormControl id="district" isRequired>
                                    <FormLabel>Quận/Huyện</FormLabel>
                                    <Select
                                        name="district"
                                        value={formData.district}
                                        onChange={handleInputChange}
                                        placeholder="Chọn Quận/Huyện"
                                    >
                                        <option value="Quận Bắc Từ Liêm">Quận Bắc Từ Liêm</option>
                                        <option value="Quận Nam Từ Liêm">Quận Nam Từ Liêm</option>
                                        <option value="Quận Cầu Giấy">Quận Cầu Giấy</option>
                                        {/* Add more districts as needed */}
                                    </Select>
                                </FormControl>
                            </HStack>

                            <FormControl id="ward" isRequired>
                                <FormLabel>Phường/Xã</FormLabel>
                                <Select
                                    name="ward"
                                    value={formData.ward}
                                    onChange={handleInputChange}
                                    placeholder="Chọn Phường/Xã"
                                >
                                    <option value="Phường Minh Khai">Phường Minh Khai</option>
                                    <option value="Phường Tây Tựu">Phường Tây Tựu</option>
                                    <option value="Phường Phú Diễn">Phường Phú Diễn</option>
                                    {/* Add more wards as needed */}
                                </Select>
                            </FormControl>

                            <FormControl id="address" isRequired>
                                <FormLabel>Địa chỉ cụ thể</FormLabel>
                                <Input
                                    name="address"
                                    value={formData.address}
                                    onChange={handleInputChange}
                                    placeholder="Số nhà, tên đường..."
                                />
                            </FormControl>

                            <FormControl id="isDefault">
                                <Checkbox
                                    name="isDefault"
                                    isChecked={formData.isDefault}
                                    onChange={handleInputChange}
                                >
                                    Đặt làm địa chỉ mặc định
                                </Checkbox>
                            </FormControl>
                        </VStack>
                    </ModalBody>
                    <ModalFooter>
                        <Button mr={3} onClick={onAddressModalClose}>
                            Hủy
                        </Button>
                        <Button colorScheme="red" onClick={handleSubmitAddress}>
                            {currentAddress ? 'Cập nhật' : 'Thêm mới'}
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>

            {/* Delete Confirmation Modal */}
            <Modal isOpen={isDeleteModalOpen} onClose={onDeleteModalClose}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Xác nhận xóa địa chỉ</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        Bạn có chắc chắn muốn xóa địa chỉ này không?
                    </ModalBody>
                    <ModalFooter>
                        <Button mr={3} onClick={onDeleteModalClose}>
                            Hủy
                        </Button>
                        <Button
                            colorScheme="red"
                            onClick={handleDeleteAddress}
                            isLoading={isDeleting}
                        >
                            Xóa
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </Box>
    );
};

export default UserAddresses;