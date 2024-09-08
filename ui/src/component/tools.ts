// @ts-expect-error not a ts module
import Paragraph from '@editorjs/paragraph'
import Table from '@editorjs/table'
import List from '@editorjs/list'
import Header from '@editorjs/header'
// @ts-expect-error not a ts module
import Delimiter from '@coolbytes/editorjs-delimiter'
// @ts-expect-error not a ts module
import AlignmentTuneTool from 'editorjs-text-alignment-blocktune'

export const EditorJSTools = {
  table: Table,
  list: {
    class: List,
    inlineToolbar: ['link', 'bold']
  },
  header: {
    class: Header,
    tunes: ['alignmentTune'],
    inlineToolbar: ['link', 'bold']
  },
  delimiter: Delimiter,
  paragraph: {
    class: Paragraph,
    tunes: ['alignmentTune'],
    inlineToolbar: ['link', 'bold']
  },
  alignmentTune: {
    class: AlignmentTuneTool,
    config: {
      default: 'center'
    }
  }
}
