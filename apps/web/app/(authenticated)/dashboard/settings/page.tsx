import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent";
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader";
import { DashboardPageLayout } from "@/components/DashboardPageLayout/DashboardPageLayout";

export default function SettingsPage() {
    return <DashboardPageLayout>
        <DashboardPageHeader title="Settings" description="GiftExchanger Settings" />
        <DashboardPageContent>
            Setting Page
        </DashboardPageContent>
    </DashboardPageLayout>
}