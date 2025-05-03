import React, { useRef } from 'react';
import {
    Box,
    IconButton,
    Badge,
    Popover,
    PopoverTrigger,
    PopoverContent,
    useDisclosure,
} from '@chakra-ui/react';
import { FaBell } from 'react-icons/fa';
import NotificationPopup from './NotificationPopup';
import useNotification from "../../hooks/useNotification.js";

const NotificationBell = () => {
    const {
        notifications,
        unreadCount,
        isLoading,
        markAllAsRead
    } = useNotification();

    const { isOpen, onOpen, onClose } = useDisclosure();
    const buttonRef = useRef(null);

    return (
        <Box position="relative">
            <Popover
                isOpen={isOpen}
                onClose={onClose}
                placement="bottom-end"
                closeOnBlur={true}
                initialFocusRef={buttonRef}
                lazyBehavior="unmount"
            >
                <PopoverTrigger>
                    <Box>
                        <IconButton
                            ref={buttonRef}
                            aria-label="Notifications"
                            icon={<FaBell />}
                            variant="ghost"
                            onClick={onOpen}
                        />
                        {unreadCount > 0 && (
                            <Badge
                                position="absolute"
                                top="-2px"
                                right="-2px"
                                colorScheme="red"
                                borderRadius="full"
                                size="xs"
                            >
                                {unreadCount > 99 ? '99+' : unreadCount}
                            </Badge>
                        )}
                    </Box>
                </PopoverTrigger>
                <NotificationPopup
                    notifications={notifications}
                    isLoading={isLoading}
                    onMarkAllAsRead={markAllAsRead}
                    onClose={onClose}
                />
            </Popover>
        </Box>
    );
};

export default NotificationBell;