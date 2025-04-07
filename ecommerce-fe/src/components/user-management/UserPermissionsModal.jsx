import React, { useState, useEffect } from 'react';
import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    Button,
    VStack,
    HStack,
    Table,
    Thead,
    Tbody,
    Tr,
    Th,
    Td,
    Checkbox,
    Text,
    useToast,
    Box,
    Flex,
    Divider,
    Badge,
    Skeleton,
    IconButton,
    useColorModeValue,
} from '@chakra-ui/react';
import { FiSave, FiCheck, FiX } from 'react-icons/fi';

const UserPermissionsModal = ({ isOpen, onClose, user, onSave, isLoading }) => {
    const [permissions, setPermissions] = useState([]);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const toast = useToast();

    // Set up initial permissions when user data is loaded
    useEffect(() => {
        if (user && user.module_permission) {
            setPermissions(user.module_permission);
        }
    }, [user]);

    // Toggle a specific permission for a module
    const togglePermission = (moduleId, permissionId) => {
        setPermissions(prevPermissions => {
            return prevPermissions.map(module => {
                if (module.module_id === moduleId) {
                    // Toggle the permission
                    const updatedPermissions = module.permissions.map(permission => {
                        if (permission.permission_id === permissionId) {
                            // Create a new "active" property or toggle the existing one
                            return {
                                ...permission,
                                active: permission.active === undefined ? true : !permission.active
                            };
                        }
                        return permission;
                    });

                    return {
                        ...module,
                        permissions: updatedPermissions
                    };
                }
                return module;
            });
        });
    };

    // Check if all permissions for a module are active
    const areAllPermissionsActive = (modulePermissions) => {
        return modulePermissions.every(permission => permission.active);
    };

    // Toggle all permissions for a module
    const toggleAllModulePermissions = (moduleId, newValue) => {
        setPermissions(prevPermissions => {
            return prevPermissions.map(module => {
                if (module.module_id === moduleId) {
                    const updatedPermissions = module.permissions.map(permission => ({
                        ...permission,
                        active: newValue
                    }));

                    return {
                        ...module,
                        permissions: updatedPermissions
                    };
                }
                return module;
            });
        });
    };

    // Handle save permissions
    const handleSavePermissions = async () => {
        if (!user || !user.id) return;

        setIsSubmitting(true);
        try {
            // Format permissions for API
            const formattedPermissions = permissions.flatMap(module =>
                module.permissions
                    .filter(permission => permission.active)
                    .map(permission => ({
                        module_id: module.module_id,
                        permission_id: permission.permission_id
                    }))
            );

            await onSave(user.id, formattedPermissions);

            toast({
                title: "Permissions updated",
                description: "User permissions have been successfully updated.",
                status: "success",
                duration: 5000,
                isClosable: true,
            });

            onClose();
        } catch (error) {
            toast({
                title: "Error updating permissions",
                description: error.message || "An error occurred while updating permissions.",
                status: "error",
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const bgColor = useColorModeValue('white', 'gray.800');
    const headerBgColor = useColorModeValue('gray.50', 'gray.900');

    if (!user) {
        return null;
    }

    return (
        <Modal isOpen={isOpen} onClose={onClose} size="xl" scrollBehavior="inside">
            <ModalOverlay bg="blackAlpha.300" backdropFilter="blur(5px)" />
            <ModalContent borderRadius="xl" shadow="xl">
                <ModalHeader borderBottomWidth="1px" borderColor={borderColor} py={4}>
                    <HStack spacing={3}>
                        <Text>User Permissions</Text>
                        <Badge colorScheme="blue">{user.fullname}</Badge>
                    </HStack>
                </ModalHeader>
                <ModalCloseButton />

                <ModalBody py={6}>
                    {isLoading ? (
                        <VStack spacing={4}>
                            {[1, 2, 3].map(i => (
                                <Skeleton key={i} height="60px" width="100%" />
                            ))}
                        </VStack>
                    ) : (
                        <Box overflowX="auto">
                            <Table variant="simple" size="sm">
                                <Thead bg={headerBgColor}>
                                    <Tr>
                                        <Th width="30%">Module</Th>
                                        <Th>Permissions</Th>
                                        <Th width="80px" textAlign="center">All</Th>
                                    </Tr>
                                </Thead>
                                <Tbody>
                                    {permissions.length > 0 ? (
                                        permissions.map((module) => (
                                            <Tr key={module.module_id}>
                                                <Td fontWeight="medium">{module.module_name}</Td>
                                                <Td>
                                                    <Flex flexWrap="wrap" gap={4}>
                                                        {module.permissions.map((permission) => (
                                                            <HStack key={permission.permission_id} spacing={2}>
                                                                <Checkbox
                                                                    isChecked={permission.active}
                                                                    onChange={() => togglePermission(module.module_id, permission.permission_id)}
                                                                    colorScheme="green"
                                                                >
                                                                    {permission.permission_name}
                                                                </Checkbox>
                                                            </HStack>
                                                        ))}
                                                    </Flex>
                                                </Td>
                                                <Td textAlign="center">
                                                    <Checkbox
                                                        isChecked={areAllPermissionsActive(module.permissions)}
                                                        onChange={(e) => toggleAllModulePermissions(module.module_id, e.target.checked)}
                                                        colorScheme="green"
                                                    />
                                                </Td>
                                            </Tr>
                                        ))
                                    ) : (
                                        <Tr>
                                            <Td colSpan={3} textAlign="center" py={6}>
                                                <Text color="gray.500">No permissions available</Text>
                                            </Td>
                                        </Tr>
                                    )}
                                </Tbody>
                            </Table>
                        </Box>
                    )}
                </ModalBody>

                <ModalFooter borderTopWidth="1px" borderColor={borderColor} py={4}>
                    <Button
                        variant="outline"
                        mr={3}
                        onClick={onClose}
                        isDisabled={isSubmitting}
                    >
                        Cancel
                    </Button>
                    <Button
                        colorScheme="blue"
                        onClick={handleSavePermissions}
                        isLoading={isSubmitting}
                        leftIcon={<FiSave />}
                    >
                        Save Changes
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default UserPermissionsModal;