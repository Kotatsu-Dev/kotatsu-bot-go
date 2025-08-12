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

const createApi = (_ctx: null) => {
  const base = `http://localhost:8006`;
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
    toaster.error({
      description: `Error parsing data:\n${z.prettifyError(e)}`,
    });
  }

  if (e instanceof AxiosError) {
    toaster.error({
      description: `HTTP Error`,
    });
  }
};
