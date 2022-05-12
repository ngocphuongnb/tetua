import { getMarkRange, Extension } from '@tiptap/core';
import { Plugin, TextSelection } from 'prosemirror-state';
import { QuickEditPlugin, QuickEditPluginProps } from './quick-edit'

export type QuickEditOptions = Omit<QuickEditPluginProps, 'editor' | 'element'>

export const QuickEdit = Extension.create<QuickEditOptions>({
  name: 'quickEdit',

  addOptions() {
    return {
      tippyOptions: {},
      pluginKey: 'quickEdit',
      shouldShow: null,
    }
  },

  addProseMirrorPlugins() {
    return [
      new Plugin({
        props: {
          handleClick(view, pos) {
            const range = getMarkRange(view.state.doc.resolve(pos), view.state.schema.marks.link);
            if (!range) return false;

            const $start = view.state.doc.resolve(range.from);
            const $end = view.state.doc.resolve(range.to);
            const transaction = view.state.tr.setSelection(new TextSelection($start, $end));
            view.dispatch(transaction);

            return true;
          },
        },
      }),
      QuickEditPlugin({
        pluginKey: this.options.pluginKey,
        editor: this.editor,
        tippyOptions: this.options.tippyOptions,
        shouldShow: this.options.shouldShow,
      }),
    ]
  },
})