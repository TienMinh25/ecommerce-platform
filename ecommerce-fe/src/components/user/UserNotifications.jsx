import React, { useEffect, useState } from 'react';
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
    HStack,
    IconButton,
} from '@chakra-ui/react';
import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import NotificationItem from '../notifications/NotificationItem.jsx';
import useNotification from "../../hooks/useNotification.js";

const UserNotifications = () => {
    const {
        notifications,
        unreadCount,
        isLoading,
        isPaginating,
        fetchNotifications,
        markAllAsRead,
        markAsRead,
        pagination
    } = useNotification();

    // Thêm state cho phân trang
    const [currentPage, setCurrentPage] = useState(1);
    // Thêm state lưu trữ dữ liệu hiện tại
    const [displayedNotifications, setDisplayedNotifications] = useState([]);
    const pageSize = 10;

    // Cập nhật displayedNotifications khi notifications thay đổi và không trong trạng thái paginating
    useEffect(() => {
        if (!isPaginating) {
            setDisplayedNotifications(notifications);
        }
    }, [notifications, isPaginating]);

    // Fetch notifications với phân trang
    useEffect(() => {
        fetchNotifications(currentPage, pageSize, currentPage > 1);
    }, [fetchNotifications, currentPage, pageSize]);

    // Group notifications by date for better display
    const groupNotificationsByDate = (notifications) => {
        const grouped = {};

        if (!notifications || !Array.isArray(notifications)) {
            return grouped;
        }

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

    // Xử lý chuyển trang
    const handleNextPage = () => {
        if (pagination && currentPage < pagination.total_pages) {
            setCurrentPage(currentPage + 1);
        }
    };

    const handlePrevPage = () => {
        if (currentPage > 1) {
            setCurrentPage(currentPage - 1);
        }
    };

    // Handle notification click
    const handleNotificationClick = (notification) => {
        // Mark as read if not already read
        if (!notification.is_read) {
            markAsRead(notification.id);
        }
    };

    // Sử dụng displayedNotifications thay vì notifications
    const groupedNotifications = groupNotificationsByDate(displayedNotifications);
    const dateGroups = Object.keys(groupedNotifications).sort((a, b) =>
        new Date(b) - new Date(a)
    );

    // Lấy giá trị totalPages từ pagination
    const totalPages = pagination?.total_pages || 1;

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
                <Box minH="500px" position="relative">
                    {/* Overlay loading khi đang chuyển trang */}
                    {isPaginating && (
                        <Box
                            position="absolute"
                            top="0"
                            left="0"
                            right="0"
                            bottom="0"
                            bg="rgba(255, 255, 255, 0.7)"
                            zIndex="10"
                            display="flex"
                            alignItems="center"
                            justifyContent="center"
                        >
                            <Spinner
                                thickness="4px"
                                speed="0.65s"
                                emptyColor="gray.200"
                                color="brand.500"
                                size="lg"
                            />
                        </Box>
                    )}

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
                    ) : !displayedNotifications || displayedNotifications.length === 0 ? (
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

                            <Flex justify="center" p={4} borderTop="1px" borderColor="gray.200">
                                <HStack spacing={4}>
                                    <IconButton
                                        icon={<ChevronLeftIcon />}
                                        aria-label="Previous page"
                                        onClick={handlePrevPage}
                                        isDisabled={currentPage === 1 || isPaginating}
                                        size="sm"
                                    />

                                    <Text fontSize="sm" fontWeight="medium">
                                        Trang {currentPage} / {totalPages}
                                    </Text>

                                    <IconButton
                                        icon={<ChevronRightIcon />}
                                        aria-label="Next page"
                                        onClick={handleNextPage}
                                        isDisabled={currentPage === totalPages || isPaginating}
                                        size="sm"
                                    />
                                </HStack>
                            </Flex>

                        </Box>
                    )}
                </Box>
            </Box>
        </Container>
    );
};

export default UserNotifications;