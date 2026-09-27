import { DashboardPageLayout } from "@/components/DashboardPageLayout/DashboardPageLayout"
import Link from "next/link"

export default function DashboardNotFound() {
    return (
        <DashboardPageLayout>
            <div className="flex flex-col items-center gap-4">
                <h2 className="text-2xl font-bold">Page not found</h2>
                <p className="text-sm text-muted-foreground">
                    That dashboard page does not exist.
                </p>
                <Link href="/dashboard" className="text-sm font-medium underline">
                    Back to dashboard
                </Link>
            </div>
        </DashboardPageLayout>
    )
}
