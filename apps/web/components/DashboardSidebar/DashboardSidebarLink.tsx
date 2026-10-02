"use client"

import { cn } from "cn";
import Link from "next/link";
import { SidebarMenuButton } from "../ui/sidebar";
import { ComponentProps } from "react";
import { useActivePathname } from "@/hooks/use-active-pathname";

type Props = Omit<ComponentProps<typeof Link>, "href"> & { href: string }

export function DashboardSidebarLink({ children, ...props }: Props) {
    const { isActive } = useActivePathname({ href: props.href })
    return <SidebarMenuButton className={cn("[&_svg]:size-6 font-semibold text-lg py-5 hover:bg-primary/10 active:bg-secondary/25 hover:text-primary active:text-primary w-full", {
        "text-primary bg-primary/10": isActive
    })} render={
        <Link {...props} />
    }>
        {children}
    </SidebarMenuButton>
}
