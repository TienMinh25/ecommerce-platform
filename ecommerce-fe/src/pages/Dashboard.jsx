import React from 'react';
import {Box, useColorModeValue} from '@chakra-ui/react';
import {motion} from 'framer-motion';
import useAuth from "../hooks/useAuth.js";
import DashboardGreeting from "../components/dashboard/admin/DashboardGreeting.jsx";

// Create a motion-enabled version of Box
const MotionBox = motion(Box);

const Dashboard = () => {
  const {user} = useAuth()

  const bgGradient = useColorModeValue(
      'linear(to-br, blue.50, purple.50)',
      'linear(to-br, blue.900, purple.900)'
  );

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
        {/* Dashboard Greeting */}
        <MotionBox
            variants={itemVariants}
            px={6}
            py={4}
            mb={4}
        >
          <DashboardGreeting
              fullName={user.fullname}
              avatarUrl={user.avatarUrl}
          />
        </MotionBox>
      </MotionBox>
  );
};

export default Dashboard;