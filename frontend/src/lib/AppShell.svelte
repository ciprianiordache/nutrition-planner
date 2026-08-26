<script lang="ts">
  import type { Snippet } from "svelte";
  import TitleBar from "./TitleBar.svelte";
  import Logo from "./Logo.svelte";
  import { currentView } from "./stores";
  import type { View } from "./types";

  let { children }: { children: Snippet } = $props();

  const nav: { to: View; label: string }[] = [
    { to: "clients", label: "Clienți" },
    { to: "plan", label: "Editor plan" },
  ];
</script>

<div class="flex h-screen flex-col overflow-hidden bg-background text-foreground">
  <TitleBar />
  <div class="flex flex-1 min-h-0">
    <aside class="hidden w-60 shrink-0 flex-col overflow-y-auto border-r border-border bg-sidebar px-5 py-7 md:flex">      <div class="mb-9 flex items-center gap-3">
        <Logo class="h-15 w-15 shrink-0" />
        <div>
          <span class="font-display text-lg font-semibold tracking-tight">NutriPlan</span>
          <p class="mt-0.5 text-xs text-muted-foreground">Planuri nutriționale</p>
        </div>
      </div>
      <nav class="flex flex-col gap-1">
        {#each nav as item (item.to)}
          <button
            onclick={() => currentView.set(item.to)}
            class="rounded-md px-3 py-2 text-left text-sm font-medium transition-colors {$currentView === item.to
              ? 'bg-accent text-accent-foreground'
              : 'text-muted-foreground hover:bg-accent/60 hover:text-foreground'}"
          >
            {item.label}
          </button>
        {/each}
      </nav>
      <div class="mt-auto rounded-lg border border-border bg-card p-3">
        <p class="text-xs font-medium">Mihai Bucur</p>
        <p class="text-xs text-muted-foreground">Antrenor Personal</p>
      </div>
    </aside>
    <main class="min-h-0 flex-1 overflow-y-auto overscroll-none">{@render children()}</main>
  </div>
</div>
