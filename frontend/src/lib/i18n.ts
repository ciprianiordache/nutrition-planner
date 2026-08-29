import { derived, writable } from "svelte/store";
import { settingsApi } from "./api";

export type Locale = "ro" | "en";

const STORAGE_KEY = "nutrition-planner:locale";

// translations — every user-facing string in the app, in one place. Keys
// are grouped by screen/component ("client.*", "plan.*", "foodForm.*"...)
// so it's obvious where a string is used without needing to search the
// codebase. Genuinely shared, single-word labels ("Cantitate"/"Quantity",
// "Ordine"/"Order") live under "common.*" instead of being repeated with
// the same value under five different feature keys.
const translations = {
  ro: {
    "app.name": "NutriPlan",
    "app.tagline": "Planuri nutriționale",

    "nav.clients": "Clienți",
    "nav.plan": "Editor plan",
    "nav.settings": "Setări",
    "nav.trainerRole": "Antrenor Personal",

    "titlebar.appName": "NutriPlan — planuri nutriționale",
    "titlebar.minimize": "Minimizează",
    "titlebar.restore": "Restaurează",
    "titlebar.maximize": "Maximizează",
    "titlebar.close": "Închide",

    "common.cancel": "Anulează",
    "common.delete": "Șterge",
    "common.deleting": "Se șterge…",
    "common.saving": "Se salvează…",
    "common.saveChanges": "Salvează modificările",
    "common.close": "Închide",
    "common.confirmDeleteTitle": "Confirmă ștergerea",
    "common.order": "Ordine",
    "common.quantity": "Cantitate",
    "common.unit": "Unitate",
    "common.edit": "Editează",
    "common.na": "—",

    "target.weight_loss": "Slăbire",
    "target.muscle_gain": "Masă musculară",
    "target.maintenance": "Menținere",
    "target.recomposition": "Recompoziție",

    "client.title": "Clienți",
    "client.subtitle": "{count} clienți · alege un client pentru a-i vedea planurile",
    "client.new": "+ Client nou",
    "client.searchPlaceholder": "Caută după nume…",
    "client.allTargets": "Toate obiectivele",
    "client.loading": "Se încarcă clienții…",
    "client.empty": "Niciun client încă. Adaugă primul client ca să începi.",
    "client.noMatch": "Niciun client nu corespunde căutării sau filtrului selectat.",
    "client.editAria": "Editează client",
    "client.deleteAria": "Șterge client",
    "client.plansOf": "Planurile lui {name}",
    "client.newPlan": "+ Plan nou",
    "client.loadingPlans": "Se încarcă planurile…",
    "client.noPlansYet": "Niciun plan creat încă pentru acest client.",
    "client.open": "Deschide",
    "client.editPlanAria": "Editează plan",
    "client.deletePlanAria": "Șterge plan",
    "client.duplicatePlanAria": "Duplică plan",
    "client.copySuffix": "(Copie)",
    "client.duplicatePlanTitle": "Duplică plan",
    "client.duplicateTargetClient": "Pentru ce client?",
    "client.duplicatePlanConfirm": "Duplică",
    "client.newTitle": "Client nou",
    "client.editTitle": "Editează client",
    "client.deleteConfirm": "Sigur vrei să ștergi clientul „{name}”? Toate planurile lui vor fi șterse.",
    "client.newPlanTitle": "Plan nou",
    "client.editPlanTitle": "Editează plan",
    "client.deletePlanConfirm": "Sigur vrei să ștergi planul „{title}”?",

    "clientForm.name": "Nume",
    "clientForm.namePlaceholder": "Ex: Andreea Popescu",
    "clientForm.weight": "Greutate (kg)",
    "clientForm.height": "Înălțime (cm)",
    "clientForm.target": "Obiectiv",
    "clientForm.chooseTarget": "Alege obiectivul",
    "clientForm.saveClient": "Salvează client",

    "planForm.title": "Titlu plan",
    "planForm.titlePlaceholder": "Ex: Plan alimentar — August",
    "planForm.recommendations": "Recomandări generale",
    "planForm.recommendationsPlaceholder": "Ex: 3 litri de apă zilnic, evită prăjelile, minim 7 ore de somn.",
    "planForm.footer": "Footer plan",
    "planForm.save": "Salvează plan",

    "dayVariantForm.name": "Nume variantă de zi",
    "dayVariantForm.namePlaceholder": "Ex: Zi de antrenament",
    "dayVariantForm.add": "Adaugă variantă",

    "mealForm.name": "Nume masă",
    "mealForm.namePlaceholder": "Ex: Mic dejun",
    "mealForm.note": "Notă",
    "mealForm.notePlaceholder": "Ex: În maxim o oră după trezire",
    "mealForm.add": "Adaugă masă",

    "foodForm.food": "Aliment",
    "foodForm.foodPlaceholder": "Ex: Piept de pui",
    "foodForm.addAlternative": "Adaugă alternativă",
    "foodForm.alternative": "Alternativă",
    "foodForm.alternativePlaceholder": "Ex: Somon",
    "foodForm.add": "Adaugă aliment",

    "supplementForm.name": "Supliment",
    "supplementForm.namePlaceholder": "Ex: Multivitamine",
    "supplementForm.add": "Adaugă supliment",

    "plan.noneSelected": "Niciun plan selectat.",
    "plan.backToClients": "← Înapoi la clienți",
    "plan.loading": "Se încarcă planul…",
    "plan.preview": "Previzualizare",
    "plan.previewing": "Se pregătește…",
    "plan.exportPdf": "Export PDF",
    "plan.exporting": "Se exportă…",
    "plan.previewOpened": "Previzualizare deschisă — dacă nu s-a deschis automat vizualizatorul de PDF, fișierul e la:",
    "plan.pdfSaved": "PDF salvat:",
    "plan.openFolder": "Deschide folderul",
    "plan.addDayVariant": "+ Variantă de zi",
    "plan.editDay": "Editează ziua",
    "plan.deleteDay": "Șterge ziua",
    "plan.addFood": "+ Aliment",
    "plan.addSupplement": "+ Supliment",
    "plan.tableFood": "Aliment",
    "plan.tableQuantity": "Cantitate",
    "plan.tableAlternative": "Alternativă",
    "plan.noFoods": "Niciun aliment adăugat.",
    "plan.supplementsLabel": "Suplimente",
    "plan.editFoodAria": "Editează aliment",
    "plan.deleteFoodAria": "Șterge aliment",
    "plan.editSupplementAria": "Editează supliment",
    "plan.deleteSupplementAria": "Șterge supliment",
    "plan.noMealsFor": "Nicio masă adăugată pentru „{name}”.",
    "plan.addMeal": "+ Adaugă masă",
    "plan.noDayVariants": "Nicio variantă de zi încă. Adaugă prima variantă ca să începi.",
    "plan.footerTitle": "Footer plan",
    "plan.newDayVariant": "Variantă de zi nouă",
    "plan.editDayVariant": "Editează variantă de zi",
    "plan.deleteDayVariantConfirm": "Sigur vrei să ștergi „{name}”? Toate mesele din această zi vor fi șterse.",
    "plan.newMeal": "Masă nouă",
    "plan.editMeal": "Editează masă",
    "plan.deleteMealConfirm": "Sigur vrei să ștergi masa „{name}”?",
    "plan.newFood": "Aliment nou",
    "plan.editFood": "Editează aliment",
    "plan.deleteFoodConfirm": "Sigur vrei să ștergi „{name}”?",
    "plan.newSupplement": "Supliment nou",
    "plan.editSupplement": "Editează supliment",
    "plan.deleteSupplementConfirm": "Sigur vrei să ștergi „{name}”?",

    "settings.title": "Setări",
    "settings.subtitle": "Personalizează exportul, tema PDF-ului și aspectul aplicației.",
    "settings.loading": "Se încarcă setările…",
    "settings.saved": "Setări salvate.",
    "settings.save": "Salvează setările",

    "settings.exportSection": "Export PDF",
    "settings.exportFolder": "Folder de export",
    "settings.exportFolderHint": "Fiecare client primește un subfolder propriu în interiorul acestui folder. Lasă gol pentru locația implicită (Desktop/NutriPlan Exports).",
    "settings.chooseFolder": "Alege folder…",
    "settings.exportFolderDefault": "(implicit)",

    "settings.pdfThemeSection": "Tema PDF-ului",
    "settings.pdfThemeHint": "Aceste culori se aplică fiecărui PDF exportat sau previzualizat.",
    "settings.colorPrimary": "Accent principal",
    "settings.colorPrimaryLight": "Accent — fundal deschis",
    "settings.colorOnPrimary": "Text pe accent",
    "settings.colorPageBackground": "Fundal pagină",
    "settings.colorCardBackground": "Fundal card",
    "settings.colorCardBorder": "Bordură card",
    "settings.colorText": "Text principal",
    "settings.colorMuted": "Text secundar",
    "settings.colorTableHeader": "Antet tabel",
    "settings.colorDivider": "Linie separatoare",
    "settings.resetColors": "Resetează la valorile implicite",

    "settings.appearanceSection": "Aspect aplicație",
    "settings.themeLight": "Deschisă",
    "settings.themeDark": "Întunecată",
    "settings.themeSystem": "Automat (sistem)",

    "settings.languageSection": "Limbă",

    "settings.backupSection": "Backup",
    "settings.backupHint": "Salvează o copie a bazei de date locale (clienți, planuri, tot) într-o arhivă .zip, ca plasă de siguranță.",
    "settings.backupNow": "Backup acum",
    "settings.backingUp": "Se salvează…",
    "settings.backupSaved": "Backup salvat:",
    "settings.openFolder": "Deschide folderul",

    "settings.profileSection": "Profilul tău",
    "settings.trainerName": "Nume",
    "settings.trainerNamePlaceholder": "Ex: Andreea Popescu",
    "settings.trainerRole": "Profesie",
    "settings.trainerRolePlaceholder": "Ex: Nutriționist",

    "onboarding.title": "Bine ai venit!",
    "onboarding.intro": "Înainte să începi, spune-ne cine ești — apare în bara laterală a aplicației.",
    "onboarding.language": "Limbă",
    "onboarding.save": "Continuă",
  },
  en: {
    "app.name": "NutriPlan",
    "app.tagline": "Nutrition plans",

    "nav.clients": "Clients",
    "nav.plan": "Plan editor",
    "nav.settings": "Settings",
    "nav.trainerRole": "Personal Trainer",

    "titlebar.appName": "NutriPlan — nutrition plans",
    "titlebar.minimize": "Minimize",
    "titlebar.restore": "Restore",
    "titlebar.maximize": "Maximize",
    "titlebar.close": "Close",

    "common.cancel": "Cancel",
    "common.delete": "Delete",
    "common.deleting": "Deleting…",
    "common.saving": "Saving…",
    "common.saveChanges": "Save changes",
    "common.close": "Close",
    "common.confirmDeleteTitle": "Confirm deletion",
    "common.order": "Order",
    "common.quantity": "Quantity",
    "common.unit": "Unit",
    "common.edit": "Edit",
    "common.na": "—",

    "target.weight_loss": "Weight loss",
    "target.muscle_gain": "Muscle gain",
    "target.maintenance": "Maintenance",
    "target.recomposition": "Recomposition",

    "client.title": "Clients",
    "client.subtitle": "{count} clients · choose a client to see their plans",
    "client.new": "+ New client",
    "client.searchPlaceholder": "Search by name…",
    "client.allTargets": "All objectives",
    "client.loading": "Loading clients…",
    "client.empty": "No clients yet. Add your first client to get started.",
    "client.noMatch": "No client matches the search or selected filter.",
    "client.editAria": "Edit client",
    "client.deleteAria": "Delete client",
    "client.plansOf": "{name}'s plans",
    "client.newPlan": "+ New plan",
    "client.loadingPlans": "Loading plans…",
    "client.noPlansYet": "No plans created yet for this client.",
    "client.open": "Open",
    "client.editPlanAria": "Edit plan",
    "client.deletePlanAria": "Delete plan",
    "client.duplicatePlanAria": "Duplicate plan",
    "client.copySuffix": "(Copy)",
    "client.duplicatePlanTitle": "Duplicate plan",
    "client.duplicateTargetClient": "For which client?",
    "client.duplicatePlanConfirm": "Duplicate",
    "client.newTitle": "New client",
    "client.editTitle": "Edit client",
    "client.deleteConfirm": "Are you sure you want to delete client \u201c{name}\u201d? All their plans will be deleted.",
    "client.newPlanTitle": "New plan",
    "client.editPlanTitle": "Edit plan",
    "client.deletePlanConfirm": "Are you sure you want to delete plan \u201c{title}\u201d?",

    "clientForm.name": "Name",
    "clientForm.namePlaceholder": "E.g. Andreea Popescu",
    "clientForm.weight": "Weight (kg)",
    "clientForm.height": "Height (cm)",
    "clientForm.target": "Objective",
    "clientForm.chooseTarget": "Choose objective",
    "clientForm.saveClient": "Save client",

    "planForm.title": "Plan title",
    "planForm.titlePlaceholder": "E.g. Meal plan — August",
    "planForm.recommendations": "General recommendations",
    "planForm.recommendationsPlaceholder": "E.g. 3 liters of water daily, avoid fried food, at least 7 hours of sleep.",
    "planForm.footer": "Plan footer",
    "planForm.save": "Save plan",

    "dayVariantForm.name": "Day variant name",
    "dayVariantForm.namePlaceholder": "E.g. Training day",
    "dayVariantForm.add": "Add variant",

    "mealForm.name": "Meal name",
    "mealForm.namePlaceholder": "E.g. Breakfast",
    "mealForm.note": "Note",
    "mealForm.notePlaceholder": "E.g. Within an hour of waking up",
    "mealForm.add": "Add meal",

    "foodForm.food": "Food",
    "foodForm.foodPlaceholder": "E.g. Chicken breast",
    "foodForm.addAlternative": "Add alternative",
    "foodForm.alternative": "Alternative",
    "foodForm.alternativePlaceholder": "E.g. Salmon",
    "foodForm.add": "Add food",

    "supplementForm.name": "Supplement",
    "supplementForm.namePlaceholder": "E.g. Multivitamins",
    "supplementForm.add": "Add supplement",

    "plan.noneSelected": "No plan selected.",
    "plan.backToClients": "← Back to clients",
    "plan.loading": "Loading plan…",
    "plan.preview": "Preview",
    "plan.previewing": "Preparing…",
    "plan.exportPdf": "Export PDF",
    "plan.exporting": "Exporting…",
    "plan.previewOpened": "Preview opened — if the PDF viewer didn't open automatically, the file is at:",
    "plan.pdfSaved": "PDF saved:",
    "plan.openFolder": "Open folder",
    "plan.addDayVariant": "+ Day variant",
    "plan.editDay": "Edit day",
    "plan.deleteDay": "Delete day",
    "plan.addFood": "+ Food",
    "plan.addSupplement": "+ Supplement",
    "plan.tableFood": "Food",
    "plan.tableQuantity": "Quantity",
    "plan.tableAlternative": "Alternative",
    "plan.noFoods": "No food added.",
    "plan.supplementsLabel": "Supplements",
    "plan.editFoodAria": "Edit food",
    "plan.deleteFoodAria": "Delete food",
    "plan.editSupplementAria": "Edit supplement",
    "plan.deleteSupplementAria": "Delete supplement",
    "plan.noMealsFor": "No meal added for \u201c{name}”.",
    "plan.addMeal": "+ Add meal",
    "plan.noDayVariants": "No day variant yet. Add the first variant to get started.",
    "plan.footerTitle": "Plan footer",
    "plan.newDayVariant": "New day variant",
    "plan.editDayVariant": "Edit day variant",
    "plan.deleteDayVariantConfirm": "Are you sure you want to delete \u201c{name}”? All meals in this day will be deleted.",
    "plan.newMeal": "New meal",
    "plan.editMeal": "Edit meal",
    "plan.deleteMealConfirm": "Are you sure you want to delete meal \u201c{name}”?",
    "plan.newFood": "New food",
    "plan.editFood": "Edit food",
    "plan.deleteFoodConfirm": "Are you sure you want to delete \u201c{name}”?",
    "plan.newSupplement": "New supplement",
    "plan.editSupplement": "Edit supplement",
    "plan.deleteSupplementConfirm": "Are you sure you want to delete \u201c{name}”?",

    "settings.title": "Settings",
    "settings.subtitle": "Customize the export location, the PDF theme, and the app's appearance.",
    "settings.loading": "Loading settings…",
    "settings.saved": "Settings saved.",
    "settings.save": "Save settings",

    "settings.exportSection": "PDF export",
    "settings.exportFolder": "Export folder",
    "settings.exportFolderHint": "Each client gets their own subfolder inside this folder. Leave empty for the default location (Desktop/NutriPlan Exports).",
    "settings.chooseFolder": "Choose folder…",
    "settings.exportFolderDefault": "(default)",

    "settings.pdfThemeSection": "PDF theme",
    "settings.pdfThemeHint": "These colors apply to every exported or previewed PDF.",
    "settings.colorPrimary": "Primary accent",
    "settings.colorPrimaryLight": "Accent — light background",
    "settings.colorOnPrimary": "Text on accent",
    "settings.colorPageBackground": "Page background",
    "settings.colorCardBackground": "Card background",
    "settings.colorCardBorder": "Card border",
    "settings.colorText": "Main text",
    "settings.colorMuted": "Muted text",
    "settings.colorTableHeader": "Table header",
    "settings.colorDivider": "Divider line",
    "settings.resetColors": "Reset to defaults",

    "settings.appearanceSection": "App appearance",
    "settings.themeLight": "Light",
    "settings.themeDark": "Dark",
    "settings.themeSystem": "Automatic (system)",

    "settings.languageSection": "Language",

    "settings.backupSection": "Backup",
    "settings.backupHint": "Save a copy of the local database (clients, plans, everything) into a .zip archive, as a safety net.",
    "settings.backupNow": "Backup now",
    "settings.backingUp": "Saving…",
    "settings.backupSaved": "Backup saved:",
    "settings.openFolder": "Open folder",

    "settings.profileSection": "Your profile",
    "settings.trainerName": "Name",
    "settings.trainerNamePlaceholder": "E.g. Andreea Popescu",
    "settings.trainerRole": "Profession",
    "settings.trainerRolePlaceholder": "E.g. Nutritionist",

    "onboarding.title": "Welcome!",
    "onboarding.intro": "Before you start, tell us who you are — this shows up in the app's sidebar.",
    "onboarding.language": "Language",
    "onboarding.save": "Continue",
  },
} as const;

