import { Editor } from '@tiptap/core';
import Image from '@tiptap/extension-image';
import { createNodeViewBlock } from '../utils';

export type ImageUploadHandler = (file: File, callback: (url: string, err?: Error) => void) => void;
export interface ImageExtensionProps {
  uploadHandler?: ImageUploadHandler;
  disableTitle?: boolean;
}

const createImageUrlElm = (editor: Editor, getPos: boolean | (() => number)) => {
  const imageUrlContainer = document.createElement('div');
  const imageUrlInput = document.createElement('input');
  const imageUrlApplyBtn = document.createElement('button');

  imageUrlContainer.className = 'mely-editor-img-url';
  imageUrlApplyBtn.innerText = 'Insert';
  imageUrlInput.setAttribute('type', 'text');
  imageUrlInput.setAttribute('placeholder', 'Enter image URL');
  imageUrlApplyBtn.addEventListener('click', (e: KeyboardEvent) => {
    e.preventDefault();
    e.stopImmediatePropagation();

    if (!imageUrlInput.value) {
      imageUrlInput.focus();
      return;
    }

    if (typeof getPos === 'function') {
      editor.view.dispatch(editor.view.state.tr.setNodeMarkup(getPos(), undefined, {
        src: imageUrlInput.value,
        alt: '',
        title: '',
      }))
      editor.commands.focus();
    }
  });

  imageUrlContainer.append(imageUrlInput, imageUrlApplyBtn);
  setTimeout(() => imageUrlInput.focus(), 0);

  return imageUrlContainer;
}

const createOrTextElm = () => {
  const orText = document.createElement('p');
  orText.innerText = 'or';
  return orText;
}

const createImageUploadElm = (dom: HTMLDivElement, uploadHandler: ImageUploadHandler, editor: Editor, getPos: boolean | (() => number)) => {
  const uploadElm = document.createElement('div');
  const uploadInput = document.createElement('input');
  const uploadBtn = document.createElement('button');

  uploadElm.className = 'mely-editor-img-upload';
  uploadInput.setAttribute('type', 'file');
  uploadInput.setAttribute('accept', 'image/*');
  uploadInput.setAttribute('name', 'file');
  uploadInput.setAttribute('id', 'file');
  uploadInput.setAttribute('style', 'display: none;');
  uploadBtn.setAttribute('type', 'button');
  uploadBtn.setAttribute('class', 'mely-editor-img-upload-btn');
  uploadBtn.innerText = 'Select file';
  uploadBtn.addEventListener('click', (e) => {
    e.preventDefault();
    uploadInput.click();
  });

  uploadInput.addEventListener('change', (e) => {
    const target = e.target as HTMLInputElement;
    const file = target.files[0];
    dom.classList.add('uploading');
    uploadHandler(file, (url, err) => {
      dom.classList.remove('uploading');
      if (err) {
        console.error(err);
        alert('Upload failed');
        return;
      }

      if (typeof getPos === 'function') {
        editor.view.dispatch(editor.view.state.tr.setNodeMarkup(getPos(), undefined, {
          src: url,
          alt: '',
          title: '',
        }))
        editor.commands.focus();
      }
    });
  });

  uploadElm.append(uploadInput, uploadBtn);
  return uploadElm;
}

export const getImageExtension = (props: ImageExtensionProps = {}) => {
  const uploadHandler = props.uploadHandler || (() => console.log('Upload handler not set'));

  return Image.extend({
    onCreate() {
      window.addEventListener('paste', (e: ClipboardEvent) => {
        if (e.clipboardData.files && e.clipboardData.files.length > 0) {
          e.preventDefault();
          e.stopImmediatePropagation();
          e.stopPropagation();
          const file = e.clipboardData.files[0];
          const selection = this.editor.view.state.tr.selection;
          const pos = selection.$anchor.pos - 1;

          if (!props.disableTitle) {
            const position = selection.$from;
            let pastedOnTitleField = false;
            let titleNode = null;
  
            position.doc.nodesBetween(selection.from, selection.to, (node) => {
              titleNode = node;
              pastedOnTitleField = node.type.name === 'heading' && node.attrs.level == 1;
            });

            if (pastedOnTitleField) {
              alert('Can\'t paste image on title field');
              if (titleNode) {
                this.editor.view.dispatch(this.editor.view.state.tr.setNodeMarkup(
                  pos,
                  this.editor.schema.nodes.paragraph
                ));
              }
              return;
            }
          }

          uploadHandler(file, (url, err) => {
            if (err) {
              console.error(err);
              alert('Upload failed');
              this.editor.view.dispatch(this.editor.view.state.tr.setNodeMarkup(
                pos,
                this.editor.schema.nodes.paragraph
              ));
              return;
            }

            this.editor.view.dispatch(this.editor.view.state.tr.setNodeMarkup(pos, undefined, {
              src: url,
              alt: '',
              title: '',
            }));
          });
        }
      });
    },
    addNodeView() {
      return ({
        editor,
        node: _node,
        getPos,
        HTMLAttributes: attrs,
        decorations: _decorations,
        extension: _extension
      }) => {

        if (!props.disableTitle && typeof getPos === 'function') {
          let pastedOnTitleField = false
          const selection = editor.view.state.selection;
          const position = selection.$from;
          const pos = getPos()
          position.doc.nodesBetween(pos, pos, (node) => {
            pastedOnTitleField = node.type.name === 'heading' && node.attrs.level == 1;
          });

          if (pastedOnTitleField) {
            alert('Can\'t insert image on title field');
            return;
          }
        }

        const contentDomElm = document.createElement('img');
        contentDomElm.setAttribute('src', attrs.src);
        contentDomElm.setAttribute('alt', attrs.alt || '');
        contentDomElm.setAttribute('title', attrs.title || '');

        const { dom, view } = createNodeViewBlock(contentDomElm, []);
        dom.classList.add('block-img');
        view.append(
          createImageUrlElm(editor, getPos),
          createOrTextElm(),
          createImageUploadElm(dom, uploadHandler, editor, getPos)
        );

        return {
          dom,
          contentDOM: contentDomElm,
          stopEvent: () => !attrs.src,
          ignoreMutation: _mutation => !attrs.src,
        }
      }
    },
  }).configure({ inline: true });
}