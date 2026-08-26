<script lang="ts">
  export type ClientFormData = {
    name: string;
    weight: number;
    height: number;
    target: string;
  };

  let {
    initial = null,
    onSubmit,
    onCancel,
    submitting = false,
    error = null,
  }: {
    initial?: ClientFormData | null;
    onSubmit: (data: ClientFormData) => void;
    onCancel: () => void;
    submitting?: boolean;
    error?: string | null;
  } = $props();

  let name = $state(initial?.name ?? "");
  let weight = $state<number | undefined>(initial?.weight);
  let height = $state<number | undefined>(initial?.height);
  let target = $state(initial?.target ?? "");

  function submit(e: SubmitEvent) {
    e.preventDefault();
    if (!name.trim() || !weight || !height) return;
    onSubmit({ name: name.trim(), weight, height, target });
  }
</script>

<form onsubmit={submit} class="space-y-4">
  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="client-name">Nume</label>
    <input
      id="client-name"
      bind:value={name}
      required
      placeholder="Ex: Andreea Popescu"
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    />
  </div>

  <div class="grid grid-cols-2 gap-3">
    <div>
      <label class="mb-1 block text-xs font-medium text-muted-foreground" for="client-weight">Greutate (kg)</label>
      <input
        id="client-weight"
        type="number"
        step="0.1"
        min="0"
        bind:value={weight}
        required
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
      />
    </div>
    <div>
      <label class="mb-1 block text-xs font-medium text-muted-foreground" for="client-height">Înălțime (cm)</label>
      <input
        id="client-height"
        type="number"
        step="0.1"
        min="0"
        bind:value={height}
        required
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
      />
    </div>
  </div>

  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="client-target">Obiectiv</label>
    <select
      id="client-target"
      bind:value={target}
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    >
      <option value="">Alege obiectivul</option>
      <option value="Slăbire">Slăbire</option>
      <option value="Masă musculară">Masă musculară</option>
      <option value="Menținere">Menținere</option>
      <option value="Recompoziție">Recompoziție</option>
    </select>
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
      {submitting ? "Se salvează…" : initial ? "Salvează modificările" : "Salvează client"}
    </button>
  </div>
</form>
