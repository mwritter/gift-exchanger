import { ReactNode } from "react";

type Props = {
    children: ReactNode,
}

export function DashboardPageLayout({ children }: Props) {

    return <div className=" relative flex flex-col w-full h-full p-10 pt-12 mb-20 gap-5">
        {children}
    </div>
}