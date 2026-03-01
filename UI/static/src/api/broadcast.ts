import type { AxiosInstance } from "axios";
import z from "zod";
import { User } from "./users";

export const BroadcastResult = z.object({
  user: User,
  success: z.boolean(),
  error_message: z.string(),
});
export type BroadcastResult = z.infer<typeof BroadcastResult>;

export interface SendBroadcastRequest {
  events: number[];
  users: number[];
  roulettes: number[];
  club_member_status: boolean | null;
  itmo_status: string[];
  message: string;
}

export const createBroadcastApi = ($: AxiosInstance, root: string) => {
  return {
    async send(props: { message: string; files?: FileList | File[] }) {
      await $.postForm(`${root}/send-message-user`, props);
    },
    async sendBroadcast(props: SendBroadcastRequest) {
      const res = await $.post(`${root}/api/broadcast/`, props);
      return BroadcastResult.array().parse(res.data.data.results);
    },
  };
};
