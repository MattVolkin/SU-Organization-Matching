<!--
@component MultiSelectDropdown
A custom multiselect dropdown component with tag display and checkbox selection.
Similar to job application multiselect dropdowns with visual tags for selected items
and a dropdown containing checkboxes for easy selection.
Only one dropdown can be open at a time across the page.

@type {string} label - The label for the dropdown
@type {string[]} options - Array of available options to select from
@type {string[]} value - Array of selected values (bindable)
@type {string} id - Unique identifier for this dropdown instance
@type {boolean} isOpen - Whether dropdown is open (managed globally)
@type {string} searchTerm - Search filter for options
-->

<script lang="ts">
  import { tick } from 'svelte';
  import { writable, type Writable } from 'svelte/store';

  interface Props {
    label: string;
    options: string[];
    value?: string[];
    id: string;
  }

  let { label, options, value = $bindable([]), id } = $props();

  // Module-level store to track which dropdown is open (shared across all instances)
  const getOpenDropdownStore = (): Writable<string | null> => {
    if (!globalThis.__multiselectOpenDropdown__) {
      globalThis.__multiselectOpenDropdown__ = writable<string | null>(null);
    }
    return globalThis.__multiselectOpenDropdown__;
  };

  const openDropdownStore = getOpenDropdownStore();
  let isOpen = $derived($openDropdownStore === id);

  let searchTerm = $state('');
  let dropdownElement = $state<HTMLDivElement>();
  let openAbove = $state(false);
  let menuMaxHeight = $state('300px');

  function openThis() {
    openDropdownStore.set(id);
  }

  function closeThis() {
    if ($openDropdownStore === id) {
      openDropdownStore.set(null);
    }
  }

  async function toggleOpen() {
    if (isOpen) {
      closeThis();
      openAbove = false;
      menuMaxHeight = '300px';
    } else {
      openThis();
      await tick();

      const bounds = dropdownElement?.getBoundingClientRect();
      if (!bounds) {
        openAbove = false;
        menuMaxHeight = '300px';
        return;
      }

      const viewportPadding = 16;
      const spaceBelow = window.innerHeight - bounds.bottom - viewportPadding;
      const spaceAbove = bounds.top - viewportPadding;
      const shouldOpenAbove = spaceBelow < 240 && spaceAbove > spaceBelow;
      const availableSpace = shouldOpenAbove ? spaceAbove : spaceBelow;

      openAbove = shouldOpenAbove;
      menuMaxHeight = `${Math.max(180, Math.min(300, availableSpace))}px`;
    }
  }

  $effect.root(() => {
    function handleClickOutside(event: MouseEvent) {
      if (isOpen && dropdownElement && !dropdownElement.contains(event.target as Node)) {
        closeThis();
      }
    }

    if (isOpen) {
      document.addEventListener('click', handleClickOutside);
      return () => document.removeEventListener('click', handleClickOutside);
    }
  });

  function toggleOption(option: string) {
    if (value.includes(option)) {
      value = value.filter((item) => item !== option);
    } else {
      value = [...value, option];
    }
  }

  function removeOption(option: string) {
    value = value.filter((item) => item !== option);
  }

  function filteredOptions() {
    return options.filter((option) => option.toLowerCase().includes(searchTerm.toLowerCase()));
  }

  function handleSearchKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      event.preventDefault();
      event.stopPropagation();
    }
  }
</script>

