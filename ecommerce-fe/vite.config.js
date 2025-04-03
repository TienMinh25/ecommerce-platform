import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0', // Lắng nghe trên tất cả các địa chỉ IP, có thể truy cập từ các thiết bị khác trong mạng
    // host: 'localhost', // Chỉ lắng nghe trên localhost, mặc định
    // host: '192.168.1.100', // Lắng nghe trên một địa chỉ IP cụ thể
    port: 5173, // Port mà server sẽ lắng nghe
    open: true, // Tự động mở browser khi khởi động server
    allowedHosts: ['client.local']
  },
});
