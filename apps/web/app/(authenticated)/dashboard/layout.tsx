import { DashboardMobileHeader } from "@/components/DashboardMobileHeader/DashboardMobileHeader";
import { DashboardMobileNav } from "@/components/DashboardMobileNav/DashboardMobileNav";
import { DashboardSidebar } from "@/components/DashboardSidebar/DashboardSidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { ReactNode } from "react";

export default function Layout({ children }: { children: ReactNode }) {
    // Show a mobile bottom nav on small screens
    return (
        <div>
            <SidebarProvider>
                <DashboardSidebar />
                <SidebarInset>
                    <DashboardMobileHeader />
                    {children}
                </SidebarInset>
            </SidebarProvider>
            <DashboardMobileNav />
        </div>
    )
}