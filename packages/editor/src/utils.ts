
export const createNodeViewBlock = (contentDom: HTMLElement, viewDomElms: HTMLElement[]) => {
  const block = document.createElement('div');
  const blockView = document.createElement('div');

  block.className = 'mely-editor-block';
  blockView.className = 'mely-editor-block-view';
  blockView.setAttribute('contenteditable', 'false');

  if (viewDomElms.length > 0) {
    blockView.append(...viewDomElms);
  }
  block.append(contentDom, blockView);


  return {
    dom: block,
    view: blockView,
  };
}