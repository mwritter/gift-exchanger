"use client"

import { cn } from "cn";
import Link, { LinkProps } from "next/link";
import { SidebarMenuButton } from "../ui/sidebar";
import { ReactNode } from "react";
import { usePathname } from "next/navigation";

export function DashboardSidebarLink({ children, ...props }: LinkProps & { children: ReactNode }) {
    const pathname = usePathname()
    return <SidebarMenuButton className={cn("[&_svg]:size-6 font-semibold text-lg py-5 hover:bg-primary/10 active:bg-secondary/25 hover:text-primary active:text-primary w-full", {
        "text-primary bg-primary/10": pathname === props.href
    })} render={
        <Link {...props} />
    }>
        {children}
    </SidebarMenuButton>
}

