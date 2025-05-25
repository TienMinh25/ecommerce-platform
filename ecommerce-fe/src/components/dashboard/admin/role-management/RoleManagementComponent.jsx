import React, { useCallback, useEffect, useState, useMemo } from 'react';
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
  FiFilter,
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
import RoleSearchFilter from "./RoleSearchFilter.jsx";
import RolePermissionPanel from "./RolePermissionPanel.jsx";
import CreateRoleModal from "./CreateRoleModal.jsx";
import EditRoleModal from "./EditRoleModal.jsx";
import DeleteConfirmationModal from "../DeleteConfirmationComponent.jsx";
import RoleFilterDropdown from './RoleFilterDropdown.jsx';

const RoleManagementComponent = () => {
  // Toast for notifications
  const toast = useToast();

  // Modal disclosures (must be at the top, before any conditional logic)
  const { isOpen: isCreateModalOpen, onOpen: onOpenCreateModal, onClose: onCloseCreateModal } = useDisclosure();
  const { isOpen: isEditModalOpen, onOpen: onOpenEditModal, onClose: onCloseEditModal } = useDisclosure();
  const { isOpen: isDeleteModalOpen, onOpen: onOpenDeleteModal, onClose: onCloseDeleteModal } = useDisclosure();

  // UI colors (define all color values at the top)
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const hoverBgColor = useColorModeValue('blue.50', 'gray.700');
  const tableBorderColor = useColorModeValue('gray.100', 'gray.800');
  const headerBgColor = useColorModeValue('gray.50', 'gray.900');
  const textColor = useColorModeValue('gray.600', 'gray.300');
  const secondaryTextColor = useColorModeValue('gray.500', 'gray.400');
  const searchIconColor = useColorModeValue('#3182CE', '#63B3ED');

  // State variables
  const [roles, setRoles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [expandedRoleId, setExpandedRoleId] = useState(null);
  const [selectedRole, setSelectedRole] = useState(null);
  const [roleToDelete, setRoleToDelete] = useState(null);
  const [isDeletingRole, setIsDeletingRole] = useState(false);

  // Module and permissions state
  const [modulesList, setModulesList] = useState([]);
  const [permissionsList, setPermissionsList] = useState([]);
  const [isLoadingModules, setIsLoadingModules] = useState(false);

  // Pagination state
  const [currentPage, setCurrentPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalPages, setTotalPages] = useState(1);

  // Filter state
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

  // Permission panel state
  const [showInlinePermissionPanel, setShowInlinePermissionPanel] = useState(false);
  const [selectedPermissionRole, setSelectedPermissionRole] = useState(null);
  const [isLoadingPermissions, setIsLoadingPermissions] = useState(false);

  // Create request parameters for API calls (must be before any conditional use)
  const createRequestParams = useCallback(() => {
    return {
      limit: rowsPerPage,
      page: currentPage,
      sortBy: activeFilterParams.sortBy || undefined,
      sortOrder: activeFilterParams.sortOrder || undefined,
      searchBy: filters.searchValue ? filters.searchBy : undefined,
      searchValue: filters.searchValue || undefined,
    };
  }, [currentPage, rowsPerPage, activeFilterParams, filters]);

  // Permission indicator component for visual representation
  const PermissionIndicator = useMemo(() => {
    return ({ isEnabled }) => (
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
  }, []);

  // Generate pagination range for display (define as a useMemo hook)
  const generatePaginationRange = useCallback((current, total) => {
    current = Math.max(1, Math.min(current, total));
    if (total <= 5) return Array.from({ length: total }, (_, i) => i + 1);
    if (current <= 3) return [1, 2, 3, 4, 5, '...', total];
    if (current >= total - 2) return [1, '...', total - 4, total - 3, total - 2, total - 1, total];
    return [1, '...', current - 1, current, current + 1, '...', total];
  }, []);

  const paginationRange = useMemo(() =>
          generatePaginationRange(currentPage, totalPages),
      [currentPage, totalPages, generatePaginationRange]
  );

  // Fetch roles from API
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

  // Map role permissions for display in the UI
  const mapRolePermissionsForDisplay = useCallback((role) => {
    if (!role || !role.permissions || !modulesList.length || !permissionsList.length) return [];

    // Create permission mapping from permissionsList
    const permissionMapping = {};
    permissionsList.forEach(permission => {
      permissionMapping[permission.name] = permission.id;
    });

    return modulesList.map((module) => {
      const modulePermissions = role.permissions?.find((p) => p.module_id === module.id);
      return {
        id: module.id,
        name: module.name,
        // Use dynamic mapping instead of hardcoded IDs
        read: modulePermissions?.permissions?.includes(permissionMapping.read) || false,
        create: modulePermissions?.permissions?.includes(permissionMapping.create) || false,
        update: modulePermissions?.permissions?.includes(permissionMapping.update) || false,
        delete: modulePermissions?.permissions?.includes(permissionMapping.delete) || false,
        approve: modulePermissions?.permissions?.includes(permissionMapping.approve) || false,
        reject: modulePermissions?.permissions?.includes(permissionMapping.reject) || false,
      };
    });
  }, [modulesList, permissionsList]);

  // Load modules and permissions on mount
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

  // Fetch roles when dependencies change
  useEffect(() => {
    if (modulesList.length > 0 && permissionsList.length > 0) {
      fetchRoles();
    }
  }, [modulesList, permissionsList, fetchRoles, currentPage, rowsPerPage, activeFilterParams]);

  // Format date for display
  const formatDate = useCallback((dateString) => {
    if (!dateString) return 'N/A';
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    } catch (error) {
      return dateString;
    }
  }, []);

  // Handle expanding/collapsing role details
  const toggleRoleExpand = useCallback((roleId) => {
    setExpandedRoleId(expandedRoleId === roleId ? null : roleId);
  }, [expandedRoleId]);

  // Handle editing a role
  const handleEditRole = useCallback((role, e) => {
    if (e) e.stopPropagation();
    setSelectedRole(role);
    onOpenEditModal();
  }, [onOpenEditModal]);

  // Handle creating a new role
  const handleCreateRole = useCallback(() => {
    onOpenCreateModal();
  }, [onOpenCreateModal]);

  // Open delete confirmation modal
  const handleOpenDeleteModal = useCallback((role, e) => {
    if (e) e.stopPropagation();
    setRoleToDelete(role);
    onOpenDeleteModal();
  }, [onOpenDeleteModal]);

  // Confirm and process role deletion
  const handleConfirmDelete = useCallback(async () => {
    if (!roleToDelete) return;

    setIsDeletingRole(true);
    try {
      await roleService.deleteRole(roleToDelete.id);
      toast({
        title: 'Role deleted successfully',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
      setRoleToDelete(null);
      onCloseDeleteModal();
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
    } finally {
      setIsDeletingRole(false);
    }
  }, [roleToDelete, toast, onCloseDeleteModal, fetchRoles]);

  // Open permissions panel for a role
  const handleOpenPermissionsPanel = useCallback((role, e) => {
    if (e) e.stopPropagation();
    setSelectedPermissionRole(role);
    setShowInlinePermissionPanel(true);
  }, []);

  // Close permissions panel
  const handleClosePermissionsPanel = useCallback(() => {
    setShowInlinePermissionPanel(false);
    setSelectedPermissionRole(null);
  }, []);

  // Handle role creation callback
  const handleRoleCreated = useCallback(() => {
    fetchRoles();
  }, [fetchRoles]);

  // Handle role update callback
  const handleRoleUpdated = useCallback(() => {
    fetchRoles();
  }, [fetchRoles]);

  // Handle filter changes from the search component
  const handleFiltersChange = useCallback((updatedFilters) => {
    setFilters(updatedFilters);
  }, []);

  // Apply filters and reset to page 1
  const handleApplyFilters = useCallback((filteredParams) => {
    setActiveFilterParams(filteredParams);
    setCurrentPage(1);
  }, []);

  // Manually reload roles data
  const handleReload = useCallback(() => {
    fetchRoles();
  }, [fetchRoles]);

  // Save permissions for a role
  const handleSavePermissions = useCallback(async (roleId, permissions) => {
    setIsLoadingPermissions(true);
    try {
      // Get the role data
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
    } finally {
      setIsLoadingPermissions(false);
    }
  }, [roles, toast, fetchRoles, handleClosePermissionsPanel]);

  return (
      <Container maxW="container.xl" py={6}>
        {/* Header with search, filter, and action buttons */}
        <Flex justifyContent="space-between" alignItems="center" p={4} mb={4} flexDir={{ base: 'column', md: 'row' }} gap={{ base: 4, md: 0 }}>
          <RoleSearchFilter
              filters={filters}
              onFiltersChange={handleFiltersChange}
              onApplyFilters={handleApplyFilters}
              onRefresh={handleReload}
              isLoading={isLoading}
          />
          <HStack spacing={3}>
            <RoleFilterDropdown
                filters={filters}
                onFiltersChange={handleFiltersChange}
                onApplyFilters={handleApplyFilters}
            />
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
              Tạo mới vai trò
            </Button>
          </HStack>
        </Flex>

        {/* Error alert if needed */}
        {isError && (
            <Alert status="error" variant="left-accent" mb={4} borderRadius="md">
              <AlertIcon />
              <Text>{errorMessage || 'An error occurred while fetching roles'}</Text>
            </Alert>
        )}

        {/* Inline permission panel */}
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

        {/* Main table container */}
        <Box width="100%" borderRadius="xl" overflow="hidden" boxShadow="lg" bg={bgColor} display="flex" flexDirection="column" borderWidth="1px" borderColor={borderColor}>
          {/* Table with scrolling */}
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
                  <Th py={4} fontWeight="bold" borderTopLeftRadius="md" fontSize="xs" color={textColor} width="40%">
                    Tên vai trò
                  </Th>
                  <Th py={4} fontWeight="bold" fontSize="xs" color={textColor} width="40%">
                    Mô tả
                  </Th>
                  <Th py={4} fontWeight="bold" fontSize="xs" color={textColor} width="15%">
                    Ngày gần nhất cập nhật
                  </Th>
                  <Th py={4} textAlign="right" borderTopRightRadius="md" fontSize="xs" color={textColor} width="5%">
                    Hành động
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
                                <FiShield size={18} color={searchIconColor} />
                                <Text fontWeight="medium">{role.name}</Text>
                                {(role.name === 'admin' || role.name === 'Admin') && <Badge colorScheme="red" ml={2}>System</Badge>}
                              </HStack>
                            </Td>
                            <Td color={textColor} fontSize="sm">{role.description || 'No description'}</Td>
                            <Td fontSize="sm" color={secondaryTextColor}>
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
                                          onClick={(e) => handleOpenDeleteModal(role, e)}
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
                                        Quản lý quyền
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
                            Không tìm thấy vai trò nào
                          </Text>
                          <Text color="gray.400" fontSize="sm" mt={1}>
                            Vui lòng thử với từ khóa hoặc bộ lọc khác
                          </Text>
                        </Flex>
                      </Td>
                    </Tr>
                )}
              </Tbody>
            </Table>
          </Box>

          {/* Pagination footer */}
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
                  Hiển thị {totalCount > 0 ? (currentPage - 1) * rowsPerPage + 1 : 0}-{Math.min(currentPage * rowsPerPage, totalCount)} trên tổng số {totalCount} vai trò
                </Text>
                <Menu>
                  <MenuButton as={Button} size="xs" variant="ghost" rightIcon={<FiChevronDown />} ml={2} fontWeight="normal" color="gray.600">
                    {rowsPerPage} dòng mỗi trang
                  </MenuButton>
                  <MenuList minW="120px" shadow="lg" borderRadius="md">
                    <MenuItem onClick={() => setRowsPerPage(10)}>10 dòng mỗi trang</MenuItem>
                    <MenuItem onClick={() => setRowsPerPage(15)}>15 dòng mỗi trang</MenuItem>
                    <MenuItem onClick={() => setRowsPerPage(20)}>20 dòng mỗi trang</MenuItem>
                    <MenuItem onClick={() => setRowsPerPage(50)}>50 dòng mỗi trang</MenuItem>
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

        {/* Modals */}
        <CreateRoleModal
            isOpen={isCreateModalOpen}
            onClose={onCloseCreateModal}
            onRoleCreated={handleRoleCreated}
            modulesList={modulesList}
            permissionsList={permissionsList}
        />

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

        {/* Delete confirmation modal */}
        <DeleteConfirmationModal
            isOpen={isDeleteModalOpen}
            onClose={onCloseDeleteModal}
            onConfirm={handleConfirmDelete}
            title="Xoá vai trò"
            message="Bạn có chắc chắn muốn xóa vai trò này không? Hành động này không thể hoàn tác."
            itemName={roleToDelete?.name}
            isLoading={isDeletingRole}
        />
      </Container>
  );
};

export default RoleManagementComponent;