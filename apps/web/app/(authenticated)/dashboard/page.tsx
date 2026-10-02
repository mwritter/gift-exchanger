import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent"
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader"
import { requireUser } from "@/lib/session"

export default async function Dashboard() {
    const user = await requireUser()

    return (
        <>
            <DashboardPageHeader title="My Dashboard" />
            <DashboardPageContent>
                <p>This is your dashboard</p>
                <p>{user.displayName || user.email}</p>
            </DashboardPageContent>
        </>
    )
}
