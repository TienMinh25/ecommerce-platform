import React, { useCallback, useEffect, useState } from 'react';
import {
  Alert,
  AlertIcon,
  Badge,
  Box,
  Button,
  Container,
  Flex,
  HStack,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Spinner,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tooltip,
  Tr,
  useColorModeValue,
  useDisclosure,
  useToast,
} from '@chakra-ui/react';
import {
  FiChevronDown,
  FiChevronLeft,
  FiChevronRight,
  FiEdit2,
  FiPlus,
  FiRefreshCw,
  FiSearch,
  FiSettings,
  FiShield,
  FiTrash2,
} from 'react-icons/fi';
import roleService from '../../../../services/roleService.js';
import moduleService from '../../../../services/moduleService.js';
import permissionService from '../../../../services/permissionService.js';
import RoleSearchFilter from './RoleSearchFilter.jsx';
import CreateRoleModal from './CreateRoleModal.jsx';
import EditRoleModal from './EditRoleModal.jsx';
import RolePermissionPanel from './RolePermissionPanel.jsx';

const RoleManagementComponent = () => {
  const [roles, setRoles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [expandedRoleId, setExpandedRoleId] = useState(null);
  const [selectedRole, setSelectedRole] = useState(null);

  const [modulesList, setModulesList] = useState([]);
  const [permissionsList, setPermissionsList] = useState([]);
  const [isLoadingModules, setIsLoadingModules] = useState(false);

  const [currentPage, setCurrentPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalPages, setTotalPages] = useState(1);

  const [filters, setFilters] = useState({
    sortBy: 'name',
    sortOrder: 'asc',
    searchBy: 'name',
    searchValue: '',
  });
  const [activeFilterParams, setActiveFilterParams] = useState({
    sortBy: 'name',
    sortOrder: 'asc',
  });

  const { isOpen: isCreateModalOpen, onOpen: onOpenCreateModal, onClose: onCloseCreateModal } = useDisclosure();
  const { isOpen: isEditModalOpen, onOpen: onOpenEditModal, onClose: onCloseEditModal } = useDisclosure();

  const [showInlinePermissionPanel, setShowInlinePermissionPanel] = useState(false);
  const [selectedPermissionRole, setSelectedPermissionRole] = useState(null);
  const [isLoadingPermissions, setIsLoadingPermissions] = useState(false);

  const toast = useToast();
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const hoverBgColor = useColorModeValue('blue.50', 'gray.700');
  const tableBorderColor = useColorModeValue('gray.100', 'gray.800');
  const headerBgColor = useColorModeValue('gray.50', 'gray.900');

  useEffect(() => {
    const loadBasicData = async () => {
      setIsLoadingModules(true);
      try {
        const [modulesResponse, permissionsResponse] = await Promise.all([
          moduleService.getModules({ getAll: true }),
          permissionService.getPermissions({ getAll: true }),
        ]);
        setModulesList(modulesResponse.data || []);
        setPermissionsList(permissionsResponse.data || []);
      } catch (error) {
        console.error('Error loading modules and permissions:', error);
        toast({
          title: 'Error loading data',
          description: 'Failed to load modules and permissions.',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      } finally {
        setIsLoadingModules(false);
      }
    };
    loadBasicData();
  }, [toast]);

  const createRequestParams = useCallback(() => {
    return {
      limit: rowsPerPage,
      page: currentPage,
      sortBy: activeFilterParams.sortBy,
      sortOrder: activeFilterParams.sortOrder,
      searchBy: filters.searchValue ? filters.searchBy : undefined,
      searchValue: filters.searchValue || undefined,
    };
  }, [currentPage, rowsPerPage, activeFilterParams, filters]);

  const mapRolePermissionsForDisplay = (role) => {
    if (!role || !role.permissions || !modulesList.length) return [];
    return modulesList.map((module) => {
      const modulePermissions = role.permissions?.find((p) => p.module_id === module.id);
      return {
        id: module.id,
        name: module.name,
        read: modulePermissions?.permissions?.includes(1) || false,
        create: modulePermissions?.permissions?.includes(2) || false,
        update: modulePermissions?.permissions?.includes(3) || false,
        delete: modulePermissions?.permissions?.includes(4) || false,
        approve: modulePermissions?.permissions?.includes(5) || false,
        reject: modulePermissions?.permissions?.includes(6) || false,
      };
    });
  };

  const PermissionIndicator = ({ isEnabled }) => (
      <Box
          w="16px"
          h="16px"
          borderRadius="sm"
          bg={isEnabled ? 'blue.500' : 'transparent'}
          borderWidth="1px"
          borderColor={isEnabled ? 'blue.500' : 'gray.300'}
          display="inline-flex"
          alignItems="center"
          justifyContent="center"
      >
        {isEnabled && <Box as="span" fontSize="xs" color="white" fontWeight="bold">✓</Box>}
      </Box>
  );

  const fetchRoles = useCallback(async () => {
    setIsLoading(true);
    setIsError(false);

    try {
      const params = createRequestParams();
      const response = await roleService.getRoles(params);

      // Check API response
      if (response && Array.isArray(response.data)) {
        setRoles(response.data);

        // Handle pagination metadata
        const metadata = response.metadata || {};
        const pagination = metadata.pagination || {};
        const totalItems = pagination.total_items || metadata.total_count || response.data.length;
        const totalPagesCalc = pagination.total_pages || Math.max(1, Math.ceil(totalItems / rowsPerPage));

        setTotalCount(totalItems);
        setTotalPages(totalPagesCalc);
      } else {
        throw new Error('Invalid response format from API');
      }
    } catch (error) {
      console.error('Error fetching roles:', error);
      setIsError(true);
      setErrorMessage(error.response?.data?.error?.message || 'Failed to fetch roles');
      toast({
        title: 'Error loading roles',
        description: error.response?.data?.error?.message || 'An error occurred while loading roles',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
      setRoles([]);
      setTotalCount(0);
      setTotalPages(1);
    } finally {
      setIsLoading(false);
    }
  }, [createRequestParams, rowsPerPage, toast]);

  useEffect(() => {
    if (modulesList.length > 0 && permissionsList.length > 0) {
      fetchRoles();
    }
  }, [modulesList, permissionsList, fetchRoles, currentPage, rowsPerPage, activeFilterParams]);

  const toggleRoleExpand = (roleId) => {
    setExpandedRoleId(expandedRoleId === roleId ? null : roleId);
  };

  const handleEditRole = (role, e) => {
    if (e) e.stopPropagation();
    setSelectedRole(role);
    onOpenEditModal();
  };

  const handleCreateRole = () => {
    onOpenCreateModal();
  };

  const handleDeleteRole = async (roleId, e) => {
    if (e) e.stopPropagation();
    if (window.confirm('Are you sure you want to delete this role?')) {
      try {
        await roleService.deleteRole(roleId);
        toast({
          title: 'Role deleted successfully',
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
        fetchRoles();
      } catch (error) {
        console.error('Error deleting role:', error);
        toast({
          title: 'Failed to delete role',
          description: error.response?.data?.error?.message || 'An unexpected error occurred',
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      }
    }
  };

  const handleOpenPermissionsPanel = (role, e) => {
    if (e) e.stopPropagation();
    setSelectedPermissionRole(role);
    setShowInlinePermissionPanel(true);
  };

  const handleClosePermissionsPanel = () => {
    setShowInlinePermissionPanel(false);
    setSelectedPermissionRole(null);
  };

  const handleRoleCreated = () => {
    fetchRoles();
  };

  const handleRoleUpdated = () => {
    fetchRoles();
  };

  const handleFiltersChange = (updatedFilters) => {
    setFilters(updatedFilters);
  };

  const handleApplyFilters = (filteredParams) => {
    setActiveFilterParams(filteredParams);
    setCurrentPage(1);
  };

  const handleReload = () => {
    fetchRoles();
  };

  const handleSavePermissions = async (roleId, permissions) => {
    setIsLoadingPermissions(true);
    try {
      // Format the data to match the API's expected structure if needed
      // Based on Swagger, needs to have role_name, modules_permissions, and optional description
      const selectedRole = roles.find(r => r.id === roleId);
      if (!selectedRole) {
        throw new Error('Role not found');
      }

      // Ensure proper format according to API schema
      const payload = {
        role_name: selectedRole.name,
        description: selectedRole.description || "",
        modules_permissions: permissions.modules_permissions || permissions.modules || []
      };

      await roleService.updateRolePermissions(roleId, payload);
      toast({
        title: 'Permissions updated successfully',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
      fetchRoles();
      handleClosePermissionsPanel();
    } catch (error) {
      console.error('Error updating permissions:', error);
      toast({
        title: 'Failed to update permissions',
        description: error.response?.data?.error?.message || 'An unexpected error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
      throw error;
    } finally {
      setIsLoadingPermissions(false);
    }
  };

  const generatePaginationRange = (current, total) => {
    current = Math.max(1, Math.min(current, total));
    if (total <= 5) return Array.from({ length: total }, (_, i) => i + 1);
    if (current <= 3) return [1, 2, 3, 4, 5, '...', total];
    if (current >= total - 2) return [1, '...', total - 4, total - 3, total - 2, total - 1, total];
    return [1, '...', current - 1, current, current + 1, '...', total];
  };

  const paginationRange = generatePaginationRange(currentPage, totalPages);

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    } catch (error) {
      return dateString;
    }
  };

  return (
      <Container maxW="container.xl" py={6}>
        <Flex justifyContent="space-between" alignItems="center" p={4} mb={4} flexDir={{ base: 'column', md: 'row' }} gap={{ base: 4, md: 0 }}>
          <Flex flex={{ md: 1 }} direction={{ base: "column", sm: "row" }} gap={3} align={{ base: "stretch", sm: "center" }}>
            <RoleSearchFilter filters={filters} onFiltersChange={handleFiltersChange} onApplyFilters={handleApplyFilters} />
            <Tooltip label="Refresh data" hasArrow>
              <IconButton
                  icon={<FiRefreshCw />}
                  onClick={handleReload}
                  aria-label="Refresh data"
                  variant="ghost"
                  colorScheme="blue"
                  size="sm"
                  isLoading={isLoading}
              />
            </Tooltip>
          </Flex>
          <HStack spacing={2}>
            <Button
                leftIcon={<FiPlus />}
                colorScheme="blue"
                size="sm"
                borderRadius="md"
                fontWeight="normal"
                px={4}
                shadow="md"
                onClick={handleCreateRole}
                bgGradient="linear(to-r, blue.400, blue.500)"
                color="white"
                _hover={{ bgGradient: 'linear(to-r, blue.500, blue.600)', shadow: 'lg', transform: 'translateY(-1px)' }}
            >
              Create Role
            </Button>
          </HStack>
        </Flex>

        {isError && (
            <Alert status="error" variant="left-accent" mb={4} borderRadius="md">
              <AlertIcon />
              <Text>{errorMessage || 'An error occurred while fetching roles'}</Text>
            </Alert>
        )}

        {showInlinePermissionPanel && selectedPermissionRole && (
            <RolePermissionPanel
                role={selectedPermissionRole}
                onSave={handleSavePermissions}
                onClose={handleClosePermissionsPanel}
                isLoading={isLoadingPermissions}
                modulesList={modulesList}
                permissionsList={permissionsList}
            />
        )}

        <Box width="100%" borderRadius="xl" overflow="hidden" boxShadow="lg" bg={bgColor} display="flex" flexDirection="column" borderWidth="1px" borderColor={borderColor}>
          <Box
              overflow="auto"
              sx={{
                '&::-webkit-scrollbar': { width: '8px', height: '8px' },
                '&::-webkit-scrollbar-track': { background: useColorModeValue('#f1f1f1', '#2d3748'), borderRadius: '4px' },
                '&::-webkit-scrollbar-thumb': { background: useColorModeValue('#c1c1c1', '#4a5568'), borderRadius: '4px' },
                '&::-webkit-scrollbar-thumb:hover': { background: useColorModeValue('#a1a1a1', '#718096') },
              }}
              flex="1"
              minH="300px"
              maxH={{ base: "60vh", lg: "calc(100vh - 250px)" }}
              borderBottomWidth="1px"
              borderColor={borderColor}
          >
            <Table variant="simple" size="md" colorScheme="gray">
              <Thead bg={headerBgColor} position="sticky" top={0} zIndex={1}>
                <Tr>
                  <Th py={4} fontWeight="bold" borderTopLeftRadius="md" fontSize="xs" color={useColorModeValue('gray.600', 'gray.300')} width="40%">
                    Role
                  </Th>
                  <Th py={4} fontWeight="bold" fontSize="xs" color={useColorModeValue('gray.600', 'gray.300')} width="40%">
                    Description
                  </Th>
                  <Th py={4} fontWeight="bold" fontSize="xs" color={useColorModeValue('gray.600', 'gray.300')} width="15%">
                    Updated
                  </Th>
                  <Th py={4} textAlign="right" borderTopRightRadius="md" fontSize="xs" color={useColorModeValue('gray.600', 'gray.300')} width="5%">
                    Actions
                  </Th>
                </Tr>
              </Thead>
              <Tbody>
                {isLoading ? (
                    <Tr>
                      <Td colSpan={4} textAlign="center" py={12}>
                        <Flex justify="center" align="center" direction="column">
                          <Box
                              h="40px"
                              w="40px"
                              borderWidth="3px"
                              borderStyle="solid"
                              borderColor="blue.500"
                              borderTopColor="transparent"
                              borderRadius="50%"
                              animation="spin 1s linear infinite"
                              mb={3}
                              sx={{ '@keyframes spin': { '0%': { transform: 'rotate(0deg)' }, '100%': { transform: 'rotate(360deg)' } } }}
                          />
                          <Text color="gray.500" fontSize="sm">Loading roles...</Text>
                        </Flex>
                      </Td>
                    </Tr>
                ) : roles.length > 0 ? (
                    roles.map((role) => (
                        <React.Fragment key={role.id}>
                          <Tr
                              onClick={() => toggleRoleExpand(role.id)}
                              _hover={{ bg: hoverBgColor }}
                              transition="background-color 0.2s"
                              borderBottomWidth="1px"
                              borderColor={borderColor}
                              cursor="pointer"
                              h="60px"
                              bg={expandedRoleId === role.id ? hoverBgColor : bgColor}
                          >
                            <Td fontWeight="medium">
                              <HStack>
                                <FiShield size={18} color={useColorModeValue('#3182CE', '#63B3ED')} />
                                <Text fontWeight="medium">{role.name}</Text>
                                {(role.name === 'admin' || role.name === 'Admin') && <Badge colorScheme="red" ml={2}>System</Badge>}
                              </HStack>
                            </Td>
                            <Td color={useColorModeValue('gray.600', 'gray.300')} fontSize="sm">{role.description || 'No description'}</Td>
                            <Td fontSize="sm" color={useColorModeValue('gray.500', 'gray.400')}>
                              {role.updated_at ? formatDate(role.updated_at) : 'N/A'}
                            </Td>
                            <Td textAlign="right">
                              <HStack spacing={1} justifyContent="flex-end">
                                <Tooltip label="Edit role" hasArrow>
                                  <IconButton
                                      icon={<FiEdit2 size={15} />}
                                      size="sm"
                                      variant="ghost"
                                      colorScheme="blue"
                                      aria-label="Edit role"
                                      onClick={(e) => handleEditRole(role, e)}
                                  />
                                </Tooltip>
                                {(role.name !== 'admin' && role.name !== 'Admin') && (
                                    <Tooltip label="Delete role" hasArrow>
                                      <IconButton
                                          icon={<FiTrash2 size={15} />}
                                          size="sm"
                                          variant="ghost"
                                          colorScheme="red"
                                          aria-label="Delete role"
                                          onClick={(e) => handleDeleteRole(role.id, e)}
                                      />
                                    </Tooltip>
                                )}
                              </HStack>
                            </Td>
                          </Tr>
                          {expandedRoleId === role.id && (
                              <Tr>
                                <Td colSpan={4} bg={useColorModeValue('gray.50', 'gray.700')} p={4}>
                                  <Box>
                                    {isLoadingModules ? (
                                        <Flex justify="center" align="center" py={4}>
                                          <Spinner size="md" color="blue.500" mr={3} />
                                          <Text color="gray.500">Loading permissions data...</Text>
                                        </Flex>
                                    ) : (
                                        <Box borderWidth="1px" borderRadius="md" borderColor={borderColor} overflow="auto">
                                          <Table variant="simple" size="sm">
                                            <Thead bg={headerBgColor} position="sticky" top={0} zIndex={1}>
                                              <Tr>
                                                <Th fontSize="xs" py={3} width="40%">
                                                  Module
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" py={3}>
                                                  Read
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" py={3}>
                                                  Create
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" py={3}>
                                                  Update
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" py={3}>
                                                  Delete
                                                </Th>
                                                <Th textAlign="center" fontSize="xs" py={3}>
                                                  Approve
                                                </Th>
                                              </Tr>
                                            </Thead>
                                            <Tbody>
                                              {modulesList.length > 0 ? (
                                                  mapRolePermissionsForDisplay(role).map((module, index) => (
                                                      <Tr
                                                          key={`permission-${role.id}-${module.id}`}
                                                          _hover={{ bg: hoverBgColor }}
                                                          bg={index % 2 === 0 ? bgColor : useColorModeValue('gray.50', 'gray.800')}
                                                      >
                                                        <Td py={2} fontWeight="medium" fontSize="sm">
                                                          {module.name}
                                                        </Td>
                                                        <Td textAlign="center" py={2}>
                                                          <PermissionIndicator isEnabled={module.read} />
                                                        </Td>
                                                        <Td textAlign="center" py={2}>
                                                          <PermissionIndicator isEnabled={module.create} />
                                                        </Td>
                                                        <Td textAlign="center" py={2}>
                                                          <PermissionIndicator isEnabled={module.update} />
                                                        </Td>
                                                        <Td textAlign="center" py={2}>
                                                          <PermissionIndicator isEnabled={module.delete} />
                                                        </Td>
                                                        <Td textAlign="center" py={2}>
                                                          <PermissionIndicator isEnabled={module.approve} />
                                                        </Td>
                                                      </Tr>
                                                  ))
                                              ) : (
                                                  <Tr>
                                                    <Td colSpan={6} textAlign="center" py={4}>
                                                      <Text color="gray.500">No permissions configured</Text>
                                                    </Td>
                                                  </Tr>
                                              )}
                                            </Tbody>
                                          </Table>
                                        </Box>
                                    )}
                                    <Flex justifyContent="flex-end" mt={3}>
                                      <Button
                                          size="sm"
                                          colorScheme="blue"
                                          leftIcon={<FiSettings />}
                                          onClick={(e) => {
                                            e.stopPropagation();
                                            handleOpenPermissionsPanel(role, e);
                                          }}
                                      >
                                        Manage Permissions
                                      </Button>
                                    </Flex>
                                  </Box>
                                </Td>
                              </Tr>
                          )}
                        </React.Fragment>
                    ))
                ) : (
                    <Tr>
                      <Td colSpan={4} textAlign="center" py={12}>
                        <Flex direction="column" align="center" justify="center" py={8}>
                          <Box color="gray.400" mb={3}>
                            <FiSearch size={36} />
                          </Box>
                          <Text fontWeight="normal" color="gray.500" fontSize="md">
                            No roles found
                          </Text>
                          <Text color="gray.400" fontSize="sm" mt={1}>
                            Try a different search term or filter
                          </Text>
                        </Flex>
                      </Td>
                    </Tr>
                )}
              </Tbody>
            </Table>
          </Box>

          <Box
              borderTop="1px"
              borderColor={borderColor}
              bg={useColorModeValue('gray.50', 'gray.800')}
              position="sticky"
              bottom="0"
              width="100%"
              zIndex="1"
              boxShadow="0 -2px 6px rgba(0,0,0,0.05)"
          >
            <Flex justifyContent="space-between" alignItems="center" py={4} px={6} flexWrap={{ base: "wrap", md: "nowrap" }} gap={4}>
              <HStack spacing={1} flexShrink={0}>
                <Text fontSize="sm" color="gray.600" fontWeight="normal">
                  Showing {totalCount > 0 ? (currentPage - 1) * rowsPerPage + 1 : 0}-{Math.min(currentPage * rowsPerPage, totalCount)} of {totalCount} roles
                </Text>
                <Menu>
                  <MenuButton as={Button} size="xs" variant="ghost" rightIcon={<FiChevronDown />} ml={2} fontWeight="normal" color="gray.600">
                    {rowsPerPage} per page
                  </MenuButton>
                  <MenuList minW="120px" shadow="lg" borderRadius="md">
                    <MenuItem onClick={() => setRowsPerPage(10)}>10 per page</MenuItem>
                    <MenuItem onClick={() => setRowsPerPage(15)}>15 per page</MenuItem>
                    <MenuItem onClick={() => setRowsPerPage(20)}>20 per page</MenuItem>
                    <MenuItem onClick={() => setRowsPerPage(50)}>50 per page</MenuItem>
                  </MenuList>
                </Menu>
              </HStack>

              {totalCount > 0 && (
                  <HStack spacing={1} justify="center" width={{ base: "100%", md: "auto" }}>
                    <IconButton
                        icon={<FiChevronLeft />}
                        size="sm"
                        variant="ghost"
                        isDisabled={currentPage === 1 || isLoading}
                        onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                        aria-label="Previous page"
                        borderRadius="md"
                    />
                    {paginationRange.map((page, index) =>
                        page === '...' ? (
                            <Text key={`ellipsis-${index}`} mx={1} color="gray.500">
                              ...
                            </Text>
                        ) : (
                            <Button
                                key={`page-${page}`}
                                size="sm"
                                variant={currentPage === page ? "solid" : "ghost"}
                                colorScheme={currentPage === page ? "blue" : "gray"}
                                onClick={() => typeof page === 'number' && setCurrentPage(page)}
                                borderRadius="md"
                                minW="32px"
                                isDisabled={isLoading}
                            >
                              {page}
                            </Button>
                        )
                    )}
                    <IconButton
                        icon={<FiChevronRight />}
                        size="sm"
                        variant="ghost"
                        isDisabled={currentPage === totalPages || isLoading}
                        onClick={() => setCurrentPage((prev) => Math.min(prev + 1, totalPages))}
                        aria-label="Next page"
                        borderRadius="md"
                    />
                  </HStack>
              )}
            </Flex>
          </Box>
        </Box>

        <CreateRoleModal isOpen={isCreateModalOpen} onClose={onCloseCreateModal} onRoleCreated={handleRoleCreated} modulesList={modulesList} permissionsList={permissionsList} />
        {selectedRole && (
            <EditRoleModal
                isOpen={isEditModalOpen}
                onClose={onCloseEditModal}
                role={selectedRole}
                onRoleUpdated={handleRoleUpdated}
                modulesList={modulesList}
                permissionsList={permissionsList}
            />
        )}
      </Container>
  );
};

export default RoleManagementComponent;