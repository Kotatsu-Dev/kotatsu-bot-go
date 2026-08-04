import { createListCollection } from "@chakra-ui/react";

export const genders = [
  { value: "male", label: "Male" },
  { value: "female", label: "Female" },
  { value: "", label: "Unknown" },
];

export const itmoStatuses = [
  { value: "guest", label: "Guest" },
  { value: "student", label: "Student" },
  { value: "graduate", label: "Graduate" },
  { value: "employee", label: "Employee" },
  { value: "student_employee", label: "Student and employee" },
  { value: "graduate_employee", label: "Graduate and employee" },
  { value: "", label: "Unknown" },
];

export const itmoStatusCollection = createListCollection({
  items: itmoStatuses,
});

export type ItmoTrait =
  | "guest"
  | "student"
  | "graduate"
  | "employee"
  | "unknown";

// Каждый itmo_status раскладывается на независимые признаки: сам статус
// скалярный, но student_employee и graduate_employee кодируют сразу два факта.
// Фильтр требует наличия всех выбранных признаков, поэтому Student + Employee
// даёт ровно student_employee.
export const ITMO_TRAITS: Record<string, ItmoTrait[]> = {
  guest: ["guest"],
  student: ["student"],
  graduate: ["graduate"],
  employee: ["employee"],
  student_employee: ["student", "employee"],
  graduate_employee: ["graduate", "employee"],
  "": ["unknown"],
};

export const itmoTraits: { value: ItmoTrait; label: string }[] = [
  { value: "guest", label: "Guest" },
  { value: "student", label: "Student" },
  { value: "graduate", label: "Graduate" },
  { value: "employee", label: "Employee" },
  { value: "unknown", label: "Unknown" },
];

export const itmoTraitCollection = createListCollection({ items: itmoTraits });

export const traitsOf = (itmoStatus: string | null | undefined): ItmoTrait[] =>
  ITMO_TRAITS[itmoStatus ?? ""] ?? ["unknown"];
