import type { AxiosInstance } from "axios";

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
      await $.post(`${root}/api/broadcast/`, props);
    },
  };
};
