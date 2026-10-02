import { User } from "@giftexchanger/types";

export type StatusResponse = {
  status: string;
};

export async function fetchApiStatus(
  path: "/api/health" | "/api/ready",
): Promise<StatusResponse> {
  const response = await fetch(path);
  const data = (await response.json()) as StatusResponse;

  if (!response.ok) {
    throw new Error(data.status ?? response.statusText);
  }

  return data;
}

export async function me(): Promise<User | null> {
  const response = await fetch("/api/me");

  if (!response.ok) {
    if (response.status === 401) {
      return null;
    }
    throw new Error("could not fetch user data");
  }

  return await response.json();
}
