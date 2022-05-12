window.addEventListener('load', function () {
  var menuTriggerElms = Array.from(document.querySelectorAll('.menu-trigger'));
  for (var menuTriggerElm of menuTriggerElms) {
    menuTriggerElm.addEventListener('click', function() {
      document.querySelector('body').classList.toggle('show-menu');
    });
  }
});