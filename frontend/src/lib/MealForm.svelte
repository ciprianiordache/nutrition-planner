<script lang="ts">
  export type MealFormData = {
    name: string;
    note: string;
    order: number;
  };

  let {
    initial = null,
    onSubmit,
    onCancel,
    submitting = false,
    error = null,
  }: {
    initial?: MealFormData | null;
    onSubmit: (data: MealFormData) => void;
    onCancel: () => void;
    submitting?: boolean;
    error?: string | null;
  } = $props();

  let name = $state(initial?.name ?? "");
  let note = $state(initial?.note ?? "");
  let order = $state(initial?.order ?? 0);

  function submit(e: SubmitEvent) {
    e.preventDefault();
    if (!name.trim()) return;
    onSubmit({ name: name.trim(), note, order });
  }
</script>

<form onsubmit={submit} class="space-y-4">
  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="meal-name">Nume masă</label>
    <input
      id="meal-name"
      bind:value={name}
      required
      placeholder="Ex: Mic dejun"
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    />
  </div>

  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="meal-note">Notă</label>
    <input
      id="meal-note"
      bind:value={note}
      placeholder="Ex: În maxim o oră după trezire"
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    />
  </div>

  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="meal-order">Ordine</label>
    <input
      id="meal-order"
      type="number"
      step="1"
      min="0"
      bind:value={order}
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    />
  </div>

  {#if error}
    <p class="text-sm text-destructive">{error}</p>
  {/if}

  <div class="flex justify-end gap-2 pt-2">
    <button
      type="button"
      onclick={onCancel}
      class="rounded-md border border-border px-4 py-2 text-sm font-medium transition-colors hover:bg-accent"
    >
      Anulează
    </button>
    <button
      type="submit"
      disabled={submitting}
      class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-60"
    >
      {submitting ? "Se salvează…" : initial ? "Salvează modificările" : "Adaugă masă"}
    </button>
  </div>
</form>
