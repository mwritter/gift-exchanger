import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent"
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader"
import { DashboardPageLayout } from "@/components/DashboardPageLayout/DashboardPageLayout"
import { LogoutButton } from "@/components/LogoutButton/LogoutButton"
import { requireUser } from "@/lib/session"

export default async function Dashboard() {
    const user = await requireUser()

    return <DashboardPageLayout>
        <DashboardPageHeader title="My Dashboard" />
        <DashboardPageContent>
            <p>This is your dashboard</p>
            <p>{user.display_name || user.email}</p>
        </DashboardPageContent>
    </DashboardPageLayout>
}
