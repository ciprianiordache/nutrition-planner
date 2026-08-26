<script lang="ts">
  export type PlanFormData = {
    title: string;
    generalRecommendations: string;
    footer: string;
  };

  let {
    initial = null,
    onSubmit,
    onCancel,
    submitting = false,
    error = null,
    defaultFooter = "",
  }: {
    initial?: PlanFormData | null;
    onSubmit: (data: PlanFormData) => void;
    onCancel: () => void;
    submitting?: boolean;
    error?: string | null;
    defaultFooter?: string;
  } = $props();

  let title = $state(initial?.title ?? "");
  let generalRecommendations = $state(initial?.generalRecommendations ?? "");
  let footer = $state(initial?.footer ?? defaultFooter);

  function submit(e: SubmitEvent) {
    e.preventDefault();
    if (!title.trim()) return;
    onSubmit({ title: title.trim(), generalRecommendations, footer });
  }
</script>

<form onsubmit={submit} class="space-y-4">
  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="plan-title">Titlu plan</label>
    <input
      id="plan-title"
      bind:value={title}
      required
      placeholder="Ex: Plan alimentar — August"
      class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    />
  </div>

  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="plan-recommendations">
      Recomandări generale
    </label>
    <textarea
      id="plan-recommendations"
      bind:value={generalRecommendations}
      rows="3"
      placeholder="Ex: 3 litri de apă zilnic, evită prăjelile, minim 7 ore de somn."
      class="w-full resize-none rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    ></textarea>
  </div>

  <div>
    <label class="mb-1 block text-xs font-medium text-muted-foreground" for="plan-footer">Footer plan</label>
    <textarea
      id="plan-footer"
      bind:value={footer}
      rows="2"
      class="w-full resize-none rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
    ></textarea>
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
      {submitting ? "Se salvează…" : initial ? "Salvează modificările" : "Salvează plan"}
    </button>
  </div>
</form>
