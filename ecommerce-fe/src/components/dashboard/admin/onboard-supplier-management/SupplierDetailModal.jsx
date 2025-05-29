import React, { useState, useEffect } from 'react';
import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalBody,
    ModalCloseButton,
    Box,
    VStack,
    HStack,
    Text,
    Badge,
    Divider,
    Icon,
    Button,
    useColorModeValue,
    Alert,
    AlertIcon,
    Spinner,
    Center,
    Tooltip,
    SimpleGrid,
    Card,
    CardBody,
    Heading,
    Avatar,
    Stack,
    StackDivider,
} from '@chakra-ui/react';

import {
    FiPhone,
    FiMapPin,
    FiCalendar,
    FiFileText,
    FiCheckCircle,
    FiXCircle,
    FiClock,
    FiShield,
    FiCheck,
    FiX,
} from 'react-icons/fi';
import { HiOutlineOfficeBuilding } from 'react-icons/hi';
import supplierService from '../../../../services/supplierService.js';

const SupplierDetailModal = ({ isOpen, onClose, supplierId }) => {
    const [supplier, setSupplier] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [actionLoading, setActionLoading] = useState(null);

    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.600');
    const textColorPrimary = useColorModeValue('gray.800', 'white');
    const textColorSecondary = useColorModeValue('gray.600', 'gray.300');

    // Document type mapping với tên tiếng Việt
    const documentTypeNames = {
        'business_license': 'Giấy phép kinh doanh',
        'tax_certificate': 'Giấy chứng nhận thuế',
        'id_card_front': 'CCCD mặt trước',
        'id_card_back': 'CCCD mặt sau'
    };

    const getStatusColor = (status) => {
        switch (status) {
            case 'active':
                return 'green';
            case 'pending':
                return 'yellow';
            case 'suspended':
                return 'red';
            default:
                return 'gray';
        }
    };

    const getStatusText = (status) => {
        switch (status) {
            case 'active':
                return 'Hoạt động';
            case 'pending':
                return 'Chờ duyệt';
            case 'suspended':
                return 'Tạm ngưng';
            default:
                return status;
        }
    };

    const getDocumentStatusColor = (status) => {
        switch (status) {
            case 'approved':
                return 'green';
            case 'pending':
                return 'yellow';
            case 'rejected':
                return 'red';
            default:
                return 'gray';
        }
    };

    const getDocumentStatusText = (status) => {
        switch (status) {
            case 'approved':
                return 'Đã duyệt';
            case 'pending':
                return 'Chờ duyệt';
            case 'rejected':
                return 'Từ chối';
            default:
                return status;
        }
    };

    const getDocumentStatusIcon = (status) => {
        switch (status) {
            case 'approved':
                return FiCheckCircle;
            case 'pending':
                return FiClock;
            case 'rejected':
                return FiXCircle;
            default:
                return FiFileText;
        }
    };

    const fetchSupplierDetail = async () => {
        if (!supplierId) return;

        setLoading(true);
        setError(null);

        try {
            const response = await supplierService.getSupplierById(supplierId);
            setSupplier(response.data);
        } catch (err) {
            setError(err.response?.data?.error?.message || 'Không thể tải thông tin nhà cung cấp');
        } finally {
            setLoading(false);
        }
    };

    const handleDocumentAction = async (documentId, action) => {
        setActionLoading(`${documentId}-${action}`);
        try {
            // Call API to approve/reject document
            await supplierService.updateDocumentStatus(documentId, action);
            // Refresh supplier data
            await fetchSupplierDetail();
        } catch (err) {
            console.error('Error updating document status:', err);
        } finally {
            setActionLoading(null);
        }
    };

    useEffect(() => {
        if (isOpen && supplierId) {
            fetchSupplierDetail();
        }
    }, [isOpen, supplierId]);

    const formatDate = (dateString) => {
        if (!dateString) return 'N/A';
        try {
            return new Date(dateString).toLocaleString('vi-VN');
        } catch (e) {
            return dateString;
        }
    };

    const openDocumentInNewTab = (url) => {
        if (url) {
            window.open(url, '_blank');
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} size="6xl" scrollBehavior="inside">
            <ModalOverlay />
            <ModalContent bg={bgColor} maxH="90vh">
                <ModalHeader borderBottomWidth="1px" borderColor={borderColor}>
                    <HStack spacing={3}>
                        <Icon as={HiOutlineOfficeBuilding} color="blue.500" />
                        <Text>Chi tiết nhà cung cấp</Text>
                    </HStack>
                </ModalHeader>
                <ModalCloseButton />

                <ModalBody p={6}>
                    {loading ? (
                        <Center py={12}>
                            <VStack spacing={4}>
                                <Spinner size="lg" color="blue.500" />
                                <Text color={textColorSecondary}>Đang tải thông tin...</Text>
                            </VStack>
                        </Center>
                    ) : error ? (
                        <Alert status="error" borderRadius="md">
                            <AlertIcon />
                            <Text>{error}</Text>
                        </Alert>
                    ) : supplier ? (
                        <VStack spacing={6} align="stretch">
                            {/* Header với logo và thông tin cơ bản */}
                            <Card>
                                <CardBody>
                                    <HStack spacing={6} align="start">
                                        <Avatar
                                            size="2xl"
                                            src={supplier.logo_thumbnail_url}
                                            name={supplier.company_name}
                                            borderWidth="2px"
                                            borderColor={borderColor}
                                        />
                                        <VStack align="start" spacing={3} flex={1}>
                                            <VStack align="start" spacing={1}>
                                                <Heading size="lg" color={textColorPrimary}>
                                                    {supplier.company_name}
                                                </Heading>
                                                <Badge
                                                    px={3}
                                                    py={1}
                                                    borderRadius="full"
                                                    colorScheme={getStatusColor(supplier.status)}
                                                    fontSize="sm"
                                                >
                                                    {getStatusText(supplier.status)}
                                                </Badge>
                                            </VStack>

                                            <Stack
                                                direction={{ base: "column", md: "row" }}
                                                spacing={6}
                                                divider={<StackDivider />}
                                                width="100%"
                                            >
                                                <VStack align="start" spacing={2} minW="200px">
                                                    <HStack spacing={2}>
                                                        <Icon as={FiPhone} color="blue.500" />
                                                        <Text fontWeight="medium" fontSize="sm" color={textColorSecondary}>
                                                            Số điện thoại
                                                        </Text>
                                                    </HStack>
                                                    <Text color={textColorPrimary} fontWeight="medium">
                                                        {supplier.contact_phone}
                                                    </Text>
                                                </VStack>

                                                <VStack align="start" spacing={2} minW="200px">
                                                    <HStack spacing={2}>
                                                        <Icon as={FiShield} color="blue.500" />
                                                        <Text fontWeight="medium" fontSize="sm" color={textColorSecondary}>
                                                            Mã số thuế
                                                        </Text>
                                                    </HStack>
                                                    <Text color={textColorPrimary} fontWeight="medium" fontFamily="mono">
                                                        {supplier.tax_id}
                                                    </Text>
                                                </VStack>

                                                <VStack align="start" spacing={2} flex={1}>
                                                    <HStack spacing={2}>
                                                        <Icon as={FiMapPin} color="blue.500" />
                                                        <Text fontWeight="medium" fontSize="sm" color={textColorSecondary}>
                                                            Địa chỉ kinh doanh
                                                        </Text>
                                                    </HStack>
                                                    <Text color={textColorPrimary} noOfLines={2}>
                                                        {supplier.business_address}
                                                    </Text>
                                                </VStack>
                                            </Stack>
                                        </VStack>
                                    </HStack>
                                </CardBody>
                            </Card>

                            {/* Thông tin thời gian */}
                            <Card>
                                <CardBody>
                                    <Heading size="md" mb={4} color={textColorPrimary}>
                                        <Icon as={FiCalendar} mr={2} />
                                        Thông tin thời gian
                                    </Heading>
                                    <SimpleGrid columns={{ base: 1, md: 2 }} spacing={4}>
                                        <VStack align="start" spacing={2}>
                                            <Text fontWeight="medium" color={textColorSecondary} fontSize="sm">
                                                Ngày đăng ký
                                            </Text>
                                            <Text color={textColorPrimary}>
                                                {formatDate(supplier.created_at)}
                                            </Text>
                                        </VStack>
                                        <VStack align="start" spacing={2}>
                                            <Text fontWeight="medium" color={textColorSecondary} fontSize="sm">
                                                Cập nhật lần cuối
                                            </Text>
                                            <Text color={textColorPrimary}>
                                                {formatDate(supplier.updated_at)}
                                            </Text>
                                        </VStack>
                                    </SimpleGrid>
                                </CardBody>
                            </Card>

                            {/* Tài liệu đăng ký */}
                            <Card>
                                <CardBody>
                                    <Heading size="md" mb={4} color={textColorPrimary}>
                                        <Icon as={FiFileText} mr={2} />
                                        Tài liệu đăng ký
                                    </Heading>

                                    {supplier.documents && supplier.documents.length > 0 ? (
                                        <VStack spacing={4} align="stretch">
                                            {supplier.documents.map((doc, index) => (
                                                <Box key={index}>
                                                    <HStack justify="space-between" mb={3}>
                                                        <HStack spacing={3}>
                                                            <Icon
                                                                as={getDocumentStatusIcon(doc.verification_status)}
                                                                color={`${getDocumentStatusColor(doc.verification_status)}.500`}
                                                            />
                                                            <VStack align="start" spacing={0}>
                                                                <Text fontWeight="medium" color={textColorPrimary}>
                                                                    Tài liệu đăng ký #{index + 1}
                                                                </Text>
                                                                <Text fontSize="sm" color={textColorSecondary}>
                                                                    Cập nhật: {formatDate(doc.updated_at)}
                                                                </Text>
                                                            </VStack>
                                                        </HStack>
                                                        <HStack spacing={2}>
                                                            <Badge
                                                                colorScheme={getDocumentStatusColor(doc.verification_status)}
                                                                px={3}
                                                                py={1}
                                                                borderRadius="md"
                                                            >
                                                                {getDocumentStatusText(doc.verification_status)}
                                                            </Badge>

                                                            {/* Approve/Reject buttons for pending documents */}
                                                            {doc.verification_status === 'pending' && (
                                                                <>
                                                                    <Button
                                                                        size="sm"
                                                                        colorScheme="green"
                                                                        variant="solid"
                                                                        leftIcon={<FiCheck />}
                                                                        onClick={() => handleDocumentAction(doc.id, 'approve')}
                                                                        isLoading={actionLoading === `${doc.id}-approve`}
                                                                        loadingText="Đang duyệt..."
                                                                    >
                                                                        Duyệt
                                                                    </Button>
                                                                    <Button
                                                                        size="sm"
                                                                        colorScheme="red"
                                                                        variant="solid"
                                                                        leftIcon={<FiX />}
                                                                        onClick={() => handleDocumentAction(doc.id, 'reject')}
                                                                        isLoading={actionLoading === `${doc.id}-reject`}
                                                                        loadingText="Đang từ chối..."
                                                                    >
                                                                        Từ chối
                                                                    </Button>
                                                                </>
                                                            )}
                                                        </HStack>
                                                    </HStack>

                                                    {doc.admin_note && (
                                                        <Alert status="info" size="sm" mb={3} borderRadius="md">
                                                            <AlertIcon boxSize="16px" />
                                                            <Text fontSize="sm">{doc.admin_note}</Text>
                                                        </Alert>
                                                    )}

                                                    {doc.document && (
                                                        <SimpleGrid columns={{ base: 1, md: 2, lg: 4 }} spacing={3}>
                                                            {Object.entries(doc.document).map(([docType, url]) => (
                                                                <Box
                                                                    key={docType}
                                                                    borderWidth="1px"
                                                                    borderColor={borderColor}
                                                                    borderRadius="md"
                                                                    overflow="hidden"
                                                                    _hover={{
                                                                        transform: 'scale(1.02)',
                                                                        boxShadow: 'md'
                                                                    }}
                                                                    transition="all 0.2s"
                                                                    cursor="pointer"
                                                                    onClick={() => openDocumentInNewTab(url)}
                                                                >
                                                                    <VStack spacing={0}>
                                                                        <Box
                                                                            width="100%"
                                                                            height="120px"
                                                                            bg={useColorModeValue('gray.100', 'gray.700')}
                                                                            display="flex"
                                                                            alignItems="center"
                                                                            justifyContent="center"
                                                                            position="relative"
                                                                        >
                                                                            <img
                                                                                src={url}
                                                                                alt={documentTypeNames[docType] || docType}
                                                                                style={{
                                                                                    width: '100%',
                                                                                    height: '100%',
                                                                                    objectFit: 'cover'
                                                                                }}
                                                                                onError={(e) => {
                                                                                    e.target.style.display = 'none';
                                                                                    e.target.nextSibling.style.display = 'flex';
                                                                                }}
                                                                            />
                                                                            <VStack
                                                                                spacing={2}
                                                                                style={{ display: 'none' }}
                                                                                position="absolute"
                                                                                top="50%"
                                                                                left="50%"
                                                                                transform="translate(-50%, -50%)"
                                                                            >
                                                                                <Icon as={FiFileText} color="gray.400" boxSize={8} />
                                                                                <Text fontSize="xs" color="gray.500">
                                                                                    Không thể tải ảnh
                                                                                </Text>
                                                                            </VStack>
                                                                        </Box>
                                                                        <Box p={3} width="100%">
                                                                            <Text
                                                                                fontSize="sm"
                                                                                fontWeight="medium"
                                                                                textAlign="center"
                                                                                color={textColorPrimary}
                                                                                noOfLines={1}
                                                                            >
                                                                                {documentTypeNames[docType] || docType}
                                                                            </Text>
                                                                        </Box>
                                                                    </VStack>
                                                                </Box>
                                                            ))}
                                                        </SimpleGrid>
                                                    )}

                                                    {index < supplier.documents.length - 1 && (
                                                        <Divider mt={4} />
                                                    )}
                                                </Box>
                                            ))}
                                        </VStack>
                                    ) : (
                                        <Alert status="warning" borderRadius="md">
                                            <AlertIcon />
                                            <Text>Chưa có tài liệu đăng ký nào</Text>
                                        </Alert>
                                    )}
                                </CardBody>
                            </Card>
                        </VStack>
                    ) : null}
                </ModalBody>
            </ModalContent>
        </Modal>
    );
};

export default SupplierDetailModal;