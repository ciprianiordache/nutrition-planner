<script lang="ts">
  import type { plan, day_variant, meal, food, supplement } from "../../wailsjs/go/models";
  import { planApi, dayVariantApi, mealApi, foodApi, supplementApi, exportApi, errorMessage } from "../lib/api";
  import { currentPlanId, currentView } from "../lib/stores";
  import { t } from "../lib/i18n";
  import Modal from "../lib/Modal.svelte";
  import ConfirmDialog from "../lib/ConfirmDialog.svelte";
  import DayVariantForm from "../lib/DayVariantForm.svelte";
  import MealForm from "../lib/MealForm.svelte";
  import FoodForm from "../lib/FoodForm.svelte";
  import SupplementForm from "../lib/SupplementForm.svelte";

  const planId = $derived($currentPlanId);

  let planFull = $state<plan.FullResponse | null>(null);
  let loading = $state(true);
  let loadError = $state<string | null>(null);

  let activeDayVariantId = $state<string | null>(null);
  const activeDayVariant = $derived(
    planFull?.dayVariants?.find((d) => d.id === activeDayVariantId) ?? planFull?.dayVariants?.[0] ?? null
  );

  $effect(() => {
    if (planId) {
      loadPlan();
    } else {
      planFull = null;
    }
  });

  async function loadPlan() {
    if (!planId) return;
    loading = true;
    loadError = null;
    try {
      planFull = await planApi.getFull(planId);
      if (!activeDayVariantId && planFull.dayVariants && planFull.dayVariants.length > 0) {
        activeDayVariantId = planFull.dayVariants[0].id;
      }
    } catch (err) {
      loadError = errorMessage(err);
    } finally {
      loading = false;
    }
  }

  function backToClients() {
    currentPlanId.set(null);
    currentView.set("clients");
  }

  // --- Export PDF ----------------------------------------------------------

  let exporting = $state(false);
  let exportError = $state<string | null>(null);
  let exportSuccess = $state<{ fileName: string; path: string } | null>(null);

  async function handleExportPDF() {
    if (!planId) return;
    exporting = true;
    exportError = null;
    exportSuccess = null;
    try {
      const res = await exportApi.exportPlanPDF(planId);
      exportSuccess = { fileName: res.fileName, path: res.path };
    } catch (err) {
      exportError = errorMessage(err);
    } finally {
      exporting = false;
    }
  }

  async function handleOpenFolder() {
    if (!exportSuccess) return;
    try {
      await exportApi.openExportFolder(exportSuccess.path);
    } catch {
      // Best-effort - the export itself already succeeded and the path is
      // shown right next to this button either way, so a failure to also
      // auto-open the folder isn't worth surfacing as an error.
    }
  }

  // --- Previzualizare PDF ----------------------------------------------

  let previewing = $state(false);
  let previewError = $state<string | null>(null);
  let previewSuccess = $state<{ fileName: string; path: string } | null>(null);

  async function handlePreviewPDF() {
    if (!planId) return;
    previewing = true;
    previewError = null;
    previewSuccess = null;
    try {
      // Backend-ul randează PDF-ul într-un folder temporar (nu în
      // exportul final al clientului) și îl deschide el însuși în
      // vizualizatorul implicit de PDF al sistemului (Wails nu permite
      // deschiderea de URL-uri file:// din frontend, din motive de
      // securitate) — aici doar arătăm rezultatul, ca backup vizual dacă
      // deschiderea automată eșuează pe mașina utilizatorului.
      const res = await exportApi.previewPlanPDF(planId);
      previewSuccess = { fileName: res.fileName, path: res.path };
    } catch (err) {
      previewError = errorMessage(err);
    } finally {
      previewing = false;
    }
  }

  // --- Day variant CRUD --------------------------------------------------

  let showDayVariantForm = $state(false);
  let editingDayVariant = $state<day_variant.Response | null>(null);
  let savingDayVariant = $state(false);
  let dayVariantFormError = $state<string | null>(null);

  let deletingDayVariant = $state<day_variant.Response | null>(null);
  let deletingDayVariantBusy = $state(false);
  let deletingDayVariantError = $state<string | null>(null);

  function openCreateDayVariant() {
    editingDayVariant = null;
    dayVariantFormError = null;
    showDayVariantForm = true;
  }

  function openEditDayVariant(dv: day_variant.Response) {
    editingDayVariant = dv;
    dayVariantFormError = null;
    showDayVariantForm = true;
  }

  async function handleDayVariantSubmit(data: { name: string; order: number }) {
    if (!planId) return;
    savingDayVariant = true;
    dayVariantFormError = null;
    try {
      if (editingDayVariant) {
        await dayVariantApi.update({ id: editingDayVariant.id, ...data });
      } else {
        const created = await dayVariantApi.create({ planId, ...data });
        activeDayVariantId = created.id;
      }
      showDayVariantForm = false;
      await loadPlan();
    } catch (err) {
      dayVariantFormError = errorMessage(err);
    } finally {
      savingDayVariant = false;
    }
  }

  function openDeleteDayVariant(dv: day_variant.Response) {
    deletingDayVariant = dv;
    deletingDayVariantError = null;
  }

  async function confirmDeleteDayVariant() {
    if (!deletingDayVariant) return;
    deletingDayVariantBusy = true;
    deletingDayVariantError = null;
    try {
      await dayVariantApi.remove(deletingDayVariant.id);
      if (activeDayVariantId === deletingDayVariant.id) activeDayVariantId = null;
      deletingDayVariant = null;
      await loadPlan();
    } catch (err) {
      deletingDayVariantError = errorMessage(err);
    } finally {
      deletingDayVariantBusy = false;
    }
  }

  // --- Meal CRUD -----------------------------------------------------------

  let showMealForm = $state(false);
  let editingMeal = $state<meal.Response | null>(null);
  let savingMeal = $state(false);
  let mealFormError = $state<string | null>(null);

  let deletingMeal = $state<meal.Response | null>(null);
  let deletingMealBusy = $state(false);
  let deletingMealError = $state<string | null>(null);

  function openCreateMeal() {
    editingMeal = null;
    mealFormError = null;
    showMealForm = true;
  }

  function openEditMeal(m: meal.Response) {
    editingMeal = m;
    mealFormError = null;
    showMealForm = true;
  }

  async function handleMealSubmit(data: { name: string; note: string; order: number }) {
    if (!activeDayVariant) return;
    savingMeal = true;
    mealFormError = null;
    try {
      if (editingMeal) {
        await mealApi.update({ id: editingMeal.id, ...data });
      } else {
        await mealApi.create({ dayVariantId: activeDayVariant.id, ...data });
      }
      showMealForm = false;
      await loadPlan();
    } catch (err) {
      mealFormError = errorMessage(err);
    } finally {
      savingMeal = false;
    }
  }

  function openDeleteMeal(m: meal.Response) {
    deletingMeal = m;
    deletingMealError = null;
  }

  async function confirmDeleteMeal() {
    if (!deletingMeal) return;
    deletingMealBusy = true;
    deletingMealError = null;
    try {
      await mealApi.remove(deletingMeal.id);
      deletingMeal = null;
      await loadPlan();
    } catch (err) {
      deletingMealError = errorMessage(err);
    } finally {
      deletingMealBusy = false;
    }
  }

  // --- Food CRUD -----------------------------------------------------------

  let showFoodForm = $state(false);
  let editingFood = $state<food.Response | null>(null);
  let activeMealForFood = $state<meal.Response | null>(null);
  let savingFood = $state(false);
  let foodFormError = $state<string | null>(null);

  let deletingFood = $state<food.Response | null>(null);
  let deletingFoodBusy = $state(false);
  let deletingFoodError = $state<string | null>(null);

  function openCreateFood(m: meal.Response) {
    activeMealForFood = m;
    editingFood = null;
    foodFormError = null;
    showFoodForm = true;
  }

  function openEditFood(m: meal.Response, f: food.Response) {
    activeMealForFood = m;
    editingFood = f;
    foodFormError = null;
    showFoodForm = true;
  }

  async function handleFoodSubmit(data: {
    baseFood: string;
    baseQuantity: number;
    baseUnit: string;
    alternativeFood?: string;
    alternativeQuantity?: number;
    alternativeUnit?: string;
    order: number;
  }) {
    if (!activeMealForFood) return;
    savingFood = true;
    foodFormError = null;
    try {
      if (editingFood) {
        await foodApi.update({ id: editingFood.id, ...data });
      } else {
        await foodApi.create({ mealId: activeMealForFood.id, ...data });
      }
      showFoodForm = false;
      await loadPlan();
    } catch (err) {
      foodFormError = errorMessage(err);
    } finally {
      savingFood = false;
    }
  }

  function openDeleteFood(f: food.Response) {
    deletingFood = f;
    deletingFoodError = null;
  }

  async function confirmDeleteFood() {
    if (!deletingFood) return;
    deletingFoodBusy = true;
    deletingFoodError = null;
    try {
      await foodApi.remove(deletingFood.id);
      deletingFood = null;
      await loadPlan();
    } catch (err) {
      deletingFoodError = errorMessage(err);
    } finally {
      deletingFoodBusy = false;
    }
  }

  // --- Supplement CRUD -------------------------------------------------------

  let showSupplementForm = $state(false);
  let editingSupplement = $state<supplement.Response | null>(null);
  let activeMealForSupplement = $state<meal.Response | null>(null);
  let savingSupplement = $state(false);
  let supplementFormError = $state<string | null>(null);

  let deletingSupplement = $state<supplement.Response | null>(null);
  let deletingSupplementBusy = $state(false);
  let deletingSupplementError = $state<string | null>(null);

  function openCreateSupplement(m: meal.Response) {
    activeMealForSupplement = m;
    editingSupplement = null;
    supplementFormError = null;
    showSupplementForm = true;
  }

  function openEditSupplement(m: meal.Response, s: supplement.Response) {
    activeMealForSupplement = m;
    editingSupplement = s;
    supplementFormError = null;
    showSupplementForm = true;
  }

  async function handleSupplementSubmit(data: { supplement: string; quantity: number; order: number }) {
    if (!activeMealForSupplement) return;
    savingSupplement = true;
    supplementFormError = null;
    try {
      if (editingSupplement) {
        await supplementApi.update({ id: editingSupplement.id, ...data });
      } else {
        await supplementApi.create({ mealId: activeMealForSupplement.id, ...data });
      }
      showSupplementForm = false;
      await loadPlan();
    } catch (err) {
      supplementFormError = errorMessage(err);
    } finally {
      savingSupplement = false;
    }
  }

  function openDeleteSupplement(s: supplement.Response) {
    deletingSupplement = s;
    deletingSupplementError = null;
  }

  async function confirmDeleteSupplement() {
    if (!deletingSupplement) return;
    deletingSupplementBusy = true;
    deletingSupplementError = null;
    try {
      await supplementApi.remove(deletingSupplement.id);
      deletingSupplement = null;
      await loadPlan();
    } catch (err) {
      deletingSupplementError = errorMessage(err);
    } finally {
      deletingSupplementBusy = false;
    }
  }
