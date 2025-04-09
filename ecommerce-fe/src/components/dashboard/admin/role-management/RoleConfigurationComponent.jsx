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
  FormControl,
  FormLabel,
  Input,
  Box,
  Text,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  IconButton,
  HStack,
  Switch,
  Flex,
  FormHelperText,
  Tooltip,
  Badge,
  useToast,
  useColorModeValue
} from '@chakra-ui/react';
import { FiTrash2, FiAlertCircle, FiLock } from 'react-icons/fi';

const RoleConfigurationComponent = ({ 
  isOpen, 
  onClose, 
  roleId = null, 
  onSave, 
  modalSize = "xl",
  disableRoleNameEdit = false  // Thêm prop mới để kiểm soát việc disable input
}) => {
  const toast = useToast();
  
  // State for role name
  const [roleName, setRoleName] = useState('');
  const [isRoleNameTouched, setIsRoleNameTouched] = useState(false);
  
  // State for permission modules
  const [modules, setModules] = useState([
    { id: 1, name: 'Dashboard', read: true, create: false, update: false, delete: false, approve: false },
    { id: 2, name: 'User Management', read: true, create: true, update: true, delete: false, approve: false },
    { id: 3, name: 'Reports', read: true, create: false, update: false, delete: false, approve: false },
  ]);
  
  // State variables for Add Module functionality have been removed
  
  // Colors
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const bgColor = useColorModeValue('white', 'gray.800');
  const hoverBgColor = useColorModeValue('gray.50', 'gray.700');
  const headerBgColor = useColorModeValue('gray.50', 'gray.900');
  const requiredColor = useColorModeValue('red.500', 'red.300');
  const activeColor = useColorModeValue('blue.50', 'blue.900');
  const tableBorderColor = useColorModeValue('gray.100', 'gray.800');
  
  // Reset form when modal opens or roleId changes
  useEffect(() => {
    if (isOpen) {
      if (roleId) {
        // If editing existing role, fetch role data
        setRoleName(`Role ${roleId}`);
      } else {
        // If creating new role, reset form
        setRoleName('');
        setModules([
          { id: 1, name: 'Dashboard', read: true, create: false, update: false, delete: false, approve: false },
          { id: 2, name: 'User Management', read: true, create: true, update: true, delete: false, approve: false },
          { id: 3, name: 'Reports', read: true, create: false, update: false, delete: false, approve: false },
        ]);
      }
      setIsRoleNameTouched(false);
    }
  }, [isOpen, roleId]);
  
  // Validate form - nếu disable chỉnh sửa tên thì luôn xem như hợp lệ
  const isRoleNameValid = disableRoleNameEdit ? true : roleName.trim() !== '';
  const formIsValid = isRoleNameValid;
  
  // Handle toggle permission
  const handleTogglePermission = (moduleId, permission) => {
    setModules(modules.map(module => 
      module.id === moduleId 
        ? { ...module, [permission]: !module[permission] } 
        : module
    ));
  };
  
  // Handle delete module
  const handleDeleteModule = (moduleId) => {
    setModules(modules.filter(module => module.id !== moduleId));
  };
  
  // The handleAddModule function has been removed
  
  // Handle save
  const handleSave = () => {
    if (!formIsValid) {
      setIsRoleNameTouched(true);
      toast({
        title: "Validation Error",
        description: "Please fill in all required fields.",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
      return;
    }
    
    // Prepare data to save
    const roleData = {
      name: roleName,
      permissions: modules.map(({ id, name, ...permissions }) => ({
        moduleId: id,
        moduleName: name,
        ...permissions
      }))
    };
    
    // Call onSave function passed from parent
    if (onSave) {
      onSave(roleData);
    }
    
    // Show success message
    toast({
      title: roleId ? "Role Updated" : "Role Created",
      description: `Role "${roleName}" has been ${roleId ? "updated" : "created"} successfully.`,
      status: "success",
      duration: 3000,
      isClosable: true,
    });
    
    // Close modal
    onClose();
  };
  
  // Custom Switch component with better styling
  const PermissionSwitch = ({ isChecked, onChange, permission }) => {
    const activeBoxShadow = useColorModeValue('0 0 0 2px rgba(49, 130, 206, 0.6)', '0 0 0 2px rgba(66, 153, 225, 0.6)');
    const switchTrackBg = useColorModeValue('gray.300', 'gray.600');
    const checkedBg = useColorModeValue('blue.500', 'blue.400');
    
    // Get tooltip text based on permission
    const getTooltipText = () => {
      switch(permission) {
        case 'read': return 'View permission';
        case 'create': return 'Create permission';
        case 'update': return 'Edit/Update permission';
        case 'delete': return 'Delete permission';
        case 'approve': return 'Approve/Reject permission';
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
          cursor="pointer"
          onClick={onChange}
          borderRadius="md"
          p={0.5}
          transition="all 0.2s"
          bg={isChecked ? activeColor : 'transparent'}
          _hover={{ 
            bg: isChecked ? activeColor : hoverBgColor
          }}
        >
          <Switch
            colorScheme="blue"
            size="sm"
            isChecked={isChecked}
            sx={{
              '& .chakra-switch__track': {
                bg: isChecked ? checkedBg : switchTrackBg
              },
              '&:focus-visible': {
                boxShadow: activeBoxShadow
              }
            }}
          />
        </Box>
      </Tooltip>
    );
  };
  
  // Function to render module name (refactored to solve the JSX conditional rendering issue)
  const renderModuleName = (moduleName) => {
    if (['Dashboard', 'User Management'].includes(moduleName)) {
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
    <Modal 
      isOpen={isOpen} 
      onClose={onClose}
      size={modalSize}
      scrollBehavior="inside"
      motionPreset="slideInBottom"
      isCentered
    >
      <ModalOverlay bg="blackAlpha.300" backdropFilter="blur(5px)" />
      <ModalContent 
        borderRadius="lg" 
        shadow="lg"
        mx={{ base: 3, md: 6 }}
        maxW="80vw"
        h={{ base: "85vh", md: "80vh" }}
        maxH="80vh"
        w="full"
      >
        <ModalHeader
          bg={headerBgColor}
          borderBottom="1px"
          borderColor={borderColor}
          py={3}
          px={4}
          borderTopRadius="lg"
          fontWeight="bold"
          display="flex"
          alignItems="center"
          position="sticky"
          top="0"
          zIndex="1"
          fontSize="md"
        >
          {roleId ? (
            <>
              <Text>EDIT PERMISSIONS</Text>
              <Badge ml={2} colorScheme="blue" variant="subtle" fontSize="2xs" borderRadius="md" px={2}>
                ID: {roleId}
              </Badge>
            </>
          ) : (
            <Text>NEW ROLE</Text>
          )}
        </ModalHeader>
        <ModalCloseButton size="sm" top={3} right={4} />
        
        <ModalBody pt={4} px={{ base: 3, md: 6 }} overflowY="auto">
          {/* Role Name Input - cập nhật để hỗ trợ disable */}
          <FormControl isRequired={!disableRoleNameEdit} isInvalid={isRoleNameTouched && !isRoleNameValid} mb={6}>
            <FormLabel fontWeight="semibold" fontSize="sm">
              Role Name
              {!disableRoleNameEdit && (
                <Text as="span" color={requiredColor} ml={1}>*</Text>
              )}
            </FormLabel>
            <Input
              placeholder="Enter role name"
              value={roleName}
              onChange={(e) => setRoleName(e.target.value)}
              onBlur={() => setIsRoleNameTouched(true)}
              size="md"
              focusBorderColor="blue.400"
              bg={bgColor}
              borderColor={borderColor}
              _hover={{ borderColor: disableRoleNameEdit ? borderColor : 'gray.300' }}
              borderRadius="md"
              shadow="sm"
              fontSize="sm"
              isDisabled={disableRoleNameEdit}
              _disabled={{
                opacity: 0.8,
                cursor: 'not-allowed',
                bg: useColorModeValue('gray.100', 'gray.700')
              }}
            />
            {isRoleNameTouched && !isRoleNameValid && !disableRoleNameEdit && (
              <FormHelperText color={requiredColor} mt={1}>
                <HStack spacing={1} align="center">
                  <FiAlertCircle size={12} />
                  <Text fontSize="xs">This field is required</Text>
                </HStack>
              </FormHelperText>
            )}
            {disableRoleNameEdit && (
              <FormHelperText color="gray.500" mt={1} fontSize="xs">
                Role name cannot be edited
              </FormHelperText>
            )}
          </FormControl>
          
          {/* System Access Section */}
          <Box mb={6}>
            <Text fontSize="md" fontWeight="semibold" mb={3}>
              System Access
              <Text as="span" fontSize="sm" fontWeight="normal" color="gray.500" ml={2}>
                Configure module permissions
              </Text>
            </Text>
            
            {/* Permissions Table */}
            <Box
              border="1px"
              borderColor={borderColor}
              borderRadius="md"
              overflow="hidden"
              shadow="sm"
              fontSize="sm"
            >
              <Table variant="simple" size="sm">
                <Thead bg={headerBgColor}>
                  <Tr>
                    <Th width="40px" px={2} py={2}></Th>
                    <Th fontSize="xs" py={2}>Module</Th>
                    <Th width="80px" textAlign="center" fontSize="xs" py={2}>Read</Th>
                    <Th width="80px" textAlign="center" fontSize="xs" py={2}>Create</Th>
                    <Th width="80px" textAlign="center" fontSize="xs" py={2}>Update</Th>
                    <Th width="80px" textAlign="center" fontSize="xs" py={2}>Delete</Th>
                    <Th width="90px" textAlign="center" fontSize="xs" py={2}>Approve</Th>
                  </Tr>
                </Thead>
                <Tbody>
                  {modules.map((module, index) => (
                    <Tr 
                      key={module.id} 
                      _hover={{ bg: hoverBgColor }}
                      bg={index % 2 === 0 ? bgColor : useColorModeValue('gray.50', 'gray.800')}
                      borderBottom="1px"
                      borderColor={tableBorderColor}
                      h="50px"
                    >
                      <Td textAlign="center" px={2} py={2}>
                        <Tooltip label="Delete module" hasArrow placement="left">
                          <IconButton
                            icon={<FiTrash2 size={14} />}
                            size="sm"
                            variant="ghost"
                            colorScheme="red"
                            onClick={() => handleDeleteModule(module.id)}
                            aria-label="Delete module"
                            _hover={{ bg: 'red.50', color: 'red.500' }}
                          />
                        </Tooltip>
                      </Td>
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
                  ))}
                  
                  {modules.length === 0 && (
                    <Tr>
                      <Td colSpan={7} textAlign="center" py={4}>
                        <Text color="gray.500" fontSize="sm">No modules found. Add a module to configure permissions.</Text>
                      </Td>
                    </Tr>
                  )}
                </Tbody>
              </Table>
            </Box>
            
            {/* Add Module section has been removed */}
          </Box>
        </ModalBody>
        
        <ModalFooter
          bg={headerBgColor}
          borderTop="1px"
          borderColor={borderColor}
          p={4}
          borderBottomRadius="lg"
          gap={2}
          justifyContent="flex-end"
          position="sticky" 
          bottom="0"
          zIndex="1"
        >
          <Button 
            variant="outline" 
            onClick={onClose}
            size="md"
            fontWeight="medium"
            borderRadius="md"
            px={6}
            fontSize="sm"
          >
            Cancel
          </Button>
          <Button 
            colorScheme="blue" 
            onClick={handleSave}
            size="md"
            fontWeight="medium"
            borderRadius="md"
            px={6}
            isDisabled={!formIsValid}
            _hover={{
              transform: 'translateY(-1px)',
              shadow: 'sm'
            }}
            _active={{
              transform: 'translateY(0)',
              shadow: 'xs'
            }}
            transition="all 0.2s"
            fontSize="sm"
          >
            {roleId ? 'Update' : 'Create'} Role
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default RoleConfigurationComponent;