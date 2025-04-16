import React, {useEffect, useState} from 'react';
import {
    Badge,
    Box,
    Button,
    Flex,
    HStack,
    IconButton,
    Spinner,
    Switch,
    Table,
    Tbody,
    Td,
    Text,
    Th,
    Thead,
    Tooltip,
    Tr,
    useColorModeValue,
    useToast
} from '@chakra-ui/react';
import {FiInfo, FiLock, FiSave, FiX} from 'react-icons/fi';
import moduleService from "../../../../services/moduleService.js";
import permissionService from "../../../../services/permissionService.js";
import PermissionSwitch from "./PermissionSwitch.jsx";

const RolePermissionPanel = ({ role, onSave, onClose, isLoading = false, modulesList = [], permissionsList = [] }) => {
    const [modules, setModules] = useState([]);
    const [loadingModules, setLoadingModules] = useState(false);
    const [savingChanges, setSavingChanges] = useState(false);
    const [hasChanges, setHasChanges] = useState(false);

    // Toast for notifications
    const toast = useToast();

    // Theme colors
    const headerBgColor = useColorModeValue('gray.50', 'gray.900');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const hoverBgColor = useColorModeValue('blue.50', 'gray.700');
    const tableBorderColor = useColorModeValue('gray.100', 'gray.800');
    const bgColor = useColorModeValue('white', 'gray.800');
    const panelBgColor = useColorModeValue('white', 'gray.800');

    // Load modules and permissions
    useEffect(() => {
        if (role) {
            // If modulesList is provided, use it
            if (modulesList && modulesList.length > 0) {
                mapModulesFromProps();
            } else {
                // Otherwise fetch from API
                fetchModulesAndPermissions();
            }
        }
    }, [role, modulesList]);

    // Map modules from props instead of fetching
    const mapModulesFromProps = () => {
        if (!role || !modulesList) return;

        setLoadingModules(true);
        try {
            // Create the module permission structure
            const formattedModules = modulesList.map(module => {
                // Check if this module has permissions in the role
                const modulePermissions = role.permissions?.find(p => p.module_id === module.id);

                // Create a permissions object for this module
                const permissionObject = {
                    id: module.id,
                    name: module.name,
                    read: false,
                    create: false,
                    update: false,
                    delete: false,
                    approve: false,
                    reject: false
                };

                // If this role has permissions for this module, mark them as true
                if (modulePermissions) {
                    const permList = modulePermissions.permissions || [];

                    // Handle numeric permission IDs
                    if (permList.includes(1)) permissionObject.read = true;
                    if (permList.includes(2)) permissionObject.create = true;
                    if (permList.includes(3)) permissionObject.update = true;
                    if (permList.includes(4)) permissionObject.delete = true;
                    if (permList.includes(5)) permissionObject.approve = true;
                    if (permList.includes(6)) permissionObject.reject = true;
                }

                return permissionObject;
            });

            setModules(formattedModules);
        } catch (error) {
            console.error('Error mapping modules:', error);
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

    // Fetch modules and permissions data from API
    const fetchModulesAndPermissions = async () => {
        setLoadingModules(true);
        try {
            // Get all modules and permissions in parallel to reduce API calls
            const [modulesResponse, permissionsResponse] = await Promise.all([
                moduleService.getModules({ getAll: true }),
                permissionService.getPermissions({ getAll: true })
            ]);

            const allModules = modulesResponse.data || [];

            // Create the module permission structure
            const formattedModules = allModules.map(module => {
                // Check if this module has permissions in the role
                const modulePermissions = role.permissions?.find(p => p.module_id === module.id);

                // Create a permissions object for this module
                const permissionObject = {
                    id: module.id,
                    name: module.name,
                    read: false,
                    create: false,
                    update: false,
                    delete: false,
                    approve: false,
                    reject: false
                };

                // If this role has permissions for this module, mark them as true
                if (modulePermissions) {
                    const permList = modulePermissions.permissions || [];

                    // Handle numeric permission IDs
                    if (permList.includes(1)) permissionObject.read = true;
                    if (permList.includes(2)) permissionObject.create = true;
                    if (permList.includes(3)) permissionObject.update = true;
                    if (permList.includes(4)) permissionObject.delete = true;
                    if (permList.includes(5)) permissionObject.approve = true;
                    if (permList.includes(6)) permissionObject.reject = true;
                }

                return permissionObject;
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
        setHasChanges(true);
    };

    // Handle save permissions
    const handleSave = async () => {
        if (!role) return;

        setSavingChanges(true);
        try {
            // Format the modules with permissions according to the API schema
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

            // Format payload for the API - based on the Swagger documentation
            const permissionsPayload = {
                role_name: role.name, // Include role_name as required by API
                modules_permissions: modulesWithPermissions, // Using the correct field name from API docs
                description: role.description // Preserve existing description
            };

            // Call the save callback with the updated permissions
            await onSave(role.id, permissionsPayload);

            toast({
                title: 'Permissions updated',
                description: 'Role permissions have been updated successfully.',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });

            setHasChanges(false);
            onClose();
        } catch (error) {
            console.error('Error saving permissions:', error);
            toast({
                title: 'Update failed',
                description: 'Failed to update role permissions.',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setSavingChanges(false);
        }
    };

    // Render module name with lock icon for system modules
    const renderModuleName = (moduleName) => {
        const systemModules = ['User Management', 'Role & Permission', 'Module Management'];
        if (systemModules.includes(moduleName)) {
            return (
                <Tooltip label="Core system module" hasArrow placement="top">
                    <HStack spacing={1}>
                        <Text>{moduleName}</Text>
                        <FiLock size={14} color="gray" />
                    </HStack>
                </Tooltip>
            );
        }
        return moduleName;
    };

    return (
        <Box
            borderWidth="1px"
            borderColor={borderColor}
            borderRadius="lg"
            bg={panelBgColor}
            shadow="lg"
            mb={4}
            mt={2}
            position="relative"
            overflow="hidden"
        >
            {/* Header with title and close button */}
            <Flex
                bg={headerBgColor}
                px={4}
                py={3}
                alignItems="center"
                justifyContent="space-between"
                borderBottomWidth="1px"
                borderColor={borderColor}
            >
                <HStack>
                    <Text fontSize="md" fontWeight="bold">Quyền của vai trò</Text>
                    {role && (
                        <Badge colorScheme="blue" ml={2} fontSize="xs">
                            {role.name}
                        </Badge>
                    )}
                </HStack>
                <IconButton
                    icon={<FiX />}
                    size="sm"
                    variant="ghost"
                    onClick={onClose}
                    aria-label="Close permissions panel"
                />
            </Flex>

            {/* Content */}
            <Box p={4}>
                {(isLoading || loadingModules) ? (
                    <Flex justify="center" align="center" direction="column" py={10}>
                        <Spinner size="xl" color="blue.500" mb={4} />
                        <Text color="gray.500">Đang tải quyền...</Text>
                    </Flex>
                ) : (
                    <Box>
                        <Box
                            borderWidth="1px"
                            borderColor={borderColor}
                            borderRadius="md"
                            overflow="auto"
                            maxH="400px"
                            mb={4}
                        >
                            <Table variant="simple" size="sm">
                                <Thead bg={headerBgColor} position="sticky" top={0} zIndex={1}>
                                    <Tr>
                                        <Th fontSize="xs" py={3} width="30%">Module</Th>
                                        <Th width="14%" textAlign="center" fontSize="xs" py={3}>Read</Th>
                                        <Th width="14%" textAlign="center" fontSize="xs" py={3}>Create</Th>
                                        <Th width="14%" textAlign="center" fontSize="xs" py={3}>Update</Th>
                                        <Th width="14%" textAlign="center" fontSize="xs" py={3}>Delete</Th>
                                        <Th width="14%" textAlign="center" fontSize="xs" py={3}>Approve</Th>
                                    </Tr>
                                </Thead>
                                <Tbody>
                                    {modules.length === 0 ? (
                                        <Tr>
                                            <Td colSpan={6} textAlign="center" py={4}>
                                                <Flex direction="column" align="center" justify="center" py={4}>
                                                    <Box color="gray.400" mb={2}>
                                                        <FiInfo size={24} />
                                                    </Box>
                                                    <Text fontWeight="normal" color="gray.500" fontSize="sm">Hiện chưa có modules nào</Text>
                                                </Flex>
                                            </Td>
                                        </Tr>
                                    ) : (
                                        modules.map((module, index) => (
                                            <Tr
                                                key={module.id}
                                                _hover={{ bg: hoverBgColor }}
                                                bg={index % 2 === 0 ? bgColor : useColorModeValue('gray.50', 'gray.800')}
                                                borderBottom="1px"
                                                borderColor={tableBorderColor}
                                            >
                                                <Td fontWeight="medium" fontSize="sm" py={2}>
                                                    {renderModuleName(module.name)}
                                                </Td>
                                                <Td textAlign="center" py={2}>
                                                    <PermissionSwitch
                                                        isChecked={module.read}
                                                        onChange={() => handleTogglePermission(module.id, 'read')}
                                                        permission="read"
                                                    />
                                                </Td>
                                                <Td textAlign="center" py={2}>
                                                    <PermissionSwitch
                                                        isChecked={module.create}
                                                        onChange={() => handleTogglePermission(module.id, 'create')}
                                                        permission="create"
                                                    />
                                                </Td>
                                                <Td textAlign="center" py={2}>
                                                    <PermissionSwitch
                                                        isChecked={module.update}
                                                        onChange={() => handleTogglePermission(module.id, 'update')}
                                                        permission="update"
                                                    />
                                                </Td>
                                                <Td textAlign="center" py={2}>
                                                    <PermissionSwitch
                                                        isChecked={module.delete}
                                                        onChange={() => handleTogglePermission(module.id, 'delete')}
                                                        permission="delete"
                                                    />
                                                </Td>
                                                <Td textAlign="center" py={2}>
                                                    <PermissionSwitch
                                                        isChecked={module.approve}
                                                        onChange={() => handleTogglePermission(module.id, 'approve')}
                                                        permission="approve"
                                                    />
                                                </Td>
                                            </Tr>
                                        ))
                                    )}
                                </Tbody>
                            </Table>
                        </Box>

                        {/* Help text */}
                        <Text fontSize="xs" color="gray.500" mb={4}>
                            <FiInfo size={14} style={{ display: 'inline', marginRight: '4px', verticalAlign: 'middle' }} />
                            Thay đổi quyền sẽ được áp dụng sau khi bạn lưu.
                        </Text>

                        {/* Footer buttons */}
                        <Flex justifyContent="flex-end" mt={2}>
                            <Button
                                variant="outline"
                                size="sm"
                                mr={2}
                                onClick={onClose}
                            >
                                Huỷ
                            </Button>
                            <Button
                                colorScheme="blue"
                                size="sm"
                                isLoading={savingChanges}
                                loadingText="Saving"
                                onClick={handleSave}
                                leftIcon={<FiSave size={14} />}
                                isDisabled={!hasChanges}
                            >
                                Lưu thay đổi
                            </Button>
                        </Flex>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

export default RolePermissionPanel;