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
    Tooltip,
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
                // Identify modules with permissions
                const modules_permissions = modules
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

                // Create the role data according to the API schema from Swagger
                const roleData = {
                    role_name: formData.name.trim(),
                    description: formData.description.trim(),
                    modules_permissions: modules_permissions
                };

                // Create the role with all data in one request
                const createdRole = await roleService.createRole(roleData);

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

    // Permission Switch component with tooltip
    const PermissionSwitch = ({ isChecked, onChange, permission }) => {
        return (
            <Tooltip
                label={isChecked ? "Enabled" : "Disabled"}
                hasArrow
                placement="top"
                openDelay={500}
            >
                <Box
                    position="relative"
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    cursor="pointer"
                    onClick={onChange}
                >
                    <Box
                        w="36px"
                        h="20px"
                        bg={isChecked ? "blue.500" : "gray.300"}
                        borderRadius="full"
                        transition="all 0.3s"
                        _hover={{
                            bg: isChecked ? "blue.600" : "gray.400"
                        }}
                    >
                        <Box
                            position="absolute"
                            top="2px"
                            left={isChecked ? "18px" : "2px"}
                            w="16px"
                            h="16px"
                            bg="white"
                            borderRadius="full"
                            transition="all 0.3s"
                            boxShadow="md"
                        />
                    </Box>
                </Box>
            </Tooltip>
        );
    };

    // Render module name with special indicator for system modules
    const renderModuleName = (moduleName) => {
        const systemModules = ['User Management', 'Role & Permission', 'Module Management'];
        if (systemModules.includes(moduleName)) {
            return (
                <Text fontWeight="medium" fontSize="sm">
                    {moduleName}{' '}
                    <Text as="span" fontSize="xs" color="red.500" fontWeight="bold">
                        (System)
                    </Text>
                </Text>
            );
        }
        return <Text fontWeight="medium" fontSize="sm">{moduleName}</Text>;
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            size="xl"
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
                            <FormLabel fontWeight="semibold" fontSize="md" color={labelColor}>Role Name</FormLabel>
                            <Input
                                name="name"
                                value={formData.name}
                                onChange={handleChange}
                                placeholder="Enter role name"
                                bg={inputBg}
                                color={textColor}
                                borderWidth="1px"
                                height="44px"
                                fontSize="md"
                                _hover={{ borderColor: 'blue.400' }}
                                _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                            />
                            {errors.name && <FormErrorMessage fontWeight="medium">{errors.name}</FormErrorMessage>}
                        </FormControl>

                        <FormControl>
                            <FormLabel fontWeight="semibold" fontSize="md" color={labelColor}>Description (Optional)</FormLabel>
                            <Textarea
                                name="description"
                                value={formData.description}
                                onChange={handleChange}
                                placeholder="Enter role description"
                                bg={inputBg}
                                color={textColor}
                                borderWidth="1px"
                                fontSize="md"
                                _hover={{ borderColor: 'blue.400' }}
                                _focus={{ borderColor: 'blue.500', boxShadow: '0 0 0 1px var(--chakra-colors-blue-500)' }}
                                minHeight="100px"
                                resize="vertical"
                            />
                        </FormControl>

                        <Divider my={4} />

                        {/* Permissions Section */}
                        <Box>
                            <Text fontSize="lg" fontWeight="semibold" color={labelColor} mb={4}>Module Permissions</Text>

                            {loadingModules ? (
                                <Flex justify="center" align="center" py={4}>
                                    <Spinner size="md" thickness="3px" color="blue.500" mr={3} />
                                    <Text color="gray.500">Loading module permissions...</Text>
                                </Flex>
                            ) : (
                                <Box
                                    borderWidth="1px"
                                    borderRadius="xl"
                                    borderColor={borderColor}
                                    overflow="hidden"
                                    maxH="300px"
                                    mb={2}
                                    boxShadow="sm"
                                >
                                    <Table variant="simple" size="md">
                                        <Thead bg={useColorModeValue('gray.50', 'gray.700')} position="sticky" top={0} zIndex={1}>
                                            <Tr>
                                                <Th fontSize="xs" fontWeight="bold" py={4} pl={6} width="30%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    MODULE
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="14%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    READ
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="14%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    CREATE
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="14%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    UPDATE
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="14%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    DELETE
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="14%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    APPROVE
                                                </Th>
                                            </Tr>
                                        </Thead>
                                        <Tbody>
                                            {modules.map((module, index) => (
                                                <Tr
                                                    key={`module-${module.id}`}
                                                    _hover={{ bg: hoverBgColor }}
                                                    bg={index % 2 === 0 ? 'transparent' : useColorModeValue('gray.50', 'gray.800')}
                                                    borderBottomWidth={index === modules.length - 1 ? "0" : "1px"}
                                                    borderColor={borderColor}
                                                >
                                                    <Td py={3} pl={6}>{renderModuleName(module.name)}</Td>
                                                    <Td textAlign="center" py={3}>
                                                        <PermissionSwitch
                                                            isChecked={module.read}
                                                            onChange={() => handleTogglePermission(module.id, 'read')}
                                                        />
                                                    </Td>
                                                    <Td textAlign="center" py={3}>
                                                        <PermissionSwitch
                                                            isChecked={module.create}
                                                            onChange={() => handleTogglePermission(module.id, 'create')}
                                                        />
                                                    </Td>
                                                    <Td textAlign="center" py={3}>
                                                        <PermissionSwitch
                                                            isChecked={module.update}
                                                            onChange={() => handleTogglePermission(module.id, 'update')}
                                                        />
                                                    </Td>
                                                    <Td textAlign="center" py={3}>
                                                        <PermissionSwitch
                                                            isChecked={module.delete}
                                                            onChange={() => handleTogglePermission(module.id, 'delete')}
                                                        />
                                                    </Td>
                                                    <Td textAlign="center" py={3}>
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
                        height="40px"
                        minWidth="120px"
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
                        height="40px"
                        minWidth="180px"
                        _hover={{
                            bg: "blue.600",
                            shadow: 'lg'
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