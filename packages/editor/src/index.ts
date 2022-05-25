import { Editor as TipTapEditor } from '@tiptap/core';
import { getExtensions } from './extensions';
import { ImageExtensionProps } from './extensions/image';
import { MarkdownEditorType } from './global';
import { Menu } from './menu';

const createMarkdownEditor = require('tiptap-markdown').createMarkdownEditor;
const MarkdownEditor: typeof TipTapEditor = createMarkdownEditor(TipTapEditor);

export interface TetuaEditorProps extends ImageExtensionProps {
  commentMode?: boolean;
  disableTitle?: boolean;
}

export class TetuaEditor {
  public tiptapEditor: TipTapEditor;
  public textareaElement: HTMLTextAreaElement;
  public tiptapEditorElement: HTMLDivElement;
  public editorElement: HTMLDivElement;
  public editorContentElement: HTMLDivElement;
  public useMarkdown = false;
  public menu: Menu = null;
  public inititalized: boolean = false;
  public props: TetuaEditorProps;

  constructor(selector: string, props?: TetuaEditorProps) {
    this.textareaElement = document.querySelector<HTMLTextAreaElement>(selector);
    this.props = props || {};

    if (!this.textareaElement) {
      throw new Error('Element not found');
    }

    if (this.textareaElement.tagName !== 'TEXTAREA') {
      throw new Error('Element must be a textarea');
    }

    this.setTextareaRows();
    this.menu = new Menu(this);
    this.createEditorElement();
    this.createTiptapEditorElement();
    this.menu.initialize();
  }
  
  private setTextareaRows() {
    const rows = this.textareaElement.value.split('\n').length;
    this.textareaElement.rows = (rows > 100 ? rows : 100) *1.2;
  }

  private createEditorElement() {
    this.editorElement = document.createElement('div');
    this.editorElement.className = 'mely-editor';
    this.insertElementAfter(this.textareaElement, this.editorElement);

    this.editorContentElement = document.createElement('div');
    this.editorContentElement.className = 'mely-editor-content';
    this.editorContentElement.appendChild(this.textareaElement);
    // this.editorContentElement.setAttribute('spellcheck', 'false');
    this.editorElement.appendChild(this.editorContentElement);


    document.addEventListener('scroll', () => {
      const menuHeight = this.menu.menuContainerElement.getBoundingClientRect().height;
      const contentTop = window.scrollY + this.editorContentElement.offsetTop;

      if (contentTop > menuHeight + 10) {
        this.menu.menuContainerElement.classList.add('mely-menu-fixed');
      } else {
        this.menu.menuContainerElement.classList.remove('mely-menu-fixed');
      }
    });
  }

  private createTiptapEditorElement() {
    this.tiptapEditorElement = document.createElement('div');
    this.tiptapEditorElement.className = 'mely-tiptap-editor';
    this.insertElementAfter(this.textareaElement, this.tiptapEditorElement);
    this.textareaElement.style.display = 'none';
    this.initialize();
  }

  private insertElementAfter(referenceNode: Node, newNode: Node) {
    referenceNode.parentNode.insertBefore(newNode, referenceNode.nextSibling);
  }

  public initialize() {
    this.tiptapEditor = new MarkdownEditor({
      markdown: {
        html: true,
        breaks: true,
      },
      element: this.tiptapEditorElement,
      extensions: getExtensions(this.props),
      content: this.textareaElement.value,
      onUpdate: ({ editor }) => {
        // const contentElmBound = this.tiptapEditorElement.getBoundingClientRect();
        if (!this.useMarkdown) {
          this.textareaElement.value = (editor as MarkdownEditorType).getMarkdown();
        }
      },
      onCreate: () => {
        this.inititalized = true;
        this.resize();
      },
      onSelectionUpdate: this.menu.update.bind(this.menu),
      onFocus: () => {
        this.editorElement.classList.add('focused');
      },
      onBlur: () => {
        this.editorElement.classList.remove('focused');
      },
    });

    const updateTiptapContent = () => {
      this.setTextareaRows();
      if (this.useMarkdown) {
        this.tiptapEditor.commands.setContent(this.textareaElement.value);
      }
    }

    this.textareaElement.addEventListener('change', updateTiptapContent);
    this.textareaElement.addEventListener('keyup', updateTiptapContent);
    window.addEventListener('resize', this.resize.bind(this));
  }

  public resize() {
    const editorWidth = this.editorElement.clientWidth;
    this.editorElement.setAttribute('data-width', editorWidth.toString());

    if (editorWidth < 620) {
      this.editorElement.classList.add('mely-editor-small');
      this.menu.menuElement.appendChild(this.menu.menuSwitcherElement);
    } else {
      this.editorElement.classList.remove('mely-editor-small');
      this.menu.menuContainerElement.appendChild(this.menu.menuSwitcherElement);
    }
  }

  public destroy() {
    this.tiptapEditor.destroy();
  }
}

window.TetuaEditor = TetuaEditor