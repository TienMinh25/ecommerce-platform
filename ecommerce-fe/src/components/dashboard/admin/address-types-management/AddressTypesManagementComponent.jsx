import React, {useEffect, useState} from 'react';
import {
    Box,
    Button,
    Container,
    Flex,
    FormControl,
    FormErrorMessage,
    FormLabel,
    HStack,
    IconButton,
    Input,
    InputGroup,
    InputLeftElement,
    Menu,
    MenuButton,
    MenuItem,
    MenuList,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Table,
    Tbody,
    Td,
    Text,
    Th,
    Thead,
    Tooltip,
    Tr,
    useColorModeValue,
    useToast,
} from '@chakra-ui/react';
import {
    FiAlertCircle,
    FiChevronDown,
    FiChevronLeft,
    FiChevronRight,
    FiEdit2,
    FiMapPin,
    FiPlus,
    FiRefreshCw,
    FiSearch,
    FiTrash2,
} from 'react-icons/fi';
import addressTypeService from "../../../../services/addressTypeService.js";
import {formatDateWithTime} from "../../../../utils/time.js";

const AddressTypesManagementComponent = () => {
    // Hook initialization - ensure consistent order
    const toast = useToast();

    // Colors - moved up to maintain Hook order
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const scrollTrackBg = useColorModeValue('#f1f1f1', '#2d3748');
    const scrollThumbBg = useColorModeValue('#c1c1c1', '#4a5568');
    const scrollThumbHoverBg = useColorModeValue('#a1a1a1', '#718096');
    const theadBg = useColorModeValue('gray.50', 'gray.900');
    const typeNameColor = useColorModeValue('gray.600', 'gray.300');
    const rowHoverBg = useColorModeValue('blue.50', 'gray.700');
    const rowEvenBg = useColorModeValue('gray.50', 'gray.800');
    const rowActiveBg = useColorModeValue('blue.100', 'gray.600');
    const textColor = useColorModeValue('gray.800', 'white');
    const updatedTextColor = useColorModeValue('gray.600', 'gray.300');
    const paginationBg = useColorModeValue('gray.50', 'gray.800');
    const paginationGradient = useColorModeValue(
        "linear(to-r, white, gray.50, white)",
        "linear(to-r, gray.800, gray.700, gray.800)"
    );

    // State variables
    const [currentPage, setCurrentPage] = useState(1);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [totalItems, setTotalItems] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [searchQuery, setSearchQuery] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [addressTypes, setAddressTypes] = useState([]);

    // Modal states
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
    const [currentAddressType, setCurrentAddressType] = useState(null);
    const [addressTypeName, setAddressTypeName] = useState('');
    const [formError, setFormError] = useState('');

    // Load address types on page load and when pagination changes
    useEffect(() => {
        fetchAddressTypes();
    }, [currentPage, rowsPerPage]);

    // Fetch address types from API
    const fetchAddressTypes = async () => {
        setIsLoading(true);
        try {
            const response = await addressTypeService.getAddressTypes(currentPage, rowsPerPage);

            if (response && response.data) {
                const formattedData = response.data.map(item => ({
                    id: item.id,
                    name: item.address_type,
                    createdAt: new Date(item.created_at).toLocaleDateString(),
                    updatedAt: formatDateWithTime(item.updated_at),
                }));

                setAddressTypes(formattedData);

                // Set pagination data from API response
                if (response.metadata && response.metadata.pagination) {
                    setTotalItems(response.metadata.pagination.total_items || formattedData.length);
                    setTotalPages(response.metadata.pagination.total_pages || 1);
                } else {
                    setTotalItems(formattedData.length);
                    setTotalPages(Math.ceil(formattedData.length / rowsPerPage));
                }
            }
        } catch (error) {
            console.error('Error fetching address types:', error);
            toast({
                title: 'Error',
                description: 'Failed to load address types',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    };

    // Handle create address type
    const handleCreateAddressType = async () => {
        if (!addressTypeName.trim()) {
            setFormError('Address type name is required');
            return;
        }

        setIsLoading(true);
        try {
            await addressTypeService.createAddressType({ address_type: addressTypeName });
            toast({
                title: 'Success',
                description: 'Address type created successfully',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
            setIsCreateModalOpen(false);
            setAddressTypeName('');
            fetchAddressTypes();
        } catch (error) {
            console.error('Error creating address type:', error);
            toast({
                title: 'Error',
                description: error.response?.data?.error?.message || 'Failed to create address type',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    };

    // Handle update address type
    const handleUpdateAddressType = async () => {
        if (!addressTypeName.trim()) {
            setFormError('Address type name is required');
            return;
        }

        setIsLoading(true);
        try {
            await addressTypeService.updateAddressType(currentAddressType.id, { address_type: addressTypeName });
            toast({
                title: 'Success',
                description: 'Address type updated successfully',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
            setIsEditModalOpen(false);
            setAddressTypeName('');
            fetchAddressTypes();
        } catch (error) {
            console.error('Error updating address type:', error);
            toast({
                title: 'Error',
                description: error.response?.data?.error?.message || 'Failed to update address type',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    };

    // Handle delete address type
    const handleDeleteAddressType = async () => {
        setIsLoading(true);
        try {
            await addressTypeService.deleteAddressType(currentAddressType.id);
            toast({
                title: 'Success',
                description: 'Address type deleted successfully',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
            setIsDeleteModalOpen(false);
            fetchAddressTypes();
        } catch (error) {
            console.error('Error deleting address type:', error);
            toast({
                title: 'Error',
                description: error.response?.data?.error?.message || 'Failed to delete address type',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    };

    // Open edit modal
    const openEditModal = (addressType) => {
        setCurrentAddressType(addressType);
        setAddressTypeName(addressType.name);
        setFormError('');
        setIsEditModalOpen(true);
    };

    // Open delete modal
    const openDeleteModal = (addressType) => {
        setCurrentAddressType(addressType);
        setIsDeleteModalOpen(true);
    };

    // Generate pagination range
    const generatePaginationRange = (current, total) => {
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

    // Filtered data based on search
    const filteredData = addressTypes.filter(item =>
        searchQuery ? item.name.toLowerCase().includes(searchQuery.toLowerCase()) : true
    );

    return (
        <Container maxW="container.xl" py={6}>
            {/* Search Bar */}
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
                                placeholder="Search by name..."
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
                                onClick={fetchAddressTypes}
                                aria-label="Refresh data"
                                variant="ghost"
                                colorScheme="blue"
                                size="sm"
                                mr={2}
                                isLoading={isLoading}
                            />
                        </Tooltip>
                    </Flex>
                </Flex>

                {/* Actions */}
                <HStack spacing={2}>
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
                        onClick={() => {
                            setAddressTypeName('');
                            setFormError('');
                            setIsCreateModalOpen(true);
                        }}
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
                {/* Data Table Container */}
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
                        <Thead bg={theadBg} position="sticky" top={0} zIndex={1}>
                            <Tr>
                                <Th
                                    py={4}
                                    borderTopLeftRadius="md"
                                    fontSize="xs"
                                    color={typeNameColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Type Name
                                </Th>
                                <Th
                                    py={4}
                                    fontSize="xs"
                                    color={typeNameColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                    display={{ base: "none", md: "table-cell" }}
                                >
                                    Last Updated
                                </Th>
                                <Th
                                    py={4}
                                    textAlign="right"
                                    borderTopRightRadius="md"
                                    fontSize="xs"
                                    color={typeNameColor}
                                    letterSpacing="0.5px"
                                    textTransform="uppercase"
                                    fontWeight="bold"
                                >
                                    Actions
                                </Th>
                            </Tr>
                        </Thead>
                        <Tbody>
                            {isLoading ? (
                                <Tr>
                                    <Td colSpan={3} textAlign="center">
                                        <Text py={4}>Loading...</Text>
                                    </Td>
                                </Tr>
                            ) : filteredData.length > 0 ? (
                                filteredData.map((row, index) => (
                                    <Tr
                                        key={row.id}
                                        _hover={{ bg: rowHoverBg }}
                                        bg={index % 2 === 0 ? bgColor : rowEvenBg}
                                        transition="background-color 0.2s"
                                        cursor="pointer"
                                        borderBottomWidth="1px"
                                        borderColor={borderColor}
                                        _active={{ bg: rowActiveBg }}
                                        h="60px"
                                    >
                                        <Td>
                                            <HStack spacing={3}>
                                                <Box p={1.5} bg="cyan.50" borderRadius="md">
                                                    <FiMapPin color="teal" size={16} />
                                                </Box>
                                                <Text
                                                    fontWeight="medium"
                                                    fontSize="sm"
                                                    color={textColor}
                                                >
                                                    {row.name}
                                                </Text>
                                            </HStack>
                                        </Td>
                                        <Td display={{ base: "none", md: "table-cell" }}>
                                            <Text
                                                fontSize="sm"
                                                color={updatedTextColor}
                                            >
                                                {row.updatedAt}
                                            </Text>
                                        </Td>
                                        <Td textAlign="right">
                                            <HStack spacing={1} justifyContent="flex-end">
                                                <Tooltip label="Edit address type" hasArrow>
                                                    <IconButton
                                                        icon={<FiEdit2 size={15} />}
                                                        size="sm"
                                                        variant="ghost"
                                                        colorScheme="blue"
                                                        aria-label="Edit address type"
                                                        borderRadius="md"
                                                        onClick={() => openEditModal(row)}
                                                    />
                                                </Tooltip>
                                                <Tooltip label="Delete address type" hasArrow>
                                                    <IconButton
                                                        icon={<FiTrash2 size={15} />}
                                                        size="sm"
                                                        variant="ghost"
                                                        colorScheme="red"
                                                        aria-label="Delete address type"
                                                        borderRadius="md"
                                                        onClick={() => openDeleteModal(row)}
                                                    />
                                                </Tooltip>
                                            </HStack>
                                        </Td>
                                    </Tr>
                                ))
                            ) : (
                                <Tr>
                                    <Td colSpan={3} textAlign="center" py={12}>
                                        <Flex direction="column" align="center" justify="center" py={8}>
                                            <Box color="gray.400" mb={3}>
                                                <FiSearch size={36} />
                                            </Box>
                                            <Text fontWeight="normal" color="gray.500" fontSize="md">No address types found</Text>
                                            <Text color="gray.400" fontSize="sm" mt={1}>Try a different search term</Text>
                                        </Flex>
                                    </Td>
                                </Tr>
                            )}
                        </Tbody>
                    </Table>
                </Box>

                {/* Pagination Section */}
                <Box
                    borderTop="1px"
                    borderColor={borderColor}
                    bg={paginationBg}
                    bgGradient={paginationGradient}
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
                                Showing {filteredData.length > 0 ? ((currentPage - 1) * rowsPerPage) + 1 : 0}-
                                {Math.min(currentPage * rowsPerPage, totalItems)} of {totalItems} types
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

                        {filteredData.length > 0 && totalPages > 1 && (
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

            {/* Create Modal */}
            <Modal isOpen={isCreateModalOpen} onClose={() => setIsCreateModalOpen(false)}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Create Address Type</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <FormControl isInvalid={!!formError}>
                            <FormLabel>Address Type Name</FormLabel>
                            <Input
                                value={addressTypeName}
                                onChange={(e) => {
                                    setAddressTypeName(e.target.value);
                                    setFormError('');
                                }}
                                placeholder="Enter address type name"
                            />
                            {formError && <FormErrorMessage>{formError}</FormErrorMessage>}
                        </FormControl>
                    </ModalBody>
                    <ModalFooter>
                        <Button mr={3} onClick={() => setIsCreateModalOpen(false)}>Cancel</Button>
                        <Button
                            colorScheme="blue"
                            onClick={handleCreateAddressType}
                            isLoading={isLoading}
                        >
                            Create
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>

            {/* Edit Modal */}
            <Modal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Edit Address Type</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <FormControl isInvalid={!!formError}>
                            <FormLabel>Address Type Name</FormLabel>
                            <Input
                                value={addressTypeName}
                                onChange={(e) => {
                                    setAddressTypeName(e.target.value);
                                    setFormError('');
                                }}
                                placeholder="Enter address type name"
                            />
                            {formError && <FormErrorMessage>{formError}</FormErrorMessage>}
                        </FormControl>
                    </ModalBody>
                    <ModalFooter>
                        <Button mr={3} onClick={() => setIsEditModalOpen(false)}>Cancel</Button>
                        <Button
                            colorScheme="blue"
                            onClick={handleUpdateAddressType}
                            isLoading={isLoading}
                        >
                            Update
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>

            {/* Delete Confirmation Modal */}
            <Modal isOpen={isDeleteModalOpen} onClose={() => setIsDeleteModalOpen(false)}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Confirm Delete</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <Flex align="center" mb={4}>
                            <Box color="red.500" mr={3}>
                                <FiAlertCircle size={24} />
                            </Box>
                            <Text>
                                Are you sure you want to delete the address type <b>{currentAddressType?.name}</b>? This action cannot be undone.
                            </Text>
                        </Flex>
                    </ModalBody>
                    <ModalFooter>
                        <Button mr={3} onClick={() => setIsDeleteModalOpen(false)}>Cancel</Button>
                        <Button
                            colorScheme="red"
                            onClick={handleDeleteAddressType}
                            isLoading={isLoading}
                        >
                            Delete
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </Container>
    );
};

export default AddressTypesManagementComponent;