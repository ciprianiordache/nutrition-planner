import { writable } from "svelte/store";
import { settingsApi } from "./api";

export type UITheme = "light" | "dark" | "system";

const STORAGE_KEY = "nutrition-planner:ui-theme";

function systemPrefersDark(): boolean {
  return typeof window !== "undefined" && window.matchMedia?.("(prefers-color-scheme: dark)").matches;
}

function detectInitialTheme(): UITheme {
  if (typeof localStorage !== "undefined") {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored === "light" || stored === "dark" || stored === "system") return stored;
  }
  return "system";
}

/** applyTheme toggles the "dark" class on <html> - every color in
 * style.css is driven by that single class (see the ".dark { ... }" block),
 * so this is the only DOM touch point the whole theme feature needs. */
function applyTheme(theme: UITheme): void {
  if (typeof document === "undefined") return;
  const isDark = theme === "dark" || (theme === "system" && systemPrefersDark());
  document.documentElement.classList.toggle("dark", isDark);
}

export const uiTheme = writable<UITheme>(detectInitialTheme());

uiTheme.subscribe((theme) => {
  applyTheme(theme);
  if (typeof localStorage !== "undefined") localStorage.setItem(STORAGE_KEY, theme);
});

// Keep "system" theme in sync if the OS preference changes while the app
// is open (e.g. the user's OS switches to dark mode at sunset).
if (typeof window !== "undefined" && window.matchMedia) {
  window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", () => {
    uiTheme.update((current) => current); // re-trigger the subscriber above with the same value
  });
}

/**
 * setUITheme updates the theme immediately and persists the choice to the
 * backend settings in the background, mirroring setLocale in i18n.ts. A
 * failed backend save doesn't roll back the UI change.
 */
export function setUITheme(theme: UITheme): void {
  uiTheme.set(theme);
  settingsApi.update({ uiTheme: theme }).catch(() => {
    // Best-effort - see i18n.ts's setLocale for the same reasoning.
  });
}
