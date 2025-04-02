import { Box, Container } from '@chakra-ui/react';
import EmailVerification from '../../components/auth/EmailVerification';

const EmailVerificationPage = () => {
    return (
        <Container maxW="container.lg" py={10}>
            <Box
                w="full"
                maxW="md"
                mx="auto"
                p={8}
                bg="white"
                borderRadius="xl"
                boxShadow="lg"
            >
                <EmailVerification />
            </Box>
        </Container>
    );
};

export default EmailVerificationPage;