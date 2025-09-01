import type { AxiosInstance } from "axios";
import z from "zod";
import { Request } from "./requests";

export const User = z.object({
  id: z.int(),
  created_at: z.iso.datetime({ offset: true }),
  step: z.int(),
  user_tg_id: z.int(), // Not particularly safe, but...
  last_message_id: z.int(),
  user_name: z.string(),
  full_tg_name: z.string(),
  isu: z.string(),
  full_name: z.string(),
  phone_number: z.string(),
  secret_code: z.string(),
  is_itmo: z.boolean(),
  is_club_member: z.boolean(),
  is_subscribe_newsletter: z.boolean(),
  is_sent_request: z.boolean(),
  is_filled_data: z.boolean(),
  temp_activity_id: z.int(),
  my_activities: z.any().array().or(z.null()), // TODO
  link_my_anime_list: z.string(),
  my_request: Request.or(z.null()),
  enigmatic_title: z.string(),
});
export type User = z.infer<typeof User>;

export const createUsersApi = ($: AxiosInstance) => {
  return {
    async getAll() {
      const res = await $.get("/users/");
      return User.array().parse(res.data.data.list_users);
    },
    async setMemberStatus(props: {
      user_tg_id: number;
      is_club_member: boolean;
    }) {
      await $.put("/users/club-member", {
        user_tg_id: props.user_tg_id.toString(),
        is_club_member: Number(props.is_club_member),
      });
    },
    async wipe() {
      await $.delete("/users/");
    },
  };
};
