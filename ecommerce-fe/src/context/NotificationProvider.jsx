import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { useToast } from '@chakra-ui/react';
import { useAuth } from '../hooks/useAuth';

// Sample notification data (replace with actual API call)
const MOCK_NOTIFICATIONS = [
    {
        id: 1,
        type: 'order_status',
        orderId: '250426AMN0CWDC',
        message: 'Your order is on the way',
        image: 'https://via.placeholder.com/60',
        timestamp: '2025-05-03T09:55:00',
        isRead: false,
        hasAction: false,
    },
    {
        id: 2,
        type: 'order_status',
        orderId: '250426AMEBNE9S',
        message: 'Your order is on the way',
        image: 'https://via.placeholder.com/60',
        timestamp: '2025-05-02T10:24:00',
        isRead: false,
        hasAction: false,
    },
    {
        id: 3,
        type: 'voucher',
        amount: 15000,
        message: 'You received a voucher',
        image: 'https://via.placeholder.com/60',
        timestamp: '2025-05-02T01:24:00',
        isRead: false,
        hasAction: true,
        actionText: 'Dùng ngay!',
    },
    {
        id: 4,
        type: 'freeship',
        expiryDate: '08-05-2025',
        message: 'Free shipping voucher',
        image: 'https://via.placeholder.com/60',
        timestamp: '2025-05-01T15:30:00',
        isRead: false,
        hasAction: true,
        actionText: 'Dùng ngay!',
    },
    {
        id: 5,
        type: 'voucher',
        amount: 15000,
        message: 'You received a voucher',
        image: 'https://via.placeholder.com/60',
        timestamp: '2025-05-02T00:54:00',
        isRead: true,
        hasAction: true,
        actionText: 'Dùng ngay!',
    },
];

// Create the notification context
const NotificationContext = createContext();

export const NotificationProvider = ({ children }) => {
    const [notifications, setNotifications] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [unreadCount, setUnreadCount] = useState(0);
    const toast = useToast();
    const { user } = useAuth();

    // Fetch notifications (mock implementation)
    const fetchNotifications = useCallback(async () => {
        if (!user) return;

        setIsLoading(true);
        try {
            // In a real app, this would be an API call
            // await api.get('/notifications')

            // Simulate API delay
            await new Promise(resolve => setTimeout(resolve, 500));

            // Mock data
            setNotifications(MOCK_NOTIFICATIONS);

            // Count unread notifications
            const unread = MOCK_NOTIFICATIONS.filter(n => !n.isRead).length;
            setUnreadCount(unread);
        } catch (error) {
            console.error('Error fetching notifications:', error);
            toast({
                title: 'Không thể tải thông báo',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    }, [user, toast]);

    // Load initial notifications
    useEffect(() => {
        if (user) {
            fetchNotifications();
        }
    }, [user, fetchNotifications]);

    // Mark a notification as read
    const markAsRead = useCallback(async (notificationId) => {
        if (!user) return;

        try {
            // In a real app, this would be an API call
            // await api.put(`/notifications/${notificationId}/read`)

            // Update local state
            setNotifications(prev =>
                prev.map(n =>
                    n.id === notificationId ? { ...n, isRead: true } : n
                )
            );

            // Update unread count
            setUnreadCount(prev => Math.max(0, prev - 1));
        } catch (error) {
            console.error('Error marking notification as read:', error);
        }
    }, [user]);

    // Mark all notifications as read
    const markAllAsRead = useCallback(async () => {
        if (!user) return;

        try {
            // In a real app, this would be an API call
            // await api.put('/notifications/read-all')

            // Update local state
            setNotifications(prev =>
                prev.map(n => ({ ...n, isRead: true }))
            );

            // Update unread count
            setUnreadCount(0);

            toast({
                title: 'Đã đánh dấu tất cả thông báo là đã đọc',
                status: 'success',
                duration: 2000,
                isClosable: true,
            });
        } catch (error) {
            console.error('Error marking all notifications as read:', error);
        }
    }, [user, toast]);

    // Handle notification actions (e.g., use voucher)
    const handleNotificationAction = useCallback((notification) => {
        switch (notification.type) {
            case 'voucher':
            case 'freeship':
                toast({
                    title: 'Đã lưu voucher',
                    description: 'Voucher đã được lưu vào ví của bạn',
                    status: 'success',
                    duration: 3000,
                    isClosable: true,
                });
                break;
            default:
                break;
        }

        // Mark as read after action
        markAsRead(notification.id);
    }, [markAsRead, toast]);

    const value = {
        notifications,
        unreadCount,
        isLoading,
        fetchNotifications,
        markAsRead,
        markAllAsRead,
        handleNotificationAction
    };

    return (
        <NotificationContext.Provider value={value}>
            {children}
        </NotificationContext.Provider>
    );
};

export default NotificationContext;