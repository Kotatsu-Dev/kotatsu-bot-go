import type { AxiosInstance } from "axios";

export const createCalendarApi = ($: AxiosInstance, root: string) => {
  return {
    imageUrl() {
      return `${root}/api/calendar/`;
    },
    async upload(props: { image: File }) {
      await $.postForm(`/calendar/`, props);
    },
  };
};
