import { writable } from "svelte/store";
import { settingsApi, errorMessage } from "./api";
import { locale, type Locale } from "./i18n";

export type Profile = {
  name: string;
  role: string;
};

export const profile = writable<Profile>({ name: "", role: "" });

/** loaded becomes true once the initial settingsApi.get() call resolves -
 * used to tell "we don't know yet" apart from "we checked and there's no
 * name", since only the latter should trigger onboarding. */
export const profileLoaded = writable(false);
export const profileLoadError = writable<string | null>(null);

/**
 * loadProfile fetches the current trainer name/role from the backend and
 * populates the store. Call once, near app startup (see App.svelte) -
 * AppShell's sidebar and the onboarding check both read from this same
 * store rather than each fetching independently.
 */
export async function loadProfile(): Promise<void> {
  try {
    const s = await settingsApi.get();
    profile.set({ name: s.trainerName, role: s.trainerRole });
    profileLoadError.set(null);
  } catch (err) {
    profileLoadError.set(errorMessage(err));
  } finally {
    profileLoaded.set(true);
  }
}

/**
 * saveProfile persists the trainer's name/role (and optionally the
 * language, used by the first-run onboarding modal so the whole "who are
 * you, and in what language" step is a single backend write) and updates
 * the relevant stores immediately on success - used by both the
 * onboarding modal and the Settings screen, so they always stay
 * consistent with each other.
 */
export async function saveProfile(name: string, role: string, language?: Locale): Promise<void> {
  await settingsApi.update({
    trainerName: name,
    trainerRole: role,
    ...(language ? { language } : {}),
  });
  profile.set({ name, role });
  if (language) locale.set(language);
}
