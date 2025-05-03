import React from 'react';
import ReactDOM from 'react-dom/client';
import { ChakraProvider } from '@chakra-ui/react';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import theme from './config/theme';
import './index.css';
import { AuthProvider } from './context/AuthProvider';
import {NotificationProvider} from "./context/NotificationProvider.jsx";

ReactDOM.createRoot(document.getElementById('root')).render(
    <BrowserRouter>
        <ChakraProvider theme={theme}>
            <AuthProvider>
                <NotificationProvider>
                    <App />
                </NotificationProvider>
            </AuthProvider>
        </ChakraProvider>
    </BrowserRouter>
);