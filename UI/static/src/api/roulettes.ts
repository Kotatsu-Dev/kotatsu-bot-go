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
  participants: User.array().or(z.null()),
  distribution: z.int().array().or(z.null()),
});
export type Roulette = z.infer<typeof Roulette>;

export const createRoulettesApi = ($: AxiosInstance) => {
  return {
    async getAll() {
      const result = await $.get("/roulettes/");
      return Roulette.array().parse(result.data.data.list_anime_roulettes);
    },

    async create(props: {
      start_date: Date;
      announce_date: Date;
      distribution_date: Date;
      end_date: Date;
      theme: string;
    }) {
      await $.post("/roulettes/", {
        stages: [
          {
            stage: 0,
            end_date: props.start_date.toISOString(),
          },
          {
            stage: 1,
            end_date: props.announce_date.toISOString(),
          },
          {
            stage: 2,
            end_date: props.distribution_date.toISOString(),
          },
          {
            stage: 3,
            end_date: props.end_date.toISOString(),
          },
        ],
        theme: props.theme,
      });
    },

    async update(props: {
      id: number;
      start_date: Date;
      announce_date: Date;
      distribution_date: Date;
      end_date: Date;
      theme: string;
    }) {
      await $.put("/roulettes/", {
        id: props.id,
        start_date: props.start_date.toISOString(),
        announce_date: props.announce_date.toISOString(),
        distribution_date: props.distribution_date.toISOString(),
        end_date: props.end_date.toISOString(),
        theme: props.theme,
      });
    },

    async wipe() {
      await $.delete("/roulettes/");
    },
  };
};
