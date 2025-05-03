// Actualización para NotificationPopup.jsx - Remover el manejo de clics en el item completo
import React from 'react';
import {
    Box,
    Text,
    Flex,
    Button,
    PopoverContent,
    PopoverBody,
    PopoverHeader,
    PopoverArrow,
    Spinner,
    Center,
} from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';
import NotificationItem from './NotificationItem.jsx';

const NotificationPopup = ({ notifications = [], isLoading = false, unreadCount = 0, onMarkAllAsRead, onClose, onAction }) => {
    const navigate = useNavigate();

    // Handle View All button click
    const handleViewAll = () => {
        // Close the popup
        if (onClose) {
            onClose();
        }

        // Navigate to the notifications page
        navigate('/user/account/notifications/see');
    };

    return (
        <PopoverContent
            borderColor="gray.200"
            boxShadow="lg"
            borderRadius="md"
            width="420px"
            maxH="600px"
            overflowY="auto"
            _focus={{ outline: "none" }}
        >
            <PopoverArrow />
            <PopoverHeader py={3} px={4} borderBottomWidth="1px">
                <Flex justifyContent="space-between" alignItems="center">
                    <Text fontWeight="bold" fontSize="md">
                        Thông Báo Mới Nhận
                        {unreadCount > 0 && ` (${unreadCount})`}
                    </Text>
                    {unreadCount > 0 && (
                        <Button
                            variant="ghost"
                            size="sm"
                            colorScheme="brand"
                            onClick={onMarkAllAsRead}
                            fontSize="sm"
                        >
                            Đánh dấu đã đọc tất cả
                        </Button>
                    )}
                </Flex>
            </PopoverHeader>
            <PopoverBody p={0}>
                {isLoading ? (
                    <Center py={8}>
                        <Spinner
                            thickness="4px"
                            speed="0.65s"
                            emptyColor="gray.200"
                            color="brand.500"
                            size="lg"
                        />
                    </Center>
                ) : notifications.length === 0 ? (
                    <Center py={8} px={4} textAlign="center">
                        <Text color="gray.500">Bạn chưa có thông báo nào</Text>
                    </Center>
                ) : (
                    <Box>
                        {notifications.map((notification) => (
                            <NotificationItem
                                key={notification.id}
                                notification={notification}
                                onAction={onAction}
                            />
                        ))}

                        <Button
                            variant="ghost"
                            size="md"
                            width="full"
                            py={3}
                            onClick={handleViewAll}
                            borderTopWidth="1px"
                            borderTopColor="gray.200"
                            borderRadius="0"
                            fontWeight="medium"
                        >
                            Xem tất cả
                        </Button>
                    </Box>
                )}
            </PopoverBody>
        </PopoverContent>
    );
};

export default NotificationPopup;