"use client"

import { cn } from "cn"
import Link, { LinkProps } from "next/link"
import { usePathname } from "next/navigation"
import { ReactNode } from "react"

export function DashboardMobileNavLink({ children, ...props }: LinkProps & { children: ReactNode }) {
    const pathname = usePathname()
    return <Link {...props} className={cn("flex-1 p-3 flex flex-col justify-center items-center hover:bg-primary/10 active:bg-secondary/25", {
        "text-primary bg-primary/10": pathname === props.href
    })}>
        {children}
    </Link>
}