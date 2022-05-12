import {
  Editor,
  posToDOMRect,
  isTextSelection,
  isNodeSelection,
} from '@tiptap/core'
import { EditorState, Plugin, PluginKey } from 'prosemirror-state'
import { Node as ProsemirrorNode } from 'prosemirror-model'
import { EditorView } from 'prosemirror-view'
import tippy, { Instance, Props } from 'tippy.js'

export interface QuickEditPluginProps {
  pluginKey: PluginKey | string,
  editor: Editor,
  tippyOptions?: Partial<Props>,
  shouldShow?: ((props: {
    editor: Editor,
    view: EditorView,
    state: EditorState,
    oldState?: EditorState,
    from: number,
    to: number,
  }) => boolean) | null,
}

export type QuickEditViewProps = QuickEditPluginProps & {
  view: EditorView,
}

export class QuickEditView {
  public editor: Editor

  public element: HTMLElement

  public view: EditorView

  public preventHide = false

  public tippy: Instance | undefined

  public tippyOptions?: Partial<Props>

  public shouldShow: Exclude<QuickEditPluginProps['shouldShow'], null> = ({
    view,
    state,
    from,
    to,
  }) => {
    const { doc, selection } = state
    const { empty } = selection

    // Sometime check for `empty` is not enough.
    // Doubleclick an empty paragraph returns a node size of 2.
    // So we check also for an empty text size.
    const isEmptyTextBlock = !doc.textBetween(from, to).length
      && isTextSelection(state.selection)

    if (
      !view.hasFocus()
      || empty
      || isEmptyTextBlock
    ) {
      return false
    }

    return true
  }

  constructor({
    editor,
    view,
    tippyOptions = {},
    shouldShow,
  }: QuickEditViewProps) {
    this.view = view
    this.editor = editor
    this.element = document.createElement('div')
    this.element.classList.add('mely-bubble');

    if (shouldShow) {
      this.shouldShow = shouldShow
    }

    this.element.addEventListener('mousedown', this.mousedownHandler, { capture: true })
    this.view.dom.addEventListener('dragstart', this.dragstartHandler)
    this.editor.on('focus', this.focusHandler)
    this.editor.on('blur', this.blurHandler)
    this.tippyOptions = tippyOptions

    // Detaches menu content from its current parent
    this.element.remove()
    this.element.style.visibility = 'visible'
  }

  mousedownHandler = () => {
    this.preventHide = true
  }

  dragstartHandler = () => {
    this.hide()
  }

  focusHandler = () => {
    // we use `setTimeout` to make sure `selection` is already updated
    setTimeout(() => this.update(this.editor.view))
  }

  blurHandler = ({ event }: { event: FocusEvent }) => {
    if (this.preventHide) {
      this.preventHide = false

      return
    }

    if (
      event?.relatedTarget
      && this.element.parentNode?.contains(event.relatedTarget as Node)
    ) {
      return
    }

    this.hide()
  }

  createTooltip() {
    const { element: editorElement } = this.editor.options
    const editorIsAttached = !!editorElement.parentElement

    if (this.tippy || !editorIsAttached) {
      return
    }

    this.tippy = tippy(editorElement, {
      duration: 0,
      getReferenceClientRect: null,
      content: this.element,
      interactive: true,
      arrow: true,
      trigger: 'manual',
      hideOnClick: 'toggle',
      theme: 'light',
      placement: 'top',
      // placement: 'auto',
      ...this.tippyOptions,
    })

    // maybe we have to hide tippy on its own blur event as well
    if (this.tippy.popper.firstChild) {
      (this.tippy.popper.firstChild as HTMLElement).addEventListener('blur', event => {
        this.blurHandler({ event })
      })
    }
  }

  showTippy(view: EditorView, from: number, to: number) {
    this.tippy?.setProps({
      getReferenceClientRect: () => {
        if (isNodeSelection(view.state.selection)) {
          const node = view.nodeDOM(from) as HTMLElement

          if (node) {
            return node.getBoundingClientRect()
          }
        }

        return posToDOMRect(view, from, to)
      },
    })

    this.show();
  }

