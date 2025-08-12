import type { AxiosInstance } from "axios";
import z from "zod";
import { User } from "./users";
import { format } from "date-fns";
import { handleZodError } from "./api";

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
      return handleZodError(() =>
        Roulette.array().parse(result.data.data.list_anime_roulettes)
      );
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
            end_date: format(props.start_date, `yyyy-MM-dd HH:mm`),
          },
          {
            stage: 1,
            end_date: format(props.announce_date, `yyyy-MM-dd HH:mm`),
          },
          {
            stage: 2,
            end_date: format(props.distribution_date, `yyyy-MM-dd HH:mm`),
          },
          {
            stage: 3,
            end_date: format(props.end_date, `yyyy-MM-dd HH:mm`),
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
        start_date: format(props.start_date, `yyyy-MM-dd HH:mm`),
        announce_date: format(props.announce_date, `yyyy-MM-dd HH:mm`),
        distribution_date: format(props.distribution_date, `yyyy-MM-dd HH:mm`),
        end_date: format(props.end_date, `yyyy-MM-dd HH:mm`),
        theme: props.theme,
      });
    },

    async wipe() {
      await $.delete("/roulettes/");
    },
  };
};
