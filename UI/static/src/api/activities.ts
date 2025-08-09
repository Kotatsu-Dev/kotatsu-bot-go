import type { AxiosInstance } from "axios";
import { format } from "date-fns";
import z from "zod";

const Activity = z.object({
  id: z.int(),
  created_at: z.iso.datetime({ offset: true }),
  title: z.string(),
  participants: z.any().array(), // TODO
  date_meeting: z.iso.datetime({ offset: true }),
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

    async create(props: {
      title: string;
      date_meeting: Date;
      description: string;
      location: string;
      send_images: FileList | File[];
    }) {
      await $.postForm("/activities/", {
        ...props,
        date_meeting: format(props.date_meeting, "yyyy-MM-dd HH:mm"),
      });
    },

    async wipe() {
      await $.delete("/activities/");
    },
  };
};
