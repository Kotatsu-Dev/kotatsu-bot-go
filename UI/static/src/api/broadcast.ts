import type { AxiosInstance } from "axios";

export const createBroadcastApi = ($: AxiosInstance, root: string) => {
  return {
    async send(props: { message: string; files?: FileList | File[] }) {
      await $.postForm(`${root}/send-message-user`, props);
    },
  };
};
