import Fuse from "fuse.js";
import type { User } from "../api/users";

/**
 * Приводит телефон к сравнимому виду: только цифры, ведущая 8 -> 7.
 * Благодаря этому +7 999 123-45-67, 89991234567 и 79991234567 совпадают.
 */
export const normalizePhone = (value: string): string =>
  value.replace(/\D/g, "").replace(/^8/, "7");

/** ИСУ всегда ровно шесть цифр — иначе по нему не ищем вообще. */
const ISU_LENGTH = 6;

const isDigitsOnly = (query: string) => /^\d+$/.test(query);

const looksLikePhone = (query: string) =>
  /^[\d\s()+.-]+$/.test(query) && normalizePhone(query).length >= 5;

const byUsername = (users: User[], query: string) => {
  const needle = query.slice(1).toLowerCase();
  if (needle.length === 0) return users;
  return users.filter((user) => user.user_name.toLowerCase().includes(needle));
};

const byPhone = (users: User[], query: string) => {
  const needle = normalizePhone(query);
  return users.filter((user) =>
    normalizePhone(user.phone_number).includes(needle),
  );
};

/**
 * Цифровой запрос неоднозначен: Telegram ID и телефон пересекаются по формату.
 * Поэтому проверяем все подходящие поля и объединяем результаты, а не гадаем.
 * ИСУ участвует только при длине ровно ISU_LENGTH.
 */
const byNumber = (users: User[], query: string) => {
  const phone = normalizePhone(query);
  return users.filter(
    (user) =>
      (query.length === ISU_LENGTH && user.isu === query) ||
      String(user.user_tg_id).startsWith(query) ||
      normalizePhone(user.phone_number).includes(phone),
  );
};

const byName = (users: User[], query: string) =>
  new Fuse(users, {
    keys: ["full_name", "full_tg_name", "user_name"],
    threshold: 0.3,
    ignoreLocation: true,
  })
    .search(query)
    .map((result) => result.item);

/**
 * Единая строка поиска: тип запроса определяется автоматически.
 * @username — по username, только цифры — по ИСУ/Telegram ID/телефону,
 * похожее на телефон — по нормализованному телефону, остальное — нечёткий
 * поиск по именам.
 */
export const searchUsers = (users: User[], rawQuery: string): User[] => {
  const query = rawQuery.trim();
  if (query.length === 0) return users;
  if (query.startsWith("@")) return byUsername(users, query);
  if (isDigitsOnly(query)) return byNumber(users, query);
  if (looksLikePhone(query)) return byPhone(users, query);
  return byName(users, query);
};
