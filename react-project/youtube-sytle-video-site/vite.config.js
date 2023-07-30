import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

export default defineConfig(() => {
    return {
        optimizeDeps: {
            force: true // 强制进行依赖预构建
        },
        server: {
            host: true // 监听所有地址
        },
        build: {
            outDir: 'build' // 打包文件的输出目录
        },
        plugins: [react()],
    };
});
