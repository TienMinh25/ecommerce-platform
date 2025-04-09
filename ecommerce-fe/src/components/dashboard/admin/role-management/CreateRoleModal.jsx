import React, {useEffect, useState} from 'react';
import {
    Box,
    Button,
    Divider,
    Flex,
    FormControl,
    FormErrorMessage,
    FormLabel,
    Input,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Table,
    Tbody,
    Td,
    Text,
    Textarea,
    Th,
    Thead,
    Tr,
    useColorModeValue,
    useToast,
    VStack,
} from '@chakra-ui/react';
import {FiPlus, FiShield} from 'react-icons/fi';
import roleService from '../../../../services/roleService.js';
import moduleService from "../../../../services/moduleService.js";
import permissionService from "../../../../services/permissionService.js";

const CreateRoleModal = ({ isOpen, onClose, onRoleCreated, modulesList = [], permissionsList = [] }) => {
    // Theme colors
    const borderColor = useColorModeValue('gray.400', 'gray.500');
    const inputBg = useColorModeValue('white', 'gray.900');
    const headerBg = useColorModeValue('blue.50', 'gray.900');
    const textColor = useColorModeValue('gray.900', 'white');
    const labelColor = useColorModeValue('gray.800', 'gray.100');
    const iconColor = useColorModeValue('blue.700', 'blue.200');
    const tableBorderColor = useColorModeValue('gray.100', 'gray.800');
    const hoverBgColor = useColorModeValue('blue.50', 'gray.700');

    // Form state
    const [formData, setFormData] = useState({
        name: '',
        description: '',
    });

    // Modules and permissions state
    const [modules, setModules] = useState([]);
    const [loadingModules, setLoadingModules] = useState(false);

    // Validation and submission state
    const [errors, setErrors] = useState({});
    const [isSubmitting, setIsSubmitting] = useState(false);

    const toast = useToast();

    // Handle input changes
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

    // Form validation
    const validateForm = () => {
        const newErrors = {};
        if (!formData.name.trim()) {
            newErrors.name = 'Role name is required';
        } else if (formData.name.trim().length < 3) {
            newErrors.name = 'Role name must be at least 3 characters';
        }
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    // Handle form submission
    const handleSubmit = async () => {
        if (validateForm()) {
            setIsSubmitting(true);
            try {
                // Step 1: Create the basic role with name and description
                const newRoleData = {
                    name: formData.name.trim(),
                    description: formData.description.trim(),
                };

                // Create the role
                const createdRole = await roleService.createRole(newRoleData);
                const newRoleId = createdRole.data?.id || createdRole.id;

                if (newRoleId) {
                    // Step 2: Identify modules with permissions
                    const modulesWithPermissions = modules
                        .filter(module =>
                            module.read || module.create || module.update ||
                            module.delete || module.approve || module.reject
                        )
                        .map(module => {
                            // Map permissions to their numeric IDs
                            const permissionIds = [];
                            if (module.read) permissionIds.push(1);
                            if (module.create) permissionIds.push(2);
                            if (module.update) permissionIds.push(3);
                            if (module.delete) permissionIds.push(4);
                            if (module.approve) permissionIds.push(5);
                            if (module.reject) permissionIds.push(6);

                            return {
                                module_id: module.id,
                                permissions: permissionIds
                            };
                        });

                    // Step 3: Format the payload according to the API requirements for CREATE
                    const permissionsPayload = {
                        role_name: formData.name.trim(),
                        modules: modulesWithPermissions
                    };

                    // Step 4: Update role permissions if any modules have permissions
                    if (modulesWithPermissions.length > 0) {
                        await roleService.updateRolePermissions(newRoleId, permissionsPayload);
                    }
                }

                toast({
                    title: 'Role created successfully',
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });

                if (onRoleCreated) {
                    onRoleCreated();
                }

                onClose();
            } catch (error) {
                console.error('Error creating role:', error);
                toast({
                    title: 'Failed to create role',
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

    // Reset form when modal closes
    const resetForm = () => {
        setFormData({
            name: '',
            description: '',
        });
        setErrors({});
    };

    // Load modules and permissions when modal opens
    useEffect(() => {
        if (isOpen) {
            // If external modules are provided, use them
            if (modulesList && modulesList.length > 0) {
                // Create the module permission structure with all permissions set to false initially
                const formattedModules = modulesList.map(module => {
                    return {
                        id: module.id,
                        name: module.name,
                        read: false,
                        create: false,
                        update: false,
                        delete: false,
                        approve: false,
                        reject: false
                    };
                });

                setModules(formattedModules);
            } else {
                // Otherwise fetch them
                fetchModulesAndPermissions();
            }
        }
    }, [isOpen, modulesList]);

    // Reset form when modal closes
    useEffect(() => {
        if (!isOpen) {
            resetForm();
        }
    }, [isOpen]);

    // Fetch modules and permissions data
    const fetchModulesAndPermissions = async () => {
        setLoadingModules(true);
        try {
            // Get all modules and permissions in parallel to reduce API calls
            const [modulesResponse, permissionsResponse] = await Promise.all([
                moduleService.getModules({ getAll: true }),
                permissionService.getPermissions({ getAll: true })
            ]);

            const allModules = modulesResponse.data || [];

            // Create the module permission structure with all permissions set to false initially
            const formattedModules = allModules.map(module => {
                return {
                    id: module.id,
                    name: module.name,
                    read: false,
                    create: false,
                    update: false,
                    delete: false,
                    approve: false,
                    reject: false
                };
            });

            setModules(formattedModules);
        } catch (error) {
            console.error('Error loading modules and permissions:', error);
            toast({
                title: 'Error loading permissions',
                description: 'Failed to load module permissions.',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setLoadingModules(false);
        }
    };

    // Handle permission toggle
    const handleTogglePermission = (moduleId, permission) => {
        setModules(modules.map(module =>
            module.id === moduleId
                ? { ...module, [permission]: !module[permission] }
                : module
        ));
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            size="md"
            motionPreset="slideInBottom"
            scrollBehavior="inside"
            isCentered
        >
            <ModalOverlay backdropFilter="blur(3px)" bg="blackAlpha.400" />
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
                        <FiShield size={24} />
                    </Box>
                    <Text fontSize="xl" fontWeight="bold" color={textColor}>Create New Role</Text>
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
                        <FormControl isRequired isInvalid={!!errors.name}>
                            <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>Role Name</FormLabel>
                            <Input
                                name="name"
                                value={formData.name}
                                onChange={handleChange}
                                placeholder="Enter role name"
                                bg={inputBg}
                                color={textColor}
                                borderWidth="1.5px"
                                _hover={{ borderColor: 'blue.400' }}
                                _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                            />
                            {errors.name && <FormErrorMessage fontWeight="medium">{errors.name}</FormErrorMessage>}
                        </FormControl>

                        <FormControl>
                            <FormLabel fontWeight="semibold" fontSize="sm" color={labelColor}>Description (Optional)</FormLabel>
                            <Textarea
                                name="description"
                                value={formData.description}
                                onChange={handleChange}
                                placeholder="Enter role description"
                                bg={inputBg}
                                color={textColor}
                                borderWidth="1.5px"
                                _hover={{ borderColor: 'blue.400' }}
                                _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                minHeight="80px"
                                resize="vertical"
                            />
                        </FormControl>

                        <Divider my={2} />

                        {/* Permissions Section */}
                        <Box>
                            <Text fontSize="md" fontWeight="semibold" color={labelColor} mb={4}>Module Permissions</Text>

                            {loadingModules ? (
                                <Flex justify="center" align="center" py={4}>
                                    <Text color="gray.500">Loading module permissions...</Text>
                                </Flex>
                            ) : (
                                <Box
                                    borderWidth="1px"
                                    borderRadius="md"
                                    borderColor={borderColor}
                                    overflow="auto"
                                    maxH="300px"
                                    mb={2}
                                >
                                    <Table variant="simple" size="sm">
                                        <Thead bg={useColorModeValue('gray.50', 'gray.700')} position="sticky" top={0} zIndex={1}>
                                            <Tr>
                                                <Th fontSize="xs" py={3} width="40%">Module</Th>
                                                <Th textAlign="center" fontSize="xs" py={3}>Read</Th>
                                                <Th textAlign="center" fontSize="xs" py={3}>Create</Th>
                                                <Th textAlign="center" fontSize="xs" py={3}>Update</Th>
                                                <Th textAlign="center" fontSize="xs" py={3}>Delete</Th>
                                                <Th textAlign="center" fontSize="xs" py={3}>Approve</Th>
                                            </Tr>
                                        </Thead>
                                        <Tbody>
                                            {modules.map((module, index) => (
                                                <Tr
                                                    key={`module-${module.id}`}
                                                    _hover={{ bg: hoverBgColor }}
                                                    bg={index % 2 === 0 ? 'transparent' : useColorModeValue('gray.50', 'gray.800')}
                                                >
                                                    <Td py={2}>{renderModuleName(module.name)}</Td>
                                                    <Td textAlign="center" py={2}>
                                                        <PermissionSwitch
                                                            isChecked={module.read}
                                                            onChange={() => handleTogglePermission(module.id, 'read')}
                                                        />
                                                    </Td>
                                                    <Td textAlign="center" py={2}>
                                                        <PermissionSwitch
                                                            isChecked={module.create}
                                                            onChange={() => handleTogglePermission(module.id, 'create')}
                                                        />
                                                    </Td>
                                                    <Td textAlign="center" py={2}>
                                                        <PermissionSwitch
                                                            isChecked={module.update}
                                                            onChange={() => handleTogglePermission(module.id, 'update')}
                                                        />
                                                    </Td>
                                                    <Td textAlign="center" py={2}>
                                                        <PermissionSwitch
                                                            isChecked={module.delete}
                                                            onChange={() => handleTogglePermission(module.id, 'delete')}
                                                        />
                                                    </Td>
                                                    <Td textAlign="center" py={2}>
                                                        <PermissionSwitch
                                                            isChecked={module.approve}
                                                            onChange={() => handleTogglePermission(module.id, 'approve')}
                                                        />
                                                    </Td>
                                                </Tr>
                                            ))}
                                        </Tbody>
                                    </Table>
                                </Box>
                            )}
                            <Text fontSize="xs" color="gray.500" mt={2}>
                                Set permissions for this role. Permissions determine what actions users with this role can perform.
                            </Text>
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
                        Create Role
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default CreateRoleModal;