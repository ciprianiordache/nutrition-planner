<script lang="ts">
  import Modal from "./Modal.svelte";

  let {
    open = false,
    title = "Confirmă ștergerea",
    message,
    confirmLabel = "Șterge",
    submitting = false,
    error = null,
    onConfirm,
    onCancel,
  }: {
    open?: boolean;
    title?: string;
    message: string;
    confirmLabel?: string;
    submitting?: boolean;
    error?: string | null;
    onConfirm: () => void;
    onCancel: () => void;
  } = $props();
</script>

<Modal {open} {title} onClose={onCancel}>
  <p class="text-sm text-muted-foreground">{message}</p>

  {#if error}
    <p class="mt-3 text-sm text-destructive">{error}</p>
  {/if}

  <div class="mt-5 flex justify-end gap-2">
    <button
      type="button"
      onclick={onCancel}
      class="rounded-md border border-border px-4 py-2 text-sm font-medium transition-colors hover:bg-accent"
    >
      Anulează
    </button>
    <button
      type="button"
      disabled={submitting}
      onclick={onConfirm}
      class="rounded-md bg-destructive px-4 py-2 text-sm font-medium text-destructive-foreground transition-colors hover:bg-destructive/90 disabled:opacity-60"
    >
      {submitting ? "Se șterge…" : confirmLabel}
    </button>
  </div>
</Modal>
