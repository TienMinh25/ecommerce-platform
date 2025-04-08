import React from 'react';
import {
    Box,
    Flex,
    Switch,
    Text,
    HStack,
    FormControl,
    FormLabel,
} from '@chakra-ui/react';

// Modern Switch Toggle Component
const ModernToggleSwitch = ({ isChecked, onChange, label, name, value }) => {
    return (
        <Flex
            align="center"
            justify="space-between"
            bg="gray.100"
            borderRadius="full"
            p={1}
            px={3}
            height="36px"
            cursor="pointer"
            onClick={() => onChange(name, value)}
        >
            <Text
                fontSize="sm"
                fontWeight="medium"
                color={isChecked ? "blue.600" : "gray.500"}
                transition="color 0.2s ease"
            >
                {label}
            </Text>
            <Switch
                isChecked={isChecked}
                size="md"
                colorScheme="blue"
                onChange={() => onChange(name, value)}
            />
        </Flex>
    );
};

// Group of toggle switches for binary options
const ModernToggleGroup = ({ options, value, onChange, name }) => {
    return (
        <HStack spacing={4} width="100%">
            {options.map((option) => (
                <ModernToggleSwitch
                    key={option.value}
                    isChecked={value === option.value}
                    onChange={onChange}
                    label={option.label}
                    name={name}
                    value={option.value}
                />
            ))}
        </HStack>
    );
};

// Example usage in a form control
const EmailVerificationToggle = ({ value, onChange }) => {
    return (
        <FormControl>
            <FormLabel fontWeight="medium" fontSize="sm" mb={2} color="gray.700">
                Email Verification
            </FormLabel>
            <ModernToggleGroup
                options={[
                    { label: 'Verified', value: 'true' },
                    { label: 'Not Verified', value: 'false' }
                ]}
                value={value}
                onChange={onChange}
                name="emailVerify"
            />
        </FormControl>
    );
};

// Example usage in a form control
const PhoneVerificationToggle = ({ value, onChange }) => {
    return (
        <FormControl>
            <FormLabel fontWeight="medium" fontSize="sm" mb={2} color="gray.700">
                Phone Verification
            </FormLabel>
            <ModernToggleGroup
                options={[
                    { label: 'Verified', value: 'true' },
                    { label: 'Not Verified', value: 'false' }
                ]}
                value={value}
                onChange={onChange}
                name="phoneVerify"
            />
        </FormControl>
    );
};

// Example usage in a form control
const StatusToggle = ({ value, onChange }) => {
    return (
        <FormControl>
            <FormLabel fontWeight="medium" fontSize="sm" mb={2} color="gray.700">
                Status
            </FormLabel>
            <ModernToggleGroup
                options={[
                    { label: 'Active', value: 'active' },
                    { label: 'Inactive', value: 'inactive' }
                ]}
                value={value}
                onChange={onChange}
                name="status"
            />
        </FormControl>
    );
};

export { ModernToggleGroup, ModernToggleSwitch, EmailVerificationToggle, PhoneVerificationToggle, StatusToggle };