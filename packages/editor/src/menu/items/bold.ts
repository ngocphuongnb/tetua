import { MenuItem } from "../menu-item";

export class MenuBold extends MenuItem {
  protected command = 'bold';
  protected label = 'Bold';
  protected isActive = false;
  protected icon = `<svg viewBox="0 0 24 24"><path fill="currentColor" d="M13.5,15.5H10V12.5H13.5A1.5,1.5 0 0,1 15,14A1.5,1.5 0 0,1 13.5,15.5M10,6.5H13A1.5,1.5 0 0,1 14.5,8A1.5,1.5 0 0,1 13,9.5H10M15.6,10.79C16.57,10.11 17.25,9 17.25,8C17.25,5.74 15.5,4 13.25,4H7V18H14.04C16.14,18 17.75,16.3 17.75,14.21C17.75,12.69 16.89,11.39 15.6,10.79Z" /></svg>`;
  // protected icon = `<svg width="24" height="24" fill="currentColor" viewBox="0 0 24 24"><g><path fill="none" d="M0 0h24v24H0z"></path><path d="M8 11h4.5a2.5 2.5 0 1 0 0-5H8v5zm10 4.5a4.5 4.5 0 0 1-4.5 4.5H6V4h6.5a4.5 4.5 0 0 1 3.256 7.606A4.498 4.498 0 0 1 18 15.5zM8 13v5h5.5a2.5 2.5 0 1 0 0-5H8z"></path></g></svg>`;

  protected handler(e: MouseEvent) {
    e.preventDefault();
    this.editor.tiptapEditor.chain().focus().toggleBold().run();
  }
}