import type {
  MagicLinkRequest,
  User,
  VerifyLoginCodeRequest,
} from "@giftexchanger/types";
import { apiErrorMessage } from "./utils";

export type { User };

export const SESSION_COOKIE_NAME = "giftexchanger_session";

export async function requestMagicLink(email: string): Promise<void> {
  const body: MagicLinkRequest = { email };
  const response = await fetch("/api/auth/magic-link", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    throw new Error(
      await apiErrorMessage(response, "We could not send your login link."),
    );
  }
}

export async function verifyLoginCode(
  email: string,
  code: string,
): Promise<void> {
  const body: VerifyLoginCodeRequest = { email, code };
  const response = await fetch("/api/auth/verify-code", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    throw new Error(
      await apiErrorMessage(response, "We could not log you in."),
    );
  }
}

export async function logout(): Promise<void> {
  const response = await fetch("/api/auth/logout", { method: "POST" });

  if (!response.ok) {
    throw new Error(
      await apiErrorMessage(response, "We could not log you out."),
    );
  }
}
