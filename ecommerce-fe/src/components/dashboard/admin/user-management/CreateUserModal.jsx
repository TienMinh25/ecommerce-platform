import React, { useEffect, useRef, useState } from 'react';
import {
    Avatar,
    Badge,
    Box,
    Button,
    Divider,
    Flex,
    FormControl,
    FormErrorMessage,
    FormLabel,
    HStack,
    IconButton,
    Input,
    InputGroup,
    InputLeftElement,
    InputRightElement,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Switch,
    Tag,
    TagCloseButton,
    TagLabel,
    Text,
    useColorModeValue,
    useToast,
    VStack,
} from '@chakra-ui/react';
import { FiCalendar, FiCheck, FiEye, FiEyeOff, FiLock, FiMail, FiPhone, FiPlus, FiUser } from 'react-icons/fi';
import roleService from '../../../../services/roleService.js';
import userService from '../../../../services/userService.js';

const CreateUserModal = ({ isOpen, onClose, onUserCreated }) => {
    const borderColor = useColorModeValue('gray.400', 'gray.500');
    const inputBg = useColorModeValue('white', 'gray.900');
    const tagBg = useColorModeValue('blue.100', 'blue.700');
    const tagColor = useColorModeValue('blue.900', 'white');
    const menuBg = useColorModeValue('white', 'gray.900');
    const menuHoverBg = useColorModeValue('blue.50', 'blue.900');
    const headerBg = useColorModeValue('blue.50', 'gray.900');
    const textColor = useColorModeValue('gray.900', 'white');
    const labelColor = useColorModeValue('gray.800', 'gray.100');
    const iconColor = useColorModeValue('blue.700', 'blue.200');

    // Form state - Thêm avatar_url với giá trị mặc định là null
    const [formData, setFormData] = useState({
        fullname: '',
        email: '',
        phone: '',
        birthdate: '',
        password: '',
        confirmPassword: '',
        status: true,
        roles: [],
        avatar_url: null,
    });

    const [errors, setErrors] = useState({});
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [availableRoles, setAvailableRoles] = useState([]);
    const [isLoadingRoles, setIsLoadingRoles] = useState(false);
    const [roleError, setRoleError] = useState(null);
    const [showRoleMenu, setShowRoleMenu] = useState(false);
    const roleInputRef = useRef(null);
    const roleMenuRef = useRef(null);
    const toast = useToast();

    useEffect(() => {
        if (isOpen) {
            fetchRoles();
        }
    }, [isOpen]);

    useEffect(() => {
        if (!isOpen) {
            resetForm();
        }
    }, [isOpen]);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (roleInputRef.current && !roleInputRef.current.contains(event.target) &&
                roleMenuRef.current && !roleMenuRef.current.contains(event.target)) {
                setShowRoleMenu(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, []);

    const fetchRoles = async () => {
        setIsLoadingRoles(true);
        setRoleError(null);
        try {
            const roles = await roleService.getRoles({getAll: true});
            if (roles && Array.isArray(roles)) {
                const formattedRoles = roles.map(role => ({
                    id: role.id,
                    name: role.name
                }));
                setAvailableRoles(formattedRoles);
            } else {
                setAvailableRoles([
                    { id: '1', name: 'admin' },
                    { id: '2', name: 'customer' },
                    { id: '3', name: 'supplier' },
                    { id: '4', name: 'deliverer' }
                ]);
            }
        } catch (error) {
            console.error("Error fetching roles:", error);
            setRoleError(error.message);
            setAvailableRoles([
                { id: '1', name: 'admin' },
                { id: '2', name: 'customer' },
                { id: '3', name: 'supplier' },
                { id: '4', name: 'deliverer' }
            ]);
        } finally {
            setIsLoadingRoles(false);
        }
    };

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

    const handleStatusToggle = () => {
        setFormData({
            ...formData,
            status: !formData.status
        });
    };

    const handleAddRole = (role) => {
        if (!formData.roles.some(r => r.id === role.id)) {
            setFormData({
                ...formData,
                roles: [...formData.roles, role]
            });
        }
        if (errors.roles) {
            setErrors({
                ...errors,
                roles: null
            });
        }
        setShowRoleMenu(false);
    };

    const handleRemoveRole = (roleId) => {
        setFormData({
            ...formData,
            roles: formData.roles.filter(role => role.id !== roleId)
        });
    };

    const validateForm = () => {
        const newErrors = {};
        if (!formData.fullname.trim()) {
            newErrors.fullname = 'Full name is required';
        }
        if (!formData.email.trim()) {
            newErrors.email = 'Email is required';
        } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
            newErrors.email = 'Email is invalid';
        }
        if (formData.phone && !/^\+?[0-9]{10,15}$/.test(formData.phone.replace(/[- ]/g, ''))) {
            newErrors.phone = 'Phone number is invalid';
        }
        if (formData.birthdate) {
            const birthDate = new Date(formData.birthdate);
            const today = new Date();
            if (birthDate > today) {
                newErrors.birthdate = 'Birth date cannot be in the future';
            }
        }
        if (!formData.password) {
            newErrors.password = 'Password is required';
        } else if (formData.password.length < 6) {
            newErrors.password = 'Password must be at least 6 characters';
        }
        if (formData.password !== formData.confirmPassword) {
            newErrors.confirmPassword = 'Passwords do not match';
        }
        if (formData.roles.length === 0) {
            newErrors.roles = 'At least one role is required';
        }
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async () => {
        if (validateForm()) {
            setIsSubmitting(true);
            try {
                // Tạo avatar_url từ ui-avatars.com dựa trên fullname
                const avatarUrl = formData.fullname
                    ? `https://ui-avatars.com/api/?name=${encodeURIComponent(formData.fullname)}&size=128`
                    : 'https://ui-avatars.com/api/?name=New+User&size=128';

                const newUserData = {
                    email: formData.email,
                    password: formData.password,
                    fullname: formData.fullname,
                    phone: formData.phone,
                    birthdate: formData.birthdate,
                    status: formData.status ? 'active' : 'inactive',
                    roles: formData.roles.map(role => parseInt(role.id)),
                    avatar_url: avatarUrl, // Sử dụng URL từ ui-avatars.com
                };

                await userService.createUser(newUserData);
                toast({
                    title: 'User created successfully',
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });

                if (onUserCreated) {
                    onUserCreated();
                }

                onClose();
            } catch (error) {
                console.error('Error creating user:', error);
                toast({
                    title: 'Failed to create user',
                    description: error.response?.data?.error?.message || 'An unexpected error occurred',
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
            fullname: '',
            email: '',
            phone: '',
            birthdate: '',
            password: '',
            confirmPassword: '',
            status: true,
            roles: [],
            avatar_url: null,
        });
        setErrors({});
        setShowPassword(false);
        setShowConfirmPassword(false);
    };

    // Phần return giữ nguyên, không thay đổi giao diện
    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            size="xl"
            motionPreset="slideInBottom"
            scrollBehavior="inside"
        >
            <ModalOverlay
                backdropFilter="blur(3px)"
                bg="blackAlpha.400"
            />
            <ModalContent
                borderRadius="xl"
                shadow="2xl"
                bg={useColorModeValue('white', 'gray.800')}
            >
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
                        <FiUser size={24} />
                    </Box>
                    <Text fontSize="xl" fontWeight="bold" color={textColor}>Tạo người dùng mới</Text>
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
                        <Box>
                            <Text fontSize="md" fontWeight="semibold" color={labelColor} mb={4}>Cài đặt tài khoản</Text>
                            <FormControl isRequired isInvalid={!!errors.roles} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>Vai trò người dùng</FormLabel>
                                <Box position="relative">
                                    <Box
                                        onClick={() => setShowRoleMenu(true)}
                                        borderWidth="1.5px"
                                        borderRadius="md"
                                        borderColor={errors.roles ? "red.500" : borderColor}
                                        p={2}
                                        minH="40px"
                                        bg={inputBg}
                                        cursor="pointer"
                                        position="relative"
                                        ref={roleInputRef}
                                        _hover={{ borderColor: useColorModeValue('blue.400', 'blue.300') }}
                                    >
                                        <Flex flexWrap="wrap" gap={2}>
                                            {formData.roles.map((role) => (
                                                <Tag
                                                    key={role.id}
                                                    size="md"
                                                    borderRadius="full"
                                                    variant="solid"
                                                    bg={tagBg}
                                                    color={tagColor}
                                                >
                                                    <TagLabel>{role.name}</TagLabel>
                                                    <TagCloseButton onClick={(e) => {
                                                        e.stopPropagation();
                                                        handleRemoveRole(role.id);
                                                    }} />
                                                </Tag>
                                            ))}
                                            {formData.roles.length === 0 && (
                                                <Text color="gray.500" fontSize="sm">Chọn vai trò người dùng</Text>
                                            )}
                                        </Flex>
                                    </Box>
                                    {showRoleMenu && (
                                        <Box
                                            ref={roleMenuRef}
                                            position="absolute"
                                            top="calc(100% + 4px)"
                                            left="0"
                                            right="0"
                                            zIndex="dropdown"
                                            borderWidth="1px"
                                            borderRadius="md"
                                            bg={menuBg}
                                            shadow="lg"
                                            maxH="200px"
                                            overflowY="auto"
                                            css={{
                                                '&::-webkit-scrollbar': {
                                                    width: '8px',
                                                },
                                                '&::-webkit-scrollbar-track': {
                                                    background: useColorModeValue('gray.100', 'gray.700'),
                                                    borderRadius: '4px',
                                                },
                                                '&::-webkit-scrollbar-thumb': {
                                                    background: useColorModeValue('blue.400', 'blue.600'),
                                                    borderRadius: '4px',
                                                },
                                                '&::-webkit-scrollbar-thumb:hover': {
                                                    background: useColorModeValue('blue.500', 'blue.500'),
                                                },
                                                scrollBehavior: 'smooth',
                                            }}
                                        >
                                            {isLoadingRoles ? (
                                                <Flex justify="center" align="center" py={3}>
                                                    <Text fontSize="sm" color="gray.500">Loading roles...</Text>
                                                </Flex>
                                            ) : roleError ? (
                                                <Flex justify="center" align="center" py={3} bg="red.50">
                                                    <Text fontSize="sm" color="red.500">{roleError}</Text>
                                                </Flex>
                                            ) : availableRoles.length === 0 ? (
                                                <Flex justify="center" align="center" py={3}>
                                                    <Text fontSize="sm" color="gray.500">No roles available</Text>
                                                </Flex>
                                            ) : (
                                                availableRoles.map((role) => (
                                                    <Flex
                                                        key={role.id}
                                                        px={3}
                                                        py={2}
                                                        align="center"
                                                        justify="space-between"
                                                        cursor="pointer"
                                                        _hover={{ bg: menuHoverBg }}
                                                        onClick={() => handleAddRole(role)}
                                                        opacity={formData.roles.some(r => r.id === role.id) ? 0.5 : 1}
                                                    >
                                                        <Box>
                                                            <Text fontSize="sm" fontWeight="medium" color={textColor}>{role.name}</Text>
                                                            {role.description && (
                                                                <Text fontSize="xs" color="gray.500">{role.description}</Text>
                                                            )}
                                                        </Box>
                                                        {formData.roles.some(r => r.id === role.id) && (
                                                            <Box color="green.500">
                                                                <FiCheck />
                                                            </Box>
                                                        )}
                                                    </Flex>
                                                ))
                                            )}
                                        </Box>
                                    )}
                                </Box>
                                {errors.roles && <FormErrorMessage fontWeight="medium">{errors.roles}</FormErrorMessage>}
                            </FormControl>
                            <FormControl mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>Trạng thái người dùng</FormLabel>
                                <Flex
                                    align="center"
                                    justify="space-between"
                                    bg={formData.status ? "green.50" : "red.50"}
                                    borderRadius="lg"
                                    p={3}
                                    px={4}
                                    borderWidth="2px"
                                    borderColor={formData.status ? "green.400" : "red.400"}
                                    transition="all 0.2s"
                                    cursor="pointer"
                                    onClick={handleStatusToggle}
                                    _hover={{
                                        borderColor: formData.status ? "green.500" : "red.500",
                                        shadow: "sm"
                                    }}
                                >
                                    <Box>
                                        <Text
                                            fontWeight="bold"
                                            color={formData.status ? "green.700" : "red.700"}
                                            mb={0.5}
                                        >
                                            {formData.status ? 'Active' : 'Inactive'}
                                        </Text>
                                        <Text fontSize="xs" color={useColorModeValue('gray.700', 'gray.300')}>
                                            {formData.status
                                                ? 'Người dùng sẽ có quyền truy cập ngay lập tức vào hệ thống'
                                                : 'Tài khoản người dùng sẽ được tạo nhưng bị tạm ngưng'}
                                        </Text>
                                    </Box>
                                    <Switch
                                        isChecked={formData.status}
                                        size="lg"
                                        colorScheme={formData.status ? "green" : "red"}
                                        onChange={handleStatusToggle}
                                    />
                                </Flex>
                            </FormControl>
                        </Box>
                        <Divider />
                        <Box>
                            <Text fontSize="md" fontWeight="semibold" color={labelColor} mb={4}>Thông tin người dùng</Text>
                            <Flex
                                align="center"
                                p={4}
                                bg={useColorModeValue('blue.50', 'blue.900')}
                                borderRadius="lg"
                                mb={4}
                                borderWidth="2px"
                                borderColor={useColorModeValue('blue.200', 'blue.700')}
                                boxShadow="sm"
                            >
                                <Avatar
                                    size="md"
                                    name={formData.fullname || 'New User'}
                                    src={formData.fullname
                                        ? `https://ui-avatars.com/api/?name=${encodeURIComponent(formData.fullname)}&background=random&color=fff&size=128`
                                        : "/api/placeholder/100/100"}
                                    mr={4}
                                    bg={useColorModeValue('blue.500', 'blue.400')}
                                    color="white"
                                />
                                <Box>
                                    <Text fontWeight="bold" color={textColor} fontSize="md">
                                        {formData.fullname || 'New User'}
                                    </Text>
                                    <Text fontSize="sm" color={useColorModeValue('gray.700', 'gray.300')}>
                                        {formData.email || 'email@example.com'}
                                    </Text>
                                    <HStack mt={2} spacing={2}>
                                        <Badge
                                            colorScheme={formData.status ? 'green' : 'red'}
                                            borderRadius="full"
                                            px={2}
                                            py={0.5}
                                            fontSize="xs"
                                            fontWeight="bold"
                                        >
                                            {formData.status ? 'Active' : 'Inactive'}
                                        </Badge>
                                        {formData.roles.map(role => (
                                            <Badge
                                                key={role.id}
                                                colorScheme="blue"
                                                borderRadius="full"
                                                px={2}
                                                fontSize="xs"
                                                fontWeight="bold"
                                            >
                                                {role.name}
                                            </Badge>
                                        ))}
                                    </HStack>
                                </Box>
                            </Flex>
                            <FormControl isRequired isInvalid={!!errors.fullname} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>Họ và tên</FormLabel>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <Box color={iconColor}>
                                            <FiUser />
                                        </Box>
                                    </InputLeftElement>
                                    <Input
                                        name="fullname"
                                        value={formData.fullname}
                                        onChange={handleChange}
                                        placeholder="Nhập họ và tên đầy đủ"
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                        _hover={{ borderColor: 'blue.400' }}
                                        _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    />
                                </InputGroup>
                                {errors.fullname && <FormErrorMessage fontWeight="medium">{errors.fullname}</FormErrorMessage>}
                            </FormControl>
                            <FormControl isRequired isInvalid={!!errors.email} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>Email</FormLabel>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <Box color={iconColor}>
                                            <FiMail />
                                        </Box>
                                    </InputLeftElement>
                                    <Input
                                        name="email"
                                        type="email"
                                        value={formData.email}
                                        onChange={handleChange}
                                        placeholder="Nhập email"
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                        _hover={{ borderColor: 'blue.400' }}
                                        _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    />
                                </InputGroup>
                                {errors.email && <FormErrorMessage fontWeight="medium">{errors.email}</FormErrorMessage>}
                            </FormControl>
                            <FormControl isInvalid={!!errors.phone} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>Số điện thoại (Tuỳ chọn)</FormLabel>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <Box color={iconColor}>
                                            <FiPhone />
                                        </Box>
                                    </InputLeftElement>
                                    <Input
                                        name="phone"
                                        value={formData.phone}
                                        onChange={handleChange}
                                        placeholder="Nhập số điện thoại"
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                        _hover={{ borderColor: 'blue.400' }}
                                        _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    />
                                </InputGroup>
                                {errors.phone && <FormErrorMessage fontWeight="medium">{errors.phone}</FormErrorMessage>}
                            </FormControl>
                            <FormControl isInvalid={!!errors.birthdate} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>Ngày sinh nhật (Tuỳ chọn)</FormLabel>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <Box color={iconColor}>
                                            <FiCalendar />
                                        </Box>
                                    </InputLeftElement>
                                    <Input
                                        name="birthdate"
                                        type="date"
                                        value={formData.birthdate}
                                        onChange={handleChange}
                                        max={new Date().toISOString().split('T')[0]}
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                        _hover={{ borderColor: 'blue.400' }}
                                        _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    />
                                </InputGroup>
                                {errors.birthdate && <FormErrorMessage fontWeight="medium">{errors.birthdate}</FormErrorMessage>}
                            </FormControl>
                        </Box>
                        <Divider />
                        <Box>
                            <Text fontSize="md" fontWeight="semibold" color={labelColor} mb={4}>Mật khẩu</Text>
                            <FormControl isRequired isInvalid={!!errors.password} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>Mật khẩu</FormLabel>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <Box color={iconColor}>
                                            <FiLock />
                                        </Box>
                                    </InputLeftElement>
                                    <Input
                                        name="password"
                                        type={showPassword ? 'text' : 'password'}
                                        value={formData.password}
                                        onChange={handleChange}
                                        placeholder="Nhập mật khẩu"
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                        _hover={{ borderColor: 'blue.400' }}
                                        _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    />
                                    <InputRightElement>
                                        <IconButton
                                            icon={showPassword ? <FiEyeOff /> : <FiEye />}
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => setShowPassword(!showPassword)}
                                            aria-label={showPassword ? "Hide password" : "Show password"}
                                            color={iconColor}
                                            _hover={{ color: 'blue.500', bg: 'transparent' }}
                                        />
                                    </InputRightElement>
                                </InputGroup>
                                {errors.password && <FormErrorMessage fontWeight="medium">{errors.password}</FormErrorMessage>}
                            </FormControl>
                            <FormControl isRequired isInvalid={!!errors.confirmPassword} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>Xác nhận mật khẩu</FormLabel>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <Box color={iconColor}>
                                            <FiLock />
                                        </Box>
                                    </InputLeftElement>
                                    <Input
                                        name="confirmPassword"
                                        type={showConfirmPassword ? 'text' : 'password'}
                                        value={formData.confirmPassword}
                                        onChange={handleChange}
                                        placeholder="Nhập mật khẩu xác nhận"
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                        _hover={{ borderColor: 'blue.400' }}
                                        _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                    />
                                    <InputRightElement>
                                        <IconButton
                                            icon={showConfirmPassword ? <FiEyeOff /> : <FiEye />}
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                            aria-label={showConfirmPassword ? "Hide password" : "Show password"}
                                            color={iconColor}
                                            _hover={{ color: 'blue.500', bg: 'transparent' }}
                                        />
                                    </InputRightElement>
                                </InputGroup>
                                {errors.confirmPassword && <FormErrorMessage fontWeight="medium">{errors.confirmPassword}</FormErrorMessage>}
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
                        leftIcon={<FiPlus />}
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
                        Tạo mới người dùng
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default CreateUserModal;