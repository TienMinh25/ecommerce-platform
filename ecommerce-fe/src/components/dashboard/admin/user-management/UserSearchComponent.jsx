import React, {useState} from 'react';
import {
    Flex,
    HStack,
    IconButton,
    Input,
    InputGroup,
    InputLeftElement,
    Select,
    Tooltip,
    useColorModeValue
} from '@chakra-ui/react';
import {FiSearch, FiX} from 'react-icons/fi';

const UserSearchComponent = ({ onSearch, isLoading }) => {
    // State for search
    const [searchField, setSearchField] = useState('');
    const [searchQuery, setSearchQuery] = useState('');
    const [typingTimeout, setTypingTimeout] = useState(0);

    // Theme colors
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');

    // Handle search input change with debounce
    const handleSearchInputChange = (e) => {
        const inputText = e.target.value;
        setSearchQuery(inputText);

        // Only trigger search if a search field is selected
        if (searchField) {
            clearTimeout(typingTimeout);

            setTypingTimeout(
                setTimeout(() => {
                    // If input is empty, clear the search
                    if (!inputText.trim()) {
                        onSearch('', '');
                    } else if (searchField && inputText.trim()) {
                        onSearch(searchField, inputText);
                    }
                }, 300)
            );
        }
    };

    // Handle search field change WITHOUT triggering search
    const handleSearchFieldChange = (e) => {
        const newField = e.target.value;
        setSearchField(newField);

        // Clear search if field is cleared
        if (newField === '') {
            setSearchQuery('');
            onSearch('', ''); // Clear search results
        }
    };

    // Clear search with immediate effect
    const clearSearch = () => {
        setSearchQuery('');
        onSearch(searchField, '');
    };

    // Search button click handler
    const handleSearchButtonClick = () => {
        if (searchField && searchQuery) {
            onSearch(searchField, searchQuery);
        }
    };

    // Check if search is active
    const hasActiveSearch = searchQuery && searchField;

    return (
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
                onChange={handleSearchFieldChange}
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
                <option value="">Select field</option>
                <option value="fullname">Name</option>
                <option value="email">Email</option>
                <option value="phone">Phone</option>
            </Select>

            <InputGroup size="md" variant="unstyled">
                <InputLeftElement pointerEvents="none" h="full" pl={3}>
                    <FiSearch color="gray.400" />
                </InputLeftElement>
                <Input
                    placeholder={searchField ? `Search by ${searchField.toLowerCase()}...` : "Select a field first"}
                    pl={10}
                    pr={2}
                    py={2.5}
                    value={searchQuery}
                    onChange={handleSearchInputChange}
                    _placeholder={{ color: "gray.400" }}
                    isDisabled={!searchField}
                />
            </InputGroup>

            {/* Search actions */}
            <HStack>
                {hasActiveSearch && (
                    <Tooltip label="Clear search" hasArrow>
                        <IconButton
                            icon={<FiX size={16} />}
                            onClick={clearSearch}
                            aria-label="Clear search"
                            variant="ghost"
                            colorScheme="red"
                            size="sm"
                        />
                    </Tooltip>
                )}
                <Tooltip label="Search" hasArrow>
                    <IconButton
                        icon={<FiSearch size={16} />}
                        onClick={handleSearchButtonClick}
                        aria-label="Search"
                        variant="ghost"
                        colorScheme="blue"
                        size="sm"
                        mr={2}
                        isDisabled={!searchField || !searchQuery}
                    />
                </Tooltip>
            </HStack>
        </Flex>
    );
};

export default UserSearchComponent;