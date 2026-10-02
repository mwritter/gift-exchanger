"use client"

import { useCurrentUser } from "@/components/CurrentUser/CurrentUserProvider"
import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent"
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader"
import { DeleteExchangeButton } from "@/components/DeleteExchangeButton/DeleteExchangeButton"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Spinner } from "@/components/ui/spinner"
import { getExchangeById, getExchangeInvites } from "@/lib/exchanges"
import { useQuery } from "@tanstack/react-query"
import { parseISO } from "date-fns"
import { ArrowLeft, Calendar, Edit, Mail, User2, Wallet } from "lucide-react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { use, type ReactNode } from "react"

export default function ExchangeDetailsPage({
    params,
}: {
    params: Promise<{ exchangeId: string }>
}) {
    const { exchangeId } = use(params)
    const router = useRouter()
    const { user } = useCurrentUser()

    const { data, error } = useQuery({
        queryKey: ["exchange", exchangeId],
        queryFn: () => getExchangeById(exchangeId)
    })

    const isOrganizer = data?.organizerId === user?.id

    // Only the organizer is allowed to read invites, so don't ask otherwise.
    const { data: invites } = useQuery({
        queryKey: ["exchange", exchangeId, "invites"],
        queryFn: () => getExchangeInvites(exchangeId),
        enabled: isOrganizer
    })

    if (error) {
        return <p className="text-sm text-destructive">{error.message}</p>
    }

    if (!data) {
        return <Spinner className="self-center" />
    }

    return <>
        <DashboardPageHeader title={data.name} description={data.description}>
            <div className="flex gap-2">
                <Button nativeButton={false} render={<Link href="/dashboard/exchanges" />} variant="outline">
                    <ArrowLeft size={12} />
                    Back
                </Button>
                {isOrganizer && <>
                    <Button nativeButton={false} render={<Link href={`/dashboard/exchanges/${exchangeId}/edit`} />} variant="outline">
                        <Edit size={12} />
                        Edit
                    </Button>
                    <DeleteExchangeButton
                        exchangeId={exchangeId}
                        exchangeName={data.name}
                        onDeleted={() => router.replace("/dashboard/exchanges")}
                    />
                </>}
            </div>
        </DashboardPageHeader>
        <DashboardPageContent>
            <dl className="flex flex-col gap-4">
                <DetailRow icon={<Calendar size={15} />} label="Exchange day">
                    {parseISO(data.exchangeDate).toDateString()}
                </DetailRow>
                <DetailRow icon={<Wallet size={15} />} label="Budget">
                    {data.budgetCents == null ? "No budget set" : `$${data.budgetCents / 100}`}
                </DetailRow>
                <DetailRow icon={<User2 size={15} />} label="Members">
                    <ul className="flex flex-col gap-2">
                        {data.members.map(member => <li key={member.userId} className="flex items-center justify-between gap-2">
                            <span className="truncate">
                                {member.displayName || member.email}
                            </span>
                            {member.isOrganizer && <Badge variant="secondary">Organizer</Badge>}
                        </li>)}
                    </ul>
                </DetailRow>
                {isOrganizer && invites && <DetailRow icon={<Mail size={15} />} label="Pending invites">
                    {invites.inviteEmails.length > 0
                        ? <ul className="flex flex-col">
                            {invites.inviteEmails.map(email => <li key={email}>{email}</li>)}
                        </ul>
                        : "Nobody invited yet"}
                </DetailRow>}
            </dl>
        </DashboardPageContent>
    </>
}

function DetailRow({ icon, label, children }: { icon: ReactNode, label: string, children: ReactNode }) {
    return <div className="flex flex-col gap-1">
        <dt className="flex items-center gap-2 text-xs text-muted-foreground">
            {icon} {label}
        </dt>
        <dd className="text-sm">{children}</dd>
    </div>
}