<div class="multiselect-container" bind:this={dropdownElement}>
  <label class="multiselect-label" for="multiselect-input-{id}">{label}</label>
  <button
    type="button"
    id="multiselect-input-{id}"
    class="multiselect-input"
    onclick={() => toggleOpen()}
    aria-haspopup="listbox"
    aria-expanded={isOpen}
  >
    <div class="selected-tags">
      {#if value.length === 0}
        <span class="placeholder">Select options...</span>
      {:else}
        {#each value as item}
          <span class="tag">
            {item}
            <div
              role="button"
              class="tag-remove"
              onmousedown={(e) => e.preventDefault()}
              onclick={(e) => {
                e.stopPropagation();
                removeOption(item);
              }}
              tabindex="0"
              onkeydown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault();
                  removeOption(item);
                }
              }}
              aria-label="Remove {item}"
            >
              ×
            </div>
          </span>
        {/each}
      {/if}
    </div>
    <svg class="dropdown-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
      <path d="M6 9l6 6 6-6" />
    </svg>
  </button>

  {#if isOpen}
    <div class={`dropdown-menu ${openAbove ? 'open-above' : 'open-below'}`} role="listbox" style={`max-height: ${menuMaxHeight};`}>
      <input
        type="text"
        class="search-input"
        placeholder="Search..."
        bind:value={searchTerm}
        onkeydown={handleSearchKeydown}
      />
      <div class="options-list">
        {#each filteredOptions() as option}
          <label class="option-item">
            <input
              type="checkbox"
              checked={value.includes(option)}
              onchange={() => toggleOption(option)}
            />
            <span>{option}</span>
          </label>
        {/each}
        {#if filteredOptions().length === 0}
          <div class="no-results">No options found</div>
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .multiselect-container {
    position: relative;
    display: grid;
    gap: 0.35rem;
  }

  .multiselect-label {
    font-weight: 700;
    font-size: 0.95rem;
    color: #132c45;
  }

  .multiselect-input {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0.75rem;
    border: 1px solid #c9d6e5;
    border-radius: 0.5rem;
    background: white;
    cursor: pointer;
    transition: border-color 0.2s, box-shadow 0.2s;
    min-height: 2.5rem;
    font-size: 0.95rem;
    font-family: inherit;
    color: inherit;
    width: 100%;
    text-align: left;
  }

  .multiselect-input:hover {
    border-color: #0f6d8c;
  }

  .multiselect-input:focus-within {
    border-color: #0f6d8c;
    box-shadow: 0 0 0 2px rgba(15, 109, 140, 0.1);
  }

  .selected-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    flex: 1;
    align-items: center;
  }

  .placeholder {
    color: #8091a5;
    font-size: 0.9rem;
  }

  .tag {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.35rem 0.6rem;
    background: #e5edf8;
    color: #132c45;
    border-radius: 0.375rem;
    font-size: 0.85rem;
    font-weight: 600;
    white-space: nowrap;
  }

  .tag-remove {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1.2rem;
    height: 1.2rem;
    padding: 0;
    margin: -0.15rem -0.3rem -0.15rem 0.15rem;
    background: transparent;
    border: none;
    color: #0f6d8c;
    cursor: pointer;
    font-size: 1.1rem;
    line-height: 1;
    border-radius: 0.25rem;
    transition: background 0.15s;
  }

  .tag-remove:hover {
    background: rgba(15, 109, 140, 0.15);
  }

  .dropdown-icon {
    width: 1.2rem;
    height: 1.2rem;
    flex: 0 0 auto;
    color: #8091a5;
    pointer-events: none;
  }

  .dropdown-menu {
    position: absolute;
    top: calc(100% + 0.35rem);
    left: 0;
    right: 0;
    background: white;
    border: 1px solid #c9d6e5;
    border-radius: 0.5rem;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    z-index: 1000;
    max-height: 300px;
    display: flex;
    flex-direction: column;
  }

  .dropdown-menu.open-above {
    top: auto;
    bottom: calc(100% + 0.35rem);
  }

  .search-input {
    padding: 0.6rem 0.75rem;
    border: none;
    border-bottom: 1px solid #e5edf8;
    border-radius: 0.5rem 0.5rem 0 0;
    font-size: 0.9rem;
    outline: none;
  }

  .search-input:focus {
    background: #f7fbff;
  }

  .options-list {
    flex: 1;
    overflow-y: auto;
    min-height: 0;
  }

  .option-item {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.6rem 0.75rem;
    cursor: pointer;
    transition: background 0.15s;
    border: none;
    background: transparent;
    width: 100%;
    font-size: 0.9rem;
    color: #132c45;
  }

  .option-item:hover {
    background: #f7fbff;
  }

  .option-item input[type='checkbox'] {
    width: 1rem;
    height: 1rem;
    margin: 0;
    accent-color: #0f6d8c;
    flex: 0 0 auto;
    cursor: pointer;
  }

  .option-item span {
    flex: 1;
    user-select: none;
  }

  .no-results {
    padding: 1rem 0.75rem;
    text-align: center;
    color: #8091a5;
    font-size: 0.9rem;
  }

  @media (prefers-color-scheme: dark) {
    .multiselect-label {
      color: #e5edf8;
    }

    .multiselect-input {
      background: #1a2332;
      border-color: #31506e;
      color: #e5edf8;
    }

    .multiselect-input:hover {
      border-color: #4db8d8;
    }

    .multiselect-input:focus-within {
      border-color: #4db8d8;
      box-shadow: 0 0 0 2px rgba(77, 184, 216, 0.1);
    }

    .placeholder {
      color: #7b8fa0;
    }

    .tag {
      background: #253448;
      color: #b6c7df;
    }

    .tag-remove {
      color: #4db8d8;
    }

    .tag-remove:hover {
      background: rgba(77, 184, 216, 0.15);
    }

    .dropdown-icon {
      color: #7b8fa0;
    }

    .dropdown-menu {
      background: #1a2332;
      border-color: #31506e;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    }

    .search-input {
      background: #1a2332;
      border-color: #31506e;
      color: #e5edf8;
    }

    .search-input:focus {
      background: #222c3a;
    }

    .option-item {
      color: #e5edf8;
    }

    .option-item:hover {
      background: #222c3a;
    }

    .option-item input[type='checkbox'] {
      accent-color: #4db8d8;
    }

    .no-results {
      color: #7b8fa0;
    }
  }
</style>
