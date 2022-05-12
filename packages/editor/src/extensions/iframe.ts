import { Node } from '@tiptap/core'
import { createNodeViewBlock } from '../utils'

export interface IframeOptions {
  allowFullscreen: boolean,
  HTMLAttributes: {
    [key: string]: any
  },
}

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    iframe: {
      setIframe: (options: { src: string, width?: number, height?: number }) => ReturnType,
    }
  }
}

export const Iframe = Node.create<IframeOptions>({
  name: 'iframe',
  group: 'block',
  atom: true,
  addOptions() {
    return {
      allowFullscreen: true,
      HTMLAttributes: {
        class: 'block-iframe',
      },
    }
  },
  addAttributes() {
    return {
      src: {
        default: null,
      },
      frameborder: {
        default: 0,
      },
      allowfullscreen: {
        default: this.options.allowFullscreen,
        parseHTML: () => this.options.allowFullscreen,
      },
      width: {
        default: 500,
      },
      height: {
        default: 315,
      },
      title: {
        default: null,
      },
      allow: {
        default: 'accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture',
      }
    }
  },
  parseHTML() {
    return [{
      tag: 'iframe',
    }]
  },
  renderHTML({ HTMLAttributes }) {
    return ['div', this.options.HTMLAttributes, ['iframe', HTMLAttributes]]
  },
  addCommands() {
    return {
      setIframe: (options: { src: string, width?: number, height?: number }) => ({ tr, dispatch }) => {
        const { selection } = tr
        const node = this.type.create(options)

        if (dispatch) {
          tr.replaceRangeWith(selection.from, selection.to, node)
        }

        return true
      },
    }
  },

  addNodeView() {
    return ({ editor, node: _node, getPos, HTMLAttributes: attrs, decorations: _decorations, extension: _extension }) => {
      const createImageUrlElm = () => {
        const iframeSrcContainer = document.createElement('div');
        const iframeSrcInput = document.createElement('input');
        const iframeSrcApplyBtn = document.createElement('button');

        iframeSrcContainer.className = 'mely-editor-iframe-src';
        iframeSrcApplyBtn.innerText = 'Insert';
        iframeSrcInput.setAttribute('type', 'text');
        iframeSrcInput.setAttribute('placeholder', 'Enter embed URL');
        iframeSrcApplyBtn.addEventListener('click', () => {
          if (typeof getPos === 'function') {
            let src = iframeSrcInput.value;
            if (src.startsWith('https://www.youtube.com/watch?v=')) {
              src = src.replace('https://www.youtube.com/watch?v=', 'https://www.youtube.com/embed/');
              src = src.split('&')[0];
            }
            editor.view.dispatch(editor.view.state.tr.setNodeMarkup(getPos(), undefined, {
              src: src,
            }))
            editor.commands.focus();
          }
        });

        iframeSrcContainer.append(iframeSrcInput, iframeSrcApplyBtn);
        setTimeout(() => {
          iframeSrcInput.focus();
        }, 0);
        return iframeSrcContainer;
      }
      const contentDomElm = document.createElement('iframe');
      contentDomElm.setAttribute('src', attrs.src);
      contentDomElm.setAttribute('width', attrs.width);
      contentDomElm.setAttribute('height', attrs.height);
      contentDomElm.setAttribute('title', attrs.title);
      contentDomElm.setAttribute('allow', attrs.allow);
      // contentDomElm.setAttribute('alt', attrs.alt);
      // contentDomElm.setAttribute('title', attrs.title);

      const viewDomElms: HTMLElement[] = [];
      const imageUrlElm = createImageUrlElm();

      viewDomElms.push(imageUrlElm);
      const { dom } = createNodeViewBlock(contentDomElm, viewDomElms);
      dom.classList.add('block-iframe');

      return {
        dom,
        contentDOM: contentDomElm,
        stopEvent: () => !attrs.src,
        ignoreMutation: _mutation => !attrs.src,
      }
    }
  },
})