import { Box, Container } from '@chakra-ui/react';
import ResetPasswordForm from "../../components/auth/ResetPasswordForm.jsx";
import PageTitle from "../PageTitle.jsx";

const ResetPasswordPage = () => {
    return (
        <>
            <PageTitle title="Đặt lại mật khẩu" />
            <Container maxW="container.lg" py={16}>
                <Box maxW="lg" mx="auto">
                    <ResetPasswordForm />
                </Box>
            </Container>
        </>
    );
};

export default ResetPasswordPage;