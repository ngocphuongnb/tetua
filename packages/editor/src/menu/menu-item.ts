import { TetuaEditor } from "..";

export class MenuItem {
  protected command = '';
  protected editor: TetuaEditor;
  protected label: string = '';
  protected icon: string = '';
  protected element: HTMLButtonElement | HTMLInputElement;
  protected subMenuitems: HTMLElement[] = [];

  constructor(editor: TetuaEditor) {
    this.editor = editor;
  }

  public init() {
    this.element = document.createElement('button');
    this.element.className = 'btn'

    if (this.icon) {
      this.element.innerHTML = this.icon;
      this.element.setAttribute('title', this.label);
    } else {
      this.element.innerText = this.label;
    }
    this.element.addEventListener('click', (e: MouseEvent) => {
      this.handler.bind(this)(e);

      if (this.command) {
        if (this.editor.tiptapEditor.isActive(this.command)) {
          this.element.classList.add('active');
        } else {
          this.element.classList.remove('active');
        }
      }
    });
  }

  public update() {
    if (this.editor.tiptapEditor.isActive(this.command)) {
      this.element.classList.add('active');
    } else {
      this.element.classList.remove('active');
    }
  }

  protected handler(e: MouseEvent): void {
    e.preventDefault();
  }

  getElement() {
    if (this.subMenuitems.length > 0) {
      const dropdown = document.createElement('div');
      dropdown.className = 'mely-editor-dropdown';
      const trigger = document.createElement('div');
      trigger.className = 'mely-editor-dropdown-trigger';
      trigger.innerHTML = '<svg viewBox="0 0 24 24"><path fill="currentColor" d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z" /></svg>';
      trigger.addEventListener('click', (e: MouseEvent) => {
        e.preventDefault();
        e.stopImmediatePropagation();
      });

      this.subMenuitems.forEach(item => {
        dropdown.appendChild(item);
      });
      trigger.append(dropdown);
      this.element.append(trigger);
      this.element.classList.add('has-sub-menu')
    }
    return this.element;
  }
}