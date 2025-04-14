import React, { useState, useRef, useEffect } from 'react';
import {
    Box,
    Flex,
    Input,
    InputGroup,
    InputRightElement,
    Icon,
    List,
    ListItem,
    Portal,
    useOutsideClick,
    useColorModeValue
} from '@chakra-ui/react';
import { ChevronDownIcon, ChevronUpIcon } from '@chakra-ui/icons';

// Custom dropdown component for date selection that matches Shopee's style
export const DateDropdown = ({ label, options, value, onChange, width = "100%" }) => {
    const [isOpen, setIsOpen] = useState(false);
    const ref = useRef();
    const inputRef = useRef();
    const [dropdownPosition, setDropdownPosition] = useState({ top: 0, left: 0, width: 0 });

    // Colors
    const borderColor = useColorModeValue('gray.200', 'gray.600');
    const borderColorActive = useColorModeValue('red.500', 'red.300');
    const hoverBgColor = useColorModeValue('gray.100', 'gray.700');

    // Close dropdown when clicking outside
    useOutsideClick({
        ref: ref,
        handler: () => setIsOpen(false),
    });

    // Toggle dropdown
    const toggleDropdown = () => {
        setIsOpen(!isOpen);
    };

    // Update dropdown position when it opens
    useEffect(() => {
        if (isOpen && inputRef.current) {
            const rect = inputRef.current.getBoundingClientRect();
            setDropdownPosition({
                top: rect.bottom + window.scrollY,
                left: rect.left + window.scrollX,
                width: rect.width,
            });
        }
    }, [isOpen]);

    // Get display value
    const getDisplayValue = () => {
        if (!value) return label;

        // For month options that have "Tháng" prefix
        if (label === "Tháng") {
            return `Tháng ${value}`;
        }

        return value.toString();
    };

    return (
        <Box
            position="relative"
            width={width}
            ref={ref}
            onClick={toggleDropdown}
            cursor="pointer"
        >
            <InputGroup>
                <Input
                    ref={inputRef}
                    value={getDisplayValue()}
                    readOnly
                    cursor="pointer"
                    borderColor={isOpen ? borderColorActive : borderColor}
                    _focus={{ borderColor: borderColorActive }}
                    _hover={{ borderColor: isOpen ? borderColorActive : borderColor }}
                    borderRadius="sm"
                    pointerEvents="none" // Make the input not capture clicks
                />
                <InputRightElement pointerEvents="none">
                    <Icon
                        as={isOpen ? ChevronUpIcon : ChevronDownIcon}
                        color={isOpen ? "red.500" : "gray.500"}
                    />
                </InputRightElement>
            </InputGroup>

            {isOpen && (
                <Portal>
                    <Box
                        position="absolute"
                        top={`${dropdownPosition.top}px`}
                        left={`${dropdownPosition.left}px`}
                        width={`${dropdownPosition.width}px`}
                        zIndex={9999}
                        mt="1px"
                        borderWidth="1px"
                        borderColor={borderColor}
                        boxShadow="sm"
                        bg="white"
                        _dark={{ bg: 'gray.800' }}
                        maxH="200px" // Reduced height
                        overflowY="auto"
                        overflowX="hidden"
                        onClick={(e) => e.stopPropagation()} // Prevent click from closing dropdown
                    >
                        <List>
                            {options.map((option) => (
                                <ListItem
                                    key={option.value}
                                    px={4}
                                    py={1.5} // Smaller padding to reduce item height
                                    cursor="pointer"
                                    fontSize="sm"
                                    onClick={(e) => {
                                        e.stopPropagation(); // Prevent click from propagating
                                        onChange(option.value);
                                        setIsOpen(false);
                                    }}
                                    bg={value === option.value ? 'gray.100' : 'transparent'}
                                    _hover={{ bg: hoverBgColor }}
                                >
                                    {option.label}
                                </ListItem>
                            ))}
                        </List>
                    </Box>
                </Portal>
            )}
        </Box>
    );
};

// Grouped component for birth date selection
export const BirthDateSelector = ({ selectedDay, selectedMonth, selectedYear, onChange }) => {
    // Generate day options (1-31)
    const dayOptions = [...Array(31)].map((_, i) => ({
        value: i + 1,
        label: i + 1
    }));

    // Generate month options
    const monthOptions = [...Array(12)].map((_, i) => ({
        value: i + 1,
        label: `Tháng ${i + 1}`
    }));

    // Generate year options (current year down to 100 years ago)
    const currentYear = new Date().getFullYear();
    const yearOptions = [...Array(100)].map((_, i) => ({
        value: currentYear - i,
        label: currentYear - i
    }));

    return (
        <Flex gap={3}>
            <DateDropdown
                label="Ngày"
                options={dayOptions}
                value={selectedDay}
                onChange={(value) => onChange('day', value)}
                width="32%"
            />
            <DateDropdown
                label="Tháng"
                options={monthOptions}
                value={selectedMonth}
                onChange={(value) => onChange('month', value)}
                width="36%"
            />
            <DateDropdown
                label="Năm"
                options={yearOptions}
                value={selectedYear}
                onChange={(value) => onChange('year', value)}
                width="32%"
            />
        </Flex>
    );
};