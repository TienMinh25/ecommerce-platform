import React, { useEffect } from 'react';
import { Box, useColorModeValue } from '@chakra-ui/react';
import { useNavigate, useLocation } from 'react-router-dom';
import { motion } from 'framer-motion';
import UserManagementComponent from './Module/Dashboard/UserManagementComponent';

// Create a motion-enabled version of Box
const MotionBox = motion(Box);

const Dashboard = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const bgGradient = useColorModeValue(
    'linear(to-br, blue.50, purple.50)',
    'linear(to-br, blue.900, purple.900)'
  );
  
  useEffect(() => {
    // Only redirect if we're exactly at /dashboard
    if (location.pathname === '/dashboard') {
      navigate('/dashboard/users', { replace: true });
    }
  }, [navigate, location]);

  // Animation variants for staggered children
  const containerVariants = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.2
      }
    }
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    show: { 
      y: 0, 
      opacity: 1,
      transition: { 
        type: "spring",
        stiffness: 300,
        damping: 24
      }
    }
  };

  return (
    <MotionBox
      as="section"
      width="100%"
      height="100%" 
      position="relative"
      overflow="auto"
      bgGradient={bgGradient}
      initial="hidden"
      animate="show"
      variants={containerVariants}
    >
      {/* Header animation */}
      <MotionBox
        variants={itemVariants}
        px={6}
        py={4}
        mb={4}
        borderRadius="md"
        bg={useColorModeValue('white', 'gray.800')}
        boxShadow="sm"
        mx={4}
        mt={4}
      >
        <motion.h1 
          style={{ 
            fontSize: '24px', 
            fontWeight: 'bold',
            margin: 0
          }}
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.3, duration: 0.5 }}
        >
          User Management
        </motion.h1>
      </MotionBox>
      
      {/* Main content animation */}
      <MotionBox
        variants={itemVariants}
        flex="1"
        m={4}
        p={6}
        borderRadius="md"
        boxShadow="md"
        bg={useColorModeValue('white', 'gray.800')}
      >
        <UserManagementComponent />
      </MotionBox>
    </MotionBox>
  );
};

export default Dashboard;