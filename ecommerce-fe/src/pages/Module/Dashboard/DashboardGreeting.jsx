import React, { useMemo } from 'react';
import {
    Box,
    Flex,
    Text,
    Avatar,
    VStack,
    HStack,
    Icon,
    useColorModeValue
} from '@chakra-ui/react';
import { motion } from 'framer-motion';
import {
    FaSun,
    FaMoon,
    FaCloud,
    FaTasks
} from 'react-icons/fa';

const MotionBox = motion(Box);

const DashboardGreeting = ({ fullName, avatarUrl }) => {
    // Get current time and generate appropriate greeting
    const { greeting, timeIcon, motivationalMessage } = useMemo(() => {
        const currentHour = new Date().getHours();
        let greetingText = 'Chào buổi tối';
        let timeIcon = FaMoon;
        let motivationalMessage = 'Hãy kiểm tra các nhiệm vụ quan trọng của bạn';

        if (currentHour >= 5 && currentHour < 12) {
            greetingText = 'Chào buổi sáng';
            timeIcon = FaSun;
            motivationalMessage = 'Bắt đầu ngày mới với năng lượng tích cực';
        } else if (currentHour >= 12 && currentHour < 18) {
            greetingText = 'Chào buổi chiều';
            timeIcon = FaCloud;
            motivationalMessage = 'Hãy tiếp tục giữ nhịp độ và đạt được mục tiêu hôm nay';
        }

        return { greeting: greetingText, timeIcon, motivationalMessage };
    }, []);

    const gradientBg = useColorModeValue(
        'linear(to-r, blue.500, purple.500)',
        'linear(to-r, blue.300, purple.300)'
    );
    const textColor = useColorModeValue('white', 'gray.100');

    return (
        <MotionBox
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{
                duration: 0.6,
                ease: "easeOut"
            }}
            bgGradient={gradientBg}
            borderRadius="xl"
            p={6}
            boxShadow="xl"
        >
            <Flex
                alignItems="center"
                justifyContent="space-between"
            >
                <VStack
                    align="start"
                    spacing={3}
                    color={textColor}
                    flex={1}
                    mr={4}
                >
                    <HStack spacing={3} alignItems="center">
                        <Icon as={timeIcon} boxSize={6} />
                        <Text
                            fontSize={{ base: 'xl', md: '2xl' }}
                            fontWeight="bold"
                        >
                            {greeting}
                        </Text>
                    </HStack>
                    <Text
                        fontSize={{ base: '2xl', md: '3xl' }}
                        fontWeight="extrabold"
                        textShadow="1px 1px 2px rgba(0,0,0,0.2)"
                    >
                        {fullName}
                    </Text>
                    <HStack spacing={3} alignItems="center">
                        <Icon as={FaTasks} boxSize={5} opacity={0.8} />
                        <Text
                            fontSize="md"
                            opacity={0.9}
                            fontStyle="italic"
                        >
                            {motivationalMessage}
                        </Text>
                    </HStack>
                </VStack>

                <Avatar
                    size={{ base: 'lg', md: 'xl' }}
                    name={fullName}
                    src={avatarUrl}
                    border="3px solid white"
                    boxShadow="md"
                />
            </Flex>
        </MotionBox>
    );
};

export default DashboardGreeting;