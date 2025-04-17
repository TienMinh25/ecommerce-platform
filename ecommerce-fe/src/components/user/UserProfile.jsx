import React, {useEffect, useState} from 'react';
import {
    Box,
    Button,
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
} from '@chakra-ui/react';
import { BirthDateSelector } from './Date';
import userMeService from "../../services/userMeService.js"; // Import our custom component

const UserProfile = () => {
    const toast = useToast();

    // Supported image types
    const SUPPORTED_IMAGE_TYPES = [
        'image/jpeg',
        'image/png',
        'image/gif',
        'image/webp',
        'image/bmp',
        'image/tiff',
        'image/svg+xml'
    ];

    // Display names for supported formats
    const SUPPORTED_FORMATS_DISPLAY = '.JPEG, .PNG, .GIF, .WEBP, .BMP, .TIFF, .SVG';

    // Form state
    const [formData, setFormData] = useState({
        fullname: '',
        email: '',
        phone: '',
        day: null,
        month: null,
        year: null,
    });

    const [isSubmitting, setIsSubmitting] = useState(false);
    const [selectedFile, setSelectedFile] = useState(null);
    const [previewUrl, setPreviewUrl] = useState('');
    const [user, setUser] = useState(null);

    // Fetch user profile on mount
    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const profileData = await userMeService.getProfile();
                setUser(profileData);
                setFormData({
                    fullname: profileData.full_name || '',
                    email: profileData.email || '',
                    phone: profileData.phone || '',
                    ...parseBirthDate(profileData.birth_date),
                });
                setPreviewUrl(profileData.avatar_url || '');
            } catch (error) {
                toast({
                    title: 'Lỗi',
                    description: 'Không thể tải thông tin hồ sơ',
                    status: 'error',
                    duration: 3000,
                    isClosable: true,
                });
            }
        };
        fetchProfile();
    }, [toast]);

    // Parse birth date if available
    const parseBirthDate = (birthDate) => {
        if (!birthDate) return { day: null, month: null, year: null };

        try {
            const date = new Date(birthDate);
            return {
                day: date.getDate(),
                month: date.getMonth() + 1, // JavaScript months are 0-indexed
                year: date.getFullYear(),
            };
        } catch (e) {
            return { day: null, month: null, year: null };
        }
    };

    // Handle input changes
    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    // Handle date changes - fixed to ensure state update triggers re-render
    const handleDateChange = (field, value) => {
        console.log("run")
        setFormData((prev) => {
            // Create a new object to ensure React detects the change
            return { ...prev, [field]: value };
        });
    };

    // Handle avatar file selection
    const handleFileSelect = (e) => {
        const file = e.target.files[0];
        if (file) {
            if (file.size > 1024 * 1024) { // 1MB limit
                toast({
                    title: 'Lỗi',
                    description: 'Dung lượng file vượt quá 1MB',
                    status: 'error',
                    duration: 3000,
                    isClosable: true,
                });
                return;
            }

            // Check if file type is supported
            if (!SUPPORTED_IMAGE_TYPES.includes(file.type)) {
                toast({
                    title: 'Lỗi',
                    description: `Chỉ hỗ trợ định dạng ${SUPPORTED_FORMATS_DISPLAY}`,
                    status: 'error',
                    duration: 3000,
                    isClosable: true,
                });
                return;
            }

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
        setIsSubmitting(true);

        try {
            // Create birth date string if all components are present
            let birthDateString = null;
            if (formData.day && formData.month && formData.year) {
                birthDateString = `${formData.year}-${formData.month
                    .toString()
                    .padStart(2, '0')}-${formData.day.toString().padStart(2, '0')}`;
            }

            // Prepare user data for update
            const userData = {
                fullname: formData.fullname,
                birth_date: birthDateString,
                phone: formData.phone || null,
                email: formData.email,
                avatar_url: user?.avatar_url || '', // Will update if new avatar is uploaded
            };

            // Handle avatar upload if a new file is selected
            if (selectedFile) {
                const presignedRequest = {
                    file_name: selectedFile.name,
                    file_size: selectedFile.size,
                    content_type: selectedFile.type,
                };

                // Get presigned URL
                const presignedResponse = await userMeService.getPresignedUrl(presignedRequest);
                const presignedUrl = presignedResponse.url;

                // Upload file to presigned URL
                await fetch(presignedUrl, {
                    method: 'PUT',
                    body: selectedFile,
                    headers: {
                        'Content-Type': selectedFile.type,
                    },
                });

                // Update avatar_url in userData (assuming the presigned URL service returns the final URL or we derive it)
                userData.avatar_url = presignedUrl.split('?')[0]; // Assuming the final URL is the base URL without query params
            }

            // Update user profile
            const updatedUser = await userMeService.updateProfile(userData);
            setUser(updatedUser);
            setPreviewUrl(updatedUser.avatar_url || '');
            setSelectedFile(null); // Clear selected file after successful upload

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

    // Force a re-render when formData changes
    useEffect(() => {
        // This empty effect ensures component re-renders when formData changes
    }, [formData]);

    return (
        <Box
            as="form"
            onSubmit={handleSubmit}
            maxH="calc(100vh - 200px)"
            overflowY="auto"
            pr={2}
            sx={{
                '&::-webkit-scrollbar': {
                    width: '4px',
                },
                '&::-webkit-scrollbar-thumb': {
                    backgroundColor: 'rgba(0,0,0,0.2)',
                    borderRadius: '2px',
                },
            }}
        >
            <Heading as="h1" size="lg" mb={4}>
                Hồ Sơ Của Tôi
            </Heading>
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
                                <Td width="180px" pr={2} pl={0} py={4} verticalAlign="center">
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
                                <Td width="180px" pr={2} pl={0} py={4} verticalAlign="center">
                                    <Text fontWeight="medium">Email</Text>
                                </Td>
                                <Td py={4} pl={4}>
                                    <Text>{formData.email}</Text>
                                </Td>
                            </Tr>

                            <Tr>
                                <Td width="180px" pr={2} pl={0} py={4} verticalAlign="center">
                                    <Text fontWeight="medium">Số điện thoại</Text>
                                </Td>
                                <Td py={4} pl={4}>
                                    {formData.phone === "" || formData.phone === null ? (
                                        formData._isAddingPhone ? (
                                            <Input
                                                name="phone"
                                                value={formData.phone || ""}
                                                onChange={handleChange}
                                                size="md"
                                                width="100%"
                                                maxW="400px"
                                                placeholder="Nhập số điện thoại"
                                            />
                                        ) : (
                                            <Button
                                                size="sm"
                                                variant="outline"
                                                colorScheme="blue"
                                                borderRadius="sm"
                                                onClick={() => {
                                                    setFormData((prev) => ({
                                                        ...prev,
                                                        _isAddingPhone: true,
                                                    }));
                                                }}
                                            >
                                                Thêm
                                            </Button>
                                        )
                                    ) : (
                                        <Flex alignItems="center" gap={4}>
                                            <Text>{formData.phone}</Text>
                                            <Button
                                                size="sm"
                                                variant="outline"
                                                colorScheme="blue"
                                                borderRadius="sm"
                                                onClick={() => {
                                                    setFormData((prev) => ({
                                                        ...prev,
                                                        _isEditingPhone: true,
                                                    }));
                                                }}
                                            >
                                                Sửa
                                            </Button>
                                        </Flex>
                                    )}
                                    {formData._isEditingPhone && (
                                        <Input
                                            name="phone"
                                            value={formData.phone || ""}
                                            onChange={handleChange}
                                            size="md"
                                            width="100%"
                                            maxW="400px"
                                            placeholder="Nhập số điện thoại"
                                            mt={2}
                                        />
                                    )}
                                </Td>
                            </Tr>
                            <Tr>
                                <Td width="180px" pr={2} pl={0} py={4} verticalAlign="center">
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
                                <Td width="180px" pr={2} pl={0} py={6} verticalAlign="top"></Td>
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
                        src={previewUrl}
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
                        Định dạng: {SUPPORTED_FORMATS_DISPLAY}
                    </Text>
                </Box>
            </Flex>
        </Box>
    );
};

export default UserProfile;