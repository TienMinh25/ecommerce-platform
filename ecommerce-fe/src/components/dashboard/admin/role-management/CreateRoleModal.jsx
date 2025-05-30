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
    Spinner,
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
import {FiInfo, FiPlus, FiShield} from 'react-icons/fi';
import roleService from '../../../../services/roleService.js';
import moduleService from "../../../../services/moduleService.js";
import permissionService from "../../../../services/permissionService.js";
import PermissionSwitch from './PermissionSwitch'; // Import the reusable component

const CreateRoleModal = ({ isOpen, onClose, onRoleCreated, modulesList = [], permissionsList = [] }) => {
    // Theme colors
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const inputBg = useColorModeValue('white', 'gray.900');
    const headerBg = useColorModeValue('gray.50', 'gray.900');
    const textColor = useColorModeValue('gray.900', 'white');
    const labelColor = useColorModeValue('gray.800', 'gray.100');
    const iconColor = useColorModeValue('blue.700', 'blue.200');
    const hoverBgColor = useColorModeValue('blue.50', 'gray.700');
    const bgColor = useColorModeValue('white', 'gray.800');

    // Form state
    const [formData, setFormData] = useState({
        name: '',
        description: '',
    });

    // Modules and permissions state
    const [modules, setModules] = useState([]);
    const [loadingModules, setLoadingModules] = useState(false);
    const [hasChanges, setHasChanges] = useState(false);
    const [permissionMapping, setPermissionMapping] = useState({});

    // Validation and submission state
    const [errors, setErrors] = useState({});
    const [isSubmitting, setIsSubmitting] = useState(false);

    const toast = useToast();

    // Create permission mapping from permissionsList
    useEffect(() => {
        if (permissionsList && permissionsList.length > 0) {
            const mapping = {};
            permissionsList.forEach(permission => {
                mapping[permission.name] = permission.id;
            });
            setPermissionMapping(mapping);
        }
    }, [permissionsList]);

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

        setHasChanges(true);
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

    // Load modules and permissions when modal opens
    useEffect(() => {
        if (isOpen && Object.keys(permissionMapping).length > 0) {
            setHasChanges(false);

            // Use provided modules and permissions or fetch them
            if (modulesList && modulesList.length > 0) {
                prepareModules(modulesList);
            } else {
                fetchModulesAndPermissions();
            }
        }
    }, [isOpen, modulesList, permissionMapping]);

    // Prepare modules from props
    const prepareModules = (modules) => {
        // Create module objects with all permissions set to false
        const formattedModules = modules.map(module => ({
            id: module.id,
            name: module.name,
            read: false,
            create: false,
            update: false,
            delete: false,
            approve: false,
            reject: false
        }));

        setModules(formattedModules);
    };

    // Fetch modules and permissions data from API
    const fetchModulesAndPermissions = async () => {
        setLoadingModules(true);

        try {
            // Get all modules and permissions in parallel
            const [modulesResponse, permissionsResponse] = await Promise.all([
                moduleService.getModules({ getAll: true }),
                permissionService.getPermissions({ getAll: true })
            ]);

            const allModules = modulesResponse.data || [];
            const allPermissions = permissionsResponse.data || [];

            // Create permission mapping if not provided in props
            if (Object.keys(permissionMapping).length === 0) {
                const mapping = {};
                allPermissions.forEach(permission => {
                    mapping[permission.name] = permission.id;
                });
                setPermissionMapping(mapping);
            }

            // Prepare modules
            prepareModules(allModules);
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
        setHasChanges(true);
    };

    // Reset form when modal closes
    const resetForm = () => {
        setFormData({
            name: '',
            description: '',
        });
        setErrors({});
        setHasChanges(false);
    };

    useEffect(() => {
        if (!isOpen) {
            resetForm();
        }
    }, [isOpen]);

    // Handle form submission
    const handleSubmit = async () => {
        if (!validateForm() || Object.keys(permissionMapping).length === 0) return;

        setIsSubmitting(true);
        try {
            // Format modules with permissions for API
            const modules_permissions = modules
                .filter(module =>
                    module.read || module.create || module.update ||
                    module.delete || module.approve || module.reject
                )
                .map(module => {
                    // Map permissions to their correct IDs using dynamic mapping
                    const permissionIds = [];
                    if (module.read) permissionIds.push(permissionMapping.read);
                    if (module.create) permissionIds.push(permissionMapping.create);
                    if (module.update) permissionIds.push(permissionMapping.update);
                    if (module.delete) permissionIds.push(permissionMapping.delete);
                    if (module.approve) permissionIds.push(permissionMapping.approve);
                    if (module.reject) permissionIds.push(permissionMapping.reject);

                    return {
                        module_id: module.id,
                        permissions: permissionIds
                    };
                });

            // Create the role data according to the API schema
            const roleData = {
                role_name: formData.name.trim(),
                description: formData.description.trim(),
                modules_permissions: modules_permissions
            };

            // Create the role with all data in one request
            await roleService.createRole(roleData);

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
            size="5xl" // Wide modal
            motionPreset="slideInBottom"
            scrollBehavior="inside"
            isCentered
        >
            <ModalOverlay backdropFilter="blur(3px)" bg="blackAlpha.400" />
            <ModalContent
                borderRadius="xl"
                shadow="2xl"
                bg={useColorModeValue('white', 'gray.800')}
                maxHeight="90vh"
                overflowY="auto"
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
                    <Text fontSize="xl" fontWeight="bold" color={textColor}>Tạo vai trò mới</Text>
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
                            <FormLabel fontWeight="semibold" fontSize="md" color={labelColor}>Tên vai trò</FormLabel>
                            <Input
                                name="name"
                                value={formData.name}
                                onChange={handleChange}
                                placeholder="Nhập tên vai trò"
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
                            <FormLabel fontWeight="semibold" fontSize="md" color={labelColor}>Mô tả (Tuỳ chọn)</FormLabel>
                            <Textarea
                                name="description"
                                value={formData.description}
                                onChange={handleChange}
                                placeholder="Nhập mô tả"
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
                            <Text fontSize="lg" fontWeight="semibold" color={labelColor} mb={4}>Quyền trên modules</Text>

                            {loadingModules ? (
                                <Flex justify="center" align="center" py={4}>
                                    <Spinner size="md" thickness="3px" color="blue.500" mr={3} />
                                    <Text color="gray.500">Đang tải dữ liệu quyền của modules...</Text>
                                </Flex>
                            ) : (
                                <Box
                                    borderWidth="1px"
                                    borderRadius="xl"
                                    borderColor={borderColor}
                                    overflow="auto"
                                    maxH="400px"
                                    mb={2}
                                    boxShadow="sm"
                                >
                                    <Table variant="simple" size="md">
                                        <Thead bg={useColorModeValue('gray.50', 'gray.700')} position="sticky" top={0} zIndex={1}>
                                            <Tr>
                                                <Th fontSize="xs" fontWeight="bold" py={4} pl={6} width="30%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    MODULE
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="11.6%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    read
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="11.6%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    CREATE
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="11.6%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    UPDATE
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="11.6%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    DELETE
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="11.6%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    APPROVE
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" fontWeight="bold" py={4} width="11.6%" borderBottomWidth="2px" borderColor={borderColor}>
                                                    REJECT
                                                </Th>
                                            </Tr>
                                        </Thead>
                                        <Tbody>
                                            {modules.length > 0 ? (
                                                modules.map((module, index) => (
                                                    <Tr
                                                        key={`module-${module.id}`}
                                                        _hover={{ bg: hoverBgColor }}
                                                        bg={index % 2 === 0 ? bgColor : useColorModeValue('gray.50', 'gray.800')}
                                                        borderBottomWidth={index === modules.length - 1 ? "0" : "1px"}
                                                        borderColor={borderColor}
                                                    >
                                                        <Td py={3} pl={6}>{renderModuleName(module.name)}</Td>
                                                        <Td textAlign="center" py={3}>
                                                            <PermissionSwitch
                                                                isChecked={module.read}
                                                                onChange={() => handleTogglePermission(module.id, 'read')}
                                                                permission="read"
                                                            />
                                                        </Td>
                                                        <Td textAlign="center" py={3}>
                                                            <PermissionSwitch
                                                                isChecked={module.create}
                                                                onChange={() => handleTogglePermission(module.id, 'create')}
                                                                permission="create"
                                                            />
                                                        </Td>
                                                        <Td textAlign="center" py={3}>
                                                            <PermissionSwitch
                                                                isChecked={module.update}
                                                                onChange={() => handleTogglePermission(module.id, 'update')}
                                                                permission="update"
                                                            />
                                                        </Td>
                                                        <Td textAlign="center" py={3}>
                                                            <PermissionSwitch
                                                                isChecked={module.delete}
                                                                onChange={() => handleTogglePermission(module.id, 'delete')}
                                                                permission="delete"
                                                            />
                                                        </Td>
                                                        <Td textAlign="center" py={3}>
                                                            <PermissionSwitch
                                                                isChecked={module.approve}
                                                                onChange={() => handleTogglePermission(module.id, 'approve')}
                                                                permission="approve"
                                                            />
                                                        </Td>
                                                        <Td textAlign="center" py={3}>
                                                            <PermissionSwitch
                                                                isChecked={module.reject}
                                                                onChange={() => handleTogglePermission(module.id, 'reject')}
                                                                permission="reject"
                                                            />
                                                        </Td>
                                                    </Tr>
                                                ))
                                            ) : (
                                                <Tr>
                                                    <Td colSpan={7} textAlign="center" py={4}>
                                                        <Flex direction="column" align="center" justify="center" py={4}>
                                                            <Box color="gray.400" mb={2}>
                                                                <FiInfo size={24} />
                                                            </Box>
                                                            <Text color="gray.500">Không có modules khả dụng. Vui lòng thêm module trước khi tiếp tục.</Text>
                                                        </Flex>
                                                    </Td>
                                                </Tr>
                                            )}
                                        </Tbody>
                                    </Table>
                                </Box>
                            )}
                            <Text fontSize="xs" color="gray.500" mt={2}>
                                <FiInfo size={14} style={{ display: 'inline', marginRight: '4px', verticalAlign: 'middle' }} />
                                Thiết lập quyền cho vai trò này. Quyền hạn quyết định những hành động mà người dùng với vai trò này có thể thực hiện.
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
                        Huỷ
                    </Button>
                    <Button
                        leftIcon={<FiPlus />}
                        colorScheme="blue"
                        onClick={handleSubmit}
                        isLoading={isSubmitting}
                        loadingText="Creating..."
                        px={8}
                        shadow="md"
                        height="40px"
                        minWidth="180px"
                        _hover={{
                            bg: "blue.600",
                            shadow: 'lg'
                        }}
                        fontWeight="bold"
                        isDisabled={loadingModules || !hasChanges}
                    >
                        Tạo vai trò
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default CreateRoleModal;