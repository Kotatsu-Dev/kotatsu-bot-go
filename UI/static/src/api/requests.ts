import type { AxiosInstance } from "axios";

export const createRequestsApi = ($: AxiosInstance) => {
  return {
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
