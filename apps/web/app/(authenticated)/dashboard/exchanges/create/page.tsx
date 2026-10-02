"use client"

import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent";
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader";
import { ExchangeForm } from "@/components/ExchangeForm/ExchangeForm";
import { Button } from "@/components/ui/button";
import { createExchange } from "@/lib/exchanges";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function CreateExchangePage() {
    const router = useRouter()
    const queryClient = useQueryClient()

    const create = useMutation({
        mutationKey: ["exchanges"],
        mutationFn: createExchange,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["exchanges"] })
            router.replace("/dashboard/exchanges")
        }
    })

    return <>
        <DashboardPageHeader title="Create Exchange" description="Set up a new exchange and invite the people you want in it.">
            <Button variant="outline" render={<Link href="/dashboard/exchanges" />} nativeButton={false}>
                Cancel
            </Button>
        </DashboardPageHeader>
        <DashboardPageContent>
            <ExchangeForm onSubmit={create.mutate} isLoading={create.isPending} />
        </DashboardPageContent>
    </>
}
