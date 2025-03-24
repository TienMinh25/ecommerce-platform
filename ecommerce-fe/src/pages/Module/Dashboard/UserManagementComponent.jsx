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
  Badge,
  Avatar,
  Tooltip,
  Divider,
  Tag,
  TagLabel,
  TagLeftIcon,
  VStack,
  Container,
  Select,
} from '@chakra-ui/react';
import { 
  FiSearch, 
  FiFilter, 
  FiChevronDown, 
  FiChevronLeft, 
  FiChevronRight,
  FiPlus,
  FiUser,
  FiEdit2,
  FiTrash2,
  FiMail,
  FiCheck,
  FiX,
  FiRefreshCw,
  FiEye,
} from 'react-icons/fi';

const UserManagementComponent = () => {
  // State for pagination
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [activeFilter, setActiveFilter] = useState(null);
  const [searchField, setSearchField] = useState('Name');
  const [searchQuery, setSearchQuery] = useState('');
  
  // Mock data for the table with more entries
  const [data, setData] = useState([
    { id: 1, name: 'Alexander Mitchell', email: 'alex.mitchell@example.com', role: 'Admin', status: 'Active', lastActive: '2 hours ago', image: 'https://i.pravatar.cc/150?img=1', phone: '(555) 123-4567', address: '123 Main St, New York, NY' },
    { id: 2, name: 'Sophia Rodriguez', email: 'sophia.r@example.com', role: 'Manager', status: 'Active', lastActive: '5 minutes ago', image: 'https://i.pravatar.cc/150?img=5', phone: '(555) 234-5678', address: '456 Park Ave, Los Angeles, CA' },
    { id: 3, name: 'James Wilson', email: 'jwilson@example.com', role: 'User', status: 'Inactive', lastActive: '3 days ago', image: 'https://i.pravatar.cc/150?img=3', phone: '(555) 345-6789', address: '789 Oak Rd, Chicago, IL' },
    { id: 4, name: 'Emma Thompson', email: 'emmathompson@example.com', role: 'Admin', status: 'Active', lastActive: '1 hour ago', image: 'https://i.pravatar.cc/150?img=4', phone: '(555) 456-7890', address: '101 Pine St, San Francisco, CA' },
    { id: 5, name: 'Daniel Carter', email: 'daniel.carter@example.com', role: 'Manager', status: 'Active', lastActive: 'Just now', image: 'https://i.pravatar.cc/150?img=11', phone: '(555) 567-8901', address: '202 Maple Dr, Boston, MA' },
    { id: 6, name: 'Olivia Martinez', email: 'olivia.m@example.com', role: 'Developer', status: 'Active', lastActive: '30 minutes ago', image: 'https://i.pravatar.cc/150?img=9', phone: '(555) 678-9012', address: '303 Cedar Ln, Seattle, WA' },
    { id: 7, name: 'William Johnson', email: 'wjohnson@example.com', role: 'User', status: 'Inactive', lastActive: '1 week ago', image: 'https://i.pravatar.cc/150?img=12', phone: '(555) 789-0123', address: '404 Birch Blvd, Miami, FL' },
    { id: 8, name: 'Ava Garcia', email: 'ava.garcia@example.com', role: 'Designer', status: 'Active', lastActive: '4 hours ago', image: 'https://i.pravatar.cc/150?img=10', phone: '(555) 890-1234', address: '505 Elm St, Austin, TX' },
    { id: 9, name: 'Ethan Brown', email: 'ethan.brown@example.com', role: 'Developer', status: 'Active', lastActive: '1 day ago', image: 'https://i.pravatar.cc/150?img=13', phone: '(555) 901-2345', address: '606 Walnut Ave, Portland, OR' },
    { id: 10, name: 'Isabella Smith', email: 'isabella.s@example.com', role: 'User', status: 'Inactive', lastActive: '2 weeks ago', image: 'https://i.pravatar.cc/150?img=16', phone: '(555) 012-3456', address: '707 Cherry Way, Denver, CO' },
    { id: 11, name: 'Mason Davis', email: 'mason.davis@example.com', role: 'Manager', status: 'Active', lastActive: '10 minutes ago', image: 'https://i.pravatar.cc/150?img=15', phone: '(555) 123-4567', address: '808 Pineapple Pkwy, Atlanta, GA' },
    { id: 12, name: 'Charlotte White', email: 'charlotte.w@example.com', role: 'Designer', status: 'Active', lastActive: '3 hours ago', image: 'https://i.pravatar.cc/150?img=7', phone: '(555) 234-5678', address: '909 Orange Ct, Philadelphia, PA' },
    { id: 13, name: 'Jacob Anderson', email: 'j.anderson@example.com', role: 'Developer', status: 'Inactive', lastActive: '5 days ago', image: 'https://i.pravatar.cc/150?img=18', phone: '(555) 345-6789', address: '1010 Lemon Rd, Phoenix, AZ' },
    { id: 14, name: 'Amelia Taylor', email: 'amelia.t@example.com', role: 'Admin', status: 'Active', lastActive: '45 minutes ago', image: 'https://i.pravatar.cc/150?img=6', phone: '(555) 456-7890', address: '1111 Lime St, San Diego, CA' },
    { id: 15, name: 'Michael Thomas', email: 'michael.t@example.com', role: 'User', status: 'Active', lastActive: '1 hour ago', image: 'https://i.pravatar.cc/150?img=20', phone: '(555) 567-8901', address: '1212 Grape Ave, Detroit, MI' },
    { id: 16, name: 'Elizabeth Clark', email: 'elizabeth.c@example.com', role: 'Manager', status: 'Inactive', lastActive: '4 days ago', image: 'https://i.pravatar.cc/150?img=8', phone: '(555) 678-9012', address: '1313 Peach Blvd, Nashville, TN' },
    { id: 17, name: 'Benjamin Lewis', email: 'ben.lewis@example.com', role: 'Developer', status: 'Active', lastActive: '20 minutes ago', image: 'https://i.pravatar.cc/150?img=17', phone: '(555) 789-0123', address: '1414 Plum Dr, Las Vegas, NV' },
    { id: 18, name: 'Mia Walker', email: 'mia.walker@example.com', role: 'Designer', status: 'Active', lastActive: '2 hours ago', image: 'https://i.pravatar.cc/150?img=14', phone: '(555) 890-1234', address: '1515 Apricot Ln, Kansas City, MO' },
    { id: 19, name: 'Alexander Hall', email: 'alex.hall@example.com', role: 'User', status: 'Inactive', lastActive: '1 month ago', image: 'https://i.pravatar.cc/150?img=19', phone: '(555) 901-2345', address: '1616 Mango St, Indianapolis, IN' },
    { id: 20, name: 'Abigail Young', email: 'abigail.y@example.com', role: 'Admin', status: 'Active', lastActive: '15 minutes ago', image: 'https://i.pravatar.cc/150?img=2', phone: '(555) 012-3456', address: '1717 Kiwi Way, Columbus, OH' },
  ]);

  // Function to refresh data
  const refreshData = () => {
    // In a real app, this would fetch new data from an API
    setData([...data].sort(() => Math.random() - 0.5));
  };

  // Filter data based on activeFilter and searchQuery
  const filteredData = data.filter(item => {
    const matchesFilter = activeFilter ? item.status === activeFilter : true;
    const matchesSearch = searchQuery 
      ? item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        item.email.toLowerCase().includes(searchQuery.toLowerCase())
      : true;
    return matchesFilter && matchesSearch;
  });
  
  // Colors - Matching the RoleManagementComponent
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const hoverBgColor = useColorModeValue('blue.50', 'gray.700');
  const boxShadow = useColorModeValue('lg', 'dark-lg');
  const scrollTrackBg = useColorModeValue('#f1f1f1', '#2d3748');
  const scrollThumbBg = useColorModeValue('#c1c1c1', '#4a5568');
  const scrollThumbHoverBg = useColorModeValue('#a1a1a1', '#718096');
  const tableBorderColor = useColorModeValue('gray.100', 'gray.800');
  
  // Pagination logic
  const totalPages = Math.max(1, Math.ceil(filteredData.length / rowsPerPage));
  
  // Make sure currentPage is within valid range
  useEffect(() => {
    if (currentPage > totalPages) {
      setCurrentPage(Math.max(1, totalPages));
    }
  }, [currentPage, totalPages]);
  
  // Reset to first page when filter changes
  useEffect(() => {
    setCurrentPage(1);
  }, [activeFilter, searchQuery, rowsPerPage]);
  
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
  
  return (
    <Container maxW="container.xl" py={6}>
      {/* Search and Filter Bar - Outside table box, similar to role management */}
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
          {/* Enhanced Search Input */}
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
            {/* Search Field Dropdown */}
            <Select
              value={searchField}
              onChange={(e) => setSearchField(e.target.value)}
              variant="unstyled"
              size="md"
              w="120px"
              pl={3}
              pr={0}
              py={2.5}
              borderRight="1px"
              borderColor={borderColor}
              borderRadius="0"
              _focus={{ boxShadow: "none" }}
              fontSize="sm"
            >
              <option value="Name">Name</option>
              <option value="Email">Email</option>
              <option value="Phone">Phone</option>
            </Select>

            <InputGroup size="md" variant="unstyled">
              <InputLeftElement pointerEvents="none" h="full" pl={3}>
                <FiSearch color="gray.400" />
              </InputLeftElement>
              <Input 
                placeholder={`Search by ${searchField.toLowerCase()}...`} 
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
        
        {/* Actions */}
        <HStack spacing={2}>
          {/* Filter Menu */}
          <Menu closeOnSelect={false}>
            <MenuButton 
              as={Button} 
              rightIcon={<FiChevronDown />} 
              leftIcon={<FiFilter />}
              variant="outline"
              size="sm"
              borderRadius="md"
              fontWeight="normal"
            >
              Filter
            </MenuButton>
            <MenuList minWidth="240px" p={2} shadow="lg" borderRadius="md">
              <VStack align="stretch" spacing={2}>
                <Text fontWeight="medium" fontSize="sm" px={3} py={1} color="gray.500">Status</Text>
                <MenuItem closeOnSelect={false} onClick={() => setActiveFilter(activeFilter === 'Active' ? null : 'Active')}>
                  <Flex align="center" justify="space-between" width="100%">
                    <Text>Status: Active</Text>
                    {activeFilter === 'Active' && <FiCheck />}
                  </Flex>
                </MenuItem>
                <MenuItem closeOnSelect={false} onClick={() => setActiveFilter(activeFilter === 'Inactive' ? null : 'Inactive')}>
                  <Flex align="center" justify="space-between" width="100%">
                    <Text>Status: Inactive</Text>
                    {activeFilter === 'Inactive' && <FiCheck />}
                  </Flex>
                </MenuItem>
              </VStack>
            </MenuList>
          </Menu>
          
          {/* Create Button - Matching RoleManagement style */}
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
          >
            Create
          </Button>
        </HStack>
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
            <Thead bg={useColorModeValue('gray.50', 'gray.900')} position="sticky" top={0} zIndex={1}>
              <Tr>
                <Th 
                  py={4} 
                  borderTopLeftRadius="md" 
                  fontSize="xs" 
                  color={useColorModeValue('gray.600', 'gray.300')}
                  letterSpacing="0.5px"
                  textTransform="uppercase"
                  fontWeight="bold"
                >
                  User
                </Th>
                <Th 
                  py={4} 
                  display={{ base: "none", md: "table-cell" }} 
                  fontSize="xs" 
                  color={useColorModeValue('gray.600', 'gray.300')}
                  letterSpacing="0.5px"
                  textTransform="uppercase"
                  fontWeight="bold"
                >
                  Contact Info
                </Th>
                <Th 
                  py={4} 
                  fontSize="xs" 
                  color={useColorModeValue('gray.600', 'gray.300')}
                  letterSpacing="0.5px"
                  textTransform="uppercase"
                  fontWeight="bold"
                >
                  Role
                </Th>
                <Th 
                  py={4} 
                  fontSize="xs" 
                  color={useColorModeValue('gray.600', 'gray.300')}
                  letterSpacing="0.5px"
                  textTransform="uppercase"
                  fontWeight="bold"
                >
                  Status
                </Th>
                <Th 
                  py={4} 
                  fontSize="xs" 
                  color={useColorModeValue('gray.600', 'gray.300')}
                  letterSpacing="0.5px"
                  textTransform="uppercase"
                  fontWeight="bold"
                >
                  Last Active
                </Th>
                <Th 
                  py={4} 
                  textAlign="right" 
                  borderTopRightRadius="md"
                  fontSize="xs" 
                  color={useColorModeValue('gray.600', 'gray.300')}
                  letterSpacing="0.5px"
                  textTransform="uppercase"
                  fontWeight="bold"
                >
                  Actions
                </Th>
              </Tr>
            </Thead>
            <Tbody>
              {currentItems.length > 0 ? (
                currentItems.map((row, index) => (
                  <Tr 
                    key={row.id}
                    _hover={{ bg: useColorModeValue('blue.50', 'gray.700') }}
                    bg={index % 2 === 0 ? bgColor : useColorModeValue('gray.50', 'gray.800')}
                    transition="background-color 0.2s"
                    cursor="pointer"
                    borderBottomWidth="1px"
                    borderColor={borderColor}
                    _active={{ bg: useColorModeValue('blue.100', 'gray.600') }}
                    h="60px"
                  >
                    <Td>
                      <HStack spacing={3}>
                        <Avatar size="sm" name={row.name} src="/api/placeholder/40/40" />
                        <Box>
                          <Text 
                            fontWeight="normal"
                            fontSize="sm"
                            color={useColorModeValue('gray.800', 'white')}
                          >
                            {row.name}
                          </Text>
                          <Text 
                            fontSize="xs" 
                            color={useColorModeValue('gray.500', 'gray.400')}
                            display={{ base: "none", lg: "block" }}
                          >
                            {row.email}
                          </Text>
                        </Box>
                      </HStack>
                    </Td>
                    <Td display={{ base: "none", md: "table-cell" }}>
                      <VStack align="start" spacing={0.5}>
                        {/* SỬA: Sử dụng HStack bên ngoài Text để tránh lỗi <div> trong <p> */}
                        <HStack spacing={1} align="center">
                          <FiMail size={12} />
                          <Text 
                            fontSize="sm" 
                            color={useColorModeValue('gray.600', 'gray.300')}
                          >
                            {row.email}
                          </Text>
                        </HStack>
                        <Text 
                          fontSize="sm" 
                          color={useColorModeValue('gray.500', 'gray.400')}
                        >
                          {row.phone}
                        </Text>
                      </VStack>
                    </Td>
                    <Td>
                      <Tag
                        size="md"
                        variant="subtle"
                        colorScheme={
                          row.role === 'Admin' ? 'purple' :
                          row.role === 'Manager' ? 'blue' :
                          row.role === 'Developer' ? 'green' :
                          row.role === 'Designer' ? 'orange' : 'gray'
                        }
                        borderRadius="md"
                      >
                        <TagLeftIcon as={FiUser} boxSize="12px" />
                        <TagLabel fontSize="xs" fontWeight="medium">{row.role}</TagLabel>
                      </Tag>
                    </Td>
                    <Td>
                      <Badge
                        px={2}
                        py={1}
                        borderRadius="full"
                        colorScheme={row.status === 'Active' ? 'green' : 'red'}
                        textTransform="capitalize"
                        fontWeight="medium"
                        fontSize="xs"
                      >
                        {row.status}
                      </Badge>
                    </Td>
                    <Td>
                      <Text 
                        fontSize="sm" 
                        color={useColorModeValue('gray.500', 'gray.400')}
                      >
                        {row.lastActive}
                      </Text>
                    </Td>
                    <Td textAlign="right">
                      <HStack spacing={1} justifyContent="flex-end">
                        <Tooltip label="View details" hasArrow>
                          <IconButton
                            icon={<FiEye size={15} />} 
                            size="sm" 
                            variant="ghost" 
                            colorScheme="gray"
                            aria-label="View details"
                            borderRadius="md"
                          />
                        </Tooltip>
                        <Tooltip label="Edit user" hasArrow>
                          <IconButton
                            icon={<FiEdit2 size={15} />} 
                            size="sm" 
                            variant="ghost" 
                            colorScheme="blue"
                            aria-label="Edit user"
                            borderRadius="md"
                          />
                        </Tooltip>
                        <Tooltip label="Delete user" hasArrow>
                          <IconButton
                            icon={<FiTrash2 size={15} />} 
                            size="sm" 
                            variant="ghost" 
                            colorScheme="red"
                            aria-label="Delete user"
                            borderRadius="md"
                          />
                        </Tooltip>
                      </HStack>
                    </Td>
                  </Tr>
                ))
              ) : (
                <Tr>
                  <Td colSpan={6} textAlign="center" py={12}>
                    <Flex direction="column" align="center" justify="center" py={8}>
                      <Box color="gray.400" mb={3}>
                        <FiSearch size={36} />
                      </Box>
                      <Text fontWeight="normal" color="gray.500" fontSize="md">No users found</Text>
                      <Text color="gray.400" fontSize="sm" mt={1}>Try a different search term or filter</Text>
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
                Showing {filteredData.length > 0 ? indexOfFirstItem + 1 : 0}-{Math.min(indexOfLastItem, filteredData.length)} of {filteredData.length} users
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

export default UserManagementComponent;