import React, { useEffect, useRef, useState } from 'react';
import {
    Box,
    Button,
    Collapse,
    Flex,
    FormControl,
    FormLabel,
    IconButton,
    Text,
    useColorModeValue,
    useDisclosure,
    VStack,
} from '@chakra-ui/react';
import { FiFilter, FiX } from 'react-icons/fi';
import {ModernToggleSwitch} from "../../../toggle/main.jsx";

const SupplierFilterDropdown = ({ filters, onFiltersChange, onApplyFilters }) => {
    const { isOpen, onToggle, onClose } = useDisclosure();
    const filterRef = useRef(null);
    const dropdownRef = useRef(null);

    const filterBgColor = useColorModeValue('white', 'gray.800');
    const filterBorderColor = useColorModeValue('gray.200', 'gray.600');
    const headerBgColor = useColorModeValue('gray.50', 'gray.700');
    const buttonHoverColor = useColorModeValue('blue.50', 'blue.700');

    const handleToggleChange = (name, value) => {
        onFiltersChange({
            ...filters,
            [name]: value
        });
    };

    const handleApplyFilters = () => {
        const activeFilters = Object.fromEntries(
            Object.entries(filters).filter(([_, v]) => v !== '')
        );
        onApplyFilters(activeFilters);
        onClose();
    };

    const handleResetFilters = () => {
        onFiltersChange({
            status: '',
        });
    };

    const statusOptions = [
        { label: 'Chờ duyệt', value: 'pending' },
        { label: 'Hoạt động', value: 'active' },
        { label: 'Tạm ngưng', value: 'suspended' }
    ];

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (
                isOpen &&
                filterRef.current &&
                !filterRef.current.contains(event.target) &&
                dropdownRef.current &&
                !dropdownRef.current.contains(event.target)
            ) {
                onClose();
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [isOpen, onClose]);

    return (
        <Box position="relative" zIndex={1000} ref={filterRef}>
            <Button
                leftIcon={<FiFilter />}
                onClick={onToggle}
                variant={isOpen ? "solid" : "outline"}
                colorScheme={isOpen ? "blue" : "gray"}
                size="sm"
                borderRadius="md"
                px={4}
            >
                Bộ lọc
            </Button>

            <Collapse in={isOpen} animateOpacity>
                <Box
                    position="absolute"
                    top="40px"
                    right="0"
                    width="400px"
                    maxW="calc(100vw - 20px)"
                    bg={filterBgColor}
                    boxShadow="0 4px 20px rgba(0,0,0,0.15)"
                    borderRadius="lg"
                    borderWidth="1px"
                    borderColor={filterBorderColor}
                    overflow="hidden"
                    zIndex={1000}
                    ref={dropdownRef}
                >
                    <Flex
                        justifyContent="space-between"
                        alignItems="center"
                        p={4}
                        borderBottomWidth="1px"
                        borderColor={filterBorderColor}
                        bg={headerBgColor}
                    >
                        <Text fontWeight="bold" fontSize="lg" color="gray.800">
                            Bộ lọc nhà cung cấp
                        </Text>
                        <IconButton
                            icon={<FiX />}
                            variant="ghost"
                            size="sm"
                            aria-label="Close filters"
                            onClick={onClose}
                            color="gray.600"
                            _hover={{ bg: "gray.200" }}
                        />
                    </Flex>

                    <Box p={4}>
                        <VStack spacing={5} align="stretch" width={'100%'}>
                            <FormControl>
                                <FormLabel fontWeight="medium" fontSize="sm" mb={1} color="gray.700">
                                    Trạng thái nhà cung cấp
                                </FormLabel>
                                <VStack spacing={3} align="stretch">
                                    {statusOptions.map((option) => (
                                        <ModernToggleSwitch
                                            key={option.value}
                                            isChecked={filters.status === option.value}
                                            onChange={handleToggleChange}
                                            label={option.label}
                                            name="status"
                                            value={option.value}
                                        />
                                    ))}
                                </VStack>
                            </FormControl>
                        </VStack>
                    </Box>

                    <Flex
                        justify="space-between"
                        p={4}
                        borderTopWidth="1px"
                        borderColor={filterBorderColor}
                        bg={headerBgColor}
                    >
                        <Button
                            variant="outline"
                            colorScheme="blue"
                            onClick={handleResetFilters}
                            size="sm"
                            width="45%"
                            borderRadius="md"
                            _hover={{ bg: buttonHoverColor }}
                            height="40px"
                        >
                            Reset bộ lọc
                        </Button>
                        <Button
                            colorScheme="blue"
                            onClick={handleApplyFilters}
                            size="sm"
                            width="45%"
                            borderRadius="md"
                            _hover={{ bg: "blue.600" }}
                            height="40px"
                        >
                            Áp dụng bộ lọc
                        </Button>
                    </Flex>
                </Box>
            </Collapse>
        </Box>
    );
};

export default SupplierFilterDropdown;