import React, { useState, useEffect } from 'react';
import {
  Box,
  Flex,
  Input,
  InputGroup,
  InputLeftElement,
  Button,
  HStack,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Text,
  IconButton,
  useColorModeValue,
  Tooltip,
  Container,
  Collapse,
  Badge,
  Divider,
} from '@chakra-ui/react';
import {
  FiSearch,
  FiChevronDown,
  FiChevronLeft,
  FiChevronRight,
  FiPlus,
  FiSettings,
  FiRefreshCw,
  FiChevronUp,
  FiCheck,
  FiMinus,
  FiLock,
} from 'react-icons/fi';
import RoleConfigurationComponent from './RoleConfigurationComponent.jsx'; // Make sure path is correct

const RoleManagementComponent = () => {
  // State for pagination
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchQuery, setSearchQuery] = useState('');

  // State for configuration modal
  const [isConfigModalOpen, setIsConfigModalOpen] = useState(false);
  const [selectedRoleId, setSelectedRoleId] = useState(null);
  const [isEditMode, setIsEditMode] = useState(false);

  // State for expanded rows
  const [expandedRows, setExpandedRows] = useState({});

  // Mock data for the table with permissions
  const [data, setData] = useState([
    {
      id: 1,
      role: 'Admin',
      permissions: [
        { id: 1, name: 'Dashboard', read: true, create: false, update: false, delete: false, approve: false },
        { id: 2, name: 'User Management', read: true, create: true, update: true, delete: true, approve: true },
        { id: 3, name: 'Reports', read: true, create: true, update: true, delete: true, approve: true },
      ]
    },
    {
      id: 2,
      role: 'Manager',
      permissions: [
        { id: 1, name: 'Dashboard', read: true, create: false, update: false, delete: false, approve: false },
        { id: 2, name: 'User Management', read: true, create: true, update: true, delete: false, approve: true },
        { id: 3, name: 'Reports', read: true, create: true, update: false, delete: false, approve: false },
      ]
    },
    {
      id: 3,
      role: 'Developer',
      permissions: [
        { id: 1, name: 'Dashboard', read: true, create: false, update: false, delete: false, approve: false },
        { id: 2, name: 'User Management', read: true, create: false, update: false, delete: false, approve: false },
        { id: 3, name: 'Reports', read: true, create: false, update: false, delete: false, approve: false },
      ]
    },
    {
      id: 4,
      role: 'Designer',
      permissions: [
        { id: 1, name: 'Dashboard', read: true, create: false, update: false, delete: false, approve: false },
        { id: 2, name: 'User Management', read: false, create: false, update: false, delete: false, approve: false },
        { id: 3, name: 'Reports', read: true, create: false, update: false, delete: false, approve: false },
      ]
    },
    {
      id: 5,
      role: 'User',
      permissions: [
        { id: 1, name: 'Dashboard', read: true, create: false, update: false, delete: false, approve: false },
        { id: 2, name: 'User Management', read: false, create: false, update: false, delete: false, approve: false },
        { id: 3, name: 'Reports', read: false, create: false, update: false, delete: false, approve: false },
      ]
    },
    {
      id: 6,
      role: 'Guest',
      permissions: [
        { id: 1, name: 'Dashboard', read: true, create: false, update: false, delete: false, approve: false },
        { id: 2, name: 'User Management', read: false, create: false, update: false, delete: false, approve: false },
        { id: 3, name: 'Reports', read: false, create: false, update: false, delete: false, approve: false },
      ]
    },
    // Additional roles can be added here
  ]);

  // Function to toggle expanded row
  const toggleRow = (id) => {
    setExpandedRows(prev => ({
      ...prev,
      [id]: !prev[id]
    }));
  };

  // Function to refresh data
  const refreshData = () => {
    // In a real app, this would fetch new data from an API
    setData([...data].sort(() => Math.random() - 0.5));
  };

  // Filter data based on search query
  const filteredData = searchQuery
    ? data.filter(item => item.role.toLowerCase().includes(searchQuery.toLowerCase()))
    : data;

  // Colors
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const hoverBgColor = useColorModeValue('gray.50', 'gray.700');
  const boxShadow = useColorModeValue('sm', 'dark-lg');
  const scrollTrackBg = useColorModeValue('#f1f1f1', '#2d3748');
  const scrollThumbBg = useColorModeValue('#c1c1c1', '#4a5568');
  const scrollThumbHoverBg = useColorModeValue('#a1a1a1', '#718096');
  const expandedBgColor = useColorModeValue('gray.50', 'gray.800');
  const tableBorderColor = useColorModeValue('gray.100', 'gray.800');
  const activePermissionColor = useColorModeValue('green.500', 'green.400');
  const inactivePermissionColor = useColorModeValue('gray.400', 'gray.600');

  // Pagination logic - safe calculation to prevent division by zero
  const totalPages = Math.max(1, Math.ceil(filteredData.length / rowsPerPage));

  // Make sure currentPage is within valid range
  useEffect(() => {
    if (currentPage > totalPages) {
      setCurrentPage(Math.max(1, totalPages));
    }
  }, [currentPage, totalPages]);

  // Function to generate pagination range
  const generatePaginationRange = (current, total) => {
    // Ensure safe values
    current = Math.max(1, Math.min(current, total));

    if (total <= 5) {
      return Array.from({ length: total }, (_, i) => i + 1);
    }

    if (current <= 3) {
      return [1, 2, 3, 4, 5, '...', total];
    }

    if (current >= total - 2) {
      return [1, '...', total - 4, total - 3, total - 2, total - 1, total];
    }

    return [1, '...', current - 1, current, current + 1, '...', total];
  };

  // Get pagination range
  const paginationRange = generatePaginationRange(currentPage, totalPages);

  // Get current page data
  const indexOfLastItem = currentPage * rowsPerPage;
  const indexOfFirstItem = indexOfLastItem - rowsPerPage;
  const currentItems = filteredData.slice(indexOfFirstItem, indexOfLastItem);

  // Reset to first page when search changes
  useEffect(() => {
    setCurrentPage(1);
  }, [searchQuery, rowsPerPage]);

  // Function to handle opening configuration modal for existing role
  const handleEditPermissions = (roleId, e) => {
    // Prevent row expansion when clicking the edit button
    e && e.stopPropagation();
    setSelectedRoleId(roleId);
    setIsEditMode(true);
    setIsConfigModalOpen(true);
  };

  // Function to handle opening configuration modal for new role
  const handleCreateRole = () => {
    setSelectedRoleId(null);
    setIsEditMode(false);
    setIsConfigModalOpen(true);
  };

  // Function to close modal and reset state
  const handleCloseModal = () => {
    setIsConfigModalOpen(false);
    setTimeout(() => {
      setIsEditMode(false);
      setSelectedRoleId(null);
    }, 300);
  };

  // Function to handle saving role data
  const handleSaveRole = (roleData) => {
    console.log('Role data saved:', roleData);

    // Create permissions array from the roleData
    const permissions = roleData.permissions.map(p => ({
      id: p.moduleId,
      name: p.moduleName,
      read: p.read,
      create: p.create,
      update: p.update,
      delete: p.delete,
      approve: p.approve
    }));

    // If editing existing role, update it in our data
    if (selectedRoleId) {
      setData(data.map(item =>
        item.id === selectedRoleId
          ? { ...item, role: roleData.name, permissions }
          : item
      ));
    }
    // If creating new role, add it to our data
    else {
      const newId = Math.max(0, ...data.map(item => item.id)) + 1;
      setData([...data, { id: newId, role: roleData.name, permissions }]);
    }
  };

  // Permission icon component
  const PermissionIcon = ({ isActive }) => {
    return isActive ? (
      <Box
        color={activePermissionColor}
        borderRadius="full"
        bg={useColorModeValue('green.50', 'green.900')}
        p={1}
        display="inline-flex"
        alignItems="center"
        justifyContent="center"
        boxSize="24px"
      >
        <FiCheck />
      </Box>
    ) : (
      <Box
        color={inactivePermissionColor}
        borderRadius="full"
        bg={useColorModeValue('gray.50', 'gray.800')}
        p={1}
        display="inline-flex"
        alignItems="center"
        justifyContent="center"
        boxSize="24px"
      >
        <FiMinus />
      </Box>
    );
  };

  // Function to render module name without tooltip
  const renderModuleName = (moduleName) => {
    return <Text>{moduleName}</Text>;
  };

  return (
    <Container maxW="container.xl" py={6}>
      {/* Role Configuration Modal */}
      <RoleConfigurationComponent
        isOpen={isConfigModalOpen}
        onClose={handleCloseModal}
        roleId={selectedRoleId}
        onSave={handleSaveRole}
        modalSize="5xl"
        disableRoleNameEdit={isEditMode}
      />

      {/* Search and Create Button - Moved outside of the table's border */}
      <Flex
        justifyContent="space-between"
        alignItems="center"
        p={4}
        mb={4}
        flexDir={{ base: 'column', md: 'row' }}
        gap={{ base: 4, md: 0 }}
      >
        <Flex
          flex={{ md: 1 }}
          direction={{ base: "column", sm: "row" }}
          gap={3}
          align={{ base: "stretch", sm: "center" }}
        >
          {/* Search Input */}
          <Flex
            borderWidth="1px"
            borderRadius="lg"
            overflow="hidden"
            align="center"
            bg={bgColor}
            shadow="sm"
            flex="1"
            maxW={{ base: "full", lg: "450px" }}
          >
            <InputGroup size="md" variant="unstyled">
              <InputLeftElement pointerEvents="none" h="full" pl={3}>
                <FiSearch color="gray.400" />
              </InputLeftElement>
              <Input
                placeholder="Search roles..."
                pl={10}
                pr={2}
                py={2.5}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                _placeholder={{ color: "gray.400" }}
              />
            </InputGroup>

            {/* Refresh Button */}
            <Tooltip label="Refresh data" hasArrow>
              <IconButton
                icon={<FiRefreshCw size={16} />}
                onClick={refreshData}
                aria-label="Refresh data"
                variant="ghost"
                colorScheme="blue"
                size="sm"
                mr={2}
              />
            </Tooltip>
          </Flex>
        </Flex>

        {/* Create Button */}
        <Button
          leftIcon={<FiPlus />}
          colorScheme="blue"
          size="sm"
          borderRadius="md"
          fontWeight="normal"
          px={4}
          shadow="md"
          bgGradient="linear(to-r, blue.400, blue.500)"
          color="white"
          _hover={{
            bgGradient: "linear(to-r, blue.500, blue.600)",
            shadow: 'lg',
            transform: 'translateY(-1px)'
          }}
          _active={{
            bgGradient: "linear(to-r, blue.600, blue.700)",
            transform: 'translateY(0)',
            shadow: 'md'
          }}
          transition="all 0.2s"
          onClick={handleCreateRole}
        >
          Create
        </Button>
      </Flex>

      {/* Table Container */}
      <Box
        width="100%"
        borderRadius="xl"
        overflow="hidden"
        boxShadow="lg"
        bg={bgColor}
        display="flex"
        flexDirection="column"
        borderWidth="1px"
        borderColor={borderColor}
      >
        {/* Data Table Container with Fixed Height */}
        <Box
          overflow="auto"
          sx={{
            '&::-webkit-scrollbar': {
              width: '8px',
              height: '8px',
            },
            '&::-webkit-scrollbar-track': {
              background: scrollTrackBg,
              borderRadius: '4px',
            },
            '&::-webkit-scrollbar-thumb': {
              background: scrollThumbBg,
              borderRadius: '4px',
            },
            '&::-webkit-scrollbar-thumb:hover': {
              background: scrollThumbHoverBg,
            },
          }}
          flex="1"
          minH="300px"
          maxH={{ base: "60vh", lg: "calc(100vh - 250px)" }}
          borderBottomWidth="1px"
          borderColor={borderColor}
        >
          <Table variant="simple" size="md" colorScheme="gray" style={{ borderCollapse: 'separate', borderSpacing: '0' }}>
            <Tbody>
              {currentItems.length > 0 ? (
                currentItems.map((row, index) => (
                  <React.Fragment key={row.id}>
                    <Tr
                      onClick={() => toggleRow(row.id)}
                      _hover={{ bg: useColorModeValue('blue.50', 'gray.700') }}
                      bg={expandedRows[row.id] ? expandedBgColor : (index % 2 === 0 ? bgColor : useColorModeValue('gray.50', 'gray.800'))}
                      transition="background-color 0.2s"
                      cursor="pointer"
                      borderBottomWidth={expandedRows[row.id] ? 0 : "1px"}
                      borderColor={borderColor}
                      _active={{ bg: useColorModeValue('blue.100', 'gray.600') }}
                      h="60px"
                    >
                      <Td
                        width="80px"
                        fontWeight="normal"
                        fontSize="sm"
                        py={4}
                        color={useColorModeValue('gray.700', 'gray.300')}
                      >
                        {indexOfFirstItem + index + 1}
                      </Td>
                      <Td py={4}>
                        <Text
                          fontWeight="normal"
                          fontSize="sm"
                          color={useColorModeValue('gray.800', 'white')}
                        >
                          {row.role}
                        </Text>
                      </Td>
                      <Td textAlign="right">
                        <Tooltip label="Edit permissions" hasArrow>
                          <Button
                            leftIcon={<FiSettings size={14} />}
                            size="sm"
                            colorScheme="blue"
                            variant="outline"
                            onClick={(e) => handleEditPermissions(row.id, e)}
                            borderRadius="md"
                            fontWeight="normal"
                            px={3}
                            py={1}
                            _hover={{
                              bg: 'blue.50',
                              borderColor: 'blue.400'
                            }}
                          >
                            Edit Permissions
                          </Button>
                        </Tooltip>
                      </Td>
                    </Tr>

                    {/* Expanded Content Row */}
                    <Tr
                      display={expandedRows[row.id] ? "table-row" : "none"}
                      bg={expandedBgColor}
                    >
                      <Td colSpan={3} py={0} px={0}>
                        <Collapse in={expandedRows[row.id]} animateOpacity>
                          <Box
                            py={3}
                            px={4}
                            bg={useColorModeValue('blue.50', 'gray.700')}
                            borderBottomWidth="1px"
                            borderColor={borderColor}
                            width="100%"  // Ensuring full width
                          >
                            <Box
                              borderRadius="lg"
                              overflow="hidden"
                              shadow="md"
                              fontSize="sm"
                              bg={bgColor}
                              transition="all 0.3s"
                              transform="translateY(0)"
                              width="100%"  // Set to 100% width instead of maxW
                              mx="0"        // Remove auto margin
                            >
                              {/* Gradient header bar */}
                              <Box
                                h="4px"
                                bgGradient="linear(to-r, blue.400, purple.500)"
                              />

                              <Box px={3} py={3} bg={useColorModeValue('blue.50', 'blue.900')} borderBottom="1px" borderColor={borderColor}>
                                <Flex justify="space-between" align="center">
                                  <Text fontWeight="normal" fontSize="sm" color={useColorModeValue('blue.600', 'blue.200')}>
                                    Permission Details
                                  </Text>
                                </Flex>
                              </Box>

                              <Table variant="simple" size="sm" colorScheme="blue" style={{ borderCollapse: 'separate', borderSpacing: '0', width: '100%' }}>
                                <Thead bg={useColorModeValue('gray.100', 'gray.700')}>
                                  <Tr>
                                    <Th
                                      fontSize="xs"
                                      py={3}
                                      pl={4}
                                      borderTopLeftRadius="md"
                                      color={useColorModeValue('gray.600', 'gray.300')}
                                      letterSpacing="0.5px"
                                      textTransform="uppercase"
                                      fontWeight="normal"
                                      width="auto"  // Let this column take remaining space
                                    >
                                      Module
                                    </Th>
                                    <Th
                                      width="80px"
                                      textAlign="center"
                                      fontSize="xs"
                                      py={3}
                                      color={useColorModeValue('gray.600', 'gray.300')}
                                      letterSpacing="0.5px"
                                      textTransform="uppercase"
                                      fontWeight="bold"
                                    >
                                      Read
                                    </Th>
                                    <Th
                                      width="80px"
                                      textAlign="center"
                                      fontSize="xs"
                                      py={3}
                                      color={useColorModeValue('gray.600', 'gray.300')}
                                      letterSpacing="0.5px"
                                      textTransform="uppercase"
                                      fontWeight="bold"
                                    >
                                      Create
                                    </Th>
                                    <Th
                                      width="80px"
                                      textAlign="center"
                                      fontSize="xs"
                                      py={3}
                                      color={useColorModeValue('gray.600', 'gray.300')}
                                      letterSpacing="0.5px"
                                      textTransform="uppercase"
                                      fontWeight="bold"
                                    >
                                      Update
                                    </Th>
                                    <Th
                                      width="80px"
                                      textAlign="center"
                                      fontSize="xs"
                                      py={3}
                                      color={useColorModeValue('gray.600', 'gray.300')}
                                      letterSpacing="0.5px"
                                      textTransform="uppercase"
                                      fontWeight="bold"
                                    >
                                      Delete
                                    </Th>
                                    <Th
                                      width="90px"
                                      textAlign="center"
                                      fontSize="xs"
                                      py={3}
                                      borderTopRightRadius="md"
                                      color={useColorModeValue('gray.600', 'gray.300')}
                                      letterSpacing="0.5px"
                                      textTransform="uppercase"
                                      fontWeight="bold"
                                    >
                                      Approve
                                    </Th>
                                  </Tr>
                                </Thead>
                                <Tbody>
                                  {row.permissions.map((permission, idx) => (
                                    <Tr
                                      key={`${row.id}-${permission.id}`}
                                      _hover={{ bg: useColorModeValue('blue.50', 'blue.900') }}
                                      bg={idx % 2 === 0 ? bgColor : useColorModeValue('gray.50', 'gray.700')}
                                      transition="all 0.2s"
                                    >
                                      <Td
                                        fontWeight="medium"
                                        fontSize="sm"
                                        py={3}
                                        pl={4}
                                        borderLeftWidth="2px"
                                        borderLeftColor="transparent"
                                        _hover={{
                                          borderLeftColor: useColorModeValue('blue.400', 'blue.300')
                                        }}
                                        width="auto"  // Let this column take remaining space
                                      >
                                        {renderModuleName(permission.name)}
                                      </Td>
                                      <Td textAlign="center" py={3} width="80px">
                                        <Box transform="scale(1.1)" transition="all 0.2s">
                                          <PermissionIcon isActive={permission.read} />
                                        </Box>
                                      </Td>
                                      <Td textAlign="center" py={3} width="80px">
                                        <Box transform="scale(1.1)" transition="all 0.2s">
                                          <PermissionIcon isActive={permission.create} />
                                        </Box>
                                      </Td>
                                      <Td textAlign="center" py={3} width="80px">
                                        <Box transform="scale(1.1)" transition="all 0.2s">
                                          <PermissionIcon isActive={permission.update} />
                                        </Box>
                                      </Td>
                                      <Td textAlign="center" py={3} width="80px">
                                        <Box transform="scale(1.1)" transition="all 0.2s">
                                          <PermissionIcon isActive={permission.delete} />
                                        </Box>
                                      </Td>
                                      <Td textAlign="center" py={3} width="90px">
                                        <Box transform="scale(1.1)" transition="all 0.2s">
                                          <PermissionIcon isActive={permission.approve} />
                                        </Box>
                                      </Td>
                                    </Tr>
                                  ))}
                                </Tbody>
                              </Table>
                            </Box>
                          </Box>
                        </Collapse>
                      </Td>
                    </Tr>
                  </React.Fragment>
                ))
              ) : (
                <Tr>
                  <Td colSpan={3} textAlign="center" py={12}>
                    <Flex direction="column" align="center" justify="center" py={8}>
                      <Box color="gray.400" mb={3}>
                        <FiSearch size={36} />
                      </Box>
                      <Text fontWeight="normal" color="gray.500" fontSize="md">No roles found</Text>
                      <Text color="gray.400" fontSize="sm" mt={1}>Try a different search term</Text>
                    </Flex>
                  </Td>
                </Tr>
              )}
            </Tbody>
          </Table>
        </Box>

        {/* Fixed Pagination Section */}
        <Box
          borderTop="1px"
          borderColor={borderColor}
          bg={useColorModeValue('gray.50', 'gray.800')}
          bgGradient={useColorModeValue(
            "linear(to-r, white, gray.50, white)",
            "linear(to-r, gray.800, gray.700, gray.800)"
          )}
          position="sticky"
          bottom="0"
          width="100%"
          zIndex="1"
          boxShadow="0 -2px 6px rgba(0,0,0,0.05)"
        >
          <Flex
            justifyContent="space-between"
            alignItems="center"
            py={4}
            px={6}
            flexWrap={{ base: "wrap", md: "nowrap" }}
            gap={4}
          >
            <HStack spacing={1} flexShrink={0}>
              <Text fontSize="sm" color="gray.600" fontWeight="normal">
                Showing {filteredData.length > 0 ? indexOfFirstItem + 1 : 0}-{Math.min(indexOfLastItem, filteredData.length)} of {filteredData.length} roles
              </Text>
              <Menu>
                <MenuButton
                  as={Button}
                  size="xs"
                  variant="ghost"
                  rightIcon={<FiChevronDown />}
                  ml={2}
                  fontWeight="normal"
                  color="gray.600"
                >
                  {rowsPerPage} per page
                </MenuButton>
                <MenuList minW="120px" shadow="lg" borderRadius="md">
                  <MenuItem onClick={() => setRowsPerPage(10)}>10 per page</MenuItem>
                  <MenuItem onClick={() => setRowsPerPage(15)}>15 per page</MenuItem>
                  <MenuItem onClick={() => setRowsPerPage(20)}>20 per page</MenuItem>
                </MenuList>
              </Menu>
            </HStack>

            {filteredData.length > 0 && (
              <HStack spacing={1} justify="center" width={{ base: "100%", md: "auto" }}>
                <IconButton
                  icon={<FiChevronLeft />}
                  size="sm"
                  variant="ghost"
                  isDisabled={currentPage === 1}
                  onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                  aria-label="Previous page"
                  borderRadius="md"
                />

                {paginationRange.map((page, index) => (
                  page === '...' ? (
                    <Text key={`ellipsis-${index}`} mx={1} color="gray.500">...</Text>
                  ) : (
                    <Button
                      key={`page-${page}`}
                      size="sm"
                      variant={currentPage === page ? "solid" : "ghost"}
                      colorScheme={currentPage === page ? "blue" : "gray"}
                      onClick={() => typeof page === 'number' && setCurrentPage(page)}
                      borderRadius="md"
                      minW="32px"
                    >
                      {page}
                    </Button>
                  )
                ))}

                <IconButton
                  icon={<FiChevronRight />}
                  size="sm"
                  variant="ghost"
                  isDisabled={currentPage === totalPages}
                  onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                  aria-label="Next page"
                  borderRadius="md"
                />
              </HStack>
            )}
          </Flex>
        </Box>
      </Box>
    </Container>
  );
};

export default RoleManagementComponent;