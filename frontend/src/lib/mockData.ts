import type { Client, DayVariant } from "./types";

// NOTE: date mock — de înlocuit cu apeluri către binding-urile Wails/Go
// (internals/app/client, .../plan, .../day_variant, .../meal, .../food, .../supplement)
// odată ce sunt expuse metodele corespunzătoare în app.go.

export const clients: Client[] = [
  { id: "1", name: "Andreea Popescu", weight: 62, height: 168, target: "Slăbire", plans: 3, updatedAt: "acum 2 zile" },
  { id: "2", name: "Mihai Ionescu", weight: 88, height: 182, target: "Masă musculară", plans: 2, updatedAt: "acum 5 zile" },
  { id: "3", name: "Elena Radu", weight: 71, height: 165, target: "Menținere", plans: 1, updatedAt: "acum 1 săptămână" },
  { id: "4", name: "Vlad Marinescu", weight: 95, height: 190, target: "Recompoziție", plans: 4, updatedAt: "acum 2 săptămâni" },
  { id: "5", name: "Ioana Dobre", weight: 55, height: 160, target: "Slăbire", plans: 1, updatedAt: "acum o lună" },
];

export const dayVariants: DayVariant[] = [
  {
    id: "d1",
    name: "Zi de antrenament",
    meals: [
      {
        id: "m1",
        name: "Mic dejun",
        note: "În maxim o oră după trezire",
        foods: [
          { id: "f1", baseFood: "Ovăz", baseQuantity: 60, baseUnit: "g", alternativeFood: "Orez brun", alternativeQuantity: 50, alternativeUnit: "g" },
          { id: "f2", baseFood: "Ouă", baseQuantity: 3, baseUnit: "buc", alternativeFood: "Brânză de vaci", alternativeQuantity: 150, alternativeUnit: "g" },
          { id: "f3", baseFood: "Afine", baseQuantity: 80, baseUnit: "g" },
        ],
        supplements: [{ id: "s1", supplement: "Multivitamine", quantity: 1 }],
      },
      {
        id: "m2",
        name: "Prânz",
        note: "La 3–4 ore după masa anterioară",
        foods: [
          { id: "f4", baseFood: "Piept de pui", baseQuantity: 180, baseUnit: "g", alternativeFood: "Somon", alternativeQuantity: 150, alternativeUnit: "g" },
          { id: "f5", baseFood: "Orez basmati", baseQuantity: 80, baseUnit: "g", alternativeFood: "Cartof dulce", alternativeQuantity: 250, alternativeUnit: "g" },
          { id: "f6", baseFood: "Fasole verde", baseQuantity: 200, baseUnit: "g" },
        ],
        supplements: [{ id: "s2", supplement: "Omega-3", quantity: 2 }],
      },
      {
        id: "m3",
        name: "Cină",
        note: "Cu 2 ore înainte de somn",
        foods: [
          { id: "f7", baseFood: "Curcan", baseQuantity: 150, baseUnit: "g", alternativeFood: "Ton", alternativeQuantity: 120, alternativeUnit: "g" },
          { id: "f8", baseFood: "Salată mixtă", baseQuantity: 200, baseUnit: "g" },
        ],
        supplements: [{ id: "s3", supplement: "Magneziu", quantity: 1 }],
      },
    ],
  },
  { id: "d2", name: "Zi de repaus", meals: [] },
  { id: "d3", name: "Weekend", meals: [] },
];
