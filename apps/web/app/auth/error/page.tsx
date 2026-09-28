import { buttonVariants } from "@/components/ui/button"
import { Gift } from "lucide-react"
import Link from "next/link"

const reasons = {
    invalid_link: {
        title: "That log in link no longer works",
        description:
            "Log in links can only be used once, and they expire a few minutes after we send them. Request a new one and it will be waiting in your inbox."
    },
    unexpected: {
        title: "We could not finish logging you in",
        description:
            "Something went wrong on our end. Requesting a new login link usually clears it up."
    }
} as const

function contentFor(reason: string | string[] | undefined) {
    if (typeof reason === "string" && reason in reasons) {
        return reasons[reason as keyof typeof reasons]
    }
    return reasons.unexpected
}

export default async function AuthError({
    searchParams
}: PageProps<"/auth/error">) {
    const { reason } = await searchParams
    const { title, description } = contentFor(reason)

    return <div className="flex flex-1 p-10">
        <Link href={"/"} className="fixed flex items-center gap-2 top-10 left-10 font-bold text-2xl max-sm:hidden">
            <Gift /> GiftExchanger
        </Link>
        <div className="flex flex-col justify-center items-center flex-1">
            <div className="flex flex-col items-center gap-6 w-full max-w-120 text-center">
                <Link href={"/"} className="flex gap-2 items-center font-bold text-xl sm:hidden">
                    <Gift /> GiftExchanger
                </Link>
                <div className="space-y-2">
                    <h1 className="text-2xl font-bold tracking-tight">{title}</h1>
                    <p className="text-sm text-muted-foreground">{description}</p>
                </div>
                <Link href={"/login"} className={buttonVariants({ size: "lg" })}>
                    Request a new link
                </Link>
            </div>
        </div>
    </div>
}
