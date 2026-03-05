export default {
  extends: [
    'stylelint-config-standard',
    'stylelint-config-tailwindcss',
    'stylelint-config-recommended-vue',
  ],
  ignoreFiles: ['dist/**', 'node_modules/**'],
  rules: {
    'color-no-hex': true,
    'color-named': 'never',
    'custom-property-pattern': '^(ref|sys|cmp|state|nav|dropdown)-[a-z0-9-]+$',
    'selector-class-pattern': null,
    'keyframes-name-pattern': '^[a-z][a-z0-9-]*$',
    'no-descending-specificity': null,
    'declaration-block-no-duplicate-properties': true,
    'value-keyword-case': null,
    'media-feature-range-notation': 'context',
    'no-empty-source': true,
    'color-function-notation': null,
    'color-function-alias-notation': null,
    'alpha-value-notation': null,
  },
  overrides: [
    {
      files: ['src/assets/main.css'],
      rules: {
        'import-notation': 'string',
      },
    },
    {
      files: ['src/assets/styles/tailwind-theme.css'],
      rules: {
        'custom-property-pattern': null,
      },
    },
    {
      files: ['src/assets/styles/tokens.ref.css'],
      rules: {
        'color-no-hex': null,
      },
    },
  ],
}
