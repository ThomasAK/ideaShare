import eslintConfigStandard from 'eslint-config-standard'
import standardts from 'eslint-config-standard-with-typescript'
import globals from 'globals'
import importPlugin from 'eslint-plugin-import'
import n from 'eslint-plugin-n'
import promise from 'eslint-plugin-promise'
import tslint from '@typescript-eslint/eslint-plugin'
import tsparser from '@typescript-eslint/parser'
delete standardts.extends

function toLanguageOptions (conf) {
  conf.languageOptions = {
    parserOptions: conf.parserOptions
  }
  delete conf.parser
  delete conf.plugins
  delete conf.globals
  delete conf.env
  delete conf.parserOptions
  // delete conf.languageOptions?.parserOptions?.project
}

toLanguageOptions(eslintConfigStandard)
toLanguageOptions(standardts)

standardts.languageOptions.parser = {
  ...tsparser,
  parse: (...args) => {
    debugger
    return tsparser.parse(...args)
  },
  parseForESLint: (...args) => {
    debugger
    return tsparser.parseForESLint(...args)
  }
}

export default [
  {
    ignores: ['**/dist/**', '**/**js']
  },
  {
    files: [
      '**/*.ts',
      '**/*.tsx',
      '**/*.jsx',
    ],
  },
  eslintConfigStandard,
  standardts,
  {
    plugins: {
      import: importPlugin,
      n,
      promise,
      '@typescript-eslint': tslint
    },
    languageOptions: {
      globals: {
        ...globals.es2021,
        ...globals.node,
        ...globals.browser
      }
    }
  }
]
