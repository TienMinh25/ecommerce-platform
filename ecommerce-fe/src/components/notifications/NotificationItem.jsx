import React from 'react';
import {
    Box,
    Flex,
    Text,
    Image,
    Icon,
    Badge,
} from '@chakra-ui/react';
import { FaShoppingBag, FaCreditCard, FaBox, FaGift, FaBell, FaCheck } from 'react-icons/fa';
import { format } from 'date-fns';
import { NOTIFICATION_TYPES } from '../../constants/notificationTypes';

// Component for individual notification items - reusable in both popup and full page
const NotificationItem = ({ notification, onAction }) => {
    // Format the timestamp
    const formatTime = (timestamp) => {
        const date = new Date(timestamp);
        return format(date, 'HH:mm dd-MM-yyyy');
    };

    // Get appropriate icon based on notification type
    const getNotificationIcon = () => {
        switch (notification.type) {
            case NOTIFICATION_TYPES.ORDER:
                return <FaShoppingBag />;
            case NOTIFICATION_TYPES.PAYMENT:
                return <FaCreditCard />;
            case NOTIFICATION_TYPES.PRODUCT:
                return <FaBox />;
            case NOTIFICATION_TYPES.PROMOTION:
                return <FaGift />;
            case NOTIFICATION_TYPES.SYSTEM:
            default:
                return <FaBell />;
        }
    };

    // Get icon color based on notification type
    const getIconColor = () => {
        switch (notification.type) {
            case NOTIFICATION_TYPES.ORDER:
                return "blue.500";
            case NOTIFICATION_TYPES.PAYMENT:
                return "green.500";
            case NOTIFICATION_TYPES.PRODUCT:
                return "purple.500";
            case NOTIFICATION_TYPES.PROMOTION:
                return "red.500";
            case NOTIFICATION_TYPES.SYSTEM:
            default:
                return "gray.500";
        }
    };

    // Generate fallback image URL from Unsplash based on notification type
    const getFallbackImage = () => {
        switch (notification.type) {
            case NOTIFICATION_TYPES.ORDER:
                return "https://images.unsplash.com/photo-1556741533-6e6a62bd8b49?w=60&h=60&fit=crop";
            case NOTIFICATION_TYPES.PAYMENT:
                return "https://images.unsplash.com/photo-1580048915913-4f8f5cb481c4?w=60&h=60&fit=crop";
            case NOTIFICATION_TYPES.PRODUCT:
                return "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=60&h=60&fit=crop";
            case NOTIFICATION_TYPES.PROMOTION:
                return "https://images.unsplash.com/photo-1607083206968-13611e3d76db?w=60&h=60&fit=crop";
            case NOTIFICATION_TYPES.SYSTEM:
            default:
                return "https://images.unsplash.com/photo-1586769852044-692d6e3703f0?w=60&h=60&fit=crop";
        }
    };

    // Handle mark as read - prevent event propagation to parent
    const handleMarkAsRead = (e) => {
        e.stopPropagation();
        if (notification.id) {
            onAction && onAction(notification);
        }
    };

    return (
        <Flex
            p={4}
            borderBottom="1px"
            borderColor="gray.200"
            bg={notification.is_read ? "white" : "gray.50"}
            _hover={{ bg: "gray.100" }}
            position="relative"
        >
            <Image
                src={notification.image_url}
                alt="Notification"
                boxSize="60px"
                objectFit="cover"
                borderRadius="md"
                mr={3}
                fallbackSrc={getFallbackImage()}
            />
            <Box flex="1">
                <Flex alignItems="center" mb={1}>
                    <Icon as={getNotificationIcon} color={getIconColor()} mr={1} />
                    <Text fontWeight="medium">
                        {notification.title}
                    </Text>
                    {!notification.is_read && (
                        <Badge ml={2} colorScheme="red" borderRadius="full">
                            Mới
                        </Badge>
                    )}
                </Flex>
                <Text fontSize="sm" color="gray.600">
                    {notification.content}
                </Text>
                <Flex justifyContent="space-between" alignItems="center" mt={1}>
                    <Text fontSize="xs" color="gray.500">
                        {formatTime(notification.created_at)}
                    </Text>
                    {/* "Mark as read" text only for unread notifications */}
                    {!notification.is_read && (
                        <Text
                            fontSize="xs"
                            color="blue.500"
                            fontWeight="medium"
                            cursor="pointer"
                            _hover={{ textDecoration: "underline" }}
                            onClick={handleMarkAsRead}
                        >
                            Đánh dấu đã đọc
                        </Text>
                    )}
                </Flex>
            </Box>
        </Flex>
    );
};

export default NotificationItem;