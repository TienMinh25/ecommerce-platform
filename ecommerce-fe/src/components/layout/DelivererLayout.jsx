import React, {createContext, useContext, useEffect, useMemo, useState} from 'react';
import {Box, Flex, useMediaQuery} from '@chakra-ui/react';
import {Navigate, Outlet, useLocation} from 'react-router-dom';
import useAuth from "../../hooks/useAuth.js";
import DelivererSidebar from "../dashboard/deliverer/DelivererSidebar.jsx";
import DelivererHeader from "../dashboard/deliverer/DelivererHeader.jsx";

// Define the header height and spacing
const headerHeight = "64px";
const defaultSidebarWidth = "256px";
const collapsedSidebarWidth = "64px";
const footerHeight = "64px";

// Create a context to share sidebar state across components
export const DelivererSidebarContext = createContext({
    isCollapsed: false,
    width: defaultSidebarWidth,
    collapsedWidth: collapsedSidebarWidth,
    updateSidebarState: () => {},
    isMobile: false
});

// Custom hook to use the sidebar context
export const useDelivererSidebar = () => useContext(DelivererSidebarContext);

const DelivererLayout = () => {
    const location = useLocation();
    const { user } = useAuth(); // Get user from auth context

    // Check if user has deliverer role
    const isDeliverer = user.hasRole('deliverer');

    // If user is not deliverer, redirect to home page
    if (!isDeliverer) {
        return <Navigate to="/" replace />;
    }

    // Check if the viewport is mobile-sized
    const [isMobileResult] = useMediaQuery("(max-width: 768px)");
    // Stabilize isMobile value with useEffect to prevent constant updates
    const [isMobile, setIsMobile] = useState(false);

    useEffect(() => {
        setIsMobile(isMobileResult);
    }, [isMobileResult]);

    // State to track sidebar configuration - initialize with default values
    const [sidebarState, setSidebarState] = useState({
        isCollapsed: false,
        width: defaultSidebarWidth,
        collapsedWidth: collapsedSidebarWidth,
        toggleableOnlyInDashboard: true
    });

    // Set initial collapsed state based on mobile detection - only once after mobile state is stable
    useEffect(() => {
        if (isMobile !== undefined) {
            setSidebarState(prev => ({
                ...prev,
                isCollapsed: isMobile
            }));
        }
    }, [isMobile]);

    // Function to update sidebar state, using functional update to avoid stale state issues
    const updateSidebarState = useMemo(() => {
        return (newState) => {
            setSidebarState(prev => {
                // Only apply isCollapsed change if toggleable is true
                if ('isCollapsed' in newState && !prev.toggleableOnlyInDashboard) {
                    const { isCollapsed, ...rest } = newState;
                    return { ...prev, ...rest };
                }
                return { ...prev, ...newState };
            });
        };
    }, []);

    // Calculate current sidebar width based on state - memoize to prevent recalculation
    const formattedSidebarWidth = useMemo(() => {
        const currentWidth = sidebarState.isCollapsed
            ? sidebarState.collapsedWidth
            : sidebarState.width;

        return typeof currentWidth === 'string'
            ? currentWidth
            : `${currentWidth}px`;
    }, [sidebarState.isCollapsed, sidebarState.width, sidebarState.collapsedWidth]);

    // Memoize context value to prevent unnecessary re-renders
    const contextValue = useMemo(() => ({
        ...sidebarState,
        updateSidebarState,
        isMobile,
        width: typeof sidebarState.width === 'string' ? sidebarState.width : `${sidebarState.width}px`,
        collapsedWidth: typeof sidebarState.collapsedWidth === 'string' ? sidebarState.collapsedWidth : `${sidebarState.collapsedWidth}px`
    }), [sidebarState, updateSidebarState, isMobile]);

    return (
        <DelivererSidebarContext.Provider value={contextValue}>
            <Box minH="100vh" w="100%" overflow="hidden" position="relative">
                {/* Header - Sticky at the top, full width */}
                <Box
                    position="fixed"
                    top="0"
                    left="0"
                    right="0"
                    height={headerHeight}
                    zIndex="30"
                    boxShadow="md"
                    bg="white"
                    width="100%"
                >
                    <DelivererHeader />
                </Box>

                {/* Deliverer Content */}
                <Flex height="100vh" width="100%" pt={headerHeight}>
                    {/* Sidebar - Fixed position with responsive behavior */}
                    <Box
                        position="fixed"
                        top={headerHeight}
                        left="0"
                        height={`calc(100vh - ${headerHeight})`}
                        width={isMobile ? "100%" : formattedSidebarWidth}
                        zIndex={isMobile ? "40" : "20"}
                        boxShadow="md"
                        transition="width 0.25s ease, transform 0.25s ease"
                        bg="white"
                        overflowX="hidden"
                        overflowY="auto"
                        transform={isMobile && sidebarState.isCollapsed ? "translateX(-100%)" : "translateX(0)"}
                        visibility={isMobile && sidebarState.isCollapsed ? "hidden" : "visible"}
                    >
                        <DelivererSidebar
                            onStateChange={updateSidebarState}
                            canToggle={true}
                        />
                    </Box>

                    {/* Overlay for mobile when sidebar is open */}
                    {isMobile && !sidebarState.isCollapsed && (
                        <Box
                            position="fixed"
                            top={headerHeight}
                            left="0"
                            right="0"
                            bottom="0"
                            bg="blackAlpha.600"
                            zIndex="25"
                            onClick={() => updateSidebarState({ isCollapsed: true })}
                        />
                    )}

                    {/* Main Content Area - Takes full remaining width */}
                    <Box
                        position="absolute"
                        top={headerHeight}
                        left={isMobile ? "0" : formattedSidebarWidth}
                        right="0"
                        bottom={footerHeight}
                        transition="left 0.25s ease"
                        bg="gray.50"
                        overflowY="auto"
                        overflowX="hidden"
                        zIndex="10"
                    >
                        {/* Content Container */}
                        <Box
                            w="100%"
                            h="100%"
                            p={{ base: "3", md: "6" }}
                            pb={{ base: "6", md: "6" }}
                        >
                            <Outlet />
                        </Box>
                    </Box>
                </Flex>
            </Box>
        </DelivererSidebarContext.Provider>
    );
};

export default DelivererLayout;