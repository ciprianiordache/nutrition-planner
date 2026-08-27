<script lang="ts">
  import { t, locale, type Locale } from "./i18n";
  import { saveProfile } from "./profile";
  import { errorMessage } from "./api";

  let name = $state("");
  let role = $state("");
  let saving = $state(false);
  let error = $state<string | null>(null);

  // Choosing a language here only updates the live preview (the rest of
  // this modal's own text re-renders immediately via $t) - it's persisted
  // to the backend together with name/role in a single write on submit,
  // not the instant it's clicked.
  function chooseLanguage(l: Locale) {
    locale.set(l);
  }

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    if (!name.trim() || !role.trim()) return;
    saving = true;
    error = null;
    try {
      await saveProfile(name.trim(), role.trim(), $locale);
      // On success, App.svelte's `needsOnboarding` derivation flips to
      // false automatically (it reads the same `profile` store this just
      // updated), which unmounts this modal - no local "close" state needed.
    } catch (err) {
      error = errorMessage(err);
    } finally {
      saving = false;
    }
  }
</script>

<!-- Deliberately not built on Modal.svelte: this step is mandatory (the
     app has no login, so this is the only place it learns who's using it)
     - no overlay-click-to-dismiss, no Escape-to-close, no visible X button
     that would otherwise do nothing. -->
<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
  <div class="w-full max-w-sm rounded-xl border border-border bg-card p-6 shadow-xl">
    <h2 class="font-display text-lg font-semibold">{$t("onboarding.title")}</h2>
    <p class="mt-1 text-sm text-muted-foreground">{$t("onboarding.intro")}</p>

    <form onsubmit={submit} class="mt-5 space-y-4">
      <div>
        <span class="mb-1 block text-xs font-medium text-muted-foreground">{$t("onboarding.language")}</span>
        <div class="flex gap-2">
          <button
            type="button"
            onclick={() => chooseLanguage("ro")}
            class="flex-1 rounded-md border px-3 py-2 text-sm font-medium transition-colors {$locale === 'ro'
              ? 'border-primary bg-accent text-accent-foreground'
              : 'border-border text-muted-foreground hover:bg-accent/60'}"
          >
            Română
          </button>
          <button
            type="button"
            onclick={() => chooseLanguage("en")}
            class="flex-1 rounded-md border px-3 py-2 text-sm font-medium transition-colors {$locale === 'en'
              ? 'border-primary bg-accent text-accent-foreground'
              : 'border-border text-muted-foreground hover:bg-accent/60'}"
          >
            English
          </button>
        </div>
      </div>

      <div>
        <label class="mb-1 block text-xs font-medium text-muted-foreground" for="onboarding-name">
          {$t("settings.trainerName")}
        </label>
        <input
          id="onboarding-name"
          bind:value={name}
          required
          placeholder={$t("settings.trainerNamePlaceholder")}
          class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
        />
      </div>
      <div>
        <label class="mb-1 block text-xs font-medium text-muted-foreground" for="onboarding-role">
          {$t("settings.trainerRole")}
        </label>
        <input
          id="onboarding-role"
          bind:value={role}
          required
          placeholder={$t("settings.trainerRolePlaceholder")}
          class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring/40"
        />
      </div>

      {#if error}
        <p class="text-sm text-destructive">{error}</p>
      {/if}

      <button
        type="submit"
        disabled={saving}
        class="w-full rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-60"
      >
        {saving ? $t("common.saving") : $t("onboarding.save")}
      </button>
    </form>
  </div>
</div>
