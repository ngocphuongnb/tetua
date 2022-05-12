import CodeBlockLowlight from './codeblock-lowlight';
import { HTMLElementEvent } from '../global';
import * as lowlight from 'lowlight/lib/core';

export const CodeBlock = CodeBlockLowlight.extend({
  addNodeView() {
    return ({ editor, node, getPos, HTMLAttributes: _HTMLAttributes, decorations: _decorations, extension }) => {
      const dom = document.createElement('div');
      dom.className = 'code-block';
      const languageSelectElm = document.createElement<'select'>('select');
      languageSelectElm.contentEditable = 'false';
      languageSelectElm.setAttribute('defaultValue', node.attrs.language);
      languageSelectElm.addEventListener('change', (e: HTMLElementEvent<HTMLSelectElement>) => {
        if (typeof getPos === 'function') {
          editor.view.dispatch(editor.view.state.tr.setNodeMarkup(
            getPos(),
            null,
            { language: e.target.value }
          ));
        }
      });
      languageSelectElm.innerHTML = `<option value="null">auto</option><option disabled>—</option>`;
      extension.options.lowlight.listLanguages().map((lang: string) => {
        const option = document.createElement<'option'>('option');
        option.value = lang;
        option.innerText = lang;

        if (node.attrs.language === lang) {
          option.selected = true;
        }

        languageSelectElm.appendChild(option);
      });

      dom.appendChild(languageSelectElm);
      const pre = document.createElement<'pre'>('pre');
      pre.setAttribute('as', 'code');
      dom.appendChild(pre);

      return {
        dom,
        contentDOM: pre,
      }
    }
  },
}).configure({ lowlight: (lowlight as any).lowlight })