import { defineConfig, mergeConfig } from 'vitest/config'
import viteConfig from './vite.config'

export default mergeConfig(
  viteConfig,
  defineConfig({
    test: {
      environment: 'happy-dom',
      globals: true,
      include: ['src/**/*.{test,spec}.{ts,tsx}']
    },
    resolve: {
      alias: {
        // cropperjs@2.1.1 doesn't ship a CSS file; the source imports
        // `cropperjs/dist/cropper.css` which doesn't exist. Stub it
        // to an empty module so the test suite can resolve imports.
        'cropperjs/dist/cropper.css': new URL('./src/__mocks__/empty.css', import.meta.url).pathname
      }
    }
  })
)