import { ReactNode } from "react"

type Props = {
    title: string,
    description?: string,
    children?: ReactNode
}

export function DashboardPageHeader({ title, description, children }: Props) {
    return <div className="flex flex-col md:flex-row gap-5 justify-between md:items-center max-w-125">
        <div>
            <h1 className="text-xl font-bold">{title}</h1>
            {description && <p className="text-xs text-muted-foreground">{description}</p>}
        </div>
        {children}
    </div>
}