export type TranslationKey = keyof (typeof translations)["ro"];

// TARGET_KEYS — the fixed set of client "objective" values, stored in the
// database as these exact language-neutral strings (not the translated
// label) so the same client record displays correctly regardless of which
// UI language is active. Existing rows seeded before this field became
// key-based (plain Romanian text like "Slăbire") won't match any of these
// keys - they're just displayed as-is (see ClientsView/ClientForm) rather
// than crashing or disappearing.
export const TARGET_KEYS = ["weight_loss", "muscle_gain", "maintenance", "recomposition"] as const;
export type TargetKey = (typeof TARGET_KEYS)[number];

function detectInitialLocale(): Locale {
  if (typeof localStorage !== "undefined") {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored === "ro" || stored === "en") return stored;
  }
  return "ro";
}

export const locale = writable<Locale>(detectInitialLocale());

if (typeof localStorage !== "undefined") {
  locale.subscribe((l) => localStorage.setItem(STORAGE_KEY, l));
}

/**
 * setLocale updates the locale immediately (so the whole UI re-renders in
 * the new language right away) and persists the choice to the backend
 * settings in the background - the same setting also picks the language
 * for generated PDFs, so switching here is a one-stop "demo this app in
 * English" toggle. A failed backend save doesn't roll back the UI change:
 * the user's immediate choice always wins over a slow/failed network call.
 */
export function setLocale(l: Locale): void {
  locale.set(l);
  settingsApi.update({ language: l }).catch(() => {
    // Best-effort - the UI already reflects the new language regardless,
    // and the next successful settings load/save will reconcile it.
  });
}

function interpolate(template: string, vars?: Record<string, string | number>): string {
  if (!vars) return template;
  let result = template;
  for (const [key, value] of Object.entries(vars)) {
    result = result.replaceAll(`{${key}}`, String(value));
  }
  return result;
}

/**
 * t — a derived store holding the current translate function. Use as
 * `$t("client.title")` in a component; it re-evaluates automatically
 * whenever `locale` changes, so switching language updates every mounted
 * component with no extra wiring.
 */
export const t = derived(locale, ($locale) => {
  return (key: TranslationKey, vars?: Record<string, string | number>): string => {
    const template = translations[$locale][key] ?? translations.ro[key] ?? key;
    return interpolate(template, vars);
  };
});
