import { TetuaEditor } from '..';
import { menuItems } from './items';
import { MenuItem } from './menu-item';

export class Menu {
  private editor: TetuaEditor;
  public menuElement: HTMLDivElement;
  public menuButtonsElement: HTMLDivElement;
  public menuContainerElement: HTMLDivElement;
  public menuSwitcherElement: HTMLLabelElement;
  private menuButtonInstances: MenuItem[] = [];

  public constructor(editor: TetuaEditor) {
    this.editor = editor;
    this.menuElement = document.createElement('div');
    this.menuElement.className = 'mely-editor-menu';

    this.menuButtonsElement = document.createElement('div');
    this.menuButtonsElement.className = 'mely-editor-menu-buttons';

    this.menuContainerElement = document.createElement('div');
    this.menuContainerElement.className = 'mely-editor-menu-container';

    this.menuButtonsElement.appendChild(this.menuElement);
    this.menuContainerElement.appendChild(this.menuButtonsElement);
  }

  public initialize() {
    this.editor.editorElement.insertBefore(this.menuContainerElement, this.editor.editorElement.firstChild);
    
    for (const MenuButton of menuItems) {
      const b = new MenuButton(this.editor);
      b.init();
      this.menuElement.appendChild(b.getElement());
      this.menuButtonInstances.push(b);
    }

    this.createSwitchEditorModeButton();
  }

  private createSwitchEditorModeButton() {
    this.menuSwitcherElement = document.createElement('label');
    const input = document.createElement('input');
    const slider = document.createElement('span');

    this.menuSwitcherElement.innerHTML = 'Markdown&nbsp;';
    this.menuSwitcherElement.className = 'switch';
    this.menuSwitcherElement.setAttribute('for', 'use-markdown');
    
    input.id = 'use-markdown';
    input.type = 'checkbox';
    input.value = '1';
    input.addEventListener('change', this.switchMode.bind(this));
    slider.className = 'slider';

    this.menuSwitcherElement.appendChild(input);
    this.menuSwitcherElement.appendChild(slider);

    this.menuSwitcherElement.appendChild(input);
    this.menuSwitcherElement.appendChild(slider);
    this.menuContainerElement.appendChild(this.menuSwitcherElement);
  }

  private switchMode(e: Event) {
    e.preventDefault();
    this.editor.useMarkdown = (e.target as HTMLInputElement).checked;
    if (this.editor.useMarkdown) {
      this.editor.tiptapEditorElement.style.display = 'none';
      this.editor.textareaElement.style.display = 'block';
      this.editor.tiptapEditorElement.setAttribute('hidden', 'true');
    } else {
      this.editor.tiptapEditorElement.style.display = 'block';
      this.editor.textareaElement.style.display = 'none';
      this.editor.tiptapEditorElement.removeAttribute('hidden');
    }
  }

  public update() {
    for (const b of this.menuButtonInstances) {
      b.update();
    }
  }
}