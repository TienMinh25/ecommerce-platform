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
  FiUsers,
  FiShield,
  FiKey,
  FiDatabase,
  FiPackage,
  FiTruck,
  FiLogOut,
  FiChevronsLeft,
  FiChevronsRight,
  FiLayers,
  FiMapPin,
} from 'react-icons/fi';
import { useNavigate, useLocation } from 'react-router-dom';
import { useSidebar } from '../../layout/DashboardLayout.jsx'; // Update with the correct path

const DashboardSidebar = ({ onStateChange }) => {
  const { updateSidebarState } = useSidebar();

  const [isCollapsed, setIsCollapsed] = useState(false);
  const [width, setWidth] = useState(256); // Default width (64 in Chakra units = 256px)
  const [isResizing, setIsResizing] = useState(false);
  const [isLogoutHovered, setIsLogoutHovered] = useState(false);
  const minWidth = 180; // Increased minimum width for better text display
  const maxWidth = 400; // Maximum width when dragging
  const collapsedWidth = 64; // Width when collapsed (16 in Chakra units)

  const sidebarRef = useRef(null);

  const navigate = useNavigate();
  const location = useLocation();

  // Responsive behavior
  const isMobile = useBreakpointValue({ base: true, md: false });
  const defaultCollapsed = useBreakpointValue({ base: true, md: false });

  // Initialize collapse state based on screen size
  useEffect(() => {
    setIsCollapsed(defaultCollapsed);
  }, [defaultCollapsed]);

  // Communicate state changes to parent component
  useEffect(() => {
    // Use the prop callback if provided (for backward compatibility)
    if (onStateChange) {
      onStateChange({
        isCollapsed,
        width,
        collapsedWidth
      });
    }

    // Use the context method if available
    if (updateSidebarState) {
      updateSidebarState({
        isCollapsed,
        width,
        collapsedWidth
      });
    }
  }, [isCollapsed, width, collapsedWidth, onStateChange, updateSidebarState]);

  // Menu configuration - Updated with new menu items
  const menuItems = [
    { icon: FiUsers, label: 'User Management', path: '/dashboard/users' },
    { icon: FiShield, label: 'Role Management', path: '/dashboard/roles' },
    { icon: FiKey, label: 'Permission Management', path: '/dashboard/permissions' },
    { icon: FiLayers, label: 'Module Management', path: '/dashboard/modules' },
    { icon: FiMapPin, label: 'Address Types Management', path: '/dashboard/address-types' },
    { icon: FiPackage, label: 'Onboarding Supplier Management', path: '/dashboard/onboarding/suppliers' },
    { icon: FiTruck, label: 'Onboarding Deliverer Management', path: '/dashboard/onboarding/deliverers' },
  ];

  // Logout item separate (for positioning at bottom)
  const logoutItem = { icon: FiLogOut, label: 'Logout', path: '/logout' };

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

  // Toggle sidebar collapsed state with animation
  const toggleSidebar = (e) => {
    if (e) {
      e.preventDefault();
      e.stopPropagation();
    }

    setIsCollapsed(prevState => !prevState);

    // Add spring effect when toggling
    if (sidebarRef.current) {
      sidebarRef.current.style.transition = "width 0.25s cubic-bezier(0.34, 1.56, 0.64, 1)";

      // Reset transition after animation completes
      setTimeout(() => {
        if (sidebarRef.current) {
          sidebarRef.current.style.transition = "width 0.2s ease";
        }
      }, 250);
    }
  };

  // Handle resize functionality
  const handleResizeStart = (e) => {
    e.preventDefault();
    setIsResizing(true);

    // If sidebar is collapsed, expand it when starting to drag
    if (isCollapsed) {
      setIsCollapsed(false);
      // Set initial width when opening from collapsed state
      setWidth(minWidth);
    }

    const startX = e.clientX;
    const startWidth = isCollapsed ? minWidth : width; // Use minWidth if collapsed

    const handleResize = (e) => {
      // Limit the width between minWidth and maxWidth
      const newWidth = Math.min(
          Math.max(startWidth + e.clientX - startX, minWidth),
          maxWidth
      );
      setWidth(newWidth);

      // Update text element states when resizing
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

  // Menu item component with enhanced hover transitions
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
            {/* Regular background fill - changes color on hover */}
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

            {/* Content container */}
            <Flex
                py={3}
                px={4}
                position="relative"
                zIndex="2"
                alignItems="center"
                justifyContent={isCollapsed ? "center" : "flex-start"}
                width="100%"
                cursor="pointer"
                onClick={handleClick}
                role="group"
            >
              {/* Icon container */}
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

              {/* Menu item text */}
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

              {/* Right arrow indicator with animation */}
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

  // Handle toggle button click - separate from navigation
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
          overflow="hidden" // Prevents overflow in the sidebar itself
          shadow={isMobile ? "xl" : "none"}
      >
        {/* CSS for menu items */}
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

        {/* Main menu items */}
        <Box
            overflowY="auto"
            overflowX="hidden" // Prevents horizontal scrolling in menu items area
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
              msOverflowStyle: 'none',  // IE and Edge
              scrollbarWidth: 'thin',   // Firefox
            }}
        >
          <VStack spacing={0} align="stretch" width="100%" maxWidth="100%">
            {menuItems.map((item, index) => {
              const isActive = location.pathname === item.path;
              return <MenuItem key={index} item={item} isActive={isActive} />;
            })}
          </VStack>
        </Box>

        {/* Bottom section with logout and collapse button */}
        <Box borderTop="1px" borderColor={borderColor} maxWidth="100%" overflow="hidden">
          {isCollapsed ? (
              // When collapsed, show collapse button in its own row
              <VStack spacing={0} align="stretch" width="100%">
                {/* Logout button */}
                <MenuItem
                    item={logoutItem}
                    isActive={location.pathname === logoutItem.path}
                    isLogoutItem={true}
                />

                {/* Collapse button in separate row below logout when collapsed */}
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
              // When expanded, show logout and collapse buttons in same row
              <Flex position="relative" alignItems="center" width="100%" maxWidth="100%">
                {/* Logout button */}
                <Box flex="1" maxWidth="100%">
                  <MenuItem
                      item={logoutItem}
                      isActive={location.pathname === logoutItem.path}
                      isLogoutItem={true}
                  />
                </Box>

                {/* Collapse button on same row as logout when expanded */}
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

        {/* Resize handle - still displayed when collapsed but visually hidden */}
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

        {/* Mobile overlay when sidebar is open */}
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

export default DashboardSidebar;