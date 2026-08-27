<script lang="ts">
  import { onMount } from "svelte";
  import type { settings } from "../../wailsjs/go/models";
  import { settingsApi, errorMessage } from "../lib/api";
  import { t, locale, setLocale, type TranslationKey } from "../lib/i18n";
  import { uiTheme, setUITheme, type UITheme } from "../lib/theme";
  import { profile } from "../lib/profile";

  // PDF_DEFAULTS mirrors config.GeneratorConfig's own defaults (config/
  // struct.go) - shown in the color pickers whenever the corresponding
  // settings field is still empty, so the picker always reflects the
  // color that's actually in effect right now, not a blank/black square.
  const PDF_DEFAULTS: Record<ColorKey, string> = {
    colorPageBackground: "F7F5EF",
    colorCardBackground: "FFFFFF",
    colorCardBorder: "E7E4D9",
    colorPrimary: "2F5233",
    colorPrimaryLight: "E3EDE1",
    colorOnPrimary: "FFFFFF",
    colorText: "2B2B28",
    colorMuted: "6B6B66",
    colorTableHeader: "8A8A82",
    colorDivider: "E5E2D8",
  };

  type ColorKey =
    | "colorPrimary"
    | "colorPrimaryLight"
    | "colorOnPrimary"
    | "colorPageBackground"
    | "colorCardBackground"
    | "colorCardBorder"
    | "colorText"
    | "colorMuted"
    | "colorTableHeader"
    | "colorDivider";

  const COLOR_FIELDS: { key: ColorKey; labelKey: TranslationKey }[] = [
    { key: "colorPrimary", labelKey: "settings.colorPrimary" },
    { key: "colorPrimaryLight", labelKey: "settings.colorPrimaryLight" },
    { key: "colorOnPrimary", labelKey: "settings.colorOnPrimary" },
    { key: "colorPageBackground", labelKey: "settings.colorPageBackground" },
    { key: "colorCardBackground", labelKey: "settings.colorCardBackground" },
    { key: "colorCardBorder", labelKey: "settings.colorCardBorder" },
    { key: "colorText", labelKey: "settings.colorText" },
    { key: "colorMuted", labelKey: "settings.colorMuted" },
    { key: "colorTableHeader", labelKey: "settings.colorTableHeader" },
    { key: "colorDivider", labelKey: "settings.colorDivider" },
  ];

  /** hexToColorInput/colorInputToHex convert between the backend's bare
   * hex format ("2F5233", no "#") and the format <input type="color">
   * requires ("#2f5233"). */
  function hexToColorInput(hex: string): string {
    return `#${hex}`;
  }
  function colorInputToHex(value: string): string {
    return value.replace(/^#/, "").toUpperCase();
  }

  let loading = $state(true);
  let loadError = $state<string | null>(null);
  let saving = $state(false);
  let saveError = $state<string | null>(null);
  let saved = $state(false);

  let exportFolder = $state("");
  let trainerName = $state("");
  let trainerRole = $state("");
  let colors = $state<Record<ColorKey, string>>({ ...PDF_DEFAULTS });

  onMount(load);

  async function load() {
    loading = true;
    loadError = null;
    try {
      const s = await settingsApi.get();
      applyLoadedSettings(s);
    } catch (err) {
      loadError = errorMessage(err);
    } finally {
      loading = false;
    }
  }

  function applyLoadedSettings(s: settings.Response) {
    exportFolder = s.exportBaseDirectory ?? "";
    trainerName = s.trainerName ?? "";
    trainerRole = s.trainerRole ?? "";
    for (const { key } of COLOR_FIELDS) {
      colors[key] = s[key] || PDF_DEFAULTS[key];
    }
    // Reconcile the language/theme stores with whatever was last saved to
    // the backend - matters the first time Settings loads on a machine
    // where localStorage is empty but a previous session already saved a
    // preference (e.g. after a reinstall, or a second Windows profile).
    if (s.language === "ro" || s.language === "en") locale.set(s.language);
    if (s.uiTheme === "light" || s.uiTheme === "dark" || s.uiTheme === "system") uiTheme.set(s.uiTheme);
  }

  async function chooseFolder() {
    try {
      const path = await settingsApi.chooseExportFolder();
      if (path) exportFolder = path;
    } catch (err) {
      saveError = errorMessage(err);
    }
  }

  function resetColors() {
    colors = { ...PDF_DEFAULTS };
  }

  async function save() {
    saving = true;
    saveError = null;
    saved = false;
    try {
      const req: settings.UpdateRequest = {
        exportBaseDirectory: exportFolder,
        trainerName,
        trainerRole,
      };
      for (const { key } of COLOR_FIELDS) {
        req[key] = colors[key];
      }
      await settingsApi.update(req);
      profile.set({ name: trainerName, role: trainerRole });
      saved = true;
    } catch (err) {
      saveError = errorMessage(err);
    } finally {
      saving = false;
    }
  }

  const themeOptions: { value: UITheme; labelKey: TranslationKey }[] = [
    { value: "light", labelKey: "settings.themeLight" },
    { value: "dark", labelKey: "settings.themeDark" },
    { value: "system", labelKey: "settings.themeSystem" },
  ];

  // --- Backup ------------------------------------------------------------

  let backingUp = $state(false);
  let backupError = $state<string | null>(null);
  let backupPath = $state<string | null>(null);

  async function handleBackup() {
    backingUp = true;
    backupError = null;
    backupPath = null;
    try {
      backupPath = await settingsApi.backupDatabase();
    } catch (err) {
      backupError = errorMessage(err);
    } finally {
      backingUp = false;
    }
  }

  async function handleOpenBackupFolder() {
    if (!backupPath) return;
    try {
      await settingsApi.openBackupFolder(backupPath);
    } catch {
      // Best-effort - the backup itself already succeeded and the path is
      // shown right next to this button either way.
    }
  }
</script>

<div class="mx-auto max-w-3xl px-6 py-10">
  <header>
    <h1 class="font-display text-3xl font-semibold tracking-tight">{$t("settings.title")}</h1>
    <p class="mt-1 text-sm text-muted-foreground">{$t("settings.subtitle")}</p>
  </header>

  {#if loading}
    <p class="mt-8 text-sm text-muted-foreground">{$t("settings.loading")}</p>
  {:else if loadError}
    <p class="mt-8 text-sm text-destructive">{loadError}</p>
  {:else}
    <div class="mt-8 space-y-8">
      <!-- Profile -->
      <section class="rounded-xl border border-border bg-card p-5">
        <h2 class="font-display text-lg font-semibold">{$t("settings.profileSection")}</h2>
        <div class="mt-4 grid grid-cols-2 gap-3">
          <div>
            <label class="mb-1 block text-xs font-medium text-muted-foreground" for="settings-trainer-name">
              {$t("settings.trainerName")}
            </label>
            <input
              id="settings-trainer-name"
              bind:value={trainerName}
              placeholder={$t("settings.trainerNamePlaceholder")}
              class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
            />
          </div>
          <div>
            <label class="mb-1 block text-xs font-medium text-muted-foreground" for="settings-trainer-role">
              {$t("settings.trainerRole")}
            </label>
            <input
              id="settings-trainer-role"
              bind:value={trainerRole}
              placeholder={$t("settings.trainerRolePlaceholder")}
              class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
            />
          </div>
        </div>
      </section>

      <!-- Export -->
      <section class="rounded-xl border border-border bg-card p-5">
        <h2 class="font-display text-lg font-semibold">{$t("settings.exportSection")}</h2>
        <p class="mt-1 text-xs text-muted-foreground">{$t("settings.exportFolderHint")}</p>
        <div class="mt-4 flex gap-2">
          <input
            readonly
            value={exportFolder}
            placeholder={$t("settings.exportFolderDefault")}
            class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm text-muted-foreground outline-none"
          />
          <button
            onclick={chooseFolder}
            class="shrink-0 rounded-md border border-border px-4 py-2 text-sm font-medium transition-colors hover:bg-accent"
          >
            {$t("settings.chooseFolder")}
          </button>
        </div>
      </section>

      <!-- PDF theme -->
      <section class="rounded-xl border border-border bg-card p-5">
        <div class="flex flex-wrap items-start justify-between gap-2">
          <div>
            <h2 class="font-display text-lg font-semibold">{$t("settings.pdfThemeSection")}</h2>
            <p class="mt-1 text-xs text-muted-foreground">{$t("settings.pdfThemeHint")}</p>
          </div>
          <button onclick={resetColors} class="text-xs font-medium text-primary hover:underline">
            {$t("settings.resetColors")}
          </button>
        </div>
        <div class="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3">
          {#each COLOR_FIELDS as field (field.key)}
            <label class="flex items-center gap-2 text-xs">
              <input
                type="color"
                value={hexToColorInput(colors[field.key])}
                oninput={(e) => (colors[field.key] = colorInputToHex(e.currentTarget.value))}
                class="h-8 w-8 shrink-0 cursor-pointer rounded border border-input bg-transparent p-0.5"
              />
              <span class="text-muted-foreground">{$t(field.labelKey)}</span>
            </label>
          {/each}
        </div>
      </section>

      <!-- Appearance -->
      <section class="rounded-xl border border-border bg-card p-5">
        <h2 class="font-display text-lg font-semibold">{$t("settings.appearanceSection")}</h2>
        <div class="mt-4 flex gap-2">
          {#each themeOptions as opt (opt.value)}
            <button
              onclick={() => setUITheme(opt.value)}
              class="flex-1 rounded-md border px-3 py-2 text-sm font-medium transition-colors {$uiTheme === opt.value
                ? 'border-primary bg-accent text-accent-foreground'
                : 'border-border text-muted-foreground hover:bg-accent/60'}"
            >
              {$t(opt.labelKey)}
            </button>
          {/each}
        </div>
      </section>

      <!-- Language -->
      <section class="rounded-xl border border-border bg-card p-5">
        <h2 class="font-display text-lg font-semibold">{$t("settings.languageSection")}</h2>
        <div class="mt-4 flex gap-2">
          <button
            onclick={() => setLocale("ro")}
            class="flex-1 rounded-md border px-3 py-2 text-sm font-medium transition-colors {$locale === 'ro'
              ? 'border-primary bg-accent text-accent-foreground'
              : 'border-border text-muted-foreground hover:bg-accent/60'}"
          >
            Română
          </button>
          <button
            onclick={() => setLocale("en")}
            class="flex-1 rounded-md border px-3 py-2 text-sm font-medium transition-colors {$locale === 'en'
              ? 'border-primary bg-accent text-accent-foreground'
              : 'border-border text-muted-foreground hover:bg-accent/60'}"
          >
            English
          </button>
        </div>
      </section>

      <!-- Backup -->
      <section class="rounded-xl border border-border bg-card p-5">
        <h2 class="font-display text-lg font-semibold">{$t("settings.backupSection")}</h2>
        <p class="mt-1 text-xs text-muted-foreground">{$t("settings.backupHint")}</p>
        <div class="mt-4 flex items-center gap-3">
          <button
            onclick={handleBackup}
            disabled={backingUp}
            class="rounded-md border border-border px-4 py-2 text-sm font-medium transition-colors hover:bg-accent disabled:opacity-60"
          >
            {backingUp ? $t("settings.backingUp") : $t("settings.backupNow")}
          </button>
        </div>
        {#if backupError}
          <p class="mt-3 text-sm text-destructive">{backupError}</p>
        {:else if backupPath}
          <p class="mt-3 text-sm text-primary">
            {$t("settings.backupSaved")} <span class="font-medium">{backupPath}</span>
            <button onclick={handleOpenBackupFolder} class="ml-2 font-medium text-primary hover:underline">
              {$t("settings.openFolder")}
            </button>
          </p>
        {/if}
      </section>

      {#if saveError}
        <p class="text-sm text-destructive">{saveError}</p>
      {:else if saved}
        <p class="text-sm text-primary">{$t("settings.saved")}</p>
      {/if}

      <div class="flex justify-end">
        <button
          onclick={save}
          disabled={saving}
          class="rounded-md bg-primary px-5 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-60"
        >
          {saving ? $t("common.saving") : $t("settings.save")}
        </button>
      </div>
    </div>
  {/if}
</div>