</script>

{#if !planId}
  <div class="mx-auto max-w-5xl px-6 py-10">
    <p class="text-sm text-muted-foreground">{$t("plan.noneSelected")}</p>
    <button
      onclick={backToClients}
      class="mt-4 rounded-md border border-border px-4 py-2 text-sm font-medium transition-colors hover:bg-accent"
    >
      {$t("plan.backToClients")}
    </button>
  </div>
{:else if loading}
  <div class="mx-auto max-w-5xl px-6 py-10">
    <p class="text-sm text-muted-foreground">{$t("plan.loading")}</p>
  </div>
{:else if loadError}
  <div class="mx-auto max-w-5xl px-6 py-10">
    <p class="text-sm text-destructive">{loadError}</p>
    <button
      onclick={backToClients}
      class="mt-4 rounded-md border border-border px-4 py-2 text-sm font-medium transition-colors hover:bg-accent"
    >
      {$t("plan.backToClients")}
    </button>
  </div>
{:else if planFull}
  <div class="mx-auto max-w-5xl px-6 py-10">
    <button onclick={backToClients} class="text-xs text-muted-foreground hover:underline">{$t("plan.backToClients")}</button>

    <header class="mt-1 flex flex-wrap items-end justify-between gap-4 border-b border-border pb-6">
      <div>
        <h1 class="font-display text-3xl font-semibold tracking-tight">{planFull.title}</h1>
        {#if planFull.generalRecommendations}
          <p class="mt-1 text-sm text-muted-foreground">{planFull.generalRecommendations}</p>
        {/if}
      </div>
      <div class="flex gap-2">
        <button
          onclick={handlePreviewPDF}
          disabled={previewing}
          class="rounded-md border border-border px-4 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent disabled:opacity-60"
        >
          {previewing ? $t("plan.previewing") : $t("plan.preview")}
        </button>
        <button
          onclick={handleExportPDF}
          disabled={exporting}
          class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-60"
        >
          {exporting ? $t("plan.exporting") : $t("plan.exportPdf")}
        </button>
      </div>
    </header>

    {#if previewError}
      <p class="mt-3 text-sm text-destructive">{previewError}</p>
    {:else if previewSuccess}
      <p class="mt-3 text-sm text-muted-foreground">
        {$t("plan.previewOpened")}
        <span class="font-medium">{previewSuccess.path}</span>
      </p>
    {/if}

    {#if exportSuccess}
      <p class="mt-3 text-sm text-primary">
        {$t("plan.pdfSaved")} <span class="font-medium">{exportSuccess.fileName}</span>
        <span class="text-muted-foreground">— {exportSuccess.path}</span>
        <button onclick={handleOpenFolder} class="ml-2 font-medium text-primary hover:underline">
          {$t("plan.openFolder")}
        </button>
      </p>
    {:else if exportError}
      <p class="mt-3 text-sm text-destructive">{exportError}</p>
    {/if}

    <div class="mt-6 flex flex-wrap items-center gap-2">
      {#each planFull.dayVariants ?? [] as dv (dv.id)}
        <button
          onclick={() => (activeDayVariantId = dv.id)}
          class="rounded-full px-4 py-1.5 text-sm font-medium transition-colors {dv.id === activeDayVariantId
            ? 'bg-primary text-primary-foreground'
            : 'border border-border text-muted-foreground hover:bg-accent hover:text-accent-foreground'}"
        >
          {dv.name}
        </button>
      {/each}
      <button
        onclick={openCreateDayVariant}
        class="rounded-full border border-dashed border-border px-4 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-accent"
      >
        {$t("plan.addDayVariant")}
      </button>
    </div>

    {#if activeDayVariant}
      <div class="mt-3 flex gap-3">
        <button onclick={() => openEditDayVariant(activeDayVariant!)} class="text-xs text-muted-foreground hover:underline">
          {$t("plan.editDay")}
        </button>
        <button onclick={() => openDeleteDayVariant(activeDayVariant!)} class="text-xs text-destructive hover:underline">
          {$t("plan.deleteDay")}
        </button>
      </div>

      <div class="mt-6 space-y-5">
        {#each activeDayVariant.meals ?? [] as m (m.id)}
          <section class="rounded-xl border border-border bg-card p-5">
            <div class="flex flex-wrap items-baseline justify-between gap-2">
              <div>
                <h2 class="font-display text-lg font-semibold">{m.name}</h2>
                {#if m.note}
                  <p class="text-xs text-muted-foreground">{m.note}</p>
                {/if}
              </div>
              <div class="flex flex-wrap gap-3">
                <button onclick={() => openCreateFood(m)} class="text-xs font-medium text-primary hover:underline">
                  {$t("plan.addFood")}
                </button>
                <button onclick={() => openCreateSupplement(m)} class="text-xs font-medium text-primary hover:underline">
                  {$t("plan.addSupplement")}
                </button>
                <button onclick={() => openEditMeal(m)} class="text-xs text-muted-foreground hover:underline">{$t("common.edit")}</button>
                <button onclick={() => openDeleteMeal(m)} class="text-xs text-destructive hover:underline">{$t("common.delete")}</button>
              </div>
            </div>

            {#if m.foods && m.foods.length > 0}
              <table class="mt-4 w-full text-sm">
                <thead>
                  <tr class="text-left text-xs uppercase tracking-wide text-muted-foreground">
                    <th class="pb-2 font-medium">{$t("plan.tableFood")}</th>
                    <th class="pb-2 font-medium">{$t("plan.tableQuantity")}</th>
                    <th class="pb-2 font-medium">{$t("plan.tableAlternative")}</th>
                    <th class="pb-2 font-medium">{$t("plan.tableQuantity")}</th>
                    <th class="pb-2 font-medium"></th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-border">
                  {#each m.foods as f (f.id)}
                    <tr>
                      <td class="py-2.5 font-medium">{f.baseFood}</td>
                      <td class="py-2.5 text-muted-foreground">{f.baseQuantity} {f.baseUnit}</td>
                      <td class="py-2.5">{f.alternativeFood ?? $t("common.na")}</td>
                      <td class="py-2.5 text-muted-foreground">
                        {f.alternativeFood ? `${f.alternativeQuantity} ${f.alternativeUnit}` : $t("common.na")}
                      </td>
                      <td class="py-2.5 text-right whitespace-nowrap">
                        <button
                          onclick={() => openEditFood(m, f)}
                          aria-label={$t("plan.editFoodAria")}
                          class="rounded px-1.5 py-0.5 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
                        >
                          ✎
                        </button>
                        <button
                          onclick={() => openDeleteFood(f)}
                          aria-label={$t("plan.deleteFoodAria")}
                          class="rounded px-1.5 py-0.5 text-muted-foreground transition-colors hover:bg-destructive hover:text-destructive-foreground"
                        >
                          🗑
                        </button>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            {:else}
              <p class="mt-3 text-sm text-muted-foreground">{$t("plan.noFoods")}</p>
            {/if}

            {#if m.supplements && m.supplements.length > 0}
              <div class="mt-4 flex flex-wrap items-center gap-2 border-t border-border pt-4">
                <span class="text-xs uppercase tracking-wide text-muted-foreground">{$t("plan.supplementsLabel")}</span>
                {#each m.supplements as s (s.id)}
                  <span class="flex items-center gap-1 rounded-full bg-secondary px-3 py-1 text-xs text-secondary-foreground">
                    {s.supplement} × {s.quantity}
                    <button
                      onclick={() => openEditSupplement(m, s)}
                      aria-label={$t("plan.editSupplementAria")}
                      class="opacity-60 transition-opacity hover:opacity-100"
                    >
                      ✎
                    </button>
                    <button
                      onclick={() => openDeleteSupplement(s)}
                      aria-label={$t("plan.deleteSupplementAria")}
                      class="opacity-60 transition-opacity hover:opacity-100"
                    >
                      ✕
                    </button>
                  </span>
                {/each}
              </div>
            {/if}
          </section>
        {/each}

        {#if !activeDayVariant.meals || activeDayVariant.meals.length === 0}
          <div class="rounded-xl border border-dashed border-border p-10 text-center text-sm text-muted-foreground">
            {$t("plan.noMealsFor", { name: activeDayVariant.name })}
          </div>
        {/if}

        <button
          onclick={openCreateMeal}
          class="w-full rounded-xl border border-dashed border-border py-3 text-sm font-medium text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
        >
          {$t("plan.addMeal")}
        </button>
      </div>
    {:else}
      <div class="mt-8 rounded-xl border border-dashed border-border p-10 text-center text-sm text-muted-foreground">
        {$t("plan.noDayVariants")}
      </div>
    {/if}

    <section class="mt-8 rounded-xl border border-border bg-card p-5">
      <h3 class="font-display font-semibold">{$t("plan.footerTitle")}</h3>
      <p class="mt-2 text-sm text-muted-foreground">{planFull.footer}</p>
    </section>
  </div>
{/if}

<Modal
  open={showDayVariantForm}
  title={editingDayVariant ? $t("plan.editDayVariant") : $t("plan.newDayVariant")}
  onClose={() => (showDayVariantForm = false)}
>
  <DayVariantForm
    initial={editingDayVariant}
    submitting={savingDayVariant}
    error={dayVariantFormError}
    onCancel={() => (showDayVariantForm = false)}
    onSubmit={handleDayVariantSubmit}
  />
</Modal>
<ConfirmDialog
  open={deletingDayVariant !== null}
  message={deletingDayVariant ? $t("plan.deleteDayVariantConfirm", { name: deletingDayVariant.name }) : ""}
  submitting={deletingDayVariantBusy}
  error={deletingDayVariantError}
  onCancel={() => (deletingDayVariant = null)}
  onConfirm={confirmDeleteDayVariant}
/>

<Modal open={showMealForm} title={editingMeal ? $t("plan.editMeal") : $t("plan.newMeal")} onClose={() => (showMealForm = false)}>
  <MealForm
    initial={editingMeal}
    submitting={savingMeal}
    error={mealFormError}
    onCancel={() => (showMealForm = false)}
    onSubmit={handleMealSubmit}
  />
</Modal>
<ConfirmDialog
  open={deletingMeal !== null}
  message={deletingMeal ? $t("plan.deleteMealConfirm", { name: deletingMeal.name }) : ""}
  submitting={deletingMealBusy}
  error={deletingMealError}
  onCancel={() => (deletingMeal = null)}
  onConfirm={confirmDeleteMeal}
/>

<Modal open={showFoodForm} title={editingFood ? $t("plan.editFood") : $t("plan.newFood")} onClose={() => (showFoodForm = false)}>
  <FoodForm
    initial={editingFood}
    submitting={savingFood}
    error={foodFormError}
    onCancel={() => (showFoodForm = false)}
    onSubmit={handleFoodSubmit}
  />
</Modal>
<ConfirmDialog
  open={deletingFood !== null}
  message={deletingFood ? $t("plan.deleteFoodConfirm", { name: deletingFood.baseFood }) : ""}
  submitting={deletingFoodBusy}
  error={deletingFoodError}
  onCancel={() => (deletingFood = null)}
  onConfirm={confirmDeleteFood}
/>

<Modal
  open={showSupplementForm}
  title={editingSupplement ? $t("plan.editSupplement") : $t("plan.newSupplement")}
  onClose={() => (showSupplementForm = false)}
>
  <SupplementForm
    initial={editingSupplement}
    submitting={savingSupplement}
    error={supplementFormError}
    onCancel={() => (showSupplementForm = false)}
    onSubmit={handleSupplementSubmit}
  />
</Modal>
<ConfirmDialog
  open={deletingSupplement !== null}
  message={deletingSupplement ? $t("plan.deleteSupplementConfirm", { name: deletingSupplement.supplement }) : ""}
  submitting={deletingSupplementBusy}
  error={deletingSupplementError}
  onCancel={() => (deletingSupplement = null)}
  onConfirm={confirmDeleteSupplement}
/>
