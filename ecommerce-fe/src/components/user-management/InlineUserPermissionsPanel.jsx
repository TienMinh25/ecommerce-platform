import React, { useState, useEffect } from 'react';
import {
    Box,
    Text,
    Heading,
    Table,
    Thead,
    Tbody,
    Tr,
    Th,
    Td,
    Switch,
    Badge,
    Button,
    Flex,
    Tooltip,
    HStack,
    Spinner,
    useColorModeValue,
    Collapse,
    IconButton,
    Divider
} from '@chakra-ui/react';
import { FiLock, FiInfo, FiX } from 'react-icons/fi';
import moduleService from "../../services/moduleService.js";
import permissionService from "../../services/permissionService.js";

const InlineUserPermissionsPanel = ({ user, onSave, onClose, isLoading = false }) => {
    const [modules, setModules] = useState([]);
    const [loadingModules, setLoadingModules] = useState(false);
    const [savingChanges, setSavingChanges] = useState(false);

    // Màu sắc
    const headerBgColor = useColorModeValue('gray.50', 'gray.900');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const hoverBgColor = useColorModeValue('blue.50', 'gray.700');
    const tableBorderColor = useColorModeValue('gray.100', 'gray.800');
    const bgColor = useColorModeValue('white', 'gray.800');
    const panelBgColor = useColorModeValue('white', 'gray.800');

    // Tạo danh sách modules với các quyền từ dữ liệu người dùng
    useEffect(() => {
        if (user) {
            setLoadingModules(true);

            // Lấy tất cả modules và permissions
            Promise.all([
                moduleService.getModules(1, 100, true),
                permissionService.getPermissions(1, 100, true)
            ])
                .then(([modulesResponse, permissionsResponse]) => {
                    const allModules = modulesResponse.data;

                    // Tạo cấu trúc dữ liệu modules với các quyền
                    const formattedModules = allModules.map(module => {
                        // Tìm module tương ứng trong quyền của người dùng
                        const userModule = user.module_permission?.find(m => m.module_id === module.id);

                        // Mặc định các quyền là false
                        const moduleData = {
                            id: module.id,
                            name: module.name,
                            read: false,
                            create: false,
                            update: false,
                            delete: false,
                            approve: false,
                            reject: false
                        };

                        // Nếu người dùng có quyền cho module này, cập nhật trạng thái
                        if (userModule) {
                            userModule.permissions.forEach(permission => {
                                const permName = permission.permission_name;
                                if (permName in moduleData) {
                                    moduleData[permName] = true;
                                }
                            });
                        }

                        return moduleData;
                    });

                    setModules(formattedModules);
                })
                .catch(error => {
                    console.error('Error loading modules and permissions:', error);
                })
                .finally(() => {
                    setLoadingModules(false);
                });
        }
    }, [user]);

    // Xử lý thay đổi quyền
    const handleTogglePermission = (moduleId, permission) => {
        setModules(modules.map(module =>
            module.id === moduleId
                ? { ...module, [permission]: !module[permission] }
                : module
        ));
    };

    // Xử lý lưu thay đổi
    const handleSave = async () => {
        if (!user) return;

        setSavingChanges(true);

        try {
            // Định dạng dữ liệu để gửi lên server
            const permissionsData = modules.map(module => {
                const permissions = [];

                if (module.read) permissions.push({ permission_name: 'read' });
                if (module.create) permissions.push({ permission_name: 'create' });
                if (module.update) permissions.push({ permission_name: 'update' });
                if (module.delete) permissions.push({ permission_name: 'delete' });
                if (module.approve) permissions.push({ permission_name: 'approve' });
                if (module.reject) permissions.push({ permission_name: 'reject' });

                return {
                    module_id: module.id,
                    module_name: module.name,
                    permissions
                };
            }).filter(module => module.permissions.length > 0);

            // Gọi callback onSave và truyền dữ liệu
            await onSave(user.id, permissionsData);
            onClose();
        } catch (error) {
            console.error('Error saving permissions:', error);
        } finally {
            setSavingChanges(false);
        }
    };

    // Hiển thị Switch với tooltip
    const PermissionSwitch = ({ isChecked, onChange, permission, isDisabled = false }) => {
        // Lấy text tooltip dựa trên loại quyền
        const getTooltipText = () => {
            switch(permission) {
                case 'read': return 'View permission';
                case 'create': return 'Create permission';
                case 'update': return 'Edit/Update permission';
                case 'delete': return 'Delete permission';
                case 'approve': return 'Approve permission';
                case 'reject': return 'Reject permission';
                default: return 'Toggle permission';
            }
        };

        return (
            <Tooltip label={getTooltipText()} hasArrow placement="top">
                <Box
                    position="relative"
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    cursor={isDisabled ? "not-allowed" : "pointer"}
                    onClick={isDisabled ? undefined : onChange}
                    borderRadius="md"
                    p={0.5}
                    opacity={isDisabled ? 0.6 : 1}
                >
                    <Switch
                        colorScheme="blue"
                        size="sm"
                        isChecked={isChecked}
                        isDisabled={isDisabled}
                    />
                </Box>
            </Tooltip>
        );
    };

    // Hiển thị tên module với icon khóa cho các module hệ thống
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
            {/* Header với tiêu đề và nút đóng */}
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
                    <Heading size="sm">User Permissions</Heading>
                    {user && (
                        <Badge colorScheme="blue" ml={2} fontSize="xs">
                            {user.fullname}
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

            {/* Nội dung */}
            <Box p={4}>
                {(isLoading || loadingModules) ? (
                    <Flex justify="center" align="center" direction="column" py={10}>
                        <Spinner size="xl" color="blue.500" mb={4} />
                        <Text color="gray.500">Loading permissions...</Text>
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
                                                    <Text fontWeight="normal" color="gray.500" fontSize="sm">No module permissions found</Text>
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

                        {/* Thông tin hướng dẫn */}
                        <Text fontSize="xs" color="gray.500" mb={4}>
                            <FiInfo size={14} style={{ display: 'inline', marginRight: '4px', verticalAlign: 'middle' }} />
                            Changes to permissions will take effect immediately after saving.
                        </Text>

                        {/* Footer buttons */}
                        <Flex justifyContent="flex-end" mt={2}>
                            <Button
                                variant="outline"
                                size="sm"
                                mr={2}
                                onClick={onClose}
                            >
                                Cancel
                            </Button>
                            <Button
                                colorScheme="blue"
                                size="sm"
                                isLoading={savingChanges}
                                loadingText="Saving"
                                onClick={handleSave}
                            >
                                Save Changes
                            </Button>
                        </Flex>
                    </Box>
                )}
            </Box>
        </Box>
    );
};

export default InlineUserPermissionsPanel;