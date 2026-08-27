<script lang="ts">
  import { t } from "./i18n";

  export type SupplementFormData = {
    supplement: string;
    quantity: number;
    order: number;
  };

  let {
    initial = null,
    onSubmit,
    onCancel,
    submitting = false,
    error = null,
  }: {
    initial?: SupplementFormData | null;
    onSubmit: (data: SupplementFormData) => void;
    onCancel: () => void;
    submitting?: boolean;
    error?: string | null;
  } = $props();

  let supplement = $state(initial?.supplement ?? "");
  let quantity = $state<number | undefined>(initial?.quantity ?? 1);
  let order = $state(initial?.order ?? 0);

  function submit(e: SubmitEvent) {
    e.preventDefault();
    if (!supplement.trim() || !quantity) return;
    onSubmit({ supplement: supplement.trim(), quantity, order });
  }
</script>

<form onsubmit={submit} class="space-y-4">
  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="supp-name">{$t("supplementForm.name")}</label>
    <input
      id="supp-name"
      bind:value={supplement}
      required
      placeholder={$t("supplementForm.namePlaceholder")}
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    />
  </div>

  <div class="grid grid-cols-2 gap-3">
    <div>
      <label class="mb-1 block text-xs font-medium text-muted-foreground" for="supp-qty">{$t("common.quantity")}</label>
      <input
        id="supp-qty"
        type="number"
        step="0.5"
        min="0"
        bind:value={quantity}
        required
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
      />
    </div>
    <div>
      <label class="mb-1 block text-xs font-medium text-muted-foreground" for="supp-order">{$t("common.order")}</label>
      <input
        id="supp-order"
        type="number"
        step="1"
        min="0"
        bind:value={order}
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
      />
    </div>
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
      {$t("common.cancel")}
    </button>
    <button
      type="submit"
      disabled={submitting}
      class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-60"
    >
      {submitting ? $t("common.saving") : initial ? $t("common.saveChanges") : $t("supplementForm.add")}
    </button>
  </div>
</form>
