<script lang="ts">
  import { onMount } from "svelte";
  import AppShell from "./lib/AppShell.svelte";
  import ClientsView from "./views/ClientsView.svelte";
  import PlanView from "./views/PlanView.svelte";
  import SettingsView from "./views/SettingsView.svelte";
  import OnboardingModal from "./lib/OnboardingModal.svelte";
  import { currentView } from "./lib/stores";
  import { profile, profileLoaded, loadProfile } from "./lib/profile";

  onMount(loadProfile);

  // Onboarding = "we've checked the backend and there's genuinely no name
  // saved yet" - profileLoaded guards against briefly showing the modal
  // before the initial settingsApi.get() call has even resolved.
  const needsOnboarding = $derived($profileLoaded && !$profile.name.trim());
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

{#if needsOnboarding}
  <OnboardingModal />
{/if}
