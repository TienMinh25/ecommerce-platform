import React from 'react';
import { Box } from '@chakra-ui/react';
import DashboardComponent from './Module/DashboardComponent';

const Dashboard = () => {
  return (
    <Box 
      width="100%" 
      height="100%" 
      position="relative"
      overflow="auto" // Changed from "hidden" to "auto"
    >
      <DashboardComponent />
    </Box>
  );
};

export default Dashboard;