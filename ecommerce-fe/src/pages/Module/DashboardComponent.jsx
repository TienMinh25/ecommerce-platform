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

const DashboardComponent = () => {
  // State for pagination
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [activeFilter, setActiveFilter] = useState(null);
  
  // Mock data for the table with more entries
  const [data, setData] = useState([
    { id: 1, name: 'Alexander Mitchell', email: 'alex.mitchell@example.com', role: 'Admin', status: 'Active', lastActive: '2 hours ago', image: 'https://i.pravatar.cc/150?img=1' },
    { id: 2, name: 'Sophia Rodriguez', email: 'sophia.r@example.com', role: 'Manager', status: 'Active', lastActive: '5 minutes ago', image: 'https://i.pravatar.cc/150?img=5' },
    { id: 3, name: 'James Wilson', email: 'jwilson@example.com', role: 'User', status: 'Inactive', lastActive: '3 days ago', image: 'https://i.pravatar.cc/150?img=3' },
    { id: 4, name: 'Emma Thompson', email: 'emmathompson@example.com', role: 'Admin', status: 'Active', lastActive: '1 hour ago', image: 'https://i.pravatar.cc/150?img=4' },
    { id: 5, name: 'Daniel Carter', email: 'daniel.carter@example.com', role: 'Manager', status: 'Active', lastActive: 'Just now', image: 'https://i.pravatar.cc/150?img=11' },
    { id: 6, name: 'Olivia Martinez', email: 'olivia.m@example.com', role: 'Developer', status: 'Active', lastActive: '30 minutes ago', image: 'https://i.pravatar.cc/150?img=9' },
    { id: 7, name: 'William Johnson', email: 'wjohnson@example.com', role: 'User', status: 'Inactive', lastActive: '1 week ago', image: 'https://i.pravatar.cc/150?img=12' },
    { id: 8, name: 'Ava Garcia', email: 'ava.garcia@example.com', role: 'Designer', status: 'Active', lastActive: '4 hours ago', image: 'https://i.pravatar.cc/150?img=10' },
    { id: 9, name: 'Ethan Brown', email: 'ethan.brown@example.com', role: 'Developer', status: 'Active', lastActive: '1 day ago', image: 'https://i.pravatar.cc/150?img=13' },
    { id: 10, name: 'Isabella Smith', email: 'isabella.s@example.com', role: 'User', status: 'Inactive', lastActive: '2 weeks ago', image: 'https://i.pravatar.cc/150?img=16' },
    { id: 11, name: 'Mason Davis', email: 'mason.davis@example.com', role: 'Manager', status: 'Active', lastActive: '10 minutes ago', image: 'https://i.pravatar.cc/150?img=15' },
    { id: 12, name: 'Charlotte White', email: 'charlotte.w@example.com', role: 'Designer', status: 'Active', lastActive: '3 hours ago', image: 'https://i.pravatar.cc/150?img=7' },
    { id: 13, name: 'Jacob Anderson', email: 'j.anderson@example.com', role: 'Developer', status: 'Inactive', lastActive: '5 days ago', image: 'https://i.pravatar.cc/150?img=18' },
    { id: 14, name: 'Amelia Taylor', email: 'amelia.t@example.com', role: 'Admin', status: 'Active', lastActive: '45 minutes ago', image: 'https://i.pravatar.cc/150?img=6' },
    { id: 15, name: 'Michael Thomas', email: 'michael.t@example.com', role: 'User', status: 'Active', lastActive: '1 hour ago', image: 'https://i.pravatar.cc/150?img=20' },
    { id: 16, name: 'Elizabeth Clark', email: 'elizabeth.c@example.com', role: 'Manager', status: 'Inactive', lastActive: '4 days ago', image: 'https://i.pravatar.cc/150?img=8' },
    { id: 17, name: 'Benjamin Lewis', email: 'ben.lewis@example.com', role: 'Developer', status: 'Active', lastActive: '20 minutes ago', image: 'https://i.pravatar.cc/150?img=17' },
    { id: 18, name: 'Mia Walker', email: 'mia.walker@example.com', role: 'Designer', status: 'Active', lastActive: '2 hours ago', image: 'https://i.pravatar.cc/150?img=14' },
    { id: 19, name: 'Alexander Hall', email: 'alex.hall@example.com', role: 'User', status: 'Inactive', lastActive: '1 month ago', image: 'https://i.pravatar.cc/150?img=19' },
    { id: 20, name: 'Abigail Young', email: 'abigail.y@example.com', role: 'Admin', status: 'Active', lastActive: '15 minutes ago', image: 'https://i.pravatar.cc/150?img=2' },
  ]);

  // Filter data based on activeFilter
  const filteredData = activeFilter ? data.filter(item => item.status === activeFilter) : data;
  
  // Colors
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const hoverBgColor = useColorModeValue('gray.50', 'gray.700');
  const activeButtonBg = useColorModeValue('blue.50', 'blue.900');
  const boxShadow = useColorModeValue('sm', 'dark-lg');
  
  // Pagination logic
  const totalPages = Math.ceil(filteredData.length / rowsPerPage);
  const paginationRange = generatePaginationRange(currentPage, totalPages);
  
  // Get current page data
  const indexOfLastItem = currentPage * rowsPerPage;
  const indexOfFirstItem = indexOfLastItem - rowsPerPage;
  const currentItems = filteredData.slice(indexOfFirstItem, indexOfLastItem);
  
  function generatePaginationRange(current, total) {
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
  }

  // Reset to first page when filter changes
  useEffect(() => {
    setCurrentPage(1);
  }, [activeFilter]);
  
  return (
    <Container maxW="container.xl" py={6}>
      <Box 
        width="100%" 
        borderRadius="xl" 
        overflow="hidden" 
        boxShadow={boxShadow}
        bg={bgColor}
        display="flex"
        flexDirection="column"
      >
        {/* Search and Filter Bar */}
        <Flex 
          justifyContent="space-between" 
          alignItems="center"
          p={4}
          borderBottom="1px"
          borderColor={borderColor}
          flexDir={{ base: 'column', md: 'row' }}
          gap={{ base: 4, md: 0 }}
          bg={useColorModeValue('gray.50', 'gray.900')}
        >
          <Flex 
            flex={{ md: 1 }} 
            direction={{ base: "column", sm: "row" }}
            gap={3}
            align={{ base: "stretch", sm: "center" }}
          >
            {/* Enhanced Search Input with Integrated Status Filter */}
            <Flex
              borderWidth="1px"
              borderRadius="lg"
              overflow="hidden"
              align="center"
              bg={bgColor}
              shadow="sm"
              flex="1"
              maxW={{ base: "full", lg: "500px" }}
              position="relative"
            >
              <InputGroup size="md" variant="unstyled">
                <InputLeftElement pointerEvents="none" h="full" pl={3}>
                  <FiSearch color="gray.400" />
                </InputLeftElement>
                <Input 
                  placeholder="Search users..." 
                  pl={10}
                  pr={2}
                  py={2.5}
                  _placeholder={{ color: "gray.400" }}
                />
              </InputGroup>
              
              <Divider orientation="vertical" h="24px" mx={2} />
              
              {/* Status Toggle Buttons */}
              <HStack spacing={1} mr={2}>
                <Tooltip label={activeFilter === 'Active' ? "Clear Active filter" : "Show only Active users"} hasArrow>
                  <Button
                    size="sm"
                    variant={activeFilter === 'Active' ? "solid" : "ghost"}
                    colorScheme={activeFilter === 'Active' ? "green" : "gray"}
                    onClick={() => setActiveFilter(activeFilter === 'Active' ? null : 'Active')}
                    px={3}
                    leftIcon={<FiCheck size={14} />}
                    fontWeight={activeFilter === 'Active' ? "semibold" : "medium"}
                  >
                    Active
                  </Button>
                </Tooltip>
                
                <Tooltip label={activeFilter === 'Inactive' ? "Clear Inactive filter" : "Show only Inactive users"} hasArrow>
                  <Button
                    size="sm"
                    variant={activeFilter === 'Inactive' ? "solid" : "ghost"}
                    colorScheme={activeFilter === 'Inactive' ? "red" : "gray"}
                    onClick={() => setActiveFilter(activeFilter === 'Inactive' ? null : 'Inactive')}
                    px={3}
                    leftIcon={<FiX size={14} />}
                    fontWeight={activeFilter === 'Inactive' ? "semibold" : "medium"}
                  >
                    Inactive
                  </Button>
                </Tooltip>
                
                {activeFilter && (
                  <Tooltip label="Clear all filters" hasArrow>
                    <IconButton
                      icon={<FiRefreshCw size={14} />}
                      onClick={() => setActiveFilter(null)}
                      aria-label="Clear filter"
                      size="sm"
                      variant="ghost"
                      colorScheme="gray"
                    />
                  </Tooltip>
                )}
              </HStack>
            </Flex>
          </Flex>
          
          {/* Filters and Actions */}
          <HStack spacing={3}>
            {/* Filter Menu */}
            <Menu closeOnSelect={false}>
              <MenuButton 
                as={Button} 
                rightIcon={<FiChevronDown />} 
                leftIcon={<FiFilter />}
                variant="outline"
                size="md"
                borderRadius="md"
              >
                Filter
              </MenuButton>
              <MenuList minWidth="240px" p={2} shadow="lg">
                <VStack align="stretch" spacing={2}>
                  <Text fontWeight="medium" fontSize="sm" px={3} py={1} color="gray.500">Status</Text>
                  <MenuItem closeOnSelect={false}>Status: Active</MenuItem>
                  <MenuItem closeOnSelect={false}>Status: Inactive</MenuItem>
                  <Divider my={1} />
                  <Text fontWeight="medium" fontSize="sm" px={3} py={1} color="gray.500">Role</Text>
                  <MenuItem closeOnSelect={false}>Role: Admin</MenuItem>
                  <MenuItem closeOnSelect={false}>Role: Manager</MenuItem>
                  <MenuItem closeOnSelect={false}>Role: User</MenuItem>
                  <MenuItem closeOnSelect={false}>Role: Developer</MenuItem>
                </VStack>
              </MenuList>
            </Menu>
            
            {/* Dropdown Menu */}
            <Menu>
              <MenuButton 
                as={Button} 
                rightIcon={<FiChevronDown />}
                variant="outline"
                size="md"
                borderRadius="md"
              >
                Options
              </MenuButton>
              <MenuList shadow="lg">
                <MenuItem icon={<FiUser size={14} />}>User Settings</MenuItem>
                <MenuItem icon={<FiMail size={14} />}>Send Email</MenuItem>
                <MenuItem icon={<FiEye size={14} />}>View Details</MenuItem>
                <Divider />
                <MenuItem icon={<FiRefreshCw size={14} />}>Refresh Data</MenuItem>
              </MenuList>
            </Menu>
            
            {/* Create Button */}
            <Button 
              leftIcon={<FiPlus />} 
              colorScheme="blue"
              size="md"
              borderRadius="md"
              fontWeight="semibold"
              px={5}
              shadow="md"
              _hover={{ 
                bg: 'blue.500',
                shadow: 'lg',
                transform: 'translateY(-1px)'
              }}
              _active={{ 
                bg: 'blue.600',
                transform: 'translateY(0)',
                shadow: 'md' 
              }}
              transition="all 0.2s"
            >
              Create
            </Button>
          </HStack>
        </Flex>
        
        {/* Data Table Container with Fixed Height */}
        <Box 
          overflow="auto"
          css={{
            '&::-webkit-scrollbar': {
              width: '8px',
              height: '8px',
            },
            '&::-webkit-scrollbar-track': {
              background: useColorModeValue('#f1f1f1', '#2d3748'),
            },
            '&::-webkit-scrollbar-thumb': {
              background: useColorModeValue('#c1c1c1', '#4a5568'),
              borderRadius: '4px',
            },
            '&::-webkit-scrollbar-thumb:hover': {
              background: useColorModeValue('#a1a1a1', '#718096'),
            },
          }}
          flex="1"
          minH="300px"
          maxH={{ base: "60vh", lg: "calc(100vh - 250px)" }}
        >
          <Table variant="simple" size="md">
            <Thead bg={useColorModeValue('gray.50', 'gray.900')} position="sticky" top={0} zIndex={1}>
              <Tr>
                <Th py={4}>User</Th>
                <Th py={4}>Role</Th>
                <Th py={4}>Status</Th>
                <Th py={4}>Last Active</Th>
                <Th py={4} textAlign="right">Actions</Th>
              </Tr>
            </Thead>
            <Tbody>
              {currentItems.map((row) => (
                <Tr 
                  key={row.id}
                  _hover={{ bg: hoverBgColor }}
                  transition="background-color 0.2s"
                >
                  <Td>
                    <HStack spacing={3}>
                      <Avatar size="sm" name={row.name} src="/api/placeholder/40/40" />
                      <Box>
                        <Text fontWeight="medium">{row.name}</Text>
                        <Text fontSize="sm" color="gray.500">
                          <HStack spacing={1} mt={0.5}>
                            <FiMail size={12} />
                            <Text>{row.email}</Text>
                          </HStack>
                        </Text>
                      </Box>
                    </HStack>
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
                    >
                      <TagLeftIcon as={FiUser} boxSize="12px" />
                      <TagLabel>{row.role}</TagLabel>
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
                    >
                      {row.status}
                    </Badge>
                  </Td>
                  <Td>
                    <Text fontSize="sm" color="gray.500">
                      {row.lastActive}
                    </Text>
                  </Td>
                  <Td textAlign="right">
                    <HStack spacing={1} justifyContent="flex-end">
                      <Tooltip label="View details" hasArrow>
                        <IconButton
                          icon={<FiEye />} 
                          size="sm" 
                          variant="ghost" 
                          colorScheme="gray"
                          aria-label="View details"
                        />
                      </Tooltip>
                      <Tooltip label="Edit user" hasArrow>
                        <IconButton
                          icon={<FiEdit2 />} 
                          size="sm" 
                          variant="ghost" 
                          colorScheme="blue"
                          aria-label="Edit user"
                        />
                      </Tooltip>
                      <Tooltip label="Delete user" hasArrow>
                        <IconButton
                          icon={<FiTrash2 />} 
                          size="sm" 
                          variant="ghost" 
                          colorScheme="red"
                          aria-label="Delete user"
                        />
                      </Tooltip>
                    </HStack>
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </Box>
        
        {/* Fixed Pagination Section */}
        <Box 
          borderTop="1px"
          borderColor={borderColor}
          bg={useColorModeValue('gray.50', 'gray.900')}
          position="sticky"
          bottom="0"
          width="100%"
          zIndex="1"
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
              <Text fontSize="sm" color="gray.600">
                Showing {indexOfFirstItem + 1}-{Math.min(indexOfLastItem, filteredData.length)} of {filteredData.length} users
              </Text>
              <Menu>
                <MenuButton 
                  as={Button} 
                  size="xs" 
                  variant="ghost" 
                  rightIcon={<FiChevronDown />}
                  ml={2}
                >
                  {rowsPerPage} per page
                </MenuButton>
                <MenuList minW="120px" shadow="lg">
                  <MenuItem onClick={() => setRowsPerPage(10)}>10 per page</MenuItem>
                  <MenuItem onClick={() => setRowsPerPage(20)}>20 per page</MenuItem>
                  <MenuItem onClick={() => setRowsPerPage(50)}>50 per page</MenuItem>
                </MenuList>
              </Menu>
            </HStack>
            
            <HStack spacing={1} justify="center" width={{ base: "100%", md: "auto" }}>
              <IconButton
                icon={<FiChevronLeft />}
                size="sm"
                variant="ghost"
                isDisabled={currentPage === 1}
                onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                aria-label="Previous page"
              />
              
              {paginationRange.map((page, index) => (
                page === '...' ? (
                  <Text key={`ellipsis-${index}`} mx={1}>...</Text>
                ) : (
                  <Button
                    key={`page-${page}`}
                    size="sm"
                    variant={currentPage === page ? "solid" : "ghost"}
                    colorScheme={currentPage === page ? "blue" : "gray"}
                    onClick={() => typeof page === 'number' && setCurrentPage(page)}
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
              />
            </HStack>
          </Flex>
        </Box>
      </Box>
    </Container>
  );
};

export default DashboardComponent;