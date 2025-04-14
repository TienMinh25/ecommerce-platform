import React, { useState } from 'react';
import {
    Box,
    Button,
    FormControl,
    Input,
    Avatar,
    Text,
    useToast,
    Heading,
    Divider,
    Flex,
    Table,
    Tbody,
    Tr,
    Td,
    Radio,
    RadioGroup,
    Stack,
    Link,
} from '@chakra-ui/react';
import { useAuth } from '../../hooks/useAuth';
import { BirthDateSelector } from './Date'; // Import our custom component

const UserProfile = () => {
    const { user } = useAuth();
    const toast = useToast();

    // Parse birth date if available
    const parseBirthDate = () => {
        if (!user?.birth_date) return { day: null, month: null, year: null };

        try {
            const date = new Date(user.birth_date);
            return {
                day: date.getDate(),
                month: date.getMonth() + 1, // JavaScript months are 0-indexed
                year: date.getFullYear()
            };
        } catch (e) {
            return { day: null, month: null, year: null };
        }
    };

    // Form state
    const [formData, setFormData] = useState({
        fullname: user?.fullname || '',
        email: user?.email || '',
        phone: user?.phone || '',
        gender: user?.gender || 'Nam',
        ...parseBirthDate()
    });

    const [isSubmitting, setIsSubmitting] = useState(false);
    const [selectedFile, setSelectedFile] = useState(null);
    const [previewUrl, setPreviewUrl] = useState(user?.avatarUrl || '');

    // Handle input changes
    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    // Handle date changes
    const handleDateChange = (field, value) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    // Handle gender change
    const handleGenderChange = (value) => {
        setFormData(prev => ({ ...prev, gender: value }));
    };

    // Handle avatar file selection
    const handleFileSelect = (e) => {
        const file = e.target.files[0];
        if (file) {
            setSelectedFile(file);

            // Create preview URL
            const reader = new FileReader();
            reader.onloadend = () => {
                setPreviewUrl(reader.result);
            };
            reader.readAsDataURL(file);
        }
    };

    // Handle form submission
    const handleSubmit = async (e) => {
        e.preventDefault();

        // Create a proper birth date string if all date components are present
        let birthDateString = null;
        if (formData.day && formData.month && formData.year) {
            birthDateString = `${formData.year}-${formData.month.toString().padStart(2, '0')}-${formData.day.toString().padStart(2, '0')}`;
        }

        setIsSubmitting(true);

        try {
            // Data to send to the server
            const userData = {
                fullname: formData.fullname,
                birth_date: birthDateString,
                gender: formData.gender,
                // Include other fields as needed
            };

            // Here you would make an API call to update the user profile
            // For now, we'll just simulate it with a timeout
            await new Promise(resolve => setTimeout(resolve, 1000));

            toast({
                title: 'Hồ sơ đã được cập nhật',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
        } catch (error) {
            toast({
                title: 'Lỗi',
                description: error.message || 'Không thể cập nhật hồ sơ',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <Box
            as="form"
            onSubmit={handleSubmit}
            maxH="calc(100vh - 200px)"  // Set a max height to prevent excessive page length
            overflowY="auto"  // Enable scrolling if content exceeds max height
            pr={2}  // Add a bit of padding for the scrollbar
            sx={{
                '&::-webkit-scrollbar': {
                    width: '4px',
                },
                '&::-webkit-scrollbar-thumb': {
                    backgroundColor: 'rgba(0,0,0,0.2)',
                    borderRadius: '2px',
                }
            }}
        >
            <Heading as="h1" size="lg" mb={4}>Hồ Sơ Của Tôi</Heading>
            <Text mb={4} color="gray.500" fontSize="sm">
                Quản lý thông tin hồ sơ để bảo mật tài khoản
            </Text>

            <Divider mb={6} />

            <Flex
                direction={{ base: 'column', md: 'row' }}
                gap={{ base: 6, md: 8 }}
                align="flex-start"
            >
                {/* Left side - Form */}
                <Box
                    flex="1"
                    borderRight={{ md: '1px' }}
                    borderColor={{ md: 'gray.200' }}
                    pr={{ md: 8 }}
                >
                    <Table variant="simple" size="md">
                        <Tbody>
                            <Tr>
                                <Td width="180px" pr={2} pl={0} py={4} verticalAlign="top">
                                    <Text fontWeight="medium">Tên đăng nhập</Text>
                                </Td>
                                <Td py={4} pl={4}>
                                    <Text>{formData.username || user?.username}</Text>
                                </Td>
                            </Tr>

                            <Tr>
                                <Td width="180px" pr={2} pl={0} py={4} verticalAlign="top">
                                    <Text fontWeight="medium">Tên</Text>
                                </Td>
                                <Td py={4} pl={4}>
                                    <Input
                                        name="fullname"
                                        value={formData.fullname}
                                        onChange={handleChange}
                                        size="md"
                                        width="100%"
                                        maxW="400px"
                                    />
                                </Td>
                            </Tr>

                            <Tr>
                                <Td width="180px" pr={2} pl={0} py={4} verticalAlign="top">
                                    <Text fontWeight="medium">Email</Text>
                                </Td>
                                <Td py={4} pl={4}>
                                    <Flex align="center">
                                        <Text>{formData.email}</Text>
                                        <Link color="blue.500" fontSize="sm" ml={3}>
                                            Thay Đổi
                                        </Link>
                                    </Flex>
                                </Td>
                            </Tr>

                            <Tr>
                                <Td width="180px" pr={2} pl={0} py={4} verticalAlign="top">
                                    <Text fontWeight="medium">Số điện thoại</Text>
                                </Td>
                                <Td py={4} pl={4}>
                                    {formData.phone ? (
                                        <Flex align="center">
                                            <Text>{formData.phone}</Text>
                                            <Link color="blue.500" fontSize="sm" ml={3}>
                                                Thay Đổi
                                            </Link>
                                        </Flex>
                                    ) : (
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            colorScheme="blue"
                                            borderRadius="sm"
                                        >
                                            Thêm
                                        </Button>
                                    )}
                                </Td>
                            </Tr>

                            <Tr>
                                <Td width="180px" pr={2} pl={0} py={4} verticalAlign="top">
                                    <Text fontWeight="medium">Giới tính</Text>
                                </Td>
                                <Td py={4} pl={4}>
                                    <RadioGroup onChange={handleGenderChange} value={formData.gender}>
                                        <Stack direction="row" spacing={6}>
                                            <Radio value="Nam">Nam</Radio>
                                            <Radio value="Nữ">Nữ</Radio>
                                            <Radio value="Khác">Khác</Radio>
                                        </Stack>
                                    </RadioGroup>
                                </Td>
                            </Tr>

                            <Tr>
                                <Td width="180px" pr={2} pl={0} py={4} verticalAlign="top">
                                    <Text fontWeight="medium">Ngày sinh</Text>
                                </Td>
                                <Td py={4} pl={4}>
                                    <Box maxW="400px">
                                        <BirthDateSelector
                                            selectedDay={formData.day}
                                            selectedMonth={formData.month}
                                            selectedYear={formData.year}
                                            onChange={handleDateChange}
                                        />
                                    </Box>
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
                                        borderRadius="sm"
                                    >
                                        Lưu
                                    </Button>
                                </Td>
                            </Tr>
                        </Tbody>
                    </Table>
                </Box>

                {/* Right side - Avatar */}
                <Box
                    width={{ base: 'full', md: '200px' }}
                    pt={{ md: 6 }}
                    display="flex"
                    flexDirection="column"
                    alignItems="center"
                >
                    <Avatar
                        size="2xl"
                        src={previewUrl || user?.avatarUrl}
                        name={formData.fullname}
                        border="1px solid"
                        borderColor="gray.200"
                        mb={4}
                    />

                    <Button
                        as="label"
                        htmlFor="avatar-upload"
                        variant="outline"
                        cursor="pointer"
                        size="sm"
                        colorScheme="blue"
                        borderRadius="sm"
                        mb={3}
                    >
                        Chọn Ảnh
                        <input
                            id="avatar-upload"
                            type="file"
                            accept="image/*"
                            onChange={handleFileSelect}
                            style={{ display: 'none' }}
                        />
                    </Button>

                    <Text fontSize="xs" color="gray.500" textAlign="center">
                        Dung lượng file tối đa 1 MB
                        <br />
                        Định dạng: .JPEG, .PNG
                    </Text>
                </Box>
            </Flex>
        </Box>
    );
};

export default UserProfile;