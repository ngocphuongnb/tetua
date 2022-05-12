import { Editor } from '@tiptap/core';
// import high from 'highlight.js/lib/core';
import { TetuaEditor } from '.';

type HTMLElementEvent<T extends HTMLElement> = Event & {
  target: T;
}

declare class MarkdownEditorType extends Editor {
  getMarkdown(): string;
  parseMarkdown(input: string): string;
}

declare global {
  interface Window {
    editor: TetuaEditor;
    hljs: any,
    TetuaEditor: typeof TetuaEditor;
  }
}

declare module '@tiptap/core' {
  interface EditorOptions {
      markdown: any;
  }
}
