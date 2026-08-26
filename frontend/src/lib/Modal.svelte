<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    open = false,
    title,
    onClose,
    children,
  }: {
    open?: boolean;
    title: string;
    onClose: () => void;
    children: Snippet;
  } = $props();

  function onOverlayClick(e: MouseEvent) {
    if (e.target === e.currentTarget) onClose();
  }

  function onOverlayKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") onClose();
  }
</script>

<svelte:window onkeydown={open ? onOverlayKeydown : undefined} />

{#if open}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4"
    onclick={onOverlayClick}
    role="presentation"
  >
    <div
      class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-xl"
      role="dialog"
      aria-modal="true"
      aria-label={title}
      tabindex="-1"
    >
      <div class="mb-4 flex items-center justify-between">
        <h2 class="font-display text-lg font-semibold">{title}</h2>
        <button
          onclick={onClose}
          aria-label="Închide"
          class="rounded-md p-1 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
        >
          ✕
        </button>
      </div>
      {@render children()}
    </div>
  </div>
{/if}
