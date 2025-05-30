import React from 'react';
import {
    Box,
    Flex,
    Text,
    Avatar,
    HStack,
    useColorModeValue,
    Menu,
    MenuButton,
    MenuList,
    MenuItem,
    MenuDivider,
} from '@chakra-ui/react';
import { FaHome } from 'react-icons/fa';
import LogoCompact from '../../ui/LogoCompact.jsx';
import { useNavigate } from "react-router-dom";
import useAuth from "../../../hooks/useAuth.js";

const SupplierHeader = () => {
    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const textColor = useColorModeValue('gray.800', 'gray.100');
    const navigate = useNavigate();
    const {user, logout} = useAuth()

    const handleLogout = async () => {
        try {
            await logout();
            navigate('/login', {replace: true});
        } catch (error) {
            console.error('Logout failed:', error);
        }
    };

    const handleLogoClick = () => {
        navigate('/');
    };

    return (
        <Box
            as="header"
            bg={bgColor}
            w="full"
            h="100%"
            py={3}
            px={6}
            borderBottom="1px"
            borderColor={borderColor}
            boxShadow="sm"
        >
            <Flex justify="space-between" align="center" h="100%">
                {/* Logo with text side by side */}
                <Flex
                    align="center"
                    gap={2}
                    cursor="pointer"
                    onClick={handleLogoClick}
                >
                    <LogoCompact size="sm" />
                    <Text
                        fontSize="xl"
                        fontWeight="bold"
                        color="gray.800"
                        _dark={{ color: 'gray.100' }}
                    >
                        Minh Plaza - Nhà Cung Cấp
                    </Text>
                </Flex>

                {/* User Profile */}
                <Menu isLazy placement="bottom-end">
                    <MenuButton
                        as={Box}
                        cursor="pointer"
                        borderRadius="md"
                        px={3}
                        py={2}
                        _hover={{ bg: 'gray.100' }}
                        _active={{ bg: 'gray.200' }}
                        transition="all 0.2s"
                    >
                        <HStack spacing={3}>
                            <Avatar
                                size="sm"
                                src={user.avatarUrl}
                                name={user.fullname}
                                bg="blue.50"
                                color="blue.700"
                                border="1px"
                                borderColor="gray.200"
                            />
                            <Text fontWeight="medium" color={textColor}>
                                {user.fullname}
                            </Text>
                        </HStack>
                    </MenuButton>
                    <MenuList
                        zIndex={1000}
                        p={0}
                        overflow="hidden"
                        borderRadius="md"
                        boxShadow="lg"
                    >
                        <MenuItem
                            py={3}
                            onClick={() => navigate("/user/account/profile")}
                            _hover={{ bg: 'gray.50' }}
                            w="full"
                        >
                            Tài khoản của tôi
                        </MenuItem>
                        <MenuDivider m={0} />
                        <MenuItem
                            onClick={() => navigate('/')}
                            py={3}
                            _hover={{ bg: 'gray.50' }}
                            w="full"
                        >
                            <Flex align="center">
                                <FaHome style={{ marginRight: '8px' }} />
                                Về trang chủ
                            </Flex>
                        </MenuItem>
                        <MenuDivider m={0} />
                        <MenuItem
                            onClick={handleLogout}
                            color="red.500"
                            py={3}
                            _hover={{ bg: 'red.50' }}
                            w="full"
                        >
                            Đăng xuất
                        </MenuItem>
                    </MenuList>
                </Menu>
            </Flex>
        </Box>
    );
};

export default SupplierHeader;