import React from 'react';
import { Outlet, NavLink, useLocation } from 'react-router-dom';
import {
    Avatar,
    Box,
    Container,
    Flex,
    Text,
    VStack,
    useColorModeValue
} from '@chakra-ui/react';
import { FaUser, FaLock, FaMapMarkerAlt, FaShoppingBag, FaHeart, FaBell } from 'react-icons/fa';
import { useAuth } from '../../hooks/useAuth';

// This component is designed to work within MainLayout's Outlet
const UserAccountLayout = () => {
    const { user } = useAuth();
    const location = useLocation();

    // Menu items with their respective icons and paths
    const menuItems = [
        { label: 'Hồ sơ của tôi', icon: <FaUser />, path: '/user/account/profile' },
        { label: 'Đổi mật khẩu', icon: <FaLock />, path: '/user/account/password' },
        { label: 'Địa chỉ', icon: <FaMapMarkerAlt />, path: '/user/account/addresses' },
        { label: 'Đơn hàng của tôi', icon: <FaShoppingBag />, path: '/user/account/orders' },
        { label: 'Cài đặt thông báo', icon: <FaBell />, path: '/user/account/notifications' },
    ];

    // Check if a menu item is active
    const isActive = (path) => location.pathname === path;
    const activeLinkBg = useColorModeValue('red.50', 'red.900');
    const activeLinkColor = useColorModeValue('red.500', 'red.200');

    return (
        <Container maxW="container.xl" py={6}>
            <Flex
                direction={{ base: 'column', md: 'row' }}
                gap={5}
            >
                {/* Left Sidebar - User Account Navigation */}
                <Box
                    w={{ base: 'full', md: '220px' }}
                    flexShrink={0}
                    bg={useColorModeValue('white', 'gray.800')}
                    borderRadius="md"
                    borderWidth="1px"
                    borderColor={useColorModeValue('gray.200', 'gray.700')}
                    boxShadow="sm"
                    overflow="hidden"
                    h="fit-content"
                >
                    {/* User Info */}
                    <Flex p={4} align="center" gap={3}>
                        <Avatar
                            size="md"
                            name={user?.fullname}
                            src={user?.avatarUrl}
                        />
                        <Box>
                            <Text fontWeight="bold" fontSize="sm" noOfLines={1} mr={2}>
                                {user?.fullname || 'User'}
                            </Text>
                            <Text fontSize="xs" color="gray.500" noOfLines={1}>
                                {user?.email && user.email.length > 16
                                    ? user.email.substring(0, 16) + '...'
                                    : user?.email || 'user@example.com'}
                            </Text>
                        </Box>
                    </Flex>

                    {/* Navigation */}
                    <VStack align="stretch" spacing={0} mt={2}>
                        {menuItems.map((item) => (
                            <NavLink
                                key={item.path}
                                to={item.path}
                                style={{ textDecoration: 'none' }}
                            >
                                <Flex
                                    px={4}
                                    py={2.5}
                                    align="center"
                                    gap={3}
                                    bg={isActive(item.path) ? activeLinkBg : 'transparent'}
                                    color={isActive(item.path) ? activeLinkColor : 'gray.700'}
                                    fontWeight={isActive(item.path) ? 'medium' : 'normal'}
                                    fontSize="sm"
                                    _hover={{ bg: 'gray.100' }}
                                    _dark={{
                                        color: isActive(item.path) ? activeLinkColor : 'gray.200',
                                        _hover: { bg: 'gray.800' }
                                    }}
                                    borderLeftWidth={isActive(item.path) ? '3px' : '0'}
                                    borderLeftColor={activeLinkColor}
                                    transition="all 0.2s"
                                >
                                    {item.icon}
                                    <Text>{item.label}</Text>
                                </Flex>
                            </NavLink>
                        ))}
                    </VStack>
                </Box>

                {/* Main Content Area */}
                <Box
                    flex="1"
                    bg={useColorModeValue('white', 'gray.800')}
                    p={5}
                    borderRadius="md"
                    boxShadow="sm"
                    borderWidth="1px"
                    borderColor={useColorModeValue('gray.200', 'gray.700')}
                    minH="500px"
                >
                    <Outlet />
                </Box>
            </Flex>
        </Container>
    );
};

export default UserAccountLayout;