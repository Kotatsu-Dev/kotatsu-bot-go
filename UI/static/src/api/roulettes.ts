import type { AxiosInstance } from "axios";
import z from "zod";
import { User } from "./users";

const Roulette = z.object({
  id: z.int().nonnegative(),
  created_at: z.iso.datetime({ offset: true }).pipe(z.coerce.date()),
  start_date: z.iso.datetime({ offset: true }).pipe(z.coerce.date()),
  announce_date: z.iso.datetime({ offset: true }).pipe(z.coerce.date()),
  distribution_date: z.iso.datetime({ offset: true }).pipe(z.coerce.date()),
  end_date: z.iso.datetime({ offset: true }).pipe(z.coerce.date()),
  theme: z.string(),
  participants: User.array(),
  distribution: z.int().array().or(z.null()),
});
export type Roulette = z.infer<typeof Roulette>;

export const createRoulettesApi = ($: AxiosInstance) => {
  return {
    async getAll() {
      const result = await $.get("/roulettes/");
      return Roulette.array().parse(result.data.data.list_anime_roulettes);
    },
  };
};
