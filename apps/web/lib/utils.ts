import { ErrorResponse } from "@giftexchanger/types";

export { cn } from "cn";

export async function apiErrorMessage(response: Response, fallback: string) {
  const data = (await response
    .json()
    .catch(() => null)) as Partial<ErrorResponse> | null;

  return data?.error ?? fallback;
}
