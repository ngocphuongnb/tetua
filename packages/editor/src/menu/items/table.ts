import { TetuaEditor } from "../..";
import { MenuItem } from "../menu-item";

interface MenuItemConfig {
  handler: string;
  check: string;
  text: string;
}

const menuItems: MenuItemConfig[] = [
  {
    handler: 'addColumnBefore',
    check: 'addColumnBefore',
    text: 'Add column before',
  },
  {
    handler: 'addColumnAfter',
    check: 'addColumnAfter',
    text: 'Add column after',
  },
  {
    handler: 'deleteColumn',
    check: 'deleteColumn',
    text: 'Delete column',
  },
  {
    handler: 'addRowBefore',
    check: 'addRowBefore',
    text: 'Add row before',
  },
  {
    handler: 'addRowAfter',
    check: 'addRowAfter',
    text: 'Add row after',
  },
  {
    handler: 'deleteRow',
    check: 'deleteRow',
    text: 'Delete row',
  },
  {
    handler: 'deleteTable',
    check: 'deleteTable',
    text: 'Delete table',
  },
  {
    handler: 'mergeCells',
    check: 'mergeCells',
    text: 'Merge cells',
  },
  {
    handler: 'splitCell',
    check: 'splitCell',
    text: 'Split cell',
  },
  {
    handler: 'toggleHeaderColumn',
    check: 'toggleHeaderColumn',
    text: 'Toggle header column',
  },
  {
    handler: 'toggleHeaderRow',
    check: 'toggleHeaderRow',
    text: 'Toggle header row',
  },
  {
    handler: 'toggleHeaderCell',
    check: 'toggleHeaderCell',
    text: 'Toggle header cell',
  },
  {
    handler: 'mergeOrSplit',
    check: 'mergeOrSplit',
    text: 'Merge or split',
  },
  // {
  //   handler: 'fixTables',
  //   check: 'fixTables',
  //   text: 'fixTables',
  // },
  // {
  //   handler: 'goToNextCell',
  //   check: 'goToNextCell',
  //   text: 'goToNextCell',
  // },
  // {
  //   handler: 'goToPreviousCell',
  //   check: 'goToPreviousCell',
  //   text: 'goToPreviousCell',
  // },
];

export class MenuTable extends MenuItem {
  protected command = 'table';
  protected label = 'Table';
  protected isActive = false;
  protected icon = `<svg viewBox="0 0 24 24"><path fill="currentColor" d="M5,4H19A2,2 0 0,1 21,6V18A2,2 0 0,1 19,20H5A2,2 0 0,1 3,18V6A2,2 0 0,1 5,4M5,8V12H11V8H5M13,8V12H19V8H13M5,14V18H11V14H5M13,14V18H19V14H13Z" /></svg>`;
  // protected icon = `<svg fill="currentColor" viewBox="0 0 24 24"><g><path fill="none" d="M0 0h24v24H0z"></path><path d="M4 8h16V5H4v3zm10 11v-9h-4v9h4zm2 0h4v-9h-4v9zm-8 0v-9H4v9h4zM3 3h18a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z"></path></g></svg>`;

  constructor(editor: TetuaEditor) {
    super(editor);
    menuItems.forEach(item => {
      const menuItem = this.createTableButton(item);
      this.subMenuitems.push(menuItem);
    });
  }

  protected handler(e: MouseEvent) {
    e.preventDefault();
    this.editor.tiptapEditor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run();
  }

  private createTableButton(config: MenuItemConfig) {
    const can = this.editor.tiptapEditor.can() as any;
    const container = document.createElement('div');
    const button = document.createElement('button');
    button.type = 'button';
    button.innerText = config.text
    button.disabled = !can[config.check]();
    button.dataset.check = config.check;
    button.addEventListener('click', (e: MouseEvent) => {
      const chain = this.editor.tiptapEditor.chain().focus() as any;
      e.preventDefault();
      e.stopImmediatePropagation();
      chain[config.handler]().run();
    });
    container.append(button);
    return container;
  }

  public update() {
    super.update();
    this.subMenuitems.forEach((item: HTMLButtonElement) => {
      const button = item.querySelector('button');
      const checkFn = button.dataset.check;
      const can = this.editor.tiptapEditor.can() as any;
      button.disabled = !can[checkFn]();
    });
  }
}