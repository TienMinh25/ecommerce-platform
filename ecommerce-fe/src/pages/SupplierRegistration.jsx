import {
    Box,
    Button,
    Container,
    FormControl,
    FormLabel,
    FormErrorMessage,
    Input,
    Textarea,
    VStack,
    HStack,
    Text,
    Alert,
    AlertIcon,
    AlertTitle,
    AlertDescription,
    Progress,
    Icon,
    useToast,
    Grid,
    GridItem,
    Card,
    CardBody,
    Heading,
    Image,
} from '@chakra-ui/react';
import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDropzone } from 'react-dropzone';
import { FaStore, FaUpload, FaFileAlt, FaCheckCircle, FaTimesCircle } from 'react-icons/fa';
import useAuth from "../hooks/useAuth.js";
import BusinessAddressSection from "../components/supplier/BusinessAddressSection.jsx";

const SupplierRegistration = () => {
    const [currentStep, setCurrentStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const [selectedAddress, setSelectedAddress] = useState(null);
    const toast = useToast();
    const navigate = useNavigate();
    const { user } = useAuth();

    const [formData, setFormData] = useState({
        companyName: '',
        contactPhone: '',
        description: '',
        logoUrl: '',
        businessAddress: {
            street: '',
            ward: '',
            district: '',
            city: '',
            country: 'Vietnam'
        },
        taxId: ''
    });

    const [uploadedDocuments, setUploadedDocuments] = useState({
        business_license: null,
        tax_certificate: null,
        id_card_front: null,
        id_card_back: null
    });

    const [errors, setErrors] = useState({});
    const totalSteps = 3;

    // Document types configuration
    const documentTypes = [
        {
            key: 'business_license',
            label: 'Giấy phép kinh doanh',
            required: true,
            description: 'Giấy chứng nhận đăng ký kinh doanh của công ty'
        },
        {
            key: 'tax_certificate',
            label: 'Chứng nhận đăng ký thuế',
            required: true,
            description: 'Chứng nhận đăng ký thuế doanh nghiệp'
        },
        {
            key: 'id_card_front',
            label: 'CCCD mặt trước',
            required: true,
            description: 'Ảnh mặt trước căn cước công dân người đại diện'
        },
        {
            key: 'id_card_back',
            label: 'CCCD mặt sau',
            required: true,
            description: 'Ảnh mặt sau căn cước công dân người đại diện'
        }
    ];

    // Handle form input changes
    const handleInputChange = (field, value) => {
        if (field.includes('.')) {
            const [parent, child] = field.split('.');
            setFormData(prev => ({
                ...prev,
                [parent]: {
                    ...prev[parent],
                    [child]: value
                }
            }));
        } else {
            setFormData(prev => ({
                ...prev,
                [field]: value
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

    // Handle address selection
    const handleSelectAddress = (address) => {
        setSelectedAddress(address);
        setFormData(prev => ({
            ...prev,
            businessAddress: {
                street: address.street,
                ward: address.ward,
                district: address.district,
                city: address.province,
                country: address.country === 'Việt Nam' ? 'Vietnam' : address.country
            }
        }));

        // Clear address-related errors
        setErrors(prev => ({
            ...prev,
            'businessAddress.street': '',
            'businessAddress.ward': '',
            'businessAddress.district': '',
            'businessAddress.city': ''
        }));
    };

    // Clear selected address
    const handleClearAddress = () => {
        setSelectedAddress(null);
        setFormData(prev => ({
            ...prev,
            businessAddress: {
                street: '',
                ward: '',
                district: '',
                city: '',
                country: 'Vietnam'
            }
        }));
    };

    // Logo upload handler
    const onLogoDropAccepted = useCallback((acceptedFiles) => {
        const file = acceptedFiles[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = () => {
                setFormData(prev => ({
                    ...prev,
                    logoUrl: reader.result
                }));
            };
            reader.readAsDataURL(file);
        }
    }, []);

    // Document upload handler
    const handleDocumentUpload = useCallback((docType, acceptedFiles) => {
        const file = acceptedFiles[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = () => {
                setUploadedDocuments(prev => ({
                    ...prev,
                    [docType]: reader.result // Will be S3 URL in production
                }));
            };
            reader.readAsDataURL(file);
        }
    }, []);

    // Remove document
    const removeDocument = (docType) => {
        setUploadedDocuments(prev => ({
            ...prev,
            [docType]: null
        }));
    };

    // Logo dropzone
    const { getRootProps: getLogoRootProps, getInputProps: getLogoInputProps, isDragActive: isLogoDragActive } = useDropzone({
        onDropAccepted: onLogoDropAccepted,
        accept: {
            'image/*': ['.jpeg', '.jpg', '.png', '.gif']
        },
        multiple: false,
        maxSize: 5 * 1024 * 1024
    });

    // Individual dropzones for each document type
    const businessLicenseDropzone = useDropzone({
        onDropAccepted: (files) => handleDocumentUpload('business_license', files),
        accept: {
            'image/*': ['.jpeg', '.jpg', '.png', '.gif'],
            'application/pdf': ['.pdf']
        },
        multiple: false,
        maxSize: 10 * 1024 * 1024
    });

    const taxCertificateDropzone = useDropzone({
        onDropAccepted: (files) => handleDocumentUpload('tax_certificate', files),
        accept: {
            'image/*': ['.jpeg', '.jpg', '.png', '.gif'],
            'application/pdf': ['.pdf']
        },
        multiple: false,
        maxSize: 10 * 1024 * 1024
    });

    const idCardFrontDropzone = useDropzone({
        onDropAccepted: (files) => handleDocumentUpload('id_card_front', files),
        accept: {
            'image/*': ['.jpeg', '.jpg', '.png', '.gif'],
            'application/pdf': ['.pdf']
        },
        multiple: false,
        maxSize: 10 * 1024 * 1024
    });

    const idCardBackDropzone = useDropzone({
        onDropAccepted: (files) => handleDocumentUpload('id_card_back', files),
        accept: {
            'image/*': ['.jpeg', '.jpg', '.png', '.gif'],
            'application/pdf': ['.pdf']
        },
        multiple: false,
        maxSize: 10 * 1024 * 1024
    });

    // Get dropzone by document type
    const getDropzoneByType = (docType) => {
        switch (docType) {
            case 'business_license': return businessLicenseDropzone;
            case 'tax_certificate': return taxCertificateDropzone;
            case 'id_card_front': return idCardFrontDropzone;
            case 'id_card_back': return idCardBackDropzone;
            default: return businessLicenseDropzone;
        }
    };

    // Validation
    const validateStep = (step) => {
        const newErrors = {};

        if (step === 1) {
            if (!formData.companyName.trim()) newErrors.companyName = 'Tên công ty là bắt buộc';
            if (!formData.contactPhone.trim()) newErrors.contactPhone = 'Số điện thoại là bắt buộc';
            if (!formData.taxId.trim()) newErrors.taxId = 'Mã số thuế là bắt buộc';
            if (!formData.businessAddress.street.trim()) newErrors['businessAddress.street'] = 'Địa chỉ là bắt buộc';
            if (!formData.businessAddress.ward.trim()) newErrors['businessAddress.ward'] = 'Phường/Xã là bắt buộc';
            if (!formData.businessAddress.district.trim()) newErrors['businessAddress.district'] = 'Quận/Huyện là bắt buộc';
            if (!formData.businessAddress.city.trim()) newErrors['businessAddress.city'] = 'Tỉnh/Thành phố là bắt buộc';
        }

        if (step === 2) {
            if (!formData.logoUrl) newErrors.logoUrl = 'Logo công ty là bắt buộc';
        }

        if (step === 3) {
            documentTypes.forEach(docType => {
                if (docType.required && !uploadedDocuments[docType.key]) {
                    newErrors[docType.key] = `${docType.label} là bắt buộc`;
                }
            });
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
            const submissionData = {
                companyName: formData.companyName,
                contactPhone: formData.contactPhone,
                taxId: formData.taxId,
                description: formData.description,
                logoUrl: formData.logoUrl,
                businessAddress: formData.businessAddress,
                documents: uploadedDocuments
            };

            console.log('Submission data:', submissionData);
            await new Promise(resolve => setTimeout(resolve, 2000));

            toast({
                title: 'Đăng ký thành công!',
                description: 'Đơn đăng ký của bạn đã được gửi và đang chờ xét duyệt.',
                status: 'success',
                duration: 5000,
                isClosable: true,
            });

            navigate('/user/account/profile');
        } catch (error) {
            toast({
                title: 'Có lỗi xảy ra',
                description: 'Vui lòng thử lại sau.',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setLoading(false);
        }
    };

    return (
        <Container maxW="4xl" py={8}>
            <VStack spacing={8} align="stretch">
                {/* Header */}
                <Box textAlign="center">
                    <Icon as={FaStore} w={16} h={16} color="blue.500" mb={4} />
                    <Heading size="xl" mb={2}>Đăng ký trở thành nhà cung cấp</Heading>
                    <Text color="gray.600" fontSize="lg">
                        Tham gia với chúng tôi để mở rộng thị trường kinh doanh của bạn
                    </Text>
                </Box>

                {/* Progress */}
                <Box>
                    <Progress value={(currentStep / totalSteps) * 100} colorScheme="blue" size="lg" borderRadius="md" />
                    <HStack justify="space-between" mt={2}>
                        <Text fontSize="sm" color="gray.600">Bước {currentStep} / {totalSteps}</Text>
                        <Text fontSize="sm" color="gray.600">
                            {currentStep === 1 && 'Thông tin công ty'}
                            {currentStep === 2 && 'Logo & Mô tả'}
                            {currentStep === 3 && 'Tài liệu chứng minh'}
                        </Text>
                    </HStack>
                </Box>

                {/* Step Content */}
                <Card>
                    <CardBody p={8}>
                        {/* Step 1: Company Info */}
                        {currentStep === 1 && (
                            <VStack spacing={6} align="stretch">
                                <Heading size="md" color="blue.600" mb={4}>
                                    <Icon as={FaStore} mr={2} />
                                    Thông tin công ty
                                </Heading>

                                <Grid templateColumns={{ base: "1fr", md: "1fr 1fr" }} gap={6}>
                                    <GridItem>
                                        <FormControl isInvalid={errors.companyName}>
                                            <FormLabel>Tên công ty *</FormLabel>
                                            <Input
                                                value={formData.companyName}
                                                onChange={(e) => handleInputChange('companyName', e.target.value)}
                                                placeholder="Nhập tên công ty của bạn"
                                            />
                                            <FormErrorMessage>{errors.companyName}</FormErrorMessage>
                                        </FormControl>
                                    </GridItem>

                                    <GridItem>
                                        <FormControl isInvalid={errors.contactPhone}>
                                            <FormLabel>Số điện thoại liên hệ *</FormLabel>
                                            <Input
                                                value={formData.contactPhone}
                                                onChange={(e) => handleInputChange('contactPhone', e.target.value)}
                                                placeholder="Số điện thoại"
                                            />
                                            <FormErrorMessage>{errors.contactPhone}</FormErrorMessage>
                                        </FormControl>
                                    </GridItem>

                                    <GridItem>
                                        <FormControl isInvalid={errors.taxId}>
                                            <FormLabel>Mã số thuế *</FormLabel>
                                            <Input
                                                value={formData.taxId}
                                                onChange={(e) => handleInputChange('taxId', e.target.value)}
                                                placeholder="Mã số thuế doanh nghiệp"
                                            />
                                            <FormErrorMessage>{errors.taxId}</FormErrorMessage>
                                        </FormControl>
                                    </GridItem>
                                </Grid>

                                <BusinessAddressSection
                                    formData={formData}
                                    errors={errors}
                                    selectedAddress={selectedAddress}
                                    onInputChange={handleInputChange}
                                    onSelectAddress={handleSelectAddress}
                                    onClearAddress={handleClearAddress}
                                />
                            </VStack>
                        )}

                        {/* Step 2: Logo & Description */}
                        {currentStep === 2 && (
                            <VStack spacing={6} align="stretch">
                                <Heading size="md" color="blue.600" mb={4}>
                                    <Icon as={FaUpload} mr={2} />
                                    Logo công ty & Mô tả
                                </Heading>

                                <FormControl isInvalid={errors.logoUrl}>
                                    <FormLabel>Logo công ty *</FormLabel>
                                    <Box
                                        {...getLogoRootProps()}
                                        border="2px dashed"
                                        borderColor={isLogoDragActive ? "blue.300" : "gray.300"}
                                        borderRadius="md"
                                        p={8}
                                        textAlign="center"
                                        cursor="pointer"
                                        _hover={{ borderColor: "blue.400", bg: "blue.50" }}
                                        bg={isLogoDragActive ? "blue.50" : "gray.50"}
                                        transition="all 0.2s"
                                    >
                                        <input {...getLogoInputProps()} />
                                        {formData.logoUrl ? (
                                            <VStack spacing={4}>
                                                <Image
                                                    src={formData.logoUrl}
                                                    alt="Logo công ty"
                                                    maxH="200px"
                                                    maxW="300px"
                                                    objectFit="contain"
                                                    borderRadius="md"
                                                />
                                                <Text color="green.600" fontWeight="medium">
                                                    <Icon as={FaCheckCircle} mr={2} />
                                                    Logo đã được tải lên
                                                </Text>
                                                <Text fontSize="sm" color="gray.600">
                                                    Kéo thả file khác để thay thế
                                                </Text>
                                            </VStack>
                                        ) : (
                                            <VStack spacing={4}>
                                                <Icon as={FaUpload} w={12} h={12} color="gray.400" />
                                                <VStack spacing={2}>
                                                    <Text fontWeight="medium">
                                                        Kéo thả file logo vào đây hoặc click để chọn
                                                    </Text>
                                                    <Text fontSize="sm" color="gray.600">
                                                        Chỉ chấp nhận file ảnh (JPG, PNG, GIF) dưới 5MB
                                                    </Text>
                                                </VStack>
                                            </VStack>
                                        )}
                                    </Box>
                                    <FormErrorMessage>{errors.logoUrl}</FormErrorMessage>
                                </FormControl>

                                <FormControl>
                                    <FormLabel>Mô tả về công ty</FormLabel>
                                    <Textarea
                                        value={formData.description}
                                        onChange={(e) => handleInputChange('description', e.target.value)}
                                        placeholder="Mô tả về công ty, lĩnh vực kinh doanh, sản phẩm chính..."
                                        rows={6}
                                        resize="vertical"
                                    />
                                    <Text fontSize="sm" color="gray.600" mt={2}>
                                        Mô tả chi tiết sẽ giúp khách hàng hiểu rõ hơn về công ty của bạn
                                    </Text>
                                </FormControl>
                            </VStack>
                        )}

                        {/* Step 3: Documents */}
                        {currentStep === 3 && (
                            <VStack spacing={6} align="stretch">
                                <Heading size="md" color="blue.600" mb={4}>
                                    <Icon as={FaFileAlt} mr={2} />
                                    Tài liệu chứng minh
                                </Heading>

                                <Alert status="info" borderRadius="md">
                                    <AlertIcon />
                                    <Box>
                                        <AlertTitle>Tài liệu bắt buộc:</AlertTitle>
                                        <AlertDescription>
                                            Vui lòng tải lên đầy đủ 4 loại tài liệu sau để hoàn tất đăng ký nhà cung cấp.
                                        </AlertDescription>
                                    </Box>
                                </Alert>

                                <VStack spacing={6}>
                                    {documentTypes.map((docType) => {
                                        const dropzone = getDropzoneByType(docType.key);
                                        const documentUrl = uploadedDocuments[docType.key];

                                        return (
                                            <Box key={docType.key} w="100%">
                                                <FormControl isInvalid={errors[docType.key]}>
                                                    <HStack mb={2} justify="space-between">
                                                        <VStack align="flex-start" spacing={0}>
                                                            <FormLabel mb={0} fontWeight="bold">
                                                                {docType.label} {docType.required && <Text as="span" color="red.500">*</Text>}
                                                            </FormLabel>
                                                            <Text fontSize="sm" color="gray.600">
                                                                {docType.description}
                                                            </Text>
                                                        </VStack>
                                                        {documentUrl && (
                                                            <Button
                                                                size="sm"
                                                                colorScheme="red"
                                                                variant="ghost"
                                                                onClick={() => removeDocument(docType.key)}
                                                            >
                                                                <Icon as={FaTimesCircle} mr={1} />
                                                                Xóa
                                                            </Button>
                                                        )}
                                                    </HStack>

                                                    <Box
                                                        {...dropzone.getRootProps()}
                                                        border="2px dashed"
                                                        borderColor={
                                                            errors[docType.key] ? "red.300" :
                                                                documentUrl ? "green.300" :
                                                                    dropzone.isDragActive ? "blue.300" : "gray.300"
                                                        }
                                                        borderRadius="md"
                                                        p={6}
                                                        textAlign="center"
                                                        cursor="pointer"
                                                        _hover={{
                                                            borderColor: documentUrl ? "green.400" : "blue.400",
                                                            bg: documentUrl ? "green.50" : "blue.50"
                                                        }}
                                                        bg={
                                                            documentUrl ? "green.50" :
                                                                dropzone.isDragActive ? "blue.50" : "gray.50"
                                                        }
                                                        transition="all 0.2s"
                                                        minH="150px"
                                                        display="flex"
                                                        alignItems="center"
                                                        justifyContent="center"
                                                    >
                                                        <input {...dropzone.getInputProps()} />

                                                        {documentUrl ? (
                                                            <VStack spacing={3}>
                                                                <Image
                                                                    src={documentUrl}
                                                                    alt={docType.label}
                                                                    maxH="120px"
                                                                    maxW="200px"
                                                                    objectFit="cover"
                                                                    borderRadius="md"
                                                                    border="2px solid"
                                                                    borderColor="green.300"
                                                                />

                                                                <VStack spacing={1}>
                                                                    <HStack>
                                                                        <Icon as={FaCheckCircle} color="green.500" />
                                                                        <Text color="green.700" fontWeight="medium" fontSize="sm">
                                                                            Đã tải lên thành công
                                                                        </Text>
                                                                    </HStack>
                                                                    <Text fontSize="xs" color="gray.500">
                                                                        Click để thay đổi file
                                                                    </Text>
                                                                </VStack>
                                                            </VStack>
                                                        ) : (
                                                            <VStack spacing={3}>
                                                                <Icon as={FaUpload} w={10} h={10} color="gray.400" />
                                                                <VStack spacing={1}>
                                                                    <Text fontWeight="medium" color="gray.700">
                                                                        Kéo thả file vào đây hoặc click để chọn
                                                                    </Text>
                                                                    <Text fontSize="sm" color="gray.500">
                                                                        Chấp nhận: JPG, PNG, PDF (tối đa 10MB)
                                                                    </Text>
                                                                </VStack>
                                                            </VStack>
                                                        )}
                                                    </Box>

                                                    <FormErrorMessage>{errors[docType.key]}</FormErrorMessage>
                                                </FormControl>
                                            </Box>
                                        );
                                    })}
                                </VStack>

                                {/* Upload Summary */}
                                <Box p={4} bg="gray.50" borderRadius="md">
                                    <Text fontWeight="medium" mb={2} fontSize="sm">
                                        Tiến độ tải tài liệu:
                                    </Text>
                                    <HStack spacing={4} flexWrap="wrap">
                                        {documentTypes.map(docType => (
                                            <HStack key={docType.key} spacing={2}>
                                                <Icon
                                                    as={uploadedDocuments[docType.key] ? FaCheckCircle : FaTimesCircle}
                                                    color={uploadedDocuments[docType.key] ? "green.500" : "gray.400"}
                                                />
                                                <Text fontSize="sm" color={uploadedDocuments[docType.key] ? "green.700" : "gray.500"}>
                                                    {docType.label}
                                                </Text>
                                            </HStack>
                                        ))}
                                    </HStack>
                                    <Progress
                                        value={(Object.values(uploadedDocuments).filter(Boolean).length / documentTypes.length) * 100}
                                        colorScheme="green"
                                        size="sm"
                                        mt={3}
                                        borderRadius="md"
                                    />
                                </Box>
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
                            colorScheme="blue"
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
                            loadingText="Đang gửi đăng ký..."
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
                            Sau khi gửi đăng ký, chúng tôi sẽ xem xét và phản hồi trong vòng 3-5 ngày làm việc.
                            Bạn sẽ nhận được thông báo qua email khi có kết quả.
                        </AlertDescription>
                    </Box>
                </Alert>
            </VStack>
        </Container>
    );
};

export default SupplierRegistration;