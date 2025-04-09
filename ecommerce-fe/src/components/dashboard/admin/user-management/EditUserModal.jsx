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
import {FiCalendar, FiCheck, FiMail, FiPhone, FiPlus, FiSave, FiUser} from 'react-icons/fi';
import roleService from '../../../../services/roleService.js';
import userService from '../../../../services/userService.js';

const EditUserModal = ({ isOpen, onClose, user, onUserUpdated }) => {
    // Theme colors
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

    // Form state
    const [formData, setFormData] = useState({
        fullname: '',
        email: '',
        phone: '',
        birthdate: '',
        status: true,
        roles: [],
    });

    // Form validation state
    const [errors, setErrors] = useState({});
    const [isSubmitting, setIsSubmitting] = useState(false);

    // Roles state
    const [availableRoles, setAvailableRoles] = useState([]);
    const [isLoadingRoles, setIsLoadingRoles] = useState(false);
    const [roleError, setRoleError] = useState(null);
    const [showRoleMenu, setShowRoleMenu] = useState(false);
    const roleInputRef = useRef(null);
    const roleMenuRef = useRef(null);

    // Toast for notifications
    const toast = useToast();

    // Load user data when modal opens
    useEffect(() => {
        if (isOpen && user) {
            setFormData({
                fullname: user.fullname || '',
                email: user.email || '',
                phone: user.phone || '',
                birthdate: user.birth_date ? user.birth_date.split('T')[0] : '',
                status: user.status === 'active',
                roles: user.roles.map(role => ({
                    id: role.id,
                    name: role.name,
                })),
            });
            fetchRoles();
        }
    }, [isOpen, user]);

    // Reset form when modal closes
    useEffect(() => {
        if (!isOpen) {
            resetForm();
        }
    }, [isOpen]);

    // Close role menu on outside click
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (
                roleInputRef.current &&
                !roleInputRef.current.contains(event.target) &&
                roleMenuRef.current &&
                !roleMenuRef.current.contains(event.target)
            ) {
                setShowRoleMenu(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, []);

    // Fetch roles from API
    const fetchRoles = async () => {
        setIsLoadingRoles(true);
        setRoleError(null);

        try {
            const roles = await roleService.getRoles();
            if (roles && Array.isArray(roles)) {
                const formattedRoles = roles.map(role => ({
                    id: role.id,
                    name: role.name,
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
            console.error('Error fetching roles:', error);
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

    // Handle status toggle
    const handleStatusToggle = () => {
        setFormData({
            ...formData,
            status: !formData.status,
        });
    };

    // Add role handler
    const handleAddRole = (role) => {
        if (!formData.roles.some(r => r.id === role.id)) {
            setFormData({
                ...formData,
                roles: [...formData.roles, role],
            });
        }
        if (errors.roles) {
            setErrors({
                ...errors,
                roles: null,
            });
        }
        setShowRoleMenu(false);
    };

    // Remove role handler
    const handleRemoveRole = (roleId) => {
        setFormData({
            ...formData,
            roles: formData.roles.filter(role => role.id !== roleId),
        });
    };

    // Form validation (chỉ validate roles)
    const validateForm = () => {
        const newErrors = {};
        if (formData.roles.length === 0) {
            newErrors.roles = 'At least one role is required';
        }
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    // Handle form submission
    const handleSubmit = async () => {
        if (validateForm()) {
            setIsSubmitting(true);
            try {
                const updatedUserData = {
                    status: formData.status ? 'active' : 'inactive',
                    roles: formData.roles.map(role => parseInt(role.id)),
                };

                await userService.updateUser(user.id, updatedUserData);
                toast({
                    title: 'User updated successfully',
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });

                if (onUserUpdated) {
                    onUserUpdated();
                }

                onClose();
            } catch (error) {
                console.error('Error updating user:', error);
                toast({
                    title: 'Failed to update user',
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

    // Reset form to initial state
    const resetForm = () => {
        setFormData({
            fullname: '',
            email: '',
            phone: '',
            birthdate: '',
            status: true,
            roles: [],
        });
        setErrors({});
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
                        <FiUser size={24} />
                    </Box>
                    <Text fontSize="xl" fontWeight="bold" color={textColor}>
                        Edit User
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
                        {/* Account Settings Section */}
                        <Box>
                            <Text fontSize="md" fontWeight="semibold" color={labelColor} mb={4}>
                                Account Settings
                            </Text>

                            {/* Role Selection */}
                            <FormControl isRequired isInvalid={!!errors.roles} mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    User Roles
                                </FormLabel>
                                <Box position="relative">
                                    <Box
                                        onClick={() => setShowRoleMenu(true)}
                                        borderWidth="1.5px"
                                        borderRadius="md"
                                        borderColor={errors.roles ? 'red.500' : borderColor}
                                        p={2}
                                        minH="40px"
                                        bg={inputBg}
                                        cursor="pointer"
                                        ref={roleInputRef}
                                        _hover={{ borderColor: useColorModeValue('blue.400', 'blue.300') }}
                                    >
                                        <Flex flexWrap="wrap" gap={2}>
                                            {formData.roles.map(role => (
                                                <Tag
                                                    key={role.id}
                                                    size="md"
                                                    borderRadius="full"
                                                    variant="solid"
                                                    bg={tagBg}
                                                    color={tagColor}
                                                >
                                                    <TagLabel>{role.name}</TagLabel>
                                                    <TagCloseButton
                                                        onClick={e => {
                                                            e.stopPropagation();
                                                            handleRemoveRole(role.id);
                                                        }}
                                                    />
                                                </Tag>
                                            ))}
                                            {formData.roles.length === 0 && (
                                                <Text color="gray.500" fontSize="sm">
                                                    Select user roles
                                                </Text>
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
                                        >
                                            {isLoadingRoles ? (
                                                <Flex justify="center" align="center" py={3}>
                                                    <Text fontSize="sm" color="gray.500">
                                                        Loading roles...
                                                    </Text>
                                                </Flex>
                                            ) : roleError ? (
                                                <Flex justify="center" align="center" py={3} bg="red.50">
                                                    <Text fontSize="sm" color="red.500">
                                                        {roleError}
                                                    </Text>
                                                </Flex>
                                            ) : availableRoles.length === 0 ? (
                                                <Flex justify="center" align="center" py={3}>
                                                    <Text fontSize="sm" color="gray.500">
                                                        No roles available
                                                    </Text>
                                                </Flex>
                                            ) : (
                                                availableRoles.map(role => (
                                                    <Flex
                                                        key={role.id}
                                                        px={3}
                                                        py={2}
                                                        align="center"
                                                        justify="space-between"
                                                        cursor="pointer"
                                                        _hover={{ bg: menuHoverBg }}
                                                        onClick={() => handleAddRole(role)}
                                                        opacity={
                                                            formData.roles.some(r => r.id === role.id)
                                                                ? 0.5
                                                                : 1
                                                        }
                                                    >
                                                        <Text fontSize="sm" fontWeight="medium" color={textColor}>
                                                            {role.name}
                                                        </Text>
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
                                {errors.roles && (
                                    <FormErrorMessage fontWeight="medium">{errors.roles}</FormErrorMessage>
                                )}
                            </FormControl>

                            {/* Status Toggle */}
                            <FormControl mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    User Status
                                </FormLabel>
                                <Flex
                                    align="center"
                                    justify="space-between"
                                    bg={formData.status ? 'green.50' : 'red.50'}
                                    borderRadius="lg"
                                    p={3}
                                    px={4}
                                    borderWidth="2px"
                                    borderColor={formData.status ? 'green.400' : 'red.400'}
                                    cursor="pointer"
                                    onClick={handleStatusToggle}
                                    _hover={{
                                        borderColor: formData.status ? 'green.500' : 'red.500',
                                        shadow: 'sm',
                                    }}
                                >
                                    <Box>
                                        <Text
                                            fontWeight="bold"
                                            color={formData.status ? 'green.700' : 'red.700'}
                                            mb={0.5}
                                        >
                                            {formData.status ? 'Active' : 'Inactive'}
                                        </Text>
                                        <Text fontSize="xs" color={useColorModeValue('gray.700', 'gray.300')}>
                                            {formData.status
                                                ? 'User will have immediate access to the system'
                                                : 'User account will be suspended'}
                                        </Text>
                                    </Box>
                                    <Switch
                                        isChecked={formData.status}
                                        size="lg"
                                        colorScheme={formData.status ? 'green' : 'red'}
                                        onChange={handleStatusToggle}
                                    />
                                </Flex>
                            </FormControl>
                        </Box>

                        <Divider />

                        {/* User Information Section */}
                        <Box>
                            <Text fontSize="md" fontWeight="semibold" color={labelColor} mb={4}>
                                User Information
                            </Text>

                            {/* User Preview */}
                            <Flex
                                align="center"
                                p={4}
                                bg={useColorModeValue('blue.50', 'blue.900')}
                                borderRadius="lg"
                                mb={4}
                                borderWidth="2px"
                                borderColor={useColorModeValue('blue.200', 'blue.700')}
                            >
                                <Avatar
                                    size="md"
                                    name={formData.fullname || 'User'}
                                    src={
                                        formData.fullname
                                            ? `https://ui-avatars.com/api/?name=${encodeURIComponent(
                                                formData.fullname
                                            )}&background=random&color=fff&size=128`
                                            : '/api/placeholder/100/100'
                                    }
                                    mr={4}
                                    bg={useColorModeValue('blue.500', 'blue.400')}
                                    color="white"
                                />
                                <Box>
                                    <Text fontWeight="bold" color={textColor} fontSize="md">
                                        {formData.fullname || 'User'}
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

                            {/* Full Name Input (Disabled) */}
                            <FormControl mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Full Name
                                </FormLabel>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <Box color={iconColor}>
                                            <FiUser />
                                        </Box>
                                    </InputLeftElement>
                                    <Input
                                        name="fullname"
                                        value={formData.fullname}
                                        isDisabled
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                    />
                                </InputGroup>
                            </FormControl>

                            {/* Email Input (Disabled) */}
                            <FormControl mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Email Address
                                </FormLabel>
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
                                        isDisabled
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                    />
                                </InputGroup>
                            </FormControl>

                            {/* Phone Input (Disabled) */}
                            <FormControl mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Phone Number
                                </FormLabel>
                                <InputGroup>
                                    <InputLeftElement pointerEvents="none">
                                        <Box color={iconColor}>
                                            <FiPhone />
                                        </Box>
                                    </InputLeftElement>
                                    <Input
                                        name="phone"
                                        value={formData.phone}
                                        isDisabled
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                    />
                                </InputGroup>
                            </FormControl>

                            {/* Birth Date Input (Disabled) */}
                            <FormControl mb={4}>
                                <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>
                                    Birth Date
                                </FormLabel>
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
                                        isDisabled
                                        bg={inputBg}
                                        color={textColor}
                                        borderWidth="1.5px"
                                    />
                                </InputGroup>
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
                        Cancel
                    </Button>
                    <Button
                        leftIcon={<FiSave />}  // Icon đã thay đổi
                        colorScheme="blue"
                        onClick={handleSubmit}
                        isLoading={isSubmitting}
                        px={8}
                        shadow="md"
                        bgGradient="linear(to-r, blue.500, blue.600)"
                        _hover={{
                            bgGradient: 'linear(to-r, blue.600, blue.700)',
                            shadow: 'lg',
                            transform: 'translateY(-1px)',
                        }}
                        _active={{
                            bgGradient: 'linear(to-r, blue.700, blue.800)',
                            transform: 'translateY(0)',
                            shadow: 'md',
                        }}
                        fontWeight="bold"
                    >
                        Update User
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default EditUserModal;