import { Gift } from "lucide-react";
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarGroupContent, SidebarHeader, SidebarMenu, SidebarMenuItem } from "../ui/sidebar";
import { requireUser } from "@/lib/session";
import { DashboardSidebarLink } from "./DashboardSidebarLink";
import { LogoutButton } from "../LogoutButton/LogoutButton";
import navLinks from "../consts/navLinks";

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
                        {navLinks.map(({ href, Icon, text }) => <SidebarMenuItem key={href}>
                            <DashboardSidebarLink href={href}>
                                <Icon />
                                {text}
                            </DashboardSidebarLink>
                        </SidebarMenuItem>)}
                    </SidebarMenu>
                </SidebarGroupContent>
            </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
            <div>
                <p>{user.displayName}</p>
                <p className="text-xs">{user.email}</p>
            </div>
            <LogoutButton />
        </SidebarFooter>
    </Sidebar>
}