import React, { useState } from 'react';
import {
    Box,
    Heading,
    Text,
    Divider,
    Switch,
    Stack,
    Flex,
    useColorModeValue,
    useToast
} from '@chakra-ui/react';

// Notification settings component based on Shopee's interface
const NotificationSettings = () => {
    const toast = useToast();

    // State for notification settings
    const [notifications, setNotifications] = useState({
        // Email notifications
        emailNotifications: true,
        emailOrderUpdates: true,
        emailPromotions: false,
        emailSurveys: true,

        // SMS notifications
        smsNotifications: true,
        smsPromotions: false,

        // Zalo notifications
        zaloNotifications: true,
        zaloPromotionsVN: true
    });

    // Handle switch toggle
    const handleToggle = (setting) => {
        setNotifications(prev => {
            const newSettings = { ...prev, [setting]: !prev[setting] };

            // If main category is turned off, turn off all sub-settings
            if (setting === 'emailNotifications' && !newSettings.emailNotifications) {
                newSettings.emailOrderUpdates = false;
                newSettings.emailPromotions = false;
                newSettings.emailSurveys = false;
            }

            if (setting === 'smsNotifications' && !newSettings.smsNotifications) {
                newSettings.smsPromotions = false;
            }

            if (setting === 'zaloNotifications' && !newSettings.zaloNotifications) {
                newSettings.zaloPromotionsVN = false;
            }

            // If any sub-setting is turned on, ensure main category is on
            if (setting === 'emailOrderUpdates' && newSettings.emailOrderUpdates) {
                newSettings.emailNotifications = true;
            }
            if (setting === 'emailPromotions' && newSettings.emailPromotions) {
                newSettings.emailNotifications = true;
            }
            if (setting === 'emailSurveys' && newSettings.emailSurveys) {
                newSettings.emailNotifications = true;
            }

            if (setting === 'smsPromotions' && newSettings.smsPromotions) {
                newSettings.smsNotifications = true;
            }

            if (setting === 'zaloPromotionsVN' && newSettings.zaloPromotionsVN) {
                newSettings.zaloNotifications = true;
            }

            // Show toast notification
            toast({
                title: 'Cài đặt đã được cập nhật',
                status: 'success',
                duration: 2000,
                isClosable: true,
            });

            return newSettings;
        });
    };

    // Styling
    const sectionBg = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');

    // Notification setting item component
    const NotificationItem = ({
                                  title,
                                  description,
                                  isChecked,
                                  onChange,
                                  isSubItem = false
                              }) => (
        <Flex
            justify="space-between"
            align="center"
            py={4}
            pl={isSubItem ? 6 : 0}
            borderBottomWidth={1}
            borderBottomColor={borderColor}
        >
            <Box>
                <Text fontWeight={isSubItem ? "normal" : "medium"} mb={1}>
                    {title}
                </Text>
                <Text fontSize="sm" color="gray.500" maxW="550px">
                    {description}
                </Text>
            </Box>
            <Switch
                colorScheme="green"
                size="lg"
                isChecked={isChecked}
                onChange={onChange}
                sx={{
                    '.chakra-switch__track': {
                        background: isChecked ? 'green.500' : 'gray.300',
                    }
                }}
            />
        </Flex>
    );

    return (
        <Box>
            <Heading as="h1" size="lg" mb={4}>Cài Đặt Thông Báo</Heading>
            <Text color="gray.500" fontSize="sm" mb={6}>
                Quản lý thông báo để không bỏ lỡ thông tin quan trọng
            </Text>

            <Stack spacing={8}>
                {/* Email Notifications */}
                <Box
                    bg={sectionBg}
                    p={6}
                    borderRadius="md"
                    boxShadow="sm"
                    borderWidth="1px"
                    borderColor={borderColor}
                >
                    <Heading as="h2" size="md" mb={4}>
                        Email thông báo
                    </Heading>
                    <Text fontSize="sm" color="gray.500" mb={4}>
                        Thông báo và nhắc nhở quan trọng về tài khoản sẽ không thể bị tắt
                    </Text>

                    <Divider mb={4} />

                    <NotificationItem
                        title="Email thông báo"
                        description="Thông báo và nhắc nhở quan trọng về tài khoản sẽ không thể bị tắt"
                        isChecked={notifications.emailNotifications}
                        onChange={() => handleToggle('emailNotifications')}
                    />

                    <NotificationItem
                        title="Cập nhật đơn hàng"
                        description="Cập nhật về tình trạng vận chuyển của tất cả các đơn hàng"
                        isChecked={notifications.emailOrderUpdates}
                        onChange={() => handleToggle('emailOrderUpdates')}
                        isSubItem
                    />

                    <NotificationItem
                        title="Khuyến mãi"
                        description="Cập nhật về các ưu đãi và khuyến mãi sắp tới"
                        isChecked={notifications.emailPromotions}
                        onChange={() => handleToggle('emailPromotions')}
                        isSubItem
                    />

                    <NotificationItem
                        title="Khảo sát"
                        description="Đồng ý nhận khảo sát để cho chúng tôi được lắng nghe bạn"
                        isChecked={notifications.emailSurveys}
                        onChange={() => handleToggle('emailSurveys')}
                        isSubItem
                    />
                </Box>

                {/* SMS Notifications */}
                <Box
                    bg={sectionBg}
                    p={6}
                    borderRadius="md"
                    boxShadow="sm"
                    borderWidth="1px"
                    borderColor={borderColor}
                >
                    <Heading as="h2" size="md" mb={4}>
                        Thông báo SMS
                    </Heading>
                    <Text fontSize="sm" color="gray.500" mb={4}>
                        Thông báo và nhắc nhở quan trọng về tài khoản sẽ không thể bị tắt
                    </Text>

                    <Divider mb={4} />

                    <NotificationItem
                        title="Thông báo SMS"
                        description="Thông báo và nhắc nhở quan trọng về tài khoản sẽ không thể bị tắt"
                        isChecked={notifications.smsNotifications}
                        onChange={() => handleToggle('smsNotifications')}
                    />

                    <NotificationItem
                        title="Khuyến mãi"
                        description="Cập nhật về các ưu đãi và khuyến mãi sắp tới"
                        isChecked={notifications.smsPromotions}
                        onChange={() => handleToggle('smsPromotions')}
                        isSubItem
                    />
                </Box>
            </Stack>
        </Box>
    );
};

export default NotificationSettings;