  showLinkSettingPopup(view: EditorView) {
    const { state } = view;
    const { selection } = state;
    const { ranges } = selection;
    const linkMarkType = state.schema.marks.link;
    const position = selection.$from;
    const from = Math.min(...ranges.map(range => range.$from.pos));
    const to = Math.max(...ranges.map(range => range.$to.pos));
    let hasLinkMark = false;
    let selectionChildCount = 0;
    let linkMark: ProsemirrorNode = null;

    for (let i = 0; !hasLinkMark && i < ranges.length; i++) {
      let { $from, $to } = ranges[i];
      hasLinkMark = state.doc.rangeHasMark($from.pos, $to.pos, linkMarkType);
    }

    position.doc.nodesBetween(selection.from, selection.to, (node, pos) => {
      if (selection.from === pos) {
        linkMark = node;
      }
      selectionChildCount++;
      return true;
    });

    if (!hasLinkMark || !linkMark || selectionChildCount > 2) {
      return this.hide();
    }

    this.element.innerHTML = '';
    this.element.classList.add('mely-bubble-inline');
    const linkInputElement = document.createElement('input');
    const linkApplyElement = document.createElement('button');
    const linkUnsetElement = document.createElement('button');
    const selectionTo = selection.to;

    const setLink = () => {
      this.editor.chain().setLink({
        href: linkInputElement.value,
      }).run();
      this.editor.chain().focus(selectionTo).run();
      this.hide();
    }

    const unsetLink = () => {
      this.editor.chain().unsetLink().run();
      this.editor.chain().focus(selectionTo).run();
      this.hide();
    }

    const onSubmmitLink = (e: KeyboardEvent) => {
      if (e.key === 'Enter') {
        e.preventDefault();
        e.stopImmediatePropagation();
        setLink();
        return false;
      }
    }


    linkApplyElement.type = 'button';
    linkApplyElement.textContent = 'Apply';
    linkUnsetElement.type = 'button';
    linkUnsetElement.textContent = 'Unset';

    linkUnsetElement.addEventListener('click', unsetLink);
    linkApplyElement.addEventListener('click', setLink);
    linkInputElement.addEventListener('keyup', onSubmmitLink);
    linkInputElement.addEventListener('keydown', onSubmmitLink);
    linkInputElement.addEventListener('keypress', onSubmmitLink);

    this.element.append(linkInputElement, linkApplyElement, linkUnsetElement);
    linkInputElement.value = this.editor.getAttributes('link').href;

    this.showTippy(view, from, to);
    // setTimeout(() => linkInputElement.focus(), 0);
  }

  showImageSettingPopup(view: EditorView) {
    const { state } = view;
    const { selection } = state;
    const { ranges } = selection;
    const from = Math.min(...ranges.map(range => range.$from.pos));
    const to = Math.max(...ranges.map(range => range.$to.pos));

    if (!selection || !isNodeSelection(selection) || selection.empty || selection.node.type.name !== 'image' || !selection.node.attrs.src) {
      return;
    }

    const setImageAttrs = () => {
      this.editor.chain().focus().setImage({
        src: imgHref.value,
        alt: imgInfoAlt.value,
        title: imgInfoTitle.value,
      }).run();
    }

    const onSubmmitImageAttrs = (e: KeyboardEvent) => {
      if (e.key === 'Enter') {
        e.preventDefault();
        e.stopImmediatePropagation();
        setImageAttrs();
        return false;
      }
    }

    const imageAttrs = selection.node.attrs || {};
    this.element.innerHTML = '';

    const imgInfo = document.createElement('div');
    imgInfo.className = 'mely-editor-bubble-content';

    const imgHref = document.createElement('input');
    imgHref.value = imageAttrs.src || '';

    const imgInfoTitle = document.createElement('input');
    imgInfoTitle.setAttribute('placeholder', 'Title');
    imgInfoTitle.value = imageAttrs.title || '';
    imgInfoTitle.addEventListener('keyup', onSubmmitImageAttrs);
    imgInfoTitle.addEventListener('keydown', onSubmmitImageAttrs);
    imgInfoTitle.addEventListener('keypress', onSubmmitImageAttrs);

    const imgInfoAlt = document.createElement('input');
    imgInfoAlt.setAttribute('placeholder', 'Alt');
    imgInfoAlt.value = imageAttrs.alt || '';
    imgInfoAlt.addEventListener('keyup', onSubmmitImageAttrs);
    imgInfoAlt.addEventListener('keydown', onSubmmitImageAttrs);
    imgInfoAlt.addEventListener('keypress', onSubmmitImageAttrs);

    const action = document.createElement('div');
    action.className = 'mely-editor-bubble-inline-actions';

    const imageApplyElement = document.createElement('button');
    const imageUnsetElement = document.createElement('button');
    imageApplyElement.type = 'button';
    imageApplyElement.textContent = 'Apply';
    imageUnsetElement.type = 'button';
    imageUnsetElement.textContent = 'Remove';

    imageApplyElement.addEventListener('click', setImageAttrs);
    imageUnsetElement.addEventListener('click', () => {

      this.editor.chain().focus().deleteRange({ from, to }).run();
    });

    action.append(imageApplyElement, imageUnsetElement);

    imgInfo.append(imgHref, imgInfoTitle, imgInfoAlt, action);

    this.element.append(imgInfo);
    this.showTippy(view, from, to);
  }

