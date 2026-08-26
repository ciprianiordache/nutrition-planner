<script lang="ts">
  import { onMount } from "svelte";
  import type { client, plan } from "../../wailsjs/go/models";
  import { clientApi, planApi, errorMessage } from "../lib/api";
  import Modal from "../lib/Modal.svelte";
  import ConfirmDialog from "../lib/ConfirmDialog.svelte";
  import ClientForm from "../lib/ClientForm.svelte";
  import PlanForm from "../lib/PlanForm.svelte";
  import { currentView, currentPlanId } from "../lib/stores";

  let clients = $state<client.Response[]>([]);
  let loadingClients = $state(true);
  let clientsError = $state<string | null>(null);

  let selectedClientId = $state<string | null>(null);
  let plans = $state<plan.Response[]>([]);
  let loadingPlans = $state(false);
  let plansError = $state<string | null>(null);

  // Client — creare/editare/ștergere
  let showClientForm = $state(false);
  let editingClient = $state<client.Response | null>(null);
  let savingClient = $state(false);
  let clientFormError = $state<string | null>(null);

  let deletingClient = $state<client.Response | null>(null);
  let deletingClientBusy = $state(false);
  let deletingClientError = $state<string | null>(null);

  // Plan — creare/editare/ștergere
  let showPlanForm = $state(false);
  let editingPlan = $state<plan.Response | null>(null);
  let savingPlan = $state(false);
  let planFormError = $state<string | null>(null);

  let deletingPlan = $state<plan.Response | null>(null);
  let deletingPlanBusy = $state(false);
  let deletingPlanError = $state<string | null>(null);

  const selectedClient = $derived(clients.find((c) => c.id === selectedClientId) ?? null);

  // Căutare + filtrare pe obiectiv — client-side, aplicate peste lista deja
  // încărcată (nu necesită un nou apel către backend).
  let searchQuery = $state("");
  let targetFilter = $state("");

  const filteredClients = $derived(
    clients.filter((c) => {
      const matchesSearch = c.name.toLowerCase().includes(searchQuery.trim().toLowerCase());
      const matchesTarget = targetFilter === "" || c.target === targetFilter;
      return matchesSearch && matchesTarget;
    })
  );

  onMount(loadClients);

  async function loadClients() {
    loadingClients = true;
    clientsError = null;
    try {
      clients = (await clientApi.list(0, 100)) ?? [];
    } catch (err) {
      clientsError = errorMessage(err);
    } finally {
      loadingClients = false;
    }
  }

  async function loadPlansForClient(id: string) {
    loadingPlans = true;
    plansError = null;
    try {
      plans = (await planApi.listByClient(id)) ?? [];
    } catch (err) {
      plansError = errorMessage(err);
      plans = [];
    } finally {
      loadingPlans = false;
    }
  }

  async function selectClient(id: string) {
    selectedClientId = selectedClientId === id ? null : id;
    if (selectedClientId === null) return;
    await loadPlansForClient(selectedClientId);
  }

  function initials(name: string): string {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("");
  }

  // --- Client CRUD -----------------------------------------------------

  function openCreateClient() {
    editingClient = null;
    clientFormError = null;
    showClientForm = true;
  }

  function openEditClient(c: client.Response, e: MouseEvent) {
    e.stopPropagation();
    editingClient = c;
    clientFormError = null;
    showClientForm = true;
  }

  async function handleClientSubmit(data: { name: string; weight: number; height: number; target: string }) {
    savingClient = true;
    clientFormError = null;
    try {
      if (editingClient) {
        await clientApi.update({ id: editingClient.id, ...data });
      } else {
        await clientApi.create(data);
      }
      showClientForm = false;
      await loadClients();
    } catch (err) {
      clientFormError = errorMessage(err);
    } finally {
      savingClient = false;
    }
  }

  function openDeleteClient(c: client.Response, e: MouseEvent) {
    e.stopPropagation();
    deletingClient = c;
    deletingClientError = null;
  }

  async function confirmDeleteClient() {
    if (!deletingClient) return;
    deletingClientBusy = true;
    deletingClientError = null;
    try {
      await clientApi.remove(deletingClient.id);
      if (selectedClientId === deletingClient.id) selectedClientId = null;
      deletingClient = null;
      await loadClients();
    } catch (err) {
      deletingClientError = errorMessage(err);
    } finally {
      deletingClientBusy = false;
    }
  }

  // --- Plan CRUD ---------------------------------------------------------

  function openCreatePlan() {
    editingPlan = null;
    planFormError = null;
    showPlanForm = true;
  }

  function openEditPlan(p: plan.Response) {
    editingPlan = p;
    planFormError = null;
    showPlanForm = true;
  }

  async function handlePlanSubmit(data: { title: string; generalRecommendations: string; footer: string }) {
    if (!selectedClientId) return;
    savingPlan = true;
    planFormError = null;
    try {
      if (editingPlan) {
        await planApi.update({ id: editingPlan.id, ...data });
      } else {
        await planApi.create({ clientId: selectedClientId, ...data });
      }
      showPlanForm = false;
      await loadPlansForClient(selectedClientId);
    } catch (err) {
      planFormError = errorMessage(err);
    } finally {
      savingPlan = false;
    }
  }

  function openDeletePlan(p: plan.Response) {
    deletingPlan = p;
    deletingPlanError = null;
  }

  async function confirmDeletePlan() {
    if (!deletingPlan || !selectedClientId) return;
    deletingPlanBusy = true;
    deletingPlanError = null;
    try {
      await planApi.remove(deletingPlan.id);
      deletingPlan = null;
      await loadPlansForClient(selectedClientId);
    } catch (err) {
      deletingPlanError = errorMessage(err);
    } finally {
      deletingPlanBusy = false;
    }
  }

  function openPlan(planId: string) {
    currentPlanId.set(planId);
    currentView.set("plan");
  }
