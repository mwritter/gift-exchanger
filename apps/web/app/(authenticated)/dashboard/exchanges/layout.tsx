"use client"

import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent"
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader"
import { DashboardPageLayout } from "@/components/DashboardPageLayout/DashboardPageLayout"
import { Button } from "@/components/ui/button"
import { Plus } from "lucide-react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { ReactNode } from "react"

type Props = { children: ReactNode }

export default function Layout({ children }: Props) {
    const pathname = usePathname()

    const isRoot = pathname.endsWith('exchanges')

    return <DashboardPageLayout>
        <div className="flex gap-2 items-center justify-between max-w-125">
            <DashboardPageHeader title="Exchanges" description="This is a list of exchanges your current a part of!" />
            {isRoot ? <Button render={<Link href="/dashboard/exchanges/create" />} nativeButton={false}>
                <Plus />
                Create
            </Button> : <Button variant='destructive' render={<Link href={"/dashboard/exchanges"} />} nativeButton={false}>
                Cancel
            </Button>}
        </div>
        <DashboardPageContent>
            {children}
        </DashboardPageContent>
    </DashboardPageLayout>
}


