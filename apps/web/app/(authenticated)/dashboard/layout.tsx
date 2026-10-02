import { DashboardMobileHeader } from "@/components/DashboardMobileHeader/DashboardMobileHeader";
import { DashboardMobileNav } from "@/components/DashboardMobileNav/DashboardMobileNav";
import { DashboardPageLayout } from "@/components/DashboardPageLayout/DashboardPageLayout";
import { DashboardSidebar } from "@/components/DashboardSidebar/DashboardSidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { requireUser } from "@/lib/session";
import { ReactNode } from "react";

export default async function Layout({ children }: { children: ReactNode }) {
    await requireUser()

    return (
        <>
            <SidebarProvider>
                <DashboardSidebar />
                <SidebarInset>
                    <DashboardMobileHeader />
                    <DashboardPageLayout>{children}</DashboardPageLayout>
                </SidebarInset>
            </SidebarProvider>
            <DashboardMobileNav />
        </>
    )
}