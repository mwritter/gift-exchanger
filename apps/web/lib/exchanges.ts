import type { ExchangeValues } from "@/components/ExchangeForm/ExchangeForm";
import type {
  CreateExchangeRequest,
  Exchange,
  ListExchangeInvitesResponse,
  ListExchangesResponse,
} from "@giftexchanger/types";
import { format, isValid, parseISO } from "date-fns";
import { apiErrorMessage } from "./utils";

// The API sends and expects calendar dates as YYYY-MM-DD. parseISO reads a
// date-only string as local midnight, so the day doesn't shift by timezone.
export function fromApiDate(value: string): Date | null {
  const date = parseISO(value);
  return isValid(date) ? date : null;
}

function toExchangeRequest(values: ExchangeValues): CreateExchangeRequest {
  if (!values.exchangeDate) {
    throw new Error("Pick an exchange date");
  }

  return {
    name: values.name,
    description: values.description ?? "",
    exchangeDate: format(values.exchangeDate, "yyyy-MM-dd"),
    budgetCents:
      values.budget == null ? null : Math.round(values.budget * 100),
    inviteEmails: values.inviteEmails,
  };
}

export async function getExchanges() {
  const response = await fetch("/api/exchanges");

  if (!response.ok) {
    throw new Error(
      await apiErrorMessage(response, "We could not get your exchanges."),
    );
  }

  return (await response.json()) as ListExchangesResponse;
}

export async function getExchangeById(exchangeId: string) {
  const response = await fetch(`/api/exchanges/${exchangeId}`);

  if (!response.ok) {
    throw new Error(
      await apiErrorMessage(response, "We could not get your exchange."),
    );
  }

  return (await response.json()) as Exchange;
}

export async function getExchangeInvites(exchangeId: string) {
  const response = await fetch(`/api/exchanges/${exchangeId}/invites`);

  if (!response.ok) {
    throw new Error(
      await apiErrorMessage(response, "We could not get your invites."),
    );
  }

  return (await response.json()) as ListExchangeInvitesResponse;
}

export async function updateExchange(
  values: ExchangeValues,
  exchangeId: string,
) {
  const response = await fetch(`/api/exchanges/${exchangeId}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(toExchangeRequest(values)),
  });

  if (!response.ok) {
    throw new Error(
      await apiErrorMessage(response, "We could not update your exchanges."),
    );
  }

  return await response.json();
}

export async function createExchange(values: ExchangeValues) {
  const response = await fetch("/api/exchanges", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(toExchangeRequest(values)),
  });

  if (!response.ok) {
    throw new Error(
      await apiErrorMessage(response, "We could not create your exchanges."),
    );
  }

  return await response.json();
}

export async function deleteExchange(exchangeId: string) {
  const response = await fetch(`/api/exchanges/${exchangeId}`, {
    method: "DELETE",
  });

  if (!response.ok) {
    throw new Error(
      await apiErrorMessage(response, "We could not delete exchange."),
    );
  }

  return await response.json();
}
