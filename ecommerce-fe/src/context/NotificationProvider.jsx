import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { useToast } from '@chakra-ui/react';
import { useAuth } from '../hooks/useAuth';
import userMeService from '../services/userMeService';
import { NOTIFICATION_TYPES } from '../constants/notificationTypes';

// Create the notification context
const NotificationContext = createContext();

export const NotificationProvider = ({ children }) => {
    const [notifications, setNotifications] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [unreadCount, setUnreadCount] = useState(0);
    const toast = useToast();
    const { user } = useAuth();

    // Fetch notifications from API
    const fetchNotifications = useCallback(async (page = 1, limit = 5) => {
        if (!user) return;

        setIsLoading(true);
        try {
            const response = await userMeService.getNotifications({ page, limit });

            // Get notifications from response
            setNotifications(response.data || []);

            // Set unread count from metadata
            setUnreadCount(response.metadata?.unread || 0);
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
            // Call API to mark notification as read
            await userMeService.markNotificationAsRead(notificationId);

            // Update local state
            setNotifications(prev =>
                prev.map(n =>
                    n.id === notificationId ? { ...n, is_read: true } : n
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
            // Call API to mark all notifications as read
            await userMeService.markAllNotificationsAsRead();

            // Update local state
            setNotifications(prev =>
                prev.map(n => ({ ...n, is_read: true }))
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

    const value = {
        notifications,
        unreadCount,
        isLoading,
        fetchNotifications,
        markAsRead,
        markAllAsRead,
    };

    return (
        <NotificationContext.Provider value={value}>
            {children}
        </NotificationContext.Provider>
    );
};

export default NotificationContext;