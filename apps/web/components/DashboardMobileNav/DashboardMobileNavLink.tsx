"use client"

import { useActivePathname } from "@/hooks/use-active-pathname"
import { cn } from "cn"
import Link from "next/link"
import { ComponentProps } from "react"

type Props = Omit<ComponentProps<typeof Link>, "href"> & { href: string }

export function DashboardMobileNavLink({ children, className, ...props }: Props) {
    const { isActive } = useActivePathname({ href: props.href })
    return <Link {...props} className={cn("flex-1 p-3 flex flex-col justify-center items-center hover:bg-primary/10 active:bg-secondary/25", {
        "text-primary bg-primary/10": isActive
    }, className)}>
        {children}
    </Link>
}
