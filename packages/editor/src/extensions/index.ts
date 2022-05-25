import Table from '@tiptap/extension-table';
import StarterKit from '@tiptap/starter-kit';
import QuickEdit from './quick-edit';
import Highlight from '@tiptap/extension-highlight';
import Typography from '@tiptap/extension-typography';
import Document from '@tiptap/extension-document';
import TableRow from '@tiptap/extension-table-row';
import TableCell from '@tiptap/extension-table-cell';
import TableHeader from '@tiptap/extension-table-header';
import TextAlign from '@tiptap/extension-text-align';
import TaskList from '@tiptap/extension-task-list';
import TaskItem from '@tiptap/extension-task-item';
import Placeholder from '@tiptap/extension-placeholder';
import { Iframe } from './iframe';
import { Link } from './link';
// import { ExtendedImage } from './image';
import { CodeBlock } from './codeblock';
import { TetuaEditorProps } from '..';
import { getImageExtension } from './image';
import { Extensions } from '@tiptap/core';

export function getExtensions(props: TetuaEditorProps) {
  const extensions: Extensions =  [
    Highlight,
    Typography,
    getImageExtension({
      uploadHandler: props.uploadHandler,
    }),
    Iframe,
    TableRow,
    TableCell,
    TableHeader,
    TextAlign,
    TaskList,
    CodeBlock,
    QuickEdit.configure(),
    TaskItem.configure({
      nested: true,
    }),
    Link.configure({
      autolink: false,
      openOnClick: false,
      linkOnPaste: false,
    }),
    StarterKit.configure({
      codeBlock: false,
      document: false,
    }),
    Table.configure({
      resizable: true,
    }),
    Placeholder.configure({
      placeholder: ({ node }) => {
        if (!props.disableTitle) {
          if (node.type.name === 'heading') {
            if (node.attrs.level == 1) {
              return 'Enter the name';
            }

            return 'Heading content';
          }
        }

        return 'Start typing something...';
      },
    })
  ];

  if (!props.disableTitle) {
    extensions.push(Document.extend({
      content: 'heading block*',
    }));
  } else {
    extensions.push(Document);
  }

  return extensions;
}

