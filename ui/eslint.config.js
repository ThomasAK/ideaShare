import tsStandardEslintFlatConfig from 'ts-standard-eslint-flat-config'
import globals from 'globals'

let [files, standard, standardts, extra ] = tsStandardEslintFlatConfig

delete standardts.rules['@typescript-eslint/strict-boolean-expressions']
export default [
  {
    ignores: ['**/dist/**', '**/**js']
  },
  files,
  standard,
  standardts,
  extra,
  {languageOptions: {
      globals: {
        ...globals.browser
      }
    }
  }
]
