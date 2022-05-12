import { MenuItem } from "../menu-item";

export class MenuItalic extends MenuItem {
  protected command = 'italic';
  protected label = 'Italic';
  protected isActive = false;
  protected icon = `<svg viewBox="0 0 24 24"><path fill="currentColor" d="M10,4V7H12.21L8.79,15H6V18H14V15H11.79L15.21,7H18V4H10Z" /></svg>`;
  // protected icon = `<svg fill="currentColor" viewBox="0 0 24 24"><g><path fill="none" d="M0 0h24v24H0z"></path><path d="M15 20H7v-2h2.927l2.116-12H9V4h8v2h-2.927l-2.116 12H15z"></path></g></svg>`;

  protected handler(e: MouseEvent) {
    e.preventDefault();
    this.editor.tiptapEditor.chain().focus().toggleItalic().run();
  }
}