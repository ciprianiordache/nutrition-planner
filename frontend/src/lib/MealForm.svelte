<script lang="ts">
  import { t } from "./i18n";

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
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="meal-name">{$t("mealForm.name")}</label>
    <input
      id="meal-name"
      bind:value={name}
      required
      placeholder={$t("mealForm.namePlaceholder")}
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    />
  </div>

  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="meal-note">{$t("mealForm.note")}</label>
    <input
      id="meal-note"
      bind:value={note}
      placeholder={$t("mealForm.notePlaceholder")}
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    />
  </div>

  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="meal-order">{$t("common.order")}</label>
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
      {$t("common.cancel")}
    </button>
    <button
      type="submit"
      disabled={submitting}
      class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-60"
    >
      {submitting ? $t("common.saving") : initial ? $t("common.saveChanges") : $t("mealForm.add")}
    </button>
  </div>
</form>
