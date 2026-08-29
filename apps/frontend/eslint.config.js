// Konfigurasi ESLint flat untuk seluruh workspace frontend.
// Aturan kunci: komponen/hook tidak boleh memanggil `fetch` atau `axios`
// langsung — semua akses jaringan lewat fungsi di `@klontong/api`.
import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import reactHooks from 'eslint-plugin-react-hooks';
import importPlugin from 'eslint-plugin-import';
import globals from 'globals';

export default tseslint.config(
  {
    ignores: [
      '**/dist/**',
      '**/node_modules/**',
      '**/coverage/**',
      '**/*.tsbuildinfo',
      'packages/api/src/schema.d.ts',
    ],
  },
  js.configs.recommended,
  ...tseslint.configs.recommended,
  {
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      globals: { ...globals.browser, ...globals.node },
    },
    plugins: {
      'react-hooks': reactHooks,
      import: importPlugin,
    },
    rules: {
      'react-hooks/rules-of-hooks': 'error',
      'react-hooks/exhaustive-deps': 'warn',
      'import/no-relative-packages': 'error',
      'no-restricted-globals': [
        'error',
        {
          name: 'fetch',
          message: 'Jangan panggil fetch langsung — gunakan fungsi dari @klontong/api.',
        },
      ],
      'no-restricted-imports': [
        'error',
        {
          paths: [
            {
              name: 'axios',
              message: 'Gunakan fungsi dari @klontong/api, bukan axios langsung.',
            },
          ],
        },
      ],
    },
  },
  {
    // Lapis API adalah satu-satunya tempat yang boleh menyentuh transport.
    files: ['packages/api/**/*.{ts,tsx}'],
    rules: {
      'no-restricted-globals': 'off',
      'no-restricted-imports': 'off',
    },
  },
  {
    files: ['**/*.config.{ts,js}'],
    rules: {
      'no-restricted-globals': 'off',
    },
  },
);
