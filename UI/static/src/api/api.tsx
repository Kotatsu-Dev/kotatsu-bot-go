import { createContext, useContext, type ReactNode } from "react";
import axios, { AxiosError } from "axios";
import { createUsersApi } from "./users";
import { createActivitiesApi } from "./activities";
import { createCalendarApi } from "./calendar";
import { createBroadcastApi } from "./broadcast";
import { createRequestsApi } from "./requests";
import { createDbApi } from "./db";
import { createRoulettesApi } from "./roulettes";
import { toaster } from "../components/ui/toaster";
import z, { ZodError } from "zod";

const ErrorData = z.object({
  status: z.object({
    code: z.int(),
    message: z.string(),
  }),
});

const createApi = (_ctx: null) => {
  const base = import.meta.env.PROD
    ? new URL("/", location.toString()).toString().slice(0, -1)
    : `http://localhost:8006`;
  const $ = axios.create({
    baseURL: `${base}/api/`,
  });

  return {
    users: createUsersApi($),
    activities: createActivitiesApi($),
    calendar: createCalendarApi($, base),
    broadcast: createBroadcastApi($, base),
    requests: createRequestsApi($),
    roulettes: createRoulettesApi($),
    db: createDbApi($, base),
  };
};

const APIContext = createContext(null);

export const APIProvider = (props: { children: ReactNode[] | ReactNode }) => {
  return (
    <APIContext.Provider value={null}>{props.children}</APIContext.Provider>
  );
};

export const useAPI = () => {
  return createApi(useContext(APIContext));
};

export const handleError = (e: unknown) => {
  if (e instanceof ZodError) {
    console.log(`Error parsing data:\n${z.prettifyError(e)}`);
    console.log(e.issues);
    toaster.error({
      description: `Error parsing data:\n${z.prettifyError(e)}`,
    });
    return;
  }

  if (e instanceof AxiosError) {
    if (e.response && e.response.data) {
      const data = ErrorData.safeParse(e.response.data);
      if (data.success) {
        toaster.error({
          description: `Error: ${data.data.status.message}`,
        });
        return;
      }
    }
    toaster.error({
      description: `HTTP Error ${e.code ?? ""}`,
    });
    return;
  }

  toaster.error({
    description: `Unknown error`,
  });
};
