import { ReactNode } from "react";

export function DashboardPageContent({ children }: { children: ReactNode }) {
    return <div className="flex flex-col w-full md:max-w-125 max-md:mx-auto md:mt-10">
        {children}
    </div>
}