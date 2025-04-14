import React, { useState } from 'react';
import {
    Box,
    Button,
    FormControl,
    Input,
    InputGroup,
    InputRightElement,
    useToast,
    Heading,
    Divider,
    Text,
    Table,
    Tbody,
    Tr,
    Td
} from '@chakra-ui/react';
import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons';
import { useNavigate } from 'react-router-dom';

const ChangePassword = () => {
    const toast = useToast();
    const navigate = useNavigate();

    // Form state
    const [formData, setFormData] = useState({
        old_password: '',
        new_password: '',
        confirm_password: ''
    });

    // Form validation state
    const [errors, setErrors] = useState({});

    // Password visibility state
    const [showOldPassword, setShowOldPassword] = useState(false);
    const [showNewPassword, setShowNewPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    // Loading state
    const [isSubmitting, setIsSubmitting] = useState(false);

    // Handle input changes
    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));

        // Clear errors for the field being edited
        if (errors[name]) {
            setErrors(prev => ({ ...prev, [name]: null }));
        }
    };

    // Validate form
    const validateForm = () => {
        const newErrors = {};

        if (!formData.old_password) {
            newErrors.old_password = 'Vui lòng nhập mật khẩu hiện tại';
        }

        if (!formData.new_password) {
            newErrors.new_password = 'Vui lòng nhập mật khẩu mới';
        } else if (formData.new_password.length < 6) {
            newErrors.new_password = 'Mật khẩu phải có ít nhất 6 ký tự';
        } else if (formData.new_password === formData.old_password) {
            newErrors.new_password = 'Mật khẩu mới không được trùng với mật khẩu cũ';
        }

        if (!formData.confirm_password) {
            newErrors.confirm_password = 'Vui lòng xác nhận mật khẩu mới';
        } else if (formData.confirm_password !== formData.new_password) {
            newErrors.confirm_password = 'Mật khẩu xác nhận không khớp';
        }

        return newErrors;
    };

    // Handle form submission
    const handleSubmit = async (e) => {
        e.preventDefault();

        // Validate form
        const formErrors = validateForm();
        if (Object.keys(formErrors).length > 0) {
            setErrors(formErrors);
            return;
        }

        setIsSubmitting(true);

        try {
            // Here you would make an API call to change the password
            // For now, we'll just simulate it with a timeout
            await new Promise(resolve => setTimeout(resolve, 1000));

            toast({
                title: 'Đổi mật khẩu thành công',
                description: 'Mật khẩu của bạn đã được cập nhật',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });

            // Reset form
            setFormData({
                old_password: '',
                new_password: '',
                confirm_password: ''
            });

        } catch (error) {
            toast({
                title: 'Lỗi',
                description: error.message || 'Không thể đổi mật khẩu',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <Box as="form" onSubmit={handleSubmit}>
            <Heading as="h1" size="lg" mb={4}>Đổi Mật Khẩu</Heading>
            <Text mb={4} color="gray.500" fontSize="sm">
                Để bảo mật tài khoản, vui lòng không chia sẻ mật khẩu cho người khác
            </Text>

            <Divider mb={6} />

            <Box maxW="600px">
                <Table variant="simple" size="md">
                    <Tbody>
                        <Tr>
                            <Td width="180px" pr={2} pl={0} py={4} verticalAlign="top">
                                <Text fontWeight="medium">Mật khẩu hiện tại</Text>
                            </Td>
                            <Td py={4} pl={4}>
                                <FormControl isInvalid={!!errors.old_password}>
                                    <InputGroup>
                                        <Input
                                            name="old_password"
                                            type={showOldPassword ? 'text' : 'password'}
                                            value={formData.old_password}
                                            onChange={handleChange}
                                            placeholder="Nhập mật khẩu hiện tại"
                                            maxW="400px"
                                        />
                                        <InputRightElement width="3rem">
                                            <Button
                                                h="1.5rem"
                                                size="sm"
                                                variant="ghost"
                                                onClick={() => setShowOldPassword(!showOldPassword)}
                                            >
                                                {showOldPassword ? <ViewOffIcon /> : <ViewIcon />}
                                            </Button>
                                        </InputRightElement>
                                    </InputGroup>
                                    {errors.old_password && (
                                        <Text color="red.500" fontSize="sm" mt={1}>
                                            {errors.old_password}
                                        </Text>
                                    )}
                                </FormControl>
                            </Td>
                        </Tr>

                        <Tr>
                            <Td width="180px" pr={2} pl={0} py={4} verticalAlign="top">
                                <Text fontWeight="medium">Mật khẩu mới</Text>
                            </Td>
                            <Td py={4} pl={4}>
                                <FormControl isInvalid={!!errors.new_password}>
                                    <InputGroup>
                                        <Input
                                            name="new_password"
                                            type={showNewPassword ? 'text' : 'password'}
                                            value={formData.new_password}
                                            onChange={handleChange}
                                            placeholder="Nhập mật khẩu mới"
                                            maxW="400px"
                                        />
                                        <InputRightElement width="3rem">
                                            <Button
                                                h="1.5rem"
                                                size="sm"
                                                variant="ghost"
                                                onClick={() => setShowNewPassword(!showNewPassword)}
                                            >
                                                {showNewPassword ? <ViewOffIcon /> : <ViewIcon />}
                                            </Button>
                                        </InputRightElement>
                                    </InputGroup>
                                    {errors.new_password && (
                                        <Text color="red.500" fontSize="sm" mt={1}>
                                            {errors.new_password}
                                        </Text>
                                    )}
                                </FormControl>
                            </Td>
                        </Tr>

                        <Tr>
                            <Td width="180px" pr={2} pl={0} py={4} verticalAlign="top">
                                <Text fontWeight="medium">Xác nhận mật khẩu</Text>
                            </Td>
                            <Td py={4} pl={4}>
                                <FormControl isInvalid={!!errors.confirm_password}>
                                    <InputGroup>
                                        <Input
                                            name="confirm_password"
                                            type={showConfirmPassword ? 'text' : 'password'}
                                            value={formData.confirm_password}
                                            onChange={handleChange}
                                            placeholder="Xác nhận mật khẩu mới"
                                            maxW="400px"
                                        />
                                        <InputRightElement width="3rem">
                                            <Button
                                                h="1.5rem"
                                                size="sm"
                                                variant="ghost"
                                                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                            >
                                                {showConfirmPassword ? <ViewOffIcon /> : <ViewIcon />}
                                            </Button>
                                        </InputRightElement>
                                    </InputGroup>
                                    {errors.confirm_password && (
                                        <Text color="red.500" fontSize="sm" mt={1}>
                                            {errors.confirm_password}
                                        </Text>
                                    )}
                                </FormControl>
                            </Td>
                        </Tr>

                        <Tr>
                            <Td width="180px" pr={2} pl={0} py={6} verticalAlign="top">
                            </Td>
                            <Td py={6} pl={4}>
                                <Button
                                    colorScheme="red"
                                    type="submit"
                                    isLoading={isSubmitting}
                                    size="md"
                                >
                                    Xác nhận
                                </Button>
                            </Td>
                        </Tr>
                    </Tbody>
                </Table>
            </Box>
        </Box>
    );
};

export default ChangePassword;