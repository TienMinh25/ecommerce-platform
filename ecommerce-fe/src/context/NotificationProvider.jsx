import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { useToast } from '@chakra-ui/react';
import { useAuth } from '../hooks/useAuth';
import userMeService from '../services/userMeService';
import { NOTIFICATION_TYPES } from '../constants/notificationTypes';

// Create the notification context
const NotificationContext = createContext();

export const NotificationProvider = ({ children }) => {
    // Trạng thái hiện tại
    const [notificationResponse, setNotificationResponse] = useState({
        data: [],
        metadata: {
            pagination: {
                page: 1,
                limit: 10,
                total_items: 0,
                total_pages: 1,
                has_next: false,
                has_previous: false
            },
            unread: 0
        }
    });
    const [isLoading, setIsLoading] = useState(false);
    const [isPaginating, setIsPaginating] = useState(false);
    const toast = useToast();
    const { user } = useAuth();

    // Getter cho notifications và unreadCount
    const notifications = notificationResponse.data || [];
    const unreadCount = notificationResponse.metadata?.unread || 0;

    // Fetch notifications from API
    const fetchNotifications = useCallback(async (page = 1, limit = 10, isPaging = false) => {
        if (!user) return;

        // Nếu là phân trang, sử dụng trạng thái isPaginating
        if (isPaging) {
            setIsPaginating(true);
        } else {
            setIsLoading(true);
        }

        try {
            const response = await userMeService.getNotifications({ page, limit });

            // Lưu toàn bộ response
            setNotificationResponse(response);
        } catch (error) {
            console.error('Error fetching notifications:', error);
            toast({
                title: 'Không thể tải thông báo',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            // Tắt loading tương ứng
            if (isPaging) {
                setIsPaginating(false);
            } else {
                setIsLoading(false);
            }
        }
    }, [user, toast]);

    // Load initial notifications
    useEffect(() => {
        if (user) {
            fetchNotifications(1, 10);
        }
    }, [user, fetchNotifications]);

    // Mark a notification as read
    const markAsRead = useCallback(async (notificationId) => {
        if (!user || !notificationId) return;

        try {
            // Call API to mark notification as read
            await userMeService.markNotificationAsRead(notificationId);

            // Update local state
            setNotificationResponse(prev => {
                // Cập nhật data
                const updatedData = prev.data.map(n =>
                    n.id === notificationId ? { ...n, is_read: true } : n
                );

                // Cập nhật metadata.unread
                const updatedUnread = Math.max(0, (prev.metadata?.unread || 0) - 1);

                return {
                    ...prev,
                    data: updatedData,
                    metadata: {
                        ...prev.metadata,
                        unread: updatedUnread
                    }
                };
            });
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
            setNotificationResponse(prev => {
                // Cập nhật tất cả thông báo là đã đọc
                const updatedData = prev.data.map(n => ({ ...n, is_read: true }));

                return {
                    ...prev,
                    data: updatedData,
                    metadata: {
                        ...prev.metadata,
                        unread: 0
                    }
                };
            });

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
        isPaginating,
        fetchNotifications,
        markAsRead,
        markAllAsRead,
        pagination: notificationResponse.metadata?.pagination
    };

    return (
        <NotificationContext.Provider value={value}>
            {children}
        </NotificationContext.Provider>
    );
};

export default NotificationContext;