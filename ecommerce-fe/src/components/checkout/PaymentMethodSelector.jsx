import React, { useState, useEffect } from 'react';
import {
    Box,
    VStack,
    HStack,
    Text,
    Icon,
    Radio,
    RadioGroup,
    useToast,
    Skeleton,
    SkeletonText
} from '@chakra-ui/react';
import { FiCreditCard, FiDollarSign, FiSmartphone } from 'react-icons/fi';
import paymentMethodService from '../../services/paymentMethodService';

const PaymentMethodSelector = ({ selectedPaymentMethod, onPaymentMethodSelect }) => {
    const [paymentMethods, setPaymentMethods] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const toast = useToast();

    useEffect(() => {
        fetchPaymentMethods();
    }, []);

    const fetchPaymentMethods = async () => {
        try {
            setIsLoading(true);
            const response = await paymentMethodService.getPaymentMethods();
            setPaymentMethods(response.data);
        } catch (error) {
            console.error('Error fetching payment methods:', error);
            toast({
                title: 'Lỗi tải phương thức thanh toán',
                description: 'Không thể tải danh sách phương thức thanh toán',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    };

    const getPaymentIcon = (code) => {
        switch (code) {
            case 'cod':
                return FiDollarSign;
            case 'momo':
                return FiSmartphone;
            case 'bank':
                return FiCreditCard;
            case 'zalopay':
                return FiSmartphone;
            default:
                return FiCreditCard;
        }
    };

    const getPaymentIconColor = (code) => {
        switch (code) {
            case 'cod':
                return 'green.500';
            case 'momo':
                return 'pink.500';
            case 'bank':
                return 'blue.500';
            case 'zalopay':
                return 'blue.400';
            default:
                return 'gray.500';
        }
    };

    if (isLoading) {
        return (
            <Box bg="white" p={4} borderRadius="md" borderWidth="1px">
                <HStack mb={4}>
                    <Skeleton height="20px" width="20px" />
                    <Skeleton height="20px" width="200px" />
                </HStack>

                <VStack spacing={3} align="stretch">
                    {[1, 2].map(i => (
                        <Box
                            key={i}
                            p={3}
                            borderWidth="1px"
                            borderRadius="md"
                        >
                            <HStack spacing={3}>
                                <Skeleton height="16px" width="16px" />
                                <VStack align="start" spacing={1} flex="1">
                                    <Skeleton height="16px" width="60%" />
                                    <SkeletonText noOfLines={1} spacing="2" width="40%" />
                                </VStack>
                            </HStack>
                        </Box>
                    ))}
                </VStack>
            </Box>
        );
    }

    return (
        <Box bg="white" p={4} borderRadius="md" borderWidth="1px">
            <HStack mb={4}>
                <Icon as={FiCreditCard} color="green.500" />
                <Text fontWeight="semibold">Phương thức thanh toán</Text>
            </HStack>

            <RadioGroup value={selectedPaymentMethod} onChange={onPaymentMethodSelect}>
                <VStack spacing={3} align="stretch">
                    {paymentMethods.map((method) => (
                        <Box
                            key={method.id}
                            p={3}
                            borderWidth="1px"
                            borderColor={selectedPaymentMethod === method.code ? "green.300" : "gray.200"}
                            borderRadius="md"
                            bg={selectedPaymentMethod === method.code ? "green.50" : "white"}
                            opacity={1}
                            cursor={"pointer"}
                            onClick={() => onPaymentMethodSelect(method.code)}
                            _hover={{ borderColor: "green.200", bg: "green.25" }}
                        >
                            <Radio
                                value={method.code}
                                colorScheme="green"
                            >
                                <HStack spacing={3}>
                                    <Icon
                                        as={getPaymentIcon(method.code)}
                                        color={getPaymentIconColor(method.code)}
                                    />
                                    <VStack align="start" spacing={0}>
                                        <Text fontWeight="medium">
                                            {method.name}
                                        </Text>
                                        {method.description && (
                                            <Text fontSize="xs" color="gray.500">
                                                {method.description}
                                            </Text>
                                        )}
                                    </VStack>
                                </HStack>
                            </Radio>
                        </Box>
                    ))}
                </VStack>
            </RadioGroup>
        </Box>
    );
};

export default PaymentMethodSelector;