<script lang="ts">
  export type FoodFormData = {
    baseFood: string;
    baseQuantity: number;
    baseUnit: string;
    alternativeFood?: string;
    alternativeQuantity?: number;
    alternativeUnit?: string;
    order: number;
  };

  let {
    initial = null,
    onSubmit,
    onCancel,
    submitting = false,
    error = null,
  }: {
    initial?: FoodFormData | null;
    onSubmit: (data: FoodFormData) => void;
    onCancel: () => void;
    submitting?: boolean;
    error?: string | null;
  } = $props();

  let baseFood = $state(initial?.baseFood ?? "");
  let baseQuantity = $state<number | undefined>(initial?.baseQuantity);
  let baseUnit = $state(initial?.baseUnit ?? "g");
  let hasAlternative = $state(!!initial?.alternativeFood);
  let alternativeFood = $state(initial?.alternativeFood ?? "");
  let alternativeQuantity = $state<number | undefined>(initial?.alternativeQuantity);
  let alternativeUnit = $state(initial?.alternativeUnit ?? "g");
  let order = $state(initial?.order ?? 0);

  function submit(e: SubmitEvent) {
    e.preventDefault();
    if (!baseFood.trim() || !baseQuantity || !baseUnit.trim()) return;

    onSubmit({
      baseFood: baseFood.trim(),
      baseQuantity,
      baseUnit: baseUnit.trim(),
      alternativeFood: hasAlternative && alternativeFood.trim() ? alternativeFood.trim() : undefined,
      alternativeQuantity: hasAlternative && alternativeQuantity ? alternativeQuantity : undefined,
      alternativeUnit: hasAlternative && alternativeUnit.trim() ? alternativeUnit.trim() : undefined,
      order,
    });
  }
</script>

<form onsubmit={submit} class="space-y-4">
  <div class="grid grid-cols-[1fr_auto_auto] gap-3">
    <div>
      <label class="mb-1 block text-xs font-medium text-muted-foreground" for="food-name">Aliment</label>
      <input
        id="food-name"
        bind:value={baseFood}
        required
        placeholder="Ex: Piept de pui"
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
      />
    </div>
    <div class="w-24">
      <label class="mb-1 block text-xs font-medium text-muted-foreground" for="food-qty">Cantitate</label>
      <input
        id="food-qty"
        type="number"
        step="0.1"
        min="0"
        bind:value={baseQuantity}
        required
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
      />
    </div>
    <div class="w-20">
      <label class="mb-1 block text-xs font-medium text-muted-foreground" for="food-unit">Unitate</label>
      <input
        id="food-unit"
        bind:value={baseUnit}
        required
        placeholder="g"
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
      />
    </div>
  </div>

  <label class="flex items-center gap-2 text-sm">
    <input type="checkbox" bind:checked={hasAlternative} class="rounded border-input" />
    Adaugă alternativă
  </label>

  {#if hasAlternative}
    <div class="grid grid-cols-[1fr_auto_auto] gap-3 rounded-md border border-dashed border-border p-3">
      <div>
        <label class="mb-1 block text-xs font-medium text-muted-foreground" for="food-alt-name">Alternativă</label>
        <input
          id="food-alt-name"
          bind:value={alternativeFood}
          placeholder="Ex: Somon"
          class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
        />
      </div>
      <div class="w-24">
        <label class="mb-1 block text-xs font-medium text-muted-foreground" for="food-alt-qty">Cantitate</label>
        <input
          id="food-alt-qty"
          type="number"
          step="0.1"
          min="0"
          bind:value={alternativeQuantity}
          class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
        />
      </div>
      <div class="w-20">
        <label class="mb-1 block text-xs font-medium text-muted-foreground" for="food-alt-unit">Unitate</label>
        <input
          id="food-alt-unit"
          bind:value={alternativeUnit}
          placeholder="g"
          class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
        />
      </div>
    </div>
  {/if}

  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="food-order">Ordine</label>
    <input
      id="food-order"
      type="number"
      step="1"
      min="0"
      bind:value={order}
      class="w-24 rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
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
      {submitting ? "Se salvează…" : initial ? "Salvează modificările" : "Adaugă aliment"}
    </button>
  </div>
</form>
