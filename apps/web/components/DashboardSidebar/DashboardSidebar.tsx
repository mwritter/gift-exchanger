import { Gift, Heart, Home, Settings } from "lucide-react";
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarGroupContent, SidebarHeader, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "../ui/sidebar";
import { requireUser } from "@/lib/session";
import { DashboardSidebarLink } from "./DashboardSidebarLink";
import { LogoutButton } from "../LogoutButton/LogoutButton";

export async function DashboardSidebar() {
    const user = await requireUser()
    return <Sidebar variant="inset">
        <SidebarHeader>
            <div className="flex gap-2 justify-center items-center py-5">
                <Gift size={40} className="text-primary" />
                <p className="font-bold text-xl">GiftExchanger</p>
            </div>
        </SidebarHeader>
        <SidebarContent>
            <SidebarGroup>
                <SidebarGroupContent>
                    <SidebarMenu className="flex flex-col gap-2">
                        <SidebarMenuItem>
                            <DashboardSidebarLink href={"/dashboard"}>
                                <Home />
                                Home
                            </DashboardSidebarLink>
                        </SidebarMenuItem>
                        <SidebarMenuItem>
                            <DashboardSidebarLink href={"/dashboard/exchanges"}>
                                <Gift />
                                Exchanges
                            </DashboardSidebarLink>
                        </SidebarMenuItem>
                        <SidebarMenuItem>
                            <DashboardSidebarLink href={"/dashboard/wishlist"}>
                                <Heart />
                                Wishlist
                            </DashboardSidebarLink>
                        </SidebarMenuItem>
                        <SidebarMenuItem>
                            <DashboardSidebarLink href={"/dashboard/settings"}>
                                <Settings />
                                Settings
                            </DashboardSidebarLink>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarGroupContent>
            </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
            <div>
                <p>{user.display_name}</p>
                <p className="text-xs">{user.email}</p>
            </div>
            <LogoutButton />
        </SidebarFooter>
    </Sidebar>
}