import React, { useState } from 'react';
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
import { FiSearch, FiX } from 'react-icons/fi';

const SupplierSearchComponent = ({ onSearch, isLoading }) => {
    const [searchField, setSearchField] = useState('');
    const [searchQuery, setSearchQuery] = useState('');
    const [typingTimeout, setTypingTimeout] = useState(0);

    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');

    const handleSearchInputChange = (e) => {
        const inputText = e.target.value;
        setSearchQuery(inputText);

        if (searchField) {
            clearTimeout(typingTimeout);
            setTypingTimeout(
                setTimeout(() => {
                    if (!inputText.trim()) {
                        onSearch('', '');
                    } else if (searchField && inputText.trim()) {
                        onSearch(searchField, inputText);
                    }
                }, 300)
            );
        }
    };

    const handleSearchFieldChange = (e) => {
        const newField = e.target.value;
        setSearchField(newField);

        if (newField === '') {
            setSearchQuery('');
            onSearch('', '');
        }
    };

    const clearSearch = () => {
        setSearchQuery('');
        onSearch(searchField, '');
    };

    const handleSearchButtonClick = () => {
        if (searchField && searchQuery) {
            onSearch(searchField, searchQuery);
        }
    };

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
            <Select
                value={searchField}
                onChange={handleSearchFieldChange}
                variant="unstyled"
                size="md"
                w="140px"
                pl={3}
                pr={0}
                py={2.5}
                borderRight="1px"
                borderColor={borderColor}
                borderRadius="0"
                _focus={{ boxShadow: "none" }}
                fontSize="sm"
            >
                <option value="">Chọn trường tìm kiếm</option>
                <option value="company_name">Tên công ty</option>
                <option value="tax_id">Mã số thuế</option>
                <option value="contact_phone">Số điện thoại</option>
            </Select>

            <InputGroup size="md" variant="unstyled">
                <InputLeftElement pointerEvents="none" h="full" pl={3}>
                    <FiSearch color="gray.400" />
                </InputLeftElement>
                <Input
                    placeholder={searchField ? `Tìm kiếm ${searchField === 'company_name' ? 'tên công ty' : searchField === 'tax_id' ? 'mã số thuế' : 'số điện thoại'}...` : "Chọn một trường tìm kiếm trước"}
                    pl={10}
                    pr={2}
                    py={2.5}
                    value={searchQuery}
                    onChange={handleSearchInputChange}
                    _placeholder={{ color: "gray.400" }}
                    isDisabled={!searchField}
                />
            </InputGroup>

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

export default SupplierSearchComponent;