import React, { useState, useRef, useEffect } from 'react';
import {
    Box,
    Flex,
    Icon,
    Text,
    IconButton,
    VStack,
    Tooltip,
    useColorModeValue,
    useBreakpointValue,
} from '@chakra-ui/react';
import {
    FiPackage,
    FiCreditCard,
    FiLogOut,
    FiChevronsLeft,
    FiChevronsRight,
} from 'react-icons/fi';
import { useNavigate, useLocation } from 'react-router-dom';
import { useSupplierSidebar } from '../../layout/SupplierLayout.jsx';
import useAuth from "../../../hooks/useAuth.js";

const SupplierSidebar = ({ onStateChange }) => {
    const { updateSidebarState } = useSupplierSidebar();
    const {logout} = useAuth()

    const [isCollapsed, setIsCollapsed] = useState(false);
    const [width, setWidth] = useState(256);
    const [isResizing, setIsResizing] = useState(false);
    const [isLogoutHovered, setIsLogoutHovered] = useState(false);
    const minWidth = 180;
    const maxWidth = 400;
    const collapsedWidth = 64;

    const sidebarRef = useRef(null);

    const navigate = useNavigate();
    const location = useLocation();

    const isMobile = useBreakpointValue({ base: true, md: false });
    const defaultCollapsed = useBreakpointValue({ base: true, md: false });

    const handleLogout = async () => {
        try {
            await logout();
            navigate('/login', {replace: true});
        } catch (error) {
            console.error('Logout failed:', error);
        }
    };

    useEffect(() => {
        setIsCollapsed(defaultCollapsed);
    }, [defaultCollapsed]);

    useEffect(() => {
        if (onStateChange) {
            onStateChange({
                isCollapsed,
                width,
                collapsedWidth
            });
        }

        if (updateSidebarState) {
            updateSidebarState({
                isCollapsed,
                width,
                collapsedWidth
            });
        }
    }, [isCollapsed, width, collapsedWidth, onStateChange, updateSidebarState]);

    // Menu items for supplier
    const menuItems = [
        { icon: FiPackage, label: 'Quản lý đơn hàng', path: '/supplier/orders' },
        { icon: FiCreditCard, label: 'Quản lý thanh toán', path: '/supplier/payments' },
    ];

    const logoutItem = { icon: FiLogOut, label: 'Đăng xuất', path: '/logout', logout: handleLogout};

    // Theme colors
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const bgColor = useColorModeValue('white', 'gray.800');
    const activeItemBg = useColorModeValue('cyan.50', 'cyan.900');
    const activeItemBorder = useColorModeValue('cyan.500', 'cyan.200');
    const hoverBg = useColorModeValue('cyan.500', 'cyan.600');
    const textColor = useColorModeValue('gray.700', 'gray.200');
    const iconColor = useColorModeValue('gray.500', 'gray.400');
    const activeIconColor = useColorModeValue('cyan.500', 'cyan.200');
    const resizeHandleColor = useColorModeValue('gray.200', 'gray.600');
    const resizeHandleHoverColor = useColorModeValue('cyan.500', 'cyan.400');

    const toggleSidebar = (e) => {
        if (e) {
            e.preventDefault();
            e.stopPropagation();
        }

        setIsCollapsed(prevState => !prevState);

        if (sidebarRef.current) {
            sidebarRef.current.style.transition = "width 0.25s cubic-bezier(0.34, 1.56, 0.64, 1)";

            setTimeout(() => {
                if (sidebarRef.current) {
                    sidebarRef.current.style.transition = "width 0.2s ease";
                }
            }, 250);
        }
    };

    const handleResizeStart = (e) => {
        e.preventDefault();
        setIsResizing(true);

        if (isCollapsed) {
            setIsCollapsed(false);
            setWidth(minWidth);
        }

        const startX = e.clientX;
        const startWidth = isCollapsed ? minWidth : width;

        const handleResize = (e) => {
            const newWidth = Math.min(
                Math.max(startWidth + e.clientX - startX, minWidth),
                maxWidth
            );
            setWidth(newWidth);

            const menuTexts = document.querySelectorAll('.menu-item-text');
            menuTexts.forEach(text => {
                if (newWidth < 220) {
                    text.style.whiteSpace = 'normal';
                    text.style.wordBreak = 'break-word';
                    text.style.maxHeight = '36px';
                } else {
                    text.style.whiteSpace = 'nowrap';
                    text.style.wordBreak = 'normal';
                    text.style.maxHeight = 'none';
                }
            });
        };

        const handleResizeEnd = () => {
            setIsResizing(false);
            document.removeEventListener('mousemove', handleResize);
            document.removeEventListener('mouseup', handleResizeEnd);
        };

        document.addEventListener('mousemove', handleResize);
        document.addEventListener('mouseup', handleResizeEnd);
    };

    const MenuItem = ({ item, isActive, isLogoutItem = false }) => {
        const [isHovered, setIsHovered] = useState(false);

        const handleClick = (e) => {
            navigate(item.path);
        };

        const handleMouseEnter = () => {
            setIsHovered(true);
            if (isLogoutItem) {
                setIsLogoutHovered(true);
            }
        };

        const handleMouseLeave = () => {
            setIsHovered(false);
            if (isLogoutItem) {
                setIsLogoutHovered(false);
            }
        };

        return (
            <Tooltip
                label={isCollapsed ? item.label : ""}
                placement="right"
                hasArrow
                bg="cyan.500"
                color="white"
                fontWeight="medium"
                px={3}
                py={2}
                borderRadius="md"
                isDisabled={!isCollapsed}
                openDelay={200}
                gutter={12}
            >
                <Box
                    position="relative"
                    width="100%"
                    mb={2}
                    className="menu-item-container"
                    onMouseEnter={handleMouseEnter}
                    onMouseLeave={handleMouseLeave}
                >
                    <Box
                        position="absolute"
                        top="0"
                        left="0"
                        right="0"
                        bottom="0"
                        borderRadius="md"
                        borderLeft="3px solid"
                        borderLeftColor={isActive ? activeItemBorder : 'transparent'}
                        bg={isActive ? activeItemBg : 'transparent'}
                        transition="all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275)"
                        zIndex="1"
                        className="menu-item-bg"
                        _groupHover={{
                            bg: hoverBg,
                            borderLeftColor: isActive ? activeItemBorder : hoverBg,
                            boxShadow: "lg",
                            transform: "translateY(-4px) scale(1.03)",
                        }}
                    />

                    <Flex
                        py={3}
                        px={4}
                        position="relative"
                        zIndex="2"
                        alignItems="center"
                        justifyContent={isCollapsed ? "center" : "flex-start"}
                        width="100%"
                        cursor="pointer"
                        onClick={isLogoutItem ? item?.logout : handleClick}
                        role="group"
                    >
                        <Box
                            position="relative"
                            zIndex="2"
                        >
                            <Icon
                                as={item.icon}
                                boxSize={5}
                                flexShrink={0}
                                color={isActive ? activeIconColor : iconColor}
                                _groupHover={{
                                    color: "white",
                                    transform: isCollapsed ? "scale(1.3)" : "scale(1.2)",
                                    filter: "drop-shadow(0 1px 2px rgba(0,0,0,0.2))"
                                }}
                                transition="all 0.3s ease"
                                transform={isCollapsed ? "scale(1.2)" : "scale(1)"}
                            />
                        </Box>

                        {!isCollapsed && (
                            <Box
                                position="relative"
                                zIndex="2"
                                ml={4}
                                width="calc(100% - 28px)"
                                overflow="hidden"
                            >
                                <Text
                                    className="menu-item-text"
                                    fontSize="sm"
                                    fontWeight={isActive ? "medium" : "normal"}
                                    color={textColor}
                                    transition="all 0.25s ease"
                                    opacity={isCollapsed ? 0 : 1}
                                    overflow="hidden"
                                    textOverflow="ellipsis"
                                    whiteSpace={width < 220 ? "normal" : "nowrap"}
                                    wordBreak={width < 220 ? "break-word" : "normal"}
                                    lineHeight="1.2"
                                    maxHeight={width < 220 ? "36px" : "none"}
                                    _groupHover={{
                                        color: "white",
                                        fontWeight: "bold",
                                        letterSpacing: "0.02em",
                                        transform: "scale(1.05)",
                                        textShadow: "0 1px 2px rgba(0,0,0,0.2)",
                                    }}
                                    style={{
                                        transformOrigin: "left center",
                                    }}
                                >
                                    {item.label}
                                </Text>
                            </Box>
                        )}

                        {!isCollapsed && (
                            <Box
                                position="absolute"
                                right="12px"
                                color="white"
                                opacity={isHovered ? 1 : 0}
                                transform={isHovered ? "translateX(0)" : "translateX(-10px)"}
                                fontWeight="bold"
                                fontSize="md"
                                zIndex="2"
                                transition="all 0.3s ease"
                                textShadow="0 1px 2px rgba(0,0,0,0.2)"
                            >
                                →
                            </Box>
                        )}
                    </Flex>
                </Box>
            </Tooltip>
        );
    };

    const handleToggleClick = (e) => {
        e.preventDefault();
        e.stopPropagation();
        toggleSidebar();
    };

    return (
        <Flex
            ref={sidebarRef}
            direction="column"
            h="full"
            borderRight="1px"
            borderColor={borderColor}
            bg={bgColor}
            w={isCollapsed ? `${collapsedWidth}px` : `${width}px`}
            maxW={isCollapsed ? `${collapsedWidth}px` : isMobile ? "80%" : "none"}
            position={isMobile ? "fixed" : "relative"}
            zIndex={isMobile ? "overlay" : "auto"}
            transition={isResizing ? "none" : "width 0.2s ease"}
            justify="space-between"
            overflow="hidden"
            shadow={isMobile ? "xl" : "none"}
        >
            <Box
                as="style"
                dangerouslySetInnerHTML={{
                    __html: `
            .menu-item-container:hover .menu-item-bg {
              background-color: var(--chakra-colors-cyan-500) !important;
              transform: translateY(-4px) scale(1.03) !important;
              box-shadow: var(--chakra-shadows-lg) !important;
            }
          `
                }}
            />

            <Box
                overflowY="auto"
                overflowX="hidden"
                py={2}
                pt={6}
                flex="1"
                css={{
                    '&::-webkit-scrollbar': {
                        width: '4px',
                    },
                    '&::-webkit-scrollbar-track': {
                        background: 'transparent',
                    },
                    '&::-webkit-scrollbar-thumb': {
                        background: 'var(--chakra-colors-gray-300)',
                        borderRadius: '4px',
                    },
                    '&::-webkit-scrollbar-thumb:hover': {
                        background: 'var(--chakra-colors-gray-400)',
                    },
                    msOverflowStyle: 'none',
                    scrollbarWidth: 'thin',
                }}
            >
                <VStack spacing={0} align="stretch" width="100%" maxWidth="100%">
                    {menuItems.map((item, index) => {
                        const isActive = location.pathname === item.path;
                        return <MenuItem key={index} item={item} isActive={isActive} />;
                    })}
                </VStack>
            </Box>

            <Box borderTop="1px" borderColor={borderColor} maxWidth="100%" overflow="hidden">
                {isCollapsed ? (
                    <VStack spacing={0} align="stretch" width="100%">
                        <MenuItem
                            item={logoutItem}
                            isActive={location.pathname === logoutItem.path}
                            isLogoutItem={true}
                        />

                        <Flex
                            justifyContent="center"
                            p={2}
                            borderTop="1px"
                            borderColor="gray.100"
                            opacity={isLogoutHovered ? 0 : 1}
                            transition="opacity 0.3s ease"
                        >
                            <Box>
                                <Tooltip
                                    label="Expand sidebar"
                                    placement="right"
                                    hasArrow
                                    bg="cyan.500"
                                    color="white"
                                    fontWeight="medium"
                                    px={3}
                                    py={2}
                                    borderRadius="md"
                                    openDelay={200}
                                    gutter={12}
                                >
                                    <IconButton
                                        size="sm"
                                        variant="ghost"
                                        icon={<FiChevronsRight />}
                                        onClick={handleToggleClick}
                                        onMouseDown={(e) => e.stopPropagation()}
                                        aria-label="Expand sidebar"
                                        color="gray.500"
                                        _hover={{
                                            color: "cyan.500",
                                            bg: "gray.50",
                                            transform: "scale(1.1)",
                                            boxShadow: "0 0 8px rgba(0, 188, 212, 0.3)",
                                        }}
                                        transition="all 0.3s ease"
                                        _active={{ transform: "scale(0.95)" }}
                                    />
                                </Tooltip>
                            </Box>
                        </Flex>
                    </VStack>
                ) : (
                    <Flex position="relative" alignItems="center" width="100%" maxWidth="100%">
                        <Box flex="1" maxWidth="100%">
                            <MenuItem
                                item={logoutItem}
                                isActive={location.pathname === logoutItem.path}
                                isLogoutItem={true}
                            />
                        </Box>

                        <Box
                            position="absolute"
                            right="2"
                            zIndex={2}
                            opacity={isLogoutHovered ? 0 : 1}
                            transition="opacity 0.3s ease"
                        >
                            <Tooltip
                                label="Collapse sidebar"
                                placement="left"
                                hasArrow
                                bg="cyan.500"
                                color="white"
                                fontWeight="medium"
                                px={3}
                                py={2}
                                borderRadius="md"
                                openDelay={200}
                                gutter={12}
                            >
                                <IconButton
                                    size="sm"
                                    variant="ghost"
                                    icon={<FiChevronsLeft />}
                                    onClick={handleToggleClick}
                                    onMouseDown={(e) => e.stopPropagation()}
                                    aria-label="Collapse sidebar"
                                    color="gray.500"
                                    _hover={{
                                        color: "cyan.500",
                                        bg: "gray.50",
                                        transform: "scale(1.1)",
                                        boxShadow: "0 0 8px rgba(0, 188, 212, 0.3)",
                                    }}
                                    _active={{ transform: "scale(0.95)" }}
                                    transition="all 0.3s ease"
                                />
                            </Tooltip>
                        </Box>
                    </Flex>
                )}
            </Box>

            {!isMobile && (
                <Box
                    position="absolute"
                    top="0"
                    right="0"
                    width="6px"
                    height="100%"
                    cursor="ew-resize"
                    zIndex="3"
                    onMouseDown={handleResizeStart}
                    _hover={{
                        _before: {
                            backgroundColor: isCollapsed ? 'transparent' : resizeHandleHoverColor,
                            transition: "background-color 0.2s",
                            width: "3px",
                            boxShadow: "0 0 6px rgba(0, 188, 212, 0.5)",
                        }
                    }}
                    _before={{
                        content: '""',
                        position: 'absolute',
                        top: 0,
                        bottom: 0,
                        left: '50%',
                        transform: 'translateX(-50%)',
                        width: '2px',
                        backgroundColor: isResizing ? resizeHandleHoverColor : 'transparent',
                        transition: 'all 0.3s ease',
                    }}
                />
            )}

            {isMobile && !isCollapsed && (
                <Box
                    position="fixed"
                    top="0"
                    left="0"
                    right="0"
                    bottom="0"
                    bg="blackAlpha.600"
                    zIndex="-1"
                    onClick={toggleSidebar}
                    transition="opacity 0.3s ease"
                />
            )}
        </Flex>
    );
};

export default SupplierSidebar;