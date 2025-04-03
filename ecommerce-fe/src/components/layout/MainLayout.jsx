import { Box, Flex } from '@chakra-ui/react';
import { Outlet } from 'react-router-dom';
import Header from './Header';
import Footer from './Footer';

const MainLayout = () => {
    return (
        <Box
            minH='100vh'
            display='flex'
            flexDirection='column'
            overflow="hidden"
            position="relative"
        >
            <Header />
            <Box
                flex='1'
                as='main'
                overflow="auto"
                width="100%"
                position="relative"
            >
                <Outlet />
            </Box>

            {/* Elegant Divider */}
            <Flex
                width="100%"
                justifyContent="center"
                alignItems="center"
                my={6}
                px={4}
            >
                <Box
                    height="2px"
                    width="90%"
                    bgGradient="linear(to-r, transparent, brand.500, transparent)"
                    opacity={0.5}
                    boxShadow="sm"
                />
            </Flex>

            <Footer
                position="relative"
                width="100%"
                zIndex="1"
            />
        </Box>
    );
};

export default MainLayout;