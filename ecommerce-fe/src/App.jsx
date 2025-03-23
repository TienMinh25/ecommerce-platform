import { Box } from '@chakra-ui/react';
import AppRoutes from './routes/AppRoutes';
import PageTitle from './pages/PageTitle';

function App() {
  return (
    <Box minH='100vh'>
      <PageTitle />
      <AppRoutes />
    </Box>
  );
}

export default App;
