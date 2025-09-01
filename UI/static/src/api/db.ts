import type { AxiosInstance } from "axios";

export const createDbApi = ($: AxiosInstance, root: string) => {
  return {
    async wipe() {
      await $.delete(`${root}/all-db`);
    },
  };
};
