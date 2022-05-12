import { MenuItem } from "../menu-item";

export class MenuhorizontalRule extends MenuItem {
  protected command = 'horizontalRule';
  protected label = 'Horizontal rule';
  protected isActive = false;
  // protected icon = `<svg viewBox="0 0 24 24"><path fill="currentColor" d="M19,21H21V19H19M15,21H17V19H15M11,17H13V15H11M19,9H21V7H19M19,5H21V3H19M3,13H21V11H3M11,21H13V19H11M19,17H21V15H19M13,3H11V5H13M13,7H11V9H13M17,3H15V5H17M9,3H7V5H9M5,3H3V5H5M7,21H9V19H7M3,17H5V15H3M5,7H3V9H5M3,21H5V19H3V21Z" /></svg>`;
  protected icon = `<svg fill="currentColor" viewBox="0 0 24 24"><g><path fill="none" d="M0 0h24v24H0z"></path><path d="M2 11h2v2H2v-2zm4 0h12v2H6v-2zm14 0h2v2h-2v-2z"></path></g></svg>`;

  protected handler(e: MouseEvent) {
    e.preventDefault();
    this.editor.tiptapEditor.chain().focus().setHorizontalRule().run();
  }
}