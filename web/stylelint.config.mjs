export default {
  extends: ['stylelint-config-standard', 'stylelint-config-recommended-vue'],
  ignoreFiles: ['dist/**', 'node_modules/**'],
  rules: {
    'color-no-hex': true,
    'color-named': 'never',
    'custom-property-pattern': '^(ref|sys|cmp|state|nav|dropdown)-[a-z0-9-]+$',
    'selector-class-pattern': null,
    'keyframes-name-pattern': null,
    'no-descending-specificity': null,
    'declaration-block-no-duplicate-properties': null,
    'value-keyword-case': null,
    'media-feature-range-notation': null,
    'no-empty-source': null,
    'color-function-notation': null,
    'color-function-alias-notation': null,
    'alpha-value-notation': null,
  },
  overrides: [
    {
      files: ['src/assets/styles/tokens.ref.css'],
      rules: {
        'color-no-hex': null,
      },
    },
  ],
}
