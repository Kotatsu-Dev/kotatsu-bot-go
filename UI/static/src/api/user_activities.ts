import z from "zod";

export const UserActivity = z.object({
  user_id: z.int(),
  activity_id: z.int(),
  visited: z.boolean().or(z.null()),
});
export type UserActivity = z.infer<typeof UserActivity>;
