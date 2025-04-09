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
  useDisclosure,
} from '@chakra-ui/react';

import LogoCompact from '../../ui/LogoCompact.jsx';
import { useNavigate } from "react-router-dom";
import useAuth from "../../../hooks/useAuth.js";

const DashboardHeader = () => {
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const textColor = useColorModeValue('gray.800', 'gray.100');
  const navigate = useNavigate();
  const {user} = useAuth()

  // Use the useDisclosure hook for controlling the menu
  const { isOpen, onOpen, onClose } = useDisclosure();

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
              Minh Plaza
            </Text>
          </Flex>

          {/* User Profile */}
          <Menu isLazy placement="bottom-end" isOpen={isOpen}>
            <MenuButton
                as={Box}
                onMouseEnter={onOpen}
                onMouseLeave={onClose}
            >
              <HStack spacing={3} cursor="pointer">
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
                shadow="md"
                minW="180px"
                onMouseEnter={onOpen}
                onMouseLeave={onClose}
            >
              <MenuItem>Profile</MenuItem>
              <MenuItem>Settings</MenuItem>
              <MenuItem>Logout</MenuItem>
            </MenuList>
          </Menu>
        </Flex>
      </Box>
  );
};

export default DashboardHeader;