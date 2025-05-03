import React, { useEffect } from 'react';
import {
    Box,
    Container,
    Flex,
    Heading,
    Text,
    Badge,
    Button,
    Divider,
    Spinner,
    Center,
} from '@chakra-ui/react';
import NotificationItem from '../notifications/NotificationItem.jsx';
import useNotification from "../../hooks/useNotification.js";

const UserNotifications = () => {
    const {
        notifications,
        unreadCount,
        isLoading,
        fetchNotifications,
        markAllAsRead,
        markAsRead,
    } = useNotification();

    // Fetch all notifications with a larger limit when component mounts
    useEffect(() => {
        fetchNotifications(1, 20); // Fetch more notifications for the full page view
    }, [fetchNotifications]);

    // Group notifications by date for better display
    const groupNotificationsByDate = (notifications) => {
        const grouped = {};

        notifications.forEach(notification => {
            const date = new Date(notification.created_at);
            const dateStr = date.toLocaleDateString('vi-VN', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });

            if (!grouped[dateStr]) {
                grouped[dateStr] = [];
            }

            grouped[dateStr].push(notification);
        });

        return grouped;
    };

    // Handle notification click
    const handleNotificationClick = (notification) => {
        // Mark as read if not already read
        if (!notification.is_read) {
            markAsRead(notification.id);
        }
    };

    const groupedNotifications = groupNotificationsByDate(notifications);
    const dateGroups = Object.keys(groupedNotifications).sort((a, b) =>
        new Date(b) - new Date(a)
    );

    return (
        <Container maxW="container.lg" py={6}>
            <Box
                bg="white"
                borderRadius="lg"
                boxShadow="sm"
                overflow="hidden"
                borderWidth="1px"
                borderColor="gray.200"
            >
                {/* Header */}
                <Flex
                    p={4}
                    bg="gray.50"
                    borderBottomWidth="1px"
                    borderColor="gray.200"
                    justifyContent="space-between"
                    alignItems="center"
                    flexWrap="wrap"
                    gap={2}
                >
                    <Heading as="h1" size="lg" fontWeight="bold">
                        Thông báo
                        {unreadCount > 0 && (
                            <Badge ml={2} colorScheme="red" borderRadius="full">
                                {unreadCount} mới
                            </Badge>
                        )}
                    </Heading>

                    {unreadCount > 0 && (
                        <Button
                            size="sm"
                            colorScheme="brand"
                            variant="outline"
                            onClick={markAllAsRead}
                        >
                            Đánh dấu tất cả đã đọc
                        </Button>
                    )}
                </Flex>

                {/* Notification Content */}
                <Box minH="500px">
                    {isLoading ? (
                        <Center py={10}>
                            <Spinner
                                thickness="4px"
                                speed="0.65s"
                                emptyColor="gray.200"
                                color="brand.500"
                                size="xl"
                            />
                        </Center>
                    ) : notifications.length === 0 ? (
                        <Center py={10} flexDirection="column">
                            <Text color="gray.500" fontSize="lg" mb={2}>
                                Không có thông báo nào
                            </Text>
                            <Text color="gray.400" fontSize="sm">
                                Bạn chưa có thông báo nào.
                            </Text>
                        </Center>
                    ) : (
                        <Box>
                            {dateGroups.map((dateGroup) => (
                                <Box key={dateGroup}>
                                    {/* Date Header */}
                                    <Box
                                        bg="gray.50"
                                        px={4}
                                        py={2}
                                        borderBottomWidth="1px"
                                        borderColor="gray.200"
                                    >
                                        <Text fontWeight="medium" fontSize="sm" color="gray.600">
                                            {dateGroup}
                                        </Text>
                                    </Box>

                                    {/* Notifications for this date */}
                                    {groupedNotifications[dateGroup].map((notification) => (
                                        <Box key={notification.id}>
                                            <NotificationItem
                                                notification={notification}
                                                onAction={handleNotificationClick}
                                            />
                                            <Divider />
                                        </Box>
                                    ))}
                                </Box>
                            ))}
                        </Box>
                    )}
                </Box>
            </Box>
        </Container>
    );
};

export default UserNotifications;