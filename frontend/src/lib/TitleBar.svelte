<script lang="ts">
  import { Minus, Square, X, Copy } from "@lucide/svelte";
  import Icon from "./Icon.svelte";
  import { WindowMinimise, WindowToggleMaximise, Quit } from "../../wailsjs/runtime/runtime";
  import { t } from "./i18n";

  let maximized = $state(false);

  function minimize() {
    WindowMinimise();
  }

  function toggleMaximize() {
    WindowToggleMaximise();
    maximized = !maximized;
  }

  function close() {
    Quit();
  }
</script>

<div
  style="--wails-draggable: drag"
  class="flex h-10 shrink-0 select-none items-center justify-between border-b border-border bg-sidebar pl-4"
>
  <div class="flex items-center gap-2.5">
    <Icon class="h-5 w-5" />
    <span class="font-display text-xs font-medium tracking-wide text-muted-foreground">
      {$t("titlebar.appName")}
    </span>
  </div>
  <div class="flex h-full" style="--wails-draggable: no-drag">
    <button
      aria-label={$t("titlebar.minimize")}
      onclick={minimize}
      class="flex h-full w-11 items-center justify-center text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
    >
      <Minus class="size-4" />
    </button>
    <button
      aria-label={maximized ? $t("titlebar.restore") : $t("titlebar.maximize")}
      onclick={toggleMaximize}
      class="flex h-full w-11 items-center justify-center text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
    >
      {#if maximized}
        <Copy class="size-3.5" />
      {:else}
        <Square class="size-3.5" />
      {/if}
    </button>
    <button
      aria-label={$t("titlebar.close")}
      onclick={close}
      class="flex h-full w-11 items-center justify-center text-muted-foreground transition-colors hover:bg-destructive hover:text-destructive-foreground"
    >
      <X class="size-4" />
    </button>
  </div>
</div>
