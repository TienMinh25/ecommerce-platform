import React from 'react';
import {
    Box,
    Button,
    Flex,
    Icon,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Text,
    useColorModeValue,
} from '@chakra-ui/react';
import {FiAlertTriangle} from 'react-icons/fi';

/**
 * Reusable delete confirmation modal component
 *
 * @param {boolean} isOpen - Controls whether the modal is displayed
 * @param {function} onClose - Function to call when the modal should close
 * @param {function} onConfirm - Function to call when delete is confirmed
 * @param {string} title - Modal title
 * @param {string} message - Confirmation message to display
 * @param {string} itemName - Name of the item being deleted (optional)
 * @param {boolean} isLoading - Whether the delete operation is in progress
 */
const DeleteConfirmationModal = ({
                                     isOpen,
                                     onClose,
                                     onConfirm,
                                     title = 'Confirm Deletion',
                                     message = 'Are you sure you want to delete this item?',
                                     itemName = '',
                                     isLoading = false,
                                 }) => {
    const borderColor = useColorModeValue('gray.200', 'gray.700');
    const headerBg = useColorModeValue('gray.50', 'gray.900');
    const iconColor = useColorModeValue('red.500', 'red.300');

    return (
        <Modal isOpen={isOpen} onClose={onClose} isCentered size="md">
            <ModalOverlay backdropFilter="blur(3px)" bg="blackAlpha.400" />
            <ModalContent borderRadius="lg" shadow="xl">
                <ModalHeader
                    py={4}
                    bg={headerBg}
                    borderTopRadius="lg"
                    borderBottom="1px solid"
                    borderColor={borderColor}
                >
                    <Flex alignItems="center">
                        <Icon as={FiAlertTriangle} color={iconColor} boxSize={5} mr={2} />
                        <Text fontSize="lg" fontWeight="bold">{title}</Text>
                    </Flex>
                </ModalHeader>
                <ModalCloseButton top={3} right={3} />

                <ModalBody py={6}>
                    <Text mb={2}>{message}</Text>
                    {itemName && (
                        <Box
                            mt={3}
                            p={3}
                            borderRadius="md"
                            borderWidth="1px"
                            borderColor={borderColor}
                            bg={useColorModeValue('gray.50', 'gray.800')}
                        >
                            <Text fontWeight="medium" fontSize="md">{itemName}</Text>
                        </Box>
                    )}
                </ModalBody>

                <ModalFooter
                    borderTop="1px solid"
                    borderColor={borderColor}
                    bg={headerBg}
                    borderBottomRadius="lg"
                >
                    <Button
                        variant="outline"
                        mr={3}
                        onClick={onClose}
                        size="md"
                        isDisabled={isLoading}
                    >
                        Cancel
                    </Button>
                    <Button
                        colorScheme="red"
                        onClick={onConfirm}
                        isLoading={isLoading}
                        loadingText="Deleting..."
                        size="md"
                    >
                        Delete
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default DeleteConfirmationModal;