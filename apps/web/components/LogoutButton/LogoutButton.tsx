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
import { useCurrentUser } from "../CurrentUser/CurrentUserProvider"

export function LogoutButton() {
    const { logout } = useCurrentUser()
    return (
        <AlertDialog>
            <AlertDialogTrigger
                render={<Button variant="outline">
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
                        disabled={logout.isPending}
                        onClick={() => logout.mutate()}
                    >{logout.isPending ? "Logging out…" : "Log out"}</AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    )
}
