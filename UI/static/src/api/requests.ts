import type { AxiosInstance } from "axios";
import z from "zod";

export const Request = z.object({
  id: z.int().nonnegative(),
  created_at: z.iso.datetime({ offset: true }),
  type: z.int(),
  status: z.int(),
  user_id: z.int().nonnegative(),
});
export type Request = z.infer<typeof Request>;

export const createRequestsApi = ($: AxiosInstance) => {
  return {
    async getAll() { 
        const res = await $.get('/requests/');
        return Request.array().parse(res.data.data.list_requests); 
    },
    async accept(props: { id: number }) {
      await $.put("/requests/choice", {
        request_id: props.id,
        status: 1,
      });
    },
    async reject(props: { id: number }) {
      await $.put("/requests/choice", {
        request_id: props.id,
        status: 2,
      });
    },
  };
};
