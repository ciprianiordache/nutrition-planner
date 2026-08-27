<script lang="ts">
  import AppShell from "./lib/AppShell.svelte";
  import ClientsView from "./views/ClientsView.svelte";
  import PlanView from "./views/PlanView.svelte";
  import SettingsView from "./views/SettingsView.svelte";
  import { currentView } from "./lib/stores";
</script>

<AppShell>
  {#if $currentView === "clients"}
    <ClientsView />
  {:else if $currentView === "settings"}
    <SettingsView />
  {:else}
    <PlanView />
  {/if}
</AppShell>