  showIframeSettingPopup(view: EditorView) {
    const { state } = view;
    const { selection } = state;
    const { ranges } = selection;
    const from = Math.min(...ranges.map(range => range.$from.pos));
    const to = Math.max(...ranges.map(range => range.$to.pos));

    if (!selection || !isNodeSelection(selection) || selection.empty || selection.node.type.name !== 'iframe' || !selection.node.attrs.src) {
      return;
    }

    const setIframeAttrs = () => {
      this.editor.chain().focus().setIframe({
        src: iframeSrc.value,
        width: parseInt(iframeWidth.value),
        height: parseInt(iframeHeight.value),
      }).run();
    }

    const onSubmmitImageAttrs = (e: KeyboardEvent) => {
      if (e.key === 'Enter') {
        e.preventDefault();
        e.stopImmediatePropagation();
        setIframeAttrs();
        return false;
      }
    }

    const iframeAttrs = selection.node.attrs || {};
    this.element.innerHTML = '';

    const iframeInfo = document.createElement('div');
    iframeInfo.className = 'mely-editor-bubble-content';

    const iframeSrc = document.createElement('input');
    iframeSrc.value = iframeAttrs.src || '';

    const iframeWidth = document.createElement('input');
    iframeWidth.setAttribute('placeholder', 'Width');
    iframeWidth.value = iframeAttrs.width || '';
    iframeWidth.type = 'number';
    iframeWidth.addEventListener('keyup', onSubmmitImageAttrs);
    iframeWidth.addEventListener('keydown', onSubmmitImageAttrs);
    iframeWidth.addEventListener('keypress', onSubmmitImageAttrs);

    const iframeHeight = document.createElement('input');
    iframeHeight.setAttribute('placeholder', 'Height');
    iframeHeight.value = iframeAttrs.height || '';
    iframeHeight.type = 'number';
    iframeHeight.addEventListener('keyup', onSubmmitImageAttrs);
    iframeHeight.addEventListener('keydown', onSubmmitImageAttrs);
    iframeHeight.addEventListener('keypress', onSubmmitImageAttrs);

    const action = document.createElement('div');
    action.className = 'mely-editor-bubble-inline-actions';

    const imageApplyElement = document.createElement('button');
    const imageUnsetElement = document.createElement('button');
    imageApplyElement.type = 'button';
    imageApplyElement.textContent = 'Apply';
    imageUnsetElement.type = 'button';
    imageUnsetElement.textContent = 'Remove';

    imageApplyElement.addEventListener('click', setIframeAttrs);
    imageUnsetElement.addEventListener('click', () => {
      this.editor.chain().focus().deleteRange({ from, to }).run();
    });

    action.append(imageApplyElement, imageUnsetElement);
    iframeInfo.append(iframeSrc, iframeWidth, iframeHeight, action);

    this.element.append(iframeInfo);
    this.showTippy(view, from, to);
  }

  update(view: EditorView, oldState?: EditorState) {
    const { state, composing } = view;
    const { doc, selection } = state;
    const { ranges } = selection;
    const isSame = oldState && oldState.doc.eq(doc) && oldState.selection.eq(selection);
    const from = Math.min(...ranges.map(range => range.$from.pos));
    const to = Math.max(...ranges.map(range => range.$to.pos));

    if (composing || isSame) {
      return
    }

    this.createTooltip()
    const shouldShow = this.shouldShow?.({
      editor: this.editor,
      view,
      state,
      oldState,
      from,
      to,
    })

    if (!shouldShow) {
      return this.hide()
    }

    this.showLinkSettingPopup(view);
    this.showImageSettingPopup(view);
    this.showIframeSettingPopup(view);
  }

  show() {
    this.tippy?.show()
  }

  hide() {
    this.element.innerHTML = '';
    this.tippy?.hide()
  }

  destroy() {
    this.tippy?.destroy()
    this.element.removeEventListener('mousedown', this.mousedownHandler, { capture: true })
    this.view.dom.removeEventListener('dragstart', this.dragstartHandler)
    this.editor.off('focus', this.focusHandler)
    this.editor.off('blur', this.blurHandler)
  }
}

export const QuickEditPlugin = (options: QuickEditPluginProps) => {
  return new Plugin({
    key: typeof options.pluginKey === 'string'
      ? new PluginKey(options.pluginKey)
      : options.pluginKey,
    view: view => new QuickEditView({ view, ...options }),
  })
}