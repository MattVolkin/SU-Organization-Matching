import { mount } from 'svelte';

export function mountPage(Component) {
  mount(Component, {
    target: document.getElementById('app'),
  });
}
