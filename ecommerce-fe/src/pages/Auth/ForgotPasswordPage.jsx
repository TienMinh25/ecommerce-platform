import { Box, Container } from '@chakra-ui/react';
import ForgotPasswordForm from "../../components/auth/ForgotPasswordForm.jsx";
import PageTitle from "../PageTitle.jsx";

const ForgotPasswordPage = () => {
    return (
        <>
            <PageTitle title="Quên mật khẩu" />
            <Container maxW="container.lg" py={16}>
                <Box maxW="lg" mx="auto">
                    <ForgotPasswordForm />
                </Box>
            </Container>
        </>
    );
};

export default ForgotPasswordPage;