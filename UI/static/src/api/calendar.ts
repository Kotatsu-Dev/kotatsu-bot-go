import type { AxiosInstance } from "axios";

export const createCalendarApi = ($: AxiosInstance, root: string) => {
  return {
    imageUrl() {
      return `${root}/get-calendar-file`;
    },
    async upload(props: { image: File }) {
      await $.postForm(`${root}/upload-file-calendar-activities`, props);
    },
  };
};
