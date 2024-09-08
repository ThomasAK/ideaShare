import tsStandardEslintFlatConfig from 'ts-standard-eslint-flat-config'
import globals from 'globals'
export default [
  {
    ignores: ['**/dist/**', '**/**js']
  },
  ...tsStandardEslintFlatConfig,
  {languageOptions: {
      globals: {
        ...globals.browser
      }
    }
  }
]
