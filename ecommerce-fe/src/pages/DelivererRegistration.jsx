import {
    Box,
    Button,
    Container,
    FormControl,
    FormLabel,
    FormErrorMessage,
    Input,
    VStack,
    HStack,
    Text,
    Alert,
    AlertIcon,
    AlertTitle,
    AlertDescription,
    Progress,
    Divider,
    Icon,
    useToast,
    Grid,
    GridItem,
    Card,
    CardBody,
    Heading,
    Image,
    Flex,
    Badge,
    SimpleGrid,
    Checkbox,
    CheckboxGroup,
    Select,
    Spinner,
} from '@chakra-ui/react';
import { useState, useCallback, useLayoutEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDropzone } from 'react-dropzone';
import { FaShippingFast, FaUpload, FaIdCard, FaCheckCircle, FaTimesCircle, FaMotorcycle, FaCar, FaTruck } from 'react-icons/fa';
import useAuth from "../hooks/useAuth.js";
import userMeService from "../services/userMeService.js";
import delivererService from "../services/delivererService.js"; // Import service để gọi API địa danh

const DelivererRegistration = () => {
    const [currentStep, setCurrentStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const [uploadedImages, setUploadedImages] = useState({
        drivingLicenseFront: null,
        drivingLicenseBack: null
    });
    const toast = useToast();
    const navigate = useNavigate();
    const { user } = useAuth();

    const [formData, setFormData] = useState({
        drivingLicenseNumber: '',
        vehicleType: '',
        vehicleLicensePlate: '',
        selectedProvince: '',
        selectedDistrict: '',
        selectedWard: '',
        selectedProvinceName: '',
        selectedDistrictName: '',
        selectedWardName: ''
    });

    // States cho địa danh
    const [provinces, setProvinces] = useState([]);
    const [districts, setDistricts] = useState([]);
    const [wards, setWards] = useState([]);
    const [isLoadingProvinces, setIsLoadingProvinces] = useState(false);
    const [isLoadingDistricts, setIsLoadingDistricts] = useState(false);
    const [isLoadingWards, setIsLoadingWards] = useState(false);

    const [errors, setErrors] = useState({});

    const totalSteps = 3;

    const vehicleTypes = [
        { value: 'Xe máy', label: 'Xe máy', icon: FaMotorcycle },
        { value: 'Ô tô', label: 'Ô tô', icon: FaCar },
        { value: 'Xe tải nhỏ', label: 'Xe tải nhỏ', icon: FaTruck }
    ];

    // Fetch provinces khi component mount
    useLayoutEffect(() => {
        fetchProvinces();
    }, []);

    // Fetch districts when province changes
    useLayoutEffect(() => {
        if (formData.selectedProvince) {
            fetchDistricts(formData.selectedProvince);
            setFormData(prev => ({ ...prev, selectedDistrict: '', selectedWard: '' }));
            setWards([]);
        } else {
            setDistricts([]);
            setWards([]);
        }
    }, [formData.selectedProvince]);

    // Fetch wards when district changes
    useLayoutEffect(() => {
        if (formData.selectedProvince && formData.selectedDistrict) {
            fetchWards(formData.selectedProvince, formData.selectedDistrict);
            setFormData(prev => ({ ...prev, selectedWard: '' }));
        } else {
            setWards([]);
        }
    }, [formData.selectedDistrict]);

    const fetchProvinces = async () => {
        setIsLoadingProvinces(true);
        try {
            const provinces = await userMeService.getProvinces();
            setProvinces(provinces || []);
        } catch (error) {
            console.error('Error fetching provinces:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể tải danh sách tỉnh/thành phố',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoadingProvinces(false);
        }
    };

    const fetchDistricts = async (provinceId) => {
        if (!provinceId) return;

        setIsLoadingDistricts(true);
        try {
            const districts = await userMeService.getDistricts(provinceId);
            setDistricts(districts || []);
        } catch (error) {
            console.error('Error fetching districts:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể tải danh sách quận/huyện',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoadingDistricts(false);
        }
    };

    const fetchWards = async (provinceId, districtId) => {
        if (!provinceId || !districtId) return;

        setIsLoadingWards(true);
        try {
            const wards = await userMeService.getWards(provinceId, districtId);
            setWards(wards || []);
        } catch (error) {
            console.error('Error fetching wards:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể tải danh sách phường/xã',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoadingWards(false);
        }
    };

    // Handle form input changes
    const handleInputChange = (field, value) => {
        setFormData(prev => ({
            ...prev,
            [field]: value
        }));

        // Update corresponding names when selecting location
        if (field === 'selectedProvince') {
            const province = provinces.find(p => p.id === value);
            setFormData(prev => ({
                ...prev,
                selectedProvince: value,
                selectedProvinceName: province ? province.name : '',
                selectedDistrict: '',
                selectedWard: '',
                selectedDistrictName: '',
                selectedWardName: ''
            }));
        } else if (field === 'selectedDistrict') {
            const district = districts.find(d => d.id === value);
            setFormData(prev => ({
                ...prev,
                selectedDistrict: value,
                selectedDistrictName: district ? district.name : '',
                selectedWard: '',
                selectedWardName: ''
            }));
        } else if (field === 'selectedWard') {
            const ward = wards.find(w => w.id === value);
            setFormData(prev => ({
                ...prev,
                selectedWard: value,
                selectedWardName: ward ? ward.name : ''
            }));
        }

        // Clear error when user starts typing
        if (errors[field]) {
            setErrors(prev => ({
                ...prev,
                [field]: ''
            }));
        }
    };

    // File upload handlers
    const createImageDropHandler = (imageType) => useCallback((acceptedFiles) => {
        const file = acceptedFiles[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = () => {
                setUploadedImages(prev => ({
                    ...prev,
                    [imageType]: {
                        file: file,
                        name: file.name,
                        size: file.size,
                        url: reader.result
                    }
                }));
            };
            reader.readAsDataURL(file);
        }
    }, []);

    const onDrivingLicenseFrontDrop = createImageDropHandler('drivingLicenseFront');
    const onDrivingLicenseBackDrop = createImageDropHandler('drivingLicenseBack');

    const { getRootProps: getDLFrontRootProps, getInputProps: getDLFrontInputProps, isDragActive: isDLFrontDragActive } = useDropzone({
        onDropAccepted: onDrivingLicenseFrontDrop,
        accept: {
            'image/*': ['.jpeg', '.jpg', '.png']
        },
        multiple: false,
        maxSize: 5 * 1024 * 1024 // 5MB
    });

    const { getRootProps: getDLBackRootProps, getInputProps: getDLBackInputProps, isDragActive: isDLBackDragActive } = useDropzone({
        onDropAccepted: onDrivingLicenseBackDrop,
        accept: {
            'image/*': ['.jpeg', '.jpg', '.png']
        },
        multiple: false,
        maxSize: 5 * 1024 * 1024 // 5MB
    });

    const removeImage = (imageType) => {
        setUploadedImages(prev => ({
            ...prev,
            [imageType]: null
        }));
    };

    // Validation
    const validateStep = (step) => {
        const newErrors = {};

        if (step === 1) {
            if (!formData.drivingLicenseNumber.trim()) newErrors.drivingLicenseNumber = 'Số bằng lái xe là bắt buộc';
            if (!formData.vehicleType) newErrors.vehicleType = 'Loại phương tiện là bắt buộc';
            if (!formData.vehicleLicensePlate.trim()) newErrors.vehicleLicensePlate = 'Biển số xe là bắt buộc';
        }

        if (step === 2) {
            if (!uploadedImages.drivingLicenseFront) newErrors.drivingLicenseFront = 'Vui lòng tải lên ảnh mặt trước bằng lái xe';
            if (!uploadedImages.drivingLicenseBack) newErrors.drivingLicenseBack = 'Vui lòng tải lên ảnh mặt sau bằng lái xe';
        }

        if (step === 3) {
            if (!formData.selectedWard) {
                newErrors.serviceAreas = 'Vui lòng chọn khu vực giao hàng';
            }
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const nextStep = () => {
        if (validateStep(currentStep)) {
            setCurrentStep(prev => Math.min(prev + 1, totalSteps));
        }
    };

    const prevStep = () => {
        setCurrentStep(prev => Math.max(prev - 1, 1));
    };

    const handleSubmit = async () => {
        if (!validateStep(currentStep)) return;

        setLoading(true);
        try {
            // Prepare files object
            const files = {
                drivingLicenseFront: uploadedImages.drivingLicenseFront?.file,
                drivingLicenseBack: uploadedImages.drivingLicenseBack?.file
            };

            // Use the service to handle complete registration process
            const result = await delivererService.completeRegistration(formData, files);

            if (result.success) {
                toast({
                    title: 'Đăng ký thành công!',
                    description: 'Đơn đăng ký của bạn đã được gửi và đang chờ xét duyệt.',
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });

                navigate('/user/account/profile');
            } else {
                throw result.error;
            }
        } catch (error) {
            console.error('Registration error:', error);

            let errorMessage = 'Vui lòng thử lại sau.';
            if (error.response?.data?.error?.message) {
                errorMessage = error.response.data.error.message;
            } else if (error.message) {
                errorMessage = error.message;
            }

            toast({
                title: 'Có lỗi xảy ra',
                description: errorMessage,
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setLoading(false);
        }
    };

    const formatFileSize = (size) => {
        if (size < 1024) return size + ' B';
        if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB';
        return (size / (1024 * 1024)).toFixed(1) + ' MB';
    };

    return (
        <Container maxW="4xl" py={8}>
            <VStack spacing={8} align="stretch">
                {/* Header */}
                <Box textAlign="center">
                    <Icon as={FaShippingFast} w={16} h={16} color="green.500" mb={4} />
                    <Heading size="xl" mb={2}>Đăng ký trở thành người giao hàng</Heading>
                    <Text color="gray.600" fontSize="lg">
                        Tham gia đội ngũ giao hàng và tạo thu nhập linh hoạt
                    </Text>
                </Box>

                {/* Progress */}
                <Box>
                    <Progress value={(currentStep / totalSteps) * 100} colorScheme="green" size="lg" borderRadius="md" />
                    <HStack justify="space-between" mt={2}>
                        <Text fontSize="sm" color="gray.600">Bước {currentStep} / {totalSteps}</Text>
                        <Text fontSize="sm" color="gray.600">
                            {currentStep === 1 && 'Thông tin bằng lái & phương tiện'}
                            {currentStep === 2 && 'Tải lên bằng lái xe'}
                            {currentStep === 3 && 'Khu vực giao hàng'}
                        </Text>
                    </HStack>
                </Box>

                {/* Step Content */}
                <Card>
                    <CardBody p={8}>
                        {currentStep === 1 && (
                            <VStack spacing={6} align="stretch">
                                <Heading size="md" color="green.600" mb={4}>
                                    <Icon as={FaIdCard} mr={2} />
                                    Thông tin bằng lái xe & phương tiện
                                </Heading>

                                <FormControl isInvalid={errors.drivingLicenseNumber}>
                                    <FormLabel>Số bằng lái xe *</FormLabel>
                                    <Input
                                        value={formData.drivingLicenseNumber}
                                        onChange={(e) => handleInputChange('drivingLicenseNumber', e.target.value)}
                                        placeholder="Nhập số bằng lái xe"
                                    />
                                    <FormErrorMessage>{errors.drivingLicenseNumber}</FormErrorMessage>
                                </FormControl>

                                <Divider />

                                <Box>
                                    <Text fontWeight="bold" mb={4}>Thông tin phương tiện *</Text>

                                    <FormControl isInvalid={errors.vehicleType} mb={6}>
                                        <FormLabel>Loại phương tiện</FormLabel>
                                        <Grid templateColumns={{ base: "1fr", md: "repeat(3, 1fr)" }} gap={4}>
                                            {vehicleTypes.map((vehicle) => (
                                                <Card
                                                    key={vehicle.value}
                                                    cursor="pointer"
                                                    border="2px solid"
                                                    borderColor={formData.vehicleType === vehicle.value ? "green.500" : "gray.200"}
                                                    bg={formData.vehicleType === vehicle.value ? "green.50" : "white"}
                                                    _hover={{ borderColor: "green.300", bg: "green.50" }}
                                                    onClick={() => handleInputChange('vehicleType', vehicle.value)}
                                                    transition="all 0.2s"
                                                >
                                                    <CardBody textAlign="center" py={6}>
                                                        <VStack spacing={3}>
                                                            <Icon
                                                                as={vehicle.icon}
                                                                w={8}
                                                                h={8}
                                                                color={formData.vehicleType === vehicle.value ? "green.500" : "gray.400"}
                                                            />
                                                            <Text
                                                                fontWeight="medium"
                                                                color={formData.vehicleType === vehicle.value ? "green.600" : "gray.600"}
                                                            >
                                                                {vehicle.label}
                                                            </Text>
                                                            {formData.vehicleType === vehicle.value && (
                                                                <Icon as={FaCheckCircle} color="green.500" />
                                                            )}
                                                        </VStack>
                                                    </CardBody>
                                                </Card>
                                            ))}
                                        </Grid>
                                        <FormErrorMessage>{errors.vehicleType}</FormErrorMessage>
                                    </FormControl>

                                    <FormControl isInvalid={errors.vehicleLicensePlate}>
                                        <FormLabel>Biển số xe *</FormLabel>
                                        <Input
                                            value={formData.vehicleLicensePlate}
                                            onChange={(e) => handleInputChange('vehicleLicensePlate', e.target.value.toUpperCase())}
                                            placeholder="VD: 30A-12345"
                                            textTransform="uppercase"
                                        />
                                        <FormErrorMessage>{errors.vehicleLicensePlate}</FormErrorMessage>
                                    </FormControl>
                                </Box>
                            </VStack>
                        )}

                        {currentStep === 2 && (
                            <VStack spacing={6} align="stretch">
                                <Heading size="md" color="green.600" mb={4}>
                                    <Icon as={FaUpload} mr={2} />
                                    Tải lên bằng lái xe
                                </Heading>

                                <Alert status="info" borderRadius="md">
                                    <AlertIcon />
                                    <Box>
                                        <AlertTitle>Yêu cầu:</AlertTitle>
                                        <AlertDescription>
                                            Vui lòng tải lên ảnh rõ nét của bằng lái xe (mặt trước và mặt sau).
                                            Đảm bảo thông tin trên ảnh có thể đọc được và bằng lái phù hợp với loại phương tiện đã chọn.
                                        </AlertDescription>
                                    </Box>
                                </Alert>

                                <Grid templateColumns={{ base: "1fr", md: "1fr 1fr" }} gap={6}>
                                    {/* Driving License Front */}
                                    <GridItem>
                                        <FormControl isInvalid={errors.drivingLicenseFront}>
                                            <FormLabel>Mặt trước bằng lái xe *</FormLabel>
                                            <Box
                                                {...getDLFrontRootProps()}
                                                border="2px dashed"
                                                borderColor={isDLFrontDragActive ? "green.300" : "gray.300"}
                                                borderRadius="md"
                                                p={6}
                                                textAlign="center"
                                                cursor="pointer"
                                                _hover={{ borderColor: "green.400", bg: "green.50" }}
                                                bg={isDLFrontDragActive ? "green.50" : "gray.50"}
                                                transition="all 0.2s"
                                                minH="200px"
                                                display="flex"
                                                alignItems="center"
                                                justifyContent="center"
                                            >
                                                <input {...getDLFrontInputProps()} />
                                                {uploadedImages.drivingLicenseFront ? (
                                                    <VStack spacing={3}>
                                                        <Image
                                                            src={uploadedImages.drivingLicenseFront.url}
                                                            alt="Mặt trước bằng lái xe"
                                                            maxH="150px"
                                                            objectFit="contain"
                                                            borderRadius="md"
                                                        />
                                                        <Text color="green.600" fontWeight="medium" fontSize="sm">
                                                            <Icon as={FaCheckCircle} mr={1} />
                                                            Đã tải lên
                                                        </Text>
                                                        <Button
                                                            size="xs"
                                                            colorScheme="red"
                                                            variant="ghost"
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                removeImage('drivingLicenseFront');
                                                            }}
                                                        >
                                                            <Icon as={FaTimesCircle} mr={1} />
                                                            Xóa
                                                        </Button>
                                                    </VStack>
                                                ) : (
                                                    <VStack spacing={3}>
                                                        <Icon as={FaUpload} w={8} h={8} color="gray.400" />
                                                        <VStack spacing={1}>
                                                            <Text fontWeight="medium" fontSize="sm">
                                                                Tải lên mặt trước
                                                            </Text>
                                                            <Text fontSize="xs" color="gray.600">
                                                                JPG, PNG dưới 5MB
                                                            </Text>
                                                        </VStack>
                                                    </VStack>
                                                )}
                                            </Box>
                                            <FormErrorMessage>{errors.drivingLicenseFront}</FormErrorMessage>
                                        </FormControl>
                                    </GridItem>

                                    {/* Driving License Back */}
                                    <GridItem>
                                        <FormControl isInvalid={errors.drivingLicenseBack}>
                                            <FormLabel>Mặt sau bằng lái xe *</FormLabel>
                                            <Box
                                                {...getDLBackRootProps()}
                                                border="2px dashed"
                                                borderColor={isDLBackDragActive ? "green.300" : "gray.300"}
                                                borderRadius="md"
                                                p={6}
                                                textAlign="center"
                                                cursor="pointer"
                                                _hover={{ borderColor: "green.400", bg: "green.50" }}
                                                bg={isDLBackDragActive ? "green.50" : "gray.50"}
                                                transition="all 0.2s"
                                                minH="200px"
                                                display="flex"
                                                alignItems="center"
                                                justifyContent="center"
                                            >
                                                <input {...getDLBackInputProps()} />
                                                {uploadedImages.drivingLicenseBack ? (
                                                    <VStack spacing={3}>
                                                        <Image
                                                            src={uploadedImages.drivingLicenseBack.url}
                                                            alt="Mặt sau bằng lái xe"
                                                            maxH="150px"
                                                            objectFit="contain"
                                                            borderRadius="md"
                                                        />
                                                        <Text color="green.600" fontWeight="medium" fontSize="sm">
                                                            <Icon as={FaCheckCircle} mr={1} />
                                                            Đã tải lên
                                                        </Text>
                                                        <Button
                                                            size="xs"
                                                            colorScheme="red"
                                                            variant="ghost"
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                removeImage('drivingLicenseBack');
                                                            }}
                                                        >
                                                            <Icon as={FaTimesCircle} mr={1} />
                                                            Xóa
                                                        </Button>
                                                    </VStack>
                                                ) : (
                                                    <VStack spacing={3}>
                                                        <Icon as={FaUpload} w={8} h={8} color="gray.400" />
                                                        <VStack spacing={1}>
                                                            <Text fontWeight="medium" fontSize="sm">
                                                                Tải lên mặt sau
                                                            </Text>
                                                            <Text fontSize="xs" color="gray.600">
                                                                JPG, PNG dưới 5MB
                                                            </Text>
                                                        </VStack>
                                                    </VStack>
                                                )}
                                            </Box>
                                            <FormErrorMessage>{errors.drivingLicenseBack}</FormErrorMessage>
                                        </FormControl>
                                    </GridItem>
                                </Grid>
                            </VStack>
                        )}

                        {currentStep === 3 && (
                            <VStack spacing={6} align="stretch">
                                <Heading size="md" color="green.600" mb={4}>
                                    <Icon as={FaShippingFast} mr={2} />
                                    Chọn khu vực giao hàng
                                </Heading>

                                <Alert status="info" borderRadius="md">
                                    <AlertIcon />
                                    <Box>
                                        <AlertTitle>Lưu ý:</AlertTitle>
                                        <AlertDescription>
                                            Chọn các khu vực mà bạn có thể giao hàng. Bạn có thể thay đổi sau khi được duyệt.
                                        </AlertDescription>
                                    </Box>
                                </Alert>

                                <FormControl isInvalid={errors.serviceAreas}>
                                    <FormLabel>Khu vực giao hàng *</FormLabel>
                                    <Text fontSize="sm" color="gray.600" mb={4}>
                                        Chọn theo thứ tự: Tỉnh/Thành phố → Quận/Huyện → Phường/Xã:
                                    </Text>

                                    <VStack spacing={4} align="stretch">
                                        {/* Province Selection */}
                                        <FormControl>
                                            <FormLabel fontSize="sm">Tỉnh/Thành phố</FormLabel>
                                            <Select
                                                placeholder="Chọn Tỉnh/Thành phố"
                                                value={formData.selectedProvince}
                                                onChange={(e) => handleInputChange('selectedProvince', e.target.value)}
                                                isDisabled={isLoadingProvinces}
                                            >
                                                {provinces.map(province => (
                                                    <option key={province.id} value={province.id}>
                                                        {province.name}
                                                    </option>
                                                ))}
                                            </Select>
                                            {isLoadingProvinces && (
                                                <Flex align="center" mt={2}>
                                                    <Spinner size="sm" mr={2} />
                                                    <Text fontSize="sm" color="gray.500">Đang tải...</Text>
                                                </Flex>
                                            )}
                                        </FormControl>

                                        {/* District Selection */}
                                        <FormControl>
                                            <FormLabel fontSize="sm">Quận/Huyện</FormLabel>
                                            <Select
                                                placeholder="Chọn Quận/Huyện"
                                                value={formData.selectedDistrict}
                                                onChange={(e) => handleInputChange('selectedDistrict', e.target.value)}
                                                isDisabled={!formData.selectedProvince || isLoadingDistricts}
                                            >
                                                {districts.map(district => (
                                                    <option key={district.id} value={district.id}>
                                                        {district.name}
                                                    </option>
                                                ))}
                                            </Select>
                                            {isLoadingDistricts && (
                                                <Flex align="center" mt={2}>
                                                    <Spinner size="sm" mr={2} />
                                                    <Text fontSize="sm" color="gray.500">Đang tải...</Text>
                                                </Flex>
                                            )}
                                        </FormControl>

                                        {/* Ward Selection */}
                                        <FormControl>
                                            <FormLabel fontSize="sm">Phường/Xã</FormLabel>
                                            <Select
                                                placeholder="Chọn Phường/Xã"
                                                value={formData.selectedWard}
                                                onChange={(e) => handleInputChange('selectedWard', e.target.value)}
                                                isDisabled={!formData.selectedDistrict || isLoadingWards}
                                            >
                                                {wards.map(ward => (
                                                    <option key={ward.id} value={ward.id}>
                                                        {ward.name}
                                                    </option>
                                                ))}
                                            </Select>
                                            {isLoadingWards && (
                                                <Flex align="center" mt={2}>
                                                    <Spinner size="sm" mr={2} />
                                                    <Text fontSize="sm" color="gray.500">Đang tải...</Text>
                                                </Flex>
                                            )}
                                        </FormControl>
                                    </VStack>

                                    <FormErrorMessage>{errors.serviceAreas}</FormErrorMessage>
                                </FormControl>

                                {/* Selected Area Summary */}
                                {formData.selectedWard && (
                                    <Box>
                                        <Text fontWeight="medium" mb={3}>
                                            Khu vực đã chọn:
                                        </Text>
                                        <Card bg="green.50" borderColor="green.200">
                                            <CardBody py={3}>
                                                <Text fontWeight="medium" color="green.800">
                                                    {wards.find(w => w.id === formData.selectedWard)?.name}, {' '}
                                                    {districts.find(d => d.id === formData.selectedDistrict)?.name}, {' '}
                                                    {provinces.find(p => p.id === formData.selectedProvince)?.name}
                                                </Text>
                                            </CardBody>
                                        </Card>
                                    </Box>
                                )}
                            </VStack>
                        )}
                    </CardBody>
                </Card>

                {/* Navigation Buttons */}
                <HStack justify="space-between">
                    <Button
                        onClick={prevStep}
                        isDisabled={currentStep === 1}
                        variant="outline"
                        size="lg"
                    >
                        Quay lại
                    </Button>

                    <HStack>
                        <Text fontSize="sm" color="gray.600">
                            {currentStep} / {totalSteps}
                        </Text>
                    </HStack>

                    {currentStep < totalSteps ? (
                        <Button
                            onClick={nextStep}
                            colorScheme="green"
                            size="lg"
                        >
                            Tiếp theo
                        </Button>
                    ) : (
                        <Button
                            onClick={handleSubmit}
                            colorScheme="green"
                            size="lg"
                            isLoading={loading}
                            loadingText="Đang upload và đăng ký..."
                        >
                            Hoàn tất đăng ký
                        </Button>
                    )}
                </HStack>

                {/* Additional Info */}
                <Alert status="warning" borderRadius="md">
                    <AlertIcon />
                    <Box>
                        <AlertTitle>Lưu ý:</AlertTitle>
                        <AlertDescription>
                            Sau khi gửi đăng ký, chúng tôi sẽ xem xét và phản hồi trong vòng 2-3 ngày làm việc.
                            Bạn sẽ nhận được thông báo qua email khi có kết quả. Đảm bảo bằng lái xe phù hợp với loại phương tiện để quá trình duyệt diễn ra nhanh chóng.
                        </AlertDescription>
                    </Box>
                </Alert>
            </VStack>
        </Container>
    );
};

export default DelivererRegistration;