import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { cache } from "react";

import { SESSION_COOKIE_NAME, type User } from "@/lib/auth";

const API_URL = process.env.API_URL ?? "http://localhost:8080";

// Server components talk to the Go API directly, so the session cookie has to be
// forwarded by hand. cache() keeps this to one /api/me call per render pass.
export const getCurrentUser = cache(async (): Promise<User | null> => {
  const session = (await cookies()).get(SESSION_COOKIE_NAME)?.value;

  if (!session) {
    return null;
  }

  const response = await fetch(`${API_URL}/api/me`, {
    headers: { cookie: `${SESSION_COOKIE_NAME}=${session}` },
    cache: "no-store",
  }).catch(() => null);

  if (!response?.ok) {
    return null;
  }

  return (await response.json()) as User;
});

export async function requireUser(): Promise<User> {
  const user = await getCurrentUser();

  if (!user) {
    redirect("/login");
  }

  return user;
}
