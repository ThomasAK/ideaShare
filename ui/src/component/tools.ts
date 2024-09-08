import Table from '@editorjs/table'
import List from '@editorjs/list'
// @ts-expect-error not a ts module
import Code from '@editorjs/code'
// @ts-expect-error not a ts module
import LinkTool from '@editorjs/link'
import Header from '@editorjs/header'
import Quote from '@editorjs/quote'
import Delimiter from '@editorjs/delimiter'

export const EditorJSTools = {
  table: Table,
  list: List,
  code: Code,
  linkTool: LinkTool,
  header: Header,
  quote: Quote,
  delimiter: Delimiter
}
