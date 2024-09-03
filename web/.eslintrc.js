module.exports = {
  root: true,
  extends: [
    require.resolve('@vercel/style-guide/eslint/browser'),
    require.resolve('@vercel/style-guide/eslint/react'),
    require.resolve('@vercel/style-guide/eslint/next'),
    require.resolve('@vercel/style-guide/eslint/typescript'),
    'next/core-web-vitals',
    'prettier',
  ],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    project: true,
    tsconfigRootDir: __dirname,
  },
  settings: {
    'import/resolver': {
      typescript: {
        project: true,
        tsconfigRootDir: __dirname,
      },
    },
  },
  rules: {
    '@typescript-eslint/explicit-function-return-type': [
      'warn',
      { allowFunctionsWithoutTypeParameters: true },
    ],
    // react-hook-form handleSubmit() bypass
    '@typescript-eslint/no-misused-promises': [
      2,
      {
        checksVoidReturn: {
          attributes: false,
        },
      },
    ],
  },
  overrides: [
    {
      extends: ['plugin:@typescript-eslint/disable-type-checked'],
      files: ['./**/*.js'],
    },
    {
      files: [
        'src/app/**/{default,error,instrumentation,layout,loading,middleware,not-found,page,route,template}.tsx',
      ],
      rules: {
        'import/no-default-export': 'off',
      },
    },
  ],
};
