import React, {useState, useCallback, useLayoutEffect} from 'react';
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
    useColorModeValue,
    Spinner,
    Center
} from '@chakra-ui/react';
import { AddIcon, EditIcon, DeleteIcon, ChevronLeftIcon, ChevronRightIcon, StarIcon } from '@chakra-ui/icons';
import userMeService from "../../services/userMeService.js";
import DeleteConfirmationModal from "../dashboard/admin/DeleteConfirmationComponent.jsx";

const UserAddresses = () => {
    const toast = useToast();
    const [addresses, setAddresses] = useState([]);
    const [currentAddress, setCurrentAddress] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 10,
        totalItems: 0,
        totalPages: 0,
        hasNext: false,
        hasPrevious: false
    });

    // Address types state
    const [addressTypes, setAddressTypes] = useState([]);
    const [isLoadingAddressTypes, setIsLoadingAddressTypes] = useState(false);

    // Province, district, ward states
    const [provinces, setProvinces] = useState([]);
    const [districts, setDistricts] = useState([]);
    const [wards, setWards] = useState([]);
    const [isLoadingProvinces, setIsLoadingProvinces] = useState(false);
    const [isLoadingDistricts, setIsLoadingDistricts] = useState(false);
    const [isLoadingWards, setIsLoadingWards] = useState(false);

    // Modal states
    const {
        isOpen: isAddressModalOpen,
        onOpen: onAddressModalOpen,
        onClose: handleAddressModalClose
    } = useDisclosure();

    // Custom onClose that resets form
    const onAddressModalClose = () => {
        handleAddressModalClose();
        setCurrentAddress(null);
        setFormData({
            recipient_name: '',
            phone: '',
            street: '',
            district: '',
            province: '',
            ward: '',
            country: 'Việt Nam',
            address_type_id: 1,
            is_default: false,
            lattitude: null,
            longtitude: null
        });
        setDistricts([]);
        setWards([]);
    };

    const {
        isOpen: isDeleteModalOpen,
        onOpen: onDeleteModalOpen,
        onClose: onDeleteModalClose
    } = useDisclosure();

    // Form state
    const [formData, setFormData] = useState({
        recipient_name: '',
        phone: '',
        street: '',
        district: '',
        province: '',
        ward: '',
        country: 'Việt Nam',
        address_type_id: 1,
        is_default: false,
        lattitude: null,
        longtitude: null
    });

    // Colors
    const cardBg = useColorModeValue('white', 'gray.800');
    const cardBorderColor = useColorModeValue('gray.200', 'gray.700');
    const defaultBadgeBg = useColorModeValue('red.50', 'red.900');
    const defaultBadgeColor = useColorModeValue('red.600', 'red.200');

    // Fetch addresses - create a memoized version so we can call it after actions
    const fetchAddresses = useCallback(async () => {
        setIsLoading(true);
        try {
            const response = await userMeService.getAddresses({
                page: pagination.page,
                limit: pagination.limit
            });

            // Set addresses directly from response
            setAddresses(response.data || []);

            // Update pagination from response metadata
            if (response && response.metadata && response.metadata.pagination) {
                setPagination(prev => ({
                    ...prev,
                    totalItems: response.metadata.pagination.totalItems,
                    totalPages: response.metadata.pagination.totalPages,
                    hasNext: response.metadata.pagination.hasNext,
                    hasPrevious: response.metadata.pagination.hasPrevious
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

    // Fetch address types
    const fetchAddressTypes = async () => {
        setIsLoadingAddressTypes(true);
        try {
            const response = await userMeService.getAddressTypes({
                page: 1,
                limit: 100
            });
            setAddressTypes(response.data || []);
        } catch (error) {
            console.error('Error fetching address types:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể tải danh sách loại địa chỉ',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoadingAddressTypes(false);
        }
    };

    // Initial load
    useLayoutEffect(() => {
        fetchAddresses();
    }, [fetchAddresses]);

    // Fetch provinces and address types on component mount
    useLayoutEffect(() => {
        fetchProvinces();
        fetchAddressTypes();
    }, []);

    // Fetch districts when province changes
    useLayoutEffect(() => {
        if (formData.province) {
            fetchDistricts(formData.province);
            setFormData(prev => ({ ...prev, district: '', ward: '' })); // Reset district and ward
            setWards([]); // Clear wards when province changes
        } else {
            setDistricts([]);
            setWards([]);
        }
    }, [formData.province]);

    // Fetch wards when district changes
    useLayoutEffect(() => {
        if (formData.province && formData.district) {
            fetchWards(formData.province, formData.district);
            setFormData(prev => ({ ...prev, ward: '' })); // Reset ward
        } else {
            setWards([]);
        }
    }, [formData.district]);

    const fetchProvinces = async () => {
        setIsLoadingProvinces(true);
        try {
            const provinces = await userMeService.getProvinces();
            setProvinces(provinces || []);
        } catch (error) {
            console.error('Error fetching provinces:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể tải danh sách tỉnh/thành phố',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoadingProvinces(false);
        }
    };

    const fetchDistricts = async (provinceId) => {
        if (!provinceId) return;

        setIsLoadingDistricts(true);
        try {
            const districts = await userMeService.getDistricts(provinceId);
            setDistricts(districts || []);
        } catch (error) {
            console.error('Error fetching districts:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể tải danh sách quận/huyện',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoadingDistricts(false);
        }
    };

    const fetchWards = async (provinceId, districtId) => {
        if (!provinceId || !districtId) return;

        setIsLoadingWards(true);
        try {
            const wards = await userMeService.getWards(provinceId, districtId);
            setWards(wards || []);
        } catch (error) {
            console.error('Error fetching wards:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể tải danh sách phường/xã',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoadingWards(false);
        }
    };

    // Handle opening address modal
    const handleAddAddress = () => {
        setCurrentAddress(null);
        setFormData({
            recipient_name: '',
            phone: '',
            street: '',
            district: '',
            province: '',
            ward: '',
            country: 'Việt Nam',
            address_type_id: addressTypes.length > 0 ? addressTypes[0].id : 1,
            is_default: false,
            lattitude: null,
            longtitude: null
        });
        setDistricts([]);
        setWards([]);
        onAddressModalOpen();
    };

    // Handle opening edit modal - optimized to reduce API calls
    const handleEditAddress = useCallback(async (address) => {
        setCurrentAddress(address);

        try {
            // Chỉ đặt dữ liệu cơ bản trước
            setFormData({
                recipient_name: address.recipient_name,
                phone: address.phone,
                street: address.street,
                country: address.country || 'Việt Nam',
                address_type_id: typeof address.address_type_id === 'string'
                    ? parseInt(address.address_type_id, 10)
                    : (address.address_type_id || 1),
                is_default: address.is_default,
                lattitude: address.lattitude || null,
                longtitude: address.longtitude || null,
                province: '',
                district: '',
                ward: ''
            });

            // Mở modal trước để người dùng thấy có phản hồi
            onAddressModalOpen();

            // Sau đó mới fetch dữ liệu

            // Fetch provinces if not already loaded
            let currentProvinces = provinces;
            if (provinces.length === 0) {
                const fetchedProvinces = await userMeService.getProvinces();
                currentProvinces = fetchedProvinces || [];
                setProvinces(currentProvinces);
            }

            // Find province by name
            const provinceObj = currentProvinces.find(p => p.name === address.province);
            if (provinceObj) {
                // Cập nhật province
                setFormData(prev => ({ ...prev, province: provinceObj.id }));

                // Fetch districts
                const fetchedDistricts = await userMeService.getDistricts(provinceObj.id);
                setDistricts(fetchedDistricts || []);

                // Find district
                const districtObj = fetchedDistricts.find(d => d.name === address.district);
                if (districtObj) {
                    setFormData(prev => ({ ...prev, district: districtObj.id }));

                    // Fetch wards
                    const fetchedWards = await userMeService.getWards(provinceObj.id, districtObj.id);
                    setWards(fetchedWards || []);

                    // Find ward
                    const wardObj = fetchedWards.find(w => w.name === address.ward);
                    if (wardObj) {
                        setFormData(prev => ({ ...prev, ward: wardObj.id }));
                    }
                }
            }
        } catch (error) {
            console.error('Error loading address data:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể tải dữ liệu địa chỉ: ' + error.message,
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        }
    }, [onAddressModalOpen, provinces]);

    // Handle opening delete modal
    const handleConfirmDelete = (address) => {
        setCurrentAddress(address);
        onDeleteModalOpen();
    };

    // Handle form input changes
    const handleInputChange = (e) => {
        const { name, value, checked, type } = e.target;

        if (name === "address_type_id") {
            // Chuyển đổi address_type_id thành số nguyên
            setFormData(prev => ({
                ...prev,
                [name]: parseInt(value, 10)
            }));
        } else {
            setFormData(prev => ({
                ...prev,
                [name]: type === 'checkbox' ? checked : value
            }));
        }
    };

    // Handle province change
    const handleProvinceChange = (e) => {
        const provinceId = e.target.value;
        setFormData(prev => ({
            ...prev,
            province: provinceId,
            district: '',
            ward: ''
        }));
    };

    // Handle district change
    const handleDistrictChange = (e) => {
        const districtId = e.target.value;
        setFormData(prev => ({
            ...prev,
            district: districtId,
            ward: ''
        }));
    };

    // Handle ward change
    const handleWardChange = (e) => {
        const wardId = e.target.value;
        setFormData(prev => ({
            ...prev,
            ward: wardId
        }));
    };

    // Validate form data
    const isFormValid = () => {
        return (
            formData.recipient_name &&
            formData.phone &&
            formData.street &&
            formData.district &&
            formData.province &&
            formData.ward &&
            formData.country &&
            formData.address_type_id
        );
    };

    // Handle form submission
    const handleSubmitAddress = async () => {
        setIsSubmitting(true);

        try {
            // Find the full names
            const selectedProvince = provinces.find(p => p.id === formData.province);
            const selectedDistrict = districts.find(d => d.id === formData.district);
            const selectedWard = wards.find(w => w.id === formData.ward);

            if (!selectedProvince || !selectedDistrict || !selectedWard) {
                throw new Error('Vui lòng chọn đầy đủ Tỉnh/Thành phố, Quận/Huyện và Phường/Xã');
            }

            // Đảm bảo address_type_id là số
            const addressTypeId = typeof formData.address_type_id === 'string'
                ? parseInt(formData.address_type_id, 10)
                : formData.address_type_id;

            // Create the data object with names instead of IDs
            const addressData = {
                recipient_name: formData.recipient_name,
                phone: formData.phone,
                street: formData.street,
                district: selectedDistrict.name,
                province: selectedProvince.name,
                ward: selectedWard.name,
                country: formData.country,
                postal_code: formData.postal_code || '',
                address_type_id: addressTypeId,
                is_default: formData.is_default,
                lattitude: formData.lattitude,
                longtitude: formData.longtitude
            };

            if (currentAddress) {
                // Update existing address
                await userMeService.updateAddress(currentAddress.id, addressData);
                toast({
                    title: 'Địa chỉ đã được cập nhật',
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });
            } else {
                // Add new address
                await userMeService.addAddress(addressData);
                toast({
                    title: 'Địa chỉ mới đã được thêm',
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });
            }

            // Refresh addresses - stay on current page
            await fetchAddresses();
            onAddressModalClose();
        } catch (error) {
            console.error('Error saving address:', error);
            toast({
                title: 'Lỗi',
                description: error.response?.data?.error?.message || error.message || 'Không thể lưu địa chỉ',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle address deletion
    const handleDeleteAddress = async () => {
        setIsDeleting(true);
        try {
            await userMeService.deleteAddress(currentAddress.id);

            toast({
                title: 'Đã xóa địa chỉ',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });

            // Refresh addresses - stay on current page
            await fetchAddresses();
            onDeleteModalClose();
        } catch (error) {
            console.error('Error deleting address:', error);
            toast({
                title: 'Lỗi',
                description: error.response?.data?.error?.message || 'Không thể xóa địa chỉ',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsDeleting(false);
        }
    };

    // Handle setting address as default
    const handleSetDefault = async (addressId) => {
        try {
            await userMeService.setDefaultAddress(addressId);

            toast({
                title: 'Đã đặt địa chỉ mặc định',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });

            // Refresh addresses - stay on current page
            await fetchAddresses();
        } catch (error) {
            console.error('Error setting default address:', error);
            toast({
                title: 'Lỗi',
                description: error.response?.data?.error?.message || 'Không thể đặt địa chỉ mặc định',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        }
    };

    // Handle pagination
    const handlePageChange = (newPage) => {
        setPagination(prev => ({
            ...prev,
            page: newPage
        }));
    };

    return (
        <Box position="relative" pb={16} minH="500px">
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

            {/* Loading State */}
            {isLoading ? (
                <Center py={10}>
                    <Spinner size="xl" color="red.500" />
                </Center>
            ) : (
                <>
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
                                    borderColor={address.is_default ? "red.300" : cardBorderColor}
                                    bg={cardBg}
                                    position="relative"
                                >
                                    <HStack spacing={4} mb={3}>
                                        <Text fontWeight="bold">{address.recipient_name}</Text>
                                        <Text color="gray.600">{address.phone}</Text>
                                        {address.is_default && (
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

                                    <Text fontSize="sm" color="gray.700" mb={1}>
                                        {address.street}
                                    </Text>
                                    <Text fontSize="sm" color="gray.700" mb={3}>{address.ward}, {address.district}, {address.province}, {address.country}</Text>

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
                                            {!address.is_default && (
                                                <Button
                                                    size="sm"
                                                    variant="outline"
                                                    leftIcon={<StarIcon />}
                                                    colorScheme="yellow"
                                                    onClick={() => handleSetDefault(address.id)}
                                                >
                                                    Đặt làm mặc định
                                                </Button>
                                            )}
                                        </HStack>
                                    </Flex>
                                </Box>
                            ))
                        )}
                    </VStack>
                </>
            )}

            {/* Enhanced Pagination Component */}
            {!isLoading && addresses.length > 0 && (
                <Flex
                    justify="center"
                    bottom={0}
                    left={0}
                    right={0}
                    marginTop={6}
                    py={4}
                    borderTopWidth="1px"
                    borderColor={cardBorderColor}
                    bg={cardBg}
                    boxShadow="0 -2px 10px rgba(0,0,0,0.05)"
                >
                    <HStack spacing={1}>
                        <Button
                            size="md"
                            variant="outline"
                            colorScheme="red"
                            onClick={() => handlePageChange(1)}
                            isDisabled={!pagination.hasPrevious}
                            borderRadius="md"
                            display={pagination.page > 2 ? "flex" : "none"}
                        >
                            <ChevronLeftIcon boxSize={5} />
                            <ChevronLeftIcon boxSize={5} marginLeft="-1.5" />
                        </Button>

                        <Button
                            size="md"
                            variant="outline"
                            colorScheme="red"
                            onClick={() => handlePageChange(pagination.page - 1)}
                            isDisabled={!pagination.hasPrevious}
                            borderRadius="md"
                        >
                            <ChevronLeftIcon boxSize={5} />
                        </Button>

                        {/* Page Numbers */}
                        {[...Array(pagination.totalPages)].map((_, idx) => {
                            const pageNum = idx + 1;
                            // Only show current page, one before, one after, first and last page
                            if (
                                pageNum === 1 ||
                                pageNum === pagination.totalPages ||
                                (pageNum >= pagination.page - 1 && pageNum <= pagination.page + 1)
                            ) {
                                return (
                                    <Button
                                        key={pageNum}
                                        size="md"
                                        variant={pageNum === pagination.page ? "solid" : "outline"}
                                        colorScheme="red"
                                        onClick={() => handlePageChange(pageNum)}
                                        borderRadius="md"
                                    >
                                        {pageNum}
                                    </Button>
                                );
                            }

                            // Show ellipsis for skipped pages
                            if (
                                (pageNum === 2 && pagination.page > 3) ||
                                (pageNum === pagination.totalPages - 1 && pagination.page < pagination.totalPages - 2)
                            ) {
                                return <Text key={pageNum} mx={2}>...</Text>;
                            }

                            return null;
                        })}

                        <Button
                            size="md"
                            variant="outline"
                            colorScheme="red"
                            onClick={() => handlePageChange(pagination.page + 1)}
                            isDisabled={!pagination.hasNext}
                            borderRadius="md"
                        >
                            <ChevronRightIcon boxSize={5} />
                        </Button>

                        <Button
                            size="md"
                            variant="outline"
                            colorScheme="red"
                            onClick={() => handlePageChange(pagination.totalPages)}
                            isDisabled={!pagination.hasNext || pagination.page === pagination.totalPages}
                            borderRadius="md"
                            display={pagination.page < pagination.totalPages - 1 ? "flex" : "none"}
                        >
                            <ChevronRightIcon boxSize={5} />
                            <ChevronRightIcon boxSize={5} marginLeft="-1.5" />
                        </Button>
                    </HStack>
                </Flex>
            )}

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
                            <FormControl id="recipient_name" isRequired>
                                <FormLabel>Họ và tên</FormLabel>
                                <Input
                                    name="recipient_name"
                                    value={formData.recipient_name}
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

                            <FormControl id="country" isRequired>
                                <FormLabel>Quốc gia</FormLabel>
                                <Select
                                    name="country"
                                    value={formData.country}
                                    onChange={handleInputChange}
                                >
                                    <option value="Việt Nam">Việt Nam</option>
                                </Select>
                            </FormControl>

                            <FormControl id="province" isRequired>
                                <FormLabel>Tỉnh/Thành phố</FormLabel>
                                <Select
                                    name="province"
                                    value={formData.province}
                                    onChange={handleProvinceChange}
                                    placeholder="Chọn Tỉnh/Thành phố"
                                    isDisabled={isLoadingProvinces}
                                >
                                    {provinces.map(province => (
                                        <option key={province.id} value={province.id}>
                                            {province.name}
                                        </option>
                                    ))}
                                </Select>
                                {isLoadingProvinces && <Spinner size="sm" ml={2} />}
                            </FormControl>

                            <FormControl id="district" isRequired>
                                <FormLabel>Quận/Huyện</FormLabel>
                                <Select
                                    name="district"
                                    value={formData.district}
                                    onChange={handleDistrictChange}
                                    placeholder="Chọn Quận/Huyện"
                                    isDisabled={isLoadingDistricts || !formData.province}
                                >
                                    {districts.map(district => (
                                        <option key={district.id} value={district.id}>
                                            {district.name}
                                        </option>
                                    ))}
                                </Select>
                                {isLoadingDistricts && <Spinner size="sm" ml={2} />}
                            </FormControl>

                            <FormControl id="ward" isRequired>
                                <FormLabel>Phường/Xã</FormLabel>
                                <Select
                                    name="ward"
                                    value={formData.ward}
                                    onChange={handleWardChange}
                                    placeholder="Chọn Phường/Xã"
                                    isDisabled={isLoadingWards || !formData.district}
                                >
                                    {wards.map(ward => (
                                        <option key={ward.id} value={ward.id}>
                                            {ward.name}
                                        </option>
                                    ))}
                                </Select>
                                {isLoadingWards && <Spinner size="sm" ml={2} />}
                            </FormControl>

                            <FormControl id="street" isRequired>
                                <FormLabel>Địa chỉ cụ thể</FormLabel>
                                <Input
                                    name="street"
                                    value={formData.street}
                                    onChange={handleInputChange}
                                    placeholder="Số nhà, tên đường..."
                                />
                            </FormControl>

                            <FormControl id="address_type_id" isRequired>
                                <FormLabel>Loại địa chỉ</FormLabel>
                                <Select
                                    name="address_type_id"
                                    value={formData.address_type_id}
                                    onChange={handleInputChange}
                                    placeholder="Chọn loại địa chỉ"
                                    isDisabled={isLoadingAddressTypes}
                                >
                                    {addressTypes.map(type => (
                                        <option key={type.id} value={type.id}>
                                            {type.address_type}
                                        </option>
                                    ))}
                                </Select>
                                {isLoadingAddressTypes && <Spinner size="sm" ml={2} />}
                            </FormControl>

                            <FormControl id="is_default">
                                <Checkbox
                                    name="is_default"
                                    isChecked={formData.is_default}
                                    onChange={handleInputChange}
                                >
                                    Đặt làm địa chỉ mặc định
                                </Checkbox>
                            </FormControl>
                        </VStack>
                    </ModalBody>
                    <ModalFooter>
                        <Button mr={3} onClick={onAddressModalClose} isDisabled={isSubmitting}>
                            Hủy
                        </Button>
                        <Button
                            colorScheme="red"
                            onClick={handleSubmitAddress}
                            isLoading={isSubmitting}
                            isDisabled={!isFormValid()}
                        >
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
                        <Button mr={3} onClick={onDeleteModalClose} isDisabled={isDeleting}>
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