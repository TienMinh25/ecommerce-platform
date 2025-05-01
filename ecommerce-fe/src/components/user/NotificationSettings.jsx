import React, { useState, useEffect } from 'react';
import {
    Box,
    Heading,
    Text,
    Divider,
    Switch,
    Stack,
    Flex,
    useColorModeValue,
    useToast,
    Spinner,
    Center,
    Alert,
    AlertIcon,
    AlertTitle,
    AlertDescription
} from '@chakra-ui/react';
import userMeService from "../../services/userMeService.js";

// Component cài đặt thông báo dựa trên tài liệu API
const NotificationSettings = () => {
    const toast = useToast();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    // State cho cài đặt thông báo
    const [settings, setSettings] = useState({
        email_setting: {
            order_status: false,
            payment_status: false,
            product_status: false,
            promotion: false
        },
        in_app_setting: {
            order_status: false,
            payment_status: false,
            product_status: false,
            promotion: false
        }
    });

    // Lấy cài đặt thông báo
    useEffect(() => {
        const fetchSettings = async () => {
            try {
                setLoading(true);
                const data = await userMeService.getNotificationSettings();
                setSettings(data);
                setError(null);
            } catch (err) {
                setError('Không thể tải cài đặt thông báo. Vui lòng thử lại sau.');
                console.error('Lỗi khi lấy cài đặt thông báo:', err);
            } finally {
                setLoading(false);
            }
        };

        fetchSettings();
    }, []);

    // Xử lý khi chuyển đổi cài đặt email
    const handleEmailToggle = (setting) => {
        const newSettings = {
            ...settings,
            email_setting: {
                ...settings.email_setting,
                [setting]: !settings.email_setting[setting]
            }
        };

        setSettings(newSettings);
        saveSettings(newSettings);
    };

    // Xử lý khi chuyển đổi cài đặt trong ứng dụng
    const handleInAppToggle = (setting) => {
        const newSettings = {
            ...settings,
            in_app_setting: {
                ...settings.in_app_setting,
                [setting]: !settings.in_app_setting[setting]
            }
        };

        setSettings(newSettings);
        saveSettings(newSettings);
    };

    // Lưu cài đặt vào API
    const saveSettings = async (newSettings) => {
        try {
            await userMeService.updateNotificationSettings(newSettings);
            toast({
                title: 'Cài đặt đã được cập nhật',
                status: 'success',
                duration: 2000,
                isClosable: true,
            });
        } catch (err) {
            toast({
                title: 'Không thể cập nhật cài đặt',
                description: 'Đã xảy ra lỗi khi lưu cài đặt của bạn.',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
            console.error('Lỗi khi lưu cài đặt thông báo:', err);
        }
    };

    // Styling
    const sectionBg = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');

    // Component cho mỗi mục cài đặt thông báo
    const NotificationItem = ({
                                  title,
                                  description,
                                  isChecked,
                                  onChange
                              }) => (
        <Flex
            justify="space-between"
            align="center"
            py={4}
            borderBottomWidth={1}
            borderBottomColor={borderColor}
        >
            <Box>
                <Text fontWeight="medium" mb={1}>
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

    if (loading) {
        return (
            <Center h="200px">
                <Spinner size="xl" color="green.500" thickness="4px" />
            </Center>
        );
    }

    if (error) {
        return (
            <Alert status="error" borderRadius="md" mb={6}>
                <AlertIcon />
                <AlertTitle mr={2}>Lỗi!</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
            </Alert>
        );
    }

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
                        Nhận thông báo qua email
                    </Text>

                    <Divider mb={4} />

                    <NotificationItem
                        title="Trạng thái đơn hàng"
                        description="Cập nhật về tình trạng vận chuyển của tất cả các đơn hàng"
                        isChecked={settings.email_setting.order_status}
                        onChange={() => handleEmailToggle('order_status')}
                    />

                    <NotificationItem
                        title="Trạng thái thanh toán"
                        description="Nhận thông báo về trạng thái thanh toán của đơn hàng"
                        isChecked={settings.email_setting.payment_status}
                        onChange={() => handleEmailToggle('payment_status')}
                    />

                    <NotificationItem
                        title="Trạng thái sản phẩm"
                        description="Cập nhật về sự thay đổi trạng thái của sản phẩm"
                        isChecked={settings.email_setting.product_status}
                        onChange={() => handleEmailToggle('product_status')}
                    />

                    <NotificationItem
                        title="Khuyến mãi"
                        description="Cập nhật về các ưu đãi và khuyến mãi sắp tới"
                        isChecked={settings.email_setting.promotion}
                        onChange={() => handleEmailToggle('promotion')}
                    />
                </Box>

                {/* In-App Notifications */}
                <Box
                    bg={sectionBg}
                    p={6}
                    borderRadius="md"
                    boxShadow="sm"
                    borderWidth="1px"
                    borderColor={borderColor}
                >
                    <Heading as="h2" size="md" mb={4}>
                        Thông báo trong ứng dụng
                    </Heading>
                    <Text fontSize="sm" color="gray.500" mb={4}>
                        Nhận thông báo trong ứng dụng
                    </Text>

                    <Divider mb={4} />

                    <NotificationItem
                        title="Trạng thái đơn hàng"
                        description="Cập nhật về tình trạng vận chuyển của tất cả các đơn hàng"
                        isChecked={settings.in_app_setting.order_status}
                        onChange={() => handleInAppToggle('order_status')}
                    />

                    <NotificationItem
                        title="Trạng thái thanh toán"
                        description="Nhận thông báo về trạng thái thanh toán của đơn hàng"
                        isChecked={settings.in_app_setting.payment_status}
                        onChange={() => handleInAppToggle('payment_status')}
                    />

                    <NotificationItem
                        title="Trạng thái sản phẩm"
                        description="Cập nhật về sự thay đổi trạng thái của sản phẩm"
                        isChecked={settings.in_app_setting.product_status}
                        onChange={() => handleInAppToggle('product_status')}
                    />

                    <NotificationItem
                        title="Khuyến mãi"
                        description="Cập nhật về các ưu đãi và khuyến mãi sắp tới"
                        isChecked={settings.in_app_setting.promotion}
                        onChange={() => handleInAppToggle('promotion')}
                    />
                </Box>
            </Stack>
        </Box>
    );
};

export default NotificationSettings;