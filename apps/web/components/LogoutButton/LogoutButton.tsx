"use client"

import { Button } from "@/components/ui/button"
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { logout } from "@/lib/auth"
import { useMutation } from "@tanstack/react-query"
import { useRouter } from "next/navigation"

export function LogoutButton() {
    const router = useRouter()

    const signOut = useMutation({
        mutationFn: logout,
        onSuccess: () => {
            router.replace("/")
        }
    })
    return (
        <AlertDialog>
            <AlertDialogTrigger
                render={<Button
                    variant="outline"

                >
                    Log Out
                </Button>}
            />
            <AlertDialogContent size="sm">
                <AlertDialogHeader>
                    <AlertDialogTitle>Are you sure you want to log out?</AlertDialogTitle>
                    <AlertDialogDescription>
                        Any unsave data will be lost, you will need to log in again to access your dashboard.
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                        disabled={signOut.isPending}
                        onClick={() => signOut.mutate()}
                    >{signOut.isPending ? "Logging out…" : "Log out"}</AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    )
}
