import type { AxiosInstance } from "axios";
import z from "zod";

const Activity = z.object({
  id: z.int(),
  created_at: z.iso.datetime(),
  title: z.string(),
  participants: z.any().array(), // TODO
  date_meeting: z.iso.datetime(),
  description: z.string(),
  location: z.string(),
  path_images: z.string().array().optional(),
  status: z.boolean(),
});
export type Activity = z.infer<typeof Activity>;

export const createActivitiesApi = ($: AxiosInstance) => {
  return {
    async getAll() {
      const res = await $.get("/activities/");
      return Activity.array().parse(res.data.data.list_activities);
    },

    async setStatus({ id, status }: { id: number; status: boolean }) {
      await $.put("/activities/", {
        activity_id: id,
        status,
      });
    },
  };
};