</script>

<div class="mx-auto max-w-5xl px-6 py-10">
  <header class="flex flex-wrap items-end justify-between gap-4">
    <div>
      <h1 class="font-display text-3xl font-semibold tracking-tight">Clienți</h1>
      <p class="mt-1 text-sm text-muted-foreground">
        {filteredClients.length} clienți · alege un client pentru a-i vedea planurile
      </p>
    </div>
    <button
      onclick={openCreateClient}
      class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
    >
      + Client nou
    </button>
  </header>

   <div class="mt-8 flex gap-3">
    <input
      bind:value={searchQuery}
      placeholder="Caută după nume…"
      class="w-full max-w-sm rounded-md border border-input bg-card px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus:ring-2 focus:ring-ring/40"
    />
    <select bind:value={targetFilter} class="rounded-md border border-input bg-card px-3 py-2 text-sm">
      <option value="">Toate obiectivele</option>
      <option value="Slăbire">Slăbire</option>
      <option value="Masă musculară">Masă musculară</option>
      <option value="Menținere">Menținere</option>
      <option value="Recompoziție">Recompoziție</option>
    </select>
  </div>

  {#if loadingClients}
    <p class="mt-8 text-sm text-muted-foreground">Se încarcă clienții…</p>
  {:else if clientsError}
    <p class="mt-8 text-sm text-destructive">{clientsError}</p>
  {:else if clients.length === 0}
    <div class="mt-8 rounded-xl border border-dashed border-border p-10 text-center text-sm text-muted-foreground">
      Niciun client încă. Adaugă primul client ca să începi.
    </div>
  {:else if filteredClients.length === 0}
    <div class="mt-8 rounded-xl border border-dashed border-border p-10 text-center text-sm text-muted-foreground">
      Niciun client nu corespunde căutării sau filtrului selectat.
    </div>
  {:else}
    <ul class="mt-6 divide-y divide-border overflow-hidden rounded-xl border border-border bg-card">
      {#each filteredClients as c (c.id)}
        <li>
          <div
            class="flex w-full flex-wrap items-center gap-4 px-5 py-4 transition-colors hover:bg-muted/60 {selectedClientId ===
            c.id
              ? 'bg-accent/50'
              : ''}"
          >
            <button onclick={() => selectClient(c.id)} class="flex min-w-0 flex-1 items-center gap-4 text-left">
              <div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-accent font-display text-sm font-semibold text-accent-foreground">
                {initials(c.name)}
              </div>
              <div class="min-w-40 flex-1">
                <p class="font-medium">{c.name}</p>
                <p class="text-xs text-muted-foreground">{c.weight} kg · {c.height} cm</p>
              </div>
              {#if c.target}
                <span class="rounded-full bg-secondary px-3 py-1 text-xs font-medium text-secondary-foreground">
                  {c.target}
                </span>
              {/if}
            </button>
            <span class="flex shrink-0 gap-1">
              <button
                onclick={(e) => openEditClient(c, e)}
                aria-label="Editează client"
                class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
              >
                ✎
              </button>
              <button
                onclick={(e) => openDeleteClient(c, e)}
                aria-label="Șterge client"
                class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-destructive hover:text-destructive-foreground"
              >
                🗑
              </button>
            </span>
          </div>

          {#if selectedClientId === c.id}
            <div class="border-t border-border bg-muted/30 px-5 py-5">
              <div class="flex flex-wrap items-center justify-between gap-3">
                <h3 class="font-display text-sm font-semibold">Planurile lui {c.name}</h3>
                <button
                  onclick={openCreatePlan}
                  class="rounded-md border border-border bg-card px-3 py-1.5 text-xs font-medium transition-colors hover:bg-accent"
                >
                  + Plan nou
                </button>
              </div>

              {#if loadingPlans}
                <p class="mt-3 text-sm text-muted-foreground">Se încarcă planurile…</p>
              {:else if plansError}
                <p class="mt-3 text-sm text-destructive">{plansError}</p>
              {:else if plans.length === 0}
                <p class="mt-3 text-sm text-muted-foreground">Niciun plan creat încă pentru acest client.</p>
              {:else}
                <ul class="mt-3 space-y-2">
                  {#each plans as p (p.id)}
                    <li class="flex flex-wrap items-center justify-between gap-3 rounded-lg border border-border bg-card px-4 py-3">
                      <div>
                        <p class="text-sm font-medium">{p.title}</p>
                        {#if p.generalRecommendations}
                          <p class="mt-0.5 text-xs text-muted-foreground">{p.generalRecommendations}</p>
                        {/if}
                      </div>
                      <div class="flex gap-1.5">
                        <button
                          onclick={() => openPlan(p.id)}
                          class="rounded-md border border-border px-3 py-1.5 text-xs font-medium transition-colors hover:bg-accent"
                        >
                          Deschide
                        </button>
                        <button
                          onclick={() => openEditPlan(p)}
                          aria-label="Editează plan"
                          class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
                        >
                          ✎
                        </button>
                        <button
                          onclick={() => openDeletePlan(p)}
                          aria-label="Șterge plan"
                          class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-destructive hover:text-destructive-foreground"
                        >
                          🗑
                        </button>
                      </div>
                    </li>
                  {/each}
                </ul>
              {/if}
            </div>
          {/if}
        </li>
      {/each}
    </ul>
  {/if}
</div>

<Modal open={showClientForm} title={editingClient ? "Editează client" : "Client nou"} onClose={() => (showClientForm = false)}>
  <ClientForm
    initial={editingClient}
    submitting={savingClient}
    error={clientFormError}
    onCancel={() => (showClientForm = false)}
    onSubmit={handleClientSubmit}
  />
</Modal>

<ConfirmDialog
  open={deletingClient !== null}
  message={deletingClient
    ? `Sigur vrei să ștergi clientul „${deletingClient.name}”? Toate planurile lui vor fi șterse.`
    : ""}
  submitting={deletingClientBusy}
  error={deletingClientError}
  onCancel={() => (deletingClient = null)}
  onConfirm={confirmDeleteClient}
/>

<Modal
  open={showPlanForm}
  title={editingPlan ? "Editează plan" : `Plan nou${selectedClient ? ` — ${selectedClient.name}` : ""}`}
  onClose={() => (showPlanForm = false)}
>
  <PlanForm
    initial={editingPlan}
    submitting={savingPlan}
    error={planFormError}
    onCancel={() => (showPlanForm = false)}
    onSubmit={handlePlanSubmit}
  />
</Modal>

<ConfirmDialog
  open={deletingPlan !== null}
  message={deletingPlan ? `Sigur vrei să ștergi planul „${deletingPlan.title}”?` : ""}
  submitting={deletingPlanBusy}
  error={deletingPlanError}
  onCancel={() => (deletingPlan = null)}
  onConfirm={confirmDeletePlan}
/>
