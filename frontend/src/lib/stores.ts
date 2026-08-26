import { writable } from "svelte/store";
import type { View } from "./types";

// Aplicația Wails este single-window, deci nu avem nevoie de un router real
// (gen tanstack-router din UI-ul React). Un store simplu e suficient pentru
// a comuta între "pagini".
export const currentView = writable<View>("clients");

// ID-ul planului deschis curent în editor (setat din ClientsView când se
// apasă "Deschide" pe un plan).
export const currentPlanId = writable<string | null>(null);
