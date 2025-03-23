import { Box } from '@chakra-ui/react';
import { Outlet } from 'react-router-dom';
import Header from './Header';
import Footer from './Footer';

const MainLayout = () => {
  return (
    <Box minH='100vh' display='flex' flexDirection='column'>
      <Header />
      <Box flex='1' as='main'>
        <Outlet />
      </Box>
      <Footer />
    </Box>
  );
};

export default MainLayout;
