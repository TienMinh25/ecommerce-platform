import {
    Box,
    Button,
    FormControl,
    FormLabel,
    FormErrorMessage,
    Input,
    Select,
    Grid,
    GridItem,
    HStack,
    Text,
    Alert,
    AlertIcon,
    AlertTitle,
    AlertDescription,
    Divider,
    Icon,
    useDisclosure,
    VStack
} from '@chakra-ui/react';
import { FaMapMarkerAlt, FaExternalLinkAlt } from 'react-icons/fa';
import AddressSelectionModal from './AddressSelectionModal';

const BusinessAddressSection = ({
                                    formData,
                                    errors,
                                    selectedAddress,
                                    onInputChange,
                                    onSelectAddress,
                                    onClearAddress
                                }) => {
    const { isOpen: isAddressModalOpen, onOpen: onAddressModalOpen, onClose: onAddressModalClose } = useDisclosure();

    const handleSelectAddress = (address) => {
        onSelectAddress(address);
        onAddressModalClose();
    };

    return (
        <>
            <Divider />

            <Box>
                <HStack justify="space-between" mb={4}>
                    <Text fontWeight="bold">Địa chỉ kinh doanh *</Text>
                    <HStack spacing={2}>
                        <Button
                            size="sm"
                            colorScheme="blue"
                            variant="outline"
                            onClick={onAddressModalOpen}
                            leftIcon={<Icon as={FaMapMarkerAlt} />}
                        >
                            Chọn từ địa chỉ có sẵn
                        </Button>
                        <Button
                            size="sm"
                            colorScheme="green"
                            variant="outline"
                            onClick={() => window.open('/user/account/addresses', '_blank')}
                            leftIcon={<Icon as={FaExternalLinkAlt} />}
                        >
                            Quản lý địa chỉ
                        </Button>
                        {selectedAddress && (
                            <Button
                                size="sm"
                                colorScheme="red"
                                variant="ghost"
                                onClick={onClearAddress}
                            >
                                Xóa
                            </Button>
                        )}
                    </HStack>
                </HStack>

                {selectedAddress && (
                    <Alert status="success" mb={4} borderRadius="md">
                        <AlertIcon />
                        <Box flex={1}>
                            <AlertTitle fontSize="sm">Đã chọn địa chỉ:</AlertTitle>
                            <AlertDescription fontSize="sm">
                                <Text fontWeight="medium">
                                    {selectedAddress.recipient_name} - {selectedAddress.phone}
                                </Text>
                                <Text>
                                    {selectedAddress.street}, {selectedAddress.ward}, {selectedAddress.district}, {selectedAddress.province}
                                </Text>
                            </AlertDescription>
                        </Box>
                    </Alert>
                )}

                <Alert status="info" mb={4} borderRadius="md">
                    <AlertIcon />
                    <Box>
                        <AlertTitle fontSize="sm">Lưu ý quan trọng:</AlertTitle>
                        <AlertDescription fontSize="sm">
                            <VStack align="flex-start" spacing={1}>
                                <Text>• Bạn chỉ có thể chọn từ các địa chỉ đã lưu trong tài khoản</Text>
                                <Text>• Nếu chưa có địa chỉ phù hợp, vui lòng tạo địa chỉ mới trước</Text>
                                <Text>• Địa chỉ này sẽ được dùng làm địa chỉ kinh doanh chính</Text>
                            </VStack>
                        </AlertDescription>
                    </Box>
                </Alert>
            </Box>

            <AddressSelectionModal
                isOpen={isAddressModalOpen}
                onClose={onAddressModalClose}
                onSelectAddress={handleSelectAddress}
                selectedAddressId={selectedAddress?.id}
            />
        </>
    );
};

export default BusinessAddressSection;