"use client"

import { useCurrentUser } from "@/components/CurrentUser/CurrentUserProvider"
import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent"
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader"
import { ExchangeForm, ExchangeValues } from "@/components/ExchangeForm/ExchangeForm"
import { Button } from "@/components/ui/button"
import { Spinner } from "@/components/ui/spinner"
import { fromApiDate, getExchangeById, getExchangeInvites, updateExchange } from "@/lib/exchanges"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { useRouter } from "next/navigation"
import { use, useEffect } from "react"

export default function EditExchange({
    params,
}: {
    params: Promise<{ exchangeId: string }>
}) {
    const router = useRouter()
    const queryClient = useQueryClient()
    const { user, isLoading: isUserLoading } = useCurrentUser()
    const { exchangeId } = use(params)

    // Fetch the exchange to fill in the form
    const { data, error: exchangeError } = useQuery({
        queryKey: ["exchange", exchangeId],
        queryFn: () => getExchangeById(exchangeId)
    })

    const isOrganizer = data?.organizerId === user?.id

    // Only the organizer is allowed to read invites, so don't ask otherwise.
    const { data: invites, error: invitesError } = useQuery({
        queryKey: ["exchange", exchangeId, "invites"],
        queryFn: () => getExchangeInvites(exchangeId),
        enabled: isOrganizer
    })

    // Only the organizer can edit, so send everyone else back to the details page.
    useEffect(() => {
        if (isUserLoading) return
        if (data && !isOrganizer) {
            router.replace(`/dashboard/exchanges/${exchangeId}`)
        }
    }, [data, isOrganizer, exchangeId, router, isUserLoading])

    const submit = useMutation({
        mutationKey: ["exchanges", exchangeId],
        mutationFn: (exchange: ExchangeValues) => updateExchange(exchange, exchangeId),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["exchanges"] })
            queryClient.invalidateQueries({ queryKey: ["exchange", exchangeId] })
            router.replace(`/dashboard/exchanges/${exchangeId}`)
        }
    })

    const handleSubmit = (exchange: ExchangeValues) => {
        submit.mutate(exchange)
    }

    const loadError = exchangeError ?? invitesError
    if (loadError) {
        return <p className="text-sm text-destructive">{loadError.message}</p>
    }

    // The form seeds its defaults once, so wait for the invites before mounting it.
    if (!data || !invites) {
        return <Spinner className="self-center" />
    }

    return <>
        <DashboardPageHeader title="Edit Exchange" description={data.name}>
            <Button variant="outline" onClick={() => router.back()}>
                Cancel
            </Button>
        </DashboardPageHeader>
        <DashboardPageContent>
            <ExchangeForm
                onSubmit={handleSubmit}
                name={data.name}
                description={data.description}
                exchangeDate={fromApiDate(data.exchangeDate)}
                budget={data.budgetCents == null ? null : data.budgetCents / 100}
                inviteEmails={invites.inviteEmails}
                isLoading={submit.isPending}
            />
        </DashboardPageContent>
    </>
}
