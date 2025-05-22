import React, { useEffect, useState } from 'react';
import {
    Box,
    Button,
    Divider,
    Flex,
    FormControl,
    FormErrorMessage,
    FormLabel,
    Input,
    InputGroup,
    InputLeftElement,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    NumberInput,
    NumberInputField,
    NumberInputStepper,
    NumberIncrementStepper,
    NumberDecrementStepper,
    Select,
    Switch,
    Text,
    Textarea,
    useColorModeValue,
    useToast,
    VStack,
    Badge,
} from '@chakra-ui/react';
import { FiPercent, FiDollarSign, FiTag, FiSave, FiCalendar } from 'react-icons/fi';
import couponService from "../../../../services/couponService.js";

const EditCouponModal = ({ isOpen, onClose, coupon, onCouponUpdated }) => {
    const borderColor = useColorModeValue('gray.400', 'gray.500');
    const inputBg = useColorModeValue('white', 'gray.900');
    const headerBg = useColorModeValue('blue.50', 'gray.900');
    const textColor = useColorModeValue('gray.900', 'white');
    const labelColor = useColorModeValue('gray.800', 'gray.100');
    const iconColor = useColorModeValue('blue.700', 'blue.200');

    // Form state
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        discount_type: 'percentage',
        discount_value: 0,
        maximum_discount_amount: 0,
        minimum_order_amount: 0,
        start_date: '',
        end_date: '',
        usage_limit: 1,
        is_active: true,
    });

    const [errors, setErrors] = useState({});
    const [isSubmitting, setIsSubmitting] = useState(false);
    const toast = useToast();

    // Load coupon data when modal opens
    useEffect(() => {
        if (isOpen && coupon) {
            setFormData({
                name: coupon.name || '',
                description: coupon.description || '',
                discount_type: coupon.discount_type || 'percentage',
                discount_value: coupon.discount_value || 0,
                maximum_discount_amount: coupon.maximum_discount_amount || 0,
                minimum_order_amount: coupon.minimum_order_amount || 0,
                start_date: coupon.start_date ? new Date(coupon.start_date).toISOString().slice(0, 16) : '',
                end_date: coupon.end_date ? new Date(coupon.end_date).toISOString().slice(0, 16) : '',
                usage_limit: coupon.usage_limit || 1,
                is_active: coupon.is_active !== undefined ? coupon.is_active : true,
            });
        }
    }, [isOpen, coupon]);

    useEffect(() => {
        if (!isOpen) {
            resetForm();
        }
    }, [isOpen]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: value
        });
        if (errors[name]) {
            setErrors({
                ...errors,
                [name]: null
            });
        }
    };

    const handleNumberChange = (name, value) => {
        setFormData({
            ...formData,
            [name]: value
        });
        if (errors[name]) {
            setErrors({
                ...errors,
                [name]: null
            });
        }
    };

    const handleSwitchChange = (checked) => {
        setFormData({
            ...formData,
            is_active: checked
        });
    };

    const validateForm = () => {
        const newErrors = {};

        if (!formData.name.trim()) {
            newErrors.name = 'Tên khuyến mãi là bắt buộc';
        }

        if (!formData.discount_value || formData.discount_value <= 0) {
            newErrors.discount_value = 'Giá trị giảm giá phải lớn hơn 0';
        }

        if (formData.discount_type === 'percentage' && formData.discount_value > 100) {
            newErrors.discount_value = 'Phần trăm giảm giá không được vượt quá 100%';
        }

        if (!formData.maximum_discount_amount || formData.maximum_discount_amount <= 0) {
            newErrors.maximum_discount_amount = 'Số tiền giảm tối đa phải lớn hơn 0';
        }

        if (formData.minimum_order_amount < 0) {
            newErrors.minimum_order_amount = 'Số tiền đơn hàng tối thiểu không được âm';
        }

        if (!formData.start_date) {
            newErrors.start_date = 'Ngày bắt đầu là bắt buộc';
        }

        if (!formData.end_date) {
            newErrors.end_date = 'Ngày kết thúc là bắt buộc';
        }

        if (formData.start_date && formData.end_date) {
            const startDate = new Date(formData.start_date);
            const endDate = new Date(formData.end_date);
            if (endDate <= startDate) {
                newErrors.end_date = 'Ngày kết thúc phải sau ngày bắt đầu';
            }
        }

        if (!formData.usage_limit || formData.usage_limit <= 0) {
            newErrors.usage_limit = 'Số lần sử dụng phải lớn hơn 0';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async () => {
        if (validateForm()) {
            setIsSubmitting(true);
            try {
                // Convert dates to ISO format
                const startDate = new Date(formData.start_date);
                const endDate = new Date(formData.end_date);
                endDate.setHours(23, 59, 59, 999); // Set to end of day

                const updatedCouponData = {
                    name: formData.name,
                    description: formData.description,
                    discount_type: formData.discount_type,
                    discount_value: parseFloat(formData.discount_value),
                    maximum_discount_amount: parseFloat(formData.maximum_discount_amount),
                    minimum_order_amount: parseFloat(formData.minimum_order_amount),
                    start_date: startDate.toISOString(),
                    end_date: endDate.toISOString(),
                    usage_limit: parseInt(formData.usage_limit),
                    is_active: formData.is_active,
                };

                await couponService.updateCoupon(coupon.id, updatedCouponData);
                toast({
                    title: 'Cập nhật mã khuyến mãi thành công',
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });

                if (onCouponUpdated) {
                    onCouponUpdated();
                }

                onClose();
            } catch (error) {
                console.error('Error updating coupon:', error);
                toast({
                    title: 'Cập nhật mã khuyến mãi thất bại',
                    description: error.response?.data?.error?.message || 'Đã xảy ra lỗi không mong muốn',
                    status: 'error',
                    duration: 5000,
                    isClosable: true,
                });
            } finally {
                setIsSubmitting(false);
            }
        }
    };

    const resetForm = () => {
        setFormData({
            name: '',
            description: '',
            discount_type: 'percentage',
            discount_value: 0,
            maximum_discount_amount: 0,
            minimum_order_amount: 0,
            start_date: '',
            end_date: '',
            usage_limit: 1,
            is_active: true,
        });
        setErrors({});
    };

    const formatCurrency = (value) => {
        return new Intl.NumberFormat('vi-VN', {
            style: 'currency',
            currency: 'VND'
        }).format(value);
    };

    const formatDate = (dateString) => {
        if (!dateString) return 'N/A';
        try {
            const date = new Date(dateString);
            return date.toLocaleString('vi-VN');
        } catch (e) {
            return dateString;
        }
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            size="xl"
            motionPreset="slideInBottom"
            scrollBehavior="inside"
        >
            <ModalOverlay backdropFilter="blur(3px)" bg="blackAlpha.400" />
            <ModalContent borderRadius="xl" shadow="2xl" bg={useColorModeValue('white', 'gray.800')}>
                <ModalHeader
                    py={6}
                    borderBottom="1px solid"
                    borderColor={borderColor}
                    bg={headerBg}
                    borderTopRadius="xl"
                    display="flex"
                    alignItems="center"
                >
                    <Box color={iconColor} mr={3}>
                        <FiTag size={24} />
                    </Box>
                    <Text fontSize="xl" fontWeight="bold" color={textColor}>
                        Chỉnh sửa mã khuyến mãi
                    </Text>
                </ModalHeader>
                <ModalCloseButton
                    size="lg"
                    top={3}
                    right={3}
                    borderRadius="full"
                    p={2}
                    m={2}
                    _hover={{ bg: useColorModeValue('gray.200', 'gray.700') }}
                />

                <ModalBody py={6}>
                    <VStack spacing={6} align="stretch">
                        {/* Coupon Preview */}
                        {coupon && (
                            <Box>
                                <Text fontSize="md" fontWeight="semibold" color={labelColor} mb={4}>
                                    Thông tin hiện tại
                                </Text>
                                <Flex
                                    align="center"
                                    p={4}
                                    bg={useColorModeValue('blue.50', 'blue.900')}
                                    borderRadius="lg"
                                    mb={4}
                                    borderWidth="2px"
                                    borderColor={useColorModeValue('blue.200', 'blue.700')}
                                >
                                    <Box flex="1">
                                        <Flex align="center" mb={2}>
                                            <Text fontWeight="bold" color={textColor} fontSize="lg" mr={2}>
                                                {coupon.code}
                                            </Text>
                                            <Badge
                                                colorScheme={coupon.is_active ? 'green' : 'red'}
                                                borderRadius="full"
                                                px={2}
                                                py={0.5}
                                                fontSize="xs"
                                                fontWeight="bold"
                                            >
                                                {coupon.is_active ? 'Hoạt động' : 'Không hoạt động'}
                                            </Badge>
                                        </Flex>
                                        <Text fontSize="sm" color={useColorModeValue('gray.700', 'gray.300')} mb={1}>
                                            Đã sử dụng: {coupon.usage_count}/{coupon.usage_limit}
                                        </Text>
                                        <Text fontSize="sm" color={useColorModeValue('gray.700', 'gray.300')}>
                                            Thời gian: {formatDate(coupon.start_date)} - {formatDate(coupon.end_date)}
                                        </Text>
                                    </Box>
                                </Flex>
                            </Box>
                        )}

                        {/* Basic Information */}
                        <Box>
                            <Text fontSize="md" fontWeight="semibold" color={labelColor} mb={4}>
                                Thông tin cơ bản
                            </Text>

                            <FormControl isRequired isInvalid={!!errors.name} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Tên khuyến mãi
                                </FormLabel>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <Box color={iconColor}>
                                            <FiTag />
                                        </Box>
                                    </InputLeftElement>
                                    <Input
                                        name="name"
                                        value={formData.name}
                                        onChange={handleChange}
                                        placeholder="Nhập tên khuyến mãi"
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                        _hover={{ borderColor: 'blue.400' }}
                                        _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    />
                                </InputGroup>
                                {errors.name && <FormErrorMessage fontWeight="medium">{errors.name}</FormErrorMessage>}
                            </FormControl>

                            <FormControl mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Mô tả
                                </FormLabel>
                                <Textarea
                                    name="description"
                                    value={formData.description}
                                    onChange={handleChange}
                                    placeholder="Nhập mô tả cho khuyến mãi"
                                    bg={inputBg}
                                    color={textColor}
                                    borderWidth="1.5px"
                                    _hover={{ borderColor: 'blue.400' }}
                                    _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    rows={3}
                                />
                            </FormControl>

                            {/* Active Status */}
                            <FormControl mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Trạng thái hoạt động
                                </FormLabel>
                                <Flex
                                    align="center"
                                    justify="space-between"
                                    bg={formData.is_active ? 'green.50' : 'red.50'}
                                    borderRadius="lg"
                                    p={3}
                                    px={4}
                                    borderWidth="2px"
                                    borderColor={formData.is_active ? 'green.400' : 'red.400'}
                                    cursor="pointer"
                                    onClick={() => handleSwitchChange(!formData.is_active)}
                                    _hover={{
                                        borderColor: formData.is_active ? 'green.500' : 'red.500',
                                        shadow: 'sm',
                                    }}
                                >
                                    <Box>
                                        <Text
                                            fontWeight="bold"
                                            color={formData.is_active ? 'green.700' : 'red.700'}
                                            mb={0.5}
                                        >
                                            {formData.is_active ? 'Đang hoạt động' : 'Không hoạt động'}
                                        </Text>
                                        <Text fontSize="xs" color={useColorModeValue('gray.700', 'gray.300')}>
                                            {formData.is_active
                                                ? 'Mã khuyến mãi có thể được sử dụng'
                                                : 'Mã khuyến mãi sẽ bị vô hiệu hóa'}
                                        </Text>
                                    </Box>
                                    <Switch
                                        isChecked={formData.is_active}
                                        size="lg"
                                        colorScheme={formData.is_active ? 'green' : 'red'}
                                        onChange={(e) => handleSwitchChange(e.target.checked)}
                                    />
                                </Flex>
                            </FormControl>
                        </Box>

                        <Divider />

                        {/* Discount Settings */}
                        <Box>
                            <Text fontSize="md" fontWeight="semibold" color={labelColor} mb={4}>
                                Cài đặt giảm giá
                            </Text>

                            <FormControl isRequired isInvalid={!!errors.discount_type} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Loại giảm giá
                                </FormLabel>
                                <Select
                                    name="discount_type"
                                    value={formData.discount_type}
                                    onChange={handleChange}
                                    bg={inputBg}
                                    color={textColor}
                                    borderWidth="1.5px"
                                    _hover={{ borderColor: 'blue.400' }}
                                    _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                >
                                    <option value="percentage">Phần trăm (%)</option>
                                    <option value="fixed_amount">Số tiền cố định</option>
                                </Select>
                            </FormControl>

                            <FormControl isRequired isInvalid={!!errors.discount_value} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Giá trị giảm giá
                                </FormLabel>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <Box color={iconColor}>
                                            {formData.discount_type === 'percentage' ? <FiPercent /> : <FiDollarSign />}
                                        </Box>
                                    </InputLeftElement>
                                    <NumberInput
                                        value={formData.discount_value}
                                        onChange={(value) => handleNumberChange('discount_value', value)}
                                        min={0}
                                        max={formData.discount_type === 'percentage' ? 100 : undefined}
                                        precision={formData.discount_type === 'percentage' ? 2 : 0}
                                        width="100%"
                                    >
                                        <NumberInputField
                                            bg={inputBg}
                                            color={textColor}
                                            borderWidth="1.5px"
                                            _hover={{ borderColor: 'blue.400' }}
                                            _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                            pl={10}
                                        />
                                        <NumberInputStepper>
                                            <NumberIncrementStepper />
                                            <NumberDecrementStepper />
                                        </NumberInputStepper>
                                    </NumberInput>
                                </InputGroup>
                                {errors.discount_value && <FormErrorMessage fontWeight="medium">{errors.discount_value}</FormErrorMessage>}
                                {formData.discount_type === 'percentage' && (
                                    <Text fontSize="xs" color="gray.500" mt={1}>
                                        Phần trăm giảm giá (0-100%)
                                    </Text>
                                )}
                            </FormControl>

                            <FormControl isRequired isInvalid={!!errors.maximum_discount_amount} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Số tiền giảm tối đa (VND)
                                </FormLabel>
                                <NumberInput
                                    value={formData.maximum_discount_amount}
                                    onChange={(value) => handleNumberChange('maximum_discount_amount', value)}
                                    min={0}
                                    precision={0}
                                    width="100%"
                                >
                                    <NumberInputField
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                        _hover={{ borderColor: 'blue.400' }}
                                        _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    />
                                    <NumberInputStepper>
                                        <NumberIncrementStepper />
                                        <NumberDecrementStepper />
                                    </NumberInputStepper>
                                </NumberInput>
                                {errors.maximum_discount_amount && <FormErrorMessage fontWeight="medium">{errors.maximum_discount_amount}</FormErrorMessage>}
                                {formData.maximum_discount_amount > 0 && (
                                    <Text fontSize="xs" color="gray.500" mt={1}>
                                        {formatCurrency(formData.maximum_discount_amount)}
                                    </Text>
                                )}
                            </FormControl>

                            <FormControl isInvalid={!!errors.minimum_order_amount} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Số tiền đơn hàng tối thiểu (VND)
                                </FormLabel>
                                <NumberInput
                                    value={formData.minimum_order_amount}
                                    onChange={(value) => handleNumberChange('minimum_order_amount', value)}
                                    min={0}
                                    precision={0}
                                    width="100%"
                                >
                                    <NumberInputField
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                        _hover={{ borderColor: 'blue.400' }}
                                        _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    />
                                    <NumberInputStepper>
                                        <NumberIncrementStepper />
                                        <NumberDecrementStepper />
                                    </NumberInputStepper>
                                </NumberInput>
                                {errors.minimum_order_amount && <FormErrorMessage fontWeight="medium">{errors.minimum_order_amount}</FormErrorMessage>}
                                {formData.minimum_order_amount > 0 && (
                                    <Text fontSize="xs" color="gray.500" mt={1}>
                                        {formatCurrency(formData.minimum_order_amount)}
                                    </Text>
                                )}
                            </FormControl>
                        </Box>

                        <Divider />

                        {/* Time and Usage Settings */}
                        <Box>
                            <Text fontSize="md" fontWeight="semibold" color={labelColor} mb={4}>
                                Thời gian và sử dụng
                            </Text>

                            <Flex gap={4} mb={4}>
                                <FormControl isRequired isInvalid={!!errors.start_date} flex="1">
                                    <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                        Ngày bắt đầu
                                    </FormLabel>
                                    <InputGroup>
                                        <InputLeftElement pointerEvents="none">
                                            <Box color={iconColor}>
                                                <FiCalendar />
                                            </Box>
                                        </InputLeftElement>
                                        <Input
                                            name="start_date"
                                            type="datetime-local"
                                            value={formData.start_date}
                                            onChange={handleChange}
                                            bg={inputBg}
                                            color={textColor}
                                            borderWidth="1.5px"
                                            _hover={{ borderColor: 'blue.400' }}
                                            _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                        />
                                    </InputGroup>
                                    {errors.start_date && <FormErrorMessage fontWeight="medium">{errors.start_date}</FormErrorMessage>}
                                </FormControl>

                                <FormControl isRequired isInvalid={!!errors.end_date} flex="1">
                                    <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                        Ngày kết thúc
                                    </FormLabel>
                                    <InputGroup>
                                        <InputLeftElement pointerEvents="none">
                                            <Box color={iconColor}>
                                                <FiCalendar />
                                            </Box>
                                        </InputLeftElement>
                                        <Input
                                            name="end_date"
                                            type="datetime-local"
                                            value={formData.end_date}
                                            onChange={handleChange}
                                            min={formData.start_date}
                                            bg={inputBg}
                                            color={textColor}
                                            borderWidth="1.5px"
                                            _hover={{ borderColor: 'blue.400' }}
                                            _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                        />
                                    </InputGroup>
                                    {errors.end_date && <FormErrorMessage fontWeight="medium">{errors.end_date}</FormErrorMessage>}
                                </FormControl>
                            </Flex>

                            <FormControl isRequired isInvalid={!!errors.usage_limit} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Số lần sử dụng tối đa
                                </FormLabel>
                                <NumberInput
                                    value={formData.usage_limit}
                                    onChange={(value) => handleNumberChange('usage_limit', value)}
                                    min={1}
                                    precision={0}
                                    width="100%"
                                >
                                    <NumberInputField
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                        _hover={{ borderColor: 'blue.400' }}
                                        _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    />
                                    <NumberInputStepper>
                                        <NumberIncrementStepper />
                                        <NumberDecrementStepper />
                                    </NumberInputStepper>
                                </NumberInput>
                                {errors.usage_limit && <FormErrorMessage fontWeight="medium">{errors.usage_limit}</FormErrorMessage>}
                            </FormControl>
                        </Box>
                    </VStack>
                </ModalBody>

                <ModalFooter
                    borderTop="1px solid"
                    borderColor={borderColor}
                    bg={headerBg}
                    borderBottomRadius="xl"
                    justifyContent="space-between"
                    py={4}
                >
                    <Button
                        onClick={onClose}
                        variant="outline"
                        colorScheme="gray"
                        px={6}
                        borderColor={borderColor}
                        _hover={{ bg: useColorModeValue('gray.200', 'gray.700') }}
                    >
                        Huỷ
                    </Button>
                    <Button
                        leftIcon={<FiSave />}
                        colorScheme="blue"
                        onClick={handleSubmit}
                        isLoading={isSubmitting}
                        px={8}
                        shadow="md"
                        bgGradient="linear(to-r, blue.500, blue.600)"
                        _hover={{
                            bgGradient: "linear(to-r, blue.600, blue.700)",
                            shadow: 'lg',
                            transform: 'translateY(-1px)'
                        }}
                        _active={{
                            bgGradient: "linear(to-r, blue.700, blue.800)",
                            transform: 'translateY(0)',
                            shadow: 'md'
                        }}
                        fontWeight="bold"
                    >
                        Cập nhật
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default EditCouponModal;