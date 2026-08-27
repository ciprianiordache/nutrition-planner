import * as ClientAPI from "../../wailsjs/go/client/Handler";
import * as PlanAPI from "../../wailsjs/go/plan/Handler";
import * as DayVariantAPI from "../../wailsjs/go/day_variant/Handler";
import * as MealAPI from "../../wailsjs/go/meal/Handler";
import * as FoodAPI from "../../wailsjs/go/food/Handler";
import * as SupplementAPI from "../../wailsjs/go/supplement/Handler";
import * as PdfExportAPI from "../../wailsjs/go/pdf_export/Handler";
import * as SettingsAPI from "../../wailsjs/go/settings/Handler";
import type { client, plan, day_variant, meal, food, supplement, settings } from "../../wailsjs/go/models";

// Strat subțire peste bindingurile generate de Wails, ca să nu importăm
// direct din wailsjs în componente — dacă se schimbă vreodată semnătura
// unei metode Go, se ajustează într-un singur loc.

export const clientApi = {
  list: (offset = 0, limit = 100) => ClientAPI.ListClients(offset, limit),
  create: (req: client.CreateRequest) => ClientAPI.CreateClient(req),
  get: (id: string) => ClientAPI.GetClient(id),
  update: (req: client.UpdateRequest) => ClientAPI.UpdateClient(req),
  remove: (id: string) => ClientAPI.Delete(id),
};

export const planApi = {
  listByClient: (clientId: string) => PlanAPI.ListPlansByClient(clientId),
  create: (req: plan.CreateRequest) => PlanAPI.CreatePlan(req),
  getFull: (id: string) => PlanAPI.GetPlanFull(id),
  update: (req: plan.UpdateRequest) => PlanAPI.UpdatePlan(req),
  remove: (id: string) => PlanAPI.DeletePlan(id),
  duplicate: (req: plan.DuplicateRequest) => PlanAPI.DuplicatePlan(req),
};

export const dayVariantApi = {
  listByPlan: (planId: string) => DayVariantAPI.ListDayVariantsByPlan(planId),
  create: (req: day_variant.CreateRequest) => DayVariantAPI.CreateDayVariant(req),
  update: (req: day_variant.UpdateRequest) => DayVariantAPI.UpdateDayVariant(req),
  remove: (id: string) => DayVariantAPI.DeleteDayVariant(id),
};

export const mealApi = {
  listByDayVariant: (dayVariantId: string) => MealAPI.ListMealsByDayVariant(dayVariantId),
  create: (req: meal.CreateRequest) => MealAPI.CreateMeal(req),
  update: (req: meal.UpdateRequest) => MealAPI.UpdateMeal(req),
  remove: (id: string) => MealAPI.DeleteMeal(id),
};

export const foodApi = {
  listByMeal: (mealId: string) => FoodAPI.ListFoodsByMeal(mealId),
  create: (req: food.CreateRequest) => FoodAPI.CreateFood(req),
  update: (req: food.UpdateRequest) => FoodAPI.UpdateFood(req),
  remove: (id: string) => FoodAPI.DeleteFood(id),
};

export const supplementApi = {
  listByMeal: (mealId: string) => SupplementAPI.ListSupplementsByMeal(mealId),
  create: (req: supplement.CreateRequest) => SupplementAPI.CreateSupplement(req),
  update: (req: supplement.UpdateRequest) => SupplementAPI.UpdateSupplement(req),
  remove: (id: string) => SupplementAPI.DeleteSupplement(id),
};

export const exportApi = {
  exportPlanPDF: (planId: string) => PdfExportAPI.ExportPlanPDF(planId),
  previewPlanPDF: (planId: string) => PdfExportAPI.PreviewPlanPDF(planId),
  openExportFolder: (filePath: string) => PdfExportAPI.OpenExportFolder(filePath),
};

export const settingsApi = {
  get: () => SettingsAPI.GetSettings(),
  update: (req: settings.UpdateRequest) => SettingsAPI.UpdateSettings(req),
  chooseExportFolder: () => SettingsAPI.ChooseExportFolder(),
  backupDatabase: () => SettingsAPI.BackupDatabase(),
  openBackupFolder: (filePath: string) => SettingsAPI.OpenBackupFolder(filePath),
};

/** Extrage un mesaj de eroare lizibil dintr-un throw venit din Go (de obicei string). */
export function errorMessage(err: unknown): string {
  if (err instanceof Error) return err.message;
  if (typeof err === "string") return err;
  return "A apărut o eroare neașteptată.";
}
