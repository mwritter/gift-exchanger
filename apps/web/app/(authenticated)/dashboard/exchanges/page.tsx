"use client"

import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent";
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader";
import { ExchangeList } from "@/components/ExchangeList/ExchangeList";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { getExchanges } from "@/lib/exchanges";
import { useQuery } from "@tanstack/react-query";
import { Plus } from "lucide-react";
import Link from "next/link";

export default function ExchangesPage() {

    const { data, error } = useQuery({
        queryKey: ['exchanges'],
        queryFn: getExchanges
    })

    return <>
        <DashboardPageHeader title="Exchanges" description="This is a list of exchanges your current a part of!">
            <Button render={<Link href="/dashboard/exchanges/create" />} nativeButton={false}>
                <Plus />
                Create
            </Button>
        </DashboardPageHeader>
        <DashboardPageContent>
            {error
                ? <p className="text-sm text-destructive">{error.message}</p>
                : !data
                ? <Spinner className="self-center" />
                : data.exchanges.length > 0
                    ? <ExchangeList exchanges={data.exchanges} />
                    : <EmptyExchangeList />}
        </DashboardPageContent>
    </>
}

function EmptyExchangeList() {
    return <div className="w-full text-center justify-center">
        <p className="text-sm">No exchanges to list yet</p>
    </div>
}
