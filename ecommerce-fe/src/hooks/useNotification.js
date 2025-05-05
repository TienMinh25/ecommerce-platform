import {useContext} from "react";
import NotificationContext from "../context/NotificationProvider.jsx";

export const useNotification = () => {
    const context = useContext(NotificationContext);

    if (!context) {
        throw new Error('useNotification must be used with a NotificationProvider');
    }
    return context;
};

export default useNotification;