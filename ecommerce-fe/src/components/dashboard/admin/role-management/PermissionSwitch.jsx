import React from 'react';
import {
    Box,
    Tooltip,
    Switch
} from '@chakra-ui/react';

/**
 * Reusable Permission Switch component with tooltip
 */
const PermissionSwitch = ({ isChecked, onChange, permission }) => {
    // Get tooltip text based on permission type
    const getTooltipText = () => {
        switch(permission) {
            case 'read': return 'View permission';
            case 'create': return 'Create permission';
            case 'update': return 'Edit/Update permission';
            case 'delete': return 'Delete permission';
            case 'approve': return 'Approve permission';
            case 'reject': return 'Reject permission';
            default: return 'Toggle permission';
        }
    };

    // Important: Only attach the onChange to the Switch, not the Box
    return (
        <Tooltip label={getTooltipText()} hasArrow placement="top">
            <Box
                position="relative"
                display="flex"
                alignItems="center"
                justifyContent="center"
                borderRadius="md"
                p={0.5}
            >
                <Switch
                    colorScheme="blue"
                    size="sm"
                    isChecked={isChecked}
                    onChange={onChange} // Add onChange here only
                    // No additional onClick needed - just use the Switch's onChange
                />
            </Box>
        </Tooltip>
    );
};

export default PermissionSwitch;