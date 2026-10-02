import { ReactNode } from "react";

type Props = {
    children: ReactNode,
}

export function DashboardPageLayout({ children }: Props) {
    return <div className="relative flex flex-col w-full h-full mt-12 px-5 mb-20 gap-5">
        {children}
    </div>
}