import React from 'react';
import {
    Box,
    Flex,
    Text,
    Image,
    Button,
    Icon,
} from '@chakra-ui/react';
import { FaHeart } from 'react-icons/fa';
import { format } from 'date-fns';

// Component for individual notification items - reusable in both popup and full page
const NotificationItem = ({ notification, onAction }) => {
    // Format the timestamp
    const formatTime = (timestamp) => {
        const date = new Date(timestamp);
        return format(date, 'HH:mm dd-MM-yyyy');
    };

    // Render different notification types
    const renderNotificationContent = () => {
        switch (notification.type) {
            case 'order_status':
                return (
                    <Box>
                        <Text fontWeight="medium" mb={1}>
                            Bạn có đơn hàng đang trên đường giao
                        </Text>
                        <Text fontSize="sm" color="gray.600">
                            🚚 Shipper báo rằng: đơn hàng {notification.orderId} của bạn vẫn đang trong quá trình vận chuyển và dự kiến được giao trong 1-2 ngày tới. Vui lòng bỏ qua thông báo này nếu bạn đã nhận được hàng nhé!😊
                        </Text>
                    </Box>
                );
            case 'voucher':
                return (
                    <Box>
                        <Flex alignItems="center" mb={1}>
                            <Icon as={FaHeart} color="red.500" mr={1} />
                            <Text fontWeight="medium">
                                Voucher dành riêng cho bạn
                            </Text>
                        </Flex>
                        <Text fontSize="sm" color="gray.600">
                            Shopee gửi bạn Voucher đ{notification.amount.toLocaleString()} thay lời xin lỗi cho đơn đã giao sau ngày Shopee đảm bảo. Lưu Voucher ngay! *Lưu ý: Nếu đơn hàng bị hủy hoặc có phát sinh yêu cầu Trả hàng/Hoàn tiền trước khi giao hàng thành công, Voucher sẽ không được áp dụng
                        </Text>
                    </Box>
                );
            case 'freeship':
                return (
                    <Box>
                        <Text fontWeight="medium" mb={1}>
                            Mã freeship cho đơn từ 0Đ có sẵn trong ví 😊
                        </Text>
                        <Text fontSize="sm" color="gray.600">
                            🎫 Mã sẽ hết hạn vào {notification.expiryDate}! Áp dụng cho đơn từ 0Đ👑
                        </Text>
                        <Text fontSize="sm" color="gray.600">
                            🏷️ Voucher freeship có sẵn trong ví, đừng ngay kẻo lỡ!
                        </Text>
                    </Box>
                );
            default:
                return (
                    <Text fontSize="sm">
                        {notification.message}
                    </Text>
                );
        }
    };

    return (
        <Flex
            p={4}
            borderBottom="1px"
            borderColor="gray.200"
            bg={notification.isRead ? "white" : "gray.50"}
            _hover={{ bg: "gray.100" }}
            cursor="pointer"
        >
            <Image
                src={notification.image}
                alt="Notification"
                boxSize="60px"
                objectFit="cover"
                borderRadius="md"
                mr={3}
                fallbackSrc="https://via.placeholder.com/60"
            />
            <Box flex="1">
                {renderNotificationContent()}
                <Text fontSize="xs" color="gray.500" mt={1}>
                    {formatTime(notification.timestamp)}
                </Text>
            </Box>
            {notification.hasAction && (
                <Button
                    size="sm"
                    colorScheme="red"
                    variant="solid"
                    height="36px"
                    mt={2}
                    fontSize="sm"
                    onClick={() => onAction && onAction(notification)}
                >
                    {notification.actionText || "Dùng ngay!"}
                </Button>
            )}
        </Flex>
    );
};

export default NotificationItem;