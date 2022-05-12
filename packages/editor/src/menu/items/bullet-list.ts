import { MenuItem } from "../menu-item";

export class MenuBulletList extends MenuItem {
  protected command = 'bulletList';
  protected label = 'Bullet List';
  protected isActive = false;
  protected icon = `<svg viewBox="0 0 24 24"><path fill="currentColor" d="M3,4H7V8H3V4M9,5V7H21V5H9M3,10H7V14H3V10M9,11V13H21V11H9M3,16H7V20H3V16M9,17V19H21V17H9" /></svg>`;
  // protected icon = `<svg fill="currentColor" viewBox="0 0 24 24"><g><path fill="none" d="M0 0h24v24H0z"></path><path d="M8 4h13v2H8V4zM4.5 6.5a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm0 7a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm0 6.9a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zM8 11h13v2H8v-2zm0 7h13v2H8v-2z"></path></g></svg>`;

  protected handler(e: MouseEvent) {
    e.preventDefault();
    this.editor.tiptapEditor.chain().focus().toggleBulletList().run();
  }
}